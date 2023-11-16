package orderpool

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/mocks/repository"
)

type TestOrderPool struct {
	pool       *OrderPool
	repository *repository.MockRepository

	bufLogger *bytes.Buffer
}

func NewTestOrderPool(t *testing.T, accrualTestAddress string) *TestOrderPool {
	buf := new(bytes.Buffer)
	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.DebugLevel),
		zap.AddCaller(),
	)

	rep := repository.NewMockRepository(gomock.NewController(t))

	config.Config.WorkersLimit = 1
	config.Config.AccrualSystemAddress = accrualTestAddress

	pool := NewOrderPool(logger, rep)

	tp := &TestOrderPool{
		pool:       pool,
		repository: rep,

		bufLogger: buf,
	}
	tp.configureRepository(t)

	return tp
}

func (tp *TestOrderPool) Get() *OrderPool {
	return tp.pool
}

func (tp *TestOrderPool) GetBufLogger() *bytes.Buffer {
	return tp.bufLogger
}

func (tp *TestOrderPool) configureRepository(t *testing.T) {
	repositoryOrder := repository.NewMockRepositoryCategoryOrders(gomock.NewController(t))
	repositoryOrder.EXPECT().Edit(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	tp.repository.EXPECT().Order().Return(repositoryOrder).AnyTimes()
}
