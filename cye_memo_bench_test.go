package holidays

// Throwaway benchmarks for bead go-holidays-cye. Measures where ResolveYear's
// per-call cost actually goes: inside registered method functions (easter,
// lunar table lookups, jp's substitute-holiday helper) versus rule-matching
// and iteration, to size the potential win of a function-result memoization
// cache (the Go equivalent of the gem's ProcResultCache).
//
// Not meant to be a permanent test; delete after the bead closes if the
// verdict is "not worth it", keep if a follow-up bead is filed.

import (
	"sync"
	"testing"

	"github.com/ppeble/go-holidays/internal/calc"
	"github.com/ppeble/go-holidays/internal/engine"
)

// BenchmarkEasterDirect isolates the cost of the closed-form easter()
// calculation itself, called the way computeDate calls it.
func BenchmarkEasterDirect(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = calc.Easter(2000 + i%100)
	}
}

// BenchmarkLunarToSolarDirect isolates the cost of a lunar table lookup, the
// other builtin that varies by more than just year (LunarToSolar also takes
// region, since cn/hk/sg share a table but kr/vn do not).
func BenchmarkLunarToSolarDirect(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = calc.LunarToSolar(1950+i%100, 1, 1, "kr")
	}
}

// BenchmarkResolveYear_AllRegions_Easter measures a full ResolveYear call
// across every registered region for a year with no informal/observed
// filtering, i.e. the worst case for redundant easter() calls: 206 rules
// across the corpus reference "easter" (grepped from internal/definitions),
// so a single ResolveYear(year, opts{}) call with no region filter invokes
// calc.Easter(year) up to 206 times for the exact same (function, year) pair.
func BenchmarkResolveYear_AllRegions(b *testing.B) {
	opts := engine.ResolveOptions{Observed: true, Informal: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.ResolveYear(2024, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkResolveYear_KR measures a lunar-heavy single region (Seollal is
// lunar_to_solar-based, computed via a large embedded table per call).
func BenchmarkResolveYear_KR(b *testing.B) {
	opts := engine.ResolveOptions{Regions: []string{"kr"}, Observed: true, Informal: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.ResolveYear(2024, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkResolveYear_JP measures the substitute-holiday-heavy region: each
// of 6+ jp_*_substitute rules independently calls jpFixedHolidaysByMonth,
// which rebuilds a nested map by re-walking rulesForCountry("jp") every time
// it's invoked (internal/engine/methods_jp.go:99-108) - this is the biggest
// non-lunar, non-easter redundancy found in the method layer.
func BenchmarkResolveYear_JP(b *testing.B) {
	opts := engine.ResolveOptions{Regions: []string{"jp"}, Observed: true, Informal: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.ResolveYear(2024, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkResolveYear_US is the non-lunar, non-jp baseline: a single
// mid-sized region with a handful of function-based rules (easter,
// to_monday_if_sunday, etc), to compare against KR/JP's heavier method cost.
func BenchmarkResolveYear_US(b *testing.B) {
	opts := engine.ResolveOptions{Regions: []string{"us"}, Observed: true, Informal: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.ResolveYear(2024, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkResolveYear_AllRegions above already IS the repeated-same-year
// case (2024 on every b.N iteration) - that's the scenario a function-result
// cache targets: repeated ResolveYear calls for the SAME year, e.g. separate
// On()/Between() calls for different single dates in the same year, each an
// uncached miss under CacheBetween, which only dedupes calls for the same
// [range, opts] key (see cache.go computeBetween, which calls ResolveYear
// once per distinct year regardless of how many dates are requested inside
// one CacheBetween call).
//
// BenchmarkResolveYear_DistinctYears simulates the OTHER real access pattern:
// a multi-year Between() call, where every ResolveYear call is for a
// different year, so a (function, year) cache never hits.
func BenchmarkResolveYear_DistinctYears_AllRegions(b *testing.B) {
	opts := engine.ResolveOptions{Observed: true, Informal: true}
	years := make([]int, b.N)
	for i := range years {
		years[i] = 1950 + i%80
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.ResolveYear(years[i], opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMutexMapLookup estimates the floor cost of a would-be
// function-result cache's lookup path (RLock + map read by a struct key +
// RUnlock), to compare against BenchmarkEasterDirect's 19.82 ns/op. If the
// cache's own bookkeeping costs more than the function it is memoizing,
// caching that function is a net loss on the common (miss-free-but-still-
// looked-up) path.
func BenchmarkMutexMapLookup(b *testing.B) {
	var mu sync.RWMutex
	type key struct {
		fn     string
		year   int
		region string
	}
	m := map[key]int{{"easter", 2024, "us"}: 42}
	k := key{"easter", 2024, "us"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mu.RLock()
		_ = m[k]
		mu.RUnlock()
	}
}
