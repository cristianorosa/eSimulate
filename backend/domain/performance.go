package domain

import "time"

// Performance representa o desempenho de um usuário em um exame
type Performance struct {
	ID                int       `json:"id"`
	UserExamID        int       `json:"user_exam_id"`
	DomainID          int       `json:"domain_id"`
	QuestionsAnswered int       `json:"questions_answered"`
	CorrectAnswers    int       `json:"correct_answers"`
	ScorePercentage   float64   `json:"score_percentage"`
	NeedsImprovement  bool      `json:"needs_improvement"`
	CreatedAt         time.Time `json:"created_at"`
}

// TopicPerformance representa o desempenho de um usuário em um tópico específico
type TopicPerformance struct {
	ID             int       `json:"id"`
	UserExamID     int       `json:"user_exam_id"`
	TopicID        int       `json:"topic_id"`
	CorrectAnswers int       `json:"correct_answers"`
	TotalQuestions int       `json:"total_questions"`
	Score          float64   `json:"score"`
	CreatedAt      time.Time `json:"created_at"`
}

// PerformanceReport representa o relatório de desempenho do usuário
type PerformanceReport struct {
	UserID         int     `json:"user_id"`
	TotalQuizzes   int     `json:"total_quizzes"`
	TotalQuestions int     `json:"total_questions"`
	CorrectAnswers int     `json:"correct_answers"`
	Accuracy       float64 `json:"accuracy"`
}

// PerformanceRepository define as operações de persistência para desempenho
type PerformanceRepository interface {
	Create(performance *Performance) error
	Update(performance *Performance) error
	ListByUserExam(userExamID int) ([]*Performance, error)
	GetReport(userID int) (*PerformanceReport, error)
}

// TopicPerformanceRepository define as operações de persistência para desempenho por tópico
type TopicPerformanceRepository interface {
	Create(performance *TopicPerformance) error
	Update(performance *TopicPerformance) error
	FindByUserExam(userExamID int) ([]*TopicPerformance, error)
}
