package usecase

import (
	"fmt"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// ExamQuestionUsecase implementa a lógica de negócio para relacionamentos exame-questão
type ExamQuestionUsecase struct {
	ExamQuestionRepo domain.ExamQuestionRepository
	ExamRepo         domain.ExamRepository
	QuestionRepo     domain.QuestionRepository
	ExamTopicRepo    domain.ExamTopicRepository
}

// NewExamQuestionUsecase cria uma nova instância do use case
func NewExamQuestionUsecase(
	examQuestionRepo domain.ExamQuestionRepository,
	examRepo domain.ExamRepository,
	questionRepo domain.QuestionRepository,
	examTopicRepo domain.ExamTopicRepository,
) *ExamQuestionUsecase {
	return &ExamQuestionUsecase{
		ExamQuestionRepo: examQuestionRepo,
		ExamRepo:         examRepo,
		QuestionRepo:     questionRepo,
		ExamTopicRepo:    examTopicRepo,
	}
}

// AssociateQuestionWithExam associa uma questão a um exame
func (uc *ExamQuestionUsecase) AssociateQuestionWithExam(examID, questionID int) error {
	// Validar se o exame existe
	_, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		return fmt.Errorf("exam not found: %w", err)
	}

	// Validar se a questão existe
	_, err = uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return fmt.Errorf("question not found: %w", err)
	}

	// Verificar se a associação já existe
	existing, err := uc.ExamQuestionRepo.FindByExamAndQuestion(examID, questionID)
	if err == nil && existing != nil {
		return fmt.Errorf("question is already associated with this exam")
	}

	// Validar regras de negócio (hierarquia tópico-exame)
	err = uc.validateExamQuestionAssociation(examID, questionID)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Obter próximo order_index
	examQuestions, err := uc.ExamQuestionRepo.ListByExam(examID)
	if err != nil {
		return fmt.Errorf("failed to get existing questions: %w", err)
	}

	nextOrderIndex := len(examQuestions) + 1

	// Criar a associação
	examQuestion := &domain.ExamQuestion{
		ExamID:     examID,
		QuestionID: questionID,
		OrderIndex: nextOrderIndex,
	}

	err = uc.ExamQuestionRepo.Create(examQuestion)
	if err != nil {
		return fmt.Errorf("failed to create exam-question association: %w", err)
	}

	// Atualizar contador de questões do exame
	err = uc.ExamRepo.UpdateQuestionsCount(examID)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to update exam questions count: %v\n", err)
	}

	return nil
}

// DisassociateQuestionFromExam remove uma questão de um exame
func (uc *ExamQuestionUsecase) DisassociateQuestionFromExam(examID, questionID int) error {
	// Verificar se a associação existe
	_, err := uc.ExamQuestionRepo.FindByExamAndQuestion(examID, questionID)
	if err != nil {
		return fmt.Errorf("exam-question association not found: %w", err)
	}

	// Remover a associação
	err = uc.ExamQuestionRepo.Delete(examID, questionID)
	if err != nil {
		return fmt.Errorf("failed to delete exam-question association: %w", err)
	}

	// Reordenar as questões restantes
	examQuestions, err := uc.ExamQuestionRepo.ListByExam(examID)
	if err != nil {
		return fmt.Errorf("failed to get remaining questions: %w", err)
	}

	questionIDs := make([]int, len(examQuestions))
	for i, eq := range examQuestions {
		questionIDs[i] = eq.QuestionID
	}

	err = uc.ExamQuestionRepo.ReorderQuestions(examID, questionIDs)
	if err != nil {
		return fmt.Errorf("failed to reorder questions: %w", err)
	}

	// Atualizar contador de questões do exame
	err = uc.ExamRepo.UpdateQuestionsCount(examID)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to update exam questions count: %v\n", err)
	}

	return nil
}

// GetExamQuestions retorna todas as questões de um exame com detalhes
func (uc *ExamQuestionUsecase) GetExamQuestions(examID int) ([]*domain.ExamQuestionWithDetails, error) {
	// Validar se o exame existe
	_, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		return nil, fmt.Errorf("exam not found: %w", err)
	}

	questions, err := uc.ExamQuestionRepo.ListByExamWithDetails(examID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam questions: %w", err)
	}

	return questions, nil
}

// GetQuestionExams retorna todos os exames que contêm uma questão
func (uc *ExamQuestionUsecase) GetQuestionExams(questionID int) ([]*domain.ExamQuestion, error) {
	// Validar se a questão existe
	_, err := uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}

	exams, err := uc.ExamQuestionRepo.ListByQuestion(questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get question exams: %w", err)
	}

	return exams, nil
}

