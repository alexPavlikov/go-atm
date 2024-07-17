package main

import (
	"log/slog"

	server "github.com/alexPavlikov/go-atm/cmd"
)

func main() {
	if err := server.Run(); err != nil {
		slog.Error("failed to start server", "error", err)
		return
	}
}
