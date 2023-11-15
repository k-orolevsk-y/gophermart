package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/mocks"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
)

func TestHandlerGetWithdrawals(t *testing.T) {
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

				userWithdrawals := []models.UserWithdraw{
					{
						ID:          uuid.Nil,
						UserID:      uuid.Nil,
						OrderID:     5081794355,
						Sum:         34.99,
						ProcessedAt: time.Date(2023, time.November, 15, 15, 23, 30, 0, time.FixedZone("UTC", 0)),
					},
				}

				api.GetPgUserWithdrawEXPECT().GetAllByUserID(gomock.Any(), userID).Return(userWithdrawals, nil)

				return tokenString, nil
			},
			[]byte(`[{"order":"5081794355","sum":34.99,"processed_at":"2023-11-15T15:23:30Z"}]`),
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

				api.GetPgUserWithdrawEXPECT().GetAllByUserID(gomock.Any(), userID).Return([]models.UserWithdraw{}, nil)

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

				api.GetPgUserWithdrawEXPECT().GetAllByUserID(gomock.Any(), userID).Return(nil, errors.New("not connected"))

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

			req := httptest.NewRequest(tt.Method, "/api/user/withdrawals", nil)
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
