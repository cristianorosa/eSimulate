package domain

import "time"

// QuestionType define os tipos de questão
type QuestionType string

const (
	QuestionTypeObjective      QuestionType = "objective"       // Objetiva (uma resposta correta)
	QuestionTypeMultipleChoice QuestionType = "multiple_choice" // Múltipla escolha (várias respostas corretas)
)

// QuestionContentType define o tipo de conteúdo da questão
type QuestionContentType string

const (
	QuestionContentTypeText QuestionContentType = "text" // Texto simples
	QuestionContentTypeCode QuestionContentType = "code" // Código com syntax highlighting
)

// Question representa uma questão de um exame
type Question struct {
	ID              int                 `json:"id"`
	TopicID         int                 `json:"topic_id"`     // Relacionamento direto apenas com tópico
	Statement       string              `json:"statement"`    // Enunciado da questão
	Problem         string              `json:"problem"`      // Problema (texto ou código)
	ContentType     QuestionContentType `json:"content_type"` // Tipo de conteúdo (texto ou código)
	Explanation     string              `json:"explanation"`
	QuestionType    QuestionType        `json:"question_type"`
	DifficultyLevel int                 `json:"difficulty_level"`
	CreatedBy       int                 `json:"created_by"`
	IsActive        bool                `json:"is_active"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at,omitempty"`
	Options         []*Option           `json:"options"`
}

// QuestionWithDetails representa uma questão com informações detalhadas do exame e tópico
type QuestionWithDetails struct {
	*Question
	ExamID    int    `json:"exam_id"`
	ExamTitle string `json:"exam_title"`
	TopicName string `json:"topic_name"`
}

// Option representa uma opção de resposta de uma questão
type Option struct {
	ID          int    `json:"id"`
	QuestionID  int    `json:"question_id"`
	Text        string `json:"text"`
	IsCorrect   bool   `json:"is_correct"`
	Explanation string `json:"explanation"`
	OrderIndex  int    `json:"order_index"`
}

// QuestionRepository define as operações de persistência para questões
type QuestionRepository interface {
	Create(question *Question) error
	Update(question *Question) error
	Delete(id int) error
	FindByID(id int) (*Question, error)
	ListByExam(examID int) ([]*Question, error)
	ListByTopic(topicID int) ([]*Question, error)
	ListAll() ([]*Question, error)
	ListAllWithDetails() ([]*QuestionWithDetails, error)
}

// OptionRepository define as operações de persistência para opções
type OptionRepository interface {
	Create(option *Option) error
	Update(option *Option) error
	Delete(id int) error
	DeleteByQuestionID(questionID int) error
	FindByQuestionID(questionID int) ([]*Option, error)
}
