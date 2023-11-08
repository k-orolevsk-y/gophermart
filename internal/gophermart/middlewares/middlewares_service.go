package middlewares

import (
	"github.com/gin-gonic/gin"
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
}

func (ms *middlewaresService) Logger(_ *gin.Context) {
	ms.logger.Info("new request")
}
