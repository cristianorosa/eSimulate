package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/domain"
	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// QuestionHandler lida com requisições HTTP relacionadas a questões.
type QuestionHandler struct {
	UC      *usecase.QuestionUsecase
	AreaUC  *usecase.AreaUsecase
	ExamUC  *usecase.ExamUsecase
	TopicUC *usecase.TopicUsecase
}

// ListHandler lista todas as questões, opcionalmente filtradas por exame.
func (h *QuestionHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extrair parâmetros de query
	examIDStr := r.URL.Query().Get("exam_id")

	var examID *int
	if examIDStr != "" {
		if id, err := strconv.Atoi(examIDStr); err == nil {
			examID = &id
		}
	}

	questions, err := h.UC.ListQuestions(r.Context(), examID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(questions)
}

// ListPaginatedHandler lista questões com paginação.
func (h *QuestionHandler) ListPaginatedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extrair parâmetros de query
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	examIDStr := r.URL.Query().Get("exam_id")
	topicIDStr := r.URL.Query().Get("topic_id")

	// Converter página e tamanho da página
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

	// Converter examID e topicID
	var examID *int
	if examIDStr != "" {
		if id, err := strconv.Atoi(examIDStr); err == nil {
			examID = &id
		}
	}

	var topicID *int
	if topicIDStr != "" {
		if id, err := strconv.Atoi(topicIDStr); err == nil {
			topicID = &id
		}
	}

	// Log para debug
	fmt.Printf("ListPaginatedHandler (questões): page=%d, pageSize=%d, examID=%v, topicID=%v\n", page, pageSize, examID, topicID)

	questions, pagination, err := h.UC.ListQuestionsPaginated(r.Context(), page, pageSize, examID, topicID)
	if err != nil {
		fmt.Printf("Erro ao buscar questões paginadas: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log para debug
	fmt.Printf("Questões encontradas: %d, Paginação: %+v\n", len(questions), pagination)

	// Garantir que questions seja sempre um array
	if questions == nil {
		questions = []*domain.QuestionWithDetails{}
	}

	response := map[string]interface{}{
		"data":       questions,
		"pagination": pagination,
	}

	json.NewEncoder(w).Encode(response)
}

// CreateHandler cria uma nova questão.
func (h *QuestionHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar JSON do request
	var requestData struct {
		TopicID         int    `json:"topic_id"`
		Statement       string `json:"statement"`
		Problem         string `json:"problem"`
		ContentType     string `json:"content_type"`
		Explanation     string `json:"explanation"`
		QuestionType    string `json:"question_type"`
		DifficultyLevel int    `json:"difficulty_level"`
		IsActive        bool   `json:"is_active"`
		Options         []struct {
			Text       string `json:"text"`
			IsCorrect  bool   `json:"is_correct"`
			OrderIndex int    `json:"order_index"`
		} `json:"options"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validar dados obrigatórios
	if requestData.TopicID == 0 {
		http.Error(w, "topic_id é obrigatório", http.StatusBadRequest)
		return
	}

	if requestData.Statement == "" || requestData.Problem == "" {
		http.Error(w, "statement e problem são obrigatórios", http.StatusBadRequest)
		return
	}

	if len(requestData.Options) < 2 {
		http.Error(w, "pelo menos 2 opções são obrigatórias", http.StatusBadRequest)
		return
	}

	// Converter opções para o formato do domínio
	options := make([]*domain.Option, len(requestData.Options))
	for i, opt := range requestData.Options {
		options[i] = &domain.Option{
			Text:       opt.Text,
			IsCorrect:  opt.IsCorrect,
			OrderIndex: opt.OrderIndex,
		}
	}

	// Criar questão usando o usecase
	question, err := h.UC.CreateQuestion(
		r.Context(),
		requestData.TopicID,
		requestData.Statement,
		requestData.Problem,
		requestData.ContentType,
		requestData.Explanation,
		requestData.QuestionType,
		requestData.DifficultyLevel,
		1, // createdBy (placeholder - deveria vir do token JWT)
		options,
	)

	if err != nil {
		http.Error(w, "Erro ao criar questão: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(question)
}

// UpdateHandler atualiza uma questão existente.
func (h *QuestionHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da questão da query
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID da questão é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Decodificar JSON do request
	var requestData struct {
		TopicID         int    `json:"topic_id"`
		Statement       string `json:"statement"`
		Problem         string `json:"problem"`
		ContentType     string `json:"content_type"`
		Explanation     string `json:"explanation"`
		QuestionType    string `json:"question_type"`
		DifficultyLevel int    `json:"difficulty_level"`
		IsActive        bool   `json:"is_active"`
		Options         []struct {
			Text       string `json:"text"`
			IsCorrect  bool   `json:"is_correct"`
			OrderIndex int    `json:"order_index"`
		} `json:"options"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Atualizar questão usando o usecase
	err = h.UC.UpdateQuestion(
		r.Context(),
		id,
		requestData.TopicID,
		requestData.Statement,
		requestData.Problem,
		requestData.ContentType,
		requestData.Explanation,
		requestData.QuestionType,
		requestData.DifficultyLevel,
		requestData.IsActive,
	)

	if err != nil {
		http.Error(w, "Erro ao atualizar questão: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Questão atualizada com sucesso"})
}

// DeleteHandler remove uma questão.
func (h *QuestionHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da questão da query
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID da questão é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Deletar questão usando o usecase
	err = h.UC.DeleteQuestion(r.Context(), id)
	if err != nil {
		http.Error(w, "Erro ao deletar questão: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Questão deletada com sucesso"})
}

// DetailHandler retorna os detalhes de uma questão específica.
func (h *QuestionHandler) DetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extrair ID da questão da query
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID da questão é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Buscar questão usando o usecase
	question, err := h.UC.GetQuestion(r.Context(), id)
	if err != nil {
		http.Error(w, "Erro ao buscar questão: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(question)
}

// ImportHandler importa um conjunto de questões via JSON
func (h *QuestionHandler) ImportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar o JSON de importação
	var importData domain.ImportQuestionsPayload

	if err := json.NewDecoder(r.Body).Decode(&importData); err != nil {
		http.Error(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validar dados obrigatórios
	if importData.Exam.Title == "" {
		http.Error(w, "Título do exame é obrigatório", http.StatusBadRequest)
		return
	}
	if importData.Area.Name == "" {
		http.Error(w, "Nome da área é obrigatório", http.StatusBadRequest)
		return
	}
	if len(importData.Topics) == 0 {
		http.Error(w, "Pelo menos um tópico é obrigatório", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 1) Área: buscar por nome (case-insensitive); se não existir, criar
	areas, err := h.AreaUC.ListAreas(ctx)
	if err != nil {
		http.Error(w, "Erro ao listar áreas: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var area *domain.Area
	for _, a := range areas {
		if strings.EqualFold(a.Name, importData.Area.Name) {
			area = a
			break
		}
	}
	if area == nil {
		area, err = h.AreaUC.CreateArea(ctx, importData.Area.Name, importData.Area.Description)
		if err != nil {
			http.Error(w, "Erro ao criar área: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 2) Exame: buscar por título na área; se não existir, criar
	exams, err := h.ExamUC.ListExamsByArea(ctx, area.ID)
	if err != nil {
		http.Error(w, "Erro ao listar exames: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var exam *domain.Exam
	for _, e := range exams {
		if strings.EqualFold(e.Title, importData.Exam.Title) {
			exam = e
			break
		}
	}
	if exam == nil {
		maxTime := importData.Exam.MaxTime
		if maxTime <= 0 {
			maxTime = 120
		}
		passing := importData.Exam.PassingScore
		if passing < 0 || passing > 100 {
			passing = 70
		}
		exam, err = h.ExamUC.CreateExam(ctx, importData.Exam.Title, importData.Exam.Description, area.ID, maxTime, passing, importData.Exam.IsActive, 1)
		if err != nil {
			http.Error(w, "Erro ao criar exame: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 3) Para cada tópico: buscar por nome no exame; se não existir, criar
	created := map[string]int{"areas_created": 0, "exams_created": 0, "topics_created": 0, "questions_created": 0}
	if area != nil && area.ID > 0 {
		// Se foi criada agora, incrementa
		// Não temos um flag fácil; ignorar contagem de criados por agora ou inferir por não encontrado
	}
	if exam != nil && exam.ID > 0 {
		// idem
	}

	topicsExisting, err := h.TopicUC.ListTopicsByExam(ctx, exam.ID)
	if err != nil {
		http.Error(w, "Erro ao listar tópicos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, t := range importData.Topics {
		var topic *domain.Topic
		for _, te := range topicsExisting {
			if strings.EqualFold(te.Name, t.Name) {
				topic = te
				break
			}
		}
		if topic == nil {
			topic = &domain.Topic{
				ExamID:           exam.ID,
				Name:             t.Name,
				WeightPercentage: 100,
				OrderIndex:       0,
				QuestionsCount:   t.QuestionsCount,
			}
			if err := h.TopicUC.CreateTopic(ctx, topic); err != nil {
				http.Error(w, "Erro ao criar tópico: "+err.Error(), http.StatusInternalServerError)
				return
			}
			topicsExisting = append(topicsExisting, topic)
			created["topics_created"]++
		}

		// Criar questões
		for qi, q := range t.Questions {
			// Mapear tipos
			qType := strings.ToLower(strings.TrimSpace(q.QuestionType))
			mappedQType := "objective"
			if qType == "multipla_escolha" || qType == "multiple_choice" {
				mappedQType = "multiple_choice"
			}
			contentType := q.ContentType
			if contentType == "" {
				contentType = "text"
			}
			// Dificuldade
			diff := 3
			switch strings.ToLower(strings.TrimSpace(q.Difficulty)) {
			case "facil":
				diff = 1
			case "medio":
				diff = 3
			case "dificil":
				diff = 5
			}

			// Opções
			options := make([]*domain.Option, 0, len(q.Options))
			for oi, opt := range q.Options {
				options = append(options, &domain.Option{
					Text:       opt.Text,
					IsCorrect:  opt.IsCorrect,
					OrderIndex: oi + 1,
				})
			}

			// Criar questão
			_, err := h.UC.CreateQuestion(ctx, topic.ID, q.Statement, q.Problem, contentType, q.Explanation, mappedQType, diff, 1, options)
			if err != nil {
				http.Error(w, fmt.Sprintf("Erro ao criar questão %d do tópico %s: %v", qi+1, t.Name, err), http.StatusInternalServerError)
				return
			}
			created["questions_created"]++
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Importação concluída com sucesso",
		"result":  created,
	})
}
