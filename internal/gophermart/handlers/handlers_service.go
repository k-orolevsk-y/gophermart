package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type handlerService struct {
	logger *zap.Logger
}

type apiService interface {
	GetRouter() *gin.Engine
	GetLogger() *zap.Logger
}

func ConfigureHandlersService(api apiService) {
	hs := &handlerService{}
	api.GetRouter().GET("/", hs.Index)
}

func (hs *handlerService) Index(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[any]any{"status": "OK"})
}
