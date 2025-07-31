package interfaces

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cristianorosa/eSimulate/backend/usecase"
)

// UserExamHandler gerencia requisições HTTP para aplicação de exames
type UserExamHandler struct {
	UC *usecase.UserExamUsecase
}

// StartExamRequest representa a requisição para iniciar um exame
type StartExamRequest struct {
	UserID int `json:"user_id"`
	ExamID int `json:"exam_id"`
}

// SubmitAnswerRequest representa a requisição para submeter uma resposta
type SubmitAnswerRequest struct {
	UserExamID      int  `json:"user_exam_id"`
	QuestionID      int  `json:"question_id"`
	OptionID        int  `json:"option_id"`
	MarkedForReview bool `json:"marked_for_review"`
}

// StartExam inicia um exame para um usuário
func (h *UserExamHandler) StartExam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StartExamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userExam, err := h.UC.StartExam(req.UserID, req.ExamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userExam)
}

// SubmitAnswer submete uma resposta do usuário
func (h *UserExamHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SubmitAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.UC.SubmitAnswer(req.UserExamID, req.QuestionID, req.OptionID, req.MarkedForReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// FinishExam finaliza um exame
func (h *UserExamHandler) FinishExam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userExam, err := h.UC.FinishExam(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userExam)
}

// GetUserExam busca um exame aplicado por ID
func (h *UserExamHandler) GetUserExam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userExam, err := h.UC.GetUserExam(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userExam)
}

// ListByUser lista todos os exames aplicados por um usuário
func (h *UserExamHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userExams, err := h.UC.ListByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userExams)
}
