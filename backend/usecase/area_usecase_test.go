package usecase

import (
	"context"
	"testing"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// MockAreaRepository implementa AreaRepository para testes
type MockAreaRepository struct {
	areas  map[int]*domain.Area
	nextID int
}

func NewMockAreaRepository() *MockAreaRepository {
	return &MockAreaRepository{
		areas:  make(map[int]*domain.Area),
		nextID: 1,
	}
}

func (m *MockAreaRepository) Create(area *domain.Area) error {
	area.ID = m.nextID
	m.areas[area.ID] = area
	m.nextID++
	return nil
}

func (m *MockAreaRepository) Update(area *domain.Area) error {
	if _, exists := m.areas[area.ID]; !exists {
		return domain.ErrNotFound
	}
	m.areas[area.ID] = area
	return nil
}

func (m *MockAreaRepository) Delete(id int) error {
	if _, exists := m.areas[id]; !exists {
		return domain.ErrNotFound
	}
	delete(m.areas, id)
	return nil
}

func (m *MockAreaRepository) FindByID(id int) (*domain.Area, error) {
	if area, exists := m.areas[id]; exists {
		return area, nil
	}
	return nil, domain.ErrNotFound
}

func (m *MockAreaRepository) ListAll() ([]*domain.Area, error) {
	var areas []*domain.Area
	for _, area := range m.areas {
		areas = append(areas, area)
	}
	return areas, nil
}

func (m *MockAreaRepository) ListPaginated(page, pageSize int) ([]*domain.Area, *domain.Pagination, error) {
	var areas []*domain.Area
	for _, area := range m.areas {
		areas = append(areas, area)
	}
	
	// Implementação simples de paginação
	totalItems := len(areas)
	totalPages := (totalItems + pageSize - 1) / pageSize
	
	pagination := &domain.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
	
	return areas, pagination, nil
}

func TestAreaUsecase_Create(t *testing.T) {
	mockRepo := NewMockAreaRepository()
	uc := &AreaUsecase{Repo: mockRepo}

	tests := []struct {
		name        string
		areaName    string
		description string
		expectError bool
	}{
		{
			name:        "Criar área válida",
			areaName:    "TI - Tecnologia da Informação",
			description: "Certificações e provas de tecnologia",
			expectError: false,
		},
		{
			name:        "Criar área sem nome",
			areaName:    "",
			description: "Descrição válida",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			area, err := uc.CreateArea(context.Background(), tt.areaName, tt.description)

			if tt.expectError {
				if err == nil {
					t.Errorf("Esperava erro, mas não houve erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas houve: %v", err)
				}
				if area.ID == 0 {
					t.Errorf("Esperava que o ID fosse definido")
				}
			}
		})
	}
}

func TestAreaUsecase_Update(t *testing.T) {
	mockRepo := NewMockAreaRepository()
	uc := &AreaUsecase{Repo: mockRepo}

	// Criar uma área primeiro
	area, _ := uc.CreateArea(context.Background(), "TI", "Tecnologia da Informação")

	tests := []struct {
		name        string
		areaID      int
		newName     string
		expectError bool
	}{
		{
			name:        "Atualizar área existente",
			areaID:      area.ID,
			newName:     "TI - Tecnologia da Informação",
			expectError: false,
		},
		{
			name:        "Atualizar área inexistente",
			areaID:      999,
			newName:     "Inexistente",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.UpdateArea(context.Background(), tt.areaID, tt.newName, "Descrição atualizada")

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

func TestAreaUsecase_Delete(t *testing.T) {
	mockRepo := NewMockAreaRepository()
	uc := &AreaUsecase{Repo: mockRepo}

	// Criar uma área primeiro
	area, _ := uc.CreateArea(context.Background(), "TI", "Tecnologia da Informação")

	tests := []struct {
		name        string
		areaID      int
		expectError bool
	}{
		{
			name:        "Deletar área existente",
			areaID:      area.ID,
			expectError: false,
		},
		{
			name:        "Deletar área inexistente",
			areaID:      999,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.DeleteArea(context.Background(), tt.areaID)

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

func TestAreaUsecase_GetByID(t *testing.T) {
	mockRepo := NewMockAreaRepository()
	uc := &AreaUsecase{Repo: mockRepo}

	// Criar uma área primeiro
	createdArea, _ := uc.CreateArea(context.Background(), "TI", "Tecnologia da Informação")

	tests := []struct {
		name        string
		areaID      int
		expectError bool
	}{
		{
			name:        "Buscar área existente",
			areaID:      createdArea.ID,
			expectError: false,
		},
		{
			name:        "Buscar área inexistente",
			areaID:      999,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			area, err := uc.GetArea(context.Background(), tt.areaID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Esperava erro, mas não houve erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro, mas houve: %v", err)
				}
				if area == nil {
					t.Errorf("Esperava área, mas recebeu nil")
				}
				if area.ID != tt.areaID {
					t.Errorf("Esperava ID %d, mas recebeu %d", tt.areaID, area.ID)
				}
			}
		})
	}
}

func TestAreaUsecase_ListAll(t *testing.T) {
	mockRepo := NewMockAreaRepository()
	uc := &AreaUsecase{Repo: mockRepo}

	// Criar algumas áreas
	uc.CreateArea(context.Background(), "TI", "Tecnologia da Informação")
	uc.CreateArea(context.Background(), "Administração", "Administração e Gestão")

	areas, err := uc.ListAreas(context.Background())

	if err != nil {
		t.Errorf("Não esperava erro, mas houve: %v", err)
	}

	if len(areas) != 2 {
		t.Errorf("Esperava 2 áreas, mas recebeu %d", len(areas))
	}
}
