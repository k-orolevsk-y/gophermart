package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/repository"
)

func (hs *handlerService) GetBalance(ctx *gin.Context) {
	tokenClaims, exists := hs.GetTokenClaims(ctx)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	sumWithdrawn, err := hs.pg.UserWithdraw().GetWithdrawnSumByUserID(ctx, tokenClaims.UserID)
	if err != nil {
		hs.logger.Error("error get sum withdrawn", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	sumAccrual, err := hs.pg.Order().GetAccrualSumByUserID(ctx, tokenClaims.UserID)
	if err != nil {
		hs.logger.Error("error get sum accrual", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"current": sumAccrual - sumWithdrawn, "withdrawn": sumWithdrawn})
}

func (hs *handlerService) NewBalanceWithdrawn(ctx *gin.Context) {
	var data struct {
		Order string  `json:"order" validate:"required|minLen:2|maxLen:19"`
		Sum   float64 `json:"sum" validate:"required|min:1"`
	}

	if err := ctx.ShouldBindWith(&data, hs.bindingWithValidation); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.NewValidationErrorResponse(err))
		return
	}

	orderNumber, err := hs.CheckNumberAlgorithmLuna(data.Order)
	if err != nil {
		if errors.Is(err, errInvalidTypeOfNumberForAlgorithmLuna) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.NewValidationErrorResponse("The order number is not a number"))
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, models.NewUnprocessableEntityErrorResponse("The order number does not match the Luna algorithm."))
		}

		return
	}

	tokenClaims, exists := hs.GetTokenClaims(ctx)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	withdraw := models.UserWithdraw{
		UserID:  tokenClaims.UserID,
		OrderID: orderNumber,
		Sum:     data.Sum,
	}
	if err = hs.pg.UserWithdraw().Create(ctx, &withdraw); err != nil {
		if errors.Is(err, repository.ErrorInsufficientFunds) {
			ctx.AbortWithStatusJSON(http.StatusPaymentRequired, models.NewPaymentRequiredErrorResponse("Insufficient funds in the account"))
		} else {
			hs.logger.Error("error create user withdraw", zap.Error(err))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		}
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}
