package middleware

import (
	"URL_SHORT/internal/interfaces/database/postgres"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type AuthMiddleware struct {
	URLRepo *postgres.URLRepository
}

func NewAuthMiddleware(urlRepo *postgres.URLRepository) *AuthMiddleware {
	return &AuthMiddleware{URLRepo: urlRepo}
}

func (amw *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		query := "SELECT password FROM users WHERE username = $1"
		row := amw.URLRepo.DB.QueryRow(query, username)
		var storedPassword []byte
		err := row.Scan(&storedPassword)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword(storedPassword, []byte(password))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
