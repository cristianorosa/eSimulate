package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// QuizUsecase implementa as regras de negócio para simulados

// QuizUsecase implementa as regras de negócio para simulados.
type QuizUsecase struct {
	Repo domain.QuizRepository
}

// CreateQuiz cria um novo simulado.
func (uc *QuizUsecase) CreateQuiz(ctx context.Context, title, description string, themeID, createdBy int, questionIDs []int) (*domain.Quiz, error) {
	title = strings.TrimSpace(title)
	if title == "" || len(questionIDs) == 0 {
		return nil, errors.New("título e pelo menos uma questão são obrigatórios")
	}
	// Monta slice de questões apenas com IDs
	var questions []*domain.Question
	for _, qid := range questionIDs {
		questions = append(questions, &domain.Question{ID: qid})
	}
	quiz := &domain.Quiz{
		Title:       title,
		Description: description,
		ThemeID:     themeID,
		CreatedBy:   createdBy,
		Questions:   questions,
	}
	if err := uc.Repo.Create(quiz); err != nil {
		return nil, err
	}
	return quiz, nil
}

// UpdateQuiz atualiza um simulado existente.
func (uc *QuizUsecase) UpdateQuiz(ctx context.Context, id int, title, description string, themeID int) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New("título é obrigatório")
	}
	quiz := &domain.Quiz{ID: id, Title: title, Description: description, ThemeID: themeID}
	return uc.Repo.Update(quiz)
}

// DeleteQuiz remove um simulado.
func (uc *QuizUsecase) DeleteQuiz(ctx context.Context, id int) error {
	return uc.Repo.Delete(id)
}

// ListQuizzes lista todos os simulados, opcionalmente filtrados por tema.
func (uc *QuizUsecase) ListQuizzes(ctx context.Context, themeID *int) ([]*domain.Quiz, error) {
	return uc.Repo.ListAll(themeID)
}

// GetQuiz busca um simulado pelo seu ID.
func (uc *QuizUsecase) GetQuiz(ctx context.Context, id int) (*domain.Quiz, error) {
	return uc.Repo.FindByID(id)
}
