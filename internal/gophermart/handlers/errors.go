package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (hs *handlerService) NoRoute(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotFound)
}

func (hs *handlerService) NoMethod(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusMethodNotAllowed)
}
