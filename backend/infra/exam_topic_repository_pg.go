package infra

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// ExamTopicRepositoryPG implementa o repositório de relacionamento exame-tópico para PostgreSQL
type ExamTopicRepositoryPG struct {
	DB *sql.DB
}

// NewExamTopicRepositoryPG cria uma nova instância do repositório
func NewExamTopicRepositoryPG(db *sql.DB) domain.ExamTopicRepository {
	return &ExamTopicRepositoryPG{DB: db}
}

// Create creates a new exam-topic association
func (r *ExamTopicRepositoryPG) Create(examTopic *domain.ExamTopic) error {
	query := `
		INSERT INTO exam_topics (
			exam_id, topic_id, questions_count, weight_percentage, order_index,
			difficulty_easy_percentage, difficulty_medium_percentage, difficulty_hard_percentage,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	now := time.Now()
	_, err := r.DB.Exec(query,
		examTopic.ExamID,
		examTopic.TopicID,
		examTopic.QuestionsCount,
		examTopic.WeightPercentage,
		examTopic.OrderIndex,
		examTopic.DifficultyEasyPercentage,
		examTopic.DifficultyMediumPercentage,
		examTopic.DifficultyHardPercentage,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create exam-topic association: %w", err)
	}

	examTopic.CreatedAt = now
	examTopic.UpdatedAt = now
	return nil
}

// Update updates an existing exam-topic association
func (r *ExamTopicRepositoryPG) Update(examTopic *domain.ExamTopic) error {
	query := `
		UPDATE exam_topics SET 
			questions_count = $1, weight_percentage = $2, order_index = $3,
			difficulty_easy_percentage = $4, difficulty_medium_percentage = $5, 
			difficulty_hard_percentage = $6, updated_at = $7
		WHERE exam_id = $8 AND topic_id = $9`

	now := time.Now()
	result, err := r.DB.Exec(query,
		examTopic.QuestionsCount,
		examTopic.WeightPercentage,
		examTopic.OrderIndex,
		examTopic.DifficultyEasyPercentage,
		examTopic.DifficultyMediumPercentage,
		examTopic.DifficultyHardPercentage,
		now,
		examTopic.ExamID,
		examTopic.TopicID,
	)

	if err != nil {
		return fmt.Errorf("failed to update exam-topic association: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("exam-topic association not found")
	}

	examTopic.UpdatedAt = now
	return nil
}

// Delete removes an exam-topic association
func (r *ExamTopicRepositoryPG) Delete(examID, topicID int) error {
	query := `DELETE FROM exam_topics WHERE exam_id = $1 AND topic_id = $2`

	result, err := r.DB.Exec(query, examID, topicID)
	if err != nil {
		return fmt.Errorf("failed to delete exam-topic association: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("exam-topic association not found")
	}

	return nil
}

// FindByExamAndTopic finds a specific exam-topic association
func (r *ExamTopicRepositoryPG) FindByExamAndTopic(examID, topicID int) (*domain.ExamTopic, error) {
	query := `
		SELECT exam_id, topic_id, questions_count, weight_percentage, order_index,
			   difficulty_easy_percentage, difficulty_medium_percentage, difficulty_hard_percentage,
			   created_at, updated_at
		FROM exam_topics
		WHERE exam_id = $1 AND topic_id = $2`

	examTopic := &domain.ExamTopic{}
	err := r.DB.QueryRow(query, examID, topicID).Scan(
		&examTopic.ExamID,
		&examTopic.TopicID,
		&examTopic.QuestionsCount,
		&examTopic.WeightPercentage,
		&examTopic.OrderIndex,
		&examTopic.DifficultyEasyPercentage,
		&examTopic.DifficultyMediumPercentage,
		&examTopic.DifficultyHardPercentage,
		&examTopic.CreatedAt,
		&examTopic.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("exam-topic association not found")
		}
		return nil, fmt.Errorf("failed to find exam-topic association: %w", err)
	}

	return examTopic, nil
}

// ListByExam returns all topics associated with an exam
func (r *ExamTopicRepositoryPG) ListByExam(examID int) ([]*domain.ExamTopic, error) {
	query := `
		SELECT exam_id, topic_id, questions_count, weight_percentage, order_index,
			   difficulty_easy_percentage, difficulty_medium_percentage, difficulty_hard_percentage,
			   created_at, updated_at
		FROM exam_topics
		WHERE exam_id = $1
		ORDER BY order_index ASC`

	rows, err := r.DB.Query(query, examID)
	if err != nil {
		return nil, fmt.Errorf("failed to list topics for exam: %w", err)
	}
	defer rows.Close()

	var examTopics []*domain.ExamTopic
	for rows.Next() {
		examTopic := &domain.ExamTopic{}
		err := rows.Scan(
			&examTopic.ExamID,
			&examTopic.TopicID,
			&examTopic.QuestionsCount,
			&examTopic.WeightPercentage,
			&examTopic.OrderIndex,
			&examTopic.DifficultyEasyPercentage,
			&examTopic.DifficultyMediumPercentage,
			&examTopic.DifficultyHardPercentage,
			&examTopic.CreatedAt,
			&examTopic.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exam-topic: %w", err)
		}
		examTopics = append(examTopics, examTopic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over exam-topic rows: %w", err)
	}

	return examTopics, nil
}

// ListByTopic returns all exams associated with a topic
func (r *ExamTopicRepositoryPG) ListByTopic(topicID int) ([]*domain.ExamTopic, error) {
	query := `
		SELECT exam_id, topic_id, questions_count, weight_percentage, order_index,
			   difficulty_easy_percentage, difficulty_medium_percentage, difficulty_hard_percentage,
			   created_at, updated_at
		FROM exam_topics
		WHERE topic_id = $1
		ORDER BY exam_id ASC`

	rows, err := r.DB.Query(query, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to list exams for topic: %w", err)
	}
	defer rows.Close()

	var examTopics []*domain.ExamTopic
	for rows.Next() {
		examTopic := &domain.ExamTopic{}
		err := rows.Scan(
			&examTopic.ExamID,
			&examTopic.TopicID,
			&examTopic.QuestionsCount,
			&examTopic.WeightPercentage,
			&examTopic.OrderIndex,
			&examTopic.DifficultyEasyPercentage,
			&examTopic.DifficultyMediumPercentage,
			&examTopic.DifficultyHardPercentage,
			&examTopic.CreatedAt,
			&examTopic.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exam-topic: %w", err)
		}
		examTopics = append(examTopics, examTopic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over exam-topic rows: %w", err)
	}

	return examTopics, nil
}

// ListByExamWithDetails returns topics for an exam with detailed information
func (r *ExamTopicRepositoryPG) ListByExamWithDetails(examID int) ([]*domain.ExamTopicWithDetails, error) {
	query := `
		SELECT 
			et.exam_id, et.topic_id, et.questions_count, et.weight_percentage, et.order_index,
			et.difficulty_easy_percentage, et.difficulty_medium_percentage, et.difficulty_hard_percentage,
			et.created_at, et.updated_at,
			e.title, e.description,
			t.name, t.description,
			COALESCE(actual_q.count, 0) as actual_questions,
			COALESCE(easy_q.count, 0) as easy_questions,
			COALESCE(medium_q.count, 0) as medium_questions,
			COALESCE(hard_q.count, 0) as hard_questions,
			CASE WHEN COALESCE(actual_q.count, 0) >= et.questions_count THEN true ELSE false END as is_complete
		FROM exam_topics et
		JOIN exams e ON et.exam_id = e.id
		JOIN topics t ON et.topic_id = t.id
		LEFT JOIN (
			SELECT eq.exam_id, q.topic_id, COUNT(*) as count
			FROM exam_questions eq
			JOIN questions q ON eq.question_id = q.id
			GROUP BY eq.exam_id, q.topic_id
		) actual_q ON (et.exam_id = actual_q.exam_id AND et.topic_id = actual_q.topic_id)
		LEFT JOIN (
			SELECT eq.exam_id, q.topic_id, COUNT(*) as count
			FROM exam_questions eq
			JOIN questions q ON eq.question_id = q.id
			WHERE q.difficulty_level = 1
			GROUP BY eq.exam_id, q.topic_id
		) easy_q ON (et.exam_id = easy_q.exam_id AND et.topic_id = easy_q.topic_id)
		LEFT JOIN (
			SELECT eq.exam_id, q.topic_id, COUNT(*) as count
			FROM exam_questions eq
			JOIN questions q ON eq.question_id = q.id
			WHERE q.difficulty_level = 2
			GROUP BY eq.exam_id, q.topic_id
		) medium_q ON (et.exam_id = medium_q.exam_id AND et.topic_id = medium_q.topic_id)
		LEFT JOIN (
			SELECT eq.exam_id, q.topic_id, COUNT(*) as count
			FROM exam_questions eq
			JOIN questions q ON eq.question_id = q.id
			WHERE q.difficulty_level >= 3
			GROUP BY eq.exam_id, q.topic_id
		) hard_q ON (et.exam_id = hard_q.exam_id AND et.topic_id = hard_q.topic_id)
		WHERE et.exam_id = $1
		ORDER BY et.order_index ASC`

	rows, err := r.DB.Query(query, examID)
	if err != nil {
		return nil, fmt.Errorf("failed to list exam topics with details: %w", err)
	}
	defer rows.Close()

	var examTopics []*domain.ExamTopicWithDetails
	for rows.Next() {
		examTopic := &domain.ExamTopicWithDetails{
			ExamTopic: &domain.ExamTopic{},
		}

		var examDesc, topicDesc sql.NullString

		err := rows.Scan(
			&examTopic.ExamID,
			&examTopic.TopicID,
			&examTopic.QuestionsCount,
			&examTopic.WeightPercentage,
			&examTopic.OrderIndex,
			&examTopic.DifficultyEasyPercentage,
			&examTopic.DifficultyMediumPercentage,
			&examTopic.DifficultyHardPercentage,
			&examTopic.CreatedAt,
			&examTopic.UpdatedAt,
			&examTopic.ExamTitle,
			&examDesc,
			&examTopic.TopicName,
			&topicDesc,
			&examTopic.ActualQuestions,
			&examTopic.EasyQuestions,
			&examTopic.MediumQuestions,
			&examTopic.HardQuestions,
			&examTopic.IsComplete,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exam-topic with details: %w", err)
		}

		// Handle nullable fields
		if examDesc.Valid {
			examTopic.ExamDescription = &examDesc.String
		}
		if topicDesc.Valid {
			examTopic.TopicDescription = &topicDesc.String
		}

		examTopics = append(examTopics, examTopic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over exam-topic details rows: %w", err)
	}

	return examTopics, nil
}

// ListByTopicWithDetails returns exams for a topic with detailed information
func (r *ExamTopicRepositoryPG) ListByTopicWithDetails(topicID int) ([]*domain.ExamTopicWithDetails, error) {
	query := `
		SELECT 
			et.exam_id, et.topic_id, et.questions_count, et.weight_percentage, et.order_index,
			et.difficulty_easy_percentage, et.difficulty_medium_percentage, et.difficulty_hard_percentage,
			et.created_at, et.updated_at,
			e.title, e.description,
			t.name, t.description,
			COALESCE(actual_q.count, 0) as actual_questions,
			COALESCE(easy_q.count, 0) as easy_questions,
			COALESCE(medium_q.count, 0) as medium_questions,
			COALESCE(hard_q.count, 0) as hard_questions,
			CASE WHEN COALESCE(actual_q.count, 0) >= et.questions_count THEN true ELSE false END as is_complete
		FROM exam_topics et
		JOIN exams e ON et.exam_id = e.id
		JOIN topics t ON et.topic_id = t.id
		LEFT JOIN (
			SELECT eq.exam_id, q.topic_id, COUNT(*) as count
			FROM exam_questions eq
			JOIN questions q ON eq.question_id = q.id
			GROUP BY eq.exam_id, q.topic_id
		) actual_q ON (et.exam_id = actual_q.exam_id AND et.topic_id = actual_q.topic_id)
		LEFT JOIN (
			SELECT eq.exam_id, q.topic_id, COUNT(*) as count
			FROM exam_questions eq
			JOIN questions q ON eq.question_id = q.id
			WHERE q.difficulty_level = 1
			GROUP BY eq.exam_id, q.topic_id
		) easy_q ON (et.exam_id = easy_q.exam_id AND et.topic_id = easy_q.topic_id)
		LEFT JOIN (
			SELECT eq.exam_id, q.topic_id, COUNT(*) as count
			FROM exam_questions eq
			JOIN questions q ON eq.question_id = q.id
			WHERE q.difficulty_level = 2
			GROUP BY eq.exam_id, q.topic_id
		) medium_q ON (et.exam_id = medium_q.exam_id AND et.topic_id = medium_q.topic_id)
		LEFT JOIN (
			SELECT eq.exam_id, q.topic_id, COUNT(*) as count
			FROM exam_questions eq
			JOIN questions q ON eq.question_id = q.id
			WHERE q.difficulty_level >= 3
			GROUP BY eq.exam_id, q.topic_id
		) hard_q ON (et.exam_id = hard_q.exam_id AND et.topic_id = hard_q.topic_id)
		WHERE et.topic_id = $1
		ORDER BY e.title ASC`

	rows, err := r.DB.Query(query, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to list topic exams with details: %w", err)
	}
	defer rows.Close()

	var examTopics []*domain.ExamTopicWithDetails
	for rows.Next() {
		examTopic := &domain.ExamTopicWithDetails{
			ExamTopic: &domain.ExamTopic{},
		}

		var examDesc, topicDesc sql.NullString

		err := rows.Scan(
			&examTopic.ExamID,
			&examTopic.TopicID,
			&examTopic.QuestionsCount,
			&examTopic.WeightPercentage,
			&examTopic.OrderIndex,
			&examTopic.DifficultyEasyPercentage,
			&examTopic.DifficultyMediumPercentage,
			&examTopic.DifficultyHardPercentage,
			&examTopic.CreatedAt,
			&examTopic.UpdatedAt,
			&examTopic.ExamTitle,
			&examDesc,
			&examTopic.TopicName,
			&topicDesc,
			&examTopic.ActualQuestions,
			&examTopic.EasyQuestions,
			&examTopic.MediumQuestions,
			&examTopic.HardQuestions,
			&examTopic.IsComplete,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exam-topic with details: %w", err)
		}

		// Handle nullable fields
		if examDesc.Valid {
			examTopic.ExamDescription = &examDesc.String
		}
		if topicDesc.Valid {
			examTopic.TopicDescription = &topicDesc.String
		}

		examTopics = append(examTopics, examTopic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over topic-exam details rows: %w", err)
	}

	return examTopics, nil
}

// GetExamTopicsCount returns the total number of topics in an exam
func (r *ExamTopicRepositoryPG) GetExamTopicsCount(examID int) (int, error) {
	query := `SELECT COUNT(*) FROM exam_topics WHERE exam_id = $1`

	var count int
	err := r.DB.QueryRow(query, examID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get exam topics count: %w", err)
	}

	return count, nil
}

// GetTopicExamsCount returns the total number of exams a topic belongs to
func (r *ExamTopicRepositoryPG) GetTopicExamsCount(topicID int) (int, error) {
	query := `SELECT COUNT(*) FROM exam_topics WHERE topic_id = $1`

	var count int
	err := r.DB.QueryRow(query, topicID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get topic exams count: %w", err)
	}

	return count, nil
}

// ValidateWeightDistribution validates that weight percentages sum to 100 for an exam
func (r *ExamTopicRepositoryPG) ValidateWeightDistribution(examID int) error {
	query := `SELECT SUM(weight_percentage) FROM exam_topics WHERE exam_id = $1`

	var totalWeight sql.NullFloat64
	err := r.DB.QueryRow(query, examID).Scan(&totalWeight)
	if err != nil {
		return fmt.Errorf("failed to calculate total weight: %w", err)
	}

	if !totalWeight.Valid {
		return fmt.Errorf("exam has no topics")
	}

	if totalWeight.Float64 != 100.0 {
		return fmt.Errorf("total weight percentage must be 100, got %.2f", totalWeight.Float64)
	}

	return nil
}

// ValidateDifficultyDistribution validates that difficulty percentages sum to 100 for each topic
func (r *ExamTopicRepositoryPG) ValidateDifficultyDistribution(examTopic *domain.ExamTopic) error {
	total := examTopic.DifficultyEasyPercentage + examTopic.DifficultyMediumPercentage + examTopic.DifficultyHardPercentage

	if total != 100.0 {
		return fmt.Errorf("difficulty percentages must sum to 100, got %.2f", total)
	}

	return nil
}

// CalculateQuestionsByDifficulty calculates how many questions of each difficulty are needed
func (r *ExamTopicRepositoryPG) CalculateQuestionsByDifficulty(examTopic *domain.ExamTopic) (*domain.DifficultyDistribution, error) {
	totalQuestions := examTopic.QuestionsCount

	easyCount := int(float64(totalQuestions) * examTopic.DifficultyEasyPercentage / 100.0)
	mediumCount := int(float64(totalQuestions) * examTopic.DifficultyMediumPercentage / 100.0)
	hardCount := totalQuestions - easyCount - mediumCount // Ensure total adds up exactly

	return &domain.DifficultyDistribution{
		EasyCount:   easyCount,
		MediumCount: mediumCount,
		HardCount:   hardCount,
		TotalCount:  totalQuestions,
	}, nil
}

// UpdateQuestionsCount updates the actual questions count for a topic in an exam
func (r *ExamTopicRepositoryPG) UpdateQuestionsCount(examID, topicID, count int) error {
	query := `
		UPDATE exam_topics 
		SET questions_count = $1, updated_at = $2 
		WHERE exam_id = $3 AND topic_id = $4`

	result, err := r.DB.Exec(query, count, time.Now(), examID, topicID)
	if err != nil {
		return fmt.Errorf("failed to update questions count: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("exam-topic association not found")
	}

	return nil
}

// BulkCreate creates multiple exam-topic associations
func (r *ExamTopicRepositoryPG) BulkCreate(examTopics []*domain.ExamTopic) error {
	if len(examTopics) == 0 {
		return nil
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO exam_topics (
			exam_id, topic_id, questions_count, weight_percentage, order_index,
			difficulty_easy_percentage, difficulty_medium_percentage, difficulty_hard_percentage,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	for _, examTopic := range examTopics {
		_, err := stmt.Exec(
			examTopic.ExamID,
			examTopic.TopicID,
			examTopic.QuestionsCount,
			examTopic.WeightPercentage,
			examTopic.OrderIndex,
			examTopic.DifficultyEasyPercentage,
			examTopic.DifficultyMediumPercentage,
			examTopic.DifficultyHardPercentage,
			now,
			now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert exam-topic association: %w", err)
		}

		examTopic.CreatedAt = now
		examTopic.UpdatedAt = now
	}

	return tx.Commit()
}

// BulkUpdate updates multiple exam-topic associations
func (r *ExamTopicRepositoryPG) BulkUpdate(examTopics []*domain.ExamTopic) error {
	if len(examTopics) == 0 {
		return nil
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE exam_topics SET 
			questions_count = $1, weight_percentage = $2, order_index = $3,
			difficulty_easy_percentage = $4, difficulty_medium_percentage = $5, 
			difficulty_hard_percentage = $6, updated_at = $7
		WHERE exam_id = $8 AND topic_id = $9`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	for _, examTopic := range examTopics {
		_, err := stmt.Exec(
			examTopic.QuestionsCount,
			examTopic.WeightPercentage,
			examTopic.OrderIndex,
			examTopic.DifficultyEasyPercentage,
			examTopic.DifficultyMediumPercentage,
			examTopic.DifficultyHardPercentage,
			now,
			examTopic.ExamID,
			examTopic.TopicID,
		)
		if err != nil {
			return fmt.Errorf("failed to update exam-topic association: %w", err)
		}

		examTopic.UpdatedAt = now
	}

	return tx.Commit()
}

// ReorderTopics reorders all topics in an exam
func (r *ExamTopicRepositoryPG) ReorderTopics(examID int, topicOrders map[int]int) error {
	if len(topicOrders) == 0 {
		return nil
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE exam_topics SET order_index = $1, updated_at = $2 WHERE exam_id = $3 AND topic_id = $4`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	for topicID, orderIndex := range topicOrders {
		_, err := stmt.Exec(orderIndex, now, examID, topicID)
		if err != nil {
			return fmt.Errorf("failed to reorder topic %d: %w", topicID, err)
		}
	}

	return tx.Commit()
}
