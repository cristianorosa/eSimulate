package domain

import "time"

// Question representa uma questão de prova
type Question struct {
	ID              int
	ExamID          int
	DomainID        int
	Statement       string
	Explanation     string
	DifficultyLevel int
	CreatedBy       int
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Options         []*Option
}

// Option representa uma opção de resposta para uma questão
type Option struct {
	ID          int
	QuestionID  int
	Text        string
	IsCorrect   bool
	Explanation string
	OrderIndex  int
}

// QuestionRepository define a interface para operações de persistência de questões
type QuestionRepository interface {
	Create(q *Question) error
	Update(q *Question) error
	Delete(id int) error
	FindByID(id int) (*Question, error)
	ListByExam(examID int) ([]*Question, error)
	ListByDomain(domainID int) ([]*Question, error)
	ListActive() ([]*Question, error)
}
