package domain

// PerformanceReport representa o relatório de desempenho do usuário.
type PerformanceReport struct {
	UserID         int
	TotalQuizzes   int
	TotalQuestions int
	CorrectAnswers int
	Accuracy       float64
}

// ReportRepository define a interface para operações de persistência de relatórios de desempenho
type ReportRepository interface {
	GetReport(userID int) (*PerformanceReport, error)
}
