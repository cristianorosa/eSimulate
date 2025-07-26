package infra

import (
	"database/sql"
	"eSimulate/backend/domain"
)

// PerformanceRepositoryPG implementa PerformanceRepository usando PostgreSQL

type PerformanceRepositoryPG struct {
	DB *sql.DB
}

func (r *PerformanceRepositoryPG) GetReport(userID int) (*domain.PerformanceReport, error) {
	var report domain.PerformanceReport
	report.UserID = userID
	// Total de simulados realizados
	err := r.DB.QueryRow(`SELECT COUNT(*) FROM user_quiz WHERE user_id = $1`, userID).Scan(&report.TotalQuizzes)
	if err != nil {
		return nil, err
	}
	// Total de questões respondidas
	err = r.DB.QueryRow(`SELECT COUNT(*) FROM user_answers ua JOIN user_quiz uq ON ua.user_quiz_id = uq.id WHERE uq.user_id = $1`, userID).Scan(&report.TotalQuestions)
	if err != nil {
		return nil, err
	}
	// Total de acertos
	err = r.DB.QueryRow(`SELECT COUNT(*) FROM user_answers ua JOIN options o ON ua.option_id = o.id JOIN user_quiz uq ON ua.user_quiz_id = uq.id WHERE uq.user_id = $1 AND o.is_correct = true`, userID).Scan(&report.CorrectAnswers)
	if err != nil {
		return nil, err
	}
	if report.TotalQuestions > 0 {
		report.Accuracy = float64(report.CorrectAnswers) / float64(report.TotalQuestions) * 100.0
	}
	return &report, nil
}
