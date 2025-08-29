package infra

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// QuestionTagRepositoryPG implementa o repositório de tags de questões para PostgreSQL
type QuestionTagRepositoryPG struct {
	DB *sql.DB
}

// NewQuestionTagRepositoryPG cria uma nova instância do repositório
func NewQuestionTagRepositoryPG(db *sql.DB) domain.QuestionTagRepository {
	return &QuestionTagRepositoryPG{DB: db}
}

// CreateTag creates a new question tag
func (r *QuestionTagRepositoryPG) CreateTag(tag *domain.QuestionTag) error {
	query := `
		INSERT INTO question_tags (name, created_at)
		VALUES ($1, $2)
		RETURNING id`

	now := time.Now()
	err := r.DB.QueryRow(query, tag.Name, now).Scan(&tag.ID)
	if err != nil {
		return fmt.Errorf("failed to create question tag: %w", err)
	}

	tag.CreatedAt = now
	return nil
}

// UpdateTag updates an existing question tag
func (r *QuestionTagRepositoryPG) UpdateTag(tag *domain.QuestionTag) error {
	query := `UPDATE question_tags SET name = $1 WHERE id = $2`

	result, err := r.DB.Exec(query, tag.Name, tag.ID)
	if err != nil {
		return fmt.Errorf("failed to update question tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("question tag not found")
	}

	return nil
}

// DeleteTag deletes a question tag
func (r *QuestionTagRepositoryPG) DeleteTag(id int) error {
	query := `DELETE FROM question_tags WHERE id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete question tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("question tag not found")
	}

	return nil
}

// FindTagByID finds a question tag by ID
func (r *QuestionTagRepositoryPG) FindTagByID(id int) (*domain.QuestionTag, error) {
	query := `SELECT id, name, created_at FROM question_tags WHERE id = $1`

	tag := &domain.QuestionTag{}

	err := r.DB.QueryRow(query, id).Scan(
		&tag.ID,
		&tag.Name,
		&tag.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("question tag not found")
		}
		return nil, fmt.Errorf("failed to find question tag: %w", err)
	}

	return tag, nil
}

// FindTagByName finds a question tag by name
func (r *QuestionTagRepositoryPG) FindTagByName(name string) (*domain.QuestionTag, error) {
	query := `SELECT id, name, created_at FROM question_tags WHERE LOWER(name) = LOWER($1)`

	tag := &domain.QuestionTag{}

	err := r.DB.QueryRow(query, name).Scan(
		&tag.ID,
		&tag.Name,
		&tag.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("question tag not found")
		}
		return nil, fmt.Errorf("failed to find question tag: %w", err)
	}

	return tag, nil
}

// SearchTagsByName searches tags by name pattern (for autocomplete)
func (r *QuestionTagRepositoryPG) SearchTagsByName(pattern string, limit int) ([]*domain.QuestionTag, error) {
	if limit <= 0 || limit > 50 {
		limit = 20 // Default limit
	}

	query := `
		SELECT id, name, created_at 
		FROM question_tags 
		WHERE LOWER(name) LIKE LOWER($1) 
		ORDER BY name ASC 
		LIMIT $2`

	rows, err := r.DB.Query(query, "%"+pattern+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search question tags: %w", err)
	}
	defer rows.Close()

	var tags []*domain.QuestionTag

	for rows.Next() {
		tag := &domain.QuestionTag{}

		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan question tag: %w", err)
		}

		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating question tags: %w", err)
	}

	return tags, nil
}

// ListTags lists all question tags with pagination
func (r *QuestionTagRepositoryPG) ListTags(page, pageSize int) ([]*domain.QuestionTag, error) {
	offset := (page - 1) * pageSize
	query := `
		SELECT id, name, created_at 
		FROM question_tags 
		ORDER BY name ASC
		LIMIT $1 OFFSET $2`

	rows, err := r.DB.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query question tags: %w", err)
	}
	defer rows.Close()

	var tags []*domain.QuestionTag

	for rows.Next() {
		tag := &domain.QuestionTag{}

		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan question tag: %w", err)
		}

		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating question tags: %w", err)
	}

	return tags, nil
}

// ListTagsWithStats lists all tags with usage statistics
func (r *QuestionTagRepositoryPG) ListTagsWithStats() ([]*domain.QuestionTagWithStats, error) {
	query := `
		SELECT 
			qt.id, qt.name, qt.created_at,
			COALESCE(COUNT(qta.question_id), 0) as question_count
		FROM question_tags qt
		LEFT JOIN question_tag_associations qta ON qt.id = qta.tag_id
		GROUP BY qt.id, qt.name, qt.created_at
		ORDER BY question_count DESC, qt.name ASC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query question tags with stats: %w", err)
	}
	defer rows.Close()

	var tags []*domain.QuestionTagWithStats

	for rows.Next() {
		tag := &domain.QuestionTagWithStats{
			QuestionTag: &domain.QuestionTag{},
		}

		err := rows.Scan(
			&tag.QuestionTag.ID,
			&tag.QuestionTag.Name,
			&tag.QuestionTag.CreatedAt,
			&tag.QuestionCount,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan question tag with stats: %w", err)
		}

		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating question tags with stats: %w", err)
	}

	return tags, nil
}

