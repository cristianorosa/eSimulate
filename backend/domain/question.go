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

	// N:N Relationships
	ExamIDs  []int          `json:"exam_ids,omitempty"`  // IDs dos exames que contêm esta questão
	Exams    []*Exam        `json:"exams,omitempty"`     // Exames que contêm esta questão
	TagIDs   []int          `json:"tag_ids,omitempty"`   // IDs das tags associadas
	Tags     []*QuestionTag `json:"tags,omitempty"`      // Tags associadas à questão
	TagNames []string       `json:"tag_names,omitempty"` // Nomes das tags (para facilitar exibição)
}

// QuestionWithDetails representa uma questão com informações detalhadas do exame e tópico
type QuestionWithDetails struct {
	*Question
	ExamCount    int      `json:"exam_count"`     // Número de exames que contêm esta questão
	ExamTitles   []string `json:"exam_titles"`    // Títulos dos exames
	TagCount     int      `json:"tag_count"`      // Número de tags associadas
	TagNamesList []string `json:"tag_names_list"` // Lista de nomes das tags
	TopicName    string   `json:"topic_name"`     // Nome do tópico

	// Campos legados para compatibilidade
	ExamID    int    `json:"exam_id,omitempty"`    // Deprecated: usar ExamIDs
	ExamTitle string `json:"exam_title,omitempty"` // Deprecated: usar ExamTitles
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
	// Basic CRUD operations
	Create(question *Question) error
	Update(question *Question) error
	Delete(id int) error
	FindByID(id int) (*Question, error)

	// Query operations
	ListByExam(examID int) ([]*Question, error)
	ListByTopic(topicID int) ([]*Question, error)
	ListAll() ([]*Question, error)
	ListAllWithDetails() ([]*QuestionWithDetails, error)
	ListPaginated(page, pageSize int, examID, topicID *int) ([]*QuestionWithDetails, *Pagination, error)

	// N:N Relationship operations
	FindByIDWithExams(id int) (*Question, error)
	ListExamsByQuestion(questionID int) ([]*Exam, error)
	ListByExamDirect(examID int) ([]*Question, error) // Direct relationship via exam_questions
	ListByExamWithDetails(examID int) ([]*QuestionWithDetails, error)

	// Question availability and filtering
	GetAvailableQuestionsForExam(examID int, topicID *int) ([]*Question, error)

	// Statistics and analytics
	GetQuestionStats(questionID int) (*QuestionStats, error)
	GetTopicQuestionCount(topicID int) (int, error)
	GetExamQuestionCount(examID int) (int, error)
}

// QuestionStats represents statistics for a question
type QuestionStats struct {
	QuestionID      int     `json:"question_id"`
	ExamCount       int     `json:"exam_count"`
	TagCount        int     `json:"tag_count"`
	AnswerCount     int     `json:"answer_count"`
	CorrectRate     float64 `json:"correct_rate"`
	AverageTime     float64 `json:"average_time"`
	DifficultyLevel int     `json:"difficulty_level"`
}

// OptionRepository define as operações de persistência para opções
type OptionRepository interface {
	Create(option *Option) error
	Update(option *Option) error
	Delete(id int) error
	DeleteByQuestionID(questionID int) error
	FindByQuestionID(questionID int) ([]*Option, error)
}
