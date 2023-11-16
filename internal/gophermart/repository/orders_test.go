//go:build integration

package repository

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
)

type OrdersTestSuite struct {
	suite.Suite
	db *pg

	rnd *rand.Rand
}

func (suite *OrdersTestSuite) SetupSuite() {
	require.NotEmpty(suite.T(), databaseDSN, "Отсутствует DSN для подключения к базе данных")

	db, err := NewTestPG(suite.T())
	require.NoError(suite.T(), err, "Не удалось подключиться к базе данных")

	suite.db = db
	suite.rnd = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
}

func (suite *OrdersTestSuite) TearDownSuite() {
	if suite.db == nil {
		return
	}

	assert.NoError(suite.T(), CloseTestPG(suite.db), "Не удалось очистить базу данных или закрыть с ней соединение")
}

func (suite *OrdersTestSuite) AfterTest(_, testName string) {
	if suite.db == nil {
		return
	}

	_, err := suite.db.db.ExecContext(context.Background(), "DELETE FROM orders;")
	require.NoErrorf(suite.T(), err, "Не удалось очистить таблицу после теста `%s`", testName)
}

func (suite *OrdersTestSuite) TestCreateOrder() {
	order := &models.Order{
		ID:         suite.rnd.Int63(),
		UserID:     uuid.New(),
		Status:     "NEW",
		Accrual:    proto.Float64(suite.rnd.Float64() * 1e3),
		UploadedAt: time.Now(),
	}

	err := suite.db.Order().Create(context.Background(), order)
	require.NoError(suite.T(), err, "Не удалось создать заказ")

	orders, err := suite.db.Order().GetAllByUserID(context.Background(), order.UserID)
	require.NoError(suite.T(), err, "Не удалось получить заказы пользователя")

	require.Len(suite.T(), orders, 1, "Получено неверное количество заказов")
	require.Equal(suite.T(), order.ID, orders[0].ID, "Найденный заказ в базе неверный")
}

func (suite *OrdersTestSuite) TestEditOrder() {
	order := &models.Order{
		ID:         suite.rnd.Int63(),
		UserID:     uuid.New(),
		Status:     "NEW",
		Accrual:    nil,
		UploadedAt: time.Now(),
	}

	err := suite.db.Order().Create(context.Background(), order)
	require.NoError(suite.T(), err, "Не удалось создать заказ")

	order.Status = "PROCESSED"
	order.Accrual = proto.Float64(suite.rnd.Float64() * 1e3)

	err = suite.db.Order().Edit(context.Background(), order)
	require.NoError(suite.T(), err, "Не удалось отредактировать заказ")

	orders, err := suite.db.Order().GetAllByUserID(context.Background(), order.UserID)
	require.NoError(suite.T(), err, "Не удалось получить заказы пользователя")

	require.Len(suite.T(), orders, 1, "Получено неверное количество заказов")
	require.Equal(suite.T(), order.ID, orders[0].ID, "Найденный заказ в базе неверный")
	require.Equal(suite.T(), order.Status, orders[0].Status, "Неверный статус заказа")
	require.Equal(suite.T(), *order.Accrual, *orders[0].Accrual, "Неверное количество баллов")
}

func (suite *OrdersTestSuite) TestGetUserOrders() {
	var orderIDs []int64
	userID := uuid.New()

	for i := 0; i < suite.rnd.Intn(7)+3; i++ {
		order := &models.Order{
			ID:         suite.rnd.Int63(),
			UserID:     userID,
			Status:     "NEW",
			Accrual:    proto.Float64(suite.rnd.Float64() * 1e3),
			UploadedAt: time.Now(),
		}

		err := suite.db.Order().Create(context.Background(), order)
		require.NoErrorf(suite.T(), err, "Не удалось создать заказ #%d", i+1)

		orderIDs = append(orderIDs, order.ID)
	}

	orders, err := suite.db.Order().GetAllByUserID(context.Background(), userID)
	require.NoError(suite.T(), err, "Не удалось получить заказы пользователя")
	require.Len(suite.T(), orders, len(orderIDs), "Получено неверное количество заказов")

	var existsOrdersIDs []int64
	for _, order := range orders {
		existsOrdersIDs = append(existsOrdersIDs, order.ID)
	}

	require.EqualValues(suite.T(), orderIDs, existsOrdersIDs, "Полученные ID из базы неверны тем, которые были добавлены")
}

func (suite *OrdersTestSuite) TestGetAccrualSum() {
	var sum float64
	userID := uuid.New()

	for i := 0; i < suite.rnd.Intn(7)+3; i++ {
		accrual := suite.rnd.Float64() * 1e3

		err := suite.db.Order().Create(context.Background(), &models.Order{
			ID:         suite.rnd.Int63(),
			UserID:     userID,
			Status:     "NEW",
			Accrual:    &accrual,
			UploadedAt: time.Now(),
		})
		require.NoErrorf(suite.T(), err, "Не удалось создать заказ #%d", i+1)

		sum += accrual
	}

	existsSum, err := suite.db.Order().GetAccrualSumByUserID(context.Background(), userID)
	require.NoError(suite.T(), err, "Не удалось получить сумму баллов")
	require.Equal(suite.T(), sum, existsSum, "Неверная сумма баллов")
}

func TestOrdersTestSuite(t *testing.T) {
	suite.Run(t, new(OrdersTestSuite))
}
