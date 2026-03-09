package api

import (
	"net/http"
	"time"

	"go_main_final_project/pkg/repeat"
)

const (
	DateFormat = "20060102"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "incorrect format of the now parameter", http.StatusBadRequest)
			return
	}
	}

	if dateStr == "" {
		http.Error(w, "the required date parameter is missing", http.StatusBadRequest)
		return
	}
	if repeatStr == "" {
		http.Error(w, "the required repeat parameter is missing
", http.StatusBadRequest)
		return
	}

	nextDate, err := repeat.NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

func RegisterNextDateHandler() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
}
