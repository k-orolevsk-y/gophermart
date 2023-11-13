package middlewares

import (
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/pkg/router"
)

type middlewaresService struct {
	logger *zap.Logger
}

type apiService interface {
	GetRouter() *router.Router
	GetLogger() *zap.Logger
}

func ConfigureMiddlewaresService(api apiService) {
	ms := &middlewaresService{
		logger: api.GetLogger(),
	}

	api.GetRouter().Use(ms.Logger)
	api.GetRouter().Use(ms.Compress)
}
