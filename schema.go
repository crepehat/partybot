package partybot

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type coOrd struct {
	x int
	y int
}

type Grid struct {
	// blocks         []*Block
	xLength        int
	yLength        int
	snakeCoord     coOrd
	snakeDirection int //0=N,1=E,2=S,3=W
	// array          [][]string
	// x and y are reversed
	blockArray [][]*Block

	seqLock   *sync.Mutex
	seqCtx    context.Context
	seqCancel context.CancelFunc

	// for monitoring changes
	changeCHAN chan Block

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// Client is a middleman between the websocket connection and the grid.
type Client struct {
	// block *Block
	grid *Grid

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

type Block struct {
	X              int     `json:"x"`
	Y              int     `json:"y"`
	Name           string  `json:"name"`
	LightMagnitude float64 `json:"light_magnitude"`
	LightState     bool    `json:"light_state"`

	// each block contains reference to grid for interactions and updating
	grid *Grid
}

// Block maintains the set of active clients and broadcasts messages to the
// clients.
// type Block struct {
// 	Name string `json:"name"`
// 	X    int    `json:"x"`
// 	Y    int    `json:"y"`

// 	state *BlockState
// }