// AssociateQuestionTag associates a tag with a question
func (r *QuestionTagRepositoryPG) AssociateQuestionTag(questionID, tagID int) error {
	query := `
		INSERT INTO question_tag_associations (question_id, tag_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (question_id, tag_id) DO NOTHING`

	_, err := r.DB.Exec(query, questionID, tagID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to associate question tag: %w", err)
	}

	return nil
}

// DisassociateQuestionTag removes association between tag and question
func (r *QuestionTagRepositoryPG) DisassociateQuestionTag(questionID, tagID int) error {
	query := `DELETE FROM question_tag_associations WHERE question_id = $1 AND tag_id = $2`

	result, err := r.DB.Exec(query, questionID, tagID)
	if err != nil {
		return fmt.Errorf("failed to disassociate question tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("association not found")
	}

	return nil
}

// ListQuestionTags lists all tags associated with a question
func (r *QuestionTagRepositoryPG) ListQuestionTags(questionID int) ([]*domain.QuestionTag, error) {
	query := `
		SELECT qt.id, qt.name, qt.created_at
		FROM question_tags qt
		INNER JOIN question_tag_associations qta ON qt.id = qta.tag_id
		WHERE qta.question_id = $1
		ORDER BY qt.name ASC`

	rows, err := r.DB.Query(query, questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query question tags: %w", err)
	}
	defer rows.Close()

	var tags []*domain.QuestionTag

	for rows.Next() {
		tag := &domain.QuestionTag{}

		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan question tag: %w", err)
		}

		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating question tags: %w", err)
	}

	return tags, nil
}

// UpdateQuestionTags updates all tags for a question (replaces existing associations)
func (r *QuestionTagRepositoryPG) UpdateQuestionTags(questionID int, tagIDs []int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Remove existing associations
	_, err = tx.Exec("DELETE FROM question_tag_associations WHERE question_id = $1", questionID)
	if err != nil {
		return fmt.Errorf("failed to remove existing tag associations: %w", err)
	}

	// Add new associations
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			_, err = tx.Exec(
				"INSERT INTO question_tag_associations (question_id, tag_id, created_at) VALUES ($1, $2, $3)",
				questionID, tagID, time.Now())
			if err != nil {
				return fmt.Errorf("failed to create tag association: %w", err)
			}
		}
	}

	return tx.Commit()
}

// BulkAssociateQuestionTags associates multiple tags with a question
func (r *QuestionTagRepositoryPG) BulkAssociateQuestionTags(questionID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO question_tag_associations (question_id, tag_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (question_id, tag_id) DO NOTHING`

	for _, tagID := range tagIDs {
		_, err = tx.Exec(query, questionID, tagID, time.Now())
		if err != nil {
			return fmt.Errorf("failed to associate tag %d: %w", tagID, err)
		}
	}

	return tx.Commit()
}

// BulkDisassociateQuestionTags removes multiple tag associations from a question
func (r *QuestionTagRepositoryPG) BulkDisassociateQuestionTags(questionID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM question_tag_associations WHERE question_id = $1 AND tag_id = $2`

	for _, tagID := range tagIDs {
		_, err = tx.Exec(query, questionID, tagID)
		if err != nil {
			return fmt.Errorf("failed to disassociate tag %d: %w", tagID, err)
		}
	}

	return tx.Commit()
}

// GetTagStats returns statistics for a specific tag
func (r *QuestionTagRepositoryPG) GetTagStats(tagID int) (*domain.QuestionTagWithStats, error) {
	query := `
		SELECT 
			qt.id, qt.name, qt.created_at,
			COALESCE(COUNT(qta.question_id), 0) as question_count
		FROM question_tags qt
		LEFT JOIN question_tag_associations qta ON qt.id = qta.tag_id
		WHERE qt.id = $1
		GROUP BY qt.id, qt.name, qt.created_at`

	tag := &domain.QuestionTagWithStats{
		QuestionTag: &domain.QuestionTag{},
	}

	err := r.DB.QueryRow(query, tagID).Scan(
		&tag.QuestionTag.ID,
		&tag.QuestionTag.Name,
		&tag.QuestionTag.CreatedAt,
		&tag.QuestionCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found")
		}
		return nil, fmt.Errorf("failed to get tag stats: %w", err)
	}

	return tag, nil
}

