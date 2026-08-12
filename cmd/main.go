package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pantan/config"
	"pantan/internal/handlers"
	"pantan/pkg/httpserver"
	"pantan/pkg/logger"
	"syscall"
	"time"
)

func main() {

	// Config
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("configuration read error: %v", err)
	}

	// Logger
	lgr := logger.NewLogger(cfg)

	lgr.Info().
		Str("version", config.Version).
		Msg("Starting pantan")

	srv := httpserver.NewServer()
	h := handlers.NewHandler(cfg)
	serverErr := make(chan error, 1)

	go func() {
		if err := srv.Run(cfg.Port, h.InitRoutes()); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case err := <-serverErr:
		lgr.Error().Err(err).Msg("Server runtime error")
	case sig := <-quit:
		lgr.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	lgr.Info().Msg("Shutting down server...")
	if err := srv.Stop(ctx); err != nil {
		lgr.Error().Err(err).Msg("Server shutdown error")
	}

	lgr.Info().Msg("Server stopped gracefully")

}
