package interfaces

import (
	"net/http"
	"strings"

	"github.com/cristianorosa/eSimulate/backend/infra"
)

// AuthMiddleware verifica o token JWT e protege rotas privadas.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Token não informado", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(token, "Bearer ")
		_, err := infra.ParseJWT(tokenStr)
		if err != nil {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}
		// Pode-se adicionar claims ao contexto se necessário
		r = r.WithContext(r.Context())
		next(w, r)
	}
}
