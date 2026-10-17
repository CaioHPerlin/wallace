package main

import (
	"context"
	"errors"
	"log"
	"os/signal"
	"syscall"

	"github.com/wallace-project/wallace/api/internal/config"
	"github.com/wallace-project/wallace/api/internal/database"
	"github.com/wallace-project/wallace/api/internal/http"
	"github.com/wallace-project/wallace/api/internal/server"
)

func main() {
	if err := run(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}

func run() error {
	appConfig, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := database.Open(ctx, appConfig.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	router := http.NewRouter(db)
	httpServer := server.New(appConfig.HTTPAddress, router)
	return httpServer.Run(ctx)
}
