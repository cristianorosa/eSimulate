package usecase

import (
	"context"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// AreaUsecase implementa as regras de negócio para áreas
type AreaUsecase struct {
	Repo domain.AreaRepository
}

// CreateArea cria uma nova área
func (uc *AreaUsecase) CreateArea(ctx context.Context, name, description string) (*domain.Area, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	area := &domain.Area{
		Name:        name,
		Description: strings.TrimSpace(description),
	}

	if err := uc.Repo.Create(area); err != nil {
		return nil, err
	}

	return area, nil
}

// UpdateArea atualiza uma área existente
func (uc *AreaUsecase) UpdateArea(ctx context.Context, id int, name, description string) error {
	if id == 0 {
		return domain.ErrInvalidInput
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return domain.ErrInvalidInput
	}

	area := &domain.Area{
		ID:          id,
		Name:        name,
		Description: strings.TrimSpace(description),
	}

	return uc.Repo.Update(area)
}

// DeleteArea remove uma área
func (uc *AreaUsecase) DeleteArea(ctx context.Context, id int) error {
	if id == 0 {
		return domain.ErrInvalidInput
	}
	return uc.Repo.Delete(id)
}

// GetArea busca uma área por ID
func (uc *AreaUsecase) GetArea(ctx context.Context, id int) (*domain.Area, error) {
	if id == 0 {
		return nil, domain.ErrInvalidInput
	}
	return uc.Repo.FindByID(id)
}

// ListAreas lista todas as áreas
func (uc *AreaUsecase) ListAreas(ctx context.Context) ([]*domain.Area, error) {
	return uc.Repo.ListAll()
}

// ListAreasPaginated lista áreas com paginação
func (uc *AreaUsecase) ListAreasPaginated(ctx context.Context, page, pageSize int) ([]*domain.Area, *domain.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return uc.Repo.ListPaginated(page, pageSize)
}
