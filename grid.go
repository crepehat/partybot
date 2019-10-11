package partybot

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func NewGrid(gridSlice [][]string) (g Grid, err error) {
	g.array = gridSlice
	g.seqCtx, g.seqCancel = context.WithCancel(context.Background())
	g.seqLock = &sync.Mutex{}
	for y, line := range gridSlice {
		for x, blockName := range line {
			b := NewBlock(blockName, x, y)
			b.Start()
			// add the block
			g.blocks = append(g.blocks, b)
			if b.X == 0 {
				g.blockArray = append(g.blockArray, []*Block{b})
			} else {
				g.blockArray[b.Y] = append(g.blockArray[b.Y], b)
			}
			if b.X > g.x {
				g.x = b.X
			}
			if b.Y > g.y {
				g.y = b.Y
			}
		}
	}
	return
}

func SleepCanBreak(ctx context.Context, sleep float64) (isBreak bool) {
	select {
	case <-ctx.Done():
		isBreak = true
	case <-time.After(time.Duration(sleep*1000) * time.Millisecond):
		isBreak = false
	}
	return
}

func (g *Grid) PrintBlock(x, y int) {
	fmt.Println(g.blockArray[y][x].Name)
}

func (g *Grid) GetMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/grid", g.handleArray())
	mux.HandleFunc("/sequence/start", g.handleSequenceStart())
	mux.HandleFunc("/sequence/stop", g.handleSequenceStop())
	mux.HandleFunc("/snake", g.handleSnake())
	for _, butt := range g.blocks {
		mux.HandleFunc(fmt.Sprintf("/block/%s", butt.Name), butt.ServeWs())
	}
	return mux
}

