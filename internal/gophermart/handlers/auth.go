package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (hs *handlerService) Register(ctx *gin.Context) {
}

func (hs *handlerService) Login(ctx *gin.Context) {
	var data struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&data); err != nil {
		hs.logger.Error("ctx.ShouldBindJSON", zap.Error(err))

		ctx.Status(http.StatusBadRequest)
		return
	}

	if data.Login != "test" && data.Password != "test" {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()

		return
	}

	token, err := hs.jwt.Encode("test")
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Abort()

		return
	}

	ctx.Header("Authorization", token)
	ctx.Status(http.StatusOK)
	ctx.Abort()
}
