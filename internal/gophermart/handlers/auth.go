package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
)

func (hs *handlerService) Register(ctx *gin.Context) {
	var data struct {
		Login    string `json:"login" validate:"required|minLen:3"`
		Password string `json:"password" validate:"required|minLen:6|maxLen:30"`
	}

	if err := ctx.ShouldBindWith(&data, hs.bindingWithValidation); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.NewValidationErrorResponse(err))
		return
	}

	user := models.User{
		Login:    data.Login,
		Password: fmt.Sprint(data.Password),
	}
	if err := hs.pg.User().Create(ctx, &user); err != nil {
		code := hs.pg.GetCodeFromError(err)
		if pgerrcode.IsIntegrityConstraintViolation(code) {
			ctx.AbortWithStatusJSON(http.StatusConflict, models.NewConflictErrorResponse("A user with this login is already registered"))
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		}

		return
	}

	tokenString, err := hs.jwt.Encode(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		return
	}

	ctx.Header("Authorization", tokenString)
	ctx.AbortWithStatus(http.StatusOK)
}

func (hs *handlerService) Login(ctx *gin.Context) {
	var data struct {
		Login    string `json:"login" validate:"required|minLen:3"`
		Password string `json:"password" validate:"required|minLen:6|maxLen:30"`
	}

	if err := ctx.ShouldBindWith(&data, hs.bindingWithValidation); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.NewValidationErrorResponse(err))
		return
	}

	user, err := hs.pg.User().GetByLogin(ctx, data.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.NewAuthorizationErrorResponse("Invalid login or password"))
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		}

		return
	}

	if !user.CheckPassword(data.Password) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.NewAuthorizationErrorResponse("Invalid login or password"))
		return
	}

	token, err := hs.jwt.Encode(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.NewInternalServerErrorResponse())
		ctx.Abort()

		return
	}

	ctx.Header("Authorization", token)
	ctx.AbortWithStatus(http.StatusOK)
}
