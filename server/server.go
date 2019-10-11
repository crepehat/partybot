package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/crepehat/partybot"
	"github.com/gobuffalo/packr/v2"
)

var (
	addr     = ":8080"
	gridFile = "grid.csv"
)

func init() {
	flag.StringVar(&addr, "addr", addr, "http service address")
	flag.StringVar(&gridFile, "gridFile", gridFile, "Layout file for site")
}

func main() {

	var gridSlice [][]string

	flag.Parse()
	fh, err := os.Open(gridFile)
	if err != nil {
		fmt.Println("Error opening gridfile. Double check it.", err)
	}
	r := csv.NewReader(fh)

	for {
		line, error := r.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		gridSlice = append(gridSlice, line)
	}

	grid, err := partybot.NewGrid(gridSlice)
	if err != nil {
		fmt.Println(err)
	}

	grid.PrintBlock(0, 8)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", grid.GetMux()))

	box := packr.New("reactAssets", "../frontend/build")
	mux.Handle("/", http.FileServer(box))

	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
