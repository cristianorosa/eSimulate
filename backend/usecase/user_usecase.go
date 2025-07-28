package usecase

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/cristianorosa/eSimulate/backend/domain"
	"golang.org/x/crypto/bcrypt"
)

// UserUsecase implementa as regras de negócio para usuários
// Inclui validação, hash de senha e tratamento seguro dos dados

type UserUsecase struct {
	Repo domain.UserRepository
}

// RegisterUser cadastra um novo usuário, validando dados e aplicando hash seguro na senha
func (uc *UserUsecase) RegisterUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))
	if name == "" || email == "" || password == "" {
		return nil, errors.New("nome, email e senha são obrigatórios")
	}
	if !isValidEmail(email) {
		return nil, errors.New("email inválido")
	}
	if len(password) < 6 {
		return nil, errors.New("a senha deve ter pelo menos 6 caracteres")
	}
	// Verifica se já existe usuário
	existing, _ := uc.Repo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("email já cadastrado")
	}
	// Hash seguro da senha
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("erro ao gerar hash da senha")
	}
	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}
	if err := uc.Repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// isValidEmail faz uma validação simples de email
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
