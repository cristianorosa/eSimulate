package interfaces

import (
	"encoding/json"
	"net/http"
	"strings"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(mux *http.ServeMux, handlers *Handlers) {
	// Autenticação e usuário
	mux.HandleFunc("/auth/register", handlers.User.RegisterHandler)
	mux.HandleFunc("/auth/login", handlers.User.LoginHandler)
	mux.HandleFunc("/users", AuthMiddleware(handlers.User.ListUsersHandler)) // Nova rota protegida
	// Esqueleto para login social
	// mux.HandleFunc("/auth/google", handlers.User.GoogleLoginHandler)
	// mux.HandleFunc("/auth/facebook", handlers.User.FacebookLoginHandler)
	// mux.HandleFunc("/users/me", interfaces.AuthMiddleware(handlers.User.MeHandler))
	// mux.HandleFunc("/users/me/update", interfaces.AuthMiddleware(handlers.User.UpdateHandler))

	// Temas
	mux.HandleFunc("/themes", handlers.Theme.ListHandler)
	mux.HandleFunc("/themes/create", AuthMiddleware(handlers.Theme.CreateHandler))
	mux.HandleFunc("/themes/update", AuthMiddleware(handlers.Theme.UpdateHandler))
	mux.HandleFunc("/themes/delete", AuthMiddleware(handlers.Theme.DeleteHandler))

	// Questões
	mux.HandleFunc("/questions/paginated", handlers.Question.ListPaginatedHandler)
	mux.HandleFunc("/questions/create", AuthMiddleware(handlers.Question.CreateHandler))
	mux.HandleFunc("/questions/update", AuthMiddleware(handlers.Question.UpdateHandler))
	mux.HandleFunc("/questions/delete", AuthMiddleware(handlers.Question.DeleteHandler))
	mux.HandleFunc("/questions/detail", handlers.Question.DetailHandler)
	mux.HandleFunc("/questions/import", AuthMiddleware(handlers.Question.ImportHandler))
	mux.HandleFunc("/questions", handlers.Question.ListHandler)

	// Simulados
	mux.HandleFunc("/quizzes", handlers.Quiz.ListHandler)
	mux.HandleFunc("/quizzes/create", AuthMiddleware(handlers.Quiz.CreateHandler))
	mux.HandleFunc("/quizzes/detail", handlers.Quiz.DetailHandler)
	mux.HandleFunc("/quizzes/start", AuthMiddleware(handlers.Quiz.StartHandler))
	mux.HandleFunc("/quizzes/answer", AuthMiddleware(handlers.Quiz.AnswerHandler))
	mux.HandleFunc("/quizzes/result", AuthMiddleware(handlers.Quiz.ResultHandler))

	// Histórico e desempenho
	mux.HandleFunc("/history", AuthMiddleware(handlers.History.ListHandler))
	mux.HandleFunc("/performance", AuthMiddleware(handlers.Performance.ReportHandler))

	// Endpoint de teste para verificar conectividade
	mux.HandleFunc("/test-db", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Backend funcionando"})
	})

	// Áreas - Novos handlers
	mux.HandleFunc("/areas/paginated", handlers.Area.ListPaginatedHandler)
	mux.HandleFunc("/areas/create", AuthMiddleware(handlers.Area.CreateHandler))
	mux.HandleFunc("/areas/update", AuthMiddleware(handlers.Area.UpdateHandler))
	mux.HandleFunc("/areas/delete", AuthMiddleware(handlers.Area.DeleteHandler))
	mux.HandleFunc("/areas/detail", handlers.Area.GetHandler)
	mux.HandleFunc("/areas", handlers.Area.ListHandler)

	// Exames - Novos handlers
	mux.HandleFunc("/exams/paginated", handlers.Exam.ListPaginatedHandler)
	mux.HandleFunc("/exams/create", AuthMiddleware(handlers.Exam.CreateHandler))
	mux.HandleFunc("/exams/update", AuthMiddleware(handlers.Exam.UpdateHandler))
	mux.HandleFunc("/exams/delete", AuthMiddleware(handlers.Exam.DeleteHandler))
	mux.HandleFunc("/exams/detail", handlers.Exam.GetHandler)
	mux.HandleFunc("/exams", handlers.Exam.ListHandler)

	// Tópicos - Novos handlers
	mux.HandleFunc("/topics/paginated", handlers.Topic.ListPaginatedHandler)
	mux.HandleFunc("/topics/create", AuthMiddleware(handlers.Topic.CreateHandler))
	mux.HandleFunc("/topics/update", AuthMiddleware(handlers.Topic.UpdateHandler))
	mux.HandleFunc("/topics/delete", AuthMiddleware(handlers.Topic.DeleteHandler))
	mux.HandleFunc("/topics/detail", handlers.Topic.GetHandler)
	mux.HandleFunc("/topics", handlers.Topic.ListHandler)

	// Tags - Seguindo o mesmo padrão dos outros endpoints
	mux.HandleFunc("/tags/create", AuthMiddleware(handlers.QuestionTag.CreateTagHandler))
	mux.HandleFunc("/tags/search", handlers.QuestionTag.SearchTagsHandler)
	mux.HandleFunc("/tags", handlers.QuestionTag.ListTagsHandler)

	// Aplicação de exames
	mux.HandleFunc("/user-exams/start", AuthMiddleware(handlers.UserExam.StartExam))
	mux.HandleFunc("/user-exams/submit-answer", AuthMiddleware(handlers.UserExam.SubmitAnswer))
	mux.HandleFunc("/user-exams/finish", AuthMiddleware(handlers.UserExam.FinishExam))
	mux.HandleFunc("/user-exams/detail", AuthMiddleware(handlers.UserExam.GetUserExam))
	mux.HandleFunc("/user-exams/list", AuthMiddleware(handlers.UserExam.ListByUser))

	// === N:N RELATIONSHIPS ENDPOINTS ===
	// Using specific paths to avoid conflicts with basic CRUD operations

	// Exam-Question Relationships - specific paths only
	mux.HandleFunc("/exams/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/questions") {
			handleExamQuestionRoutes(handlers.ExamQuestion)(w, r)
		} else if strings.Contains(r.URL.Path, "/topics") {
			handleExamTopicRoutes(handlers.ExamTopic)(w, r)
		}
	})

	// Question Tags Management - usando rotas diretas acima
	// mux.HandleFunc("/tags/", handleTagRoutes(handlers.QuestionTag))

	// Question-Tag relationships - only specific tag paths to avoid CRUD conflicts
	mux.HandleFunc("/questions/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/tags") {
			handleQuestionTagRoutes(handlers.QuestionTag)(w, r)
		}
		// Don't handle other /questions/ paths - let them fall through to basic CRUD
	})

	// Topic-Exam Relationships - specific paths only
	mux.HandleFunc("/topics/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/exams") {
			handleTopicExamRoutes(handlers.ExamTopic)(w, r)
		}
	})
}

