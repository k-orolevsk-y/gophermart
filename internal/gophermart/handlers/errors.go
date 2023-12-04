package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
)

func (hs *handlerService) NoRoute(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotFound)
}

func (hs *handlerService) NoMethod(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusMethodNotAllowed)
}

func (hs *handlerService) Recovery(ctx *gin.Context, err any) {
	hs.api.GetLogger().Error("Panic on request", zap.Any("error", err))
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
}
