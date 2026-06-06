package main

import (
	"log/slog"
	"os"

	"github.com/gamis65/tinyfilehost/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	router := server.NewRouter(logger)

	srv := server.New("localhost:3000", router)

	logger.Info("Server running")
	srv.Start()
}
