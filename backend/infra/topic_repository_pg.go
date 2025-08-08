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
		INSERT INTO topics (exam_id, name, weight_percentage, order_index, questions_count, created_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		RETURNING id, created_at`

	err := r.DB.QueryRow(
		query,
		topic.ExamID,
		topic.Name,
		topic.WeightPercentage,
		topic.OrderIndex,
		topic.QuestionsCount,
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
		SET name = $1, weight_percentage = $2, order_index = $3, questions_count = $4
		WHERE id = $5`

	result, err := r.DB.Exec(
		query,
		t.Name,
		t.WeightPercentage,
		t.OrderIndex,
		t.QuestionsCount,
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
		SELECT id, exam_id, name, weight_percentage, order_index, questions_count, created_at
		FROM topics
		WHERE id = $1`

	t := &domain.Topic{}
	err := r.DB.QueryRow(query, id).Scan(
		&t.ID,
		&t.ExamID,
		&t.Name,
		&t.WeightPercentage,
		&t.OrderIndex,
		&t.QuestionsCount,
		&t.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return t, nil
}

// ListByExam lista todos os tópicos de um exame
func (r *TopicRepositoryPG) ListByExam(examID int) ([]*domain.Topic, error) {
	query := `
		SELECT id, exam_id, name, weight_percentage, order_index, questions_count, created_at
		FROM topics
		WHERE exam_id = $1
		ORDER BY order_index, name`

	rows, err := r.DB.Query(query, examID)
	if err != nil {
		return []*domain.Topic{}, err
	}
	defer rows.Close()

	var topics []*domain.Topic
	for rows.Next() {
		t := &domain.Topic{}
		err := rows.Scan(
			&t.ID,
			&t.ExamID,
			&t.Name,
			&t.WeightPercentage,
			&t.OrderIndex,
			&t.QuestionsCount,
			&t.CreatedAt,
		)
		if err != nil {
			return []*domain.Topic{}, err
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
		SELECT id, exam_id, name, weight_percentage, order_index, questions_count, created_at
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
		err := rows.Scan(
			&t.ID,
			&t.ExamID,
			&t.Name,
			&t.WeightPercentage,
			&t.OrderIndex,
			&t.QuestionsCount,
			&t.CreatedAt,
		)
		if err != nil {
			return []*domain.Topic{}, err
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

	// Construir query base
	baseQuery := `
		SELECT id, exam_id, name, weight_percentage, order_index, questions_count, created_at
		FROM topics`

	countQuery := `SELECT COUNT(*) FROM topics`

	// Adicionar filtro por exame se fornecido
	whereClause := ""
	if examID != nil {
		whereClause = " WHERE exam_id = $1"
	}

	// Query para contar total de registros
	var totalItems int
	var countErr error
	if examID != nil {
		countErr = r.DB.QueryRow(countQuery+whereClause, *examID).Scan(&totalItems)
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
		dataQuery = baseQuery + whereClause + ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		rows, queryErr = r.DB.Query(dataQuery, *examID, pageSize, offset)
	} else {
		dataQuery = baseQuery + whereClause + ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		rows, queryErr = r.DB.Query(dataQuery, pageSize, offset)
	}

	if queryErr != nil {
		return nil, nil, queryErr
	}
	defer rows.Close()

	var topics []*domain.Topic
	for rows.Next() {
		t := &domain.Topic{}
		err := rows.Scan(
			&t.ID,
			&t.ExamID,
			&t.Name,
			&t.WeightPercentage,
			&t.OrderIndex,
			&t.QuestionsCount,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, nil, err
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
