package handlers

import (
	"io"
	"net/http"
	"strconv"

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

//Хендлер: `POST /api/user/orders`.
//
//Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.
//
//Номер заказа может быть проверен на корректность ввода с помощью [алгоритма Луна](https://ru.wikipedia.org/wiki/Алгоритм_Луна){target="_blank"}.
//
//Формат запроса:
//
//```
//POST /api/user/orders HTTP/1.1
//Content-Type: text/plain
//...
//
//12345678903
//```
//Возможные коды ответа:
//
//- `200` — номер заказа уже был загружен этим пользователем;
//- `202` — новый номер заказа принят в обработку;
//- `400` — неверный формат запроса;
//- `401` — пользователь не аутентифицирован;
//- `409` — номер заказа уже был загружен другим пользователем;
//- `422` — неверный формат номера заказа;
//- `500` — внутренняя ошибка сервера.

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

	orderNumber, err := strconv.ParseInt(body, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.NewValidationErrorResponse("The order number is not a number"))
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

	if !order.CheckValidID() {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, models.NewUnprocessableEntityErrorResponse("The order number does not match the Luna algorithm."))
		return
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

	ctx.AbortWithStatus(http.StatusCreated)
}
