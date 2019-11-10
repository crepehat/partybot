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
