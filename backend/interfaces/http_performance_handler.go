package interfaces

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// PerformanceHandler lida com requisições HTTP relacionadas a relatórios de desempenho
type PerformanceHandler struct {
	UC *usecase.PerformanceUsecase
}

// ReportHandler gera um relatório de desempenho para um usuário
func (h *PerformanceHandler) ReportHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("PerformanceHandler: Iniciando relatório de desempenho")

	userIDStr := r.URL.Query().Get("user_id")
	log.Printf("PerformanceHandler: user_id recebido: %s", userIDStr)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		log.Printf("PerformanceHandler: ID de usuário inválido: %s", userIDStr)
		http.Error(w, "ID de usuário inválido", http.StatusBadRequest)
		return
	}

	log.Printf("PerformanceHandler: Chamando usecase para usuário %d", userID)
	report, err := h.UC.GetReport(userID)
	if err != nil {
		log.Printf("PerformanceHandler: Erro ao obter performance: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("PerformanceHandler: Relatório gerado com sucesso para usuário %d", userID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}
