package domain

import "time"

// UserExam representa a aplicação de uma prova por um usuário
type UserExam struct {
	ID                int
	UserID            int
	ExamID            int
	StartedAt         time.Time
	FinishedAt        *time.Time
	TotalScore        *float64
	Passed            *bool
	TimeSpentMinutes  *int
	Answers           []*UserAnswer
	TopicPerformance  []*TopicPerformance
}

// UserAnswer representa uma resposta do usuário em uma prova
type UserAnswer struct {
	ID                int
	UserExamID        int
	QuestionID        int
	OptionID          *int
	IsCorrect         *bool
	IsMarkedForReview bool
	AnsweredAt        time.Time
}

// UserExamRepository define a interface para operações de persistência de provas aplicadas
type UserExamRepository interface {
	Create(ue *UserExam) error
	Update(ue *UserExam) error
	FindByID(id int) (*UserExam, error)
	ListByUser(userID int) ([]*UserExam, error)
	ListByExam(examID int) ([]*UserExam, error)
	FindActiveByUser(userID int, examID int) (*UserExam, error)
}

// UserAnswerRepository define a interface para operações de persistência de respostas
type UserAnswerRepository interface {
	Create(ua *UserAnswer) error
	Update(ua *UserAnswer) error
	ListByUserExam(userExamID int) ([]*UserAnswer, error)
	FindByQuestion(userExamID int, questionID int) (*UserAnswer, error)
}
