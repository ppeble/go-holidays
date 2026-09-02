package calc

import "time"

// DayOfMonth returns the date of the Nth occurrence of wday in (year, month).
// Pass week=-1 for the last occurrence; week=1..5 for the Nth from the start.
// If no such date exists in the month (e.g. 5th Monday when there are only 4),
// returns the zero time.
func DayOfMonth(year int, month time.Month, week int, wday time.Weekday) time.Time {
	if week == -1 {
		return lastWeekdayOfMonth(year, month, wday)
	}
	return nthWeekdayOfMonth(year, month, week, wday)
}

func nthWeekdayOfMonth(year int, month time.Month, n int, wday time.Weekday) time.Time {
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	offset := (int(wday) - int(first.Weekday()) + 7) % 7
	day := 1 + offset + (n-1)*7
	candidate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if candidate.Month() != month {
		return time.Time{}
	}
	return candidate
}

func lastWeekdayOfMonth(year int, month time.Month, wday time.Weekday) time.Time {
	firstNext := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstNext.AddDate(0, 0, -1)
	offset := (int(lastDay.Weekday()) - int(wday) + 7) % 7
	return lastDay.AddDate(0, 0, -offset)
}
