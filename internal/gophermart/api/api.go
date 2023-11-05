package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
)

type APIService struct {
	router *gin.Engine
	logger *zap.Logger

	srv       *http.Server
	srvClosed bool
}

func New(logger *zap.Logger) *APIService {
	if config.Config.ProductionMode {
		gin.SetMode(gin.ReleaseMode)
	}

	return &APIService{
		router: gin.New(),
		logger: logger,
	}
}

func (a *APIService) GetRouter() *gin.Engine {
	return a.router
}

func (a *APIService) GetLogger() *zap.Logger {
	return a.logger
}

func (a *APIService) GetRepository() {
	return
}

func (a *APIService) Run() {
	a.router.GET("/", func(c *gin.Context) {
		time.Sleep(time.Second * 15)
		c.JSON(http.StatusOK, "OK")
	})

	a.srv = &http.Server{
		Addr:    config.Config.RunAddress,
		Handler: a.router,
	}

	a.logger.Info("starting HTTP server", zap.String("addr", a.srv.Addr))
	go func() {
		if err := a.srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				a.srvClosed = true
			} else {
				a.logger.Panic("error starting HTTP server", zap.Error(err))
			}
		}
	}()
}

func (a *APIService) Shutdown(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}

func (a *APIService) WaitShutdown(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		if a.srvClosed {
			return true
		}

		time.Sleep(time.Second)
	}

	return false
}
