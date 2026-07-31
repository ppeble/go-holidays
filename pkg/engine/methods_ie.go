package engine

import "time"

func init() {
	// Ireland: St Brigid's Day is the first Monday in February, except when
	// February 1 is itself a Friday, in which case it falls on that Friday.
	RegisterMethod("ie_st_brigids_day", func(a MethodArgs) (time.Time, error) {
		d := time.Date(a.Year, time.February, 1, 0, 0, 0, 0, time.UTC)
		if d.Weekday() == time.Friday {
			return d, nil
		}
		return d.AddDate(0, 0, (1-int(d.Weekday())+7)%7), nil
	})
}
