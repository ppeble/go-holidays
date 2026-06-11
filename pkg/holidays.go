// Package holidays is a Go port of the top-level public API of the Ruby
// `holidays` gem. It answers three questions: "what holidays fall on this
// date?", "what holidays fall in this range?", and "what holidays are in
// this year?", scoped to a set of regions.
package holidays

import (
	"fmt"
	"sort"
	"time"

	"github.com/ppeble/go-holidays/pkg/engine"
)

// On returns every holiday matching the given options on the calendar date of `date`.
// Mirrors Ruby's Holidays.on, which delegates to Holidays.between(date, date),
// so a cached range covering `date` will satisfy this call from the cache.
func On(date time.Time, opts Options) ([]Holiday, error) {
	return Between(date, date, opts)
}

// Between returns every holiday matching the given options whose date falls in
// [start, end] (inclusive on both ends, compared by calendar day). If a previous
// CacheBetween call covers [start, end] for the same options, the result is
// returned from the cache.
func Between(start, end time.Time, opts Options) ([]Holiday, error) {
	if end.Before(start) {
		return nil, fmt.Errorf("holidays.Between: end %s is before start %s",
			end.Format("2006-01-02"), start.Format("2006-01-02"))
	}
	if hs, ok := cacheFind(truncateToDay(start), truncateToDay(end), opts); ok {
		return hs, nil
	}
	return computeBetween(start, end, opts)
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

// YearHolidaysFrom returns every holiday matching the given options from `from`
// (truncated to its UTC calendar day) through Dec 31 of `from`'s year, sorted by
// date ascending. Mirrors the upstream Ruby gem's Holidays.year_holidays(options,
// from_date), which clips a 12-month forward window to Dec 31: this includes
// next-year holidays whose observed date shifts back on or before Dec 31 of
// `from`'s year (for example New Year's Day observed on Dec 31). Resolving both
// `from`'s year and the following year, then applying the [fromDay, Dec31] clip,
// captures those without duplication: a given holiday instance appears in only
// one year's ResolveYear.
func YearHolidaysFrom(from time.Time, opts Options) ([]Holiday, error) {
	fromDay := truncateToDay(from)
	upper := time.Date(fromDay.Year(), 12, 31, 0, 0, 0, 0, fromDay.Location())
	resolveOpts := engine.ResolveOptions{
		Regions:  opts.Regions,
		Informal: opts.Informal,
		Observed: opts.Observed,
	}
	var out []Holiday
	for _, year := range []int{fromDay.Year(), fromDay.Year() + 1} {
		resolved, err := engine.ResolveYear(year, resolveOpts)
		if err != nil {
			return nil, err
		}
		for _, r := range resolved {
			if r.Date.Before(fromDay) || r.Date.After(upper) {
				continue
			}
			out = append(out, Holiday(r))
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Date.Before(out[j].Date)
	})
	return out, nil
}

// nextHolidaysMaxForwardYears bounds the forward scan in NextHolidays so a
// region with too few holidays cannot loop unbounded; the loop stops once it
// has collected `count` holidays or has scanned this many years past `from`.
const nextHolidaysMaxForwardYears = 100

// NextHolidays returns the next `count` holidays on or after `from`, sorted by
// date ascending. Mirrors the upstream Ruby gem's Holidays.next_holidays: it
// keeps resolving forward, year by year, accumulating holidays with date >=
// `from` until at least `count` are gathered, then truncates to `count`. Unlike
// a fixed 12-month window, this returns the full `count` even when the count-th
// holiday lands more than a year out (for example a January holiday relative to
// a January start). If the region has fewer than `count` holidays within the
// safety cap, it returns however many were found.
func NextHolidays(from time.Time, count int, opts Options) ([]Holiday, error) {
	if count <= 0 {
		return nil, fmt.Errorf("holidays.NextHolidays: count must be positive, got %d", count)
	}
	fromDay := truncateToDay(from)
	resolveOpts := engine.ResolveOptions{
		Regions:  opts.Regions,
		Informal: opts.Informal,
		Observed: opts.Observed,
	}
	var collected []Holiday
	startYear := fromDay.Year()
	for offset := 0; offset <= nextHolidaysMaxForwardYears; offset++ {
		resolved, err := engine.ResolveYear(startYear+offset, resolveOpts)
		if err != nil {
			return nil, err
		}
		for _, r := range resolved {
			if r.Date.Before(fromDay) {
				continue
			}
			collected = append(collected, Holiday(r))
		}
		if len(collected) >= count {
			break
		}
	}
	sort.SliceStable(collected, func(i, j int) bool {
		return collected[i].Date.Before(collected[j].Date)
	})
	if len(collected) > count {
		collected = collected[:count]
	}
	return collected, nil
}

// AnyHolidaysDuringWorkWeek reports whether any holiday matching opts falls
// during the Mon-Fri work week containing `date`. Mirrors the upstream Ruby
// gem's Holidays.any_holidays_during_work_week?: for a Saturday input, the
// work week is the preceding Mon-Fri; for a Sunday input, the following.
func AnyHolidaysDuringWorkWeek(date time.Time, opts Options) (bool, error) {
	d := truncateToDay(date)
	wday := int(d.Weekday())
	monday := d.AddDate(0, 0, -(wday - 1))
	friday := d.AddDate(0, 0, 5-wday)
	hs, err := Between(monday, friday, opts)
	if err != nil {
		return false, err
	}
	return len(hs) > 0, nil
}

// AvailableRegions returns every region code registered, sorted lexicographically.
func AvailableRegions() []string {
	return engine.AvailableRegions()
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
