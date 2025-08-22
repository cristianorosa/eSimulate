package usecase

import (
	"fmt"

	"github.com/cristianorosa/eSimulate/backend/domain"
)

// ExamTopicUsecase implementa a lógica de negócio para relacionamentos exame-tópico
type ExamTopicUsecase struct {
	ExamTopicRepo domain.ExamTopicRepository
	ExamRepo      domain.ExamRepository
	TopicRepo     domain.TopicRepository
}

// NewExamTopicUsecase cria uma nova instância do use case
func NewExamTopicUsecase(
	examTopicRepo domain.ExamTopicRepository,
	examRepo domain.ExamRepository,
	topicRepo domain.TopicRepository,
) *ExamTopicUsecase {
	return &ExamTopicUsecase{
		ExamTopicRepo: examTopicRepo,
		ExamRepo:      examRepo,
		TopicRepo:     topicRepo,
	}
}

// AssociateTopicWithExam associa um tópico a um exame com configurações
func (uc *ExamTopicUsecase) AssociateTopicWithExam(examTopic *domain.ExamTopic) error {
	// Validar se o exame existe
	_, err := uc.ExamRepo.FindByID(examTopic.ExamID)
	if err != nil {
		return fmt.Errorf("exam not found: %w", err)
	}

	// Validar se o tópico existe
	_, err = uc.TopicRepo.FindByID(examTopic.TopicID)
	if err != nil {
		return fmt.Errorf("topic not found: %w", err)
	}

	// Verificar se a associação já existe
	existing, err := uc.ExamTopicRepo.FindByExamAndTopic(examTopic.ExamID, examTopic.TopicID)
	if err == nil && existing != nil {
		return fmt.Errorf("topic is already associated with this exam")
	}

	// Validar distribuição de dificuldade
	err = uc.ValidateDifficultyDistribution(examTopic)
	if err != nil {
		return fmt.Errorf("difficulty distribution validation failed: %w", err)
	}

	// Validar peso percentual
	if examTopic.WeightPercentage <= 0 || examTopic.WeightPercentage > 100 {
		return fmt.Errorf("weight percentage must be between 0 and 100")
	}

	// Validar número de questões
	if examTopic.QuestionsCount < 0 {
		return fmt.Errorf("questions count cannot be negative")
	}

	// Obter próximo order_index
	examTopics, err := uc.ExamTopicRepo.ListByExam(examTopic.ExamID)
	if err != nil {
		return fmt.Errorf("failed to get existing topics: %w", err)
	}

	if examTopic.OrderIndex <= 0 {
		examTopic.OrderIndex = len(examTopics) + 1
	}

	// Criar a associação
	err = uc.ExamTopicRepo.Create(examTopic)
	if err != nil {
		return fmt.Errorf("failed to create exam-topic association: %w", err)
	}

	return nil
}

// DisassociateTopicFromExam remove um tópico de um exame
func (uc *ExamTopicUsecase) DisassociateTopicFromExam(examID, topicID int) error {
	// Verificar se a associação existe
	_, err := uc.ExamTopicRepo.FindByExamAndTopic(examID, topicID)
	if err != nil {
		return fmt.Errorf("exam-topic association not found: %w", err)
	}

	// TODO: Verificar se há questões associadas que dependem desta associação
	// Isso pode ser implementado quando integrarmos com ExamQuestionRepository

	// Remover a associação
	err = uc.ExamTopicRepo.Delete(examID, topicID)
	if err != nil {
		return fmt.Errorf("failed to delete exam-topic association: %w", err)
	}

	return nil
}

// UpdateExamTopic atualiza uma associação exame-tópico existente
func (uc *ExamTopicUsecase) UpdateExamTopic(examTopic *domain.ExamTopic) error {
	// Verificar se a associação existe
	existing, err := uc.ExamTopicRepo.FindByExamAndTopic(examTopic.ExamID, examTopic.TopicID)
	if err != nil {
		return fmt.Errorf("exam-topic association not found: %w", err)
	}

	// Validar distribuição de dificuldade
	err = uc.ValidateDifficultyDistribution(examTopic)
	if err != nil {
		return fmt.Errorf("difficulty distribution validation failed: %w", err)
	}

	// Validar peso percentual
	if examTopic.WeightPercentage <= 0 || examTopic.WeightPercentage > 100 {
		return fmt.Errorf("weight percentage must be between 0 and 100")
	}

	// Validar número de questões
	if examTopic.QuestionsCount < 0 {
		return fmt.Errorf("questions count cannot be negative")
	}

	// Preservar timestamps
	examTopic.CreatedAt = existing.CreatedAt

	// Atualizar a associação
	err = uc.ExamTopicRepo.Update(examTopic)
	if err != nil {
		return fmt.Errorf("failed to update exam-topic association: %w", err)
	}

	return nil
}

