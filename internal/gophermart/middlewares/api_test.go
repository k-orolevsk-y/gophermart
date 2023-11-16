package middlewares

import (
	"testing"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/mocks/api"
)

func NewTestAPI(t *testing.T, logger *zap.Logger) *api.TestAPI {
	api := api.NewTestAPI(t)
	if logger != nil {
		api.SetNewLogger(logger)
	}

	ConfigureMiddlewaresService(api)

	return api
}
