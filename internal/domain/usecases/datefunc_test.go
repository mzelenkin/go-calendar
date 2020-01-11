package usecases

import (
	"testing"
	"time"
)

func TestBod(t *testing.T) {
	now := time.Date(2020, 01, 01, 15, 30, 25, 100, time.Local)
	beginOfDate := bod(now)

	goodDate := time.Date(2020, 01, 01, 0, 0, 0, 0, time.Local)
	if !beginOfDate.Equal(goodDate) {
		t.Fail()
	}
}

func TestEod(t *testing.T) {
	now := time.Date(2020, 01, 01, 15, 30, 25, 100, time.Local)
	endOfDate := eod(now)
	// Конец дня, это на 1 наносекунду раньше начала следующего
	goodDate := time.Date(2020, 01, 02, 0, 0, 0, 0, time.Local).Add(-time.Nanosecond)
	if !endOfDate.Equal(goodDate) {
		t.Fail()
	}
}

func TestWeekRange(t *testing.T) {
	now := time.Date(2020, 01, 01, 15, 30, 25, 100, time.Local)
	beginOfWeek, endOfWeek := weekRange(now)

	// Начало недели 31 декабря 2019 00:00:00
	// Конец недели 5 января 2020 23:59:59
	goodBeginDate := time.Date(2019, 12, 30, 0, 0, 0, 0, time.Local)
	goodEndDate := time.Date(2020, 01, 6, 0, 0, 0, 0, time.Local).Add(-time.Nanosecond)
	if !beginOfWeek.Equal(goodBeginDate) {
		t.Fail()
	}

	if !endOfWeek.Equal(goodEndDate) {
		t.Fail()
	}
}
