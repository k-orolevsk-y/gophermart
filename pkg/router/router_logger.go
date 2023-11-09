package router

import (
	"strings"

	"go.uber.org/zap"
)

type RouterLogger struct {
	log *zap.Logger
}

func NewRouterLogger(logger *zap.Logger) *RouterLogger {
	return &RouterLogger{
		log: logger,
	}
}

func (logger *RouterLogger) Write(p []byte) (n int, err error) {
	logger.log.Info(strings.TrimSuffix(string(p), "\n"))
	return 0, nil
}
