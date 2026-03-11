package api

import (
	"net/http"
	"os"

	"go_main_final_project/pkg/auth"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")

	if password != "" {
		cookie, err := r.Cookie("token")
		if err != nil {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}

	if !auth.ValidateToken(cookie.Value, password) {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	}

	next(w, r)
	}
}
