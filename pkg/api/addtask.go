package api

import (
	"encoding/json"
	"net/http"
	"fmt"

	"go_main_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "error deserializing JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "the issue title is not specified"})
		return
	}

	if err := db.CheckDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}
