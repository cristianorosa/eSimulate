package domain

// Theme representa um tema ou subtema de simulado
// Estrutura hierárquica (parent_id)
type Theme struct {
	ID       int
	Name     string
	ParentID *int
}

type ThemeRepository interface {
	Create(theme *Theme) error
	Update(theme *Theme) error
	Delete(id int) error
	FindByID(id int) (*Theme, error)
	ListAll() ([]*Theme, error)
}
