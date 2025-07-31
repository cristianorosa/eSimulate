package domain

import "time"

// Area representa uma área de conhecimento
type Area struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// Pagination representa informações de paginação
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse representa uma resposta paginada
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// AreaRepository define as operações de persistência para áreas
type AreaRepository interface {
	Create(area *Area) error
	Update(area *Area) error
	Delete(id int) error
	FindByID(id int) (*Area, error)
	ListAll() ([]*Area, error)
	ListPaginated(page, pageSize int) ([]*Area, *Pagination, error)
}
