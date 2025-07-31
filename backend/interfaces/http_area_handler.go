package interfaces

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// AreaHandler lida com requisições HTTP relacionadas a áreas
type AreaHandler struct {
	UC *usecase.AreaUsecase
}

// CreateHandler cria uma nova área
func (h *AreaHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	area, err := h.UC.CreateArea(context.Background(), req.Name, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(area)
}

// UpdateHandler atualiza uma área existente
func (h *AreaHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	err = h.UC.UpdateArea(context.Background(), id, req.Name, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteHandler remove uma área
func (h *AreaHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.UC.DeleteArea(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListHandler lista todas as áreas (mantido para compatibilidade)
func (h *AreaHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	areas, err := h.UC.ListAreas(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(areas)
}

// ListPaginatedHandler lista áreas com paginação
func (h *AreaHandler) ListPaginatedHandler(w http.ResponseWriter, r *http.Request) {
	// Obter parâmetros de paginação da query string
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

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

	areas, pagination, err := h.UC.ListAreasPaginated(context.Background(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Data       []interface{} `json:"data"`
		Pagination interface{}   `json:"pagination"`
	}{
		Data:       make([]interface{}, len(areas)),
		Pagination: pagination,
	}

	for i, area := range areas {
		response.Data[i] = area
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetHandler retorna uma área específica
func (h *AreaHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	area, err := h.UC.GetArea(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(area)
}
