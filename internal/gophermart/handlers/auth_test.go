package handlers

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/mocks"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
)

func TestHandlerRegister(t *testing.T) {
	tests := []struct {
		Name             string
		Method           string
		Body             []byte
		RepositoryFunc   func(api *mocks.TestAPI)
		WantedAuthHeader bool
		WantedStatusCode int
	}{
		{
			"Positive",
			http.MethodPost,
			[]byte(`{"login":"userWantRegister","password":"strongPassword"}`),
			func(api *mocks.TestAPI) {
				api.GetPgUserEXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ interface{}, user *models.User) error {
						user.ID = uuid.New()
						user.CreatedAt = time.Now()
						user.EncryptPassword()

						return nil
					})
			},
			true,
			http.StatusOK,
		},
		{
			"Negative/BadRequest/Login",
			http.MethodPost,
			[]byte(`{"login":"","password":"superPassword"}`),
			nil,
			false,
			http.StatusBadRequest,
		},
		{
			"Negative/BadRequest/Password",
			http.MethodPost,
			[]byte(`{"login":"rootUser_kk","password":"STRNG"}`),
			nil,
			false,
			http.StatusBadRequest,
		},
		{
			"Negative/User already created",
			http.MethodPost,
			[]byte(`{"login":"userAlreadyCreated","password":"strongPassword"}`),
			func(api *mocks.TestAPI) {
				api.GetPgUserEXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&pgconn.PgError{Code: "23000"})
			},
			false,
			http.StatusConflict,
		},
		{
			"Negative/Repository error",
			http.MethodPost,
			[]byte(`{"login":"loser","password":"strongPassword"}`),
			func(api *mocks.TestAPI) {
				api.GetPgUserEXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("not connected"))
			},
			false,
			http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			api := NewTestAPI(t)

			if tt.RepositoryFunc != nil {
				tt.RepositoryFunc(api)
			}

			w := httptest.NewRecorder()

			req := httptest.NewRequest(tt.Method, "/api/user/register", bytes.NewReader(tt.Body))
			req.Header.Set("Content-Type", "application/json")

			api.GetRouter().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			reqSuccess := assert.Equal(t, tt.WantedStatusCode, res.StatusCode, "Статус код не соответствует ожидаемому")
			if tt.WantedAuthHeader {
				reqSuccess = reqSuccess && assert.NotEmpty(t, res.Header.Get("Authorization"), "В ответе отсутствует JWT токен")
			}

			if !reqSuccess {
				dump, _ := httputil.DumpRequest(req, true)
				body, _ := io.ReadAll(res.Body)

				t.Logf("\nЗапрос:\n%sОтвет:\n%s", dump, body)
			}
		})
	}
}

func TestHandlerLogin(t *testing.T) {
	tests := []struct {
		Name             string
		Method           string
		Body             []byte
		RepositoryFunc   func(api *mocks.TestAPI)
		WantedAuthHeader bool
		WantedStatusCode int
	}{
		{
			"Positive",
			http.MethodPost,
			[]byte(`{"login":"user-123","password":"strongPassword"}`),
			func(api *mocks.TestAPI) {
				user := models.User{
					ID:        uuid.New(),
					Login:     "user-123",
					Password:  "strongPassword",
					CreatedAt: time.Now(),
				}
				user.EncryptPassword()

				api.GetPgUserEXPECT().GetByLogin(gomock.Any(), "user-123").Return(&user, nil)
			},
			true,
			http.StatusOK,
		},
		{
			"Negative/BadRequest/Login",
			http.MethodPost,
			[]byte(`{"login":"","password":"superPassword"}`),
			nil,
			false,
			http.StatusBadRequest,
		},
		{
			"Negative/BadRequest/Password",
			http.MethodPost,
			[]byte(`{"login":"rootUser_kk","password":"STRNG"}`),
			nil,
			false,
			http.StatusBadRequest,
		},
		{
			"Negative/Invalid/Login",
			http.MethodPost,
			[]byte(`{"login":"userNotExist","password":"strongPassword"}`),
			func(api *mocks.TestAPI) {
				api.GetPgUserEXPECT().GetByLogin(gomock.Any(), "userNotExist").Return(nil, sql.ErrNoRows)
			},
			false,
			http.StatusUnauthorized,
		},
		{
			"Negative/Invalid/Password",
			http.MethodPost,
			[]byte(`{"login":"geniusUser","password":"strongPassword"}`),
			func(api *mocks.TestAPI) {
				user := models.User{
					ID:        uuid.New(),
					Login:     "geniusUser",
					Password:  "otherStrongPassword",
					CreatedAt: time.Now(),
				}
				user.EncryptPassword()

				api.GetPgUserEXPECT().GetByLogin(gomock.Any(), "geniusUser").Return(&user, nil)
			},
			false,
			http.StatusUnauthorized,
		},
		{
			"Negative/Repository error",
			http.MethodPost,
			[]byte(`{"login":"errorUser","password":"strongPassword"}`),
			func(api *mocks.TestAPI) {
				api.GetPgUserEXPECT().
					GetByLogin(gomock.Any(), "errorUser").
					Return(nil, errors.New("not connected"))
			},
			false,
			http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			api := NewTestAPI(t)

			if tt.RepositoryFunc != nil {
				tt.RepositoryFunc(api)
			}

			w := httptest.NewRecorder()

			req := httptest.NewRequest(tt.Method, "/api/user/login", bytes.NewReader(tt.Body))
			req.Header.Set("Content-Type", "application/json")

			api.GetRouter().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			reqSuccess := assert.Equal(t, tt.WantedStatusCode, res.StatusCode, "Статус код не соответствует ожидаемому")
			if tt.WantedAuthHeader {
				reqSuccess = reqSuccess && assert.NotEmpty(t, res.Header.Get("Authorization"), "В ответе отсутствует JWT токен")
			}

			if !reqSuccess {
				dump, _ := httputil.DumpRequest(req, true)
				body, _ := io.ReadAll(res.Body)

				t.Logf("\nЗапрос:\n%sОтвет:\n%s", dump, body)
			}
		})
	}
}
