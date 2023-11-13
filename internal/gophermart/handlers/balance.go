package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/repository"
)

func (hs *handlerService) GetBalance(ctx *gin.Context) {
	tokenClaims, exists := hs.GetTokenClaims(ctx)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	user, err := hs.pg.User().GetByID(ctx, tokenClaims.UserID)
	if err != nil {

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	sumWithdrawn, err := hs.pg.UserWithdraw().GetWithdrawnSumByUserID(ctx, user.ID)
	if err != nil {

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"balance": user.Balance, "withdrawn": sumWithdrawn})
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
		if errors.Is(err, errInvalidTypeOfNumberForAlogirthmLuna) {
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

	user, err := hs.pg.User().GetByID(ctx, tokenClaims.UserID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	withdraw := models.UserWithdraw{
		UserID:  user.ID,
		OrderID: orderNumber,
		Sum:     data.Sum,
	}
	if err = hs.pg.UserWithdraw().Create(ctx, &withdraw, user); err != nil {
		if errors.Is(err, repository.ErrorInsufficientFunds) {
			ctx.AbortWithStatusJSON(http.StatusPaymentRequired, models.NewPaymentRequiredErrorResponse("Insufficient funds in the account"))
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		}
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}
