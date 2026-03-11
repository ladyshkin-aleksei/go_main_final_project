package validation

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go_main_final_project/pkg/models"
)

func CheckDate(task *models.Task) error {
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

func ValidateTask(task *models.Task) error {
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
