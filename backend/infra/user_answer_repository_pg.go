package infra

import (
	"database/sql"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// UserAnswerRepositoryPG implementa UserAnswerRepository para PostgreSQL
type UserAnswerRepositoryPG struct {
	DB *sql.DB
}

// Create cria uma nova resposta de usuário
func (r *UserAnswerRepositoryPG) Create(userAnswer *domain.UserAnswer) error {
	query := `
		INSERT INTO user_answers (user_exam_id, question_id, option_id, is_correct, is_marked_for_review, answered_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := r.DB.QueryRow(
		query,
		userAnswer.UserExamID,
		userAnswer.QuestionID,
		userAnswer.OptionID,
		userAnswer.IsCorrect,
		userAnswer.IsMarkedForReview,
		userAnswer.AnsweredAt,
	).Scan(&userAnswer.ID)

	if err != nil {
		return err
	}

	return nil
}

// Update atualiza uma resposta de usuário
func (r *UserAnswerRepositoryPG) Update(userAnswer *domain.UserAnswer) error {
	query := `
		UPDATE user_answers 
		SET option_id = $1, is_correct = $2, is_marked_for_review = $3
		WHERE id = $4`

	result, err := r.DB.Exec(
		query,
		userAnswer.OptionID,
		userAnswer.IsCorrect,
		userAnswer.IsMarkedForReview,
		userAnswer.ID,
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

// ListByUserExam lista todas as respostas de um exame aplicado
func (r *UserAnswerRepositoryPG) ListByUserExam(userExamID int) ([]*domain.UserAnswer, error) {
	query := `
		SELECT id, user_exam_id, question_id, option_id, is_correct, is_marked_for_review, answered_at
		FROM user_answers
		WHERE user_exam_id = $1
		ORDER BY answered_at`

	rows, err := r.DB.Query(query, userExamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAnswers []*domain.UserAnswer
	for rows.Next() {
		userAnswer := &domain.UserAnswer{}
		var optionID sql.NullInt32
		var isCorrect sql.NullBool

		err := rows.Scan(
			&userAnswer.ID,
			&userAnswer.UserExamID,
			&userAnswer.QuestionID,
			&optionID,
			&isCorrect,
			&userAnswer.IsMarkedForReview,
			&userAnswer.AnsweredAt,
		)
		if err != nil {
			return nil, err
		}

		if optionID.Valid {
			optionIDInt := int(optionID.Int32)
			userAnswer.OptionID = &optionIDInt
		}
		if isCorrect.Valid {
			userAnswer.IsCorrect = &isCorrect.Bool
		}

		userAnswers = append(userAnswers, userAnswer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userAnswers, nil
}

// FindByQuestion busca uma resposta específica de uma questão
func (r *UserAnswerRepositoryPG) FindByQuestion(userExamID int, questionID int) (*domain.UserAnswer, error) {
	query := `
		SELECT id, user_exam_id, question_id, option_id, is_correct, is_marked_for_review, answered_at
		FROM user_answers
		WHERE user_exam_id = $1 AND question_id = $2`

	userAnswer := &domain.UserAnswer{}
	var optionID sql.NullInt32
	var isCorrect sql.NullBool

	err := r.DB.QueryRow(query, userExamID, questionID).Scan(
		&userAnswer.ID,
		&userAnswer.UserExamID,
		&userAnswer.QuestionID,
		&optionID,
		&isCorrect,
		&userAnswer.IsMarkedForReview,
		&userAnswer.AnsweredAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if optionID.Valid {
		optionIDInt := int(optionID.Int32)
		userAnswer.OptionID = &optionIDInt
	}
	if isCorrect.Valid {
		userAnswer.IsCorrect = &isCorrect.Bool
	}

	return userAnswer, nil
}
