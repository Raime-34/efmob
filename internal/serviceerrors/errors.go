package serviceerrors

import (
	"errors"
)

var (
	ErrSubscriptionNotFound = errors.New("Переданная комбианция пользователь-подписка не найдена")

	// MonthYear (формат MM-YYYY для util.MonthYearToTime).
	ErrMonthYearInvalidFormat   = errors.New("ожидается формат даты MM-YYYY (месяц-год через дефис)")
	ErrMonthYearMonthTwoDigits  = errors.New("месяц должен быть из двух цифр (01–12)")
	ErrMonthYearInvalidMonthNum = errors.New("месяц не является числом")
	ErrMonthYearInvalidYearNum  = errors.New("год не является числом")
	ErrMonthYearMonthOutOfRange = errors.New("месяц вне диапазона 1–12")
)

// IsMonthYearError — любая ошибка разбора MM-YYYY
func IsMonthYearError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrMonthYearInvalidFormat) ||
		errors.Is(err, ErrMonthYearMonthTwoDigits) ||
		errors.Is(err, ErrMonthYearInvalidMonthNum) ||
		errors.Is(err, ErrMonthYearInvalidYearNum) ||
		errors.Is(err, ErrMonthYearMonthOutOfRange)
}
