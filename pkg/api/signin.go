package api

import (
	"encoding/json"
	"net/http"

	"go_main_final_project/pkg/auth"
	"go_main_final_project/pkg/config"
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
		writeJSON(w, SigninResponse{Error: "the method is not supported"}, http.StatusMethodNotAllowed)
		return
	}

	var req SigninRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, SigninResponse{Error: "invalid JSON format"}, http.StatusBadRequest)
		return
	}

	cfg := config.Get()

	if cfg.Password == "" {
		writeJSON(w, SigninResponse{Token: "no-auth-required"}, http.StatusOK)
		return
	}

	if req.Password != cfg.Password {
		writeJSON(w, SigninResponse{Error: "invalid password"}, http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(req.Password)
	if err != nil {
		writeJSON(w, SigninResponse{Error: "token generation error"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, SigninResponse{Token: token}, http.StatusOK)
}