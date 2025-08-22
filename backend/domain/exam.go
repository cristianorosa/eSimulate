package domain

import "time"

// Exam representa uma prova/exame
type Exam struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	AreaID         int       `json:"area_id"`
	MaxTimeMinutes int       `json:"max_time_minutes"`
	PassingScore   float64   `json:"passing_score"`
	QuestionsCount int       `json:"questions_count"`
	IsActive       bool      `json:"is_active"`
	CreatedBy      int       `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	
	// N:N Relationships
	Topics      []*ExamTopic `json:"topics,omitempty"`       // Tópicos associados ao exame
	Questions   []*Question  `json:"questions,omitempty"`    // Questões associadas ao exame
	QuestionIDs []int        `json:"question_ids,omitempty"` // IDs das questões (para facilitar operações)
}

// Topic representa um tópico (agora independente de exames)
type Topic struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	
	// Campos legados para compatibilidade (deprecated)
	ExamID           int     `json:"exam_id,omitempty"`           // Deprecated: usar ExamTopic
	WeightPercentage float64 `json:"weight_percentage,omitempty"` // Deprecated: usar ExamTopic
	OrderIndex       int     `json:"order_index,omitempty"`       // Deprecated: usar ExamTopic
	QuestionsCount   int     `json:"questions_count,omitempty"`   // Deprecated: usar ExamTopic
}

// ExamRepository define as operações de persistência para exames
type ExamRepository interface {
	// Basic CRUD operations
	Create(exam *Exam) error
	Update(exam *Exam) error
	Delete(id int) error
	FindByID(id int) (*Exam, error)
	ListAll() ([]*Exam, error)
	ListByArea(areaID int) ([]*Exam, error)
	ListPaginated(page, pageSize int, areaID *int) ([]*Exam, *Pagination, error)
	
	// N:N Relationship operations
	FindByIDWithTopics(id int) (*Exam, error)
	ListByAreaWithTopics(areaID int) ([]*Exam, error)
	FindByIDWithQuestions(id int) (*Exam, error)
	FindByIDWithTopicsAndQuestions(id int) (*Exam, error)
	ListByAreaWithQuestions(areaID int) ([]*Exam, error)
	
	// Statistics and maintenance
	UpdateQuestionsCount(examID int) error
}

// TopicRepository define a interface para operações de persistência de tópicos
type TopicRepository interface {
	// Basic CRUD operations
	Create(t *Topic) error
	Update(t *Topic) error
	Delete(id int) error
	FindByID(id int) (*Topic, error)
	ListAll() ([]*Topic, error)
	ListPaginated(page, pageSize int, examID *int) ([]*Topic, *Pagination, error)
	
	// Query operations (legacy support)
	ListByExam(examID int) ([]*Topic, error) // Deprecated: usar ExamTopicRepository
	
	// New N:N relationship operations
	ListByExamWithDetails(examID int) ([]*ExamTopicWithDetails, error)
}
