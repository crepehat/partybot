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
	type sequenceCommand struct {
		Name         string  `json:"name"`
		CycleSeconds float64 `json:"cycle_seconds"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var command sequenceCommand
		err := json.NewDecoder(r.Body).Decode(&command)
		if err != nil {
			fmt.Println(err)
			return
		}
		g.seqLock.Lock()
		// json.NewEncoder(w).Encode(g.array)
		g.seqCancel()
		g.seqCtx, g.seqCancel = context.WithCancel(context.Background())
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
			g.Snake(command.CycleSeconds)
		}
	}
}

func (g *Grid) handleSequenceStop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		g.seqLock.Lock()
		// json.NewEncoder(w).Encode(g.array)
		g.seqCancel()
		g.seqLock.Unlock()
	}
}

func (g *Grid) handleArray() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(g.array)
	}
}

func (g *Grid) RandomSnake(cycleTime float64) {
	var x, y int
	x = 0
	y = 0
	var step, axis int
	fmt.Println(g.x, g.y)
	go func() {
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
			select {
			case <-g.seqCtx.Done():
				return
			default:
				g.blockArray[y][x].LightOn()
				time.Sleep(time.Duration(cycleTime*1000) * time.Millisecond)
				g.blockArray[y][x].LightOff()
			}
		}
	}()
}

func (g *Grid) handleSnake() http.HandlerFunc {
	type snakeCommand struct {
		Direction string `json:"direction"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var command snakeCommand
		err := json.NewDecoder(r.Body).Decode(&command)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		switch command.Direction {
		case "up":
			g.snakeDirection = 0
		case "right":
			g.snakeDirection = 1
		case "down":
			g.snakeDirection = 2
		case "left":
			g.snakeDirection = 3
		}
		fmt.Println("direction:", g.snakeDirection)
	}
}

func (g *Grid) Snake(cycleTime float64) {
	go func() {
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
			select {
			case <-g.seqCtx.Done():
				return
			default:
				g.blockArray[g.snakeCoord.y][g.snakeCoord.x].LightOn()
				time.Sleep(time.Duration(cycleTime*1000) * time.Millisecond)
				g.blockArray[g.snakeCoord.y][g.snakeCoord.x].LightOff()
			}
		}
	}()

}

func (g Grid) MexicanWave(cycleTime float64) {
	for _, block := range g.blocks {
		go func(block *Block) {
			block.LightPulse(g.seqCtx, cycleTime, float64(block.X)/float64(g.x))
			for {
				select {
				case <-g.seqCtx.Done():
					return
				default:
					block.LightPulse(g.seqCtx, cycleTime, 0.0)
				}
			}
		}(block)
	}
}

func (g Grid) AlternatingMexicanWave(cycleTime float64) {
	for _, block := range g.blocks {
		go func(block *Block) {
			if block.Y%2 == 0 {
				block.LightPulse(g.seqCtx, cycleTime, float64(block.X)/float64(g.x))
			} else {
				block.LightPulse(g.seqCtx, cycleTime, 1.0-float64(block.X)/float64(g.x))
			}
			for {
				select {
				case <-g.seqCtx.Done():
					return
				default:
					block.LightPulse(g.seqCtx, cycleTime, 0.0)
				}
			}
		}(block)
	}
}

func (g Grid) Wave(cycleTime float64) {
	for _, block := range g.blocks {
		go func(block *Block) {
			// fmt.Println(float64(block.X) / float64(g.x) * cycleTime)
			time.Sleep(time.Duration(float64(block.X)/float64(g.x)*cycleTime*1000) * time.Millisecond)
			fmt.Println("starting", block.Name)
			for {
				select {
				case <-g.seqCtx.Done():
					return
				default:
					block.LightOn()
					time.Sleep(time.Duration(cycleTime*1000) * time.Millisecond)
					block.LightOff()
					time.Sleep(time.Duration(cycleTime*1000) * time.Millisecond)
				}
			}
		}(block)
	}
}

func (g Grid) AlternatingWave(cycleTime float64) {
	for _, block := range g.blocks {
		go func(block *Block) {
			// fmt.Println(float64(block.X) / float64(g.x) * cycleTime)
			if block.Y%2 == 0 {
				time.Sleep(time.Duration(float64(block.X)/float64(g.x)*cycleTime*1000) * time.Millisecond)
			} else {
				time.Sleep(time.Duration((1-float64(block.X)/float64(g.x))*cycleTime*1000) * time.Millisecond)

			}
			fmt.Println("starting", block.Name)
			for {
				select {
				case <-g.seqCtx.Done():
					return
				default:
					block.LightOn()
					time.Sleep(time.Duration(cycleTime*1000) * time.Millisecond)
					block.LightOff()
					time.Sleep(time.Duration(cycleTime*1000) * time.Millisecond)
				}
			}
		}(block)
	}
}
