package interfaces

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"eSimulate/backend/domain"
	"eSimulate/backend/usecase"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	users map[string]*domain.User
}

func (m *mockUserRepo) Create(user *domain.User) error {
	if _, exists := m.users[user.Email]; exists {
		return domain.ErrUserExists
	}
	user.ID = len(m.users) + 1
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) FindByEmail(email string) (*domain.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, nil
}

func TestLoginHandler(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*domain.User)}
	uc := &usecase.UserUsecase{Repo: repo}
	h := &UserHandler{UC: uc}

	// Cria usuário de teste
	hash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	repo.users["teste@exemplo.com"] = &domain.User{
		ID:           1,
		Name:         "Teste",
		Email:        "teste@exemplo.com",
		PasswordHash: string(hash),
	}

	// Teste de login válido
	body := map[string]string{"email": "teste@exemplo.com", "password": "senha123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.LoginHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("esperado status 200, obteve %d", w.Code)
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["token"] == nil {
		t.Error("esperado token JWT na resposta")
	}

	// Teste de login inválido
	body = map[string]string{"email": "teste@exemplo.com", "password": "errada"}
	b, _ = json.Marshal(body)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(b))
	w = httptest.NewRecorder()
	h.LoginHandler(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado status 401 para senha errada, obteve %d", w.Code)
	}
}
