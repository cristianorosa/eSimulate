package interfaces

import (
	"log"
	"net/http"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/infra"
)

// AuthMiddleware verifica o token JWT e protege rotas privadas.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AuthMiddleware: Verificando requisição para %s", r.URL.Path)

		token := r.Header.Get("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			log.Printf("AuthMiddleware: Token não informado para %s", r.URL.Path)
			http.Error(w, "Token não informado", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(token, "Bearer ")
		claims, err := infra.ParseJWT(tokenStr)
		if err != nil {
			log.Printf("AuthMiddleware: Token inválido para %s: %v", r.URL.Path, err)
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		log.Printf("AuthMiddleware: Token válido para usuário %v em %s", claims, r.URL.Path)
		// Pode-se adicionar claims ao contexto se necessário
		r = r.WithContext(r.Context())
		next(w, r)
	}
}
