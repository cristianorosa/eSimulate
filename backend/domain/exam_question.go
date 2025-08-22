package domain

import "time"

// ExamQuestion represents the N:N relationship between exams and questions
type ExamQuestion struct {
	ExamID     int       `json:"exam_id" db:"exam_id"`
	QuestionID int       `json:"question_id" db:"question_id"`
	OrderIndex int       `json:"order_index" db:"order_index"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// ExamQuestionWithDetails represents exam-question relationship with detailed information
type ExamQuestionWithDetails struct {
	*ExamQuestion
	ExamTitle          string  `json:"exam_title" db:"exam_title"`
	ExamDescription    *string `json:"exam_description" db:"exam_description"`
	QuestionStatement  string  `json:"question_statement" db:"question_statement"`
	QuestionProblem    string  `json:"question_problem" db:"question_problem"`
	DifficultyLevel    int     `json:"difficulty_level" db:"difficulty_level"`
	TopicID            int     `json:"topic_id" db:"topic_id"`
	TopicName          string  `json:"topic_name" db:"topic_name"`
	TopicDescription   *string `json:"topic_description" db:"topic_description"`
	TopicWeightInExam  *float64 `json:"topic_weight_in_exam" db:"topic_weight_in_exam"`
	TopicQuestionsCount *int    `json:"topic_questions_count_in_exam" db:"topic_questions_count_in_exam"`
}

// ExamQuestionRepository defines the interface for exam-question relationship operations
type ExamQuestionRepository interface {
	// Create creates a new exam-question association
	Create(examQuestion *ExamQuestion) error
	
	// Delete removes an exam-question association
	Delete(examID, questionID int) error
	
	// FindByExamAndQuestion finds a specific exam-question association
	FindByExamAndQuestion(examID, questionID int) (*ExamQuestion, error)
	
	// ListByExam returns all questions associated with an exam
	ListByExam(examID int) ([]*ExamQuestion, error)
	
	// ListByQuestion returns all exams associated with a question
	ListByQuestion(questionID int) ([]*ExamQuestion, error)
	
	// ListByExamWithDetails returns questions for an exam with detailed information
	ListByExamWithDetails(examID int) ([]*ExamQuestionWithDetails, error)
	
	// UpdateOrder updates the order of questions in an exam
	UpdateOrder(examID int, questionOrders map[int]int) error
	
	// GetExamQuestionCount returns the total number of questions in an exam
	GetExamQuestionCount(examID int) (int, error)
	
	// GetQuestionExamCount returns the total number of exams a question belongs to
	GetQuestionExamCount(questionID int) (int, error)
	
	// BulkCreate creates multiple exam-question associations
	BulkCreate(examQuestions []*ExamQuestion) error
	
	// BulkDelete removes multiple exam-question associations for an exam
	BulkDelete(examID int, questionIDs []int) error
	
	// ReorderQuestions reorders all questions in an exam starting from index 1
	ReorderQuestions(examID int, questionIDs []int) error
	
	// ValidateExamQuestionAssociation validates if a question can be associated with an exam
	// This checks the business rules defined in the database trigger
	ValidateExamQuestionAssociation(examID, questionID int) error
}

// ExamQuestionFilters represents filters for querying exam-question relationships
type ExamQuestionFilters struct {
	ExamID         *int    `json:"exam_id,omitempty"`
	QuestionID     *int    `json:"question_id,omitempty"`
	TopicID        *int    `json:"topic_id,omitempty"`
	DifficultyLevel *int   `json:"difficulty_level,omitempty"`
	OrderFrom      *int    `json:"order_from,omitempty"`
	OrderTo        *int    `json:"order_to,omitempty"`
	Limit          int     `json:"limit"`
	Offset         int     `json:"offset"`
}

// ExamQuestionStats represents statistics for exam-question relationships
type ExamQuestionStats struct {
	ExamID              int                    `json:"exam_id"`
	TotalQuestions      int                    `json:"total_questions"`
	QuestionsByTopic    map[int]int           `json:"questions_by_topic"`
	QuestionsByDifficulty map[int]int         `json:"questions_by_difficulty"`
	TopicDistribution   []*TopicDistribution  `json:"topic_distribution"`
}

// TopicDistribution represents the distribution of questions by topic in an exam
type TopicDistribution struct {
	TopicID         int     `json:"topic_id"`
	TopicName       string  `json:"topic_name"`
	QuestionCount   int     `json:"question_count"`
	WeightPercentage float64 `json:"weight_percentage"`
	ExpectedCount   int     `json:"expected_count"`
	IsComplete      bool    `json:"is_complete"`
}
