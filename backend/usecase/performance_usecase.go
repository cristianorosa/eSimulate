package usecase

import (
	"errors"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// PerformanceUsecase implementa as regras de negócio para relatórios de desempenho
type PerformanceUsecase struct {
	Repo domain.ReportRepository
}

// GetReport obtém o relatório de desempenho de um usuário
func (uc *PerformanceUsecase) GetReport(userID int) (*domain.PerformanceReport, error) {
	if userID <= 0 {
		return nil, errors.New("user_id deve ser maior que zero")
	}

	return uc.Repo.GetReport(userID)
}
