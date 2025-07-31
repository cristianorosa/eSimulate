package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cristianorosa/eSimulate/backend/filesystem-api/infra"
	"github.com/cristianorosa/eSimulate/backend/filesystem-api/interfaces"
	"github.com/cristianorosa/eSimulate/backend/filesystem-api/usecase"
)

func main() {
	// Define o diretório base para operações (por segurança)
	baseDir := os.Getenv("FILESYSTEM_API_BASE_DIR")
	if baseDir == "" {
		// Se não definido, usa o diretório atual
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Erro ao obter diretório atual: %v", err)
		}
		baseDir = filepath.Join(currentDir, "sandbox")
	}
	
	// Cria o diretório sandbox se não existir
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Fatalf("Erro ao criar diretório base: %v", err)
	}
	
	log.Printf("API de Sistema de Arquivos iniciada")
	log.Printf("Diretório base: %s", baseDir)
	
	// Injeta dependências
	repo := infra.NewFileSystemRepositoryOS(baseDir)
	uc := usecase.NewFileSystemUsecase(repo)
	handler := interfaces.NewFileSystemHandler(uc)
	
	// Configura rotas
	mux := http.NewServeMux()
	interfaces.SetupFileSystemRoutes(mux, handler)
	
	// Middleware de CORS para desenvolvimento
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
	
	// Middleware de logging
	loggingHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	}
	
	// Aplica middlewares
	finalHandler := corsHandler(loggingHandler(mux))
	
	// Define a porta
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	
	log.Printf("Servidor iniciado na porta %s", port)
	log.Printf("Endpoints disponíveis:")
	log.Printf("  GET    /read-file/{path}     - Ler arquivo")
	log.Printf("  PUT    /create-file/{path}   - Criar arquivo")
	log.Printf("  PATCH  /edit-file/{path}     - Editar arquivo")
	log.Printf("  DELETE /delete-file/{path}   - Deletar arquivo")
	log.Printf("  GET    /read-dir/{path}      - Ler diretório")
	log.Printf("  PUT    /create-dir/{path}    - Criar diretório")
	log.Printf("  PATCH  /edit-dir/{path}      - Renomear diretório")
	log.Printf("  DELETE /delete-dir/{path}    - Deletar diretório")
	log.Printf("  POST   /execute-command/{cmd} - Executar comando")
	log.Printf("  GET    /status               - Status da API")
	
	log.Fatal(http.ListenAndServe(":"+port, finalHandler))
}