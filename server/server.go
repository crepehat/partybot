package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/crepehat/partybot"
	"github.com/gobuffalo/packr/v2"
)

var (
	addr     = ":8080"
	gridFile = "./grid.csv"
)

func init() {
	flag.StringVar(&addr, "addr", addr, "http service address")
	flag.StringVar(&gridFile, "gridFile", gridFile, "Layout file for site")
}

func main() {

	var nameGrid [][]string

	flag.Parse()

	nameGrid, err := partybot.ReadGridFile(gridFile)

	grid, err := partybot.NewGrid(nameGrid)
	if err != nil {
		fmt.Println(err)
	}

	grid.StartWebsocketServer()
	grid.StartMonitor()

	// go func() {
	// 	counter := 0
	// 	for {
	// 		time.Sleep(250 * time.Millisecond)
	// 		grid.Broadcast(fmt.Sprintf("%d", counter))
	// 		counter++
	// 		if counter > 100 {
	// 			counter = 0
	// 		}
	// 	}
	// }()

	go func() {
		grid.Test()
	}()

	// go func() {
	// 	type rng struct {
	// 		Key   int `json:"key"`
	// 		Value int `json:"value"`
	// 	}
	// 	for i := 0; i < 300; i++ {
	// 		go func(i int) {
	// 			counter := 0
	// 			for {
	// 				rngSus := rng{
	// 					Key:   i,
	// 					Value: rand.Intn(10),
	// 				}
	// 				time.Sleep(250 * time.Millisecond)
	// 				rngSusBytes, err := json.Marshal(rngSus)
	// 				if err != nil {
	// 					fmt.Println(err)
	// 				}
	// 				grid.Broadcast(string(rngSusBytes))
	// 				counter++
	// 				if counter > 100 {
	// 					counter = 0
	// 				}
	// 			}
	// 		}(i)
	// 	}
	// }()

	// grid.PrintBlock(0, 8)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", grid.GetMux()))

	box := packr.New("reactAssets", "../frontend/build")
	mux.Handle("/", http.FileServer(box))

	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
