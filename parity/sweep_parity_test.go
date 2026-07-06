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
	// Reuse the single shared oracle started in TestMain. Older gems (9.1.0)
	// accumulated cross-region global-state pollution across many load_custom
	// queries, which once forced a fresh per-region oracle (spawn ruby + reload
	// all 80 YAML, ~287x) and made this sweep take minutes. Gem 9.1.2 (#344,
	// region load-order) and 10.0.0 (#352, function_modifier merge) fixed that,
	// so one long-lived oracle produces the same tally in a fraction of the time.
	for _, region := range serviceable {
		sweepRegion(t, ora, region, &c)
	}

	t.Logf("exhaustive year sweep: %d comparisons across %d serviceable regions x %d years x %d flag combos; %d mismatches, %d known divergences, %d mutual lunar-boundary confirmations",
		c.comparisons, len(serviceable), years, len(corpusFlags), c.mismatches, c.knownDivs, c.lunarBoundary)
	if c.mismatches > maxReportedMismatches {
		t.Errorf("year sweep: %d mismatches total (%d reported above); see summary", c.mismatches, maxReportedMismatches)
	}
}

// sweepCounters accumulates the sweep tally across all regions.
type sweepCounters struct {
	comparisons, mismatches, lunarBoundary, knownDivs int
}

// sweepRegion compares Go against the shared oracle for one region across the
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
		// Go's lunar tables stop at 2049, so lunar regions (hk, kr, vn) cannot
		// resolve the primary year 2050. This is a genuine MUTUAL boundary: the
		// gem's tables end at the same year, so Ruby is equally undefined here.
		// Rather than blindly skipping, confirm it: query the oracle and require
		// it to ALSO error. If Ruby instead resolves the year, Go has a real gap.
		// (go-holidays-egm; the year+1 look-ahead at 2049 is handled in Go's
		// YearHolidaysFrom, so only the 2050 primary year reaches here.)
		if isLunarRangeError(err) {
			if _, oerr := ro.holidayList(request{
				Func: "year_holidays", From: from, Regions: []string{region},
				Informal: f.informal, Observed: f.observed,
			}); oerr != nil {
				c.lunarBoundary++
				return
			}
			t.Errorf("lunar boundary [region=%s year=%d %s]: Go errored (%v) but Ruby resolved the year",
				region, year, f.name, err)
			return
		}
		t.Errorf("YearHolidaysFrom(%s, %s, %s) Go error: %v", from, region, f.name, err)
		return
	}
	// Compare Go and oracle output as (date, name) SETS, not multisets. The gem
	// emits benign duplicate rows for some region/flag combos (e.g. br informal
	// years: Ruby yields 13 rows to Go's 12 for the same distinct holidays), an
	// artifact of load_custom merge semantics, not a genuinely missing or extra
	// holiday. This dedupe is empirically required: dropping it surfaces 7290 such
	// spurious mismatches. It is unrelated to the (now removed) per-region oracle.
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
// limit (LunarToSolar out of range), which lunar regions hit at the primary year
// 2050. The gem's tables end at the same year, so compareYear confirms it as a
// mutual boundary rather than treating it as a Go-only failure.
func isLunarRangeError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "calc.LunarToSolar") &&
		strings.Contains(err.Error(), "out of range")
}

// dedupePairs collapses a sorted pair slice to distinct (date, name) entries, so
// compareYear can compare Go and oracle output as sets; see that call site for
// why the gem's benign duplicate rows make this necessary.
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
// between Go and the gem over identical v8.0.0 data, found by this sweep and each
// tracked by a follow-up bug bead. These reproduce against the shared oracle on
// gem 11.0.0 (they are genuine engine divergences, not a gem global-state artifact).
func knownDivergentHoliday(region, name string) bool {
	switch {
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
