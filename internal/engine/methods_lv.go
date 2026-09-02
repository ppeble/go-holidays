package engine

import "time"

func init() {
	// Latvian Song and Dance Festival end date. The Latvian government
	// announces specific dates in advance; only the known years are filled in.
	RegisterMethod("lv_song_and_dance_festival_end_date", func(a MethodArgs) (time.Time, error) {
		switch a.Year {
		case 2018:
			return time.Date(2018, time.July, 8, 0, 0, 0, 0, time.UTC), nil
		case 2023:
			return time.Date(2023, time.July, 9, 0, 0, 0, 0, time.UTC), nil
		}
		return time.Time{}, nil
	})
}
