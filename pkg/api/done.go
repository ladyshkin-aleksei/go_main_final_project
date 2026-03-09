package api

import (
	"net/http"
	"time"

	"go_main_final_project/pkg/db"
	"go_main_final_project/pkg/repeat"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "the method is not supported"})
		return
	}

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

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
	}
	} else {
		now := time.Now()
		nextDate, err := repeat.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}

		task.Date = nextDate
		err = db.UpdateTask(task)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	}

	writeJSON(w, map[string]interface{}{})
}
