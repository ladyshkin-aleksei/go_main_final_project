package api

import (
	"encoding/json"
	"net/http"

	"go_main_final_project/pkg/db"
)

func Init() error {
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/task", AuthMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", AuthMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", AuthMiddleware(doneHandler))
	return nil
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "ID not specified"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, task)
}

func handlePostTask(w http.ResponseWriter, r *http.Request) {
	addTaskHandler(w, r)
}

func handlePutTask(w http.ResponseWriter, r *http.Request) {
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

	err := db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]interface{}{})
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "ID not specified"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]interface{}{})
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTask(w, r)
	case http.MethodPost:
		handlePostTask(w, r)
	case http.MethodPut:
		handlePutTask(w, r)
	case http.MethodDelete:
		handleDeleteTask(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}