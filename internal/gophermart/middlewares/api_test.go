package middlewares

import (
	"testing"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/mocks"
)

func NewTestAPI(t *testing.T, logger *zap.Logger) *mocks.TestAPI {
	api := mocks.NewTestAPI(t)
	if logger != nil {
		api.SetNewLogger(logger)
	}

	ConfigureMiddlewaresService(api)

	return api
}
