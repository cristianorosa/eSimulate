package usecase

import (
	"context"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// QuestionUsecase implementa as regras de negócio para questões.
type QuestionUsecase struct {
	Repo domain.QuestionRepository
}

// CreateQuestion cria uma nova questão com suas opções.
func (uc *QuestionUsecase) CreateQuestion(ctx context.Context, examID, topicID int, statement, explanation string, difficultyLevel, createdBy int, options []*domain.Option) (*domain.Question, error) {
	return &domain.Question{}, nil
}

// UpdateQuestion atualiza uma questão existente.
func (uc *QuestionUsecase) UpdateQuestion(ctx context.Context, id, examID, topicID int, statement, explanation string, difficultyLevel int) error {
	return nil
}

// DeleteQuestion remove uma questão.
func (uc *QuestionUsecase) DeleteQuestion(ctx context.Context, id int) error {
	return nil
}

// ListQuestions lista todas as questões, opcionalmente filtradas por exame.
func (uc *QuestionUsecase) ListQuestions(ctx context.Context, examID *int) ([]*domain.Question, error) {
	return []*domain.Question{}, nil
}

// GetQuestion busca uma questão pelo seu ID.
func (uc *QuestionUsecase) GetQuestion(ctx context.Context, id int) (*domain.Question, error) {
	return &domain.Question{}, nil
}
