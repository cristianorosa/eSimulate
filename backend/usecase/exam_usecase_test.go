package usecase

import (
	"context"
	"testing"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// MockExamRepository implementa ExamRepository para testes
type MockExamRepository struct {
	exams  map[int]*domain.Exam
	nextID int
}

func NewMockExamRepository() *MockExamRepository {
	return &MockExamRepository{
		exams:  make(map[int]*domain.Exam),
		nextID: 1,
	}
}

func (m *MockExamRepository) Create(exam *domain.Exam) error {
	exam.ID = m.nextID
	m.exams[exam.ID] = exam
	m.nextID++
	return nil
}

func (m *MockExamRepository) Update(exam *domain.Exam) error {
	if _, exists := m.exams[exam.ID]; !exists {
		return domain.ErrNotFound
	}
	m.exams[exam.ID] = exam
	return nil
}

func (m *MockExamRepository) Delete(id int) error {
	if _, exists := m.exams[id]; !exists {
		return domain.ErrNotFound
	}
	delete(m.exams, id)
	return nil
}

func (m *MockExamRepository) FindByID(id int) (*domain.Exam, error) {
	if exam, exists := m.exams[id]; exists {
		return exam, nil
	}
	return nil, domain.ErrNotFound
}

func (m *MockExamRepository) ListAll() ([]*domain.Exam, error) {
	var exams []*domain.Exam
	for _, exam := range m.exams {
		exams = append(exams, exam)
	}
	return exams, nil
}

func (m *MockExamRepository) ListByArea(areaID int) ([]*domain.Exam, error) {
	var exams []*domain.Exam
	for _, exam := range m.exams {
		if exam.AreaID == areaID {
			exams = append(exams, exam)
		}
	}
	return exams, nil
}

func (m *MockExamRepository) ListPaginated(page, pageSize int, areaID *int) ([]*domain.Exam, *domain.Pagination, error) {
	var exams []*domain.Exam
	for _, exam := range m.exams {
		if areaID == nil || exam.AreaID == *areaID {
			exams = append(exams, exam)
		}
	}
	
	// Implementação simples de paginação
	totalItems := len(exams)
	totalPages := (totalItems + pageSize - 1) / pageSize
	
	pagination := &domain.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
	
	return exams, pagination, nil
}

func TestExamUsecase_Create(t *testing.T) {
	mockRepo := NewMockExamRepository()
	uc := &ExamUsecase{Repo: mockRepo}

	tests := []struct {
		name           string
		title          string
		areaID         int
		maxTimeMinutes int
		expectError    bool
	}{
		{
			name:           "Criar exame válido",
			title:          "AWS Certified Cloud Practitioner",
			areaID:         1,
			maxTimeMinutes: 90,
			expectError:    false,
		},
		{
			name:           "Criar exame sem título",
			title:          "",
			areaID:         1,
			maxTimeMinutes: 90,
			expectError:    true,
		},
		{
			name:           "Criar exame sem área",
			title:          "Exame válido",
			areaID:         0,
			maxTimeMinutes: 90,
			expectError:    true,
		},
		{
			name:           "Criar exame com tempo inválido",
			title:          "Exame válido",
			areaID:         1,
			maxTimeMinutes: 0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exam, err := uc.CreateExam(context.Background(), tt.title, "Descrição do exame", tt.areaID, tt.maxTimeMinutes, 70.0, true, 1)

			if tt.expectError {
				if err == nil {
					t.Errorf("Esperava erro, mas não houve erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas houve: %v", err)
				}
				if exam.ID == 0 {
					t.Errorf("Esperava que o ID fosse definido")
				}
			}
		})
	}
}

