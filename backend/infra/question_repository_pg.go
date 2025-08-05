package infra

import (
	"database/sql"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// QuestionRepositoryPG implementa QuestionRepository para PostgreSQL
type QuestionRepositoryPG struct {
	DB *sql.DB
}

// Create cria uma nova questão
func (r *QuestionRepositoryPG) Create(question *domain.Question) error {
	query := `INSERT INTO questions (topic_id, statement, problem, content_type, explanation, question_type, difficulty_level, created_by, is_active, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP) RETURNING id, created_at`

	err := r.DB.QueryRow(
		query,
		question.TopicID,
		question.Statement,
		question.Problem,
		question.ContentType,
		question.Explanation,
		question.QuestionType,
		question.DifficultyLevel,
		question.CreatedBy,
		question.IsActive,
	).Scan(&question.ID, &question.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Update atualiza uma questão existente
func (r *QuestionRepositoryPG) Update(question *domain.Question) error {
	query := `
		UPDATE questions 
		SET topic_id = $1, statement = $2, problem = $3, content_type = $4, explanation = $5, question_type = $6, difficulty_level = $7, is_active = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9`

	result, err := r.DB.Exec(
		query,
		question.TopicID,
		question.Statement,
		question.Problem,
		question.ContentType,
		question.Explanation,
		question.QuestionType,
		question.DifficultyLevel,
		question.IsActive,
		question.ID,
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

// Delete remove uma questão
func (r *QuestionRepositoryPG) Delete(id int) error {
	query := `DELETE FROM questions WHERE id = $1`

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

// FindByID busca uma questão por ID
func (r *QuestionRepositoryPG) FindByID(id int) (*domain.Question, error) {
	query := `
		SELECT id, topic_id, statement, problem, content_type, explanation, question_type, difficulty_level, created_by, is_active, created_at, updated_at
		FROM questions
		WHERE id = $1`

	q := &domain.Question{}
	err := r.DB.QueryRow(query, id).Scan(
		&q.ID,
		&q.TopicID,
		&q.Statement,
		&q.Problem,
		&q.ContentType,
		&q.Explanation,
		&q.QuestionType,
		&q.DifficultyLevel,
		&q.CreatedBy,
		&q.IsActive,
		&q.CreatedAt,
		&q.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return q, nil
}

// ListByExam lista todas as questões de um exame (via tópicos)
func (r *QuestionRepositoryPG) ListByExam(examID int) ([]*domain.Question, error) {
	query := `
		SELECT q.id, q.topic_id, q.statement, q.problem, q.content_type, q.explanation, q.question_type, q.difficulty_level, q.created_by, q.is_active, q.created_at, q.updated_at
		FROM questions q
		JOIN topics t ON q.topic_id = t.id
		WHERE t.exam_id = $1 AND q.is_active = true
		ORDER BY q.created_at DESC`

	rows, err := r.DB.Query(query, examID)
	if err != nil {
		return []*domain.Question{}, err
	}
	defer rows.Close()

	var questions []*domain.Question
	for rows.Next() {
		q := &domain.Question{}
		err := rows.Scan(
			&q.ID,
			&q.TopicID,
			&q.Statement,
			&q.Problem,
			&q.ContentType,
			&q.Explanation,
			&q.QuestionType,
			&q.DifficultyLevel,
			&q.CreatedBy,
			&q.IsActive,
			&q.CreatedAt,
			&q.UpdatedAt,
		)
		if err != nil {
			return []*domain.Question{}, err
		}
		questions = append(questions, q)
	}

	if err = rows.Err(); err != nil {
		return []*domain.Question{}, err
	}

	// Garantir que sempre retorne um array, mesmo que vazio
	if questions == nil {
		questions = []*domain.Question{}
	}

	return questions, nil
}

// ListByTopic lista todas as questões de um tópico
func (r *QuestionRepositoryPG) ListByTopic(topicID int) ([]*domain.Question, error) {
	query := `
		SELECT id, topic_id, statement, problem, content_type, explanation, question_type, difficulty_level, created_by, is_active, created_at, updated_at
		FROM questions
		WHERE topic_id = $1
		ORDER BY created_at DESC`

	rows, err := r.DB.Query(query, topicID)
	if err != nil {
		return []*domain.Question{}, err
	}
	defer rows.Close()

	var questions []*domain.Question
	for rows.Next() {
		q := &domain.Question{}
		err := rows.Scan(
			&q.ID,
			&q.TopicID,
			&q.Statement,
			&q.Problem,
			&q.ContentType,
			&q.Explanation,
			&q.QuestionType,
			&q.DifficultyLevel,
			&q.CreatedBy,
			&q.IsActive,
			&q.CreatedAt,
			&q.UpdatedAt,
		)
		if err != nil {
			return []*domain.Question{}, err
		}
		questions = append(questions, q)
	}

	if err = rows.Err(); err != nil {
		return []*domain.Question{}, err
	}

	// Garantir que sempre retorne um array, mesmo que vazio
	if questions == nil {
		questions = []*domain.Question{}
	}

	return questions, nil
}

// ListAll lista todas as questões
func (r *QuestionRepositoryPG) ListAll() ([]*domain.Question, error) {
	query := `
		SELECT id, topic_id, statement, problem, content_type, explanation, question_type, difficulty_level, created_by, is_active, created_at, updated_at
		FROM questions
		WHERE is_active = true
		ORDER BY created_at DESC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return []*domain.Question{}, err
	}
	defer rows.Close()

	var questions []*domain.Question
	for rows.Next() {
		q := &domain.Question{}
		err := rows.Scan(
			&q.ID,
			&q.TopicID,
			&q.Statement,
			&q.Problem,
			&q.ContentType,
			&q.Explanation,
			&q.QuestionType,
			&q.DifficultyLevel,
			&q.CreatedBy,
			&q.IsActive,
			&q.CreatedAt,
			&q.UpdatedAt,
		)
		if err != nil {
			return []*domain.Question{}, err
		}
		questions = append(questions, q)
	}

	if err = rows.Err(); err != nil {
		return []*domain.Question{}, err
	}

	// Garantir que sempre retorne um array, mesmo que vazio
	if questions == nil {
		questions = []*domain.Question{}
	}

	return questions, nil
}

// ListAllWithDetails lista todas as questões com informações do exame e tópico
func (r *QuestionRepositoryPG) ListAllWithDetails() ([]*domain.QuestionWithDetails, error) {
	query := `
		SELECT 
			q.id, q.topic_id, q.statement, q.problem, q.content_type, q.explanation, 
			q.question_type, q.difficulty_level, q.created_by, q.is_active, q.created_at, q.updated_at,
			t.name as topic_name, t.exam_id,
			e.title as exam_title
		FROM questions q
		JOIN topics t ON q.topic_id = t.id
		JOIN exams e ON t.exam_id = e.id
		WHERE q.is_active = true
		ORDER BY q.created_at DESC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return []*domain.QuestionWithDetails{}, err
	}
	defer rows.Close()

	var questions []*domain.QuestionWithDetails
	for rows.Next() {
		q := &domain.Question{}
		var topicName, examTitle string
		var examID int
		
		err := rows.Scan(
			&q.ID,
			&q.TopicID,
			&q.Statement,
			&q.Problem,
			&q.ContentType,
			&q.Explanation,
			&q.QuestionType,
			&q.DifficultyLevel,
			&q.CreatedBy,
			&q.IsActive,
			&q.CreatedAt,
			&q.UpdatedAt,
			&topicName,
			&examID,
			&examTitle,
		)
		if err != nil {
			return []*domain.QuestionWithDetails{}, err
		}
		
		// Criar QuestionWithDetails
		questionWithDetails := &domain.QuestionWithDetails{
			Question:   q,
			ExamID:     examID,
			ExamTitle:  examTitle,
			TopicName:  topicName,
		}
		
		questions = append(questions, questionWithDetails)
	}

	if err = rows.Err(); err != nil {
		return []*domain.QuestionWithDetails{}, err
	}

	// Garantir que sempre retorne um array, mesmo que vazio
	if questions == nil {
		questions = []*domain.QuestionWithDetails{}
	}

	return questions, nil
}
