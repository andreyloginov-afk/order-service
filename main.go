package main

import (
	"context"
	"log"

	"github.com/andreyloginov-afk/order-service/internal/app/config"
	rhealth "github.com/andreyloginov-afk/order-service/internal/app/handler/http/health"
	rprocessor "github.com/andreyloginov-afk/order-service/internal/app/processor/http"
	rcpostgres "github.com/andreyloginov-afk/order-service/internal/app/repository/conn/postgres"
)

func main() {
	config.Load()

	cfg := config.Root

	db, err := rcpostgres.NewClient(context.Background(), cfg.Repository.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	log.Printf("connected to postgres at %s", cfg.Repository.Postgres.Address)

	hHealth := rhealth.NewHandler()

	proc := rprocessor.NewHTTP(hHealth, cfg.Processor.WebServer)

	serveErr := proc.Serve()

	if err := db.Close(); err != nil {
		log.Printf("failed to close db: %v", err)
	}

	if serveErr != nil {
		log.Fatalf("http server error: %v", serveErr)
	}
}
