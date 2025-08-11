package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/cristianorosa/eSimulate/backend/domain"
	"github.com/cristianorosa/eSimulate/backend/infra"
	"github.com/cristianorosa/eSimulate/backend/interfaces"
	"github.com/cristianorosa/eSimulate/backend/usecase"
)

func main() {
	// Configuração do banco (ajuste conforme necessário)
	dsn := "host=localhost port=5432 user=esimulate password=DbaInv=2025 dbname=esimulate sslmode=disable"
	log.Printf("Tentando conectar ao banco com DSN: %s", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}

	// Testar a conexão
	err = db.Ping()
	if err != nil {
		log.Fatalf("Erro ao fazer ping no banco: %v", err)
	}
	log.Printf("Conexão com banco estabelecida com sucesso")

	defer db.Close()

	// Injeta dependências
	userRepo := &infra.UserRepositoryPG{DB: db}
	userUC := &usecase.UserUsecase{Repo: userRepo}

	// Repositories e Usecases para History e Performance
	userQuizRepo := &infra.UserQuizRepositoryPG{DB: db}
	userQuizUC := &usecase.UserQuizUsecase{Repo: userQuizRepo}

	performanceRepo := &infra.PerformanceRepositoryPG{DB: db}
	performanceUC := &usecase.PerformanceUsecase{Repo: performanceRepo}

	// Novos repositories e usecases para v2
	areaRepo := &infra.AreaRepositoryPG{DB: db}
	areaUC := &usecase.AreaUsecase{Repo: areaRepo}

	examRepo := &infra.ExamRepositoryPG{DB: db}
	examUC := &usecase.ExamUsecase{Repo: examRepo}

	topicRepo := &infra.TopicRepositoryPG{DB: db}
	topicUC := &usecase.TopicUsecase{Repo: topicRepo}

	questionRepo := &infra.QuestionRepositoryPG{DB: db}
	optionRepo := &infra.OptionRepositoryPG{DB: db}
	questionUC := &usecase.QuestionUsecase{Repo: questionRepo, OptionRepo: optionRepo}
	log.Printf("QuestionUsecase inicializado: %+v", questionUC)

	quizRepo := &infra.QuizRepositoryPG{DB: db}
	quizUC := &usecase.QuizUsecase{Repo: quizRepo}
	log.Printf("QuizUsecase inicializado: %+v", quizUC)

	userExamRepo := &infra.UserExamRepositoryPG{DB: db}
	userAnswerRepo := &infra.UserAnswerRepositoryPG{DB: db}
	topicPerformanceRepo := &infra.TopicPerformanceRepositoryPG{DB: db}
	userExamUC := &usecase.UserExamUsecase{
		Repo:            userExamRepo,
		AnswerRepo:      userAnswerRepo,
		PerformanceRepo: topicPerformanceRepo,
	}

	// Cria usuários iniciais se não existirem
	createInitialUsers(userUC)

	// Theme usecase
	themeRepo := &infra.ThemeRepositoryPG{DB: db}
	themeUC := &usecase.ThemeUsecase{Repo: themeRepo}

	handlers := &interfaces.Handlers{
		User:        &interfaces.UserHandler{UC: userUC},
		Theme:       &interfaces.ThemeHandler{UC: themeUC},
		Question:    &interfaces.QuestionHandler{UC: questionUC, AreaUC: areaUC, ExamUC: examUC, TopicUC: topicUC},
		Quiz:        &interfaces.QuizHandler{UC: quizUC},
		History:     &interfaces.HistoryHandler{UC: userQuizUC},
		Performance: &interfaces.PerformanceHandler{UC: performanceUC},
		Area:        &interfaces.AreaHandler{UC: areaUC},
		Exam:        &interfaces.ExamHandler{UC: examUC},
		Topic:       &interfaces.TopicHandler{UC: topicUC},
		UserExam:    &interfaces.UserExamHandler{UC: userExamUC},
	}
	log.Printf("QuestionHandler inicializado: %+v", handlers.Question)

	mux := http.NewServeMux()
	interfaces.SetupRoutes(mux, handlers)

	log.Println("Servidor iniciado em :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// createInitialUsers cria os usuários iniciais do sistema
func createInitialUsers(userUC *usecase.UserUsecase) {
	users := []struct {
		name     string
		email    string
		password string
		roleID   int
	}{
		{"Administrador", "admin@esimulate.com", "password", domain.RoleIDAdmin},
		{"Redator", "redator@esimulate.com", "password", domain.RoleIDRedator},
		{"Usuário", "user@esimulate.com", "password", domain.RoleIDUser},
	}

	for _, u := range users {
		_, err := userUC.Repo.FindByEmail(u.email)
		if err != nil {
			if err == sql.ErrNoRows {
				// Usuário não encontrado, criar
				user := &domain.User{
					Name:         u.name,
					Email:        u.email,
					PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // hash de "password"
					RoleID:       u.roleID,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				if err := userUC.Repo.Create(user); err != nil {
					log.Printf("Erro ao criar usuário %s: %v", u.name, err)
				} else {
					log.Printf("Usuário %s criado com sucesso", u.name)
				}
			} else {
				log.Printf("Erro ao verificar usuário %s: %v", u.name, err)
			}
		} else {
			log.Printf("Usuário %s já existe", u.name)
		}
	}
}