func TestExamUsecase_Update(t *testing.T) {
	mockRepo := NewMockExamRepository()
	uc := &ExamUsecase{Repo: mockRepo}

	// Criar um exame primeiro
	exam, _ := uc.CreateExam(context.Background(), "Exame Original", "Descrição original", 1, 90, 70.0, true, 1)

	tests := []struct {
		name        string
		examID      int
		newTitle    string
		expectError bool
	}{
		{
			name:        "Atualizar exame existente",
			examID:      exam.ID,
			newTitle:    "Exame Atualizado",
			expectError: false,
		},
		{
			name:        "Atualizar exame inexistente",
			examID:      999,
			newTitle:    "Inexistente",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.UpdateExam(context.Background(), tt.examID, tt.newTitle, "Descrição atualizada", 1, 90, 70.0, true)

			if tt.expectError {
				if err == nil {
					t.Errorf("Esperava erro, mas não houve erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas houve: %v", err)
				}
			}
		})
	}
}

func TestExamUsecase_Delete(t *testing.T) {
	mockRepo := NewMockExamRepository()
	uc := &ExamUsecase{Repo: mockRepo}

	// Criar um exame primeiro
	exam, _ := uc.CreateExam(context.Background(), "Exame para deletar", "Descrição", 1, 90, 70.0, true, 1)

	tests := []struct {
		name        string
		examID      int
		expectError bool
	}{
		{
			name:        "Deletar exame existente",
			examID:      exam.ID,
			expectError: false,
		},
		{
			name:        "Deletar exame inexistente",
			examID:      999,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.DeleteExam(context.Background(), tt.examID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Esperava erro, mas não houve erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas houve: %v", err)
				}
			}
		})
	}
}

func TestExamUsecase_GetByID(t *testing.T) {
	mockRepo := NewMockExamRepository()
	uc := &ExamUsecase{Repo: mockRepo}

	// Criar um exame primeiro
	createdExam, _ := uc.CreateExam(context.Background(), "Exame para buscar", "Descrição", 1, 90, 70.0, true, 1)

	tests := []struct {
		name        string
		examID      int
		expectError bool
	}{
		{
			name:        "Buscar exame existente",
			examID:      createdExam.ID,
			expectError: false,
		},
		{
			name:        "Buscar exame inexistente",
			examID:      999,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exam, err := uc.GetExam(context.Background(), tt.examID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Esperava erro, mas não houve erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas houve: %v", err)
				}
				if exam == nil {
					t.Errorf("Esperava exame, mas recebeu nil")
				}
				if exam.ID != tt.examID {
					t.Errorf("Esperava ID %d, mas recebeu %d", tt.examID, exam.ID)
				}
			}
		})
	}
}

func TestExamUsecase_ListAll(t *testing.T) {
	mockRepo := NewMockExamRepository()
	uc := &ExamUsecase{Repo: mockRepo}

	// Criar alguns exames
	uc.CreateExam(context.Background(), "Exame 1", "Descrição 1", 1, 90, 70.0, true, 1)
	uc.CreateExam(context.Background(), "Exame 2", "Descrição 2", 1, 90, 70.0, true, 1)

	exams, err := uc.ListExams(context.Background())

	if err != nil {
		t.Errorf("Não esperava erro, mas houve: %v", err)
	}

	if len(exams) != 2 {
		t.Errorf("Esperava 2 exames, mas recebeu %d", len(exams))
	}
}

func TestExamUsecase_ListActive(t *testing.T) {
	mockRepo := NewMockExamRepository()
	uc := &ExamUsecase{Repo: mockRepo}

	// Criar exames ativos e inativos
	uc.CreateExam(context.Background(), "Exame Ativo 1", "Descrição", 1, 90, 70.0, true, 1)
	uc.CreateExam(context.Background(), "Exame Ativo 2", "Descrição", 1, 90, 70.0, true, 1)
	uc.CreateExam(context.Background(), "Exame Inativo", "Descrição", 1, 90, 70.0, false, 1)

	exams, err := uc.ListExams(context.Background())

	if err != nil {
		t.Errorf("Não esperava erro, mas houve: %v", err)
	}

	if len(exams) != 3 {
		t.Errorf("Esperava 3 exames, mas recebeu %d", len(exams))
	}
}
