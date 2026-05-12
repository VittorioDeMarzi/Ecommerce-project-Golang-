package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/VittorioDeMarzi/Ecommerce-project-Golang-/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()


	cfg := config{
		addr: "127.0.0.1:8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost port=5432 dbname=ecommerce user=user password=password sslmode=disable"),
		},
	}

	// structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		slog.Error("DB connection failed", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	api := application{
		config: cfg,
		db: conn,

	}

	logger.Info("Connected to db", "dsn", cfg.db.dsn)

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server has failed to start", "error", err)
		os.Exit(1)
	}

}
