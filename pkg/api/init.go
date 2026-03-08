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


func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	
case http.MethodDelete:

	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]interface{}{})
	
	case http.MethodGet:
		id := r.FormValue("id")
		if id == "" {
			writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
			return
		}

		task, err := db.GetTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
	}

		writeJSON(w, task)

	case http.MethodPut:
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
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
	case http.MethodPost:
		addTaskHandler(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}