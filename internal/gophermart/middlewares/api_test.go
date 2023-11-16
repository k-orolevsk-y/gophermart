package middlewares

import (
	"testing"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/mocks"
)

func NewTestAPI(t *testing.T) *mocks.TestAPI {
	api := mocks.NewTestAPI(t)

	ConfigureMiddlewaresService(api)

	return api
}
