package infra

import (
	"database/sql"
	"github.com/cristianorosa/eSimulate/backend/domain"
)

// UserQuizRepositoryPG implementa UserQuizRepository usando PostgreSQL

// UserQuizRepositoryPG implementa UserQuizRepository usando PostgreSQL.
type UserQuizRepositoryPG struct {
	DB *sql.DB
}

// Create cria um novo registro de simulado de usuário no banco de dados.
func (r *UserQuizRepositoryPG) Create(uq *domain.UserQuiz) error {
	query := `INSERT INTO user_quiz (user_id, quiz_id, started_at, finished_at) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.DB.QueryRow(query, uq.UserID, uq.QuizID, uq.StartedAt, uq.FinishedAt).Scan(&uq.ID)
}

// ListByUser lista o histórico de simulados de um usuário.
func (r *UserQuizRepositoryPG) ListByUser(userID int) ([]*domain.UserQuiz, error) {
	query := `SELECT id, user_id, quiz_id, started_at, finished_at FROM user_quiz WHERE user_id = $1 ORDER BY started_at DESC`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var history []*domain.UserQuiz
	for rows.Next() {
		uq := &domain.UserQuiz{}
		var finishedAt sql.NullTime
		if err := rows.Scan(&uq.ID, &uq.UserID, &uq.QuizID, &uq.StartedAt, &finishedAt); err != nil {
			return nil, err
		}
		if finishedAt.Valid {
			t := finishedAt.Time
			uq.FinishedAt = &t
		}
		history = append(history, uq)
	}
	return history, nil
}
