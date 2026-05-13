package engine

import "time"

// Ramadan and Sacrifice feasts are tied to the Hijri calendar, so the upstream
// gem hard-codes start dates per year (2014–2020 only).
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
		}
		return dates[a.Year], nil
	})
}
