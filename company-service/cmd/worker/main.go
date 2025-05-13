package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"company-service/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)
	_ = config.NewDatabase(viperConfig, logger)

	logger.Info("Starting worker service")

	_, cancel := context.WithCancel(context.Background())

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	stop := false
	for !stop {
		s := <-terminateSignals
		logger.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)
		cancel()
		stop = true
	}

	time.Sleep(5 * time.Second) // wait for all consumers to finish processing
}
