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
	UsageCount    int `json:"usage_count" db:"usage_count"`
}

// QuestionWithTags represents a question with its associated tags
type QuestionWithTags struct {
	QuestionID int            `json:"question_id"`
	Tags       []*QuestionTag `json:"tags"`
	TagNames   []string       `json:"tag_names"`
}

// QuestionTagRepository defines the interface for question tag operations
type QuestionTagRepository interface {
	// Tag Management
	CreateTag(tag *QuestionTag) error
	UpdateTag(tag *QuestionTag) error
	DeleteTag(id int) error
	FindTagByID(id int) (*QuestionTag, error)
	FindTagByName(name string) (*QuestionTag, error)
	ListAllTags() ([]*QuestionTag, error)
	ListTagsWithStats() ([]*QuestionTagWithStats, error)

	// Tag Association Management
	AssociateQuestionTag(questionID, tagID int) error
	DisassociateQuestionTag(questionID, tagID int) error

	// Query Operations
	ListTagsByQuestion(questionID int) ([]*QuestionTag, error)
	ListQuestionsByTag(tagID int) ([]*Question, error)
	ListQuestionsByTags(tagIDs []int) ([]*Question, error)

	// Bulk Operations
	BulkAssociateQuestionTags(questionID int, tagIDs []int) error
	BulkDisassociateQuestionTags(questionID int, tagIDs []int) error
	UpdateQuestionTags(questionID int, tagIDs []int) error

	// Statistics and Analytics
	GetTagUsageStats(tagID int) (*TagUsageStats, error)
	GetQuestionTagStats(questionID int) (*QuestionTagStats, error)
	GetMostUsedTags(limit int) ([]*QuestionTagWithStats, error)
	GetLeastUsedTags(limit int) ([]*QuestionTagWithStats, error)

	// Search and Filter
	SearchTags(query string, limit int) ([]*QuestionTag, error)
	FilterQuestionsByTags(filters *QuestionTagFilters) ([]*Question, error)
}

// QuestionTagFilters represents filters for querying questions by tags
type QuestionTagFilters struct {
	TagIDs          []int    `json:"tag_ids,omitempty"`
	TagNames        []string `json:"tag_names,omitempty"`
	ExamID          *int     `json:"exam_id,omitempty"`
	TopicID         *int     `json:"topic_id,omitempty"`
	DifficultyLevel *int     `json:"difficulty_level,omitempty"`
	MatchAll        bool     `json:"match_all"` // true = AND logic, false = OR logic
	Limit           int      `json:"limit"`
	Offset          int      `json:"offset"`
}

// TagUsageStats represents usage statistics for a specific tag
type TagUsageStats struct {
	TagID                 int                `json:"tag_id"`
	TagName               string             `json:"tag_name"`
	TotalQuestions        int                `json:"total_questions"`
	QuestionsByTopic      map[int]int        `json:"questions_by_topic"`
	QuestionsByDifficulty map[int]int        `json:"questions_by_difficulty"`
	QuestionsByExam       map[int]int        `json:"questions_by_exam"`
	MostUsedWith          []*TagCooccurrence `json:"most_used_with"`
}

// QuestionTagStats represents tag statistics for a specific question
type QuestionTagStats struct {
	QuestionID       int      `json:"question_id"`
	TotalTags        int      `json:"total_tags"`
	TagNames         []string `json:"tag_names"`
	TagCategories    []string `json:"tag_categories"`
	SimilarQuestions []int    `json:"similar_questions"`
}

// TagCooccurrence represents how often two tags appear together
type TagCooccurrence struct {
	TagID       int     `json:"tag_id"`
	TagName     string  `json:"tag_name"`
	Occurrences int     `json:"occurrences"`
	Percentage  float64 `json:"percentage"`
}

// TagSuggestion represents a tag suggestion for a question
type TagSuggestion struct {
	TagID      int     `json:"tag_id"`
	TagName    string  `json:"tag_name"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}

// QuestionTagService defines business logic for question tags
type QuestionTagService interface {
	// Suggest tags for a question based on content, topic, difficulty, etc.
	SuggestTagsForQuestion(questionID int) ([]*TagSuggestion, error)

	// Find similar questions based on tag overlap
	FindSimilarQuestionsByTags(questionID int, limit int) ([]*Question, error)

	// Validate tag associations (business rules)
	ValidateTagAssociation(questionID, tagID int) error

	// Auto-tag questions based on content analysis
	AutoTagQuestion(questionID int) error

	// Clean up unused tags
	CleanupUnusedTags() (int, error)

	// Merge duplicate tags
	MergeTags(sourceTagID, targetTagID int) error
}
