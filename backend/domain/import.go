package domain

// ImportQuestionsPayload representa o payload de importação de questões via JSON
type ImportQuestionsPayload struct {
	Exam struct {
		Title        string  `json:"title"`
		Description  string  `json:"description"`
		MaxTime      int     `json:"max_time"`
		PassingScore float64 `json:"passing_score"`
		IsActive     bool    `json:"is_active"`
	} `json:"exam"`
	Area struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"area"`
	Topics []struct {
		Name           string `json:"name"`
		QuestionsCount int    `json:"questions_count"`
		Questions      []struct {
			Statement    string `json:"statement"`
			Problem      string `json:"problem"`
			ContentType  string `json:"content_type"`
			Explanation  string `json:"explanation"`
			QuestionType string `json:"question_type"`
			Difficulty   string `json:"difficulty"`
			IsActive     bool   `json:"is_active"`
			Options      []struct {
				Text      string `json:"text"`
				IsCorrect bool   `json:"is_correct"`
			} `json:"options"`
		} `json:"questions"`
	} `json:"topics"`
}
