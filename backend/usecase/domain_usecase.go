package usecase

import (
	"context"
	"errors"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// DomainUsecase define as operações de negócio para domínios
type DomainUsecase struct {
	Repo domain.DomainRepository
}

// NewDomainUsecase cria uma nova instância do DomainUsecase
func NewDomainUsecase(repo domain.DomainRepository) *DomainUsecase {
	return &DomainUsecase{
		Repo: repo,
	}
}

// ListDomains lista todos os domínios, opcionalmente filtrados por exame
func (uc *DomainUsecase) ListDomains(ctx context.Context, examID int) ([]*domain.Domain, error) {
	if examID > 0 {
		return uc.Repo.ListByExam(examID)
	}
	// Para listar todos, precisamos implementar um método ListAll ou usar uma abordagem diferente
	// Por enquanto, retornamos erro se não especificar examID
	return nil, errors.New("exam_id é obrigatório para listar domínios")
}

// GetDomain obtém um domínio específico por ID
func (uc *DomainUsecase) GetDomain(ctx context.Context, id int) (*domain.Domain, error) {
	return uc.Repo.FindByID(id)
}

// CreateDomain cria um novo domínio
func (uc *DomainUsecase) CreateDomain(ctx context.Context, examID int, name, description string, weightPercentage, orderIndex int) (*domain.Domain, error) {
	// Validações
	if name == "" {
		return nil, errors.New("nome é obrigatório")
	}
	if len(name) < 3 {
		return nil, errors.New("nome deve ter pelo menos 3 caracteres")
	}
	if description == "" {
		return nil, errors.New("descrição é obrigatória")
	}
	if len(description) < 10 {
		return nil, errors.New("descrição deve ter pelo menos 10 caracteres")
	}
	if weightPercentage < 0 || weightPercentage > 100 {
		return nil, errors.New("peso deve estar entre 0 e 100")
	}
	if orderIndex < 0 {
		return nil, errors.New("ordem deve ser maior ou igual a 0")
	}

	domain := &domain.Domain{
		ExamID:           examID,
		Name:             name,
		Description:      description,
		WeightPercentage: float64(weightPercentage),
		OrderIndex:       orderIndex,
	}

	err := uc.Repo.Create(domain)
	if err != nil {
		return nil, err
	}
	return domain, nil
}

// UpdateDomain atualiza um domínio existente
func (uc *DomainUsecase) UpdateDomain(ctx context.Context, id, examID int, name, description string, weightPercentage, orderIndex int) error {
	// Validações
	if name == "" {
		return errors.New("nome é obrigatório")
	}
	if len(name) < 3 {
		return errors.New("nome deve ter pelo menos 3 caracteres")
	}
	if description == "" {
		return errors.New("descrição é obrigatória")
	}
	if len(description) < 10 {
		return errors.New("descrição deve ter pelo menos 10 caracteres")
	}
	if weightPercentage < 0 || weightPercentage > 100 {
		return errors.New("peso deve estar entre 0 e 100")
	}
	if orderIndex < 0 {
		return errors.New("ordem deve ser maior ou igual a 0")
	}

	domain := &domain.Domain{
		ID:               id,
		ExamID:           examID,
		Name:             name,
		Description:      description,
		WeightPercentage: float64(weightPercentage),
		OrderIndex:       orderIndex,
	}

	return uc.Repo.Update(domain)
}

// DeleteDomain exclui um domínio
func (uc *DomainUsecase) DeleteDomain(ctx context.Context, id int) error {
	return uc.Repo.Delete(id)
}
