package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// QuestionUsecase implementa as regras de negócio para questões.
type QuestionUsecase struct {
	Repo       domain.QuestionRepository
	OptionRepo domain.OptionRepository
}

// CreateQuestion cria uma nova questão com suas opções.
func (uc *QuestionUsecase) CreateQuestion(ctx context.Context, topicID int, statement, problem, contentType, explanation, questionType string, difficultyLevel, createdBy int, options []*domain.Option) (*domain.Question, error) {
	log.Printf("Criando questão: topicID=%d, statement=%s, problem=%s, contentType=%s, questionType=%s, difficultyLevel=%d, options=%d",
		topicID, statement[:min(50, len(statement))], problem[:min(50, len(problem))], contentType, questionType, difficultyLevel, len(options))

	// Validar regras de negócio
	if len(options) < 2 {
		return nil, domain.ErrInvalidData
	}

	// Validar tipo de questão
	correctOptions := 0
	for _, opt := range options {
		if opt.IsCorrect {
			correctOptions++
		}
	}

	if questionType == string(domain.QuestionTypeObjective) && correctOptions != 1 {
		return nil, domain.ErrInvalidData
	}

	if questionType == string(domain.QuestionTypeMultipleChoice) && correctOptions < 1 {
		return nil, domain.ErrInvalidData
	}

	// Criar questão
	question := &domain.Question{
		TopicID:         topicID,
		Statement:       statement,
		Problem:         problem,
		ContentType:     domain.QuestionContentType(contentType),
		Explanation:     explanation,
		QuestionType:    domain.QuestionType(questionType),
		DifficultyLevel: difficultyLevel,
		CreatedBy:       createdBy,
		IsActive:        true,
	}

	// Salvar questão no repositório
	err := uc.Repo.Create(question)
	if err != nil {
		log.Printf("Erro ao criar questão: %v", err)
		return nil, err
	}

	// Salvar opções
	for _, option := range options {
		option.QuestionID = question.ID
		err := uc.OptionRepo.Create(option)
		if err != nil {
			log.Printf("Erro ao criar opção: %v", err)
			return nil, err
		}
	}

	// Carregar opções salvas
	savedOptions, err := uc.OptionRepo.FindByQuestionID(question.ID)
	if err != nil {
		log.Printf("Erro ao carregar opções: %v", err)
		return nil, err
	}
	question.Options = savedOptions

	log.Printf("Questão criada com sucesso: ID=%d, opções=%d", question.ID, len(savedOptions))
	return question, nil
}

// UpdateQuestion atualiza uma questão existente.
func (uc *QuestionUsecase) UpdateQuestion(ctx context.Context, id, topicID int, statement, problem, contentType, explanation, questionType string, difficultyLevel int, isActive bool) error {
	log.Printf("Atualizando questão: ID=%d, topicID=%d", id, topicID)

	// Buscar questão existente
	existingQuestion, err := uc.Repo.FindByID(id)
	if err != nil {
		log.Printf("Erro ao buscar questão para atualização: %v", err)
		return err
	}

	// Atualizar campos
	existingQuestion.TopicID = topicID
	existingQuestion.Statement = statement
	existingQuestion.Problem = problem
	existingQuestion.ContentType = domain.QuestionContentType(contentType)
	existingQuestion.Explanation = explanation
	existingQuestion.QuestionType = domain.QuestionType(questionType)
	existingQuestion.DifficultyLevel = difficultyLevel
	existingQuestion.IsActive = isActive

	// Salvar atualização
	err = uc.Repo.Update(existingQuestion)
	if err != nil {
		log.Printf("Erro ao atualizar questão: %v", err)
		return err
	}

	log.Printf("Questão atualizada com sucesso: ID=%d", id)
	return nil
}

// DeleteQuestion remove uma questão.
func (uc *QuestionUsecase) DeleteQuestion(ctx context.Context, id int) error {
	log.Printf("Deletando questão: ID=%d", id)

	// Deletar opções primeiro (devido à foreign key)
	err := uc.OptionRepo.DeleteByQuestionID(id)
	if err != nil {
		log.Printf("Erro ao deletar opções da questão: %v", err)
		return err
	}

	// Deletar questão
	err = uc.Repo.Delete(id)
	if err != nil {
		log.Printf("Erro ao deletar questão: %v", err)
		return err
	}

	log.Printf("Questão deletada com sucesso: ID=%d", id)
	return nil
}

