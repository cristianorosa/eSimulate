package interfaces

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// ExamQuestionHandler gerencia as requisições HTTP para relacionamentos exame-questão
type ExamQuestionHandler struct {
	ExamQuestionUC *usecase.ExamQuestionUsecase
}

// NewExamQuestionHandler cria uma nova instância do handler
func NewExamQuestionHandler(examQuestionUC *usecase.ExamQuestionUsecase) *ExamQuestionHandler {
	return &ExamQuestionHandler{
		ExamQuestionUC: examQuestionUC,
	}
}

// AssociateQuestionHandler associa uma questão a um exame
// POST /exams/{examId}/questions/{questionId}
func (h *ExamQuestionHandler) AssociateQuestionHandler(w http.ResponseWriter, r *http.Request) {
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

	questionID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	// Associar questão ao exame
	err = h.ExamQuestionUC.AssociateQuestionWithExam(examID, questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Question associated with exam successfully",
	})
}

// DisassociateQuestionHandler remove uma questão de um exame
// DELETE /exams/{examId}/questions/{questionId}
func (h *ExamQuestionHandler) DisassociateQuestionHandler(w http.ResponseWriter, r *http.Request) {
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

	questionID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	// Desassociar questão do exame
	err = h.ExamQuestionUC.DisassociateQuestionFromExam(examID, questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Question disassociated from exam successfully",
	})
}

// GetExamQuestionsHandler retorna todas as questões de um exame
// GET /exams/{examId}/questions
func (h *ExamQuestionHandler) GetExamQuestionsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Obter questões do exame
	questions, err := h.ExamQuestionUC.GetExamQuestions(examID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  questions,
		"count": len(questions),
	})
}

// GetQuestionExamsHandler retorna todos os exames que contêm uma questão
// GET /questions/{questionId}/exams
func (h *ExamQuestionHandler) GetQuestionExamsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair question ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	questionID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	// Obter exames da questão
	exams, err := h.ExamQuestionUC.GetQuestionExams(questionID)
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

// ReorderExamQuestionsHandler reordena as questões de um exame
// PUT /exams/{examId}/questions/reorder
func (h *ExamQuestionHandler) ReorderExamQuestionsHandler(w http.ResponseWriter, r *http.Request) {
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
		QuestionIDs []int `json:"question_ids"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Reordenar questões
	err = h.ExamQuestionUC.ReorderExamQuestions(examID, requestBody.QuestionIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Questions reordered successfully",
	})
}

// BulkAssociateQuestionsHandler associa múltiplas questões a um exame
// POST /exams/{examId}/questions/bulk
func (h *ExamQuestionHandler) BulkAssociateQuestionsHandler(w http.ResponseWriter, r *http.Request) {
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
		QuestionIDs []int `json:"question_ids"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if len(requestBody.QuestionIDs) == 0 {
		http.Error(w, "Question IDs are required", http.StatusBadRequest)
		return
	}

	// Associar questões ao exame
	err = h.ExamQuestionUC.BulkAssociateQuestionsWithExam(examID, requestBody.QuestionIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Questions associated with exam successfully",
		"count":   len(requestBody.QuestionIDs),
	})
}

// BulkDisassociateQuestionsHandler remove múltiplas questões de um exame
// DELETE /exams/{examId}/questions/bulk
func (h *ExamQuestionHandler) BulkDisassociateQuestionsHandler(w http.ResponseWriter, r *http.Request) {
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
		QuestionIDs []int `json:"question_ids"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if len(requestBody.QuestionIDs) == 0 {
		http.Error(w, "Question IDs are required", http.StatusBadRequest)
		return
	}

	// Desassociar questões do exame
	err = h.ExamQuestionUC.BulkDisassociateQuestionsFromExam(examID, requestBody.QuestionIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Questions disassociated from exam successfully",
		"count":   len(requestBody.QuestionIDs),
	})
}

// GetExamQuestionStatsHandler retorna estatísticas das questões de um exame
// GET /exams/{examId}/questions/stats
func (h *ExamQuestionHandler) GetExamQuestionStatsHandler(w http.ResponseWriter, r *http.Request) {
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
	stats, err := h.ExamQuestionUC.GetExamQuestionStats(examID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": stats,
	})
}

// GetAvailableQuestionsHandler retorna questões disponíveis para associar a um exame
// GET /exams/{examId}/questions/available?topic_id={topicId}
func (h *ExamQuestionHandler) GetAvailableQuestionsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Extrair topic_id opcional dos query parameters
	var topicID *int
	if topicIDStr := r.URL.Query().Get("topic_id"); topicIDStr != "" {
		tid, err := strconv.Atoi(topicIDStr)
		if err != nil {
			http.Error(w, "Invalid topic ID", http.StatusBadRequest)
			return
		}
		topicID = &tid
	}

	// Obter questões disponíveis
	questions, err := h.ExamQuestionUC.GetAvailableQuestionsForExam(examID, topicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  questions,
		"count": len(questions),
	})
}

// ValidateExamConfigurationHandler valida a configuração de um exame
// GET /exams/{examId}/validate
func (h *ExamQuestionHandler) ValidateExamConfigurationHandler(w http.ResponseWriter, r *http.Request) {
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
	isValid, errors := h.ExamQuestionUC.ValidateExamConfiguration(examID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"is_valid": isValid,
		"errors":   errors,
	})
}