// ReorderExamQuestions reordena as questões de um exame
func (uc *ExamQuestionUsecase) ReorderExamQuestions(examID int, questionIDs []int) error {
	// Validar se o exame existe
	_, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		return fmt.Errorf("exam not found: %w", err)
	}

	// Validar se todas as questões pertencem ao exame
	for _, questionID := range questionIDs {
		_, err := uc.ExamQuestionRepo.FindByExamAndQuestion(examID, questionID)
		if err != nil {
			return fmt.Errorf("question %d is not associated with exam %d: %w", questionID, examID, err)
		}
	}

	// Reordenar as questões
	err = uc.ExamQuestionRepo.ReorderQuestions(examID, questionIDs)
	if err != nil {
		return fmt.Errorf("failed to reorder questions: %w", err)
	}

	return nil
}

// BulkAssociateQuestionsWithExam associa múltiplas questões a um exame
func (uc *ExamQuestionUsecase) BulkAssociateQuestionsWithExam(examID int, questionIDs []int) error {
	if len(questionIDs) == 0 {
		return nil
	}

	// Validar se o exame existe
	_, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		return fmt.Errorf("exam not found: %w", err)
	}

	// Validar cada questão e regras de negócio
	for _, questionID := range questionIDs {
		// Verificar se a questão existe
		_, err := uc.QuestionRepo.FindByID(questionID)
		if err != nil {
			return fmt.Errorf("question %d not found: %w", questionID, err)
		}

		// Verificar se já existe associação
		existing, err := uc.ExamQuestionRepo.FindByExamAndQuestion(examID, questionID)
		if err == nil && existing != nil {
			return fmt.Errorf("question %d is already associated with exam %d", questionID, examID)
		}

		// Validar regras de negócio
		err = uc.validateExamQuestionAssociation(examID, questionID)
		if err != nil {
			return fmt.Errorf("validation failed for question %d: %w", questionID, err)
		}
	}

	// Obter order_index inicial
	examQuestions, err := uc.ExamQuestionRepo.ListByExam(examID)
	if err != nil {
		return fmt.Errorf("failed to get existing questions: %w", err)
	}

	startOrderIndex := len(examQuestions) + 1

	// Criar as associações
	var newExamQuestions []*domain.ExamQuestion
	for i, questionID := range questionIDs {
		examQuestion := &domain.ExamQuestion{
			ExamID:     examID,
			QuestionID: questionID,
			OrderIndex: startOrderIndex + i,
		}
		newExamQuestions = append(newExamQuestions, examQuestion)
	}

	err = uc.ExamQuestionRepo.BulkCreate(newExamQuestions)
	if err != nil {
		return fmt.Errorf("failed to create exam-question associations: %w", err)
	}

	// Atualizar contador de questões do exame
	err = uc.ExamRepo.UpdateQuestionsCount(examID)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to update exam questions count: %v\n", err)
	}

	return nil
}

// BulkDisassociateQuestionsFromExam remove múltiplas questões de um exame
func (uc *ExamQuestionUsecase) BulkDisassociateQuestionsFromExam(examID int, questionIDs []int) error {
	if len(questionIDs) == 0 {
		return nil
	}

	// Validar se o exame existe
	_, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		return fmt.Errorf("exam not found: %w", err)
	}

	// Validar se todas as questões estão associadas ao exame
	for _, questionID := range questionIDs {
		_, err := uc.ExamQuestionRepo.FindByExamAndQuestion(examID, questionID)
		if err != nil {
			return fmt.Errorf("question %d is not associated with exam %d: %w", questionID, examID, err)
		}
	}

	// Remover as associações
	err = uc.ExamQuestionRepo.BulkDelete(examID, questionIDs)
	if err != nil {
		return fmt.Errorf("failed to delete exam-question associations: %w", err)
	}

	// Reordenar as questões restantes
	examQuestions, err := uc.ExamQuestionRepo.ListByExam(examID)
	if err != nil {
		return fmt.Errorf("failed to get remaining questions: %w", err)
	}

	remainingQuestionIDs := make([]int, len(examQuestions))
	for i, eq := range examQuestions {
		remainingQuestionIDs[i] = eq.QuestionID
	}

	err = uc.ExamQuestionRepo.ReorderQuestions(examID, remainingQuestionIDs)
	if err != nil {
		return fmt.Errorf("failed to reorder questions: %w", err)
	}

	// Atualizar contador de questões do exame
	err = uc.ExamRepo.UpdateQuestionsCount(examID)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to update exam questions count: %v\n", err)
	}

	return nil
}

