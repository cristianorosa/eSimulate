package domain

import (
	"errors"
	"time"
)

// User representa um usuário do sistema.
type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	GoogleID     *string
	FacebookID   *string
	CreatedAt    time.Time
}

// ErrUserExists indica que um usuário com o email fornecido já existe.
var ErrUserExists = errors.New("user already exists")
