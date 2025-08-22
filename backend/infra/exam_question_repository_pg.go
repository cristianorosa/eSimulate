package infra

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// ExamQuestionRepositoryPG implementa o repositório de relacionamento exame-questão para PostgreSQL
type ExamQuestionRepositoryPG struct {
	DB *sql.DB
}

// NewExamQuestionRepositoryPG cria uma nova instância do repositório
func NewExamQuestionRepositoryPG(db *sql.DB) domain.ExamQuestionRepository {
	return &ExamQuestionRepositoryPG{DB: db}
}

// Create creates a new exam-question association
func (r *ExamQuestionRepositoryPG) Create(examQuestion *domain.ExamQuestion) error {
	query := `
		INSERT INTO exam_questions (exam_id, question_id, order_index, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	
	now := time.Now()
	_, err := r.DB.Exec(query, 
		examQuestion.ExamID, 
		examQuestion.QuestionID, 
		examQuestion.OrderIndex,
		now,
		now,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create exam-question association: %w", err)
	}
	
	examQuestion.CreatedAt = now
	examQuestion.UpdatedAt = now
	return nil
}

// Delete removes an exam-question association
func (r *ExamQuestionRepositoryPG) Delete(examID, questionID int) error {
	query := `DELETE FROM exam_questions WHERE exam_id = $1 AND question_id = $2`
	
	result, err := r.DB.Exec(query, examID, questionID)
	if err != nil {
		return fmt.Errorf("failed to delete exam-question association: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("exam-question association not found")
	}
	
	return nil
}

// FindByExamAndQuestion finds a specific exam-question association
func (r *ExamQuestionRepositoryPG) FindByExamAndQuestion(examID, questionID int) (*domain.ExamQuestion, error) {
	query := `
		SELECT exam_id, question_id, order_index, created_at, updated_at
		FROM exam_questions
		WHERE exam_id = $1 AND question_id = $2`
	
	examQuestion := &domain.ExamQuestion{}
	err := r.DB.QueryRow(query, examID, questionID).Scan(
		&examQuestion.ExamID,
		&examQuestion.QuestionID,
		&examQuestion.OrderIndex,
		&examQuestion.CreatedAt,
		&examQuestion.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("exam-question association not found")
		}
		return nil, fmt.Errorf("failed to find exam-question association: %w", err)
	}
	
	return examQuestion, nil
}

// ListByExam returns all questions associated with an exam
func (r *ExamQuestionRepositoryPG) ListByExam(examID int) ([]*domain.ExamQuestion, error) {
	query := `
		SELECT exam_id, question_id, order_index, created_at, updated_at
		FROM exam_questions
		WHERE exam_id = $1
		ORDER BY order_index ASC`
	
	rows, err := r.DB.Query(query, examID)
	if err != nil {
		return nil, fmt.Errorf("failed to list questions for exam: %w", err)
	}
	defer rows.Close()
	
	var examQuestions []*domain.ExamQuestion
	for rows.Next() {
		examQuestion := &domain.ExamQuestion{}
		err := rows.Scan(
			&examQuestion.ExamID,
			&examQuestion.QuestionID,
			&examQuestion.OrderIndex,
			&examQuestion.CreatedAt,
			&examQuestion.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exam-question: %w", err)
		}
		examQuestions = append(examQuestions, examQuestion)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over exam-question rows: %w", err)
	}
	
	return examQuestions, nil
}

// ListByQuestion returns all exams associated with a question
func (r *ExamQuestionRepositoryPG) ListByQuestion(questionID int) ([]*domain.ExamQuestion, error) {
	query := `
		SELECT exam_id, question_id, order_index, created_at, updated_at
		FROM exam_questions
		WHERE question_id = $1
		ORDER BY exam_id ASC`
	
	rows, err := r.DB.Query(query, questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list exams for question: %w", err)
	}
	defer rows.Close()
	
	var examQuestions []*domain.ExamQuestion
	for rows.Next() {
		examQuestion := &domain.ExamQuestion{}
		err := rows.Scan(
			&examQuestion.ExamID,
			&examQuestion.QuestionID,
			&examQuestion.OrderIndex,
			&examQuestion.CreatedAt,
			&examQuestion.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exam-question: %w", err)
		}
		examQuestions = append(examQuestions, examQuestion)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over exam-question rows: %w", err)
	}
	
	return examQuestions, nil
}

// ListByExamWithDetails returns questions for an exam with detailed information
func (r *ExamQuestionRepositoryPG) ListByExamWithDetails(examID int) ([]*domain.ExamQuestionWithDetails, error) {
	query := `
		SELECT 
			eq.exam_id, eq.question_id, eq.order_index, eq.created_at,
			e.title, e.description,
			q.statement, q.problem, q.difficulty_level, q.topic_id,
			t.name, t.description,
			et.weight_percentage, et.questions_count
		FROM exam_questions eq
		JOIN exams e ON eq.exam_id = e.id
		JOIN questions q ON eq.question_id = q.id
		JOIN topics t ON q.topic_id = t.id
		LEFT JOIN exam_topics et ON (et.exam_id = eq.exam_id AND et.topic_id = q.topic_id)
		WHERE eq.exam_id = $1
		ORDER BY eq.order_index ASC`
	
	rows, err := r.DB.Query(query, examID)
	if err != nil {
		return nil, fmt.Errorf("failed to list exam questions with details: %w", err)
	}
	defer rows.Close()
	
	var examQuestions []*domain.ExamQuestionWithDetails
	for rows.Next() {
		examQuestion := &domain.ExamQuestionWithDetails{
			ExamQuestion: &domain.ExamQuestion{},
		}
		
		var examDesc, topicDesc sql.NullString
		var topicWeight sql.NullFloat64
		var topicQuestionsCount sql.NullInt64
		
		err := rows.Scan(
			&examQuestion.ExamID,
			&examQuestion.QuestionID,
			&examQuestion.OrderIndex,
			&examQuestion.CreatedAt,
			&examQuestion.ExamTitle,
			&examDesc,
			&examQuestion.QuestionStatement,
			&examQuestion.QuestionProblem,
			&examQuestion.DifficultyLevel,
			&examQuestion.TopicID,
			&examQuestion.TopicName,
			&topicDesc,
			&topicWeight,
			&topicQuestionsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exam-question with details: %w", err)
		}
		
		// Handle nullable fields
		if examDesc.Valid {
			examQuestion.ExamDescription = &examDesc.String
		}
		if topicDesc.Valid {
			examQuestion.TopicDescription = &topicDesc.String
		}
		if topicWeight.Valid {
			weight := topicWeight.Float64
			examQuestion.TopicWeightInExam = &weight
		}
		if topicQuestionsCount.Valid {
			count := int(topicQuestionsCount.Int64)
			examQuestion.TopicQuestionsCount = &count
		}
		
		examQuestions = append(examQuestions, examQuestion)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over exam-question details rows: %w", err)
	}
	
	return examQuestions, nil
}

// UpdateOrder updates the order of questions in an exam
func (r *ExamQuestionRepositoryPG) UpdateOrder(examID int, questionOrders map[int]int) error {
	if len(questionOrders) == 0 {
		return nil
	}
	
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	query := `UPDATE exam_questions SET order_index = $1, updated_at = $2 WHERE exam_id = $3 AND question_id = $4`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()
	
	now := time.Now()
	for questionID, orderIndex := range questionOrders {
		_, err := stmt.Exec(orderIndex, now, examID, questionID)
		if err != nil {
			return fmt.Errorf("failed to update order for question %d: %w", questionID, err)
		}
	}
	
	return tx.Commit()
}

// GetExamQuestionCount returns the total number of questions in an exam
func (r *ExamQuestionRepositoryPG) GetExamQuestionCount(examID int) (int, error) {
	query := `SELECT COUNT(*) FROM exam_questions WHERE exam_id = $1`
	
	var count int
	err := r.DB.QueryRow(query, examID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get exam question count: %w", err)
	}
	
	return count, nil
}

// GetQuestionExamCount returns the total number of exams a question belongs to
func (r *ExamQuestionRepositoryPG) GetQuestionExamCount(questionID int) (int, error) {
	query := `SELECT COUNT(*) FROM exam_questions WHERE question_id = $1`
	
	var count int
	err := r.DB.QueryRow(query, questionID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get question exam count: %w", err)
	}
	
	return count, nil
}

// BulkCreate creates multiple exam-question associations
func (r *ExamQuestionRepositoryPG) BulkCreate(examQuestions []*domain.ExamQuestion) error {
	if len(examQuestions) == 0 {
		return nil
	}
	
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	query := `
		INSERT INTO exam_questions (exam_id, question_id, order_index, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()
	
	now := time.Now()
	for _, examQuestion := range examQuestions {
		_, err := stmt.Exec(
			examQuestion.ExamID,
			examQuestion.QuestionID,
			examQuestion.OrderIndex,
			now,
			now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert exam-question association: %w", err)
		}
		
		examQuestion.CreatedAt = now
		examQuestion.UpdatedAt = now
	}
	
	return tx.Commit()
}

// BulkDelete removes multiple exam-question associations for an exam
func (r *ExamQuestionRepositoryPG) BulkDelete(examID int, questionIDs []int) error {
	if len(questionIDs) == 0 {
		return nil
	}
	
	placeholders := make([]string, len(questionIDs))
	args := make([]interface{}, len(questionIDs)+1)
	args[0] = examID
	
	for i, questionID := range questionIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = questionID
	}
	
	query := fmt.Sprintf(
		`DELETE FROM exam_questions WHERE exam_id = $1 AND question_id IN (%s)`,
		strings.Join(placeholders, ","),
	)
	
	_, err := r.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to bulk delete exam-question associations: %w", err)
	}
	
	return nil
}

