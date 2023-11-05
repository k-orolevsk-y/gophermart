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
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/handlers"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/middlewares"
	"github.com/k-orolevsk-y/gophermart/pkg/logger"
)

func main() {
	if err := config.ParseConfig(); err != nil {
		panic(err)
	}

	log, err := logger.New()
	if err != nil {
		panic(err)
	}

	log.Info("init logger")
	log.Info("parsed config", zap.Any("Config", config.Config))

	apiService := api.New(log)
	middlewares.ConfigureMiddlewaresService(apiService)
	handlers.ConfigureHandlersService(apiService)
	apiService.Run()

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("shutting down gracefully", zap.Any("signal", <-quitSignal))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err = apiService.Shutdown(ctx); err != nil {
		log.Fatal("fatal error shutdown HTTP server", zap.Error(err))
	}

	log.Info("http server status", zap.Bool("shutdown", apiService.WaitShutdown(ctx)))
	log.Info("successfully shutdown server gracefully")
}
