package engine

import "time"

func init() {
	// New Zealand: roll a date to the closest Monday. Mon..Thu roll backwards
	// to the prior Monday; Fri/Sat/Sun roll forward to the next Monday.
	RegisterMethod("closest_monday", func(a MethodArgs) (time.Time, error) {
		d := a.Date
		switch w := int(d.Weekday()); w {
		case 0:
			return d.AddDate(0, 0, 1), nil
		case 1, 2, 3, 4:
			return d.AddDate(0, 0, -(w - 1)), nil
		default:
			return d.AddDate(0, 0, 8-w), nil
		}
	})
	RegisterMethod("previous_friday", func(a MethodArgs) (time.Time, error) {
		return a.Date.AddDate(0, 0, -3), nil
	})
	RegisterMethod("next_week", func(a MethodArgs) (time.Time, error) {
		return a.Date.AddDate(0, 0, 7), nil
	})
	RegisterMethod("nz_canterbury_anniversary", func(a MethodArgs) (time.Time, error) {
		d := time.Date(a.Year, time.November, 1, 0, 0, 0, 0, time.UTC)
		for d.Weekday() != time.Tuesday {
			d = d.AddDate(0, 0, 1)
		}
		d = d.AddDate(0, 0, 1)
		for d.Weekday() != time.Friday {
			d = d.AddDate(0, 0, 1)
		}
		return d.AddDate(0, 0, 7), nil
	})
	// Matariki is set in New Zealand legislation each year (lunar/astronomical).
	// The matarikiDates table lists the legislated dates verbatim.
	RegisterMethod("matariki", func(a MethodArgs) (time.Time, error) {
		d, ok := matarikiDates[a.Year]
		if !ok {
			return time.Time{}, nil
		}
		return d, nil
	})
}

var matarikiDates = map[int]time.Time{
	2022: time.Date(2022, time.June, 24, 0, 0, 0, 0, time.UTC),
	2023: time.Date(2023, time.July, 14, 0, 0, 0, 0, time.UTC),
	2024: time.Date(2024, time.June, 28, 0, 0, 0, 0, time.UTC),
	2025: time.Date(2025, time.June, 20, 0, 0, 0, 0, time.UTC),
	2026: time.Date(2026, time.July, 10, 0, 0, 0, 0, time.UTC),
	2027: time.Date(2027, time.June, 25, 0, 0, 0, 0, time.UTC),
	2028: time.Date(2028, time.July, 14, 0, 0, 0, 0, time.UTC),
	2029: time.Date(2029, time.July, 6, 0, 0, 0, 0, time.UTC),
	2030: time.Date(2030, time.June, 21, 0, 0, 0, 0, time.UTC),
	2031: time.Date(2031, time.July, 11, 0, 0, 0, 0, time.UTC),
	2032: time.Date(2032, time.July, 2, 0, 0, 0, 0, time.UTC),
	2033: time.Date(2033, time.June, 24, 0, 0, 0, 0, time.UTC),
	2034: time.Date(2034, time.July, 7, 0, 0, 0, 0, time.UTC),
	2035: time.Date(2035, time.June, 29, 0, 0, 0, 0, time.UTC),
	2036: time.Date(2036, time.July, 18, 0, 0, 0, 0, time.UTC),
	2037: time.Date(2037, time.July, 10, 0, 0, 0, 0, time.UTC),
	2038: time.Date(2038, time.June, 25, 0, 0, 0, 0, time.UTC),
	2039: time.Date(2039, time.July, 15, 0, 0, 0, 0, time.UTC),
	2040: time.Date(2040, time.July, 6, 0, 0, 0, 0, time.UTC),
	2041: time.Date(2041, time.July, 19, 0, 0, 0, 0, time.UTC),
	2042: time.Date(2042, time.July, 11, 0, 0, 0, 0, time.UTC),
	2043: time.Date(2043, time.July, 3, 0, 0, 0, 0, time.UTC),
	2044: time.Date(2044, time.June, 24, 0, 0, 0, 0, time.UTC),
	2045: time.Date(2045, time.July, 7, 0, 0, 0, 0, time.UTC),
	2046: time.Date(2046, time.June, 29, 0, 0, 0, 0, time.UTC),
	2047: time.Date(2047, time.July, 19, 0, 0, 0, 0, time.UTC),
	2048: time.Date(2048, time.July, 3, 0, 0, 0, 0, time.UTC),
	2049: time.Date(2049, time.June, 25, 0, 0, 0, 0, time.UTC),
	2050: time.Date(2050, time.July, 15, 0, 0, 0, 0, time.UTC),
	2051: time.Date(2051, time.June, 30, 0, 0, 0, 0, time.UTC),
	2052: time.Date(2052, time.June, 21, 0, 0, 0, 0, time.UTC),
}
