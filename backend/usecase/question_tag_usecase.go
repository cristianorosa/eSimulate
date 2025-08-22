package usecase

import (
	"fmt"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// QuestionTagUsecase implementa a lógica de negócio para tags de questões
type QuestionTagUsecase struct {
	QuestionTagRepo domain.QuestionTagRepository
	QuestionRepo    domain.QuestionRepository
}

// NewQuestionTagUsecase cria uma nova instância do use case
func NewQuestionTagUsecase(
	questionTagRepo domain.QuestionTagRepository,
	questionRepo domain.QuestionRepository,
) *QuestionTagUsecase {
	return &QuestionTagUsecase{
		QuestionTagRepo: questionTagRepo,
		QuestionRepo:    questionRepo,
	}
}

// CreateTag cria uma nova tag
func (uc *QuestionTagUsecase) CreateTag(name string) (*domain.QuestionTag, error) {
	// Validar nome da tag
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return nil, fmt.Errorf("tag name must be at least 2 characters long")
	}

	// Verificar se a tag já existe
	existing, err := uc.QuestionTagRepo.FindTagByName(name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("tag with name '%s' already exists", name)
	}

	// Criar a tag
	tag := &domain.QuestionTag{
		Name: name,
	}

	err = uc.QuestionTagRepo.CreateTag(tag)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}

// UpdateTag atualiza uma tag existente
func (uc *QuestionTagUsecase) UpdateTag(id int, name string) (*domain.QuestionTag, error) {
	// Buscar a tag existente
	tag, err := uc.QuestionTagRepo.FindTagByID(id)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	// Validar nome da tag
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return nil, fmt.Errorf("tag name must be at least 2 characters long")
	}

	// Verificar se outro tag já usa esse nome
	existing, err := uc.QuestionTagRepo.FindTagByName(name)
	if err == nil && existing != nil && existing.ID != id {
		return nil, fmt.Errorf("tag with name '%s' already exists", name)
	}

	// Atualizar o campo
	tag.Name = name

	err = uc.QuestionTagRepo.UpdateTag(tag)
	if err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return tag, nil
}

// DeleteTag remove uma tag
func (uc *QuestionTagUsecase) DeleteTag(id int) error {
	// Verificar se a tag existe
	_, err := uc.QuestionTagRepo.FindTagByID(id)
	if err != nil {
		return fmt.Errorf("tag not found: %w", err)
	}

	// Verificar se a tag está sendo usada
	questions, err := uc.QuestionTagRepo.ListQuestionsByTag(id)
	if err != nil {
		return fmt.Errorf("failed to check tag usage: %w", err)
	}

	if len(questions) > 0 {
		return fmt.Errorf("cannot delete tag: it is associated with %d question(s)", len(questions))
	}

	// Deletar a tag
	err = uc.QuestionTagRepo.DeleteTag(id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// GetTag retorna uma tag por ID
func (uc *QuestionTagUsecase) GetTag(id int) (*domain.QuestionTag, error) {
	tag, err := uc.QuestionTagRepo.FindTagByID(id)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	return tag, nil
}

// ListAllTags retorna todas as tags
func (uc *QuestionTagUsecase) ListAllTags() ([]*domain.QuestionTag, error) {
	tags, err := uc.QuestionTagRepo.ListAllTags()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return tags, nil
}

// ListTagsWithStats retorna todas as tags com estatísticas de uso
func (uc *QuestionTagUsecase) ListTagsWithStats() ([]*domain.QuestionTagWithStats, error) {
	tags, err := uc.QuestionTagRepo.ListTagsWithStats()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags with stats: %w", err)
	}

	return tags, nil
}

// AssociateQuestionTag associa uma questão com uma tag
func (uc *QuestionTagUsecase) AssociateQuestionTag(questionID, tagID int) error {
	// Verificar se a questão existe
	_, err := uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return fmt.Errorf("question not found: %w", err)
	}

	// Verificar se a tag existe
	_, err = uc.QuestionTagRepo.FindTagByID(tagID)
	if err != nil {
		return fmt.Errorf("tag not found: %w", err)
	}

	// Criar a associação
	err = uc.QuestionTagRepo.AssociateQuestionTag(questionID, tagID)
	if err != nil {
		return fmt.Errorf("failed to associate question with tag: %w", err)
	}

	return nil
}

// DisassociateQuestionTag remove a associação entre uma questão e uma tag
func (uc *QuestionTagUsecase) DisassociateQuestionTag(questionID, tagID int) error {
	// Verificar se a associação existe (implicitamente verificado no repositório)
	err := uc.QuestionTagRepo.DisassociateQuestionTag(questionID, tagID)
	if err != nil {
		return fmt.Errorf("failed to disassociate question from tag: %w", err)
	}

	return nil
}

// GetQuestionTags retorna todas as tags de uma questão
func (uc *QuestionTagUsecase) GetQuestionTags(questionID int) ([]*domain.QuestionTag, error) {
	// Verificar se a questão existe
	_, err := uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}

	tags, err := uc.QuestionTagRepo.ListTagsByQuestion(questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get question tags: %w", err)
	}

	return tags, nil
}

