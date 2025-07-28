package usecase

import (
	"context"
	"github.com/cristianorosa/eSimulate/backend/domain"
	"time"
)

// UserQuizUsecase implementa as regras de negócio para histórico de simulados

// UserQuizUsecase implementa as regras de negócio para histórico de simulados.
type UserQuizUsecase struct {
	Repo domain.UserQuizRepository
}

// RegisterUserQuiz registra um novo simulado realizado por um usuário.
func (uc *UserQuizUsecase) RegisterUserQuiz(ctx context.Context, userID, quizID int) (*domain.UserQuiz, error) {
	uq := &domain.UserQuiz{
		UserID:    userID,
		QuizID:    quizID,
		StartedAt: time.Now(),
	}
	if err := uc.Repo.Create(uq); err != nil {
		return nil, err
	}
	return uq, nil
}

// ListUserHistory lista o histórico de simulados de um usuário.
func (uc *UserQuizUsecase) ListUserHistory(ctx context.Context, userID int) ([]*domain.UserQuiz, error) {
	return uc.Repo.ListByUser(userID)
}
