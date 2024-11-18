package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chpter/shared/config"
	"github.com/chpter/shared/logger"
)

func main() {
	// setup logging
	logHandler := logger.New(config.GlobalConfig.LogLevel)

	slog.SetDefault(logHandler)

	ctx := context.Background()
	cfg := config.GlobalConfig

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// no need to run this in a goroutine, since it already starts the server in a separate goroutine
	srv, err := startGrpcServer(ctx, cfg)
	if err != nil {
		os.Exit(1)
	}

	// block main goroutine until a signal is received
	<-sigChan

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Close(); err != nil {
		slog.Error("Graceful shutdown failed with: ", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Server shutdown successfully")
}
