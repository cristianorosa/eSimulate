package infra

import (
	"database/sql"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// DomainPerformanceRepositoryPG implementa PerformanceRepository para PostgreSQL
type DomainPerformanceRepositoryPG struct {
	DB *sql.DB
}

// Create cria um novo registro de desempenho por domínio
func (r *DomainPerformanceRepositoryPG) Create(performance *domain.Performance) error {
	query := `
		INSERT INTO domain_performance (user_exam_id, domain_id, questions_answered, correct_answers, score_percentage, needs_improvement)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := r.DB.QueryRow(
		query,
		performance.UserExamID,
		performance.DomainID,
		performance.QuestionsAnswered,
		performance.CorrectAnswers,
		performance.ScorePercentage,
		performance.NeedsImprovement,
	).Scan(&performance.ID)

	if err != nil {
		return err
	}

	return nil
}

// Update atualiza um registro de desempenho por domínio
func (r *DomainPerformanceRepositoryPG) Update(performance *domain.Performance) error {
	query := `
		UPDATE domain_performance 
		SET questions_answered = $1, correct_answers = $2, score_percentage = $3, needs_improvement = $4
		WHERE id = $5`

	result, err := r.DB.Exec(
		query,
		performance.QuestionsAnswered,
		performance.CorrectAnswers,
		performance.ScorePercentage,
		performance.NeedsImprovement,
		performance.ID,
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

// ListByUserExam lista todos os desempenhos por domínio de um exame aplicado
func (r *DomainPerformanceRepositoryPG) ListByUserExam(userExamID int) ([]*domain.Performance, error) {
	query := `
		SELECT id, user_exam_id, domain_id, questions_answered, correct_answers, score_percentage, needs_improvement
		FROM domain_performance
		WHERE user_exam_id = $1
		ORDER BY domain_id`

	rows, err := r.DB.Query(query, userExamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []*domain.Performance
	for rows.Next() {
		performance := &domain.Performance{}
		err := rows.Scan(
			&performance.ID,
			&performance.UserExamID,
			&performance.DomainID,
			&performance.QuestionsAnswered,
			&performance.CorrectAnswers,
			&performance.ScorePercentage,
			&performance.NeedsImprovement,
		)
		if err != nil {
			return nil, err
		}
		performances = append(performances, performance)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return performances, nil
}
