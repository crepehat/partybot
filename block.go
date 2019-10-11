package partybot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func NewBlock(name string, x, y int) *Block {

	return &Block{
		Name: name,
		X:    x,
		Y:    y,
		state: &BlockState{
			LightMagnitude: 0.0,
			LightOn:        false,
		},
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (b *Block) SetLight(on bool, mag float64) (err error) {
	b.state.LightMagnitude = mag
	b.state.LightOn = on
	b.Send()
	return nil
}

func (b *Block) LightFadeIn(ctx context.Context, duration, start float64) {
	// Do 250ms pulses
	totalSteps := duration * 8
	for i := start * totalSteps; i <= totalSteps; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			b.SetLight(true, 1/totalSteps*i)
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

// serveWs handles websocket requests from the peer.
func (b *Block) ServeWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("New connection")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		c := &Client{block: b, conn: conn, send: make(chan []byte, 256)}
		c.block.register <- c

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go c.writePump()
		go c.readPump()
	}
}

func (b *Block) Send() {
	state, err := json.Marshal(b.state)
	if err != nil {
		fmt.Println(err)
	}
	b.broadcast <- state
}

func (b *Block) Start() {
	go func() {
		for {
			select {
			case client := <-b.register:
				b.clients[client] = true
			case client := <-b.unregister:
				if _, ok := b.clients[client]; ok {
					delete(b.clients, client)
					close(client.send)
				}
			case message := <-b.broadcast:
				for client := range b.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(b.clients, client)
					}
				}
			}
		}
	}()
}
