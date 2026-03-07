package repeat

import (
	"strconv"
	"strings"
	"time"
	"fmt"
)

func NextDate(now time.Time, dateStr, repeatStr string) (string, error) {
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", err
	}

	if repeatStr == "" {
		return dateStr, nil
	}

	parts := strings.Split(repeatStr, " ")
	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("некорректный формат правила 'y'")
		}

		next := date.AddDate(1, 0, 0)
		return next.Format("20060102"), nil

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("некорректный формат правила 'd'")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("некорректное количество дней в правиле 'd'")
		}

		next := date.AddDate(0, 0, days)
		return next.Format("20060102"), nil

	default:
		return "", fmt.Errorf("неизвестное правило повторения: %s", parts[0])
	}
}
