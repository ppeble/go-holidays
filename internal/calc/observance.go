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

// ToTheWeekdayAfter returns the day after date, treating Sunday as a non-working
// day to skip: to_monday_if_sunday(to_monday_if_sunday(date) + 1). Saturday is
// not skipped.
func ToTheWeekdayAfter(date time.Time) time.Time {
	d := date
	if d.Weekday() == time.Sunday {
		d = d.AddDate(0, 0, 1)
	}
	d = d.AddDate(0, 0, 1)
	if d.Weekday() == time.Sunday {
		d = d.AddDate(0, 0, 1)
	}
	return d
}

// ToTheSecondWeekdayAfter is to_monday_if_sunday(to_the_weekday_after(date) + 1).
func ToTheSecondWeekdayAfter(date time.Time) time.Time {
	d := ToTheWeekdayAfter(date).AddDate(0, 0, 1)
	if d.Weekday() == time.Sunday {
		d = d.AddDate(0, 0, 1)
	}
	return d
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
