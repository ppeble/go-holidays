package engine

import "time"

func init() {
	// Argentina: roll a date to the nearest Monday. Tuesday/Wednesday roll back,
	// Thursday/Friday roll forward. Saturday, Sunday, and Monday stay put.
	RegisterMethod("to_nearest_monday", func(a MethodArgs) (time.Time, error) {
		switch a.Date.Weekday() {
		case time.Tuesday:
			return a.Date.AddDate(0, 0, -1), nil
		case time.Wednesday:
			return a.Date.AddDate(0, 0, -2), nil
		case time.Thursday:
			return a.Date.AddDate(0, 0, 4), nil
		case time.Friday:
			return a.Date.AddDate(0, 0, 3), nil
		}
		return a.Date, nil
	})
}
