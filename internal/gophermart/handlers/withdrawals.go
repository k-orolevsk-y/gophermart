package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
)

func (hs *handlerService) GetWithdrawals(ctx *gin.Context) {
	tokenClaims, exists := hs.GetTokenClaims(ctx)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	withdrawals, err := hs.pg.UserWithdraw().GetAllByUserID(ctx, tokenClaims.UserID)
	if err != nil {
		hs.logger.Error("error get user withdrawals", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, withdrawals)
}
