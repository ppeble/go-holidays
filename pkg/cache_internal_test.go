package holidays

import (
	"testing"
	"time"
)

// Test cache internals directly. Lives in `package holidays` (not _test) so it
// can poke the unexported store.

func TestCacheStoreAndFind_ExactRange(t *testing.T) {
	ResetCache()
	defer ResetCache()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	opts := Options{Regions: []string{"us"}}
	sentinel := []Holiday{
		{Date: time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC), Name: "Test 4", Regions: []string{"us"}},
		{Date: time.Date(2026, 12, 25, 0, 0, 0, 0, time.UTC), Name: "Test 25", Regions: []string{"us"}},
	}
	cacheStore(from, to, opts, sentinel)
	got, ok := cacheFind(from, to, opts)
	if !ok {
		t.Fatal("expected cache hit on exact range")
	}
	if len(got) != 2 {
		t.Fatalf("got %d holidays, want 2", len(got))
	}
}

func TestCacheFind_ContainedRangeHits(t *testing.T) {
	ResetCache()
	defer ResetCache()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	opts := Options{Regions: []string{"us"}}
	sentinel := []Holiday{
		{Date: time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC), Name: "Test 4"},
		{Date: time.Date(2026, 12, 25, 0, 0, 0, 0, time.UTC), Name: "Test 25"},
	}
	cacheStore(from, to, opts, sentinel)
	// Narrower range entirely inside the cached one: should still hit.
	got, ok := cacheFind(
		time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC),
		opts,
	)
	if !ok {
		t.Fatal("expected cache hit on contained range")
	}
	if len(got) != 1 || got[0].Name != "Test 4" {
		t.Fatalf("got %+v, want only Test 4", got)
	}
}

func TestCacheFind_OutOfRangeMisses(t *testing.T) {
	ResetCache()
	defer ResetCache()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	opts := Options{Regions: []string{"us"}}
	cacheStore(from, to, opts, nil)
	// Request range extending past cached upper bound.
	if _, ok := cacheFind(
		time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2027, 1, 31, 0, 0, 0, 0, time.UTC),
		opts,
	); ok {
		t.Fatal("expected cache miss for range extending past cache")
	}
	// Request range starting before cached lower bound.
	if _, ok := cacheFind(
		time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
		opts,
	); ok {
		t.Fatal("expected cache miss for range starting before cache")
	}
}

func TestCacheFind_DifferentOptionsMiss(t *testing.T) {
	ResetCache()
	defer ResetCache()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	cacheStore(from, to, Options{Regions: []string{"us"}}, nil)
	// Different region: miss.
	if _, ok := cacheFind(from, to, Options{Regions: []string{"gb"}}); ok {
		t.Fatal("different region should miss")
	}
	// Same region but informal flipped: miss.
	if _, ok := cacheFind(from, to, Options{Regions: []string{"us"}, Informal: true}); ok {
		t.Fatal("flipping informal should miss")
	}
	// Same region but observed flipped: miss.
	if _, ok := cacheFind(from, to, Options{Regions: []string{"us"}, Observed: true}); ok {
		t.Fatal("flipping observed should miss")
	}
}

func TestOptionsKey_RegionOrderInsensitive(t *testing.T) {
	a := optionsKey(Options{Regions: []string{"us", "gb"}})
	b := optionsKey(Options{Regions: []string{"gb", "us"}})
	if a != b {
		t.Fatalf("options key should be order-insensitive; got %+v vs %+v", a, b)
	}
}

func TestResetCache_Clears(t *testing.T) {
	ResetCache()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	opts := Options{Regions: []string{"us"}}
	cacheStore(from, to, opts, []Holiday{{Date: from, Name: "x"}})
	if _, ok := cacheFind(from, to, opts); !ok {
		t.Fatal("expected hit before reset")
	}
	ResetCache()
	if _, ok := cacheFind(from, to, opts); ok {
		t.Fatal("expected miss after reset")
	}
}
