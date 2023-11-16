package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestMiddlewareLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.InfoLevel),
		zap.AddCaller(),
	)

	api := NewTestAPI(t, logger)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	api.GetRouter().ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	reqSuccess := assert.Contains(t, buf.String(), "request", "Отсутствует лог о выполненном запросе")
	if !reqSuccess {
		dump, _ := httputil.DumpRequest(req, true)
		t.Logf("\nЗапрос:\n%s", dump)
	}
}
