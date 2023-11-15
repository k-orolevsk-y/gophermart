package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandlerNoRoute(t *testing.T) {
	api := NewTestAPI(t)

	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/api/undefined", nil)
	api.GetRouter().ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	reqSuccess := assert.Equal(t, http.StatusNotFound, res.StatusCode, "Статус код не соответствует ожидаемому")
	if !reqSuccess {
		dump, _ := httputil.DumpRequest(req, true)
		t.Logf("\nЗапрос:\n%s", dump)
	}
}

func TestHandlerNoMethod(t *testing.T) {
	api := NewTestAPI(t)

	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/api/user/register", nil)
	api.GetRouter().ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	reqSuccess := assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode, "Статус код не соответствует ожидаемому")
	if !reqSuccess {
		dump, _ := httputil.DumpRequest(req, true)
		t.Logf("\nЗапрос:\n%s", dump)
	}
}

func TestHandlerRecovery(t *testing.T) {
	gin.DefaultErrorWriter = nil
	api := NewTestAPI(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/panic", nil)

	api.GetRouter().GET("/api/panic", func(ctx *gin.Context) {
		panic("critical error!")
	})
	api.GetRouter().ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	reqSuccess := assert.Equal(t, http.StatusInternalServerError, res.StatusCode, "Статус код не соответствует ожидаемому")
	if !reqSuccess {
		dump, _ := httputil.DumpRequest(req, true)
		t.Logf("\nЗапрос:\n%s", dump)
	}
}
