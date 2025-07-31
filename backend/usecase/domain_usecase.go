package usecase

import (
	"context"
	"errors"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// DomainUsecase implementa as regras de negócio para domínios
type DomainUsecase struct {
	Repo domain.Repository
}

// NewDomainUsecase cria uma nova instância de DomainUsecase
func NewDomainUsecase(repo domain.Repository) *DomainUsecase {
	return &DomainUsecase{Repo: repo}
}

// CreateDomain cria um novo domínio
func (uc *DomainUsecase) CreateDomain(ctx context.Context, examID int, name, description string, weightPercentage float64, orderIndex int) (*domain.Domain, error) {
	if name == "" {
		return nil, errors.New("nome do domínio é obrigatório")
	}

	if weightPercentage < 0 || weightPercentage > 100 {
		return nil, errors.New("peso do domínio deve estar entre 0 e 100")
	}

	domain := &domain.Domain{
		ExamID:           examID,
		Name:             name,
		Description:      description,
		WeightPercentage: weightPercentage,
		OrderIndex:       orderIndex,
	}

	err := uc.Repo.Create(domain)
	if err != nil {
		return nil, err
	}

	return domain, nil
}

// UpdateDomain atualiza um domínio existente
func (uc *DomainUsecase) UpdateDomain(ctx context.Context, id int, name, description string, weightPercentage float64, orderIndex int) error {
	if name == "" {
		return errors.New("nome do domínio é obrigatório")
	}

	if weightPercentage < 0 || weightPercentage > 100 {
		return errors.New("peso do domínio deve estar entre 0 e 100")
	}

	domain, err := uc.Repo.FindByID(id)
	if err != nil {
		return err
	}

	domain.Name = name
	domain.Description = description
	domain.WeightPercentage = weightPercentage
	domain.OrderIndex = orderIndex

	return uc.Repo.Update(domain)
}

// DeleteDomain remove um domínio
func (uc *DomainUsecase) DeleteDomain(ctx context.Context, id int) error {
	return uc.Repo.Delete(id)
}

// GetDomain obtém um domínio por ID
func (uc *DomainUsecase) GetDomain(ctx context.Context, id int) (*domain.Domain, error) {
	return uc.Repo.FindByID(id)
}

// ListDomainsByExam lista todos os domínios de um exame
func (uc *DomainUsecase) ListDomainsByExam(ctx context.Context, examID int) ([]*domain.Domain, error) {
	return uc.Repo.ListByExam(examID)
}
