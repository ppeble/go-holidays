package engine

import (
	"fmt"
	"time"

	"github.com/ppeble/go-holidays/pkg/definition"
)

// Resolved is the internal result of resolving a single rule for a single year.
type Resolved struct {
	Date    time.Time
	Name    string
	Regions []string
}

// ResolveOptions controls filtering during resolution.
type ResolveOptions struct {
	Regions  []string
	Informal bool
	Observed bool
}

// ResolveYear walks the registered rules and returns holidays for the given year.
func ResolveYear(year int, opts ResolveOptions) ([]Resolved, error) {
	rules := rulesFor(opts.Regions)
	out := make([]Resolved, 0, len(rules))
	for _, rule := range rules {
		if !rule.AppliesIn(year) {
			continue
		}
		if rule.Type == definition.Informal && !opts.Informal {
			continue
		}
		if !ruleMatchesRequested(rule, opts.Regions) {
			continue
		}
		date, err := computeDate(rule, year)
		if err != nil {
			return nil, fmt.Errorf("rule %q: %w", rule.Name, err)
		}
		if date.IsZero() {
			continue
		}
		if opts.Observed && rule.HasObserved() {
			date, err = applyObserved(rule.Observed, date)
			if err != nil {
				return nil, fmt.Errorf("rule %q observed=%s: %w", rule.Name, rule.Observed, err)
			}
		}
		out = append(out, Resolved{Date: date, Name: rule.Name, Regions: rule.Regions})
	}
	return out, nil
}

func computeDate(rule definition.HolidayRule, year int) (time.Time, error) {
	var base time.Time
	switch {
	case rule.HasMday():
		base = time.Date(year, time.Month(rule.Month), rule.Mday, 0, 0, 0, 0, time.UTC)
	case rule.HasWday():
		base = nthOrLastWeekday(year, time.Month(rule.Month), rule.Week, time.Weekday(rule.Wday))
		if base.IsZero() {
			return time.Time{}, nil
		}
	}
	if rule.HasFunction() {
		fn, ok := LookupMethod(rule.Function)
		if !ok {
			return time.Time{}, fmt.Errorf("unregistered method %q", rule.Function)
		}
		args := MethodArgs{Year: year, Month: rule.Month, Date: base}
		if len(rule.Regions) > 0 {
			args.Region = rule.Regions[0]
		}
		d, err := fn(args)
		if err != nil {
			return time.Time{}, err
		}
		base = d
	}
	if base.IsZero() {
		if !rule.HasMday() && !rule.HasWday() && !rule.HasFunction() {
			return time.Time{}, fmt.Errorf("rule has no mday, wday/week, or function")
		}
		return time.Time{}, nil
	}
	if rule.FunctionModifier != 0 {
		base = base.AddDate(0, 0, rule.FunctionModifier)
	}
	return base, nil
}

func applyObserved(methodName string, date time.Time) (time.Time, error) {
	fn, ok := LookupMethod(methodName)
	if !ok {
		return time.Time{}, fmt.Errorf("unregistered observed method %q", methodName)
	}
	return fn(MethodArgs{Date: date, Year: date.Year(), Month: int(date.Month()), Day: date.Day()})
}

func nthOrLastWeekday(year int, month time.Month, week int, wday time.Weekday) time.Time {
	if week == -1 {
		firstNext := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
		last := firstNext.AddDate(0, 0, -1)
		offset := (int(last.Weekday()) - int(wday) + 7) % 7
		return last.AddDate(0, 0, -offset)
	}
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	offset := (int(wday) - int(first.Weekday()) + 7) % 7
	day := 1 + offset + (week-1)*7
	candidate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if candidate.Month() != month {
		return time.Time{}
	}
	return candidate
}
