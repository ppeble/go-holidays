package engine

import "time"

func init() {
	// Sweden: Midsummer Day — first Saturday on or after June 20.
	RegisterMethod("se_midsommardagen", func(a MethodArgs) (time.Time, error) {
		return nextSaturdayOnOrAfter(a.Year, time.June, 20), nil
	})
	// Sweden: All Saints' Day — first Saturday on or after October 31.
	RegisterMethod("se_alla_helgons_dag", func(a MethodArgs) (time.Time, error) {
		return nextSaturdayOnOrAfter(a.Year, time.October, 31), nil
	})
}

func nextSaturdayOnOrAfter(year int, month time.Month, day int) time.Time {
	d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return d.AddDate(0, 0, (6-int(d.Weekday())+7)%7)
}
