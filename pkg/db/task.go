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
		return fmt.Errorf("дата представлена в формате, отличном от 20060102")
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("правило повторения указано в неправильном формате")
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

func NextDate(now time.Time, date, repeat string) (string, error) {
	return date, nil
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
		return nil, fmt.Errorf("ошибка выполнения запроса к БД: %w", err)
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки БД: %w", err)
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
			return nil, fmt.Errorf("Задача не найдена")
		}
		return nil, fmt.Errorf("ошибка получения задачи: %w", err)
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
		return fmt.Errorf("ошибка обновления задачи: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки количества изменённых записей: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

func ValidateTask(task *Task) error {
	if task.Title == "" {
		return fmt.Errorf("Не указан заголовок задачи")
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты: %s", task.Date)
	}

	if task.Repeat != "" {
		if !isValidRepeat(task.Repeat) {
			return fmt.Errorf("некорректное правило повторения: %s", task.Repeat)
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