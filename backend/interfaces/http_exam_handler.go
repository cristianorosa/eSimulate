package interfaces

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// ExamHandler lida com requisições HTTP relacionadas a exames
type ExamHandler struct {
	UC *usecase.ExamUsecase
}

// CreateHandler cria um novo exame
func (h *ExamHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title          string  `json:"title"`
		Description    string  `json:"description"`
		AreaID         int     `json:"area_id"`
		MaxTimeMinutes int     `json:"max_time_minutes"`
		PassingScore   float64 `json:"passing_score"`
		IsActive       bool    `json:"is_active"`
		CreatedBy      int     `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	exam, err := h.UC.CreateExam(context.Background(), req.Title, req.Description, req.AreaID, req.MaxTimeMinutes, req.PassingScore, req.IsActive, req.CreatedBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exam)
}

// UpdateHandler atualiza um exame existente
func (h *ExamHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Title          string  `json:"title"`
		Description    string  `json:"description"`
		AreaID         int     `json:"area_id"`
		MaxTimeMinutes int     `json:"max_time_minutes"`
		PassingScore   float64 `json:"passing_score"`
		IsActive       bool    `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	err = h.UC.UpdateExam(context.Background(), id, req.Title, req.Description, req.AreaID, req.MaxTimeMinutes, req.PassingScore, req.IsActive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteHandler remove um exame
func (h *ExamHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.UC.DeleteExam(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListHandler lista todos os exames (mantido para compatibilidade)
func (h *ExamHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	areaIDStr := r.URL.Query().Get("area_id")
	
	var exams []interface{}

	if areaIDStr != "" {
		areaID, err := strconv.Atoi(areaIDStr)
		if err != nil {
			http.Error(w, "ID da área inválido", http.StatusBadRequest)
			return
		}
		
		examList, err := h.UC.ListExamsByArea(context.Background(), areaID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		exams = make([]interface{}, len(examList))
		for i, exam := range examList {
			exams[i] = exam
		}
	} else {
		examList, err := h.UC.ListExams(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		exams = make([]interface{}, len(examList))
		for i, exam := range examList {
			exams[i] = exam
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exams)
}

// ListPaginatedHandler lista exames com paginação
func (h *ExamHandler) ListPaginatedHandler(w http.ResponseWriter, r *http.Request) {
	// Obter parâmetros de paginação da query string
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	areaIDStr := r.URL.Query().Get("area_id")

	page := 1
	pageSize := 10
	var areaID *int

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

	if areaIDStr != "" {
		if aid, err := strconv.Atoi(areaIDStr); err == nil && aid > 0 {
			areaID = &aid
		}
	}

	exams, pagination, err := h.UC.ListExamsPaginated(context.Background(), page, pageSize, areaID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Data       []interface{} `json:"data"`
		Pagination interface{}   `json:"pagination"`
	}{
		Data:       make([]interface{}, len(exams)),
		Pagination: pagination,
	}

	for i, exam := range exams {
		response.Data[i] = exam
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetHandler retorna um exame específico
func (h *ExamHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	exam, err := h.UC.GetExam(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exam)
}
