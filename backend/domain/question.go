package domain

// Question representa uma questão de simulado
type Question struct {
	ID          int
	ThemeID     int
	Statement   string
	Explanation string
	CreatedBy   int
	Options     []*Option
}

// Option representa uma opção de resposta para uma questão
type Option struct {
	ID          int
	QuestionID  int
	Text        string
	IsCorrect   bool
	Explanation string
}

// QuestionRepository define a interface para operações de persistência de questões
type QuestionRepository interface {
	Create(q *Question) error
	Update(q *Question) error
	Delete(id int) error
	FindByID(id int) (*Question, error)
	ListAll(themeID *int) ([]*Question, error)
}
