package interfaces

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"eSimulate/backend/usecase"
)

// PerformanceHandler lida com relatórios de desempenho

type PerformanceHandler struct{
	UC *usecase.PerformanceUsecase
}

func (h *PerformanceHandler) ReportHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "ID de usuário inválido", http.StatusBadRequest)
		return
	}
	report, err := h.UC.GetUserPerformance(context.Background(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}
