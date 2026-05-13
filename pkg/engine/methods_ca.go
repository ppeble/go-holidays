package engine

import "time"

// ca_victoria_day is shared with tsx.yaml — register once here.
func init() {
	// Monday on or before May 24.
	RegisterMethod("ca_victoria_day", func(a MethodArgs) (time.Time, error) {
		d := time.Date(a.Year, time.May, 24, 0, 0, 0, 0, time.UTC)
		for d.Weekday() != time.Monday {
			d = d.AddDate(0, 0, -1)
		}
		return d, nil
	})
}
