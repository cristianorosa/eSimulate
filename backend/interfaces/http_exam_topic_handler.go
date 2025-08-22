package interfaces

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/domain"
	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// ExamTopicHandler gerencia as requisições HTTP para relacionamentos exame-tópico
type ExamTopicHandler struct {
	ExamTopicUC *usecase.ExamTopicUsecase
}

// NewExamTopicHandler cria uma nova instância do handler
func NewExamTopicHandler(examTopicUC *usecase.ExamTopicUsecase) *ExamTopicHandler {
	return &ExamTopicHandler{
		ExamTopicUC: examTopicUC,
	}
}

// AssociateTopicHandler associa um tópico a um exame
// POST /exams/{examId}/topics
func (h *ExamTopicHandler) AssociateTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair exam ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	examID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid exam ID", http.StatusBadRequest)
		return
	}

	// Decodificar body da requisição
	var requestBody struct {
		TopicID                    int     `json:"topic_id"`
		QuestionsCount             int     `json:"questions_count"`
		WeightPercentage           float64 `json:"weight_percentage"`
		OrderIndex                 int     `json:"order_index,omitempty"`
		DifficultyEasyPercentage   float64 `json:"difficulty_easy_percentage,omitempty"`
		DifficultyMediumPercentage float64 `json:"difficulty_medium_percentage,omitempty"`
		DifficultyHardPercentage   float64 `json:"difficulty_hard_percentage,omitempty"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validações básicas
	if requestBody.TopicID <= 0 {
		http.Error(w, "Topic ID is required", http.StatusBadRequest)
		return
	}

	if requestBody.WeightPercentage <= 0 {
		http.Error(w, "Weight percentage must be greater than 0", http.StatusBadRequest)
		return
	}

	if requestBody.QuestionsCount < 0 {
		http.Error(w, "Questions count cannot be negative", http.StatusBadRequest)
		return
	}

	// Definir percentuais padrão de dificuldade se não fornecidos
	if requestBody.DifficultyEasyPercentage == 0 && requestBody.DifficultyMediumPercentage == 0 && requestBody.DifficultyHardPercentage == 0 {
		requestBody.DifficultyEasyPercentage = 33.33
		requestBody.DifficultyMediumPercentage = 33.33
		requestBody.DifficultyHardPercentage = 33.34
	}

	// Criar associação
	examTopic := &domain.ExamTopic{
		ExamID:                     examID,
		TopicID:                    requestBody.TopicID,
		QuestionsCount:             requestBody.QuestionsCount,
		WeightPercentage:           requestBody.WeightPercentage,
		OrderIndex:                 requestBody.OrderIndex,
		DifficultyEasyPercentage:   requestBody.DifficultyEasyPercentage,
		DifficultyMediumPercentage: requestBody.DifficultyMediumPercentage,
		DifficultyHardPercentage:   requestBody.DifficultyHardPercentage,
	}

	err = h.ExamTopicUC.AssociateTopicWithExam(examTopic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Topic associated with exam successfully",
		"data":    examTopic,
	})
}

// DisassociateTopicHandler remove um tópico de um exame
// DELETE /exams/{examId}/topics/{topicId}
func (h *ExamTopicHandler) DisassociateTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair IDs da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	examID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid exam ID", http.StatusBadRequest)
		return
	}

	topicID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	// Desassociar tópico do exame
	err = h.ExamTopicUC.DisassociateTopicFromExam(examID, topicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Topic disassociated from exam successfully",
	})
}

// UpdateExamTopicHandler atualiza uma associação exame-tópico
// PUT /exams/{examId}/topics/{topicId}
func (h *ExamTopicHandler) UpdateExamTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair IDs da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	examID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid exam ID", http.StatusBadRequest)
		return
	}

	topicID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	// Decodificar body da requisição
	var requestBody struct {
		QuestionsCount             int     `json:"questions_count"`
		WeightPercentage           float64 `json:"weight_percentage"`
		OrderIndex                 int     `json:"order_index"`
		DifficultyEasyPercentage   float64 `json:"difficulty_easy_percentage"`
		DifficultyMediumPercentage float64 `json:"difficulty_medium_percentage"`
		DifficultyHardPercentage   float64 `json:"difficulty_hard_percentage"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validações básicas
	if requestBody.WeightPercentage <= 0 {
		http.Error(w, "Weight percentage must be greater than 0", http.StatusBadRequest)
		return
	}

	if requestBody.QuestionsCount < 0 {
		http.Error(w, "Questions count cannot be negative", http.StatusBadRequest)
		return
	}

	// Atualizar associação
	examTopic := &domain.ExamTopic{
		ExamID:                     examID,
		TopicID:                    topicID,
		QuestionsCount:             requestBody.QuestionsCount,
		WeightPercentage:           requestBody.WeightPercentage,
		OrderIndex:                 requestBody.OrderIndex,
		DifficultyEasyPercentage:   requestBody.DifficultyEasyPercentage,
		DifficultyMediumPercentage: requestBody.DifficultyMediumPercentage,
		DifficultyHardPercentage:   requestBody.DifficultyHardPercentage,
	}

	err = h.ExamTopicUC.UpdateExamTopic(examTopic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Exam-topic association updated successfully",
		"data":    examTopic,
	})
}

