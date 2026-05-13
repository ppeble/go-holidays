package engine

import "time"

func init() {
	// Wednesday strictly before November 23 (Buß- und Bettag, Germany).
	// If Nov 23 is Wednesday, returns Nov 16.
	RegisterMethod("de_buss_und_bettag", func(a MethodArgs) (time.Time, error) {
		d := time.Date(a.Year, time.November, 23, 0, 0, 0, 0, time.UTC)
		w := int(d.Weekday())
		if w > 3 {
			d = d.AddDate(0, 0, -(w - 3))
		} else {
			d = d.AddDate(0, 0, -(w + 4))
		}
		return d, nil
	})
}