// ReorderQuestions reorders all questions in an exam starting from index 1
func (r *ExamQuestionRepositoryPG) ReorderQuestions(examID int, questionIDs []int) error {
	if len(questionIDs) == 0 {
		return nil
	}
	
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	query := `UPDATE exam_questions SET order_index = $1, updated_at = $2 WHERE exam_id = $3 AND question_id = $4`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()
	
	now := time.Now()
	for i, questionID := range questionIDs {
		orderIndex := i + 1 // Start from 1
		_, err := stmt.Exec(orderIndex, now, examID, questionID)
		if err != nil {
			return fmt.Errorf("failed to reorder question %d: %w", questionID, err)
		}
	}
	
	return tx.Commit()
}

// ValidateExamQuestionAssociation validates if a question can be associated with an exam
func (r *ExamQuestionRepositoryPG) ValidateExamQuestionAssociation(examID, questionID int) error {
	// This validation is implemented as a database trigger
	// We can perform a test insertion to validate the association
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin validation transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Try to insert with a temporary high order_index
	query := `
		INSERT INTO exam_questions (exam_id, question_id, order_index, created_at, updated_at)
		VALUES ($1, $2, 999999, $3, $4)`
	
	now := time.Now()
	_, err = tx.Exec(query, examID, questionID, now, now)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	
	// Remove the test insertion
	deleteQuery := `DELETE FROM exam_questions WHERE exam_id = $1 AND question_id = $2 AND order_index = 999999`
	_, err = tx.Exec(deleteQuery, examID, questionID)
	if err != nil {
		return fmt.Errorf("failed to cleanup validation test: %w", err)
	}
	
	return tx.Commit()
}
