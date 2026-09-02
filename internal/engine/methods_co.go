package engine

import "time"

// Colombia rolls many fixed-date Catholic holidays to the following Monday
// when they don't already fall on a Monday. Each named method targets a
// specific anchor date and shares the same rolling logic.
func init() {
	RegisterMethod("to_following_monday_if_not_monday", func(a MethodArgs) (time.Time, error) {
		return rollToFollowingMonday(a.Date), nil
	})
	RegisterMethod("epiphany", coAnchor(time.January, 6))
	RegisterMethod("saint_josephs_day", coAnchor(time.March, 19))
	RegisterMethod("saint_peter_and_saint_paul", coAnchor(time.June, 29))
	RegisterMethod("assumption_of_mary", coAnchor(time.August, 15))
	RegisterMethod("columbus_day", coAnchor(time.October, 12))
	RegisterMethod("all_saints_day", coAnchor(time.November, 1))
	RegisterMethod("independence_of_cartagena", coAnchor(time.November, 11))
}

func coAnchor(month time.Month, day int) Method {
	return func(a MethodArgs) (time.Time, error) {
		return rollToFollowingMonday(time.Date(a.Year, month, day, 0, 0, 0, 0, time.UTC)), nil
	}
}

func rollToFollowingMonday(d time.Time) time.Time {
	switch w := int(d.Weekday()); {
	case w == 0:
		return d.AddDate(0, 0, 1)
	case w > 1:
		return d.AddDate(0, 0, 8-w)
	}
	return d
}
