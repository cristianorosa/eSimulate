package domain

import "time"

// QuestionTag represents a tag that can be associated with questions
type QuestionTag struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// QuestionTagAssociation represents the N:N relationship between questions and tags
type QuestionTagAssociation struct {
	QuestionID int       `json:"question_id" db:"question_id"`
	TagID      int       `json:"tag_id" db:"tag_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// QuestionTagWithStats represents a tag with usage statistics
type QuestionTagWithStats struct {
	*QuestionTag
	QuestionCount int `json:"question_count" db:"question_count"`
}

// QuestionTagRepository define as operações de persistência para tags de questões
type QuestionTagRepository interface {
	// Tag Management
	CreateTag(tag *QuestionTag) error
	UpdateTag(tag *QuestionTag) error
	DeleteTag(id int) error
	FindTagByID(id int) (*QuestionTag, error)
	FindTagByName(name string) (*QuestionTag, error)
	ListTags(page, pageSize int) ([]*QuestionTag, error)
	ListTagsWithStats() ([]*QuestionTagWithStats, error)

	// Tag Association Management
	AssociateQuestionTag(questionID, tagID int) error
	DisassociateQuestionTag(questionID, tagID int) error

	// Query Operations
	ListQuestionTags(questionID int) ([]*QuestionTag, error)
	ListQuestionsByTag(tagID int) ([]int, error)

	// Bulk Operations
	BulkAssociateQuestionTags(questionID int, tagIDs []int) error
	BulkDisassociateQuestionTags(questionID int, tagIDs []int) error
	UpdateQuestionTags(questionID int, tagIDs []int) error

	// Statistics and Analytics
	GetTagStats(tagID int) (*QuestionTagWithStats, error)
	GetTagsByQuestions(questionIDs []int) (map[int][]*QuestionTag, error)

	// Search and Utilities
	SearchTags(searchTerm string, limit int) ([]*QuestionTag, error)
	CleanupUnusedTags() (int, error)
	CountTags() (int, error)
	CountTagsByQuestion(questionID int) (int, error)
}
