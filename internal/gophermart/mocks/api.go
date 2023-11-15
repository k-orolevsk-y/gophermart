package mocks

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
	orderpool "github.com/k-orolevsk-y/gophermart/internal/gophermart/order_pool"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/repository"
	"github.com/k-orolevsk-y/gophermart/pkg/router"
)

type TestAPI struct {
	t *testing.T

	router    *router.Router
	logger    *zap.Logger
	pg        repository.Repository
	orderPool *orderpool.OrderPool

	mockRepositoryCategoryUser         *MockRepositoryCategoryUser
	mockRepositoryCategoryOrders       *MockRepositoryCategoryOrders
	mockRepositoryCategoryUserWithdraw *MockRepositoryCategoryUserWithdraw
}

func NewTestAPI(t *testing.T) *TestAPI {
	gin.SetMode(gin.TestMode)

	if config.Config.HmacTokenSecret == "" {
		generateHmacTokenSecret(t)
	}
	if config.Config.WorkersLimit == 0 {
		config.Config.WorkersLimit = 1
	}

	api := &TestAPI{
		t:      t,
		router: router.New(),
		logger: zaptest.NewLogger(t, zaptest.Level(zapcore.PanicLevel)),
		pg:     NewMockRepository(gomock.NewController(t)),
	}

	api.configureRepository()
	api.orderPool = orderpool.NewOrderPool(api.GetLogger(), api.GetPg())

	return api
}

func (api *TestAPI) GetRouter() *router.Router {
	return api.router
}

func (api *TestAPI) GetLogger() *zap.Logger {
	return api.logger
}

func (api *TestAPI) GetPg() repository.Repository {
	return api.pg
}

func (api *TestAPI) GetOrderPool() *orderpool.OrderPool {
	return api.orderPool
}

func (api *TestAPI) GetPgEXPECT() *MockRepositoryMockRecorder {
	return api.pg.(*MockRepository).EXPECT()
}

func (api *TestAPI) GetPgUserEXPECT() *MockRepositoryCategoryUserMockRecorder {
	api.GetPgEXPECT().User().Return(api.mockRepositoryCategoryUser)
	return api.mockRepositoryCategoryUser.EXPECT()
}

func (api *TestAPI) GetPgOrderEXPECT() *MockRepositoryCategoryOrdersMockRecorder {
	api.GetPgEXPECT().Order().Return(api.mockRepositoryCategoryOrders)
	return api.mockRepositoryCategoryOrders.EXPECT()
}

func (api *TestAPI) GetPgUserWithdrawEXPECT() *MockRepositoryCategoryUserWithdrawMockRecorder {
	api.GetPgEXPECT().UserWithdraw().Return(api.mockRepositoryCategoryUserWithdraw)
	return api.mockRepositoryCategoryUserWithdraw.EXPECT()
}

func (api *TestAPI) configureRepository() {
	api.mockRepositoryCategoryUser = NewMockRepositoryCategoryUser(gomock.NewController(api.t))
	api.mockRepositoryCategoryOrders = NewMockRepositoryCategoryOrders(gomock.NewController(api.t))
	api.mockRepositoryCategoryUserWithdraw = NewMockRepositoryCategoryUserWithdraw(gomock.NewController(api.t))

	api.GetPgEXPECT().ParsePgError(gomock.Any()).DoAndReturn(func(err error) *pgconn.PgError {
		var pgError *pgconn.PgError
		if !errors.As(err, &pgError) {
			return &pgconn.PgError{}
		}

		return pgError
	}).AnyTimes()
}

func generateHmacTokenSecret(t *testing.T) {
	secret := make([]byte, 16)

	_, err := rand.Read(secret)
	require.NoError(t, err, "Не удалось сгенерировать secret-ключ для JWT")

	config.Config.HmacTokenSecret = hex.EncodeToString(secret)
}
