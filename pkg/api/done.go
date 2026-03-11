package api

import (
	"net/http"
	"time"

	"go_main_final_project/pkg/db"
	"go_main_final_project/pkg/repeat"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, errorResponse("the method is not supported", http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

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

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			if err == db.ErrTaskNotFound {
				writeJSON(w, errorResponse("task not found", http.StatusNotFound), http.StatusNotFound)
		} else {
			writeJSON(w, errorResponse(err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
		}
			return
		}
	} else {
		now := time.Now()
		nextDate, err := repeat.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, errorResponse("invalid repeat pattern: "+err.Error(), http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		task.Date = nextDate
		err = db.UpdateTask(task)
		if err != nil {
			writeJSON(w, errorResponse(err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, map[string]interface{}{}, http.StatusNoContent)
}
