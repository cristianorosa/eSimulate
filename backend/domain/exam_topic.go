package domain

import "time"

// ExamTopic represents the N:N relationship between exams and topics
type ExamTopic struct {
	ExamID                     int       `json:"exam_id" db:"exam_id"`
	TopicID                    int       `json:"topic_id" db:"topic_id"`
	QuestionsCount             int       `json:"questions_count" db:"questions_count"`
	WeightPercentage           float64   `json:"weight_percentage" db:"weight_percentage"`
	OrderIndex                 int       `json:"order_index" db:"order_index"`
	DifficultyEasyPercentage   float64   `json:"difficulty_easy_percentage" db:"difficulty_easy_percentage"`
	DifficultyMediumPercentage float64   `json:"difficulty_medium_percentage" db:"difficulty_medium_percentage"`
	DifficultyHardPercentage   float64   `json:"difficulty_hard_percentage" db:"difficulty_hard_percentage"`
	CreatedAt                  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at" db:"updated_at"`
}

// ExamTopicWithDetails represents exam-topic relationship with detailed information
type ExamTopicWithDetails struct {
	*ExamTopic
	ExamTitle          string  `json:"exam_title" db:"exam_title"`
	ExamDescription    *string `json:"exam_description" db:"exam_description"`
	TopicName          string  `json:"topic_name" db:"topic_name"`
	TopicDescription   *string `json:"topic_description" db:"topic_description"`
	ActualQuestions    int     `json:"actual_questions" db:"actual_questions"`
	EasyQuestions      int     `json:"easy_questions" db:"easy_questions"`
	MediumQuestions    int     `json:"medium_questions" db:"medium_questions"`
	HardQuestions      int     `json:"hard_questions" db:"hard_questions"`
	IsComplete         bool    `json:"is_complete" db:"is_complete"`
}

// ExamTopicRepository defines the interface for exam-topic relationship operations
type ExamTopicRepository interface {
	// Create creates a new exam-topic association
	Create(examTopic *ExamTopic) error
	
	// Update updates an existing exam-topic association
	Update(examTopic *ExamTopic) error
	
	// Delete removes an exam-topic association
	Delete(examID, topicID int) error
	
	// FindByExamAndTopic finds a specific exam-topic association
	FindByExamAndTopic(examID, topicID int) (*ExamTopic, error)
	
	// ListByExam returns all topics associated with an exam
	ListByExam(examID int) ([]*ExamTopic, error)
	
	// ListByTopic returns all exams associated with a topic
	ListByTopic(topicID int) ([]*ExamTopic, error)
	
	// ListByExamWithDetails returns topics for an exam with detailed information
	ListByExamWithDetails(examID int) ([]*ExamTopicWithDetails, error)
	
	// ListByTopicWithDetails returns exams for a topic with detailed information
	ListByTopicWithDetails(topicID int) ([]*ExamTopicWithDetails, error)
	
	// GetExamTopicsCount returns the total number of topics in an exam
	GetExamTopicsCount(examID int) (int, error)
	
	// GetTopicExamsCount returns the total number of exams a topic belongs to
	GetTopicExamsCount(topicID int) (int, error)
	
	// ValidateWeightDistribution validates that weight percentages sum to 100 for an exam
	ValidateWeightDistribution(examID int) error
	
	// ValidateDifficultyDistribution validates that difficulty percentages sum to 100 for each topic
	ValidateDifficultyDistribution(examTopic *ExamTopic) error
	
	// CalculateQuestionsByDifficulty calculates how many questions of each difficulty are needed
	CalculateQuestionsByDifficulty(examTopic *ExamTopic) (*DifficultyDistribution, error)
	
	// UpdateQuestionsCount updates the actual questions count for a topic in an exam
	UpdateQuestionsCount(examID, topicID, count int) error
	
	// BulkCreate creates multiple exam-topic associations
	BulkCreate(examTopics []*ExamTopic) error
	
	// BulkUpdate updates multiple exam-topic associations
	BulkUpdate(examTopics []*ExamTopic) error
	
	// ReorderTopics reorders all topics in an exam
	ReorderTopics(examID int, topicOrders map[int]int) error
}

