package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/repository"
)

func (hs *handlerService) GetOrders(ctx *gin.Context) {
	tokenClaims, exists := hs.GetTokenClaims(ctx)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	orders, err := hs.pg.Order().GetAllByUserID(ctx, tokenClaims.UserID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, orders)
}

func (hs *handlerService) NewOrder(ctx *gin.Context) {
	if ctx.ContentType() != "text/plain" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.NewBadRequestErrorResponse("An unsupported request body was sent"))
		return
	}

	bs, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.NewValidationErrorResponse("The order number was not transmitted"))
		return
	}

	body := string(bs)
	if body == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.NewBadRequestErrorResponse("An empty request body was sent"))
		return
	}

	orderNumber, err := hs.CheckNumberAlgorithmLuna(body)
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

	order := models.Order{
		ID:      orderNumber,
		UserID:  tokenClaims.UserID,
		Status:  "NEW",
		Accrual: nil,
	}

	if err = hs.pg.Order().Create(ctx, &order); err != nil {
		pgError := hs.pg.ParsePgError(err)

		switch pgError.Message {
		case repository.ErrorOrderByThisUser:
			ctx.AbortWithStatus(http.StatusOK)
		case repository.ErrorOrderByOtherUser:
			ctx.AbortWithStatusJSON(http.StatusConflict, models.NewConflictErrorResponse("An order with this number has already been uploaded by another user"))
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		}

		return
	}

	ctx.AbortWithStatus(http.StatusAccepted)
}
