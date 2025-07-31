package interfaces

import (
	"encoding/json"
	"net/http"
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
	mux.HandleFunc("/questions", handlers.Question.ListHandler)
	mux.HandleFunc("/questions/create", AuthMiddleware(handlers.Question.CreateHandler))
	mux.HandleFunc("/questions/update", AuthMiddleware(handlers.Question.UpdateHandler))
	mux.HandleFunc("/questions/delete", AuthMiddleware(handlers.Question.DeleteHandler))
	mux.HandleFunc("/questions/detail", handlers.Question.DetailHandler)

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
	mux.HandleFunc("/areas", handlers.Area.ListHandler)
	mux.HandleFunc("/areas/paginated", handlers.Area.ListPaginatedHandler)
	mux.HandleFunc("/areas/create", AuthMiddleware(handlers.Area.CreateHandler))
	mux.HandleFunc("/areas/update", AuthMiddleware(handlers.Area.UpdateHandler))
	mux.HandleFunc("/areas/delete", AuthMiddleware(handlers.Area.DeleteHandler))
	mux.HandleFunc("/areas/detail", handlers.Area.GetHandler)

	// Exames - Novos handlers
	mux.HandleFunc("/exams", handlers.Exam.ListHandler)
	mux.HandleFunc("/exams/paginated", handlers.Exam.ListPaginatedHandler)
	mux.HandleFunc("/exams/create", AuthMiddleware(handlers.Exam.CreateHandler))
	mux.HandleFunc("/exams/update", AuthMiddleware(handlers.Exam.UpdateHandler))
	mux.HandleFunc("/exams/delete", AuthMiddleware(handlers.Exam.DeleteHandler))
	mux.HandleFunc("/exams/detail", handlers.Exam.GetHandler)

	// Domínios - Novos handlers
	mux.HandleFunc("/domains", handlers.Domain.ListHandler)
	mux.HandleFunc("/domains/create", AuthMiddleware(handlers.Domain.CreateHandler))
	mux.HandleFunc("/domains/update", AuthMiddleware(handlers.Domain.UpdateHandler))
	mux.HandleFunc("/domains/delete", AuthMiddleware(handlers.Domain.DeleteHandler))
	mux.HandleFunc("/domains/detail", handlers.Domain.GetHandler)

	// Aplicação de exames
	mux.HandleFunc("/user-exams/start", AuthMiddleware(handlers.UserExam.StartExam))
	mux.HandleFunc("/user-exams/submit-answer", AuthMiddleware(handlers.UserExam.SubmitAnswer))
	mux.HandleFunc("/user-exams/finish", AuthMiddleware(handlers.UserExam.FinishExam))
	mux.HandleFunc("/user-exams/detail", AuthMiddleware(handlers.UserExam.GetUserExam))
	mux.HandleFunc("/user-exams/list", AuthMiddleware(handlers.UserExam.ListByUser))
}

// Handlers agrupa todos os handlers da aplicação
// Facilita a injeção de dependências

// Handlers agrupa todos os handlers da aplicação.
// Facilita a injeção de dependências.
type Handlers struct {
	User        *UserHandler
	Theme       *ThemeHandler
	Question    *QuestionHandler
	Quiz        *QuizHandler
	History     *HistoryHandler
	Performance *PerformanceHandler
	Area        *AreaHandler
	Exam        *ExamHandler
	Domain      *DomainHandler
	UserExam    *UserExamHandler
}
