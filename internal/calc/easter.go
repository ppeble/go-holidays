package calc

import "time"

// Easter returns Gregorian Easter Sunday for the given year using the
// anonymous Gregorian algorithm (Butcher's algorithm).
func Easter(year int) time.Time {
	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// OrthodoxEasterJulian returns Eastern Orthodox Easter Sunday in the JULIAN
// calendar (Meeus's Julian algorithm). Used by churches that observe Easter
// on the unconverted Julian date.
func OrthodoxEasterJulian(year int) time.Time {
	a := year % 4
	b := year % 7
	c := year % 19
	d := (19*c + 15) % 30
	e := (2*a + 4*b - d + 34) % 7
	month := (d + e + 114) / 31
	day := ((d + e + 114) % 31) + 1
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// OrthodoxEaster returns Eastern Orthodox Easter Sunday in the GREGORIAN
// calendar. Computes the Julian date and shifts by the Julian-to-Gregorian
// offset for the year's century.
func OrthodoxEaster(year int) time.Time {
	j := OrthodoxEasterJulian(year)
	return j.AddDate(0, 0, julianToGregorianOffset(year))
}

// julianToGregorianOffset returns the number of days to add to a Julian-calendar
// date to obtain the equivalent Gregorian-calendar date for the given year.
// For 1900-2099 this is 13; the formula generalizes across centuries.
func julianToGregorianOffset(year int) int {
	century := year / 100
	return century - century/4 - 2
}
