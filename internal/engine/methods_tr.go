package engine

import "time"

// Ramadan and Sacrifice feasts are tied to the Hijri calendar, so proclaimed
// (Diyanet) start dates are hard-coded per year. The table is extended
// periodically; it currently runs from 2014 to 2030.
func init() {
	RegisterMethod("ramadan_feast", func(a MethodArgs) (time.Time, error) {
		dates := map[int]time.Time{
			2014: time.Date(2014, time.July, 28, 0, 0, 0, 0, time.UTC),
			2015: time.Date(2015, time.July, 17, 0, 0, 0, 0, time.UTC),
			2016: time.Date(2016, time.July, 5, 0, 0, 0, 0, time.UTC),
			2017: time.Date(2017, time.June, 25, 0, 0, 0, 0, time.UTC),
			2018: time.Date(2018, time.June, 15, 0, 0, 0, 0, time.UTC),
			2019: time.Date(2019, time.June, 4, 0, 0, 0, 0, time.UTC),
			2020: time.Date(2020, time.May, 24, 0, 0, 0, 0, time.UTC),
			2021: time.Date(2021, time.May, 13, 0, 0, 0, 0, time.UTC),
			2022: time.Date(2022, time.May, 2, 0, 0, 0, 0, time.UTC),
			2023: time.Date(2023, time.April, 21, 0, 0, 0, 0, time.UTC),
			2024: time.Date(2024, time.April, 10, 0, 0, 0, 0, time.UTC),
			2025: time.Date(2025, time.March, 30, 0, 0, 0, 0, time.UTC),
			2026: time.Date(2026, time.March, 20, 0, 0, 0, 0, time.UTC),
			2027: time.Date(2027, time.March, 9, 0, 0, 0, 0, time.UTC),
			2028: time.Date(2028, time.February, 26, 0, 0, 0, 0, time.UTC),
			2029: time.Date(2029, time.February, 15, 0, 0, 0, 0, time.UTC),
			2030: time.Date(2030, time.February, 4, 0, 0, 0, 0, time.UTC),
		}
		return dates[a.Year], nil
	})
	RegisterMethod("sacrifice_feast", func(a MethodArgs) (time.Time, error) {
		dates := map[int]time.Time{
			2014: time.Date(2014, time.October, 4, 0, 0, 0, 0, time.UTC),
			2015: time.Date(2015, time.September, 24, 0, 0, 0, 0, time.UTC),
			2016: time.Date(2016, time.September, 12, 0, 0, 0, 0, time.UTC),
			2017: time.Date(2017, time.September, 1, 0, 0, 0, 0, time.UTC),
			2018: time.Date(2018, time.August, 21, 0, 0, 0, 0, time.UTC),
			2019: time.Date(2019, time.August, 11, 0, 0, 0, 0, time.UTC),
			2020: time.Date(2020, time.July, 31, 0, 0, 0, 0, time.UTC),
			2021: time.Date(2021, time.July, 20, 0, 0, 0, 0, time.UTC),
			2022: time.Date(2022, time.July, 9, 0, 0, 0, 0, time.UTC),
			2023: time.Date(2023, time.June, 28, 0, 0, 0, 0, time.UTC),
			2024: time.Date(2024, time.June, 16, 0, 0, 0, 0, time.UTC),
			2025: time.Date(2025, time.June, 6, 0, 0, 0, 0, time.UTC),
			2026: time.Date(2026, time.May, 27, 0, 0, 0, 0, time.UTC),
			2027: time.Date(2027, time.May, 16, 0, 0, 0, 0, time.UTC),
			2028: time.Date(2028, time.May, 5, 0, 0, 0, 0, time.UTC),
			2029: time.Date(2029, time.April, 24, 0, 0, 0, 0, time.UTC),
			2030: time.Date(2030, time.April, 13, 0, 0, 0, 0, time.UTC),
		}
		return dates[a.Year], nil
	})
}
