package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type pgCategoryUser struct {
	db postgres.PgSQL
}

func (pgCU *pgCategoryUser) Create(ctx context.Context, user *models.User) error {
	user.EncryptPassword()
	user.CreatedAt = time.Now()

	id, err := pgCU.db.ExecContextWithReturnID(
		ctx,
		"INSERT INTO users (login, password, balance, created_at) VALUES ($1, $2, $3, $4)",
		user.Login, user.Password, user.Balance, user.CreatedAt,
	)
	if err != nil {
		return err
	}

	user.ID, _ = uuid.Parse(id.(string))
	return nil
}

func (pgCU *pgCategoryUser) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := pgCU.db.GetContext(ctx, &user, "SELECT * FROM users WHERE lower(login) = lower($1)", login)

	return &user, err
}
