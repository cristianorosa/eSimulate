package interfaces

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cristianorosa/eSimulate/backend/domain"
	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// TopicHandler lida com requisições HTTP relacionadas a tópicos
type TopicHandler struct {
	UC *usecase.TopicUsecase
}

// ListHandler lista todos os tópicos de um exame
func (h *TopicHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	examIDStr := r.URL.Query().Get("exam_id")
	
	// Se exam_id não foi fornecido, retornar todos os tópicos
	if examIDStr == "" {
		topics, err := h.UC.ListAllTopics(context.Background())
		if err != nil {
			http.Error(w, "Erro ao buscar tópicos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Garantir que sempre retorne um array, mesmo que vazio
		if topics == nil {
			topics = []*domain.Topic{}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(topics)
		return
	}

	examID, err := strconv.Atoi(examIDStr)
	if err != nil {
		http.Error(w, "exam_id inválido", http.StatusBadRequest)
		return
	}

	topics, err := h.UC.ListTopicsByExam(context.Background(), examID)
	if err != nil {
		http.Error(w, "Erro ao buscar tópicos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Garantir que sempre retorne um array, mesmo que vazio
	if topics == nil {
		topics = []*domain.Topic{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topics)
}

// GetHandler trata o endpoint para obter um tópico específico
func (h *TopicHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
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

	topic, err := h.UC.GetTopic(context.Background(), id)
	if err != nil {
		http.Error(w, "Erro ao buscar tópico: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topic)
}

// CreateHandler trata o endpoint para criar um novo tópico
func (h *TopicHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var topic domain.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		http.Error(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.UC.CreateTopic(context.Background(), &topic); err != nil {
		http.Error(w, "Erro ao criar tópico: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topic)
}

// UpdateHandler trata o endpoint para atualizar um tópico existente
func (h *TopicHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	var topic domain.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		http.Error(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	topic.ID = id

	if err := h.UC.UpdateTopic(context.Background(), &topic); err != nil {
		http.Error(w, "Erro ao atualizar tópico: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tópico atualizado com sucesso"})
}

// DeleteHandler trata o endpoint para excluir um tópico
func (h *TopicHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := h.UC.DeleteTopic(context.Background(), id); err != nil {
		http.Error(w, "Erro ao excluir tópico: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tópico excluído com sucesso"})
}
