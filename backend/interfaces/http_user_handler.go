package interfaces

import (
	"context"
	"encoding/json"
	"net/http"

	"eSimulate/backend/usecase"
)

// UserHandler lida com requisições HTTP relacionadas a usuários

import (
	"eSimulate/backend/infra"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UC *usecase.UserUsecase
}

// RegisterHandler trata o endpoint de cadastro de usuário
func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	user, err := h.UC.RegisterUser(context.Background(), req.Name, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"created_at": user.CreatedAt,
	})
}

// LoginHandler trata o endpoint de login tradicional
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	user, err := h.UC.Repo.FindByEmail(req.Email)
	if err != nil || user == nil {
		http.Error(w, "Usuário ou senha inválidos", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Usuário ou senha inválidos", http.StatusUnauthorized)
		return
	}
	token, err := infra.GenerateJWT(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// GoogleLoginHandler trata o endpoint de login via Google (esqueleto)
func (h *UserHandler) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Aqui deve ser implementado o fluxo OAuth2 do Google
	// Receber o token do Google, validar, buscar dados do usuário
	// Se usuário não existir, criar; se existir, autenticar
	// Gerar e retornar JWT do sistema
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Login via Google ainda não implementado"))
}

// FacebookLoginHandler trata o endpoint de login via Facebook (esqueleto)
func (h *UserHandler) FacebookLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Aqui deve ser implementado o fluxo OAuth2 do Facebook
	// Receber o token do Facebook, validar, buscar dados do usuário
	// Se usuário não existir, criar; se existir, autenticar
	// Gerar e retornar JWT do sistema
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Login via Facebook ainda não implementado"))
}
