package interfaces

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// HistoryHandler lida com histórico de simulados realizados

// HistoryHandler lida com requisições HTTP relacionadas ao histórico de simulados.
type HistoryHandler struct{
	UC *usecase.UserQuizUsecase
}

// ListHandler lista o histórico de simulados de um usuário.
func (h *HistoryHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "ID de usuário inválido", http.StatusBadRequest)
		return
	}
	history, err := h.UC.ListUserHistory(context.Background(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}
