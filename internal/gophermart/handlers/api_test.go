package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/middlewares"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/mocks"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/pkg/jwt"
)

func NewTestAPI(t *testing.T) *mocks.TestAPI {
	if config.Config.HmacTokenSecret == "" {
		GenerateHmacTokenSecret(t)
	}

	api := mocks.NewTestAPI(t)

	ConfigureHandlersService(api)
	middlewares.ConfigureMiddlewaresService(api)

	return api
}

func GenerateHmacTokenSecret(t *testing.T) {
	secret := make([]byte, 16)

	_, err := rand.Read(secret)
	require.NoError(t, err, "Не удалось сгенерировать secret-ключ для JWT")

	config.Config.HmacTokenSecret = hex.EncodeToString(secret)
}

func GetUserIDWithToken() (string, uuid.UUID, error) {
	id := uuid.New()
	j := jwt.New(config.Config.HmacTokenSecret)

	tokenString, err := j.Encode(&models.User{ID: id})
	return tokenString, id, err
}
