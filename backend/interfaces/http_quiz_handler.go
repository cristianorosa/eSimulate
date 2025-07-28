package interfaces

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// QuizHandler lida com requisições HTTP relacionadas a simulados

// QuizHandler lida com requisições HTTP relacionadas a simulados.
type QuizHandler struct{
	UC *usecase.QuizUsecase
}

// ListHandler lista todos os simulados, opcionalmente filtrados por tema.
func (h *QuizHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	var themeID *int
	if tid := r.URL.Query().Get("theme_id"); tid != "" {
		id, err := strconv.Atoi(tid)
		if err == nil {
			themeID = &id
		}
	}
	quizzes, err := h.UC.ListQuizzes(context.Background(), themeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quizzes)
}
// CreateHandler cria um novo simulado.
func (h *QuizHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ThemeID     int    `json:"theme_id"`
		CreatedBy   int    `json:"created_by"`
		QuestionIDs []int  `json:"question_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	quiz, err := h.UC.CreateQuiz(context.Background(), req.Title, req.Description, req.ThemeID, req.CreatedBy, req.QuestionIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(quiz)
}
// DetailHandler retorna os detalhes de um simulado específico.
func (h *QuizHandler) DetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	quiz, err := h.UC.GetQuiz(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quiz)
}
// StartHandler inicia um simulado.
func (h *QuizHandler) StartHandler(w http.ResponseWriter, r *http.Request) {
	// Implementação: iniciar simulado
	w.WriteHeader(http.StatusOK)
}
// AnswerHandler registra a resposta de uma questão do simulado.
func (h *QuizHandler) AnswerHandler(w http.ResponseWriter, r *http.Request) {
	// Implementação: responder questão do simulado
	w.WriteHeader(http.StatusOK)
}
// ResultHandler retorna o resultado de um simulado.
func (h *QuizHandler) ResultHandler(w http.ResponseWriter, r *http.Request) {
	// Implementação: resultado do simulado
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
