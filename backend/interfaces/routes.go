package interfaces

import (
	"net/http"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(mux *http.ServeMux, handlers *Handlers) {
	// Autenticação e usuário
	mux.HandleFunc("/auth/register", handlers.User.RegisterHandler)
	mux.HandleFunc("/auth/login", handlers.User.LoginHandler)
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
}

// Handlers agrupa todos os handlers da aplicação
// Facilita a injeção de dependências

type Handlers struct {
	User        *UserHandler
	Theme       *ThemeHandler
	Question    *QuestionHandler
	Quiz        *QuizHandler
	History     *HistoryHandler
	Performance *PerformanceHandler
}
