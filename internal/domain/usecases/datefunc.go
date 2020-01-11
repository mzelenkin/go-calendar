package usecases

import "time"

// bod функция, возвращающая начало дня
func bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// eod функция, возвращающая конец дня
func eod(t time.Time) time.Time {
	return bod(t).AddDate(0, 0, 1).Add(-time.Nanosecond)
}

// weekStart находит день начала недели
func weekStart(day time.Time) time.Time {
	year, week := day.ISOWeek()
	loc := day.Location()

	// Начинаем с середины года
	t := time.Date(year, 7, 1, 0, 0, 0, 0, loc)

	// Откатываемся к понедельнику
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Разница в неделях
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

// WeekRange Выдает диапазон дат начала и конца недели по году year и номеру недели week
func weekRange(day time.Time) (start, end time.Time) {
	start = weekStart(day)
	end = start.AddDate(0, 0, 7).Add(-time.Nanosecond)

	return
}

// monthStart находит начало месяца по дате t
func monthStart(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// monthRange выдает диапазон дат начала и конца месяца по дате t
func monthRange(t time.Time) (start, end time.Time) {
	year, month, _ := t.Date()

	start = time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return
}
