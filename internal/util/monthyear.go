// Пакет util — общие вспомогательные функции приложения.
package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"efmob/internal/constants"
)

// MonthYearToTime разбирает строку формата «месяц-год» (MM-YYYY), например "07-2025".
// Месяц — строго две цифры (01–12). День всегда 1-е число; время — 00:00 UTC.
func MonthYearToTime(s string) (time.Time, error) {
	parts := strings.Split(strings.TrimSpace(s), "-")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("%w: получено %q", constants.ErrMonthYearInvalidFormat, s)
	}

	if len(parts[0]) != 2 {
		return time.Time{}, fmt.Errorf("%w: получено %q", constants.ErrMonthYearMonthTwoDigits, parts[0])
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %w", constants.ErrMonthYearInvalidMonthNum, err)
	}
	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %w", constants.ErrMonthYearInvalidYearNum, err)
	}

	if month < 1 || month > 12 {
		return time.Time{}, fmt.Errorf("%w: %d", constants.ErrMonthYearMonthOutOfRange, month)
	}

	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}

// TimeToMonthYear форматирует дату в строку MM-YYYY (как в MonthYearToTime).
func TimeToMonthYear(t time.Time) string {
	return fmt.Sprintf("%02d-%04d", int(t.Month()), t.Year())
}
