package interfaces

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cristianorosa/eSimulate/backend/domain"
	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// DomainHandler lida com requisições HTTP relacionadas a domínios
type DomainHandler struct {
	UC *usecase.DomainUsecase
}

// ListHandler lista todos os domínios de um exame
func (h *DomainHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	examIDStr := r.URL.Query().Get("exam_id")
	
	// Se exam_id não foi fornecido, retornar todos os domínios
	if examIDStr == "" {
		domains, err := h.UC.ListAllDomains(context.Background())
		if err != nil {
			http.Error(w, "Erro ao buscar domínios: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Garantir que sempre retorne um array, mesmo que vazio
		if domains == nil {
			domains = []*domain.Domain{}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(domains)
		return
	}

	examID, err := strconv.Atoi(examIDStr)
	if err != nil {
		http.Error(w, "exam_id inválido", http.StatusBadRequest)
		return
	}

	domains, err := h.UC.ListDomainsByExam(context.Background(), examID)
	if err != nil {
		http.Error(w, "Erro ao buscar domínios: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Garantir que sempre retorne um array, mesmo que vazio
	if domains == nil {
		domains = []*domain.Domain{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domains)
}

// GetHandler trata o endpoint para obter um domínio específico
func (h *DomainHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	domain, err := h.UC.GetDomain(context.Background(), id)
	if err != nil {
		http.Error(w, "Erro ao buscar domínio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain)
}

// CreateHandler trata o endpoint para criar um novo domínio
func (h *DomainHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ExamID           int     `json:"exam_id"`
		Name             string  `json:"name"`
		Description      string  `json:"description"`
		WeightPercentage float64 `json:"weight_percentage"`
		OrderIndex       int     `json:"order_index"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	domain, err := h.UC.CreateDomain(context.Background(), req.ExamID, req.Name, req.Description, req.WeightPercentage, req.OrderIndex)
	if err != nil {
		http.Error(w, "Erro ao criar domínio: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(domain)
}

// UpdateHandler trata o endpoint para atualizar um domínio
func (h *DomainHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Name             string  `json:"name"`
		Description      string  `json:"description"`
		WeightPercentage float64 `json:"weight_percentage"`
		OrderIndex       int     `json:"order_index"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	err = h.UC.UpdateDomain(context.Background(), id, req.Name, req.Description, req.WeightPercentage, req.OrderIndex)
	if err != nil {
		http.Error(w, "Erro ao atualizar domínio: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteHandler trata o endpoint para excluir um domínio
func (h *DomainHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.UC.DeleteDomain(context.Background(), id)
	if err != nil {
		http.Error(w, "Erro ao excluir domínio: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
