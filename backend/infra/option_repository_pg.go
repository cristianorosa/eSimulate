package infra

import (
	"database/sql"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// OptionRepositoryPG implementa OptionRepository para PostgreSQL
type OptionRepositoryPG struct {
	DB *sql.DB
}

// Create cria uma nova opção
func (r *OptionRepositoryPG) Create(option *domain.Option) error {
	query := `INSERT INTO options (question_id, text, is_correct, explanation, order_index) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.DB.QueryRow(
		query,
		option.QuestionID,
		option.Text,
		option.IsCorrect,
		option.Explanation,
		option.OrderIndex,
	).Scan(&option.ID)

	if err != nil {
		return err
	}

	return nil
}

// Update atualiza uma opção existente
func (r *OptionRepositoryPG) Update(option *domain.Option) error {
	query := `
		UPDATE options 
		SET text = $1, is_correct = $2, explanation = $3, order_index = $4
		WHERE id = $5`

	result, err := r.DB.Exec(
		query,
		option.Text,
		option.IsCorrect,
		option.Explanation,
		option.OrderIndex,
		option.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Delete remove uma opção
func (r *OptionRepositoryPG) Delete(id int) error {
	query := `DELETE FROM options WHERE id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// DeleteByQuestionID remove todas as opções de uma questão
func (r *OptionRepositoryPG) DeleteByQuestionID(questionID int) error {
	query := `DELETE FROM options WHERE question_id = $1`

	_, err := r.DB.Exec(query, questionID)
	if err != nil {
		return err
	}

	return nil
}

// FindByQuestionID busca todas as opções de uma questão
func (r *OptionRepositoryPG) FindByQuestionID(questionID int) ([]*domain.Option, error) {
	query := `
		SELECT id, question_id, text, is_correct, explanation, order_index
		FROM options
		WHERE question_id = $1
		ORDER BY order_index`

	rows, err := r.DB.Query(query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []*domain.Option
	for rows.Next() {
		option := &domain.Option{}
		err := rows.Scan(
			&option.ID,
			&option.QuestionID,
			&option.Text,
			&option.IsCorrect,
			&option.Explanation,
			&option.OrderIndex,
		)
		if err != nil {
			return nil, err
		}
		options = append(options, option)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return options, nil
} 