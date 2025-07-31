package usecase

import (
	"context"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// ExamUsecase implementa as regras de negócio para exames
type ExamUsecase struct {
	Repo domain.ExamRepository
}

// CreateExam cria um novo exame
func (uc *ExamUsecase) CreateExam(ctx context.Context, title, description string, areaID int, maxTimeMinutes int, passingScore float64, isActive bool, createdBy int) (*domain.Exam, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, domain.ErrInvalidInput
	}

	if areaID == 0 {
		return nil, domain.ErrInvalidInput
	}

	if maxTimeMinutes <= 0 {
		return nil, domain.ErrInvalidInput
	}

	if passingScore < 0 || passingScore > 100 {
		return nil, domain.ErrInvalidInput
	}

	exam := &domain.Exam{
		Title:          title,
		Description:    strings.TrimSpace(description),
		AreaID:         areaID,
		MaxTimeMinutes: maxTimeMinutes,
		PassingScore:   passingScore,
		IsActive:       isActive,
		CreatedBy:      createdBy,
	}

	if err := uc.Repo.Create(exam); err != nil {
		return nil, err
	}

	return exam, nil
}

// UpdateExam atualiza um exame existente
func (uc *ExamUsecase) UpdateExam(ctx context.Context, id int, title, description string, areaID int, maxTimeMinutes int, passingScore float64, isActive bool) error {
	if id == 0 {
		return domain.ErrInvalidInput
	}

	title = strings.TrimSpace(title)
	if title == "" {
		return domain.ErrInvalidInput
	}

	if areaID == 0 {
		return domain.ErrInvalidInput
	}

	if maxTimeMinutes <= 0 {
		return domain.ErrInvalidInput
	}

	if passingScore < 0 || passingScore > 100 {
		return domain.ErrInvalidInput
	}

	exam := &domain.Exam{
		ID:             id,
		Title:          title,
		Description:    strings.TrimSpace(description),
		AreaID:         areaID,
		MaxTimeMinutes: maxTimeMinutes,
		PassingScore:   passingScore,
		IsActive:       isActive,
	}

	return uc.Repo.Update(exam)
}

// DeleteExam remove um exame
func (uc *ExamUsecase) DeleteExam(ctx context.Context, id int) error {
	if id == 0 {
		return domain.ErrInvalidInput
	}
	return uc.Repo.Delete(id)
}

// GetExam busca um exame por ID
func (uc *ExamUsecase) GetExam(ctx context.Context, id int) (*domain.Exam, error) {
	if id == 0 {
		return nil, domain.ErrInvalidInput
	}
	return uc.Repo.FindByID(id)
}

// ListExams lista todos os exames
func (uc *ExamUsecase) ListExams(ctx context.Context) ([]*domain.Exam, error) {
	return uc.Repo.ListAll()
}

// ListExamsByArea lista exames por área
func (uc *ExamUsecase) ListExamsByArea(ctx context.Context, areaID int) ([]*domain.Exam, error) {
	if areaID == 0 {
		return nil, domain.ErrInvalidInput
	}
	return uc.Repo.ListByArea(areaID)
}

// ListExamsPaginated lista exames com paginação
func (uc *ExamUsecase) ListExamsPaginated(ctx context.Context, page, pageSize int, areaID *int) ([]*domain.Exam, *domain.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return uc.Repo.ListPaginated(page, pageSize, areaID)
}
