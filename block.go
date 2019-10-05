package partybot

import (
	"fmt"
	"log"
	"net/http"
)

func NewBlock() *Block {

	return &Block{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// serveWs handles websocket requests from the peer.
func (block *Block) ServeWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("New connection")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		c := &Client{block: block, conn: conn, send: make(chan []byte, 256)}
		c.block.register <- c

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go c.writePump()
		go c.readPump()
	}
}

func (b *Block) Send(sendable string) {
	b.broadcast <- []byte(sendable)
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
