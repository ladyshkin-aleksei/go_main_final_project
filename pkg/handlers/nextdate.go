package handlers

import (
	"fmt"
	"net/http"
	"time"

	"go_main_final_project/pkg/repeat"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	if nowStr == "" {
		http.Error(w, "missing required parameter 'now'", http.StatusBadRequest)
		return
	}
	if dateStr == "" {
		http.Error(w, "missing required parameter 'date'", http.StatusBadRequest)
		return
	}
	if repeatStr == "" {
		http.Error(w, "missing required parameter 'repeat'", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "invalid 'now' format: expected YYYYMMDD", http.StatusBadRequest)
		return
	}

	result, err := repeat.NextDate(now, dateStr, repeatStr)
	if err != nil {
		fmt.Printf("NextDate error for date=%s, repeat=%s: %v\n", dateStr, repeatStr, err)

		fmt.Fprint(w, "")
		return
	}

	fmt.Fprint(w, result)
}
