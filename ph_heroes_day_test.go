package holidays_test

import (
	"time"

	holidays "github.com/ppeble/go-holidays"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// assertHeroesDay checks whether "National Heroes Day" resolves on the given
// date for the ph region, asserting presence (want=true) or absence (want=false).
func assertHeroesDay(opts holidays.Options, date time.Time, want bool) {
	hs, err := holidays.On(date, opts)
	Expect(err).NotTo(HaveOccurred())
	Expect(hasNamedHoliday(hs, "National Heroes Day")).To(Equal(want),
		"National Heroes Day on %s: want present=%v", date.Format("2006-01-02"), want)
}

// ph National Heroes Day is the last Monday of August. This is the Go engine's
// authoritative behavior and is asserted here directly, independent of the Ruby
// gem. Upstream ph.yaml used to carry an off-by-one (holidays/definitions#345)
// that emitted September 1 in years where August 31 is a Sunday; Go always
// refused to reproduce that quirk, which made it the parity sweep's last known
// divergence. The fix shipped in definitions v8.0.2 as function: ph_heroes_day(year)
// with a matching implementation in gem 11.2.0, so both engines now agree and the
// sweep allowlist is gone. These cases stay to lock the Go side so a future
// refactor cannot silently regress it to the buggy date.
var _ = Describe("PH National Heroes Day", func() {
	It("falls on the last Monday of August", func() {
		opts := holidays.Options{Regions: []string{"ph"}}

		// The eight years in 1970-2050 where August 31 falls on a Sunday: the
		// correct last Monday of August is the 25th (the buggy Ruby rule yields
		// September 1 for exactly these years).
		sundayAug31Years := []int{1975, 1980, 1986, 1997, 2003, 2008, 2014, 2025}
		for _, y := range sundayAug31Years {
			assertHeroesDay(opts, time.Date(y, time.August, 25, 0, 0, 0, 0, time.UTC), true)
			assertHeroesDay(opts, time.Date(y, time.September, 1, 0, 0, 0, 0, time.UTC), false)
		}

		// Control years where August 31 is not a Sunday: last Monday of August is
		// unambiguous and both engines already agree.
		controls := map[int]int{2020: 31, 2021: 30, 2022: 29, 2023: 28, 2024: 26}
		for y, mday := range controls {
			assertHeroesDay(opts, time.Date(y, time.August, mday, 0, 0, 0, 0, time.UTC), true)
		}
	})
})