// DifficultyDistribution represents the calculated distribution of questions by difficulty
type DifficultyDistribution struct {
	EasyCount   int `json:"easy_count"`
	MediumCount int `json:"medium_count"`
	HardCount   int `json:"hard_count"`
	TotalCount  int `json:"total_count"`
}

// ExamTopicFilters represents filters for querying exam-topic relationships
type ExamTopicFilters struct {
	ExamID           *int     `json:"exam_id,omitempty"`
	TopicID          *int     `json:"topic_id,omitempty"`
	WeightFrom       *float64 `json:"weight_from,omitempty"`
	WeightTo         *float64 `json:"weight_to,omitempty"`
	QuestionsFrom    *int     `json:"questions_from,omitempty"`
	QuestionsTo      *int     `json:"questions_to,omitempty"`
	IncludeComplete  *bool    `json:"include_complete,omitempty"`
	IncludeIncomplete *bool   `json:"include_incomplete,omitempty"`
	Limit            int      `json:"limit"`
	Offset           int      `json:"offset"`
}

// ExamTopicStats represents statistics for exam-topic relationships
type ExamTopicStats struct {
	ExamID              int                      `json:"exam_id"`
	TotalTopics         int                      `json:"total_topics"`
	TotalWeight         float64                  `json:"total_weight"`
	TotalQuestions      int                      `json:"total_questions"`
	CompletedTopics     int                      `json:"completed_topics"`
	TopicDistribution   []*TopicDistribution     `json:"topic_distribution"`
	DifficultyBreakdown *DifficultyBreakdown     `json:"difficulty_breakdown"`
	IsValid             bool                     `json:"is_valid"`
	ValidationErrors    []string                 `json:"validation_errors,omitempty"`
}

// DifficultyBreakdown represents the overall difficulty distribution for an exam
type DifficultyBreakdown struct {
	EasyQuestions   int     `json:"easy_questions"`
	MediumQuestions int     `json:"medium_questions"`
	HardQuestions   int     `json:"hard_questions"`
	EasyPercentage  float64 `json:"easy_percentage"`
	MediumPercentage float64 `json:"medium_percentage"`
	HardPercentage  float64 `json:"hard_percentage"`
}

// ExamTopicService defines business logic for exam-topic relationships
type ExamTopicService interface {
	// AssociateTopicWithExam creates a new exam-topic association with validation
	AssociateTopicWithExam(examTopic *ExamTopic) error
	
	// DisassociateTopicFromExam removes an exam-topic association
	DisassociateTopicFromExam(examID, topicID int) error
	
	// UpdateExamTopic updates an existing exam-topic association with validation
	UpdateExamTopic(examTopic *ExamTopic) error
	
	// GetExamTopics returns all topics for an exam with detailed information
	GetExamTopics(examID int) ([]*ExamTopicWithDetails, error)
	
	// ValidateDifficultyDistribution validates that difficulty percentages are valid
	ValidateDifficultyDistribution(examTopic *ExamTopic) error
	
	// CalculateQuestionsByDifficulty calculates the number of questions needed for each difficulty level
	CalculateQuestionsByDifficulty(examTopic *ExamTopic) (*DifficultyDistribution, error)
	
	// GetExamStats returns comprehensive statistics for an exam
	GetExamStats(examID int) (*ExamTopicStats, error)
	
	// ValidateExamConfiguration validates the complete exam configuration
	ValidateExamConfiguration(examID int) (bool, []string)
	
	// AutoDistributeWeights automatically distributes weights evenly among topics
	AutoDistributeWeights(examID int) error
	
	// SuggestOptimalDistribution suggests optimal weight and difficulty distribution
	SuggestOptimalDistribution(examID int) (*ExamTopicStats, error)
}
