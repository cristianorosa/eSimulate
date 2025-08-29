package interfaces

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// QuestionTagHandler lida com requisições HTTP relacionadas a tags de questões
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
		Name string `json:"name"`
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
	tag, err := h.QuestionTagUC.CreateTag(requestBody.Name)
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

// FindOrCreateTagHandler busca uma tag por nome ou cria uma nova se não existir
// POST /tags/find-or-create
func (h *QuestionTagHandler) FindOrCreateTagHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name string `json:"name"`
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

	// Buscar ou criar tag
	tag, err := h.QuestionTagUC.FindOrCreateTag(requestBody.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tag found or created successfully",
		"data":    tag,
	})
}

// SearchTagsHandler busca tags por padrão de nome (para autocomplete)
// GET /tags/search?q=pattern&limit=20
func (h *QuestionTagHandler) SearchTagsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20 // Default limit

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	// Buscar tags
	tags, err := h.QuestionTagUC.SearchTagsByName(query, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retornar apenas os nomes para otimizar o payload
	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tags found successfully",
		"data":    tagNames,
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

	tagID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Name string `json:"name"`
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
	tag, err := h.QuestionTagUC.UpdateTag(tagID, requestBody.Name)
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

// DeleteTagHandler deleta uma tag
// DELETE /tags/{id}
func (h *QuestionTagHandler) DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(pathParts[len(pathParts)-1])
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
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tag deleted successfully",
	})
}

// GetTagHandler busca uma tag por ID
// GET /tags/{id}
func (h *QuestionTagHandler) GetTagHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	tagID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	// Buscar tag
	tag, err := h.QuestionTagUC.GetTagByID(tagID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tag retrieved successfully",
		"data":    tag,
	})
}

// ListTagsHandler lista todas as tags
// GET /tags
func (h *QuestionTagHandler) ListTagsHandler(w http.ResponseWriter, r *http.Request) {
	// Parâmetros de paginação
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// Listar tags
	tags, err := h.QuestionTagUC.ListTags(page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tags retrieved successfully",
		"data":    tags,
	})
}

// ListTagsWithStatsHandler lista tags com estatísticas
// GET /tags/stats
func (h *QuestionTagHandler) ListTagsWithStatsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := h.QuestionTagUC.ListTagsWithStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tags with stats retrieved successfully",
		"data":    tags,
	})
}



// AssociateQuestionTagHandler associa uma tag a uma questão
// POST /questions/{question_id}/tags/{tag_id}
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

	// Associar
	err = h.QuestionTagUC.AssociateQuestionTag(questionID, tagID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Question tag associated successfully",
	})
}

// DisassociateQuestionTagHandler remove associação entre questão e tag
// DELETE /questions/{question_id}/tags/{tag_id}
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

	// Desassociar
	err = h.QuestionTagUC.DisassociateQuestionTag(questionID, tagID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Question tag disassociated successfully",
	})
}

// GetQuestionTagsHandler lista todas as tags de uma questão
// GET /questions/{question_id}/tags
func (h *QuestionTagHandler) GetQuestionTagsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
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

	// Buscar tags da questão
	tags, err := h.QuestionTagUC.GetQuestionTags(questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Question tags retrieved successfully",
		"data":    tags,
	})
}

// UpdateQuestionTagsHandler atualiza todas as tags de uma questão
// PUT /questions/{question_id}/tags
func (h *QuestionTagHandler) UpdateQuestionTagsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
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
	})
}

// BulkAssociateQuestionTagsHandler associa múltiplas tags a uma questão
// POST /questions/{question_id}/tags/bulk
func (h *QuestionTagHandler) BulkAssociateQuestionTagsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
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

	// Associar tags em lote
	err = h.QuestionTagUC.BulkAssociateQuestionTags(questionID, requestBody.TagIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Question tags bulk associated successfully",
	})
}

// GetTagStatsHandler retorna estatísticas de uma tag
// GET /tags/{id}/stats
func (h *QuestionTagHandler) GetTagStatsHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
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

	// Buscar estatísticas
	stats, err := h.QuestionTagUC.GetTagStats(tagID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tag stats retrieved successfully",
		"data":    stats,
	})
}
