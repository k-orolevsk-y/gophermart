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

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
)

type UserWithdrawTestSuite struct {
	suite.Suite
	db *pg

	rnd *rand.Rand
}

func (suite *UserWithdrawTestSuite) SetupSuite() {
	require.NotEmpty(suite.T(), databaseDSN, "Отсутствует DSN для подключения к базе данных")

	db, err := NewTestPG(suite.T())
	require.NoError(suite.T(), err, "Не удалось подключиться к базе данных")

	suite.db = db
	suite.rnd = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
}

func (suite *UserWithdrawTestSuite) TearDownSuite() {
	if suite.db == nil {
		return
	}

	assert.NoError(suite.T(), CloseTestPG(suite.db), "Не удалось очистить базу данных или закрыть с ней соединение")
}

func (suite *UserWithdrawTestSuite) AfterTest(_, testName string) {
	if suite.db == nil {
		return
	}

	_, err := suite.db.db.ExecContext(context.Background(), "DELETE FROM users_withdrawals;")
	require.NoErrorf(suite.T(), err, "Не удалось очистить таблицу после теста `%s`", testName)
}

func (suite *UserWithdrawTestSuite) TestCreateUserWithdraw() {
	userID := uuid.New()

	err := suite.createTestOrder(userID, 1e10)
	require.NoError(suite.T(), err, "Не удалось создать тестовый заказ с начисленными баллами")

	userWithdraw := &models.UserWithdraw{
		UserID:  userID,
		OrderID: suite.rnd.Int63(),
		Sum:     suite.rnd.Float64(),
	}

	err = suite.db.UserWithdraw().Create(context.Background(), userWithdraw)
	require.NoError(suite.T(), err, "Не удалось создать объект списание баллов")

	userWithdrawals, err := suite.db.UserWithdraw().GetAllByUserID(context.Background(), userWithdraw.UserID)
	require.NoError(suite.T(), err, "Не удалось получить объект списания баллов")
	require.Len(suite.T(), userWithdrawals, 1, "Получено неверное количество объектов списания баллов")
	require.EqualValues(suite.T(), userWithdraw.ID, userWithdrawals[0].ID, "ID объекта не совпадает")
}

func (suite *UserWithdrawTestSuite) TestNegativeCreateUserWithdraw() {
	userWithdraw := &models.UserWithdraw{
		UserID:  uuid.New(),
		OrderID: suite.rnd.Int63(),
		Sum:     suite.rnd.Float64(),
	}

	err := suite.db.UserWithdraw().Create(context.Background(), userWithdraw)
	require.Error(suite.T(), err, "Создался объект списания баллов, хотя у пользователя их нет")
}

func (suite *UserWithdrawTestSuite) TestGetWithdrawnSum() {
	userID := uuid.New()

	err := suite.createTestOrder(userID, 1e10)
	require.NoError(suite.T(), err, "Не удалось создать тестовый заказ с начисленными баллами")

	var sum float64
	for i := 0; i < suite.rnd.Intn(7)+3; i++ {
		sumWithdrawn := suite.rnd.Float64() * 1e3

		err = suite.db.UserWithdraw().Create(context.Background(), &models.UserWithdraw{
			UserID:  userID,
			OrderID: suite.rnd.Int63(),
			Sum:     sumWithdrawn,
		})
		require.NoErrorf(suite.T(), err, "Не удалось создать объект списания баллов #%d", i+1)

		sum += sumWithdrawn
	}

	existsSum, err := suite.db.UserWithdraw().GetWithdrawnSumByUserID(context.Background(), userID)
	require.NoError(suite.T(), err, "Не удалось получить сумму списанных баллов")
	require.Equal(suite.T(), sum, existsSum, "Неверная сумма списанных баллов")
}

func (suite *UserWithdrawTestSuite) createTestOrder(userID uuid.UUID, sum float64) error {
	return suite.db.Order().Create(context.Background(), &models.Order{
		ID:      suite.rnd.Int63(),
		UserID:  userID,
		Status:  "PROCESSED",
		Accrual: &sum,
	})
}

func TestUserWithdrawTestSuite(t *testing.T) {
	suite.Run(t, new(UserWithdrawTestSuite))
}
