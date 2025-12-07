package main

import (
	"log/slog"
	"net/http"

	"github.com/inahym196/gameHub/backend/pkg/hub"
	gamehub "github.com/inahym196/reversi/pkg/gameHub"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	hub := hub.NewHub(nil)
	http.Handle("/", gamehub.NewServer(hub))
	if err := http.ListenAndServe("localhost:3000", nil); err != nil {
		slog.Error(err.Error())
	}
}
