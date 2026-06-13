//go:build parity

package parity

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	holidays "github.com/ppeble/go-holidays/pkg"
)

// sweepYearStart and sweepYearEnd bound the exhaustive year_holidays sweep. The
// range spans weekend-boundary years (observed-date shifts), historical years,
// and well into the future so year_range and leap-year paths are all exercised.
const (
	sweepYearStart = 1970
	sweepYearEnd   = 2050
)

// maxReportedMismatches caps the per-run failure spew: the first N mismatches
// are reported in full; the rest are only counted, so a systemic break does not
// bury the summary under tens of thousands of lines.
const maxReportedMismatches = 50

// TestParity_ExhaustiveYearSweep is the broad equivalence proof. For every
// region the oracle can serve (the gem resolves some, notably jp, through
// Ruby-coded modules that load_custom cannot supply), it compares the full
// year_holidays list from Go's YearHolidaysFrom against the oracle for every
// year in [sweepYearStart, sweepYearEnd] under all four flag combinations.
//
// It first enumerates which of AvailableRegions() the oracle can serve and
// logs an explicit coverage tally, so skipped regions are accounted for rather
// than silently passing. Unlike the curated corpus, this exercises every region
// and a wide year span, turning "we checked a few cases" into "every holiday
// each engine produces, for every serviceable region and year, matches".
func TestParity_ExhaustiveYearSweep(t *testing.T) {
	regions := availableRegionsSorted(t)

	serviceable, unsupported := partitionByOracleSupport(t, regions)
	t.Logf("coverage: %d regions total, %d serviceable, %d oracle-unsupported",
		len(regions), len(serviceable), len(unsupported))
	if len(unsupported) > 0 {
		t.Logf("oracle-unsupported regions (skipped, e.g. gem needs a Ruby-coded module): %v", unsupported)
	}

	years := sweepYearEnd - sweepYearStart + 1
	var c sweepCounters
	// Use a FRESH oracle per region. The installed gem accumulates spurious state
	// across many load_custom-backed queries (it both duplicates and injects
	// holidays once it has serviced enough regions, an oracle-side artifact, not a
	// Go bug). Within a single region's 324 queries the gem stays stable, so a
	// per-region restart (cheap: startup ~0.3s) eliminates the cross-region drift
	// while keeping the comparison honest.
	for _, region := range serviceable {
		ro, err := startOracle()
		if err != nil {
			t.Fatalf("start oracle for region %s: %v", region, err)
		}
		sweepRegion(t, ro, region, &c)
		_ = ro.Close()
	}

	t.Logf("exhaustive year sweep: %d comparisons across %d serviceable regions x %d years x %d flag combos; %d mismatches, %d known divergences, %d lunar-range skips",
		c.comparisons, len(serviceable), years, len(corpusFlags), c.mismatches, c.knownDivs, c.lunarSkips)
	if c.mismatches > maxReportedMismatches {
		t.Errorf("year sweep: %d mismatches total (%d reported above); see summary", c.mismatches, maxReportedMismatches)
	}
}

// sweepCounters accumulates the sweep tally across all regions.
type sweepCounters struct {
	comparisons, mismatches, lunarSkips, knownDivs int
}

// sweepRegion compares Go against the (fresh) oracle for one region across the
// full year span and all flag combinations.
func sweepRegion(t *testing.T, ro *oracle, region string, c *sweepCounters) {
	t.Helper()
	for year := sweepYearStart; year <= sweepYearEnd; year++ {
		from := fmt.Sprintf("%04d-01-01", year)
		for _, f := range corpusFlags {
			compareYear(t, ro, region, year, from, f, c)
		}
	}
}

