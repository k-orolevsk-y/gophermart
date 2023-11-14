package handlers

import (
	"testing"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/middlewares"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/mocks"
)

func NewTestAPI(t *testing.T) *mocks.TestAPI {
	api := mocks.NewTestAPI(t)

	ConfigureHandlersService(api)
	middlewares.ConfigureMiddlewaresService(api)

	return api
}
