package infra

import (
	"database/sql"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// TopicPerformanceRepositoryPG implementa TopicPerformanceRepository para PostgreSQL
type TopicPerformanceRepositoryPG struct {
	DB *sql.DB
}

// Create cria um novo registro de desempenho por tópico
func (r *TopicPerformanceRepositoryPG) Create(performance *domain.TopicPerformance) error {
	query := `
		INSERT INTO topic_performance (user_exam_id, topic_id, correct_answers, total_questions, score)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := r.DB.QueryRow(
		query,
		performance.UserExamID,
		performance.TopicID,
		performance.CorrectAnswers,
		performance.TotalQuestions,
		performance.Score,
	).Scan(&performance.ID)

	if err != nil {
		return err
	}

	return nil
}

// Update atualiza um registro de desempenho por tópico
func (r *TopicPerformanceRepositoryPG) Update(performance *domain.TopicPerformance) error {
	query := `
		UPDATE topic_performance 
		SET correct_answers = $1, total_questions = $2, score = $3
		WHERE id = $4`

	result, err := r.DB.Exec(
		query,
		performance.CorrectAnswers,
		performance.TotalQuestions,
		performance.Score,
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

// FindByUserExam busca desempenho por tópico de um exame do usuário
func (r *TopicPerformanceRepositoryPG) FindByUserExam(userExamID int) ([]*domain.TopicPerformance, error) {
	query := `
		SELECT id, user_exam_id, topic_id, correct_answers, total_questions, score
		FROM topic_performance
		WHERE user_exam_id = $1`

	rows, err := r.DB.Query(query, userExamID)
	if err != nil {
		return []*domain.TopicPerformance{}, err
	}
	defer rows.Close()

	var performances []*domain.TopicPerformance
	for rows.Next() {
		p := &domain.TopicPerformance{}
		err := rows.Scan(
			&p.ID,
			&p.UserExamID,
			&p.TopicID,
			&p.CorrectAnswers,
			&p.TotalQuestions,
			&p.Score,
		)
		if err != nil {
			return []*domain.TopicPerformance{}, err
		}
		performances = append(performances, p)
	}

	if err = rows.Err(); err != nil {
		return []*domain.TopicPerformance{}, err
	}

	// Garantir que sempre retorne um array, mesmo que vazio
	if performances == nil {
		performances = []*domain.TopicPerformance{}
	}

	return performances, nil
}
