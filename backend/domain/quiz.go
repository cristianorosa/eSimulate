package domain

import "time"

// Quiz representa um simulado
// UserQuiz representa o histórico de simulados realizados

// Quiz representa um simulado
type Quiz struct {
	ID          int
	Title       string
	Description string
	ThemeID     int
	CreatedBy   int
	Questions   []*Question
}

// QuizRepository define a interface para operações de persistência de simulados
type QuizRepository interface {
	Create(q *Quiz) error
	Update(q *Quiz) error
	Delete(id int) error
	FindByID(id int) (*Quiz, error)
	ListAll(themeID *int) ([]*Quiz, error)
}

// UserQuiz representa o histórico de simulados realizados por um usuário
type UserQuiz struct {
	ID        int
	UserID    int
	QuizID    int
	StartedAt time.Time
	FinishedAt *time.Time
}

// UserQuizRepository define a interface para operações de persistência de histórico de simulados de usuários
type UserQuizRepository interface {
	Create(uq *UserQuiz) error
	ListByUser(userID int) ([]*UserQuiz, error)
}
