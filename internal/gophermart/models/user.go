package models

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

func (user *User) EncryptPassword() {
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
