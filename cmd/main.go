// cmd/main.go
package main

import (
	"context"
	"errors"
	"github.com/chatbox/whatsapp/internal"
	"github.com/chatbox/whatsapp/internal/config"
	"github.com/chatbox/whatsapp/internal/handler"
	"github.com/chatbox/whatsapp/internal/repository"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chatbox/whatsapp/internal/infrastructure"
)

func main() {
	ctx := context.Background()

	// Load config (env, flags, etc.)
	cfg := config.LoadConfig("config/config.yaml")
	logger := infrastructure.NewLogger()

	db, err := infrastructure.NewDB(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}

	repo := repository.NewWuzRepo(db)
	waClient := infrastructure.NewWhatsAppClient(cfg.DBConfig.DatabaseUrl, logger)

	service := internal.NewWuzApi(repo, waClient, &logger)

	// Reconnect sessions from DB
	if err := service.ReconnectActiveSessions(ctx); err != nil {
		logger.Error().Err(err).Msg("Failed to reconnect sessions")
	}

	router := handler.NewHandler(&logger, service)

	srv := &http.Server{
		Addr:         ":" + cfg.AppEnv.Port,
		Handler:      router.RegisterRoutes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		logger.Info().Str("addr", srv.Addr).Msg("HTTP server listening")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	// wait for SIGINT or SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logger.Info().Msg("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server shutdown failed")
	} else {
		logger.Info().Msg("Server exited cleanly")
	}
}