// GetExamTopics retorna todos os tópicos de um exame com detalhes
func (uc *ExamTopicUsecase) GetExamTopics(examID int) ([]*domain.ExamTopicWithDetails, error) {
	// Validar se o exame existe
	_, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		return nil, fmt.Errorf("exam not found: %w", err)
	}

	topics, err := uc.ExamTopicRepo.ListByExamWithDetails(examID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam topics: %w", err)
	}

	return topics, nil
}

// GetTopicExams retorna todos os exames que contêm um tópico
func (uc *ExamTopicUsecase) GetTopicExams(topicID int) ([]*domain.ExamTopicWithDetails, error) {
	// Validar se o tópico existe
	_, err := uc.TopicRepo.FindByID(topicID)
	if err != nil {
		return nil, fmt.Errorf("topic not found: %w", err)
	}

	exams, err := uc.ExamTopicRepo.ListByTopicWithDetails(topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic exams: %w", err)
	}

	return exams, nil
}

// ValidateDifficultyDistribution valida se a distribuição de dificuldade é válida
func (uc *ExamTopicUsecase) ValidateDifficultyDistribution(examTopic *domain.ExamTopic) error {
	return uc.ExamTopicRepo.ValidateDifficultyDistribution(examTopic)
}

// CalculateQuestionsByDifficulty calcula quantas questões são necessárias para cada nível de dificuldade
func (uc *ExamTopicUsecase) CalculateQuestionsByDifficulty(examTopic *domain.ExamTopic) (*domain.DifficultyDistribution, error) {
	// Validar primeiro
	err := uc.ValidateDifficultyDistribution(examTopic)
	if err != nil {
		return nil, err
	}

	return uc.ExamTopicRepo.CalculateQuestionsByDifficulty(examTopic)
}

// GetExamStats retorna estatísticas completas de um exame
func (uc *ExamTopicUsecase) GetExamStats(examID int) (*domain.ExamTopicStats, error) {
	// Validar se o exame existe
	exam, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		return nil, fmt.Errorf("exam not found: %w", err)
	}

	// Obter tópicos do exame com detalhes
	topics, err := uc.ExamTopicRepo.ListByExamWithDetails(examID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam topics: %w", err)
	}

	// Calcular estatísticas
	stats := &domain.ExamTopicStats{
		ExamID:              examID,
		TotalTopics:         len(topics),
		TotalWeight:         0,
		TotalQuestions:      0,
		CompletedTopics:     0,
		TopicDistribution:   []*domain.TopicDistribution{},
		DifficultyBreakdown: &domain.DifficultyBreakdown{},
		IsValid:             true,
		ValidationErrors:    []string{},
	}

	var totalEasyQuestions, totalMediumQuestions, totalHardQuestions int

	for _, topic := range topics {
		// Somar pesos
		stats.TotalWeight += topic.WeightPercentage

		// Somar questões
		stats.TotalQuestions += topic.ActualQuestions

		// Contar tópicos completos
		if topic.IsComplete {
			stats.CompletedTopics++
		}

		// Criar distribuição de tópicos
		distribution := &domain.TopicDistribution{
			TopicID:          topic.TopicID,
			TopicName:        topic.TopicName,
			QuestionCount:    topic.ActualQuestions,
			WeightPercentage: topic.WeightPercentage,
			ExpectedCount:    topic.QuestionsCount,
			IsComplete:       topic.IsComplete,
		}
		stats.TopicDistribution = append(stats.TopicDistribution, distribution)

		// Somar questões por dificuldade
		totalEasyQuestions += topic.EasyQuestions
		totalMediumQuestions += topic.MediumQuestions
		totalHardQuestions += topic.HardQuestions
	}

	// Calcular breakdown de dificuldade
	if stats.TotalQuestions > 0 {
		stats.DifficultyBreakdown.EasyQuestions = totalEasyQuestions
		stats.DifficultyBreakdown.MediumQuestions = totalMediumQuestions
		stats.DifficultyBreakdown.HardQuestions = totalHardQuestions
		stats.DifficultyBreakdown.EasyPercentage = float64(totalEasyQuestions) / float64(stats.TotalQuestions) * 100
		stats.DifficultyBreakdown.MediumPercentage = float64(totalMediumQuestions) / float64(stats.TotalQuestions) * 100
		stats.DifficultyBreakdown.HardPercentage = float64(totalHardQuestions) / float64(stats.TotalQuestions) * 100
	}

	// Validar configuração
	isValid, validationErrors := uc.ValidateExamConfiguration(examID)
	stats.IsValid = isValid
	stats.ValidationErrors = validationErrors

	return stats, nil
}

