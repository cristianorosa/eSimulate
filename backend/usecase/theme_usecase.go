package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// ThemeUsecase implementa as regras de negócio para temas

type ThemeUsecase struct {
	Repo domain.ThemeRepository
}

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

func (uc *ThemeUsecase) UpdateTheme(ctx context.Context, id int, name string, parentID *int) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nome do tema é obrigatório")
	}
	theme := &domain.Theme{ID: id, Name: name, ParentID: parentID}
	return uc.Repo.Update(theme)
}

func (uc *ThemeUsecase) DeleteTheme(ctx context.Context, id int) error {
	return uc.Repo.Delete(id)
}

func (uc *ThemeUsecase) ListThemes(ctx context.Context) ([]*domain.Theme, error) {
	return uc.Repo.ListAll()
}

func (uc *ThemeUsecase) GetTheme(ctx context.Context, id int) (*domain.Theme, error) {
	return uc.Repo.FindByID(id)
}
