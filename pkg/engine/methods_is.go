package engine

import "time"

func init() {
	// First day of summer (Iceland): Thursday strictly after April 18.
	// If April 18 is a Thursday, returns April 25.
	RegisterMethod("is_sumardagurinn_fyrsti", func(a MethodArgs) (time.Time, error) {
		d := time.Date(a.Year, time.April, 18, 0, 0, 0, 0, time.UTC)
		w := int(d.Weekday())
		if w < 4 {
			d = d.AddDate(0, 0, 4-w)
		} else {
			d = d.AddDate(0, 0, 11-w)
		}
		return d, nil
	})
}
