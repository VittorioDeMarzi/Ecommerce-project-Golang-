package main

import (
	"log/slog"
	"os"
)

func main() {
	cfg := config{
		addr: "127.0.0.1:8080",
		db: dbConfig{},
	}

	api := application{
		config: cfg,

	}

	// structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)


	if err := api.run(api.mount()); err != nil {
		slog.Error("Server has failed to start", "error", err)
		os.Exit(1)
	}

}
