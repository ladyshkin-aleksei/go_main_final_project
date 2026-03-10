
package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"strings"
	"strconv"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX idx_scheduler_date ON scheduler (date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func CheckDate(task *Task) error {
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
		return nil
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}
	return nil
}

func AddTask(task *Task) error {
	if err := ValidateTask(task); err != nil {
		return err
	}

	today := time.Now().Format("20060102")
	if task.Date < today {
		task.Date = today
	}

	result, err := db.Exec(`
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = fmt.Sprintf("%d", id)
	return nil
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

	if err := CheckDate(task); err != nil {
		return err
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
	if len(parts) == 0 {
		return false
	}

	unit := parts[0]
	switch unit {
	case "y":
		return len(parts) == 1
	case "d":
		if len(parts) != 2 {
			return false
		}
		days, err := strconv.Atoi(parts[1])
		return err == nil && days > 0 && days <= 3650
	default:
		return false
	}
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

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}