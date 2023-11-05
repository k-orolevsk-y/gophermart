package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type middlewaresService struct {
	logger *zap.Logger
}

type apiService interface {
	GetRouter() *gin.Engine
	GetLogger() *zap.Logger
}

func ConfigureMiddlewaresService(api apiService) {
	ms := &middlewaresService{}
	api.GetRouter().Use(ms.Logger)
}

func (ms *middlewaresService) Logger(ctx *gin.Context) {
	ms.logger.Info("new request", zap.Any("ctx", ctx))
}
