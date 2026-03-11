package api

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func errorResponse(message string, code int) map[string]interface{} {
	return map[string]interface{}{
		"error": message,
		"code":  code,
	}
}