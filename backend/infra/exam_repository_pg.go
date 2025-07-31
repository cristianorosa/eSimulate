package infra

import (
	"database/sql"
	"math"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// ExamRepositoryPG implementa ExamRepository para PostgreSQL
type ExamRepositoryPG struct {
	DB *sql.DB
}

// Create cria um novo exame
func (r *ExamRepositoryPG) Create(exam *domain.Exam) error {
	query := `
		INSERT INTO exams (title, description, area_id, max_time_minutes, passing_score, is_active, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
		RETURNING id, created_at`

	err := r.DB.QueryRow(query, exam.Title, exam.Description, exam.AreaID, exam.MaxTimeMinutes, exam.PassingScore, exam.IsActive, exam.CreatedBy).Scan(&exam.ID, &exam.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Update atualiza um exame existente
func (r *ExamRepositoryPG) Update(exam *domain.Exam) error {
	query := `
		UPDATE exams 
		SET title = $1, description = $2, area_id = $3, max_time_minutes = $4, passing_score = $5, is_active = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7`

	result, err := r.DB.Exec(query, exam.Title, exam.Description, exam.AreaID, exam.MaxTimeMinutes, exam.PassingScore, exam.IsActive, exam.ID)
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

// Delete remove um exame
func (r *ExamRepositoryPG) Delete(id int) error {
	query := `DELETE FROM exams WHERE id = $1`

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

// FindByID busca um exame por ID
func (r *ExamRepositoryPG) FindByID(id int) (*domain.Exam, error) {
	query := `
		SELECT id, title, description, area_id, max_time_minutes, passing_score, is_active, created_by, created_at
		FROM exams
		WHERE id = $1`

	exam := &domain.Exam{}
	err := r.DB.QueryRow(query, id).Scan(
		&exam.ID,
		&exam.Title,
		&exam.Description,
		&exam.AreaID,
		&exam.MaxTimeMinutes,
		&exam.PassingScore,
		&exam.IsActive,
		&exam.CreatedBy,
		&exam.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return exam, nil
}

// ListAll lista todos os exames
func (r *ExamRepositoryPG) ListAll() ([]*domain.Exam, error) {
	query := `
		SELECT id, title, description, area_id, max_time_minutes, passing_score, is_active, created_by, created_at
		FROM exams
		ORDER BY title`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exams []*domain.Exam
	for rows.Next() {
		exam := &domain.Exam{}
		err := rows.Scan(
			&exam.ID,
			&exam.Title,
			&exam.Description,
			&exam.AreaID,
			&exam.MaxTimeMinutes,
			&exam.PassingScore,
			&exam.IsActive,
			&exam.CreatedBy,
			&exam.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		exams = append(exams, exam)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return exams, nil
}

// ListByArea lista exames por área
func (r *ExamRepositoryPG) ListByArea(areaID int) ([]*domain.Exam, error) {
	query := `
		SELECT id, title, description, area_id, max_time_minutes, passing_score, is_active, created_by, created_at
		FROM exams
		WHERE area_id = $1
		ORDER BY title`

	rows, err := r.DB.Query(query, areaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exams []*domain.Exam
	for rows.Next() {
		exam := &domain.Exam{}
		err := rows.Scan(
			&exam.ID,
			&exam.Title,
			&exam.Description,
			&exam.AreaID,
			&exam.MaxTimeMinutes,
			&exam.PassingScore,
			&exam.IsActive,
			&exam.CreatedBy,
			&exam.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		exams = append(exams, exam)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return exams, nil
}

// ListPaginated lista exames com paginação
func (r *ExamRepositoryPG) ListPaginated(page, pageSize int, areaID *int) ([]*domain.Exam, *domain.Pagination, error) {
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

	// Construir query base
	baseQuery := `FROM exams`
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if areaID != nil {
		whereClause = `WHERE area_id = $` + string(rune(argIndex+'0'))
		args = append(args, *areaID)
		argIndex++
	}

	// Contar total de registros
	countQuery := `SELECT COUNT(*) ` + baseQuery + whereClause
	var totalItems int
	err := r.DB.QueryRow(countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, nil, err
	}

	// Calcular total de páginas
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	// Buscar dados paginados
	query := `
		SELECT id, title, description, area_id, max_time_minutes, passing_score, is_active, created_by, created_at
		` + baseQuery + whereClause + `
		ORDER BY title
		LIMIT $` + string(rune(argIndex+'0')) + ` OFFSET $` + string(rune(argIndex+1+'0'))

	args = append(args, pageSize, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var exams []*domain.Exam
	for rows.Next() {
		exam := &domain.Exam{}
		err := rows.Scan(
			&exam.ID,
			&exam.Title,
			&exam.Description,
			&exam.AreaID,
			&exam.MaxTimeMinutes,
			&exam.PassingScore,
			&exam.IsActive,
			&exam.CreatedBy,
			&exam.CreatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		exams = append(exams, exam)
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

	return exams, pagination, nil
}
