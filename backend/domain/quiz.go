package domain

import "time"

// Quiz representa um simulado
// UserQuiz representa o histórico de simulados realizados

type Quiz struct {
	ID          int
	Title       string
	Description string
	ThemeID     int
	CreatedBy   int
	Questions   []*Question
}

type QuizRepository interface {
	Create(q *Quiz) error
	Update(q *Quiz) error
	Delete(id int) error
	FindByID(id int) (*Quiz, error)
	ListAll(themeID *int) ([]*Quiz, error)
}

type UserQuiz struct {
	ID        int
	UserID    int
	QuizID    int
	StartedAt time.Time
	FinishedAt *time.Time
}

type UserQuizRepository interface {
	Create(uq *UserQuiz) error
	ListByUser(userID int) ([]*UserQuiz, error)
}