// ListQuestionsByTag lists all questions that have a specific tag
func (r *QuestionTagRepositoryPG) ListQuestionsByTag(tagID int) ([]int, error) {
	query := `
		SELECT question_id
		FROM question_tag_associations
		WHERE tag_id = $1
		ORDER BY question_id ASC`

	rows, err := r.DB.Query(query, tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions by tag: %w", err)
	}
	defer rows.Close()

	var questionIDs []int

	for rows.Next() {
		var questionID int
		err := rows.Scan(&questionID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question ID: %w", err)
		}
		questionIDs = append(questionIDs, questionID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating questions by tag: %w", err)
	}

	return questionIDs, nil
}

// GetTagsByQuestions returns all tags for multiple questions
func (r *QuestionTagRepositoryPG) GetTagsByQuestions(questionIDs []int) (map[int][]*domain.QuestionTag, error) {
	if len(questionIDs) == 0 {
		return make(map[int][]*domain.QuestionTag), nil
	}

	// Build placeholders for IN clause
	placeholders := ""
	args := make([]interface{}, len(questionIDs))
	for i, id := range questionIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT 
			qta.question_id,
			qt.id, qt.name, qt.created_at
		FROM question_tag_associations qta
		INNER JOIN question_tags qt ON qta.tag_id = qt.id
		WHERE qta.question_id IN (%s)
		ORDER BY qta.question_id, qt.name`, placeholders)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags by questions: %w", err)
	}
	defer rows.Close()

	result := make(map[int][]*domain.QuestionTag)

	for rows.Next() {
		var questionID int
		tag := &domain.QuestionTag{}

		err := rows.Scan(
			&questionID,
			&tag.ID,
			&tag.Name,
			&tag.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}

		result[questionID] = append(result[questionID], tag)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tags by questions: %w", err)
	}

	return result, nil
}

// SearchTags searches for tags by name
func (r *QuestionTagRepositoryPG) SearchTags(searchTerm string, limit int) ([]*domain.QuestionTag, error) {
	query := `
		SELECT id, name, created_at 
		FROM question_tags 
		WHERE name ILIKE $1
		ORDER BY name ASC
		LIMIT $2`

	searchPattern := "%" + searchTerm + "%"
	rows, err := r.DB.Query(query, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search question tags: %w", err)
	}
	defer rows.Close()

	var tags []*domain.QuestionTag

	for rows.Next() {
		tag := &domain.QuestionTag{}

		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan question tag: %w", err)
		}

		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	return tags, nil
}

// CleanupUnusedTags removes tags that are not associated with any questions
func (r *QuestionTagRepositoryPG) CleanupUnusedTags() (int, error) {
	query := `
		DELETE FROM question_tags 
		WHERE id NOT IN (
			SELECT DISTINCT tag_id 
			FROM question_tag_associations
		)`

	result, err := r.DB.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup unused tags: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

// CountTags returns the total number of tags
func (r *QuestionTagRepositoryPG) CountTags() (int, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM question_tags").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count tags: %w", err)
	}
	return count, nil
}

// CountTagsByQuestion returns the number of tags for a specific question
func (r *QuestionTagRepositoryPG) CountTagsByQuestion(questionID int) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM question_tag_associations 
		WHERE question_id = $1`

	err := r.DB.QueryRow(query, questionID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count tags by question: %w", err)
	}

	return count, nil
}
