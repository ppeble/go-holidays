package holidays_test

import (
	"testing"
	"time"

	holidays "github.com/ppeble/go-holidays/pkg"
	_ "github.com/ppeble/go-holidays/pkg/definitions"
)

func TestNextHolidays_BasicUS(t *testing.T) {
	from := time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.NextHolidays(from, 5, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("NextHolidays: %v", err)
	}
	if len(hs) != 5 {
		t.Fatalf("want 5 holidays, got %d", len(hs))
	}
	want := []struct {
		date string
		name string
	}{
		{"2026-06-19", "Juneteenth National Independence Day"},
		{"2026-07-04", "Independence Day"},
		{"2026-09-07", "Labor Day"},
		{"2026-11-11", "Veterans Day"},
		{"2026-11-26", "Thanksgiving"},
	}
	for i, w := range want {
		if got := hs[i].Date.Format("2006-01-02"); got != w.date {
			t.Errorf("[%d] date: got %s, want %s", i, got, w.date)
		}
		if hs[i].Name != w.name {
			t.Errorf("[%d] name: got %q, want %q", i, hs[i].Name, w.name)
		}
	}
}

func TestNextHolidays_CrossesYearBoundary(t *testing.T) {
	from := time.Date(2026, 12, 20, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.NextHolidays(from, 3, holidays.Options{
		Regions:  []string{"us"},
		Observed: true,
	})
	if err != nil {
		t.Fatalf("NextHolidays: %v", err)
	}
	if len(hs) != 3 {
		t.Fatalf("want 3 holidays, got %d", len(hs))
	}
	if got, want := hs[0].Date.Format("2006-01-02"), "2026-12-25"; got != want {
		t.Errorf("hs[0] date: got %s, want %s", got, want)
	}
	if got, want := hs[2].Date.Year(), 2027; got != want {
		t.Errorf("hs[2] year: got %d, want %d", got, want)
	}
}

func TestNextHolidays_StartsOnFromDate(t *testing.T) {
	// July 4 is Independence Day; from=July 4 must include it.
	from := time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.NextHolidays(from, 1, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("NextHolidays: %v", err)
	}
	if len(hs) != 1 || hs[0].Date.Format("2006-01-02") != "2026-07-04" {
		t.Fatalf("expected July 4, 2026; got %+v", hs)
	}
}

func TestNextHolidays_ReturnsFewerWhenWindowExhausted(t *testing.T) {
	// Asking for 1000 holidays from a single region caps out at whatever the
	// 12-month window provides; must not error, must not exceed cap.
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.NextHolidays(from, 1000, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("NextHolidays: %v", err)
	}
	if len(hs) >= 1000 {
		t.Fatalf("expected fewer than 1000 holidays in a year for us; got %d", len(hs))
	}
	if len(hs) == 0 {
		t.Fatalf("expected some us holidays in 2026")
	}
}

func TestNextHolidays_CountMustBePositive(t *testing.T) {
	from := time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC)
	for _, c := range []int{0, -1, -100} {
		if _, err := holidays.NextHolidays(from, c, holidays.Options{Regions: []string{"us"}}); err == nil {
			t.Errorf("count=%d: expected error, got nil", c)
		}
	}
}

func TestAnyHolidaysDuringWorkWeek_True(t *testing.T) {
	// Week containing July 4, 2026 (Saturday) — the preceding Mon-Fri is
	// Jun 29-Jul 3, which contains no US holidays. Use July 3 (Friday) to
	// land in the same Mon-Fri window without Independence Day. Pick a
	// week we know has a holiday: Thanksgiving 2026 falls Thursday Nov 26.
	wed := time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC)
	got, err := holidays.AnyHolidaysDuringWorkWeek(wed, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("AnyHolidaysDuringWorkWeek: %v", err)
	}
	if !got {
		t.Fatal("expected true for week containing Thanksgiving 2026")
	}
}

func TestAnyHolidaysDuringWorkWeek_False(t *testing.T) {
	// Week of Aug 10-14, 2026 has no US federal holidays.
	tue := time.Date(2026, 8, 11, 0, 0, 0, 0, time.UTC)
	got, err := holidays.AnyHolidaysDuringWorkWeek(tue, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("AnyHolidaysDuringWorkWeek: %v", err)
	}
	if got {
		t.Fatal("expected false for Aug 10-14 2026 (us)")
	}
}

func TestAnyHolidaysDuringWorkWeek_SaturdayUsesPriorWeek(t *testing.T) {
	// Saturday Nov 28, 2026 — Ruby maps this to Mon Nov 23 - Fri Nov 27,
	// which contains Thanksgiving (Thu Nov 26).
	sat := time.Date(2026, 11, 28, 0, 0, 0, 0, time.UTC)
	got, err := holidays.AnyHolidaysDuringWorkWeek(sat, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("AnyHolidaysDuringWorkWeek: %v", err)
	}
	if !got {
		t.Fatal("Saturday should use prior Mon-Fri; expected Thanksgiving hit")
	}
}

