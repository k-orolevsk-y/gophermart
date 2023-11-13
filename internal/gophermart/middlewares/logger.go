package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ms *middlewaresService) Logger(ctx *gin.Context) {
	start := time.Now()

	ctx.Next()

	uri := ctx.Request.URL
	size := ctx.Writer.Size()
	method := ctx.Request.Method
	duration := time.Since(start)
	statusCode := ctx.Writer.Status()
	userAgent := ctx.Request.UserAgent()

	ms.logger.Info(
		"Request",
		zap.Stringer("uri", uri),
		zap.String("method", method),
		zap.Duration("duration", duration),
		zap.String("userAgent", userAgent),
		zap.Int("statusCode", statusCode),
		zap.Int("size", size),
	)
}
