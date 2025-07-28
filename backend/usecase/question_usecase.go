package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// QuestionUsecase implementa as regras de negócio para questões

// QuestionUsecase implementa as regras de negócio para questões.
type QuestionUsecase struct {
	Repo domain.QuestionRepository
}

// CreateQuestion cria uma nova questão com suas opções.
func (uc *QuestionUsecase) CreateQuestion(ctx context.Context, themeID int, statement, explanation string, createdBy int, options []*domain.Option) (*domain.Question, error) {
	statement = strings.TrimSpace(statement)
	if statement == "" || len(options) < 2 {
		return nil, errors.New("enunciado e pelo menos duas opções são obrigatórios")
	}
	var hasCorrect bool
	for _, o := range options {
		if o.IsCorrect {
			hasCorrect = true
			break
		}
	}
	if !hasCorrect {
		return nil, errors.New("pelo menos uma opção deve ser correta")
	}
	q := &domain.Question{
		ThemeID:     themeID,
		Statement:   statement,
		Explanation: explanation,
		CreatedBy:   createdBy,
		Options:     options,
	}
	if err := uc.Repo.Create(q); err != nil {
		return nil, err
	}
	return q, nil
}

// UpdateQuestion atualiza uma questão existente.
func (uc *QuestionUsecase) UpdateQuestion(ctx context.Context, id, themeID int, statement, explanation string) error {
	statement = strings.TrimSpace(statement)
	if statement == "" {
		return errors.New("enunciado é obrigatório")
	}
	q := &domain.Question{ID: id, ThemeID: themeID, Statement: statement, Explanation: explanation}
	return uc.Repo.Update(q)
}

// DeleteQuestion remove uma questão.
func (uc *QuestionUsecase) DeleteQuestion(ctx context.Context, id int) error {
	return uc.Repo.Delete(id)
}

// ListQuestions lista todas as questões, opcionalmente filtradas por tema.
func (uc *QuestionUsecase) ListQuestions(ctx context.Context, themeID *int) ([]*domain.Question, error) {
	return uc.Repo.ListAll(themeID)
}

// GetQuestion busca uma questão pelo seu ID.
func (uc *QuestionUsecase) GetQuestion(ctx context.Context, id int) (*domain.Question, error) {
	return uc.Repo.FindByID(id)
}