// Handlers agrupa todos os handlers da aplicação
// Facilita a injeção de dependências

// Handlers agrupa todos os handlers da aplicação.
// Facilita a injeção de dependências.
type Handlers struct {
	User         *UserHandler
	Theme        *ThemeHandler
	Question     *QuestionHandler
	Quiz         *QuizHandler
	History      *HistoryHandler
	Performance  *PerformanceHandler
	Area         *AreaHandler
	Exam         *ExamHandler
	Topic        *TopicHandler
	UserExam     *UserExamHandler
	ExamQuestion *ExamQuestionHandler
	QuestionTag  *QuestionTagHandler
	ExamTopic    *ExamTopicHandler
}

// === ROUTE HANDLERS FOR N:N RELATIONSHIPS ===

// handleExamQuestionRoutes handles exam-question relationship routes
func handleExamQuestionRoutes(handler *ExamQuestionHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		switch {
		case matchRoute(path, "/exams/*/questions/bulk") && method == "POST":
			handler.BulkAssociateQuestionsHandler(w, r)
		case matchRoute(path, "/exams/*/questions/bulk") && method == "DELETE":
			handler.BulkDisassociateQuestionsHandler(w, r)
		case matchRoute(path, "/exams/*/questions/reorder") && method == "PUT":
			handler.ReorderExamQuestionsHandler(w, r)
		case matchRoute(path, "/exams/*/questions/stats") && method == "GET":
			handler.GetExamQuestionStatsHandler(w, r)
		case matchRoute(path, "/exams/*/questions/available") && method == "GET":
			handler.GetAvailableQuestionsHandler(w, r)
		case matchRoute(path, "/exams/*/questions/*/*") && method == "POST":
			handler.AssociateQuestionHandler(w, r)
		case matchRoute(path, "/exams/*/questions/*/*") && method == "DELETE":
			handler.DisassociateQuestionHandler(w, r)
		case matchRoute(path, "/exams/*/questions") && method == "GET":
			handler.GetExamQuestionsHandler(w, r)
		case matchRoute(path, "/exams/*/validate") && method == "GET":
			handler.ValidateExamConfigurationHandler(w, r)
		}
	}
}

// handleQuestionExamRoutes handles question-exam relationship routes
func handleQuestionExamRoutes(handler *ExamQuestionHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		switch {
		case matchRoute(path, "/questions/*/exams") && method == "GET":
			handler.GetQuestionExamsHandler(w, r)
		}
	}
}

