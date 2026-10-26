package main

import (
	"context"
	"errors"
	"log"
	"os/signal"
	"syscall"

	"github.com/caiohperlin/wallace/api/internal/config"
	"github.com/caiohperlin/wallace/api/internal/database"
	"github.com/caiohperlin/wallace/api/internal/http"
	"github.com/caiohperlin/wallace/api/internal/server"
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

	hub := http.NewHub()
	router := http.NewRouter(http.Dependencies{
		Pinger:   db,
		Store:    db,
		Token:    appConfig.APIToken,
		Notifier: hub,
		Hub:      hub,
	})
	httpServer := server.New(appConfig.HTTPAddress, router)
	return httpServer.Run(ctx)
}
