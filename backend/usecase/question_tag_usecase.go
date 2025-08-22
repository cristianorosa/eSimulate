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
	count, err := uc.QuestionTagRepo.CountTagsByQuestion(id)
	if err != nil {
		return fmt.Errorf("failed to check tag usage: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete tag: it is associated with %d questions", count)
	}

	// Deletar a tag
	err = uc.QuestionTagRepo.DeleteTag(id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// GetTagByID busca uma tag por ID
func (uc *QuestionTagUsecase) GetTagByID(id int) (*domain.QuestionTag, error) {
	tag, err := uc.QuestionTagRepo.FindTagByID(id)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}
	return tag, nil
}

// GetTagByName busca uma tag por nome
func (uc *QuestionTagUsecase) GetTagByName(name string) (*domain.QuestionTag, error) {
	tag, err := uc.QuestionTagRepo.FindTagByName(name)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}
	return tag, nil
}

// ListTags lista todas as tags com paginação
func (uc *QuestionTagUsecase) ListTags(page, pageSize int) ([]*domain.QuestionTag, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	tags, err := uc.QuestionTagRepo.ListTags(page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return tags, nil
}

// ListTagsWithStats lista tags com estatísticas de uso
func (uc *QuestionTagUsecase) ListTagsWithStats() ([]*domain.QuestionTagWithStats, error) {
	tags, err := uc.QuestionTagRepo.ListTagsWithStats()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags with stats: %w", err)
	}

	return tags, nil
}

// AssociateQuestionTag associa uma tag a uma questão
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

	// Associar
	err = uc.QuestionTagRepo.AssociateQuestionTag(questionID, tagID)
	if err != nil {
		return fmt.Errorf("failed to associate question tag: %w", err)
	}

	return nil
}

// DisassociateQuestionTag remove associação entre tag e questão
func (uc *QuestionTagUsecase) DisassociateQuestionTag(questionID, tagID int) error {
	err := uc.QuestionTagRepo.DisassociateQuestionTag(questionID, tagID)
	if err != nil {
		return fmt.Errorf("failed to disassociate question tag: %w", err)
	}

	return nil
}

// GetQuestionTags lista todas as tags de uma questão
func (uc *QuestionTagUsecase) GetQuestionTags(questionID int) ([]*domain.QuestionTag, error) {
	// Verificar se a questão existe
	_, err := uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}

	tags, err := uc.QuestionTagRepo.ListQuestionTags(questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get question tags: %w", err)
	}

	return tags, nil
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
		_, err = uc.QuestionTagRepo.FindTagByID(tagID)
		if err != nil {
			return fmt.Errorf("tag %d not found: %w", tagID, err)
		}
	}

	// Atualizar associações
	err = uc.QuestionTagRepo.UpdateQuestionTags(questionID, tagIDs)
	if err != nil {
		return fmt.Errorf("failed to update question tags: %w", err)
	}

	return nil
}

// BulkAssociateQuestionTags associa múltiplas tags a uma questão
func (uc *QuestionTagUsecase) BulkAssociateQuestionTags(questionID int, tagIDs []int) error {
	// Verificar se a questão existe
	_, err := uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return fmt.Errorf("question not found: %w", err)
	}

	// Verificar se todas as tags existem
	for _, tagID := range tagIDs {
		_, err = uc.QuestionTagRepo.FindTagByID(tagID)
		if err != nil {
			return fmt.Errorf("tag %d not found: %w", tagID, err)
		}
	}

	// Associar em lote
	err = uc.QuestionTagRepo.BulkAssociateQuestionTags(questionID, tagIDs)
	if err != nil {
		return fmt.Errorf("failed to bulk associate question tags: %w", err)
	}

	return nil
}

// BulkDisassociateQuestionTags remove múltiplas associações de uma questão
func (uc *QuestionTagUsecase) BulkDisassociateQuestionTags(questionID int, tagIDs []int) error {
	err := uc.QuestionTagRepo.BulkDisassociateQuestionTags(questionID, tagIDs)
	if err != nil {
		return fmt.Errorf("failed to bulk disassociate question tags: %w", err)
	}

	return nil
}

// SearchTags busca tags por nome
func (uc *QuestionTagUsecase) SearchTags(searchTerm string, limit int) ([]*domain.QuestionTag, error) {
	searchTerm = strings.TrimSpace(searchTerm)
	if len(searchTerm) < 2 {
		return nil, fmt.Errorf("search term must be at least 2 characters long")
	}

	if limit <= 0 || limit > 50 {
		limit = 20
	}

	tags, err := uc.QuestionTagRepo.SearchTags(searchTerm, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search tags: %w", err)
	}

	return tags, nil
}

// GetTagStats retorna estatísticas de uma tag
func (uc *QuestionTagUsecase) GetTagStats(tagID int) (*domain.QuestionTagWithStats, error) {
	stats, err := uc.QuestionTagRepo.GetTagStats(tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag stats: %w", err)
	}

	return stats, nil
}

// CleanupUnusedTags remove tags que não estão sendo usadas
func (uc *QuestionTagUsecase) CleanupUnusedTags() (int, error) {
	count, err := uc.QuestionTagRepo.CleanupUnusedTags()
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup unused tags: %w", err)
	}

	return count, nil
}

// GetTagsByQuestions retorna tags para múltiplas questões
func (uc *QuestionTagUsecase) GetTagsByQuestions(questionIDs []int) (map[int][]*domain.QuestionTag, error) {
	if len(questionIDs) == 0 {
		return make(map[int][]*domain.QuestionTag), nil
	}

	tags, err := uc.QuestionTagRepo.GetTagsByQuestions(questionIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags by questions: %w", err)
	}

	return tags, nil
}
