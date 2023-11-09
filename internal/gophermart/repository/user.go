package repository

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/gophermart/pkg/database/postgres"
)

type (
	pgCategoryUser struct {
		db postgres.PgSQL
	}

	User struct {
		ID        uuid.UUID `db:"id"`
		Login     string    `db:"login"`
		Password  string    `db:"password"`
		Balance   int       `db:"balance"`
		CreatedAt time.Time `db:"created_at"`
	}
)

func (pgCU *pgCategoryUser) Create(ctx context.Context, user *User) error {
	user.encryptPassword()
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

func (pgCU *pgCategoryUser) GetByLogin(ctx context.Context, login string) (*User, error) {
	var user User
	err := pgCU.db.GetContext(ctx, &user, "SELECT * FROM users WHERE lower(login) = lower($1)", login)

	return &user, err
}

func (user *User) encryptPassword() {
	h := sha256.New()

	h.Write([]byte(user.Password))
	dst := h.Sum(nil)

	user.Password = base64.StdEncoding.EncodeToString(dst)
}

func (user *User) CheckPassword(password string) bool {
	bytesPassword, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		return false
	}

	sum := sha256.Sum256([]byte(password))

	return bytes.Equal(bytesPassword, sum[:])
}
