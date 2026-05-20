package calc

import (
	"fmt"
	"time"
)

// solarStartDate is the Gregorian date corresponding to lunar 1900-01-01
// (Lunar New Year, January 31 1900). All conversions are deltas off this anchor.
var solarStartDate = time.Date(1900, time.January, 31, 0, 0, 0, 0, time.UTC)

// LunarToSolar converts a lunar calendar date (year/month/day in the region's
// lunar calendar) to its Gregorian equivalent. The lookup tables ported from
// the upstream Ruby gem cover hk, kr, and vi from 1900 onward.
func LunarToSolar(year, month, day int, region string) (time.Time, error) {
	info := calendarYearInfoForRegion(region)
	if info == nil {
		return time.Time{}, fmt.Errorf("calc.LunarToSolar: no lunar table for region %q", region)
	}
	yearDiff := year - 1900
	if yearDiff < 0 || yearDiff >= len(info) {
		return time.Time{}, fmt.Errorf("calc.LunarToSolar: year %d out of range for region %q", year, region)
	}
	if month < 1 || month > 12 {
		return time.Time{}, fmt.Errorf("calc.LunarToSolar: invalid lunar month %d", month)
	}

	days := 0
	for i := 0; i < yearDiff; i++ {
		days += info[i][0]
	}
	for i := 0; i < month-1; i++ {
		days += lunarDaysForMonthType[info[yearDiff][i+1]][0]
	}
	days += day - 1
	return solarStartDate.AddDate(0, 0, days), nil
}
