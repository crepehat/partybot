package partybot

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
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

		// prepend to maintain correct order
		nameGrid = append([][]string{row}, nameGrid...)
	}
	return
}

// NewGrid accepts a gird of names with coordinates [y][x]
// y,x gels better both with reading csvs and html tables
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
		xLength:    len(nameGrid[0]),
		yLength:    len(nameGrid),
	}
	g.seqCtx, g.seqCancel = context.WithCancel(context.Background())

	// Create new block for each value of the supplied grid
	for y, row := range nameGrid {
		for x, name := range row {
			b := g.NewBlock(name, x, y)
			if len(g.blockArray) <= y {
				g.blockArray = append(g.blockArray, []*Block{b})
			} else {
				g.blockArray[y] = append(g.blockArray[y], b)
			}
		}
	}
	fmt.Printf("Initialised grid with dimensions %dx%d\n", g.xLength, g.yLength)
	return
}

func (g *Grid) StartWebsocketServer() {
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
				// fmt.Println(string(message))
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

func (g *Grid) StartMonitor() {
	changeMap := make(map[string]Block)
	g.changeCHAN = make(chan Block)
	var sendLock sync.Mutex
	// get all changes and write them to the map
	go func() {
		for change := range g.changeCHAN {
			sendLock.Lock()
			changeMap[change.Name] = change
			sendLock.Unlock()
		}
	}()

	go func() {
		sendTick := time.NewTicker(250 * time.Millisecond)
		for {
			sendLock.Lock()
			changes, err := json.Marshal(changeMap)
			if err != nil {
				fmt.Printf("Error marshalling changes:")
			}
			g.Broadcast(changes)
			changeMap = make(map[string]Block)

			sendLock.Unlock()
			<-sendTick.C
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

func (g *Grid) Broadcast(payload []byte) {
	// fmt.Println("broadcasting message:", payload)
	g.broadcast <- payload
}