// GetTagQuestions retorna todas as questões de uma tag
func (uc *QuestionTagUsecase) GetTagQuestions(tagID int) ([]*domain.Question, error) {
	// Verificar se a tag existe
	_, err := uc.QuestionTagRepo.FindTagByID(tagID)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	questions, err := uc.QuestionTagRepo.ListQuestionsByTag(tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag questions: %w", err)
	}

	return questions, nil
}

// UpdateQuestionTags atualiza todas as tags de uma questão
func (uc *QuestionTagUsecase) UpdateQuestionTags(questionID int, tagIDs []int) error {
	// Verificar se a questão existe
	_, err := uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return fmt.Errorf("question not found: %w", err)
	}

	// Verificar se todas as tags existem
	for _, tagID := range tagIDs {
		_, err := uc.QuestionTagRepo.FindTagByID(tagID)
		if err != nil {
			return fmt.Errorf("tag %d not found: %w", tagID, err)
		}
	}

	// Atualizar as associações
	err = uc.QuestionTagRepo.UpdateQuestionTags(questionID, tagIDs)
	if err != nil {
		return fmt.Errorf("failed to update question tags: %w", err)
	}

	return nil
}

// BulkAssociateQuestionTags associa uma questão com múltiplas tags
func (uc *QuestionTagUsecase) BulkAssociateQuestionTags(questionID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}

	// Verificar se a questão existe
	_, err := uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return fmt.Errorf("question not found: %w", err)
	}

	// Verificar se todas as tags existem
	for _, tagID := range tagIDs {
		_, err := uc.QuestionTagRepo.FindTagByID(tagID)
		if err != nil {
			return fmt.Errorf("tag %d not found: %w", tagID, err)
		}
	}

	// Criar as associações
	err = uc.QuestionTagRepo.BulkAssociateQuestionTags(questionID, tagIDs)
	if err != nil {
		return fmt.Errorf("failed to bulk associate question with tags: %w", err)
	}

	return nil
}

// BulkDisassociateQuestionTags remove múltiplas tags de uma questão
func (uc *QuestionTagUsecase) BulkDisassociateQuestionTags(questionID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}

	// Verificar se a questão existe
	_, err := uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return fmt.Errorf("question not found: %w", err)
	}

	// Remover as associações
	err = uc.QuestionTagRepo.BulkDisassociateQuestionTags(questionID, tagIDs)
	if err != nil {
		return fmt.Errorf("failed to bulk disassociate question from tags: %w", err)
	}

	return nil
}

// SearchTags busca tags por nome
func (uc *QuestionTagUsecase) SearchTags(query string, limit int) ([]*domain.QuestionTag, error) {
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return nil, fmt.Errorf("search query must be at least 2 characters long")
	}

	if limit <= 0 {
		limit = 20
	}

	tags, err := uc.QuestionTagRepo.SearchTags(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search tags: %w", err)
	}

	return tags, nil
}

// FilterQuestionsByTags filtra questões baseado em critérios de tags
func (uc *QuestionTagUsecase) FilterQuestionsByTags(filters *domain.QuestionTagFilters) ([]*domain.Question, error) {
	if filters == nil {
		return nil, fmt.Errorf("filters cannot be nil")
	}

	// Validar filtros
	if len(filters.TagIDs) == 0 && len(filters.TagNames) == 0 {
		return nil, fmt.Errorf("at least one tag ID or tag name must be provided")
	}

	// Verificar se as tags existem
	for _, tagID := range filters.TagIDs {
		_, err := uc.QuestionTagRepo.FindTagByID(tagID)
		if err != nil {
			return nil, fmt.Errorf("tag %d not found: %w", tagID, err)
		}
	}

	for _, tagName := range filters.TagNames {
		_, err := uc.QuestionTagRepo.FindTagByName(tagName)
		if err != nil {
			return nil, fmt.Errorf("tag '%s' not found: %w", tagName, err)
		}
	}

	// Aplicar filtros padrão
	if filters.Limit <= 0 {
		filters.Limit = 50
	}

	questions, err := uc.QuestionTagRepo.FilterQuestionsByTags(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to filter questions by tags: %w", err)
	}

	return questions, nil
}

// GetTagUsageStats retorna estatísticas de uso de uma tag
func (uc *QuestionTagUsecase) GetTagUsageStats(tagID int) (*domain.TagUsageStats, error) {
	// Verificar se a tag existe
	_, err := uc.QuestionTagRepo.FindTagByID(tagID)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	stats, err := uc.QuestionTagRepo.GetTagUsageStats(tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag usage stats: %w", err)
	}

	return stats, nil
}

// GetQuestionTagStats retorna estatísticas de tags de uma questão
func (uc *QuestionTagUsecase) GetQuestionTagStats(questionID int) (*domain.QuestionTagStats, error) {
	// Verificar se a questão existe
	_, err := uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}

	stats, err := uc.QuestionTagRepo.GetQuestionTagStats(questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get question tag stats: %w", err)
	}

	return stats, nil
}

