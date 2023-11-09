package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	pgCategoryUser struct {
		db *sqlx.DB
	}

	User struct {
		ID uuid.UUID
	}
)

func (pgCU *pgCategoryUser) Create(user User) error {
	// TODO
	panic("todo")
}

func (pgCU *pgCategoryUser) GetByID(id uuid.UUID) (*User, error) {
	var user User
	err := pgCU.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)

	return &user, err
}

func (pgCU *pgCategoryUser) GetByLogin(login string) (*User, error) {
	var user User
	err := pgCU.db.Get(&user, "SELECT * FROM users WHERE login = $1", login)

	return &user, err
}

func (user *User) CheckPassword() bool {
	// todo
	panic("todo")
}
