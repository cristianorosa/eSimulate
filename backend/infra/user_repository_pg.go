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
	query := `INSERT INTO users (name, email, password_hash, role_id, created_at, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`
	return r.DB.QueryRow(query, user.Name, user.Email, user.PasswordHash, user.RoleID).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// FindByEmail busca um usuário pelo seu email.
func (r *UserRepositoryPG) FindByEmail(email string) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, role_id, created_at, updated_at FROM users WHERE email = $1`
	row := r.DB.QueryRow(query, email)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.RoleID, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

// ListAll retorna todos os usuários cadastrados.
func (r *UserRepositoryPG) ListAll() ([]*domain.User, error) {
	query := `SELECT id, name, email, password_hash, role_id, created_at, updated_at FROM users ORDER BY created_at DESC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.RoleID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}
