package domain

// PerformanceReport representa o relatório de desempenho do usuário
// Pode ser expandido conforme necessidade

type PerformanceReport struct {
	UserID         int
	TotalQuizzes   int
	TotalQuestions int
	CorrectAnswers int
	Accuracy       float64
}

type PerformanceRepository interface {
	GetReport(userID int) (*PerformanceReport, error)
}
