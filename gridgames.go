package partybot

import (
	"context"
	"math/rand"
)

// func (g *Grid) Snake() {
// 	go func(ctx context.Context) {
// 		score := 1
// 		target := g.GetRandomOffBlock()
// 		target.LightOn()
// 		for {
// 			if g.snakeDirection == 0 {
// 				g.snakeCoord.y = (g.snakeCoord.y + 1) % (g.yLength + 1)
// 			} else if g.snakeDirection == 1 {
// 				g.snakeCoord.x = (g.snakeCoord.x + 1) % (g.xLength + 1)
// 			} else if g.snakeDirection == 2 {
// 				g.snakeCoord.y = (g.snakeCoord.y - 1) % (g.yLength + 1)
// 			} else {
// 				g.snakeCoord.x = (g.snakeCoord.x - 1) % (g.xLength + 1)
// 			}
// 			if g.snakeCoord.y < 0 {
// 				g.snakeCoord.y = g.yLength
// 			}
// 			if g.snakeCoord.x < 0 {
// 				g.snakeCoord.x = g.xLength
// 			}
// 			// check if hit target or crashed
// 			if target.X == g.snakeCoord.x && target.Y == g.snakeCoord.y {
// 				score++
// 				target = g.GetRandomOffBlock()
// 				if target != nil {
// 					target.LightOn()
// 				} else {
// 					g.FadeAll()
// 					return
// 				}
// 			} else if g.blockArray[g.snakeCoord.y][g.snakeCoord.x].state.LightOn {
// 				g.TurnAllOff()
// 				return
// 			}
// 			g.blockArray[g.snakeCoord.y][g.snakeCoord.x].LightOn()
// 			// The bit that controls the light staying on
// 			go func(co coOrd) {
// 				SleepCanBreak(ctx, 1.0)
// 				g.blockArray[co.y][co.x].LightOff()

// 			}(g.snakeCoord)
// 			// The bit that waits to move the head forward
// 			if SleepCanBreak(ctx, 1.0/float64(score)) {
// 				return
// 			}
// 		}
// 	}(g.seqCtx)
// }

func (g *Grid) RandomSnake(cycleTime float64) {
	var x, y int
	x = 0
	y = 0
	var step, axis int
	go func(ctx context.Context) {
		for {
			step = 1 - 2*rand.Intn(2)
			axis = rand.Intn(2)
			if axis == 0 {
				x = (x + step) % g.xLength
				if x < 0 {
					x = g.xLength - 1
				}
			} else {
				y = (y + step) % g.yLength
				if y < 0 {
					y = g.yLength - 1
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