// ValidateExamConfiguration valida se a configuração completa do exame está correta
func (uc *ExamTopicUsecase) ValidateExamConfiguration(examID int) (bool, []string) {
	var errors []string

	// Verificar se o exame existe
	_, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Exam not found: %v", err))
		return false, errors
	}

	// Verificar se há tópicos
	topicCount, err := uc.ExamTopicRepo.GetExamTopicsCount(examID)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to count exam topics: %v", err))
		return false, errors
	}

	if topicCount == 0 {
		errors = append(errors, "Exam has no topics configured")
		return false, errors
	}

	// Verificar distribuição de pesos
	err = uc.ExamTopicRepo.ValidateWeightDistribution(examID)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Weight distribution error: %v", err))
	}

	// Verificar configuração de cada tópico
	topics, err := uc.ExamTopicRepo.ListByExamWithDetails(examID)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to get exam topics: %v", err))
		return false, errors
	}

	for _, topic := range topics {
		// Validar distribuição de dificuldade
		examTopic := &domain.ExamTopic{
			DifficultyEasyPercentage:   topic.DifficultyEasyPercentage,
			DifficultyMediumPercentage: topic.DifficultyMediumPercentage,
			DifficultyHardPercentage:   topic.DifficultyHardPercentage,
		}

		err := uc.ValidateDifficultyDistribution(examTopic)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Topic '%s' difficulty distribution error: %v", topic.TopicName, err))
		}

		// Verificar se há questões suficientes
		if topic.ActualQuestions < topic.QuestionsCount {
			errors = append(errors, fmt.Sprintf("Topic '%s' has only %d questions but needs %d",
				topic.TopicName, topic.ActualQuestions, topic.QuestionsCount))
		}
	}

	return len(errors) == 0, errors
}

// AutoDistributeWeights distribui pesos automaticamente entre os tópicos de um exame
func (uc *ExamTopicUsecase) AutoDistributeWeights(examID int) error {
	// Validar se o exame existe
	_, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		return fmt.Errorf("exam not found: %w", err)
	}

	// Obter tópicos do exame
	topics, err := uc.ExamTopicRepo.ListByExam(examID)
	if err != nil {
		return fmt.Errorf("failed to get exam topics: %w", err)
	}

	if len(topics) == 0 {
		return fmt.Errorf("exam has no topics to distribute weights")
	}

	// Calcular peso igual para todos os tópicos
	weightPerTopic := 100.0 / float64(len(topics))
	remainder := 100.0 - (weightPerTopic * float64(len(topics)))

	// Atualizar pesos
	for i, topic := range topics {
		topic.WeightPercentage = weightPerTopic

		// Adicionar o restante ao primeiro tópico para garantir total de 100%
		if i == 0 {
			topic.WeightPercentage += remainder
		}
	}

	// Atualizar no banco
	err = uc.ExamTopicRepo.BulkUpdate(topics)
	if err != nil {
		return fmt.Errorf("failed to update topic weights: %w", err)
	}

	return nil
}