// handleTagRoutes handles tag management routes
func handleTagRoutes(handler *QuestionTagHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		switch {
		case matchRoute(path, "/tags/search") && method == "GET":
			handler.SearchTagsHandler(w, r)
		case matchRoute(path, "/tags/stats") && method == "GET":
			handler.ListTagsWithStatsHandler(w, r)
		case matchRoute(path, "/tags/find-or-create") && method == "POST":
			handler.FindOrCreateTagHandler(w, r) // Temporariamente sem auth para teste
		case matchRoute(path, "/tags/*/stats") && method == "GET":
			handler.GetTagStatsHandler(w, r)
		case (path == "/tags" || path == "/tags/") && method == "GET":
			handler.ListTagsHandler(w, r)
		case (path == "/tags" || path == "/tags/") && method == "POST":
			AuthMiddleware(handler.CreateTagHandler)(w, r)
		case matchRoute(path, "/tags/*") && method == "GET" && len(strings.Split(path, "/")) > 3:
			handler.GetTagHandler(w, r)
		case matchRoute(path, "/tags/*") && method == "PUT":
			AuthMiddleware(handler.UpdateTagHandler)(w, r)
		case matchRoute(path, "/tags/*") && method == "DELETE":
			AuthMiddleware(handler.DeleteTagHandler)(w, r)
		}
	}
}

// handleQuestionTagRoutes handles question-tag relationship routes
func handleQuestionTagRoutes(handler *QuestionTagHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		switch {
		case matchRoute(path, "/questions/*/tags/bulk") && method == "POST":
			AuthMiddleware(handler.BulkAssociateQuestionTagsHandler)(w, r)
		case matchRoute(path, "/questions/*/tags/*/*") && method == "POST":
			AuthMiddleware(handler.AssociateQuestionTagHandler)(w, r)
		case matchRoute(path, "/questions/*/tags/*/*") && method == "DELETE":
			AuthMiddleware(handler.DisassociateQuestionTagHandler)(w, r)
		case matchRoute(path, "/questions/*/tags") && method == "GET":
			handler.GetQuestionTagsHandler(w, r)
		case matchRoute(path, "/questions/*/tags") && method == "PUT":
			AuthMiddleware(handler.UpdateQuestionTagsHandler)(w, r)

		}
	}
}

// handleExamTopicRoutes handles exam-topic relationship routes
func handleExamTopicRoutes(handler *ExamTopicHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		switch {
		case matchRoute(path, "/exams/*/topics/bulk") && method == "PUT":
			AuthMiddleware(handler.BulkUpdateExamTopicsHandler)(w, r)
		case matchRoute(path, "/exams/*/topics/reorder") && method == "PUT":
			AuthMiddleware(handler.ReorderTopicsHandler)(w, r)
		case matchRoute(path, "/exams/*/topics/auto-distribute-weights") && method == "POST":
			AuthMiddleware(handler.AutoDistributeWeightsHandler)(w, r)
		case matchRoute(path, "/exams/*/topics/suggest-distribution") && method == "GET":
			handler.SuggestOptimalDistributionHandler(w, r)
		case matchRoute(path, "/exams/*/topics/validate") && method == "GET":
			handler.ValidateExamConfigurationHandler(w, r)
		case matchRoute(path, "/exams/*/topics/*/calculate-difficulty") && method == "POST":
			handler.CalculateQuestionsByDifficultyHandler(w, r)
		case matchRoute(path, "/exams/*/topics/*/*") && method == "PUT":
			AuthMiddleware(handler.UpdateExamTopicHandler)(w, r)
		case matchRoute(path, "/exams/*/topics/*/*") && method == "DELETE":
			AuthMiddleware(handler.DisassociateTopicHandler)(w, r)
		case matchRoute(path, "/exams/*/topics") && method == "GET":
			handler.GetExamTopicsHandler(w, r)
		case matchRoute(path, "/exams/*/topics") && method == "POST":
			AuthMiddleware(handler.AssociateTopicHandler)(w, r)
		case matchRoute(path, "/exams/*/stats") && method == "GET":
			handler.GetExamStatsHandler(w, r)
		}
	}
}

// handleTopicExamRoutes handles topic-exam relationship routes
func handleTopicExamRoutes(handler *ExamTopicHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		switch {
		case matchRoute(path, "/topics/*/exams") && method == "GET":
			handler.GetTopicExamsHandler(w, r)
		}
	}
}

// matchRoute checks if a path matches a pattern with wildcards (*)
func matchRoute(path, pattern string) bool {
	pathParts := strings.Split(path, "/")
	patternParts := strings.Split(pattern, "/")

	if len(pathParts) != len(patternParts) {
		return false
	}

	for i, part := range patternParts {
		if part != "*" && part != pathParts[i] {
			return false
		}
	}

	return true
}