// GetMostUsedTags retorna as tags mais utilizadas
func (uc *QuestionTagUsecase) GetMostUsedTags(limit int) ([]*domain.QuestionTagWithStats, error) {
	if limit <= 0 {
		limit = 10
	}

	tags, err := uc.QuestionTagRepo.GetMostUsedTags(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get most used tags: %w", err)
	}

	return tags, nil
}

// GetLeastUsedTags retorna as tags menos utilizadas
func (uc *QuestionTagUsecase) GetLeastUsedTags(limit int) ([]*domain.QuestionTagWithStats, error) {
	if limit <= 0 {
		limit = 10
	}

	tags, err := uc.QuestionTagRepo.GetLeastUsedTags(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get least used tags: %w", err)
	}

	return tags, nil
}

// SuggestTagsForQuestion sugere tags para uma questão baseado no conteúdo
func (uc *QuestionTagUsecase) SuggestTagsForQuestion(questionID int) ([]*domain.TagSuggestion, error) {
	// Buscar a questão
	question, err := uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}

	// Implementação simplificada de sugestões
	// Em uma implementação completa, isso poderia usar ML/NLP
	suggestions := []*domain.TagSuggestion{}

	// Sugestões baseadas no nível de dificuldade
	switch question.DifficultyLevel {
	case 1:
		suggestions = append(suggestions, &domain.TagSuggestion{
			TagName:    "Básico",
			Confidence: 0.9,
			Reason:     "Baseado no nível de dificuldade",
		})
	case 2:
		suggestions = append(suggestions, &domain.TagSuggestion{
			TagName:    "Intermediário",
			Confidence: 0.9,
			Reason:     "Baseado no nível de dificuldade",
		})
	case 3, 4, 5:
		suggestions = append(suggestions, &domain.TagSuggestion{
			TagName:    "Avançado",
			Confidence: 0.9,
			Reason:     "Baseado no nível de dificuldade",
		})
	}

	// Sugestões baseadas no tipo de conteúdo
	if question.ContentType == domain.QuestionContentTypeCode {
		suggestions = append(suggestions, &domain.TagSuggestion{
			TagName:    "Prático",
			Confidence: 0.8,
			Reason:     "Questão contém código",
		})
	} else {
		suggestions = append(suggestions, &domain.TagSuggestion{
			TagName:    "Teórico",
			Confidence: 0.7,
			Reason:     "Questão é baseada em texto",
		})
	}

	// Buscar tags existentes para validar sugestões
	existingTags, err := uc.QuestionTagRepo.ListAllTags()
	if err != nil {
		return suggestions, nil // Retornar sugestões mesmo se não conseguir validar
	}

	// Filtrar apenas tags que existem
	var validSuggestions []*domain.TagSuggestion
	for _, suggestion := range suggestions {
		for _, tag := range existingTags {
			if tag.Name == suggestion.TagName {
				suggestion.TagID = tag.ID
				validSuggestions = append(validSuggestions, suggestion)
				break
			}
		}
	}

	return validSuggestions, nil
}

// FindSimilarQuestionsByTags encontra questões similares baseado na sobreposição de tags
func (uc *QuestionTagUsecase) FindSimilarQuestionsByTags(questionID int, limit int) ([]*domain.Question, error) {
	// Verificar se a questão existe
	_, err := uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}

	// Obter tags da questão
	tags, err := uc.QuestionTagRepo.ListTagsByQuestion(questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get question tags: %w", err)
	}

	if len(tags) == 0 {
		return []*domain.Question{}, nil
	}

	// Obter IDs das tags
	tagIDs := make([]int, len(tags))
	for i, tag := range tags {
		tagIDs[i] = tag.ID
	}

	// Buscar questões com tags similares
	similarQuestions, err := uc.QuestionTagRepo.ListQuestionsByTags(tagIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to find similar questions: %w", err)
	}

	// Filtrar a questão original e aplicar limite
	var result []*domain.Question
	for _, q := range similarQuestions {
		if q.ID != questionID {
			result = append(result, q)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}

	return result, nil
}

// CleanupUnusedTags remove tags que não estão sendo utilizadas
func (uc *QuestionTagUsecase) CleanupUnusedTags() (int, error) {
	// Buscar tags não utilizadas
	tags, err := uc.QuestionTagRepo.GetLeastUsedTags(1000) // Buscar muitas para encontrar as não utilizadas
	if err != nil {
		return 0, fmt.Errorf("failed to get unused tags: %w", err)
	}

	deletedCount := 0
	for _, tag := range tags {
		if tag.QuestionCount == 0 {
			err := uc.QuestionTagRepo.DeleteTag(tag.ID)
			if err != nil {
				fmt.Printf("Warning: failed to delete unused tag %d: %v\n", tag.ID, err)
				continue
			}
			deletedCount++
		}
	}

	return deletedCount, nil
}
