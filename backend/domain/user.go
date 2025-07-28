package domain

import (
	"errors"
	"time"
)

// User representa um usuário do sistema
// Dados pessoais devem ser tratados conforme a LGPD
// Armazenar apenas o necessário para funcionamento e segurança
// Permitir exclusão e anonimização conforme solicitado

type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	GoogleID     *string
	FacebookID   *string
	CreatedAt    time.Time
}

var ErrUserExists = errors.New("user already exists")
