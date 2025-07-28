package infra

import (
	"database/sql"
	"github.com/cristianorosa/eSimulate/backend/domain"
)

// UserRepositoryPG implementa UserRepository usando PostgreSQL

// UserRepositoryPG implementa UserRepository usando PostgreSQL.
type UserRepositoryPG struct {
	DB *sql.DB
}

// Create cria um novo usuário no banco de dados.
func (r *UserRepositoryPG) Create(user *domain.User) error {
	query := `INSERT INTO users (name, email, password_hash, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.DB.QueryRow(query, user.Name, user.Email, user.PasswordHash, user.CreatedAt).Scan(&user.ID)
}

// FindByEmail busca um usuário pelo seu email.
func (r *UserRepositoryPG) FindByEmail(email string) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`
	row := r.DB.QueryRow(query, email)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
