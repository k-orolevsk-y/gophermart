package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/mocks"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/repository"
)

func TestHandlerGetBalance(t *testing.T) {
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

				api.GetPgOrderEXPECT().GetAccrualSumByUserID(gomock.Any(), userID).Return(float64(319), nil)
				api.GetPgUserWithdrawEXPECT().GetWithdrawnSumByUserID(gomock.Any(), userID).Return(float64(90), nil)

				return tokenString, nil
			},
			[]byte(`{"current":229,"withdrawn":90}`),
			http.StatusOK,
		},
		{
			"Positive/Float64",
			http.MethodGet,
			func(api *mocks.TestAPI) (string, error) {
				tokenString, userID, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				api.GetPgOrderEXPECT().GetAccrualSumByUserID(gomock.Any(), userID).Return(132.99, nil)
				api.GetPgUserWithdrawEXPECT().GetWithdrawnSumByUserID(gomock.Any(), userID).Return(84.845, nil)

				return tokenString, nil
			},
			[]byte(`{"current":48.14500000000001,"withdrawn":84.845}`),
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
			"Negative/RepositoryError/SumWithdraw",
			http.MethodGet,
			func(api *mocks.TestAPI) (string, error) {
				tokenString, userID, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				api.GetPgUserWithdrawEXPECT().GetWithdrawnSumByUserID(gomock.Any(), userID).Return(float64(0), errors.New("not connected"))

				return tokenString, nil
			},
			nil,
			http.StatusInternalServerError,
		},
		{
			"Negative/RepositoryError/SumCurrent",
			http.MethodGet,
			func(api *mocks.TestAPI) (string, error) {
				tokenString, userID, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				api.GetPgUserWithdrawEXPECT().GetWithdrawnSumByUserID(gomock.Any(), userID).Return(float64(0), nil)
				api.GetPgOrderEXPECT().GetAccrualSumByUserID(gomock.Any(), userID).Return(float64(0), errors.New("not connected"))

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

			req := httptest.NewRequest(tt.Method, "/api/user/balance/", nil)
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

func TestHandlerNewBalanceWithdrawn(t *testing.T) {
	tests := []struct {
		Name             string
		Method           string
		Body             []byte
		BeforeFunc       func(api *mocks.TestAPI) (tokenString string, err error)
		WantedStatusCode int
	}{
		{
			"Positive",
			http.MethodPost,
			[]byte(`{"order":"5081794355", "sum": 100}`),
			func(api *mocks.TestAPI) (string, error) {
				tokenString, _, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				api.GetPgUserWithdrawEXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

				return tokenString, nil
			},
			http.StatusOK,
		},
		{
			"Negative/WithoutAuthorization",
			http.MethodPost,
			[]byte(``),
			nil,
			http.StatusUnauthorized,
		},
		{
			"Negative/InvalidOrderNumber",
			http.MethodPost,
			[]byte(`{"order":"NotANumber", "sum": 300}`),
			func(api *mocks.TestAPI) (tokenString string, err error) {
				tokenString, _, err = GetUserIDWithToken()
				return
			},
			http.StatusBadRequest,
		},
		{
			"Negative/InvalidNumberOfAlgorithmLuna",
			http.MethodPost,
			[]byte(`{"order":"1234", "sum": 300}`),
			func(api *mocks.TestAPI) (tokenString string, err error) {
				tokenString, _, err = GetUserIDWithToken()
				return
			},
			http.StatusUnprocessableEntity,
		},
		{
			"Negative/InvalidSum",
			http.MethodPost,
			[]byte(`{"order":"1234", "sum": true}`),
			func(api *mocks.TestAPI) (tokenString string, err error) {
				tokenString, _, err = GetUserIDWithToken()
				return
			},
			http.StatusBadRequest,
		},
		{
			"Negative/SmallSum",
			http.MethodPost,
			[]byte(`{"order":"5081794355", "sum": 0.3}`),
			func(api *mocks.TestAPI) (tokenString string, err error) {
				tokenString, _, err = GetUserIDWithToken()
				return
			},
			http.StatusBadRequest,
		},
		{
			"Negative/UserHasNoFunds",
			http.MethodPost,
			[]byte(`{"order":"5081794355", "sum": 2}`),
			func(api *mocks.TestAPI) (string, error) {
				tokenString, _, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				api.GetPgUserWithdrawEXPECT().Create(gomock.Any(), gomock.Any()).Return(repository.ErrorInsufficientFunds)

				return tokenString, nil
			},
			http.StatusPaymentRequired,
		},
		{
			"Negative/RepositoryError",
			http.MethodPost,
			[]byte(`{"order":"5081794355", "sum": 20}`),
			func(api *mocks.TestAPI) (string, error) {
				tokenString, _, err := GetUserIDWithToken()
				if err != nil {
					return "", err
				}

				api.GetPgUserWithdrawEXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("not connected"))

				return tokenString, nil
			},
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

			req := httptest.NewRequest(tt.Method, "/api/user/balance/withdraw", bytes.NewReader(tt.Body))
			req.Header.Set("Authorization", tokenString)
			req.Header.Set("Content-Type", "application/json")

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
