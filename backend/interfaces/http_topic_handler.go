package interfaces

import (
	"context"
	"encoding/json"
	"fmt"
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

// ListPaginatedHandler lista tópicos com paginação
func (h *TopicHandler) ListPaginatedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair parâmetros de query
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	examIDStr := r.URL.Query().Get("exam_id")

	// Converter parâmetros
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	var examID *int
	if examIDStr != "" {
		if id, err := strconv.Atoi(examIDStr); err == nil {
			examID = &id
		}
	}

	// Log para debug
	fmt.Printf("ListPaginatedHandler: page=%d, pageSize=%d, examID=%v\n", page, pageSize, examID)

	// Buscar tópicos paginados
	topics, pagination, err := h.UC.ListTopicsPaginated(context.Background(), page, pageSize, examID)
	if err != nil {
		fmt.Printf("Erro ao buscar tópicos paginados: %v\n", err)
		http.Error(w, "Erro ao buscar tópicos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log para debug
	fmt.Printf("Tópicos encontrados: %d, Paginação: %+v\n", len(topics), pagination)

	// Criar resposta paginada
	response := map[string]interface{}{
		"data":       topics,
		"pagination": pagination,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
