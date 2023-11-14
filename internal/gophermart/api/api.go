package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/handlers"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/middlewares"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/order_pool"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/repository"
	"github.com/k-orolevsk-y/gophermart/pkg/router"
)

type APIService struct {
	router    *router.Router
	logger    *zap.Logger
	orderPool *orderpool.OrderPool
	pg        *repository.Pg

	srv       *http.Server
	srvClosed bool
}

func New(logger *zap.Logger, orderPool *orderpool.OrderPool, pg *repository.Pg) *APIService {
	if config.Config.ProductionMode {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = router.NewRouterLogger(logger)

	return &APIService{
		router:    router.New(),
		logger:    logger,
		orderPool: orderPool,
		pg:        pg,
	}
}

func (a *APIService) Configure() {
	middlewares.ConfigureMiddlewaresService(a)
	handlers.ConfigureHandlersService(a)
}

func (a *APIService) GetRouter() *router.Router {
	return a.router
}

func (a *APIService) GetLogger() *zap.Logger {
	return a.logger
}

func (a *APIService) GetOrderPool() *orderpool.OrderPool {
	return a.orderPool
}

func (a *APIService) GetPg() *repository.Pg {
	return a.pg
}

func (a *APIService) Run() {
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
