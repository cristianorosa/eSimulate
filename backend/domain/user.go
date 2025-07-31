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
	RoleID       int
	GoogleID     *string
	FacebookID   *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Role representa um papel/role do sistema
type Role struct {
	ID   int
	Name string
}

// ErrUserExists indica que um usuário com o email fornecido já existe.
var ErrUserExists = errors.New("user already exists")

// ErrInvalidInput indica que os dados de entrada são inválidos
var ErrInvalidInput = errors.New("invalid input")

// ErrNotFound indica que o recurso não foi encontrado
var ErrNotFound = errors.New("not found")

// Constantes para IDs dos roles (baseado nos inserts do schema)
const (
	RoleIDUser    = 1
	RoleIDRedator = 2
	RoleIDAdmin   = 3
)

// IsAdmin verifica se o usuário é administrador
func (u *User) IsAdmin() bool {
	return u.RoleID == RoleIDAdmin
}

// IsRedator verifica se o usuário é redator
func (u *User) IsRedator() bool {
	return u.RoleID == RoleIDRedator
}

// CanCreateContent verifica se o usuário pode criar conteúdo (redator ou admin)
func (u *User) CanCreateContent() bool {
	return u.RoleID == RoleIDRedator || u.RoleID == RoleIDAdmin
}
