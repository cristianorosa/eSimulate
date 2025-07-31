package infra

import (
	"database/sql"
	"math"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// AreaRepositoryPG implementa AreaRepository para PostgreSQL
type AreaRepositoryPG struct {
	DB *sql.DB
}

// Create cria uma nova área
func (r *AreaRepositoryPG) Create(area *domain.Area) error {
	query := `
		INSERT INTO areas (name, description, created_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		RETURNING id, created_at`

	err := r.DB.QueryRow(query, area.Name, area.Description).Scan(&area.ID, &area.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Update atualiza uma área existente
func (r *AreaRepositoryPG) Update(area *domain.Area) error {
	query := `
		UPDATE areas 
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3`

	result, err := r.DB.Exec(query, area.Name, area.Description, area.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Delete remove uma área
func (r *AreaRepositoryPG) Delete(id int) error {
	query := `DELETE FROM areas WHERE id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// FindByID busca uma área por ID
func (r *AreaRepositoryPG) FindByID(id int) (*domain.Area, error) {
	query := `
		SELECT id, name, description, created_at
		FROM areas
		WHERE id = $1`

	area := &domain.Area{}
	err := r.DB.QueryRow(query, id).Scan(
		&area.ID,
		&area.Name,
		&area.Description,
		&area.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return area, nil
}

// ListAll lista todas as áreas
func (r *AreaRepositoryPG) ListAll() ([]*domain.Area, error) {
	query := `
		SELECT id, name, description, created_at
		FROM areas
		ORDER BY name`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var areas []*domain.Area
	for rows.Next() {
		area := &domain.Area{}
		err := rows.Scan(
			&area.ID,
			&area.Name,
			&area.Description,
			&area.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		areas = append(areas, area)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return areas, nil
}

// ListPaginated lista áreas com paginação
func (r *AreaRepositoryPG) ListPaginated(page, pageSize int) ([]*domain.Area, *domain.Pagination, error) {
	// Validar parâmetros
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// Contar total de registros
	var totalItems int
	countQuery := `SELECT COUNT(*) FROM areas`
	err := r.DB.QueryRow(countQuery).Scan(&totalItems)
	if err != nil {
		return nil, nil, err
	}

	// Calcular total de páginas
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	// Buscar dados paginados
	query := `
		SELECT id, name, description, created_at
		FROM areas
		ORDER BY name
		LIMIT $1 OFFSET $2`

	rows, err := r.DB.Query(query, pageSize, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var areas []*domain.Area
	for rows.Next() {
		area := &domain.Area{}
		err := rows.Scan(
			&area.ID,
			&area.Name,
			&area.Description,
			&area.CreatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		areas = append(areas, area)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	pagination := &domain.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return areas, pagination, nil
}
