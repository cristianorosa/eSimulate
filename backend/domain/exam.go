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
}

// Topic representa um tópico de uma prova
type Topic struct {
	ID               int       `json:"id"`
	ExamID           int       `json:"exam_id"`
	Name             string    `json:"name"`
	WeightPercentage float64   `json:"weight_percentage"`
	OrderIndex       int       `json:"order_index"`
	QuestionsCount   int       `json:"questions_count"`
	CreatedAt        time.Time `json:"created_at"`
}

// ExamRepository define as operações de persistência para exames
type ExamRepository interface {
	Create(exam *Exam) error
	Update(exam *Exam) error
	Delete(id int) error
	FindByID(id int) (*Exam, error)
	ListAll() ([]*Exam, error)
	ListByArea(areaID int) ([]*Exam, error)
	ListPaginated(page, pageSize int, areaID *int) ([]*Exam, *Pagination, error)
}

// TopicRepository define a interface para operações de persistência de tópicos
type TopicRepository interface {
	Create(t *Topic) error
	Update(t *Topic) error
	Delete(id int) error
	FindByID(id int) (*Topic, error)
	ListByExam(examID int) ([]*Topic, error)
	ListAll() ([]*Topic, error)
	ListPaginated(page, pageSize int, examID *int) ([]*Topic, *Pagination, error)
}
