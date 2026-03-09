package db

import (
	"database/sql"
	"fmt"
	"time"
	"strconv"
	"strings"
)

type Task struct {
	ID      string `json:"id"`  
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func CheckDate(task *Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
		return nil
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("the date is presented in a format other than 20060102")
	}

	if task.Repeat != "" {
		next, err := NextDate(task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("the repetition rule is in the wrong format")
		}
	if afterNow(now, t) {
		task.Date = next
	}
	} else {
		if afterNow(now, t) {
			task.Date = now.Format("20060102")
	}
	}

	return nil
}

func NextDate(date, repeat string) (string, error) {
	t, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("incorrect date format: %s", date)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("incorrect repetition rule: %s", repeat)
	}

	unit := parts[0]
	count, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("incorrect number in the repetition rule: %s", parts[1])
	}

	var next time.Time
	switch unit {
	case "d":
		next = t.AddDate(0, 0, count)
	case "w":
		next = t.AddDate(0, 0, count*7)
	case "m":
		next = t.AddDate(0, count, 0)
	case "y":
		next = t.AddDate(count, 0, 0)
	default:
		return "", fmt.Errorf("unknown type of repetition: %s", unit)
	}

	return next.Format("20060102"), nil
}

func afterNow(now, t time.Time) bool {
	return t.Before(now) 
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func Tasks(limit int) ([]*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date ASC
		LIMIT ?
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("error in executing a database request: %w", err)
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error scanning a database row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`

	task := &Task{}
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("issue not found")
		}
		return nil, fmt.Errorf("issue receipt error: %w", err)
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	if err := ValidateTask(task); err != nil {
		return err
	}

	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("issue update error: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking the number of modified records: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("issue not found")
	}

	return nil
}

func ValidateTask(task *Task) error {
	if task.Title == "" {
		return fmt.Errorf("the issue title is not specified")
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("incorrect date format: %s", task.Date)
	}

	if task.Repeat != "" {
		if !isValidRepeat(task.Repeat) {
			return fmt.Errorf("incorrect repetition rule: %s", task.Repeat)
		}
	}

	return nil
}

func isValidRepeat(repeat string) bool {
	parts := strings.Split(repeat, " ")
	if len(parts) != 2 {
		return false
	}

	unit := parts[0]
	validUnits := map[string]bool{"d": true, "w": true, "m": true, "y": true}
	if !validUnits[unit] {
		return false
	}

	_, err := strconv.Atoi(parts[1])
	return err == nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("issue deletion error: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking the number of deleted records: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("issue not found")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := db.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("issue date update error: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking the number of modified records: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("issue not found")
	}

	return nil
}