// GetExamTopicsHandler retorna todos os tópicos de um exame
// GET /exams/{examId}/topics
func (h *ExamTopicHandler) GetExamTopicsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair exam ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	examID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid exam ID", http.StatusBadRequest)
		return
	}

	// Obter tópicos do exame
	topics, err := h.ExamTopicUC.GetExamTopics(examID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  topics,
		"count": len(topics),
	})
}

// GetTopicExamsHandler retorna todos os exames que contêm um tópico
// GET /topics/{topicId}/exams
func (h *ExamTopicHandler) GetTopicExamsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair topic ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	topicID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	// Obter exames do tópico
	exams, err := h.ExamTopicUC.GetTopicExams(topicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  exams,
		"count": len(exams),
	})
}

// CalculateQuestionsByDifficultyHandler calcula distribuição de questões por dificuldade
// POST /exams/{examId}/topics/{topicId}/calculate-difficulty
func (h *ExamTopicHandler) CalculateQuestionsByDifficultyHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair IDs da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	examID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid exam ID", http.StatusBadRequest)
		return
	}

	topicID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	// Decodificar body da requisição
	var requestBody struct {
		QuestionsCount             int     `json:"questions_count"`
		DifficultyEasyPercentage   float64 `json:"difficulty_easy_percentage"`
		DifficultyMediumPercentage float64 `json:"difficulty_medium_percentage"`
		DifficultyHardPercentage   float64 `json:"difficulty_hard_percentage"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Criar estrutura temporária para cálculo
	examTopic := &domain.ExamTopic{
		ExamID:                     examID,
		TopicID:                    topicID,
		QuestionsCount:             requestBody.QuestionsCount,
		DifficultyEasyPercentage:   requestBody.DifficultyEasyPercentage,
		DifficultyMediumPercentage: requestBody.DifficultyMediumPercentage,
		DifficultyHardPercentage:   requestBody.DifficultyHardPercentage,
	}

	// Calcular distribuição
	distribution, err := h.ExamTopicUC.CalculateQuestionsByDifficulty(examTopic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": distribution,
	})
}

// GetExamStatsHandler retorna estatísticas completas de um exame
// GET /exams/{examId}/stats
func (h *ExamTopicHandler) GetExamStatsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair exam ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	examID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid exam ID", http.StatusBadRequest)
		return
	}

	// Obter estatísticas
	stats, err := h.ExamTopicUC.GetExamStats(examID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": stats,
	})
}

// ValidateExamConfigurationHandler valida a configuração completa de um exame
// GET /exams/{examId}/topics/validate
func (h *ExamTopicHandler) ValidateExamConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair exam ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	examID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid exam ID", http.StatusBadRequest)
		return
	}

	// Validar configuração
	isValid, errors := h.ExamTopicUC.ValidateExamConfiguration(examID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"is_valid": isValid,
		"errors":   errors,
	})
}

// AutoDistributeWeightsHandler distribui pesos automaticamente entre os tópicos
// POST /exams/{examId}/topics/auto-distribute-weights
func (h *ExamTopicHandler) AutoDistributeWeightsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair exam ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	examID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid exam ID", http.StatusBadRequest)
		return
	}

	// Auto-distribuir pesos
	err = h.ExamTopicUC.AutoDistributeWeights(examID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Weights distributed automatically",
	})
}

// SuggestOptimalDistributionHandler sugere uma distribuição otimizada
// GET /exams/{examId}/topics/suggest-distribution
func (h *ExamTopicHandler) SuggestOptimalDistributionHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair exam ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	examID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid exam ID", http.StatusBadRequest)
		return
	}

	// Obter sugestões
	suggestions, err := h.ExamTopicUC.SuggestOptimalDistribution(examID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": suggestions,
	})
}

// ReorderTopicsHandler reordena os tópicos de um exame
// PUT /exams/{examId}/topics/reorder
func (h *ExamTopicHandler) ReorderTopicsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair exam ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	examID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid exam ID", http.StatusBadRequest)
		return
	}

	// Decodificar body da requisição
	var requestBody struct {
		TopicOrders map[int]int `json:"topic_orders"` // topicID -> orderIndex
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if len(requestBody.TopicOrders) == 0 {
		http.Error(w, "Topic orders are required", http.StatusBadRequest)
		return
	}

	// Reordenar tópicos
	err = h.ExamTopicUC.ReorderTopics(examID, requestBody.TopicOrders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Topics reordered successfully",
	})
}

// BulkUpdateExamTopicsHandler atualiza múltiplas associações exame-tópico
// PUT /exams/{examId}/topics/bulk
func (h *ExamTopicHandler) BulkUpdateExamTopicsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair exam ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	examID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid exam ID", http.StatusBadRequest)
		return
	}

	// Decodificar body da requisição
	var requestBody struct {
		Topics []struct {
			TopicID                    int     `json:"topic_id"`
			QuestionsCount             int     `json:"questions_count"`
			WeightPercentage           float64 `json:"weight_percentage"`
			OrderIndex                 int     `json:"order_index"`
			DifficultyEasyPercentage   float64 `json:"difficulty_easy_percentage"`
			DifficultyMediumPercentage float64 `json:"difficulty_medium_percentage"`
			DifficultyHardPercentage   float64 `json:"difficulty_hard_percentage"`
		} `json:"topics"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if len(requestBody.Topics) == 0 {
		http.Error(w, "Topics are required", http.StatusBadRequest)
		return
	}

	// Converter para estruturas de domínio
	var examTopics []*domain.ExamTopic
	for _, topic := range requestBody.Topics {
		examTopic := &domain.ExamTopic{
			ExamID:                     examID,
			TopicID:                    topic.TopicID,
			QuestionsCount:             topic.QuestionsCount,
			WeightPercentage:           topic.WeightPercentage,
			OrderIndex:                 topic.OrderIndex,
			DifficultyEasyPercentage:   topic.DifficultyEasyPercentage,
			DifficultyMediumPercentage: topic.DifficultyMediumPercentage,
			DifficultyHardPercentage:   topic.DifficultyHardPercentage,
		}
		examTopics = append(examTopics, examTopic)
	}

	// Atualizar associações
	err = h.ExamTopicUC.BulkUpdateExamTopics(examTopics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Exam-topic associations updated successfully",
		"count":   len(examTopics),
	})
}
