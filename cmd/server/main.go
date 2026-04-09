package main

import (
	"log"
	"net/http"

	"github.com/xiaofeihuang85/go-crypto-product-service/internal/api"
	"github.com/xiaofeihuang85/go-crypto-product-service/internal/config"
)

func main() {
	cfg := config.Load()
	server := &http.Server{
		Addr:         cfg.Address(),
		Handler:      api.NewRouter(cfg.ServiceName),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Printf("starting %s on %s", cfg.ServiceName, cfg.Address())

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