// SuggestOptimalDistribution sugere uma distribuição otimizada de pesos e dificuldade
func (uc *ExamTopicUsecase) SuggestOptimalDistribution(examID int) (*domain.ExamTopicStats, error) {
	// Implementação simplificada - em uma versão completa, isso poderia usar
	// algoritmos mais sofisticados baseados em dados históricos

	// Obter estatísticas atuais
	currentStats, err := uc.GetExamStats(examID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current stats: %w", err)
	}

	// Criar sugestões baseadas em boas práticas
	suggestions := &domain.ExamTopicStats{
		ExamID:            examID,
		TotalTopics:       currentStats.TotalTopics,
		TotalWeight:       100.0,
		TotalQuestions:    currentStats.TotalQuestions,
		CompletedTopics:   currentStats.CompletedTopics,
		TopicDistribution: []*domain.TopicDistribution{},
		DifficultyBreakdown: &domain.DifficultyBreakdown{
			// Distribuição padrão sugerida: 40% fácil, 40% médio, 20% difícil
			EasyPercentage:   40.0,
			MediumPercentage: 40.0,
			HardPercentage:   20.0,
		},
		IsValid:          true,
		ValidationErrors: []string{},
	}

	// Distribuir pesos igualmente entre os tópicos
	if currentStats.TotalTopics > 0 {
		weightPerTopic := 100.0 / float64(currentStats.TotalTopics)

		for _, topic := range currentStats.TopicDistribution {
			suggestedTopic := &domain.TopicDistribution{
				TopicID:          topic.TopicID,
				TopicName:        topic.TopicName,
				QuestionCount:    topic.QuestionCount,
				WeightPercentage: weightPerTopic,
				ExpectedCount:    topic.ExpectedCount,
				IsComplete:       topic.IsComplete,
			}
			suggestions.TopicDistribution = append(suggestions.TopicDistribution, suggestedTopic)
		}
	}

	return suggestions, nil
}

// ReorderTopics reordena os tópicos de um exame
func (uc *ExamTopicUsecase) ReorderTopics(examID int, topicOrders map[int]int) error {
	// Validar se o exame existe
	_, err := uc.ExamRepo.FindByID(examID)
	if err != nil {
		return fmt.Errorf("exam not found: %w", err)
	}

	// Validar se todos os tópicos pertencem ao exame
	for topicID := range topicOrders {
		_, err := uc.ExamTopicRepo.FindByExamAndTopic(examID, topicID)
		if err != nil {
			return fmt.Errorf("topic %d is not associated with exam %d: %w", topicID, examID, err)
		}
	}

	// Reordenar os tópicos
	err = uc.ExamTopicRepo.ReorderTopics(examID, topicOrders)
	if err != nil {
		return fmt.Errorf("failed to reorder topics: %w", err)
	}

	return nil
}

// BulkUpdateExamTopics atualiza múltiplas associações exame-tópico
func (uc *ExamTopicUsecase) BulkUpdateExamTopics(examTopics []*domain.ExamTopic) error {
	if len(examTopics) == 0 {
		return nil
	}

	// Validar cada associação
	for _, examTopic := range examTopics {
		// Verificar se a associação existe
		_, err := uc.ExamTopicRepo.FindByExamAndTopic(examTopic.ExamID, examTopic.TopicID)
		if err != nil {
			return fmt.Errorf("exam-topic association %d-%d not found: %w", examTopic.ExamID, examTopic.TopicID, err)
		}

		// Validar distribuição de dificuldade
		err = uc.ValidateDifficultyDistribution(examTopic)
		if err != nil {
			return fmt.Errorf("difficulty distribution validation failed for exam %d topic %d: %w", examTopic.ExamID, examTopic.TopicID, err)
		}

		// Validar peso percentual
		if examTopic.WeightPercentage <= 0 || examTopic.WeightPercentage > 100 {
			return fmt.Errorf("weight percentage must be between 0 and 100 for exam %d topic %d", examTopic.ExamID, examTopic.TopicID)
		}
	}

	// Validar distribuição total de pesos por exame
	examWeights := make(map[int]float64)
	for _, examTopic := range examTopics {
		examWeights[examTopic.ExamID] += examTopic.WeightPercentage
	}

	for examID, totalWeight := range examWeights {
		// Obter tópicos existentes para calcular peso total
		existingTopics, err := uc.ExamTopicRepo.ListByExam(examID)
		if err != nil {
			return fmt.Errorf("failed to get existing topics for exam %d: %w", examID, err)
		}

		// Calcular peso total considerando apenas os tópicos que não estão sendo atualizados
		currentTotalWeight := totalWeight
		updateMap := make(map[int]bool)
		for _, examTopic := range examTopics {
			if examTopic.ExamID == examID {
				updateMap[examTopic.TopicID] = true
			}
		}

		for _, existing := range existingTopics {
			if !updateMap[existing.TopicID] {
				currentTotalWeight += existing.WeightPercentage
			}
		}

		if currentTotalWeight != 100.0 {
			return fmt.Errorf("total weight percentage for exam %d would be %.2f, must be 100", examID, currentTotalWeight)
		}
	}

	// Atualizar as associações
	err := uc.ExamTopicRepo.BulkUpdate(examTopics)
	if err != nil {
		return fmt.Errorf("failed to bulk update exam-topic associations: %w", err)
	}

	return nil
}
