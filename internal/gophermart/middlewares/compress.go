package middlewares

import (
	"compress/gzip"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer io.Writer
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w *gzipWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

func (ms *middlewaresService) Compress(ctx *gin.Context) {
	if !strings.Contains(ctx.GetHeader("Accept-Encoding"), "gzip") {
		return
	}

	gz, err := gzip.NewWriterLevel(ctx.Writer, gzip.BestSpeed)
	if err != nil {
		ms.logger.Error("failed to create writer with compression", zap.Error(err))
		return
	}
	defer gz.Close()

	ctx.Header("Content-Encoding", "gzip")
	ctx.Writer = &gzipWriter{ResponseWriter: ctx.Writer, writer: gz}

	ctx.Next()
}
