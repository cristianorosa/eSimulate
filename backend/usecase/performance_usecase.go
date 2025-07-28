package usecase

import (
	"context"
	"github.com/cristianorosa/eSimulate/backend/domain"
)

// PerformanceUsecase implementa a lógica de geração de relatório de desempenho

// PerformanceUsecase implementa a lógica de geração de relatório de desempenho.
type PerformanceUsecase struct {
	Repo domain.PerformanceRepository
}

// GetUserPerformance obtém o relatório de desempenho de um usuário.
func (uc *PerformanceUsecase) GetUserPerformance(ctx context.Context, userID int) (*domain.PerformanceReport, error) {
	return uc.Repo.GetReport(userID)
}
