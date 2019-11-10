package partybot

import (
	"context"
	"fmt"
	"math/rand"
)

func (g *Grid) TurnAllOff() {
	for x, _ := range g.blockArray {
		for y, _ := range g.blockArray[x] {
			g.blockArray[x][y].LightOff()
		}
	}
}

func (g *Grid) FadeAll() {
	for x, _ := range g.blockArray {
		for y, _ := range g.blockArray[x] {
			g.blockArray[x][y].LightFadeOut(g.seqCtx, 1.0, 0)
		}
	}
}

func (g *Grid) GetOffBlocks() (emptyBlocks []*Block) {
	for x, _ := range g.blockArray {
		for y, _ := range g.blockArray[x] {
			if !g.blockArray[x][y].LightState {
				emptyBlocks = append(emptyBlocks, g.blockArray[x][y])
			}
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

func (g Grid) MexicanWave(cycleTime float64) {
	for x, _ := range g.blockArray {
		for y, _ := range g.blockArray[x] {
			go func(ctx context.Context, block *Block) {
				block.LightPulse(ctx, cycleTime, float64(x)/float64(g.xLength))
				for {
					select {
					case <-ctx.Done():
						g.TurnAllOff()
						return
					default:
						block.LightPulse(ctx, cycleTime, 0.0)
					}
				}
			}(g.seqCtx, g.blockArray[x][y])
		}
	}
}

// func (g Grid) AlternatingMexicanWave(cycleTime float64) {
// 	for _, block := range g.blocks {
// 		go func(ctx context.Context, block *Block) {
// 			if block.Y%2 == 0 {
// 				block.LightPulse(ctx, cycleTime, float64(block.X)/float64(g.xLength))
// 			} else {
// 				block.LightPulse(ctx, cycleTime, 1.0-float64(block.X)/float64(g.xLength))
// 			}
// 			for {
// 				select {
// 				case <-ctx.Done():
// 					g.TurnAllOff()
// 					return
// 				default:
// 					block.LightPulse(ctx, cycleTime, 0.0)
// 				}
// 			}
// 		}(g.seqCtx, block)
// 	}
// }

// func (g Grid) Wave(cycleTime float64) {
// 	for _, block := range g.blocks {
// 		go func(ctx context.Context, block *Block) {
// 			time.Sleep(time.Duration(float64(block.X)/float64(g.xLength)*cycleTime*1000) * time.Millisecond)
// 			for {
// 				block.LightOn()
// 				if SleepCanBreak(ctx, cycleTime) {
// 					return
// 				}
// 				block.LightOff()
// 				if SleepCanBreak(ctx, cycleTime) {
// 					return
// 				}
// 			}
// 		}(g.seqCtx, block)
// 	}
// }

// func (g Grid) AlternatingWave(cycleTime float64) {
// 	for _, block := range g.blocks {
// 		go func(ctx context.Context, block *Block) {
// 			if block.Y%2 == 0 {
// 				time.Sleep(time.Duration(float64(block.X)/float64(g.xLength)*cycleTime*1000) * time.Millisecond)
// 			} else {
// 				time.Sleep(time.Duration((1-float64(block.X)/float64(g.xLength))*cycleTime*1000) * time.Millisecond)

// 			}
// 			for {
// 				select {
// 				case <-ctx.Done():
// 					g.TurnAllOff()
// 					return
// 				default:
// 					block.LightOn()
// 					if SleepCanBreak(ctx, cycleTime) {
// 						return
// 					}
// 					block.LightOff()
// 					if SleepCanBreak(ctx, cycleTime) {
// 						return
// 					}
// 				}
// 			}
// 		}(g.seqCtx, block)
// 	}
// }
