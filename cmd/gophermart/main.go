package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/api"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/repository"
	"github.com/k-orolevsk-y/gophermart/pkg/log"
)

func main() {
	if err := config.ParseConfig(); err != nil {
		panic(err)
	}

	logger, err := log.New()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // nolint

	logger.Info("parsed config", zap.Any("config", config.Config))

	rep, err := repository.NewPG(logger)
	if err != nil {
		logger.Panic("error connect database", zap.Error(err))
	}
	defer rep.Close()

	apiService := api.New(logger, rep)
	apiService.Configure()
	apiService.Run()

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("shutting down gracefully", zap.Any("signal", <-quitSignal))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err = apiService.Shutdown(ctx); err != nil {
		logger.Fatal("fatal error shutdown HTTP server", zap.Error(err))
	}

	logger.Info("http server status", zap.Bool("shutdown", apiService.WaitShutdown(ctx)))
	logger.Info("successfully shutdown server gracefully")
}
