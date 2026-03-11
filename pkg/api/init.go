package api

import (
	"encoding/json"
	"net/http"

	"go_main_final_project/pkg/db"
	"go_main_final_project/pkg/models"
	"go_main_final_project/pkg/validation"
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
		writeJSON(w, errorResponse("ID not specified", http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err == db.ErrTaskNotFound {
			writeJSON(w, errorResponse("task not found", http.StatusNotFound), http.StatusNotFound)
		} else {
			writeJSON(w, errorResponse(err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}
		return
	}

	writeJSON(w, task, http.StatusOK)
}

func handlePostTask(w http.ResponseWriter, r *http.Request) {
	addTaskHandler(w, r)
}

func handlePutTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, errorResponse("error deserializing JSON", http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateTask(&task); err != nil {
		writeJSON(w, errorResponse(err.Error(), http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	err := db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, errorResponse(err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusNoContent)
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, errorResponse("ID not specified", http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		if err == db.ErrTaskNotFound {
			writeJSON(w, errorResponse("task not found", http.StatusNotFound), http.StatusNotFound)
	} else {
			writeJSON(w, errorResponse(err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusNoContent)
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
		writeJSON(w, errorResponse("method not allowed", http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}