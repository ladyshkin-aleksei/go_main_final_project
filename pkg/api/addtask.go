package api

import (
	"encoding/json"
	"net/http"

	"go_main_final_project/pkg/db"
	"go_main_final_project/pkg/models" 
	"go_main_final_project/pkg/validation"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task 

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "error deserializing JSON"})
		return
	}

	if err := validation.ValidateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": task.ID})
}