// ListQuestions lista todas as questões com detalhes.
func (uc *QuestionUsecase) ListQuestions(ctx context.Context, examID *int) ([]*domain.QuestionWithDetails, error) {
	log.Printf("Listando questões: examID=%v", examID)

	questions, err := uc.Repo.ListAllWithDetails()
	if err != nil {
		log.Printf("Erro ao listar questões: %v", err)
		return nil, err
	}

	// Filtrar por exame se especificado
	if examID != nil {
		filteredQuestions := make([]*domain.QuestionWithDetails, 0)
		for _, q := range questions {
			if q.ExamID == *examID {
				filteredQuestions = append(filteredQuestions, q)
			}
		}
		questions = filteredQuestions
	}

	// Carregar opções para cada questão
	for _, question := range questions {
		options, err := uc.OptionRepo.FindByQuestionID(question.ID)
		if err != nil {
			log.Printf("Erro ao carregar opções da questão %d: %v", question.ID, err)
			continue
		}
		question.Options = options
	}

	log.Printf("Questões listadas com sucesso: %d questões", len(questions))
	return questions, nil
}

// ListQuestionsPaginated lista questões com paginação.
func (uc *QuestionUsecase) ListQuestionsPaginated(ctx context.Context, page, pageSize int, examID, topicID *int) ([]*domain.QuestionWithDetails, *domain.Pagination, error) {
	log.Printf("Listando questões paginadas: page=%d, pageSize=%d, examID=%v, topicID=%v", page, pageSize, examID, topicID)

	questions, pagination, err := uc.Repo.ListPaginated(page, pageSize, examID, topicID)
	if err != nil {
		log.Printf("Erro ao listar questões paginadas: %v", err)
		return nil, nil, err
	}

	// Carregar opções para cada questão
	for _, question := range questions {
		options, err := uc.OptionRepo.FindByQuestionID(question.ID)
		if err != nil {
			log.Printf("Erro ao carregar opções da questão %d: %v", question.ID, err)
			continue
		}
		question.Options = options
	}

	log.Printf("Questões paginadas listadas com sucesso: %d questões", len(questions))
	return questions, pagination, nil
}

// GetQuestion busca uma questão pelo seu ID.
func (uc *QuestionUsecase) GetQuestion(ctx context.Context, id int) (*domain.Question, error) {
	log.Printf("Buscando questão: ID=%d", id)

	question, err := uc.Repo.FindByID(id)
	if err != nil {
		log.Printf("Erro ao buscar questão: %v", err)
		return nil, err
	}

	// Carregar opções da questão
	options, err := uc.OptionRepo.FindByQuestionID(id)
	if err != nil {
		log.Printf("Erro ao carregar opções da questão %d: %v", id, err)
		return nil, err
	}
	question.Options = options

	log.Printf("Questão encontrada: ID=%d, opções=%d", id, len(options))
	return question, nil
}

// Função auxiliar para min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ImportQuestions importa um conjunto de questões via JSON
func (uc *QuestionUsecase) ImportQuestions(ctx context.Context, importData interface{}) (map[string]interface{}, error) {
	log.Printf("Iniciando importação de questões")

	payload, ok := importData.(domain.ImportQuestionsPayload)
	if !ok {
		return nil, domain.ErrInvalidData
	}

	result := map[string]interface{}{
		"areas_created":     0,
		"exams_created":     0,
		"topics_created":    0,
		"questions_created": 0,
		"errors":            []string{},
	}

	// Apenas contabiliza e valida minimamente por enquanto
	if payload.Area.Name == "" || payload.Exam.Title == "" {
		return nil, domain.ErrInvalidData
	}
	for _, t := range payload.Topics {
		for _, q := range t.Questions {
			if q.Statement == "" || q.Problem == "" || len(q.Options) < 1 {
				return nil, fmt.Errorf("questão inválida no tópico %s", t.Name)
			}
			result["questions_created"] = result["questions_created"].(int) + 1
		}
		result["topics_created"] = result["topics_created"].(int) + 1
	}

	log.Printf("Importação concluída (pré-persistência): %+v", result)
	return result, nil
}
