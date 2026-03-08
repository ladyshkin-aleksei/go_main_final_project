package api

import (
	"encoding/json"
	"net/http"
	"os"

	"go_main_final_project/pkg/auth"
)

type SigninRequest struct {
	Password string `json:"password"`
}

type SigninResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, SigninResponse{Error: "Метод не поддерживается"})
		return
	}

	var req SigninRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, SigninResponse{Error: "Неверный формат JSON"})
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		writeJSON(w, SigninResponse{Token: "no-auth-required"})
		return
	}

	if req.Password != expectedPassword {
		writeJSON(w, SigninResponse{Error: "Неверный пароль"})
		return
	}

	token, err := auth.GenerateToken(req.Password)
	if err != nil {
		writeJSON(w, SigninResponse{Error: "Ошибка генерации токена"})
		return
	}

	writeJSON(w, SigninResponse{Token: token})
}