// GetExamQuestionStats retorna estatísticas das questões de um exame
func (uc *ExamQuestionUsecase) GetExamQuestionStats(examID int) (*domain.ExamQuestionStats, error) {
	// Validar se o exame existe
	_, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		return nil, fmt.Errorf("exam not found: %w", err)
	}

	// Obter questões do exame
	questions, err := uc.ExamQuestionRepo.ListByExamWithDetails(examID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam questions: %w", err)
	}

	// Calcular estatísticas
	stats := &domain.ExamQuestionStats{
		ExamID:                examID,
		TotalQuestions:        len(questions),
		QuestionsByTopic:      make(map[int]int),
		QuestionsByDifficulty: make(map[int]int),
		TopicDistribution:     []*domain.TopicDistribution{},
	}

	// Contar questões por tópico e dificuldade
	topicCounts := make(map[int]int)
	topicNames := make(map[int]string)

	for _, q := range questions {
		stats.QuestionsByTopic[q.TopicID]++
		stats.QuestionsByDifficulty[q.DifficultyLevel]++
		topicCounts[q.TopicID]++
		topicNames[q.TopicID] = q.TopicName
	}

	// Obter configuração dos tópicos do exame
	examTopics, err := uc.ExamTopicRepo.ListByExam(examID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam topics: %w", err)
	}

	// Criar distribuição de tópicos
	for _, et := range examTopics {
		actualCount := topicCounts[et.TopicID]
		topicName := topicNames[et.TopicID]

		distribution := &domain.TopicDistribution{
			TopicID:          et.TopicID,
			TopicName:        topicName,
			QuestionCount:    actualCount,
			WeightPercentage: et.WeightPercentage,
			ExpectedCount:    et.QuestionsCount,
			IsComplete:       actualCount >= et.QuestionsCount,
		}

		stats.TopicDistribution = append(stats.TopicDistribution, distribution)
	}

	return stats, nil
}

// GetAvailableQuestionsForExam retorna questões disponíveis para associar a um exame
func (uc *ExamQuestionUsecase) GetAvailableQuestionsForExam(examID int, topicID *int) ([]*domain.Question, error) {
	// Validar se o exame existe
	_, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		return nil, fmt.Errorf("exam not found: %w", err)
	}

	// Obter questões disponíveis
	questions, err := uc.QuestionRepo.GetAvailableQuestionsForExam(examID, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available questions: %w", err)
	}

	return questions, nil
}

// ValidateExamConfiguration valida se a configuração do exame está correta
func (uc *ExamQuestionUsecase) ValidateExamConfiguration(examID int) (bool, []string) {
	var errors []string

	// Verificar se o exame existe
	_, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Exam not found: %v", err))
		return false, errors
	}

	// Verificar distribuição de pesos dos tópicos
	err = uc.ExamTopicRepo.ValidateWeightDistribution(examID)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Weight distribution error: %v", err))
	}

	// Verificar se há questões suficientes para cada tópico
	examTopics, err := uc.ExamTopicRepo.ListByExamWithDetails(examID)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to get exam topics: %v", err))
		return false, errors
	}

	for _, et := range examTopics {
		if et.ActualQuestions < et.QuestionsCount {
			errors = append(errors, fmt.Sprintf("Topic '%s' has only %d questions but needs %d",
				et.TopicName, et.ActualQuestions, et.QuestionsCount))
		}
	}

	// Verificar total de questões
	totalQuestions, err := uc.ExamQuestionRepo.GetExamQuestionCount(examID)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to count exam questions: %v", err))
		return false, errors
	}

	if totalQuestions == 0 {
		errors = append(errors, "Exam has no questions")
	}

	return len(errors) == 0, errors
}

// validateExamQuestionAssociation validates business rules for exam-question association
func (uc *ExamQuestionUsecase) validateExamQuestionAssociation(examID, questionID int) error {
	// Get the question to check its topic
	question, err := uc.QuestionRepo.FindByID(questionID)
	if err != nil {
		return fmt.Errorf("question not found: %w", err)
	}

	// Check if question has a topic
	if question.TopicID == 0 {
		return fmt.Errorf("question %d does not have a topic assigned", questionID)
	}

	// Get exam topics
	examTopics, err := uc.ExamTopicRepo.ListByExam(examID)
	if err != nil {
		return fmt.Errorf("failed to get exam topics: %w", err)
	}

	// If exam has no topics, reject the association
	if len(examTopics) == 0 {
		return fmt.Errorf("exam %d has no topics configured", examID)
	}

	// Check if exam has only one topic with 100% weight (exception case)
	if len(examTopics) == 1 && examTopics[0].WeightPercentage == 100.0 {
		// Allow any question for single topic with 100% weight
		return nil
	}

	// Check if the question's topic is associated with the exam
	for _, examTopic := range examTopics {
		if examTopic.TopicID == question.TopicID {
			return nil // Valid association
		}
	}

	// Topic is not associated with exam, reject
	return fmt.Errorf("question %d belongs to topic %d which is not associated with exam %d",
		questionID, question.TopicID, examID)
}
