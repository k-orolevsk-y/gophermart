package handlers

import (
	"testing"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/middlewares"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/mocks/api"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/pkg/jwt"
)

func NewTestAPI(t *testing.T) *api.TestAPI {
	testAPI := api.NewTestAPI(t)

	ConfigureHandlersService(testAPI)
	middlewares.ConfigureMiddlewaresService(testAPI)

	return testAPI
}

func GetUserIDWithToken() (string, uuid.UUID, error) {
	id := uuid.New()
	j := jwt.New(config.Config.HmacTokenSecret)

	tokenString, err := j.Encode(&models.User{ID: id})
	return tokenString, id, err
}