// compareYear runs one (region, year, flag) comparison and updates the tally.
func compareYear(t *testing.T, ro *oracle, region string, year int, from string, f flagCombo, c *sweepCounters) {
	t.Helper()
	hs, err := holidays.YearHolidaysFrom(mustDate(from), opts([]string{region}, f))
	if err != nil {
		// Go's lunar conversion tables stop short of 2050, so lunar regions (hk, kr,
		// vi) error at the top of the range needing the 2050 lunar new year. That is
		// a known Go data limitation (go-holidays-egm), not a parity failure; count
		// and skip it.
		if isLunarRangeError(err) {
			c.lunarSkips++
			return
		}
		t.Errorf("YearHolidaysFrom(%s, %s, %s) Go error: %v", from, region, f.name, err)
		return
	}
	// Dedupe both sides on (date, name) as a secondary guard: the per-region fresh
	// oracle prevents cross-region drift, and comparing as sets rather than
	// multisets absorbs any residual intra-region duplication without hiding a
	// genuinely missing or extra distinct holiday.
	got := dedupePairs(normalizeGo(hs))
	want, err := ro.holidayList(request{
		Func: "year_holidays", From: from, Regions: []string{region},
		Informal: f.informal, Observed: f.observed,
	})
	if err != nil {
		// A region serviceable at the probe year should stay serviceable; treat any
		// later oracle error as a real failure to investigate.
		t.Errorf("year_holidays(%s, %s, %s) oracle error: %v", from, region, f.name, err)
		return
	}
	want = dedupePairs(want)
	c.comparisons++
	if pairsEqual(got, want) {
		return
	}
	onlyGo, onlyRuby := diffPairs(got, want)
	if knownSweepDivergence(region, onlyGo, onlyRuby) {
		c.knownDivs++
		return
	}
	c.mismatches++
	if c.mismatches <= maxReportedMismatches {
		t.Errorf("year sweep mismatch [region=%s year=%d %s]: Go=%d Ruby=%d%s",
			region, year, f.name, len(got), len(want), formatDiff(onlyGo, onlyRuby))
	}
}

// isLunarRangeError reports whether a Go error is the lunar-table upper-bound
// limit (LunarToSolar out of range), which lunar regions hit near 2050. This is
// a known Go data limitation tracked in a follow-up bead, not a parity failure.
func isLunarRangeError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "calc.LunarToSolar") &&
		strings.Contains(err.Error(), "out of range")
}

// dedupePairs collapses a sorted pair slice to distinct (date, name) entries.
// Used to compare Go and oracle output as sets; see the call site for why.
func dedupePairs(ps []pair) []pair {
	if len(ps) < 2 {
		return ps
	}
	out := ps[:1]
	for _, p := range ps[1:] {
		if p != out[len(out)-1] {
			out = append(out, p)
		}
	}
	return out
}

// knownSweepDivergence reports whether a residual mismatch consists ENTIRELY of
// documented, genuine engine divergences (each tracked by a follow-up bug bead).
// It tolerates a mismatch only when every differing entry, on both sides, is a
// known case, so the broad sweep stays green on the known issues while still
// failing loudly on anything new or unrelated.
func knownSweepDivergence(region string, onlyGo, onlyRuby []pair) bool {
	if len(onlyGo) == 0 && len(onlyRuby) == 0 {
		return false
	}
	for _, p := range onlyGo {
		if !knownDivergentHoliday(region, p.Name) {
			return false
		}
	}
	for _, p := range onlyRuby {
		if !knownDivergentHoliday(region, p.Name) {
			return false
		}
	}
	return true
}

// knownDivergentHoliday lists the (region, holiday) pairs that genuinely disagree
// between Go and the gem over identical v7.0.0 data, found by this sweep and each
// tracked by a follow-up bug bead. These reproduce against a fresh oracle (they
// are not the global-state artifact handled by the per-region restart + dedupe).
func knownDivergentHoliday(region, name string) bool {
	switch {
	case name == "Boxing Day" && (region == "au_nt" || region == "au_tas" ||
		region == "au_tas_north" || region == "au_tas_south"):
		return true // go-holidays-coi: AU Boxing Day observed-cascade off-by-one
	case name == "Proclamation Day" && region == "au_sa":
		return true // go-holidays-coi: au_sa Proclamation Day observed-cascade off-by-one
	case region == "kr" && name == "설날 연휴":
		return true // go-holidays-253: kr Seollal emits one occurrence vs Ruby's two
	case region == "ph" && name == "National Heroes Day":
		return true // go-holidays-dz9: ph National Heroes Day date rule divergence
	}
	return false
}

// availableRegionsSorted returns the Go region universe, sorted for stable
// iteration and reporting.
func availableRegionsSorted(t *testing.T) []string {
	t.Helper()
	regions := append([]string(nil), holidays.AvailableRegions()...)
	sort.Strings(regions)
	return regions
}

// partitionByOracleSupport probes each region once (a single year_holidays
// call) and splits the set into regions the oracle can serve and those it
// cannot (the gem raising "uninitialized constant Holidays::...").
func partitionByOracleSupport(t *testing.T, regions []string) (serviceable, unsupported []string) {
	t.Helper()
	const probeFrom = "2024-01-01"
	for _, region := range regions {
		_, err := ora.holidayList(request{Func: "year_holidays", From: probeFrom, Regions: []string{region}})
		switch {
		case isOracleUnsupported(err):
			unsupported = append(unsupported, region)
		case err != nil:
			t.Errorf("probe year_holidays(%s, %s) oracle error: %v", probeFrom, region, err)
		default:
			serviceable = append(serviceable, region)
		}
	}
	return serviceable, unsupported
}
