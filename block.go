package partybot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func NewBlock(name string, x, y int) *Block {

	return &Block{
		Name:           name,
		LightMagnitude: 0.0,
		LightState:     false,
		X:              x,
		Y:              y,
	}
}

func (b *Block) UpdateClients() (err error) {
	blockJSON, err := json.Marshal(b)
	if err != nil {
		return err
	}

	b.grid.broadcast <- blockJSON
	return
}

func (b *Block) SetLight(on bool, mag float64) (err error) {
	b.LightMagnitude = mag
	b.LightState = on
	return b.UpdateClients()
}

func (b *Block) LightFadeIn(ctx context.Context, duration, start float64) {
	// Do 250ms pulses
	totalSteps := duration * 8
	for i := start * totalSteps; i <= totalSteps; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			err := b.SetLight(true, 1/totalSteps*i)
			if err != nil {
				fmt.Println("error setting light: ", err)
				return
			}
			time.Sleep(125 * time.Millisecond)
		}
	}
}

// LightFadeOut fades the light out over [seconds]
// starting at [start] as a portion of the total time
func (b *Block) LightFadeOut(ctx context.Context, duration, start float64) {
	// Do 250ms pulses
	totalSteps := duration * 8
	for i := start * totalSteps; i <= totalSteps; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			b.SetLight(true, 1-1/totalSteps*i)
			time.Sleep(125 * time.Millisecond)
		}
	}
}

func (b *Block) LightOn() {
	b.SetLight(true, 1)
}

func (b *Block) LightOff() {
	b.SetLight(false, 0)
}

func (b Block) LightPulse(ctx context.Context, duration, start float64) {
	b.LightFadeIn(ctx, duration/2, start)
	b.LightFadeOut(ctx, duration/2, 0.0)
}