func (g *Grid) handleSequenceStart() http.HandlerFunc {
	type startCommand struct {
		Name         string  `json:"name"`
		CycleSeconds float64 `json:"cycle_seconds"`
	}
	type startResponse struct {
		Response string `json:"response"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var command startCommand
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&command)
		if err != nil {
			fmt.Println(err)
			return
		}
		g.seqLock.Lock()
		g.seqCancel()
		g.TurnAllOff()
		// time.Sleep(1 * time.Second)
		g.seqCtx, g.seqCancel = context.WithCancel(context.Background())
		g.seqLock.Unlock()
		switch command.Name {
		case "wave":
			g.Wave(command.CycleSeconds)
		case "mexican_wave":
			g.MexicanWave(command.CycleSeconds)
		case "alt_mexican_wave":
			g.AlternatingMexicanWave(command.CycleSeconds)
		case "alt_wave":
			g.AlternatingWave(command.CycleSeconds)
		case "random_snake":
			g.RandomSnake(command.CycleSeconds)
		case "snake":
			g.Snake()
		}
		res := startResponse{
			Response: fmt.Sprintf("Started %s", command.Name),
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (g *Grid) handleSequenceStop() http.HandlerFunc {
	type stopResponse struct {
		Response string `json:"response"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		g.seqLock.Lock()
		g.seqCancel()
		g.TurnAllOff()
		g.seqLock.Unlock()
		res := stopResponse{
			Response: "Stopped sequence",
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (g *Grid) handleArray() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(g.array)
	}
}

func (g *Grid) TurnAllOff() {
	for _, b := range g.blocks {
		b.LightOff()
	}
}

func (g *Grid) FadeAll() {
	for _, b := range g.blocks {
		b.LightFadeOut(g.seqCtx, 1.0, 0)
	}
}
func (g *Grid) RandomSnake(cycleTime float64) {
	var x, y int
	x = 0
	y = 0
	var step, axis int
	fmt.Println(g.x, g.y)
	go func(ctx context.Context) {
		for {
			step = 1 - 2*rand.Intn(2)
			axis = rand.Intn(2)
			if axis == 0 {
				x = (x + step) % (g.x + 1)
				if x < 0 {
					x = g.x
				}
			} else {
				y = (y + step) % (g.y + 1)
				if y < 0 {
					y = g.y
				}
			}
			g.blockArray[y][x].LightOn()
			if SleepCanBreak(ctx, cycleTime) {
				return
			}
			g.blockArray[y][x].LightOff()
		}
	}(g.seqCtx)
}

func (g *Grid) handleSnake() http.HandlerFunc {
	type snakeCommand struct {
		Direction string `json:"direction"`
	}
	type snakeResponse struct {
		Response string `json:"response"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var command snakeCommand
		w.Header().Set("Access-Control-Allow-Origin", "*")
		err := json.NewDecoder(r.Body).Decode(&command)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch command.Direction {
		case "ArrowDown":
			g.snakeDirection = 0
		case "ArrowRight":
			g.snakeDirection = 1
		case "ArrowUp":
			g.snakeDirection = 2
		case "ArrowLeft":
			g.snakeDirection = 3
		}
		resp := snakeResponse{Response: command.Direction}
		json.NewEncoder(w).Encode(resp)
	}
}

func (g *Grid) GetOffBlocks() (emptyBlocks []*Block) {
	for _, block := range g.blocks {
		if !block.state.LightOn {
			emptyBlocks = append(emptyBlocks, block)
		}
	}
	return
}

func (g *Grid) GetRandomOffBlock() *Block {
	offBlocks := g.GetOffBlocks()
	if len(offBlocks) == 0 {
		fmt.Println("No off blocks")
		return nil
	}
	position := rand.Int() % len(offBlocks)
	return offBlocks[position]
}

func (g *Grid) Snake() {
	go func(ctx context.Context) {
		score := 1
		target := g.GetRandomOffBlock()
		target.LightOn()
		for {
			fmt.Println(g.snakeDirection)
			if g.snakeDirection == 0 {
				g.snakeCoord.y = (g.snakeCoord.y + 1) % (g.y + 1)
			} else if g.snakeDirection == 1 {
				g.snakeCoord.x = (g.snakeCoord.x + 1) % (g.x + 1)
			} else if g.snakeDirection == 2 {
				g.snakeCoord.y = (g.snakeCoord.y - 1) % (g.y + 1)
			} else {
				g.snakeCoord.x = (g.snakeCoord.x - 1) % (g.x + 1)
			}
			if g.snakeCoord.y < 0 {
				g.snakeCoord.y = g.y
			}
			if g.snakeCoord.x < 0 {
				g.snakeCoord.x = g.x
			}
			// check if hit target or crashed
			if target.X == g.snakeCoord.x && target.Y == g.snakeCoord.y {
				score++
				target = g.GetRandomOffBlock()
				if target != nil {
					target.LightOn()
				} else {
					g.FadeAll()
					return
				}
			} else if g.blockArray[g.snakeCoord.y][g.snakeCoord.x].state.LightOn {
				g.TurnAllOff()
				return
			}
			g.blockArray[g.snakeCoord.y][g.snakeCoord.x].LightOn()
			// The bit that controls the light staying on
			go func(co coOrd) {
				SleepCanBreak(ctx, 1.0)
				g.blockArray[co.y][co.x].LightOff()

			}(g.snakeCoord)
			// The bit that waits to move the head forward
			if SleepCanBreak(ctx, 1.0/float64(score)) {
				return
			}
		}
	}(g.seqCtx)
}

func (g Grid) MexicanWave(cycleTime float64) {
	for _, block := range g.blocks {
		go func(ctx context.Context, block *Block) {
			block.LightPulse(ctx, cycleTime, float64(block.X)/float64(g.x))
			for {
				select {
				case <-ctx.Done():
					g.TurnAllOff()
					return
				default:
					block.LightPulse(ctx, cycleTime, 0.0)
				}
			}
		}(g.seqCtx, block)
	}
}

func (g Grid) AlternatingMexicanWave(cycleTime float64) {
	for _, block := range g.blocks {
		go func(ctx context.Context, block *Block) {
			if block.Y%2 == 0 {
				block.LightPulse(ctx, cycleTime, float64(block.X)/float64(g.x))
			} else {
				block.LightPulse(ctx, cycleTime, 1.0-float64(block.X)/float64(g.x))
			}
			for {
				select {
				case <-ctx.Done():
					g.TurnAllOff()
					return
				default:
					block.LightPulse(ctx, cycleTime, 0.0)
				}
			}
		}(g.seqCtx, block)
	}
}

func (g Grid) Wave(cycleTime float64) {
	for _, block := range g.blocks {
		go func(ctx context.Context, block *Block) {
			// fmt.Println(float64(block.X) / float64(g.x) * cycleTime)
			time.Sleep(time.Duration(float64(block.X)/float64(g.x)*cycleTime*1000) * time.Millisecond)
			fmt.Println("starting", block.Name)
			for {
				block.LightOn()
				if SleepCanBreak(ctx, cycleTime) {
					return
				}
				block.LightOff()
				if SleepCanBreak(ctx, cycleTime) {
					return
				}
			}
		}(g.seqCtx, block)
	}
}

func (g Grid) AlternatingWave(cycleTime float64) {
	for _, block := range g.blocks {
		go func(ctx context.Context, block *Block) {
			// fmt.Println(float64(block.X) / float64(g.x) * cycleTime)
			if block.Y%2 == 0 {
				time.Sleep(time.Duration(float64(block.X)/float64(g.x)*cycleTime*1000) * time.Millisecond)
			} else {
				time.Sleep(time.Duration((1-float64(block.X)/float64(g.x))*cycleTime*1000) * time.Millisecond)

			}
			fmt.Println("starting", block.Name)
			for {
				select {
				case <-ctx.Done():
					g.TurnAllOff()
					return
				default:
					block.LightOn()
					if SleepCanBreak(ctx, cycleTime) {
						return
					}
					block.LightOff()
					if SleepCanBreak(ctx, cycleTime) {
						return
					}
				}
			}
		}(g.seqCtx, block)
	}
}
