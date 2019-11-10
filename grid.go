package partybot

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func ReadGridFile(gridFile string) (nameGrid [][]string, err error) {
	fh, err := os.Open(gridFile)
	if err != nil {
		return nil, fmt.Errorf("Error opening gridfile. Double check it. %s", err.Error())
	}

	r := csv.NewReader(fh)
	for {
		row, error := r.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			return nil, err
		}

		// transpose
		for x, column := range row {
			// prepending so not inverted
			if len(nameGrid) <= x {
				nameGrid = append(nameGrid, []string{column})
			} else {
				nameGrid[x] = append([]string{column}, nameGrid[x]...)
			}
		}
	}
	return
}

// NewGrid accepts a gird of names with coordinates [x][y]
func NewGrid(nameGrid [][]string) (g *Grid, err error) {
	if len(nameGrid) == 0 {
		return nil, fmt.Errorf("empty namegrid supplied")
	}

	g = &Grid{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		seqLock:    &sync.Mutex{},
		xLength:    len(nameGrid),
		yLength:    len(nameGrid[0]),
	}
	g.seqCtx, g.seqCancel = context.WithCancel(context.Background())

	// Create new block for each value of the supplied grid
	for x, column := range nameGrid {
		for y, name := range column {
			b := NewBlock(name, x, y)
			if len(g.blockArray) <= x {
				g.blockArray = append(g.blockArray, []*Block{b})
			} else {
				g.blockArray[x] = append(g.blockArray[x], b)
			}
		}
	}
	return
}

func (g *Grid) Start() {
	// Monitor for web clients, add to pool if new, send data to all when received
	go func() {
		for {
			select {
			case client := <-g.register:
				g.clients[client] = true
			case client := <-g.unregister:
				if _, ok := g.clients[client]; ok {
					delete(g.clients, client)
					close(client.send)
				}
			case message := <-g.broadcast:
				fmt.Println(message)
				for client := range g.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(g.clients, client)
					}
				}
			}
		}
	}()
}

// serveWs handles websocket requests from the peer.
func (g *Grid) ServeWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("New connection")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		c := &Client{grid: g, conn: conn, send: make(chan []byte, 256)}
		c.grid.register <- c

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go c.writePump()
		go c.readPump()
	}
}

func (g *Grid) GetMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/grid", g.handleArray())
	mux.HandleFunc("/sequence/start", g.handleSequenceStart())
	mux.HandleFunc("/sequence/stop", g.handleSequenceStop())
	mux.HandleFunc("/snake", g.handleSnake())
	mux.HandleFunc("/socket", g.ServeWs())

	return mux
}
