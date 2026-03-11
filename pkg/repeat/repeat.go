package repeat

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, currentDate, repeat string) (string, error) {
	current, err := time.Parse(DateFormat, currentDate)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %s", currentDate)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("empty repeat rule")
	}

	unit := parts[0]

	switch unit {
	case "y":
		next := current.AddDate(1, 0, 0)
		return next.Format(DateFormat), nil

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid repeat format for days: %s. Expected 'd N'", repeat)
	}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid number of days: %s", parts[1])
	}
		if days <= 0 {
			return "", fmt.Errorf("number of days must be positive: %d", days)
	}
		if days > 3650 {
			return "", fmt.Errorf("number of days too large: %d (max 3650)", days)
	}
		next := current.AddDate(0, 0, days)
		return next.Format(DateFormat), nil

	default:
		return "", fmt.Errorf("unknown repeat unit: %s. Supported: 'y', 'd N'", unit)
	}
}