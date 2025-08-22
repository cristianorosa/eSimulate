package interfaces

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/domain"
	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// QuestionTagHandler gerencia as requisições HTTP para tags de questões
type QuestionTagHandler struct {
	QuestionTagUC *usecase.QuestionTagUsecase
}

// NewQuestionTagHandler cria uma nova instância do handler
func NewQuestionTagHandler(questionTagUC *usecase.QuestionTagUsecase) *QuestionTagHandler {
	return &QuestionTagHandler{
		QuestionTagUC: questionTagUC,
	}
}

// CreateTagHandler cria uma nova tag
// POST /tags
func (h *QuestionTagHandler) CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if requestBody.Name == "" {
		http.Error(w, "Tag name is required", http.StatusBadRequest)
		return
	}

	// Criar tag
	tag, err := h.QuestionTagUC.CreateTag(requestBody.Name, requestBody.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tag created successfully",
		"data":    tag,
	})
}

// UpdateTagHandler atualiza uma tag existente
// PUT /tags/{id}
func (h *QuestionTagHandler) UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if requestBody.Name == "" {
		http.Error(w, "Tag name is required", http.StatusBadRequest)
		return
	}

	// Atualizar tag
	tag, err := h.QuestionTagUC.UpdateTag(tagID, requestBody.Name, requestBody.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tag updated successfully",
		"data":    tag,
	})
}

// DeleteTagHandler remove uma tag
// DELETE /tags/{id}
func (h *QuestionTagHandler) DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	// Deletar tag
	err = h.QuestionTagUC.DeleteTag(tagID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Tag deleted successfully",
	})
}

// GetTagHandler retorna uma tag por ID
// GET /tags/{id}
func (h *QuestionTagHandler) GetTagHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	// Obter tag
	tag, err := h.QuestionTagUC.GetTag(tagID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": tag,
	})
}

// ListTagsHandler retorna todas as tags
// GET /tags?with_stats=true
func (h *QuestionTagHandler) ListTagsHandler(w http.ResponseWriter, r *http.Request) {
	withStats := r.URL.Query().Get("with_stats") == "true"

	if withStats {
		// Retornar tags com estatísticas
		tags, err := h.QuestionTagUC.ListTagsWithStats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":  tags,
			"count": len(tags),
		})
	} else {
		// Retornar tags simples
		tags, err := h.QuestionTagUC.ListAllTags()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":  tags,
			"count": len(tags),
		})
	}
}

// AssociateQuestionTagHandler associa uma questão com uma tag
// POST /questions/{questionId}/tags/{tagId}
func (h *QuestionTagHandler) AssociateQuestionTagHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair IDs da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	questionID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	// Associar questão com tag
	err = h.QuestionTagUC.AssociateQuestionTag(questionID, tagID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Question associated with tag successfully",
	})
}

// DisassociateQuestionTagHandler remove a associação entre uma questão e uma tag
// DELETE /questions/{questionId}/tags/{tagId}
func (h *QuestionTagHandler) DisassociateQuestionTagHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair IDs da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	questionID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	// Desassociar questão da tag
	err = h.QuestionTagUC.DisassociateQuestionTag(questionID, tagID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Question disassociated from tag successfully",
	})
}

// GetQuestionTagsHandler retorna todas as tags de uma questão
// GET /questions/{questionId}/tags
func (h *QuestionTagHandler) GetQuestionTagsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Obter tags da questão
	tags, err := h.QuestionTagUC.GetQuestionTags(questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  tags,
		"count": len(tags),
	})
}

// GetTagQuestionsHandler retorna todas as questões de uma tag
// GET /tags/{tagId}/questions
func (h *QuestionTagHandler) GetTagQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair tag ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	// Obter questões da tag
	questions, err := h.QuestionTagUC.GetTagQuestions(tagID)
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

