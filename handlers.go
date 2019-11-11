package partybot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (g *Grid) handleSequenceStart() http.HandlerFunc {
	type startCommand struct {
		Name         string  `json:"name"`
		CycleSeconds float64 `json:"cycle_seconds"`
	}
	type startResponse struct {
		Response string `json:"response"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var command startCommand
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&command)
		if err != nil {
			fmt.Println(err)
			return
		}
		g.seqLock.Lock()
		g.seqCancel()
		g.TurnAllOff()
		// time.Sleep(1 * time.Second)
		g.seqCtx, g.seqCancel = context.WithCancel(context.Background())
		g.seqLock.Unlock()
		switch command.Name {
		// case "wave":
		// 	g.Wave(command.CycleSeconds)
		case "mexican_wave":
			g.MexicanWave(command.CycleSeconds)
		// case "alt_mexican_wave":
		// 	g.AlternatingMexicanWave(command.CycleSeconds)
		// case "alt_wave":
		// 	g.AlternatingWave(command.CycleSeconds)
		case "random_snake":
			g.RandomSnake(command.CycleSeconds)
			// case "snake":
			// 	g.Snake()
		}
		res := startResponse{
			Response: fmt.Sprintf("Started %s", command.Name),
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (g *Grid) handleSequenceStop() http.HandlerFunc {
	type stopResponse struct {
		Response string `json:"response"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		g.seqLock.Lock()
		g.seqCancel()
		g.TurnAllOff()
		g.seqLock.Unlock()
		res := stopResponse{
			Response: "Stopped sequence",
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (g *Grid) handleArray() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(g.blockArray)
	}
}

func (g *Grid) handleSnake() http.HandlerFunc {
	type snakeCommand struct {
		Direction string `json:"direction"`
	}
	type snakeResponse struct {
		Response string `json:"response"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var command snakeCommand
		w.Header().Set("Access-Control-Allow-Origin", "*")
		err := json.NewDecoder(r.Body).Decode(&command)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch command.Direction {
		case "ArrowDown":
			g.snakeDirection = 0
		case "ArrowRight":
			g.snakeDirection = 1
		case "ArrowUp":
			g.snakeDirection = 2
		case "ArrowLeft":
			g.snakeDirection = 3
		}
		resp := snakeResponse{Response: command.Direction}
		json.NewEncoder(w).Encode(resp)
	}
}
