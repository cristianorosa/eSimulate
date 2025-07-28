package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// ThemeUsecase implementa as regras de negócio para temas

// ThemeUsecase implementa as regras de negócio para temas.
type ThemeUsecase struct {
	Repo domain.ThemeRepository
}

// CreateTheme cria um novo tema.
func (uc *ThemeUsecase) CreateTheme(ctx context.Context, name string, parentID *int) (*domain.Theme, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("nome do tema é obrigatório")
	}
	theme := &domain.Theme{Name: name, ParentID: parentID}
	if err := uc.Repo.Create(theme); err != nil {
		return nil, err
	}
	return theme, nil
}

// UpdateTheme atualiza um tema existente.
func (uc *ThemeUsecase) UpdateTheme(ctx context.Context, id int, name string, parentID *int) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nome do tema é obrigatório")
	}
	theme := &domain.Theme{ID: id, Name: name, ParentID: parentID}
	return uc.Repo.Update(theme)
}

// DeleteTheme remove um tema.
func (uc *ThemeUsecase) DeleteTheme(ctx context.Context, id int) error {
	return uc.Repo.Delete(id)
}

// ListThemes lista todos os temas.
func (uc *ThemeUsecase) ListThemes(ctx context.Context) ([]*domain.Theme, error) {
	return uc.Repo.ListAll()
}

// GetTheme busca um tema pelo seu ID.
func (uc *ThemeUsecase) GetTheme(ctx context.Context, id int) (*domain.Theme, error) {
	return uc.Repo.FindByID(id)
}
