package infra

import (
	"database/sql"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// TopicRepositoryPG implementa TopicRepository para PostgreSQL
type TopicRepositoryPG struct {
	DB *sql.DB
}

// Create cria um novo tópico
func (r *TopicRepositoryPG) Create(topic *domain.Topic) error {
	query := `
		INSERT INTO topics (name, description, created_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		RETURNING id, created_at`

	var description interface{}
	if topic.Description != nil {
		description = *topic.Description
	} else {
		description = nil
	}

	err := r.DB.QueryRow(
		query,
		topic.Name,
		description,
	).Scan(&topic.ID, &topic.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Update atualiza um tópico existente
func (r *TopicRepositoryPG) Update(t *domain.Topic) error {
	query := `
		UPDATE topics 
		SET name = $1, description = $2
		WHERE id = $3`

	var description interface{}
	if t.Description != nil {
		description = *t.Description
	} else {
		description = nil
	}

	result, err := r.DB.Exec(
		query,
		t.Name,
		description,
		t.ID,
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

// Delete remove um tópico
func (r *TopicRepositoryPG) Delete(id int) error {
	query := `DELETE FROM topics WHERE id = $1`

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

// FindByID busca um tópico por ID
func (r *TopicRepositoryPG) FindByID(id int) (*domain.Topic, error) {
	query := `
		SELECT id, name, description, created_at
		FROM topics
		WHERE id = $1`

	t := &domain.Topic{}
	var description sql.NullString
	err := r.DB.QueryRow(query, id).Scan(
		&t.ID,
		&t.Name,
		&description,
		&t.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if description.Valid {
		t.Description = &description.String
	}

	return t, nil
}

// ListByExam lista todos os tópicos de um exame
func (r *TopicRepositoryPG) ListByExam(examID int) ([]*domain.Topic, error) {
	query := `
		SELECT t.id, t.name, t.description, t.created_at
		FROM topics t
		JOIN exam_topics et ON t.id = et.topic_id
		WHERE et.exam_id = $1
		ORDER BY et.order_index, t.name`

	rows, err := r.DB.Query(query, examID)
	if err != nil {
		return []*domain.Topic{}, err
	}
	defer rows.Close()

	var topics []*domain.Topic
	for rows.Next() {
		t := &domain.Topic{}
		var description sql.NullString
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&description,
			&t.CreatedAt,
		)
		if err != nil {
			return []*domain.Topic{}, err
		}

		if description.Valid {
			t.Description = &description.String
		}

		topics = append(topics, t)
	}

	if err = rows.Err(); err != nil {
		return []*domain.Topic{}, err
	}

	// Garantir que sempre retorne um array, mesmo que vazio
	if topics == nil {
		topics = []*domain.Topic{}
	}

	return topics, nil
}

// ListAll lista todos os tópicos
func (r *TopicRepositoryPG) ListAll() ([]*domain.Topic, error) {
	query := `
		SELECT id, name, description, created_at
		FROM topics
		ORDER BY created_at DESC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return []*domain.Topic{}, err
	}
	defer rows.Close()

	var topics []*domain.Topic
	for rows.Next() {
		t := &domain.Topic{}
		var description sql.NullString
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&description,
			&t.CreatedAt,
		)
		if err != nil {
			return []*domain.Topic{}, err
		}

		if description.Valid {
			t.Description = &description.String
		}

		topics = append(topics, t)
	}

	if err = rows.Err(); err != nil {
		return []*domain.Topic{}, err
	}

	if topics == nil {
		topics = []*domain.Topic{}
	}

	return topics, nil
}

// ListPaginated lista tópicos com paginação
func (r *TopicRepositoryPG) ListPaginated(page, pageSize int, examID *int) ([]*domain.Topic, *domain.Pagination, error) {
	// Calcular offset
	offset := (page - 1) * pageSize

	// Construir query base - nova estrutura sem exam_id
	baseQuery := `
		SELECT t.id, t.name, t.description, t.created_at
		FROM topics t`

	countQuery := `SELECT COUNT(*) FROM topics t`

	// Adicionar filtro por exame se fornecido (usando exam_topics)
	whereClause := ""
	joinClause := ""
	if examID != nil {
		joinClause = " JOIN exam_topics et ON t.id = et.topic_id"
		whereClause = " WHERE et.exam_id = $1"
	}

	// Query para contar total de registros
	var totalItems int
	var countErr error
	if examID != nil {
		countErr = r.DB.QueryRow(countQuery+joinClause+whereClause, *examID).Scan(&totalItems)
	} else {
		countErr = r.DB.QueryRow(countQuery + whereClause).Scan(&totalItems)
	}

	if countErr != nil {
		return nil, nil, countErr
	}

	// Query para buscar dados paginados
	var dataQuery string
	var rows *sql.Rows
	var queryErr error

	if examID != nil {
		dataQuery = baseQuery + joinClause + whereClause + ` ORDER BY t.created_at DESC LIMIT $2 OFFSET $3`
		rows, queryErr = r.DB.Query(dataQuery, *examID, pageSize, offset)
	} else {
		dataQuery = baseQuery + whereClause + ` ORDER BY t.created_at DESC LIMIT $1 OFFSET $2`
		rows, queryErr = r.DB.Query(dataQuery, pageSize, offset)
	}

	if queryErr != nil {
		return nil, nil, queryErr
	}
	defer rows.Close()

	var topics []*domain.Topic
	for rows.Next() {
		t := &domain.Topic{}
		var description sql.NullString
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&description,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, nil, err
		}

		if description.Valid {
			t.Description = &description.String
		}

		topics = append(topics, t)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	// Calcular informações de paginação
	totalPages := (totalItems + pageSize - 1) / pageSize

	pagination := &domain.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	if topics == nil {
		topics = []*domain.Topic{}
	}

	return topics, pagination, nil
}

// ListByExamWithDetails returns topics for an exam with detailed information
func (r *TopicRepositoryPG) ListByExamWithDetails(examID int) ([]*domain.ExamTopicWithDetails, error) {
	// This method would typically be implemented using the exam_topics table
	// For now, we'll return an empty slice to satisfy the interface
	return []*domain.ExamTopicWithDetails{}, nil
}
