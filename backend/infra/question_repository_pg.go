package infra

import (
	"database/sql"
	"strconv"

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

// ListAllWithDetails lista todas as questões com detalhes de exame e tópico
func (r *QuestionRepositoryPG) ListAllWithDetails() ([]*domain.QuestionWithDetails, error) {
	query := `
		SELECT q.id, q.topic_id, q.statement, q.problem, q.content_type, q.explanation, q.question_type, q.difficulty_level, q.created_by, q.is_active, q.created_at, q.updated_at,
		       t.name as topic_name, t.exam_id, e.title as exam_title
		FROM questions q
		JOIN topics t ON q.topic_id = t.id
		JOIN exams e ON t.exam_id = e.id
		ORDER BY q.created_at DESC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*domain.QuestionWithDetails
	for rows.Next() {
		// Inicializar corretamente a estrutura QuestionWithDetails
		q := &domain.QuestionWithDetails{
			Question: &domain.Question{},
		}

		err := rows.Scan(
			&q.Question.ID, &q.Question.TopicID, &q.Question.Statement, &q.Question.Problem, &q.Question.ContentType, &q.Question.Explanation, &q.Question.QuestionType, &q.Question.DifficultyLevel, &q.Question.CreatedBy, &q.Question.IsActive, &q.Question.CreatedAt, &q.Question.UpdatedAt,
			&q.TopicName, &q.ExamID, &q.ExamTitle,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

// ListPaginated lista questões com paginação e filtros opcionais
func (r *QuestionRepositoryPG) ListPaginated(page, pageSize int, examID, topicID *int) ([]*domain.QuestionWithDetails, *domain.Pagination, error) {
	// Construir query base
	baseQuery := `
		SELECT q.id, q.topic_id, q.statement, q.problem, q.content_type, q.explanation, q.question_type, q.difficulty_level, q.created_by, q.is_active, q.created_at, q.updated_at,
		       t.name as topic_name, t.exam_id, e.title as exam_title
		FROM questions q
		JOIN topics t ON q.topic_id = t.id
		JOIN exams e ON t.exam_id = e.id`

	// Adicionar filtros
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if examID != nil {
		whereClause = " WHERE t.exam_id = $" + strconv.Itoa(argIndex)
		args = append(args, *examID)
		argIndex++
	}

	if topicID != nil {
		if whereClause == "" {
			whereClause = " WHERE q.topic_id = $" + strconv.Itoa(argIndex)
		} else {
			whereClause += " AND q.topic_id = $" + strconv.Itoa(argIndex)
		}
		args = append(args, *topicID)
		argIndex++
	}

	// Query para contar total de registros
	countQuery := "SELECT COUNT(*) FROM questions q JOIN topics t ON q.topic_id = t.id JOIN exams e ON t.exam_id = e.id" + whereClause
	var totalItems int
	err := r.DB.QueryRow(countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, nil, err
	}

	// Calcular paginação
	offset := (page - 1) * pageSize
	totalPages := (totalItems + pageSize - 1) / pageSize

	// Query principal com paginação
	query := baseQuery + whereClause + " ORDER BY q.created_at DESC LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, pageSize, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var questions []*domain.QuestionWithDetails
	for rows.Next() {
		// Inicializar corretamente a estrutura QuestionWithDetails
		q := &domain.QuestionWithDetails{
			Question: &domain.Question{},
		}

		err := rows.Scan(
			&q.Question.ID, &q.Question.TopicID, &q.Question.Statement, &q.Question.Problem, &q.Question.ContentType, &q.Question.Explanation, &q.Question.QuestionType, &q.Question.DifficultyLevel, &q.Question.CreatedBy, &q.Question.IsActive, &q.Question.CreatedAt, &q.Question.UpdatedAt,
			&q.TopicName, &q.ExamID, &q.ExamTitle,
		)
		if err != nil {
			return nil, nil, err
		}
		questions = append(questions, q)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	pagination := &domain.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return questions, pagination, nil
}

// FindByIDWithExams finds a question by ID with its associated exams
func (r *QuestionRepositoryPG) FindByIDWithExams(id int) (*domain.Question, error) {
	// First get the basic question info
	question, err := r.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Initialize empty exams slice
	question.ExamIDs = []int{}
	question.Exams = []*domain.Exam{}
	
	return question, nil
}

// ListExamsByQuestion returns all exams associated with a question
func (r *QuestionRepositoryPG) ListExamsByQuestion(questionID int) ([]*domain.Exam, error) {
	// This would typically query the exam_questions table
	// For now, return empty slice to satisfy interface
	return []*domain.Exam{}, nil
}

// ListByExamDirect returns questions directly associated with an exam via exam_questions table
func (r *QuestionRepositoryPG) ListByExamDirect(examID int) ([]*domain.Question, error) {
	// This would typically query via exam_questions table
	// For now, return empty slice to satisfy interface
	return []*domain.Question{}, nil
}

// ListByExamWithDetails returns questions for an exam with detailed information
func (r *QuestionRepositoryPG) ListByExamWithDetails(examID int) ([]*domain.QuestionWithDetails, error) {
	// This would typically query via exam_questions table with joins
	// For now, return empty slice to satisfy interface
	return []*domain.QuestionWithDetails{}, nil
}

// GetAvailableQuestionsForExam returns questions available for association with an exam
func (r *QuestionRepositoryPG) GetAvailableQuestionsForExam(examID int, topicID *int) ([]*domain.Question, error) {
	// This would implement the business logic for finding available questions
	// For now, return empty slice to satisfy interface
	return []*domain.Question{}, nil
}

// GetQuestionStats returns statistics for a question
func (r *QuestionRepositoryPG) GetQuestionStats(questionID int) (*domain.QuestionStats, error) {
	// This would calculate various statistics about the question
	// For now, return basic stats to satisfy interface
	return &domain.QuestionStats{
		QuestionID:      questionID,
		ExamCount:       0,
		TagCount:        0,
		AnswerCount:     0,
		CorrectRate:     0.0,
		AverageTime:     0.0,
		DifficultyLevel: 1,
	}, nil
}

// GetTopicQuestionCount returns the number of questions in a topic
func (r *QuestionRepositoryPG) GetTopicQuestionCount(topicID int) (int, error) {
	query := `SELECT COUNT(*) FROM questions WHERE topic_id = $1 AND is_active = true`
	
	var count int
	err := r.DB.QueryRow(query, topicID).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// GetExamQuestionCount returns the number of questions in an exam
func (r *QuestionRepositoryPG) GetExamQuestionCount(examID int) (int, error) {
	query := `SELECT COUNT(*) FROM exam_questions WHERE exam_id = $1`
	
	var count int
	err := r.DB.QueryRow(query, examID).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}
