package handlers

import (
	"fmt"
	"net/http"
	"time"

	"go_main_final_project/pkg/repeat"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "invalid now format", http.StatusBadRequest)
		return
	}

	result, err := repeat.NextDate(now, dateStr, repeatStr)
	if err != nil {
		fmt.Fprint(w, "")
		return
	}

	fmt.Fprint(w, result)
}