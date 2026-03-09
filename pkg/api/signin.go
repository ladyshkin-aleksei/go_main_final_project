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
		writeJSON(w, SigninResponse{Error: "the method is not supported"})
		return
	}

	var req SigninRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, SigninResponse{Error: "invalid JSON format"})
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		writeJSON(w, SigninResponse{Token: "no-auth-required"})
		return
	}

	if req.Password != expectedPassword {
		writeJSON(w, SigninResponse{Error: "invalid password"})
		return
	}

	token, err := auth.GenerateToken(req.Password)
	if err != nil {
		writeJSON(w, SigninResponse{Error: "token generation error"})
		return
	}

	writeJSON(w, SigninResponse{Token: token})
}
