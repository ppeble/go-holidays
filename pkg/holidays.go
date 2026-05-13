// Package holidays is a Go port of the top-level public API of the Ruby
// `holidays` gem. It answers three questions: "what holidays fall on this
// date?", "what holidays fall in this range?", and "what holidays are in
// this year?", scoped to a set of regions.
package holidays

import (
	"fmt"
	"time"

	"github.com/ppeble/go-holidays/pkg/engine"
)

// On returns every holiday matching the given options on the calendar date of `date`.
func On(date time.Time, opts Options) ([]Holiday, error) {
	resolved, err := engine.ResolveYear(date.Year(), engine.ResolveOptions{
		Regions:  opts.Regions,
		Informal: opts.Informal,
		Observed: opts.Observed,
	})
	if err != nil {
		return nil, err
	}
	var out []Holiday
	for _, r := range resolved {
		if sameDay(r.Date, date) {
			out = append(out, Holiday(r))
		}
	}
	return out, nil
}

// Between returns every holiday matching the given options whose date falls in
// [start, end] (inclusive on both ends, compared by calendar day).
func Between(start, end time.Time, opts Options) ([]Holiday, error) {
	if end.Before(start) {
		return nil, fmt.Errorf("holidays.Between: end %s is before start %s",
			end.Format("2006-01-02"), start.Format("2006-01-02"))
	}
	var out []Holiday
	for y := start.Year(); y <= end.Year(); y++ {
		resolved, err := engine.ResolveYear(y, engine.ResolveOptions{
			Regions:  opts.Regions,
			Informal: opts.Informal,
			Observed: opts.Observed,
		})
		if err != nil {
			return nil, err
		}
		for _, r := range resolved {
			if inRange(r.Date, start, end) {
				out = append(out, Holiday(r))
			}
		}
	}
	return out, nil
}

// YearHolidays returns every holiday matching the given options in the given year.
func YearHolidays(year int, opts Options) ([]Holiday, error) {
	resolved, err := engine.ResolveYear(year, engine.ResolveOptions{
		Regions:  opts.Regions,
		Informal: opts.Informal,
		Observed: opts.Observed,
	})
	if err != nil {
		return nil, err
	}
	out := make([]Holiday, len(resolved))
	for i, r := range resolved {
		out[i] = Holiday(r)
	}
	return out, nil
}

// AvailableRegions returns every region code registered, sorted lexicographically.
func AvailableRegions() []string {
	return engine.AvailableRegions()
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func inRange(d, start, end time.Time) bool {
	if d.Before(truncateToDay(start)) {
		return false
	}
	if d.After(truncateToDay(end)) {
		return false
	}
	return true
}

func truncateToDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
