package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// QuestionHandler lida com requisições HTTP relacionadas a questões.
type QuestionHandler struct {
	UC *usecase.QuestionUsecase
}

// ListHandler lista todas as questões, opcionalmente filtradas por exame.
func (h *QuestionHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	// Retorna uma resposta simples para testar
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]interface{}{})
}

// CreateHandler cria uma nova questão.
func (h *QuestionHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Create handler"})
}

// UpdateHandler atualiza uma questão existente.
func (h *QuestionHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Update handler"})
}

// DeleteHandler remove uma questão.
func (h *QuestionHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Delete handler"})
}

// DetailHandler retorna os detalhes de uma questão específica.
func (h *QuestionHandler) DetailHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Detail handler"})
}
