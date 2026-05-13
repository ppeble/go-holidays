package engine

import "time"

func init() {
	// Finland: Midsummer Eve — Friday between June 19 and 25. If June 19 is
	// a Saturday, returns the Friday after (June 25).
	RegisterMethod("fi_juhannusaatto", func(a MethodArgs) (time.Time, error) {
		d := time.Date(a.Year, time.June, 19, 0, 0, 0, 0, time.UTC)
		w := int(d.Weekday())
		if w > 5 {
			return d.AddDate(0, 0, 6), nil
		}
		return d.AddDate(0, 0, 5-w), nil
	})
	// Finland: Midsummer Day — first Saturday on or after June 20.
	RegisterMethod("fi_juhannuspaiva", func(a MethodArgs) (time.Time, error) {
		return nextSaturdayOnOrAfter(a.Year, time.June, 20), nil
	})
	// Finland: All Saints' Day — first Saturday on or after October 31.
	RegisterMethod("fi_pyhainpaiva", func(a MethodArgs) (time.Time, error) {
		return nextSaturdayOnOrAfter(a.Year, time.October, 31), nil
	})
}
