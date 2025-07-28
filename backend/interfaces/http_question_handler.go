package interfaces

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cristianorosa/eSimulate/backend/domain"
	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// QuestionHandler lida com requisições HTTP relacionadas a questões

// QuestionHandler lida com requisições HTTP relacionadas a questões.
type QuestionHandler struct{
	UC *usecase.QuestionUsecase
}

// ListHandler lista todas as questões, opcionalmente filtradas por tema.
func (h *QuestionHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	var themeID *int
	if tid := r.URL.Query().Get("theme_id"); tid != "" {
		id, err := strconv.Atoi(tid)
		if err == nil {
			themeID = &id
		}
	}
	questions, err := h.UC.ListQuestions(context.Background(), themeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(questions)
}
// CreateHandler cria uma nova questão.
func (h *QuestionHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ThemeID     int                `json:"theme_id"`
		Statement   string             `json:"statement"`
		Explanation string             `json:"explanation"`
		CreatedBy   int                `json:"created_by"`
		Options     []domain.Option    `json:"options"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	// Converte slice de valor para slice de ponteiro
	var opts []*domain.Option
	for i := range req.Options {
		opts = append(opts, &req.Options[i])
	}
	q, err := h.UC.CreateQuestion(context.Background(), req.ThemeID, req.Statement, req.Explanation, req.CreatedBy, opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}
// UpdateHandler atualiza uma questão existente.
func (h *QuestionHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	var req struct {
		ThemeID     int    `json:"theme_id"`
		Statement   string `json:"statement"`
		Explanation string `json:"explanation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	err = h.UC.UpdateQuestion(context.Background(), id, req.ThemeID, req.Statement, req.Explanation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
// DeleteHandler remove uma questão.
func (h *QuestionHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	err = h.UC.DeleteQuestion(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
// DetailHandler retorna os detalhes de uma questão específica.
func (h *QuestionHandler) DetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	q, err := h.UC.GetQuestion(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(q)
}
