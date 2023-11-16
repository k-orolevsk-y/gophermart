//go:build integration

package repository

import (
	"context"
	"encoding/hex"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
)

type UserTestSuite struct {
	suite.Suite
	db *pg

	rnd *rand.Rand
}

func (suite *UserTestSuite) SetupSuite() {
	require.NotEmpty(suite.T(), databaseDSN, "Отсутствует DSN для подключения к базе данных")

	db, err := NewTestPG(suite.T())
	require.NoError(suite.T(), err, "Не удалось подключиться к базе данных")

	suite.db = db
	suite.rnd = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
}

func (suite *UserTestSuite) TearDownSuite() {
	if suite.db == nil {
		return
	}

	assert.NoError(suite.T(), CloseTestPG(suite.db), "Не удалось очистить базу данных или закрыть с ней соединение")
}

func (suite *UserTestSuite) AfterTest(_, testName string) {
	if suite.db == nil {
		return
	}

	_, err := suite.db.db.ExecContext(context.Background(), "DELETE FROM users;")
	require.NoErrorf(suite.T(), err, "Не удалось очистить таблицу после теста `%s`", testName)
}

func (suite *UserTestSuite) TestUsers() {
	login := suite.generateString(30)
	pass := suite.generateString(30)

	user := &models.User{
		Login:    login,
		Password: pass,
	}

	err := suite.db.User().Create(context.Background(), user)
	require.NoError(suite.T(), err, "Не удалось создать пользователя")

	userByLogin, err := suite.db.User().GetByLogin(context.Background(), login)
	require.NoError(suite.T(), err, "Не удалось получить пользователя по логину")
	require.EqualValues(suite.T(), user.ID, userByLogin.ID, "ID пользователя полученного по логину не совпадают")

	userByID, err := suite.db.User().GetByID(context.Background(), user.ID)
	require.NoError(suite.T(), err, "Не удалось получить пользователя по ID")
	require.Equal(suite.T(), user.ID, userByID.ID, "ID пользователя полученного по ID не совпадают")
}

func (suite *UserTestSuite) generateString(n int) string {
	bs := make([]byte, n)
	suite.rnd.Read(bs)

	return hex.EncodeToString(bs)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
