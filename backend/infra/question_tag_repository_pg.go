package infra

import (
	"database/sql"
	"fmt"
	"strings"
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
		INSERT INTO question_tags (name, description, created_at)
		VALUES ($1, $2, $3)
		RETURNING id`
	
	now := time.Now()
	err := r.DB.QueryRow(query, tag.Name, tag.Description, now).Scan(&tag.ID)
	if err != nil {
		return fmt.Errorf("failed to create question tag: %w", err)
	}
	
	tag.CreatedAt = now
	return nil
}

// UpdateTag updates an existing question tag
func (r *QuestionTagRepositoryPG) UpdateTag(tag *domain.QuestionTag) error {
	query := `UPDATE question_tags SET name = $1, description = $2 WHERE id = $3`
	
	result, err := r.DB.Exec(query, tag.Name, tag.Description, tag.ID)
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
	query := `SELECT id, name, description, created_at FROM question_tags WHERE id = $1`
	
	tag := &domain.QuestionTag{}
	var description sql.NullString
	
	err := r.DB.QueryRow(query, id).Scan(
		&tag.ID,
		&tag.Name,
		&description,
		&tag.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("question tag not found")
		}
		return nil, fmt.Errorf("failed to find question tag: %w", err)
	}
	
	if description.Valid {
		tag.Description = &description.String
	}
	
	return tag, nil
}

// FindTagByName finds a question tag by name
func (r *QuestionTagRepositoryPG) FindTagByName(name string) (*domain.QuestionTag, error) {
	query := `SELECT id, name, description, created_at FROM question_tags WHERE name = $1`
	
	tag := &domain.QuestionTag{}
	var description sql.NullString
	
	err := r.DB.QueryRow(query, name).Scan(
		&tag.ID,
		&tag.Name,
		&description,
		&tag.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("question tag not found")
		}
		return nil, fmt.Errorf("failed to find question tag: %w", err)
	}
	
	if description.Valid {
		tag.Description = &description.String
	}
	
	return tag, nil
}

// ListAllTags returns all question tags
func (r *QuestionTagRepositoryPG) ListAllTags() ([]*domain.QuestionTag, error) {
	query := `
		SELECT id, name, description, created_at 
		FROM question_tags 
		ORDER BY name ASC`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list question tags: %w", err)
	}
	defer rows.Close()
	
	var tags []*domain.QuestionTag
	for rows.Next() {
		tag := &domain.QuestionTag{}
		var description sql.NullString
		
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&description,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question tag: %w", err)
		}
		
		if description.Valid {
			tag.Description = &description.String
		}
		
		tags = append(tags, tag)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over question tag rows: %w", err)
	}
	
	return tags, nil
}

// ListTagsWithStats returns all tags with usage statistics
func (r *QuestionTagRepositoryPG) ListTagsWithStats() ([]*domain.QuestionTagWithStats, error) {
	query := `
		SELECT 
			qt.id, qt.name, qt.description, qt.created_at,
			COUNT(qta.question_id) as question_count
		FROM question_tags qt
		LEFT JOIN question_tag_associations qta ON qt.id = qta.tag_id
		GROUP BY qt.id, qt.name, qt.description, qt.created_at
		ORDER BY question_count DESC, qt.name ASC`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list question tags with stats: %w", err)
	}
	defer rows.Close()
	
	var tags []*domain.QuestionTagWithStats
	for rows.Next() {
		tag := &domain.QuestionTagWithStats{
			QuestionTag: &domain.QuestionTag{},
		}
		var description sql.NullString
		
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&description,
			&tag.CreatedAt,
			&tag.QuestionCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question tag with stats: %w", err)
		}
		
		if description.Valid {
			tag.Description = &description.String
		}
		
		tag.UsageCount = tag.QuestionCount // For now, same as question count
		tags = append(tags, tag)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over question tag stats rows: %w", err)
	}
	
	return tags, nil
}

// AssociateQuestionTag creates an association between a question and a tag
func (r *QuestionTagRepositoryPG) AssociateQuestionTag(questionID, tagID int) error {
	query := `
		INSERT INTO question_tag_associations (question_id, tag_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (question_id, tag_id) DO NOTHING`
	
	_, err := r.DB.Exec(query, questionID, tagID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to associate question with tag: %w", err)
	}
	
	return nil
}

// DisassociateQuestionTag removes an association between a question and a tag
func (r *QuestionTagRepositoryPG) DisassociateQuestionTag(questionID, tagID int) error {
	query := `DELETE FROM question_tag_associations WHERE question_id = $1 AND tag_id = $2`
	
	result, err := r.DB.Exec(query, questionID, tagID)
	if err != nil {
		return fmt.Errorf("failed to disassociate question from tag: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("question-tag association not found")
	}
	
	return nil
}

// ListTagsByQuestion returns all tags associated with a question
func (r *QuestionTagRepositoryPG) ListTagsByQuestion(questionID int) ([]*domain.QuestionTag, error) {
	query := `
		SELECT qt.id, qt.name, qt.description, qt.created_at
		FROM question_tags qt
		JOIN question_tag_associations qta ON qt.id = qta.tag_id
		WHERE qta.question_id = $1
		ORDER BY qt.name ASC`
	
	rows, err := r.DB.Query(query, questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags for question: %w", err)
	}
	defer rows.Close()
	
	var tags []*domain.QuestionTag
	for rows.Next() {
		tag := &domain.QuestionTag{}
		var description sql.NullString
		
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&description,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question tag: %w", err)
		}
		
		if description.Valid {
			tag.Description = &description.String
		}
		
		tags = append(tags, tag)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over question tag rows: %w", err)
	}
	
	return tags, nil
}

// ListQuestionsByTag returns all questions associated with a tag
func (r *QuestionTagRepositoryPG) ListQuestionsByTag(tagID int) ([]*domain.Question, error) {
	query := `
		SELECT q.id, q.topic_id, q.statement, q.problem, q.content_type, 
			   q.explanation, q.question_type, q.difficulty_level, q.created_by,
			   q.is_active, q.created_at, q.updated_at
		FROM questions q
		JOIN question_tag_associations qta ON q.id = qta.question_id
		WHERE qta.tag_id = $1
		ORDER BY q.created_at DESC`
	
	rows, err := r.DB.Query(query, tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to list questions for tag: %w", err)
	}
	defer rows.Close()
	
	var questions []*domain.Question
	for rows.Next() {
		question := &domain.Question{}
		var explanation sql.NullString
		var updatedAt sql.NullTime
		
		err := rows.Scan(
			&question.ID,
			&question.TopicID,
			&question.Statement,
			&question.Problem,
			&question.ContentType,
			&explanation,
			&question.QuestionType,
			&question.DifficultyLevel,
			&question.CreatedBy,
			&question.IsActive,
			&question.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		
		if explanation.Valid {
			question.Explanation = explanation.String
		}
		if updatedAt.Valid {
			question.UpdatedAt = updatedAt.Time
		}
		
		questions = append(questions, question)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over question rows: %w", err)
	}
	
	return questions, nil
}

// ListQuestionsByTags returns all questions associated with any of the given tags
func (r *QuestionTagRepositoryPG) ListQuestionsByTags(tagIDs []int) ([]*domain.Question, error) {
	if len(tagIDs) == 0 {
		return []*domain.Question{}, nil
	}
	
	placeholders := make([]string, len(tagIDs))
	args := make([]interface{}, len(tagIDs))
	
	for i, tagID := range tagIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = tagID
	}
	
	query := fmt.Sprintf(`
		SELECT DISTINCT q.id, q.topic_id, q.statement, q.problem, q.content_type, 
			   q.explanation, q.question_type, q.difficulty_level, q.created_by,
			   q.is_active, q.created_at, q.updated_at
		FROM questions q
		JOIN question_tag_associations qta ON q.id = qta.question_id
		WHERE qta.tag_id IN (%s)
		ORDER BY q.created_at DESC`,
		strings.Join(placeholders, ","))
	
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list questions for tags: %w", err)
	}
	defer rows.Close()
	
	var questions []*domain.Question
	for rows.Next() {
		question := &domain.Question{}
		var explanation sql.NullString
		var updatedAt sql.NullTime
		
		err := rows.Scan(
			&question.ID,
			&question.TopicID,
			&question.Statement,
			&question.Problem,
			&question.ContentType,
			&explanation,
			&question.QuestionType,
			&question.DifficultyLevel,
			&question.CreatedBy,
			&question.IsActive,
			&question.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		
		if explanation.Valid {
			question.Explanation = explanation.String
		}
		if updatedAt.Valid {
			question.UpdatedAt = updatedAt.Time
		}
		
		questions = append(questions, question)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over question rows: %w", err)
	}
	
	return questions, nil
}

// BulkAssociateQuestionTags associates a question with multiple tags
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
	
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()
	
	now := time.Now()
	for _, tagID := range tagIDs {
		_, err := stmt.Exec(questionID, tagID, now)
		if err != nil {
			return fmt.Errorf("failed to associate question with tag %d: %w", tagID, err)
		}
	}
	
	return tx.Commit()
}

// BulkDisassociateQuestionTags removes associations between a question and multiple tags
func (r *QuestionTagRepositoryPG) BulkDisassociateQuestionTags(questionID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}
	
	placeholders := make([]string, len(tagIDs))
	args := make([]interface{}, len(tagIDs)+1)
	args[0] = questionID
	
	for i, tagID := range tagIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = tagID
	}
	
	query := fmt.Sprintf(
		`DELETE FROM question_tag_associations WHERE question_id = $1 AND tag_id IN (%s)`,
		strings.Join(placeholders, ","),
	)
	
	_, err := r.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to bulk disassociate question from tags: %w", err)
	}
	
	return nil
}

// UpdateQuestionTags replaces all tag associations for a question
func (r *QuestionTagRepositoryPG) UpdateQuestionTags(questionID int, tagIDs []int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Remove all existing associations
	deleteQuery := `DELETE FROM question_tag_associations WHERE question_id = $1`
	_, err = tx.Exec(deleteQuery, questionID)
	if err != nil {
		return fmt.Errorf("failed to remove existing tag associations: %w", err)
	}
	
	// Add new associations
	if len(tagIDs) > 0 {
		insertQuery := `
			INSERT INTO question_tag_associations (question_id, tag_id, created_at)
			VALUES ($1, $2, $3)`
		
		stmt, err := tx.Prepare(insertQuery)
		if err != nil {
			return fmt.Errorf("failed to prepare insert statement: %w", err)
		}
		defer stmt.Close()
		
		now := time.Now()
		for _, tagID := range tagIDs {
			_, err := stmt.Exec(questionID, tagID, now)
			if err != nil {
				return fmt.Errorf("failed to associate question with tag %d: %w", tagID, err)
			}
		}
	}
	
	return tx.Commit()
}

// GetTagUsageStats returns usage statistics for a tag
func (r *QuestionTagRepositoryPG) GetTagUsageStats(tagID int) (*domain.TagUsageStats, error) {
	// This is a simplified implementation
	// A full implementation would include more detailed statistics
	
	tag, err := r.FindTagByID(tagID)
	if err != nil {
		return nil, err
	}
	
	// Count total questions
	countQuery := `
		SELECT COUNT(*) 
		FROM question_tag_associations 
		WHERE tag_id = $1`
	
	var totalQuestions int
	err = r.DB.QueryRow(countQuery, tagID).Scan(&totalQuestions)
	if err != nil {
		return nil, fmt.Errorf("failed to get question count for tag: %w", err)
	}
	
	stats := &domain.TagUsageStats{
		TagID:               tagID,
		TagName:             tag.Name,
		TotalQuestions:      totalQuestions,
		QuestionsByTopic:    make(map[int]int),
		QuestionsByDifficulty: make(map[int]int),
		QuestionsByExam:     make(map[int]int),
		MostUsedWith:        []*domain.TagCooccurrence{},
	}
	
	return stats, nil
}

// GetQuestionTagStats returns tag statistics for a question
func (r *QuestionTagRepositoryPG) GetQuestionTagStats(questionID int) (*domain.QuestionTagStats, error) {
	// Get tags for the question
	tags, err := r.ListTagsByQuestion(questionID)
	if err != nil {
		return nil, err
	}
	
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}
	
	stats := &domain.QuestionTagStats{
		QuestionID:       questionID,
		TotalTags:        len(tags),
		TagNames:         tagNames,
		TagCategories:    []string{}, // Could be implemented based on tag naming conventions
		SimilarQuestions: []int{},    // Could be implemented based on tag overlap
	}
	
	return stats, nil
}

// GetMostUsedTags returns the most frequently used tags
func (r *QuestionTagRepositoryPG) GetMostUsedTags(limit int) ([]*domain.QuestionTagWithStats, error) {
	query := `
		SELECT 
			qt.id, qt.name, qt.description, qt.created_at,
			COUNT(qta.question_id) as question_count
		FROM question_tags qt
		LEFT JOIN question_tag_associations qta ON qt.id = qta.tag_id
		GROUP BY qt.id, qt.name, qt.description, qt.created_at
		ORDER BY question_count DESC
		LIMIT $1`
	
	rows, err := r.DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get most used tags: %w", err)
	}
	defer rows.Close()
	
	var tags []*domain.QuestionTagWithStats
	for rows.Next() {
		tag := &domain.QuestionTagWithStats{
			QuestionTag: &domain.QuestionTag{},
		}
		var description sql.NullString
		
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&description,
			&tag.CreatedAt,
			&tag.QuestionCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question tag with stats: %w", err)
		}
		
		if description.Valid {
			tag.Description = &description.String
		}
		
		tag.UsageCount = tag.QuestionCount
		tags = append(tags, tag)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over most used tags rows: %w", err)
	}
	
	return tags, nil
}

// GetLeastUsedTags returns the least frequently used tags
func (r *QuestionTagRepositoryPG) GetLeastUsedTags(limit int) ([]*domain.QuestionTagWithStats, error) {
	query := `
		SELECT 
			qt.id, qt.name, qt.description, qt.created_at,
			COUNT(qta.question_id) as question_count
		FROM question_tags qt
		LEFT JOIN question_tag_associations qta ON qt.id = qta.tag_id
		GROUP BY qt.id, qt.name, qt.description, qt.created_at
		ORDER BY question_count ASC, qt.name ASC
		LIMIT $1`
	
	rows, err := r.DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get least used tags: %w", err)
	}
	defer rows.Close()
	
	var tags []*domain.QuestionTagWithStats
	for rows.Next() {
		tag := &domain.QuestionTagWithStats{
			QuestionTag: &domain.QuestionTag{},
		}
		var description sql.NullString
		
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&description,
			&tag.CreatedAt,
			&tag.QuestionCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question tag with stats: %w", err)
		}
		
		if description.Valid {
			tag.Description = &description.String
		}
		
		tag.UsageCount = tag.QuestionCount
		tags = append(tags, tag)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over least used tags rows: %w", err)
	}
	
	return tags, nil
}

// SearchTags searches for tags by name
func (r *QuestionTagRepositoryPG) SearchTags(query string, limit int) ([]*domain.QuestionTag, error) {
	searchQuery := `
		SELECT id, name, description, created_at 
		FROM question_tags 
		WHERE name ILIKE $1 OR description ILIKE $1
		ORDER BY name ASC
		LIMIT $2`
	
	searchPattern := "%" + query + "%"
	rows, err := r.DB.Query(searchQuery, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search question tags: %w", err)
	}
	defer rows.Close()
	
	var tags []*domain.QuestionTag
	for rows.Next() {
		tag := &domain.QuestionTag{}
		var description sql.NullString
		
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&description,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question tag: %w", err)
		}
		
		if description.Valid {
			tag.Description = &description.String
		}
		
		tags = append(tags, tag)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over search results: %w", err)
	}
	
	return tags, nil
}

// FilterQuestionsByTags filters questions based on tag criteria
func (r *QuestionTagRepositoryPG) FilterQuestionsByTags(filters *domain.QuestionTagFilters) ([]*domain.Question, error) {
	if len(filters.TagIDs) == 0 && len(filters.TagNames) == 0 {
		return []*domain.Question{}, nil
	}
	
	var whereConditions []string
	var args []interface{}
	argCount := 0
	
	// Build tag conditions
	if len(filters.TagIDs) > 0 {
		placeholders := make([]string, len(filters.TagIDs))
		for i, tagID := range filters.TagIDs {
			argCount++
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, tagID)
		}
		whereConditions = append(whereConditions, fmt.Sprintf("qt.id IN (%s)", strings.Join(placeholders, ",")))
	}
	
	if len(filters.TagNames) > 0 {
		placeholders := make([]string, len(filters.TagNames))
		for i, tagName := range filters.TagNames {
			argCount++
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, tagName)
		}
		whereConditions = append(whereConditions, fmt.Sprintf("qt.name IN (%s)", strings.Join(placeholders, ",")))
	}
	
	// Additional filters
	if filters.TopicID != nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("q.topic_id = $%d", argCount))
		args = append(args, *filters.TopicID)
	}
	
	if filters.DifficultyLevel != nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("q.difficulty_level = $%d", argCount))
		args = append(args, *filters.DifficultyLevel)
	}
	
	// Build the query
	baseQuery := `
		SELECT DISTINCT q.id, q.topic_id, q.statement, q.problem, q.content_type, 
			   q.explanation, q.question_type, q.difficulty_level, q.created_by,
			   q.is_active, q.created_at, q.updated_at
		FROM questions q
		JOIN question_tag_associations qta ON q.id = qta.question_id
		JOIN question_tags qt ON qta.tag_id = qt.id`
	
	if len(whereConditions) > 0 {
		baseQuery += " WHERE " + strings.Join(whereConditions, " AND ")
	}
	
	baseQuery += " ORDER BY q.created_at DESC"
	
	if filters.Limit > 0 {
		argCount++
		baseQuery += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.Limit)
		
		if filters.Offset > 0 {
			argCount++
			baseQuery += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, filters.Offset)
		}
	}
	
	rows, err := r.DB.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to filter questions by tags: %w", err)
	}
	defer rows.Close()
	
	var questions []*domain.Question
	for rows.Next() {
		question := &domain.Question{}
		var explanation sql.NullString
		var updatedAt sql.NullTime
		
		err := rows.Scan(
			&question.ID,
			&question.TopicID,
			&question.Statement,
			&question.Problem,
			&question.ContentType,
			&explanation,
			&question.QuestionType,
			&question.DifficultyLevel,
			&question.CreatedBy,
			&question.IsActive,
			&question.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		
		if explanation.Valid {
			question.Explanation = explanation.String
		}
		if updatedAt.Valid {
			question.UpdatedAt = updatedAt.Time
		}
		
		questions = append(questions, question)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over filtered questions: %w", err)
	}
	
	return questions, nil
}
