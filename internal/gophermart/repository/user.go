package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type (
	RepositoryCategoryUser interface {
		Create(ctx context.Context, user *models.User) error
		GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
		GetByLogin(ctx context.Context, login string) (*models.User, error)
	}

	pgCategoryUser struct {
		db postgres.PgSQL
	}
)

func (pgCU *pgCategoryUser) Create(ctx context.Context, user *models.User) error {
	user.EncryptPassword()
	user.CreatedAt = time.Now()

	id, err := pgCU.db.ExecContextWithReturnID(
		ctx,
		"INSERT INTO users (login, password, created_at) VALUES ($1, $2, $3)",
		user.Login, user.Password, user.CreatedAt,
	)
	if err != nil {
		return err
	}

	user.ID, _ = uuid.Parse(id.(string))
	return nil
}

func (pgCU *pgCategoryUser) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := pgCU.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)

	return &user, err
}

func (pgCU *pgCategoryUser) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := pgCU.db.GetContext(ctx, &user, "SELECT * FROM users WHERE lower(login) = lower($1)", login)

	return &user, err
}
