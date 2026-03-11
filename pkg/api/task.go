package api

import (
	"net/http"

	"go_main_final_project/pkg/db"
	"go_main_final_project/pkg/models"
)

const DefaultTasksLimit = 50

type TasksResp struct {
	Tasks []*models.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, errorResponse("the method is not supported", http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(DefaultTasksLimit)
	if err != nil {
		writeJSON(w, errorResponse(err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	}, http.StatusOK)
}