// UpdateQuestionTagsHandler atualiza todas as tags de uma questão
// PUT /questions/{questionId}/tags
func (h *QuestionTagHandler) UpdateQuestionTagsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Decodificar body da requisição
	var requestBody struct {
		TagIDs []int `json:"tag_ids"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Atualizar tags da questão
	err = h.QuestionTagUC.UpdateQuestionTags(questionID, requestBody.TagIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Question tags updated successfully",
		"count":   len(requestBody.TagIDs),
	})
}

// BulkAssociateQuestionTagsHandler associa uma questão com múltiplas tags
// POST /questions/{questionId}/tags/bulk
func (h *QuestionTagHandler) BulkAssociateQuestionTagsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Decodificar body da requisição
	var requestBody struct {
		TagIDs []int `json:"tag_ids"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if len(requestBody.TagIDs) == 0 {
		http.Error(w, "Tag IDs are required", http.StatusBadRequest)
		return
	}

	// Associar questão com tags
	err = h.QuestionTagUC.BulkAssociateQuestionTags(questionID, requestBody.TagIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Question associated with tags successfully",
		"count":   len(requestBody.TagIDs),
	})
}

// SearchTagsHandler busca tags por nome
// GET /tags/search?q={query}&limit={limit}
func (h *QuestionTagHandler) SearchTagsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20 // padrão
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 20
		}
	}

	// Buscar tags
	tags, err := h.QuestionTagUC.SearchTags(query, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  tags,
		"count": len(tags),
		"query": query,
	})
}

// FilterQuestionsByTagsHandler filtra questões baseado em critérios de tags
// POST /questions/filter-by-tags
func (h *QuestionTagHandler) FilterQuestionsByTagsHandler(w http.ResponseWriter, r *http.Request) {
	var filters domain.QuestionTagFilters

	err := json.NewDecoder(r.Body).Decode(&filters)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Filtrar questões
	questions, err := h.QuestionTagUC.FilterQuestionsByTags(&filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    questions,
		"count":   len(questions),
		"filters": filters,
	})
}

// GetTagUsageStatsHandler retorna estatísticas de uso de uma tag
// GET /tags/{tagId}/stats
func (h *QuestionTagHandler) GetTagUsageStatsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair tag ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	// Obter estatísticas
	stats, err := h.QuestionTagUC.GetTagUsageStats(tagID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": stats,
	})
}

// GetQuestionTagStatsHandler retorna estatísticas de tags de uma questão
// GET /questions/{questionId}/tags/stats
func (h *QuestionTagHandler) GetQuestionTagStatsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Obter estatísticas
	stats, err := h.QuestionTagUC.GetQuestionTagStats(questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": stats,
	})
}

// GetMostUsedTagsHandler retorna as tags mais utilizadas
// GET /tags/most-used?limit={limit}
func (h *QuestionTagHandler) GetMostUsedTagsHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // padrão
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10
		}
	}

	// Obter tags mais usadas
	tags, err := h.QuestionTagUC.GetMostUsedTags(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  tags,
		"count": len(tags),
	})
}

// SuggestTagsHandler sugere tags para uma questão
// GET /questions/{questionId}/tags/suggestions
func (h *QuestionTagHandler) SuggestTagsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Obter sugestões
	suggestions, err := h.QuestionTagUC.SuggestTagsForQuestion(questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  suggestions,
		"count": len(suggestions),
	})
}

// FindSimilarQuestionsHandler encontra questões similares baseado nas tags
// GET /questions/{questionId}/similar?limit={limit}
func (h *QuestionTagHandler) FindSimilarQuestionsHandler(w http.ResponseWriter, r *http.Request) {
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

	limitStr := r.URL.Query().Get("limit")
	limit := 10 // padrão
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10
		}
	}

	// Encontrar questões similares
	questions, err := h.QuestionTagUC.FindSimilarQuestionsByTags(questionID, limit)
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

// CleanupUnusedTagsHandler remove tags não utilizadas
// DELETE /tags/cleanup
func (h *QuestionTagHandler) CleanupUnusedTagsHandler(w http.ResponseWriter, r *http.Request) {
	// Limpar tags não utilizadas
	deletedCount, err := h.QuestionTagUC.CleanupUnusedTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Unused tags cleaned up successfully",
		"deleted_count": deletedCount,
	})
}
