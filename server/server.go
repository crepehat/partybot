package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/crepehat/partybot"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	block1 := partybot.NewBlock()
	block2 := partybot.NewBlock()
	block1.Start()
	block2.Start()
	mux := http.NewServeMux()
	mux.HandleFunc("/1", block1.ServeWs())
	mux.HandleFunc("/2", block2.ServeWs())

	go func() {
		for {
			block1.Send("Swoop")
			block2.Send("Sweep")
			time.Sleep(time.Second)
		}
	}()
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
