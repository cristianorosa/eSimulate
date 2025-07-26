package domain

// UserRepository define as operações de persistência para usuários
// Interface para facilitar testes e desacoplamento (Clean Architecture)
type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
}
