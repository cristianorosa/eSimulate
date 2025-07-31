package infra

import (
	"database/sql"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// DomainRepositoryPG implementa Repository para PostgreSQL
type DomainRepositoryPG struct {
	DB *sql.DB
}

// Create cria um novo domínio
func (r *DomainRepositoryPG) Create(domain *domain.Domain) error {
	query := `
		INSERT INTO domains (exam_id, name, description, weight_percentage, order_index, created_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		RETURNING id, created_at`

	err := r.DB.QueryRow(
		query,
		domain.ExamID,
		domain.Name,
		domain.Description,
		domain.WeightPercentage,
		domain.OrderIndex,
	).Scan(&domain.ID, &domain.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Update atualiza um domínio existente
func (r *DomainRepositoryPG) Update(d *domain.Domain) error {
	query := `
		UPDATE domains 
		SET name = $1, description = $2, weight_percentage = $3, order_index = $4
		WHERE id = $5`

	result, err := r.DB.Exec(
		query,
		d.Name,
		d.Description,
		d.WeightPercentage,
		d.OrderIndex,
		d.ID,
	)
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

// Delete remove um domínio
func (r *DomainRepositoryPG) Delete(id int) error {
	query := `DELETE FROM domains WHERE id = $1`

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

// FindByID busca um domínio por ID
func (r *DomainRepositoryPG) FindByID(id int) (*domain.Domain, error) {
	query := `
		SELECT id, exam_id, name, description, weight_percentage, order_index, created_at
		FROM domains
		WHERE id = $1`

	d := &domain.Domain{}
	err := r.DB.QueryRow(query, id).Scan(
		&d.ID,
		&d.ExamID,
		&d.Name,
		&d.Description,
		&d.WeightPercentage,
		&d.OrderIndex,
		&d.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return d, nil
}

// ListByExam lista todos os domínios de um exame
func (r *DomainRepositoryPG) ListByExam(examID int) ([]*domain.Domain, error) {
	query := `
		SELECT id, exam_id, name, description, weight_percentage, order_index, created_at
		FROM domains
		WHERE exam_id = $1
		ORDER BY order_index, name`

	rows, err := r.DB.Query(query, examID)
	if err != nil {
		return []*domain.Domain{}, err
	}
	defer rows.Close()

	var domains []*domain.Domain
	for rows.Next() {
		d := &domain.Domain{}
		err := rows.Scan(
			&d.ID,
			&d.ExamID,
			&d.Name,
			&d.Description,
			&d.WeightPercentage,
			&d.OrderIndex,
			&d.CreatedAt,
		)
		if err != nil {
			return []*domain.Domain{}, err
		}
		domains = append(domains, d)
	}

	if err = rows.Err(); err != nil {
		return []*domain.Domain{}, err
	}

	// Garantir que sempre retorne um array, mesmo que vazio
	if domains == nil {
		domains = []*domain.Domain{}
	}

	return domains, nil
}

// ListAll lista todos os domínios
func (r *DomainRepositoryPG) ListAll() ([]*domain.Domain, error) {
	query := `
		SELECT id, exam_id, name, description, weight_percentage, order_index, created_at
		FROM domains
		ORDER BY exam_id, order_index, name`

	rows, err := r.DB.Query(query)
	if err != nil {
		return []*domain.Domain{}, err
	}
	defer rows.Close()

	var domains []*domain.Domain
	for rows.Next() {
		d := &domain.Domain{}
		err := rows.Scan(
			&d.ID,
			&d.ExamID,
			&d.Name,
			&d.Description,
			&d.WeightPercentage,
			&d.OrderIndex,
			&d.CreatedAt,
		)
		if err != nil {
			return []*domain.Domain{}, err
		}
		domains = append(domains, d)
	}

	if err = rows.Err(); err != nil {
		return []*domain.Domain{}, err
	}

	// Garantir que sempre retorne um array, mesmo que vazio
	if domains == nil {
		domains = []*domain.Domain{}
	}

	return domains, nil
}