func TestAnyHolidaysDuringWorkWeek_SundayUsesFollowingWeek(t *testing.T) {
	// Sunday July 5, 2026 — Ruby maps to Mon Jul 6 - Fri Jul 10, which has
	// no US federal holidays. (July 4 itself is excluded since it's Saturday.)
	sun := time.Date(2026, 7, 5, 0, 0, 0, 0, time.UTC)
	got, err := holidays.AnyHolidaysDuringWorkWeek(sun, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("AnyHolidaysDuringWorkWeek: %v", err)
	}
	if got {
		t.Fatal("Sunday should use following Mon-Fri (Jul 6-10), which has no us holidays")
	}
}

func TestCacheBetween_PopulatesAndReturnsSameResults(t *testing.T) {
	holidays.ResetCache()
	defer holidays.ResetCache()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	opts := holidays.Options{Regions: []string{"us"}}

	uncached, err := holidays.Between(from, to, opts)
	if err != nil {
		t.Fatalf("uncached Between: %v", err)
	}

	if err := holidays.CacheBetween(from, to, opts); err != nil {
		t.Fatalf("CacheBetween: %v", err)
	}

	cached, err := holidays.Between(from, to, opts)
	if err != nil {
		t.Fatalf("cached Between: %v", err)
	}
	if len(cached) != len(uncached) {
		t.Fatalf("cached len=%d, uncached len=%d", len(cached), len(uncached))
	}
}

func TestCacheBetween_OnHitsCache(t *testing.T) {
	// CacheBetween a year, then On for a known holiday date should resolve
	// without recomputing — verified indirectly by correctness, since Between
	// is what On delegates to.
	holidays.ResetCache()
	defer holidays.ResetCache()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	opts := holidays.Options{Regions: []string{"us"}}
	if err := holidays.CacheBetween(from, to, opts); err != nil {
		t.Fatalf("CacheBetween: %v", err)
	}
	july4 := time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.On(july4, opts)
	if err != nil {
		t.Fatalf("On: %v", err)
	}
	if len(hs) == 0 || hs[0].Name != "Independence Day" {
		t.Fatalf("expected Independence Day, got %+v", hs)
	}
}

func TestCacheBetween_NarrowerRangeReturnsSubset(t *testing.T) {
	holidays.ResetCache()
	defer holidays.ResetCache()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	opts := holidays.Options{Regions: []string{"us"}}
	if err := holidays.CacheBetween(from, to, opts); err != nil {
		t.Fatalf("CacheBetween: %v", err)
	}
	jul := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	julEnd := time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.Between(jul, julEnd, opts)
	if err != nil {
		t.Fatalf("Between: %v", err)
	}
	for _, h := range hs {
		if h.Date.Before(jul) || h.Date.After(julEnd) {
			t.Fatalf("date %s outside requested range", h.Date.Format("2006-01-02"))
		}
	}
}

func TestYearHolidaysFrom_FromJan1ReturnsFullYear(t *testing.T) {
	// Mirrors gem: year_holidays(["us"], Jan 1 2024) => all 10 US holidays.
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.YearHolidaysFrom(from, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("YearHolidaysFrom: %v", err)
	}
	if len(hs) != 10 {
		t.Fatalf("want 10 US holidays for 2024, got %d", len(hs))
	}
	for _, name := range []string{"New Year's Day", "Independence Day", "Thanksgiving", "Christmas Day"} {
		if !hasNamedHoliday(hs, name) {
			t.Errorf("expected %q in full-year results", name)
		}
	}
	for i := 1; i < len(hs); i++ {
		if hs[i].Date.Before(hs[i-1].Date) {
			t.Fatalf("results not sorted ascending at index %d", i)
		}
	}
}

func TestYearHolidaysFrom_FromMidYearReturnsRemainder(t *testing.T) {
	// Mirrors gem: year_holidays(["us"], Nov 15 2024) => only Thanksgiving
	// (2024-11-28) and Christmas (2024-12-25).
	from := time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.YearHolidaysFrom(from, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("YearHolidaysFrom: %v", err)
	}
	if len(hs) != 2 {
		t.Fatalf("want exactly 2 holidays after Nov 15 2024, got %d: %+v", len(hs), hs)
	}
	if got, want := hs[0].Date.Format("2006-01-02"), "2024-11-28"; got != want {
		t.Errorf("hs[0] date: got %s, want %s", got, want)
	}
	if got, want := hs[1].Date.Format("2006-01-02"), "2024-12-25"; got != want {
		t.Errorf("hs[1] date: got %s, want %s", got, want)
	}
	if hasNamedHoliday(hs, "New Year's Day") {
		t.Error("New Year's Day must be absent when starting from Nov 15")
	}
}

func hasNamedHoliday(hs []holidays.Holiday, name string) bool {
	for _, h := range hs {
		if h.Name == name {
			return true
		}
	}
	return false
}

func TestNextHolidays_ResultsSortedAscending(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.NextHolidays(from, 50, holidays.Options{Regions: []string{"us", "ca"}})
	if err != nil {
		t.Fatalf("NextHolidays: %v", err)
	}
	for i := 1; i < len(hs); i++ {
		if hs[i].Date.Before(hs[i-1].Date) {
			t.Fatalf("results not sorted at index %d: %s before %s",
				i, hs[i].Date.Format("2006-01-02"), hs[i-1].Date.Format("2006-01-02"))
		}
	}
}
