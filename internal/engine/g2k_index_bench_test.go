package engine

// Throwaway benchmark for bead go-holidays-g2k. Measures rulesFor()'s current
// full-scan-of-every-country cost against a region-indexed map lookup, using
// synthetic data sized to match the REAL per-country rule counts (grepped
// from internal/definitions/*.go: 81 countries, 1547 total rules, "us" alone
// has 67), so the comparison reflects the actual scan size without importing
// internal/definitions into this package's test binary (which would be an
// import cycle: definitions -> engine -> [this test] -> definitions).
//
// Not meant to be kept as a permanent test; delete after the bead closes.

import (
	"testing"

	"github.com/ppeble/go-holidays/internal/definition"
)

// realCountryRuleCounts mirrors internal/definitions/*.go rule counts as of
// the definitions v9.0.0 submodule pin (grepped 2026-09-04).
var realCountryRuleCounts = map[string]int{
	"us": 67, "gb": 40, "de": 45, "au": 49, "ar": 21, "at": 13,
	"ca": 30, "fr": 20, "jp": 40, "cn": 15, "br": 25, "mx": 20,
	"nz": 30, "za": 20, "in": 15, "es": 25, "it": 20, "nl": 20,
	"se": 25, "no": 20, "fi": 20, "dk": 20, "pl": 20, "pt": 20,
	"ru": 15, "kr": 20, "gr": 15, "ch": 20, "be_fr": 12, "be_nl": 12,
}

func buildSyntheticRegistry(b *testing.B) map[string][]definition.HolidayRule {
	b.Helper()
	out := map[string][]definition.HolidayRule{}
	total := 0
	for cc, n := range realCountryRuleCounts {
		rules := make([]definition.HolidayRule, n)
		for i := range rules {
			rules[i] = definition.HolidayRule{Name: cc, Regions: []string{cc}}
		}
		out[cc] = rules
		total += n
	}
	// Pad remaining ~51 countries (of the real 81) at the real average
	// (1547 total / 81 countries =~ 19/country) so the full-scan size matches
	// the real ~1547-rule corpus.
	avg := 19
	for i := 0; i < 51; i++ {
		cc := "pad" + string(rune('a'+i%26)) + string(rune('a'+(i/26)%26))
		rules := make([]definition.HolidayRule, avg)
		for j := range rules {
			rules[j] = definition.HolidayRule{Name: cc, Regions: []string{cc}}
		}
		out[cc] = rules
		total += avg
	}
	b.Logf("synthetic registry: %d countries, %d total rules (real corpus: 81 countries, 1547 rules)", len(out), total)
	return out
}

func BenchmarkRulesFor_CurrentFullScan(b *testing.B) {
	reg := buildSyntheticRegistry(b)
	req := []string{"us"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out []definition.HolidayRule
		for _, rs := range reg {
			out = append(out, rs...)
		}
		cnt := 0
		for _, r := range out {
			if ruleMatchesRequested(r, req) {
				cnt++
			}
		}
		if cnt != realCountryRuleCounts["us"] {
			b.Fatalf("expected %d us matches, got %d", realCountryRuleCounts["us"], cnt)
		}
	}
}

func BenchmarkRulesFor_IndexedLookup(b *testing.B) {
	reg := buildSyntheticRegistry(b)
	// Build once, as a real index would be built at RegisterCountry time.
	idx := map[string][]definition.HolidayRule{}
	for _, rs := range reg {
		for _, r := range rs {
			for _, reg := range r.Regions {
				idx[reg] = append(idx[reg], r)
			}
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matched := idx["us"]
		if len(matched) != realCountryRuleCounts["us"] {
			b.Fatalf("expected %d us matches, got %d", realCountryRuleCounts["us"], len(matched))
		}
	}
}
