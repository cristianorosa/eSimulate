package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"eSimulate/backend/infra"
	"eSimulate/backend/interfaces"
	"eSimulate/backend/usecase"
)

func main() {
	// Configuração do banco (ajuste conforme necessário)
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=esimulate sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	// Injeta dependências
	repo := &infra.UserRepositoryPG{DB: db}
	uc := &usecase.UserUsecase{Repo: repo}

	handlers := &interfaces.Handlers{
		User:        &interfaces.UserHandler{UC: uc},
		Theme:       &interfaces.ThemeHandler{},
		Question:    &interfaces.QuestionHandler{},
		Quiz:        &interfaces.QuizHandler{},
		History:     &interfaces.HistoryHandler{},
		Performance: &interfaces.PerformanceHandler{},
	}

	mux := http.NewServeMux()
	interfaces.SetupRoutes(mux, handlers)

	log.Println("Servidor iniciado em :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
