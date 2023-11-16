package middlewares

import (
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddlewareCompress(t *testing.T) {
	api := NewTestAPI(t, nil)

	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	api.GetRouter().ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	contentEncoding := res.Header.Get("Content-Encoding")
	_, err := gzip.NewReader(res.Body)

	reqSuccess := assert.Contains(t, contentEncoding, "gzip", "Отсутствует header о сжатии через gzip")
	reqSuccess = reqSuccess && assert.NoError(t, err, "Не удалось расшифровать тело ответа сжатое через gzip")

	if !reqSuccess {
		dump, _ := httputil.DumpRequest(req, true)
		t.Logf("\nЗапрос:\n%s", dump)
	}
}
