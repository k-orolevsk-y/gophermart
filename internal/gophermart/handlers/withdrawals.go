package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, withdrawals)
}
