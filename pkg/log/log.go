package log

import (
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
)

func New() (*zap.Logger, error) {
	if config.Config.ProductionMode {
		return zap.NewProduction()
	}

	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	l.Core().Enabled(zap.DebugLevel)

	return l, nil
}
