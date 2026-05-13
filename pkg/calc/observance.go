package calc

import "time"

// ToMondayIfSunday returns the next Monday if date is a Sunday; otherwise date.
func ToMondayIfSunday(date time.Time) time.Time {
	if date.Weekday() == time.Sunday {
		return date.AddDate(0, 0, 1)
	}
	return date
}

// ToMondayIfWeekend rolls Saturday/Sunday observations forward to Monday.
func ToMondayIfWeekend(date time.Time) time.Time {
	switch date.Weekday() {
	case time.Saturday:
		return date.AddDate(0, 0, 2)
	case time.Sunday:
		return date.AddDate(0, 0, 1)
	}
	return date
}

// ToWeekdayIfWeekend rolls Saturday back to Friday and Sunday forward to Monday.
func ToWeekdayIfWeekend(date time.Time) time.Time {
	switch date.Weekday() {
	case time.Saturday:
		return date.AddDate(0, 0, -1)
	case time.Sunday:
		return date.AddDate(0, 0, 1)
	}
	return date
}

// ToWeekdayIfBoxingWeekend handles Boxing Day observance. Saturday and Sunday
// roll forward by 2 days (giving Monday and Tuesday respectively). Monday rolls
// forward by 1 day (Tuesday) — used when Christmas Day observance has already
// landed on the Monday, pushing Boxing Day's observance one further.
func ToWeekdayIfBoxingWeekend(date time.Time) time.Time {
	switch date.Weekday() {
	case time.Saturday, time.Sunday:
		return date.AddDate(0, 0, 2)
	case time.Monday:
		return date.AddDate(0, 0, 1)
	}
	return date
}

// ToTuesdayIfSundayOrMondayIfSaturday observes Saturday -> Monday, Sunday -> Tuesday.
// Used where Saturday rolls to Monday but Monday is the actual holiday, pushing Sunday to Tuesday.
func ToTuesdayIfSundayOrMondayIfSaturday(date time.Time) time.Time {
	switch date.Weekday() {
	case time.Saturday:
		return date.AddDate(0, 0, 2)
	case time.Sunday:
		return date.AddDate(0, 0, 2)
	}
	return date
}

// ToTheWeekdayAfter returns the next non-weekend day strictly after date.
func ToTheWeekdayAfter(date time.Time) time.Time {
	d := date.AddDate(0, 0, 1)
	for d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
		d = d.AddDate(0, 0, 1)
	}
	return d
}

// ToTheSecondWeekdayAfter returns the second non-weekend day strictly after date.
func ToTheSecondWeekdayAfter(date time.Time) time.Time {
	return ToTheWeekdayAfter(ToTheWeekdayAfter(date))
}

// ToPreviousDayIfLeapYear shifts date back one day in leap years; otherwise unchanged.
func ToPreviousDayIfLeapYear(date time.Time) time.Time {
	if IsLeapYear(date.Year()) {
		return date.AddDate(0, 0, -1)
	}
	return date
}

// IsLeapYear reports whether year is a Gregorian leap year.
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
