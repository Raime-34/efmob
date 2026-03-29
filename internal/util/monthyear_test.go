package util

import (
	"errors"
	"testing"
	"time"

	"efmob/internal/constants"

	"github.com/go-openapi/testify/v2/assert"
)

func TestMonthYearToTime(t *testing.T) {
	t.Parallel()

	got, err := MonthYearToTime("07-2025")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC), got)

	_, err = MonthYearToTime("7-2025")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, constants.ErrMonthYearMonthTwoDigits))

	_, err = MonthYearToTime("13-2025")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, constants.ErrMonthYearMonthOutOfRange))

	_, err = MonthYearToTime("2025-07")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, constants.ErrMonthYearMonthTwoDigits))

	_, err = MonthYearToTime("2005")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, constants.ErrMonthYearInvalidFormat))
}
