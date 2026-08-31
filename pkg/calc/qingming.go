package calc

import (
	"fmt"
	"math"
	"time"
)

// Qingming returns the Gregorian date of the Qingming solar term (清明), which
// anchors China's Tomb-Sweeping Day. It ports the astronomical approximation
// used by the upstream Ruby gem: a per-century base constant plus a linear
// drift term, corrected for leap years. Qingming always falls on April 4 or 5
// within the supported range.
func Qingming(year int) (time.Time, error) {
	var constant float64
	switch {
	case year >= 1900 && year <= 1999:
		constant = 5.59
	case year >= 2000 && year <= 2099:
		constant = 4.81
	default:
		return time.Time{}, fmt.Errorf("calc.Qingming: year %d out of supported range", year)
	}
	y := year % 100
	day := int(math.Floor(float64(y)*0.2422+constant)) - y/4
	return time.Date(year, time.April, day, 0, 0, 0, 0, time.UTC), nil
}
