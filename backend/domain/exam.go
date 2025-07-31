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
	IsActive       bool      `json:"is_active"`
	CreatedBy      int       `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// Domain representa um domínio de uma prova
type Domain struct {
	ID               int
	ExamID           int
	Name             string
	Description      string
	WeightPercentage float64
	OrderIndex       int
	CreatedAt        time.Time
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

// Repository define a interface para operações de persistência de domínios
type Repository interface {
	Create(d *Domain) error
	Update(d *Domain) error
	Delete(id int) error
	FindByID(id int) (*Domain, error)
	ListByExam(examID int) ([]*Domain, error)
	ListAll() ([]*Domain, error)
}
