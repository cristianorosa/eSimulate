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
	DomainPerformance []*DomainPerformance
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

// DomainPerformance representa o desempenho do usuário em um domínio específico
type DomainPerformance struct {
	ID                int
	UserExamID        int
	DomainID          int
	QuestionsAnswered int
	CorrectAnswers    int
	ScorePercentage   float64
	NeedsImprovement  bool
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

// DomainPerformanceRepository define a interface para operações de persistência de desempenho por domínio
type DomainPerformanceRepository interface {
	Create(dp *DomainPerformance) error
	Update(dp *DomainPerformance) error
	ListByUserExam(userExamID int) ([]*DomainPerformance, error)
}
