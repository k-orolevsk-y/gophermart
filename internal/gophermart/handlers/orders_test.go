package handlers

import (
	"bytes"
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
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/mocks"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/repository"
)

func TestHandlerGetOrders(t *testing.T) {
	tests := []struct {
		Name             string
		Method           string
		BeforeFunc       func(api *mocks.TestAPI) (tokenString string, err error)
		WantedBody       []byte
		WantedStatusCode int
	}{
		{
			"Positive",
			http.MethodGet,
			func(api *mocks.TestAPI) (string, error) {
				tokenString, userID, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				orders := []models.Order{
					{
						ID:         100,
						UserID:     uuid.Nil,
						Status:     "NEW",
						Accrual:    proto.Float64(30.123),
						UploadedAt: time.Date(2023, time.November, 15, 15, 23, 30, 0, time.FixedZone("UTC", 0)),
					},
				}
				api.GetPgOrderEXPECT().GetAllByUserID(gomock.Any(), userID).Return(orders, nil)

				return tokenString, nil
			},
			[]byte(`[{"number":"100","status":"NEW","accrual":30.123,"uploaded_at":"2023-11-15T15:23:30Z"}]`),
			http.StatusOK,
		},
		{
			"Positive/WithoutItems",
			http.MethodGet,
			func(api *mocks.TestAPI) (string, error) {
				tokenString, userID, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				api.GetPgOrderEXPECT().GetAllByUserID(gomock.Any(), userID).Return([]models.Order{}, nil)

				return tokenString, nil
			},
			[]byte(`[]`),
			http.StatusOK,
		},
		{
			"Negative/WithoutAuthorization",
			http.MethodGet,
			nil,
			nil,
			http.StatusUnauthorized,
		},
		{
			"Negative/RepositoryError",
			http.MethodGet,
			func(api *mocks.TestAPI) (string, error) {
				tokenString, userID, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				api.GetPgOrderEXPECT().GetAllByUserID(gomock.Any(), userID).Return(nil, errors.New("not connected"))

				return tokenString, nil
			},
			nil,
			http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			api := NewTestAPI(t)

			var tokenString string
			if tt.BeforeFunc != nil {
				tknString, err := tt.BeforeFunc(api)
				require.NoError(t, err, "Ошибка при выполнении функции до теста")

				tokenString = tknString
			}

			w := httptest.NewRecorder()

			req := httptest.NewRequest(tt.Method, "/api/user/orders", nil)
			req.Header.Set("Authorization", tokenString)

			api.GetRouter().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err, "Не удалось прочитать ответ от запроса")

			reqSuccess := assert.Equal(t, tt.WantedStatusCode, res.StatusCode, "Статус код не соответствует ожидаемому")
			if tt.WantedBody != nil {
				reqSuccess = reqSuccess && assert.Equal(t, string(tt.WantedBody), string(body), "Тело ответа не соответствует ожидаемому")
			}

			if !reqSuccess {
				dump, _ := httputil.DumpRequest(req, true)
				t.Logf("\nЗапрос:\n%s\nОтвет:\n%s", dump, body)
			}
		})
	}
}

func TestHandlerNewOrder(t *testing.T) {
	tests := []struct {
		Name             string
		Method           string
		OrderID          []byte
		BeforeFunc       func(api *mocks.TestAPI) (tokenString string, err error)
		WantedStatusCode int
	}{
		{
			"Positive",
			http.MethodPost,
			[]byte(`5081794355`),
			func(api *mocks.TestAPI) (string, error) {
				tokenString, _, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				api.GetPgOrderEXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)

				return tokenString, nil
			},
			http.StatusAccepted,
		},
		{
			"Positive/AlreadyCreated",
			http.MethodPost,
			[]byte(`5081794355`),
			func(api *mocks.TestAPI) (string, error) {
				tokenString, _, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				api.GetPgOrderEXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&pgconn.PgError{Message: repository.ErrorOrderByThisUser})

				return tokenString, nil
			},
			http.StatusOK,
		},
		{
			"Negative/WithoutBody",
			http.MethodPost,
			[]byte(``),
			func(api *mocks.TestAPI) (tokenString string, err error) {
				tokenString, _, err = GetUserIDWithToken()
				return
			},
			http.StatusBadRequest,
		},
		{
			"Negative/InvalidOrderNumber",
			http.MethodPost,
			[]byte(`invalid_body`),
			func(api *mocks.TestAPI) (tokenString string, err error) {
				tokenString, _, err = GetUserIDWithToken()
				return
			},
			http.StatusBadRequest,
		},
		{
			"Negative/InvalidNumberOfAlgorithmLuna",
			http.MethodPost,
			[]byte(`1234`),
			func(api *mocks.TestAPI) (tokenString string, err error) {
				tokenString, _, err = GetUserIDWithToken()
				return
			},
			http.StatusUnprocessableEntity,
		},
		{
			"Negative/AlreadyCreatedByOtherUser",
			http.MethodPost,
			[]byte(`5081794355`),
			func(api *mocks.TestAPI) (string, error) {
				tokenString, _, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				api.GetPgOrderEXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&pgconn.PgError{Message: repository.ErrorOrderByOtherUser})

				return tokenString, nil
			},
			http.StatusConflict,
		},
		{
			"Negative/RepositoryError",
			http.MethodPost,
			[]byte(`5081794355`),
			func(api *mocks.TestAPI) (string, error) {
				tokenString, _, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				api.GetPgOrderEXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("not connected"))

				return tokenString, nil
			},
			http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			api := NewTestAPI(t)

			api.GetOrderPool().Run()
			defer api.GetOrderPool().Close()

			var tokenString string
			if tt.BeforeFunc != nil {
				tknString, err := tt.BeforeFunc(api)
				require.NoError(t, err, "Ошибка при выполнении функции до теста")

				tokenString = tknString
			}

			w := httptest.NewRecorder()

			req := httptest.NewRequest(tt.Method, "/api/user/orders", bytes.NewReader(tt.OrderID))
			req.Header.Set("Authorization", tokenString)
			req.Header.Set("Content-Type", "text/plain")

			api.GetRouter().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err, "Не удалось прочитать ответ от запроса")

			reqSuccess := assert.Equal(t, tt.WantedStatusCode, res.StatusCode, "Статус код не соответствует ожидаемому")
			if !reqSuccess {
				dump, _ := httputil.DumpRequest(req, true)
				t.Logf("\nЗапрос:\n%s\nОтвет:\n%s", dump, body)
			}
		})
	}
}
