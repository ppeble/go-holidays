package calc

import "time"

// LunarToSolar converts a lunar calendar date to its Gregorian equivalent for
// the given region. The upstream definitions repo ships a multi-region lunar
// table; we will populate this as regions that depend on it (cn, kr, hk, vn)
// come online. Until then this returns the zero time so callers treat the
// holiday as not applicable in the requested year, rather than propagating
// an error that would mask other holidays in the same lookup.
func LunarToSolar(year, month, day int, region string) (time.Time, error) {
	return time.Time{}, nil
}
