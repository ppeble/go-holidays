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

func TestNextHolidays_ReturnsAtMostCount(t *testing.T) {
	// The scan expands forward across years until count are collected; a large
	// count is satisfied from multiple years without error or exceeding count.
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.NextHolidays(from, 1000, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("NextHolidays: %v", err)
	}
	if len(hs) > 1000 {
		t.Fatalf("expected at most 1000 holidays; got %d", len(hs))
	}
	if len(hs) == 0 {
		t.Fatalf("expected some us holidays from 2026")
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

func TestYearHolidaysFrom_IncludesNextYearObservedBackIntoRange(t *testing.T) {
	// Mirrors gem: year_holidays(["us"], Jan 1 2021, observed) => 11 holidays,
	// including New Year's Day observed on 2021-12-31 (Jan 1 2022 is a Saturday,
	// observed back to Friday Dec 31 2021).
	from := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.YearHolidaysFrom(from, holidays.Options{Regions: []string{"us"}, Observed: true})
	if err != nil {
		t.Fatalf("YearHolidaysFrom: %v", err)
	}
	if len(hs) != 11 {
		t.Fatalf("want 11 US holidays for 2021 observed, got %d: %+v", len(hs), hs)
	}
	found := false
	for _, h := range hs {
		if h.Name == "New Year's Day" && h.Date.Format("2006-01-02") == "2021-12-31" {
			found = true
		}
	}
	if !found {
		t.Error("expected New Year's Day observed on 2021-12-31 in results")
	}
}

// go-holidays-egm: a lunar region at the top of the lunar-table range must not
// fail its year query just because the year+1 look-ahead cannot resolve. kr's
// lunar tables end at 2049, so ResolveYear(2050) errors; but every 2050 result
// is clipped out of the [Jan 1, Dec 31 2049] window anyway, so YearHolidaysFrom
// for 2049 must succeed rather than propagating the look-ahead error.
//
// The assertion stays focused on egm: the call succeeds, clips to 2049, and
// yields both a lunar (설날) and a gregorian (크리스마스) holiday. It deliberately
// does NOT pin an exact count: kr 2049 also carries the separate go-holidays-253
// 설날 연휴 divergence, whose resolution would change the count without bearing on
// this look-ahead fix.
func TestYearHolidaysFrom_LunarBoundaryLookAheadTolerant(t *testing.T) {
	from := time.Date(2049, 1, 1, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.YearHolidaysFrom(from, holidays.Options{Regions: []string{"kr"}})
	if err != nil {
		t.Fatalf("YearHolidaysFrom(kr, 2049): unexpected error: %v", err)
	}
	if !hasNamedHoliday(hs, "설날") || !hasNamedHoliday(hs, "크리스마스") {
		t.Fatalf("expected both 설날 and 크리스마스 in kr 2049 results, got %+v", hs)
	}
	for _, h := range hs {
		if h.Date.Year() != 2049 {
			t.Errorf("result outside 2049 window: %s %s", h.Date.Format("2006-01-02"), h.Name)
		}
	}
}

// go-holidays-253: Seollal (설날) spans three lunar days, the eve being a
// month:12/mday:30 rule. lunar_to_solar(Y,12,30) lands in January of Y+1, and
// the gem rebuilds it as Date.civil(Y, month, day) to pull that eve back into
// year Y. So Ruby emits two 설날 연휴 per year (Jan eve + Feb day-after); Go must
// match. For 2022 the eve is 2022-01-21 and the day-after is 2022-02-02.
func TestYearHolidaysFrom_KRSeollalEmitsBothHolidayDays(t *testing.T) {
	from := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.YearHolidaysFrom(from, holidays.Options{Regions: []string{"kr"}})
	if err != nil {
		t.Fatalf("YearHolidaysFrom(kr, 2022): %v", err)
	}
	var got []string
	for _, h := range hs {
		if h.Name == "설날 연휴" {
			got = append(got, h.Date.Format("2006-01-02"))
		}
	}
	want := map[string]bool{"2022-01-21": false, "2022-02-02": false}
	for _, d := range got {
		if _, ok := want[d]; ok {
			want[d] = true
		}
	}
	for d, found := range want {
		if !found {
			t.Errorf("expected 설날 연휴 on %s; got 설날 연휴 dates %v", d, got)
		}
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

// go-holidays-wdp: On/Between must include a next-year holiday whose observed
// date shifts back into the requested range. New Year's Day Jan 1 2022 is a
// Saturday, observed back to Friday Dec 31 2021.
func TestOn_NextYearObservedBackIntoRange(t *testing.T) {
	date := time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.On(date, holidays.Options{Regions: []string{"us"}, Observed: true})
	if err != nil {
		t.Fatalf("On: %v", err)
	}
	if !hasNamedHoliday(hs, "New Year's Day") {
		t.Fatalf("expected New Year's Day observed on 2021-12-31, got %+v", hs)
	}
}

func TestBetween_IncludesNextYearObservedBackIntoRange(t *testing.T) {
	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.Between(start, end, holidays.Options{Regions: []string{"us"}, Observed: true})
	if err != nil {
		t.Fatalf("Between: %v", err)
	}
	found := false
	for _, h := range hs {
		if h.Name == "New Year's Day" && h.Date.Format("2006-01-02") == "2021-12-31" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected New Year's Day observed on 2021-12-31 within 2021 range, got %+v", hs)
	}
}

// go-holidays-gmc: a multi-region request must not duplicate a holiday defined
// for more than one of the requested regions. Informal Easter Sunday exists in
// both us and ca; the gem emits it once.
func TestBetween_MultiRegionDeduplicatesSharedHoliday(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.Between(start, end, holidays.Options{Regions: []string{"us", "ca"}, Informal: true})
	if err != nil {
		t.Fatalf("Between: %v", err)
	}
	easter, goodFriday := 0, 0
	for _, h := range hs {
		switch h.Name {
		case "Easter Sunday":
			easter++
		case "Good Friday":
			goodFriday++
		}
	}
	// Easter Sunday is one identical informal def in both regions: merged to one.
	if easter != 1 {
		t.Fatalf("expected Easter Sunday exactly once for [us ca] informal 2024, got %d", easter)
	}
	// Good Friday has two distinct defs (us informal vs untagged us-states/ca):
	// they differ on type, so the gem keeps both. Dedup must not collapse them.
	if goodFriday != 2 {
		t.Fatalf("expected Good Friday twice for [us ca] informal 2024, got %d", goodFriday)
	}
}

// go-holidays-6vj: NextHolidays must return the full count even when the
// count-th holiday lands more than 12 months past `from`. From 2024-01-01 the
// 12th US holiday is MLK Day 2025-01-20, beyond a fixed 12-month window.
func TestNextHolidays_ExpandsBeyondTwelveMonths(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.NextHolidays(from, 12, holidays.Options{Regions: []string{"us"}})
	if err != nil {
		t.Fatalf("NextHolidays: %v", err)
	}
	if len(hs) != 12 {
		t.Fatalf("want 12 US holidays from 2024-01-01, got %d: %+v", len(hs), hs)
	}
	last := hs[11].Date.Format("2006-01-02")
	if last != "2025-01-20" {
		t.Fatalf("expected 12th holiday on 2025-01-20, got %s", last)
	}
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

// go-holidays-2fc: upstream holidays master 272b549 (PR#343 / issue#396, not yet
// in the pinned v8.0.0 definitions tag) fixes gb.yaml Christmas Day / Boxing Day
// observed rules. In 2022 Christmas Day (Dec 25) fell on a Sunday. With the
// pinned v8.0.0 rules (to_monday_if_weekend for Christmas, to_weekday_if_boxing_weekend
// for Boxing Day) both the substitute Christmas Day and the actual Boxing Day
// collide on Monday Dec 26, and no holiday is observed on Dec 27. Upstream's
// fix (to_tuesday_if_sunday_or_monday_if_saturday for both) observes Boxing Day
// on its actual date (Mon Dec 26, since Monday needs no substitution) and moves
// Christmas Day's substitute to Tue Dec 27.
func TestOn_GB_ChristmasSunday2022_ObservedDatesDoNotCollide(t *testing.T) {
	t.Skip("go-holidays-2fc: known bug against pinned v8.0.0 gb.yaml; fails until " +
		"the submodule is bumped past v8.0.0 (upstream master 272b549) and " +
		"pkg/definitions/gb.go is regenerated. Remove this Skip once that lands.")
	dec26 := time.Date(2022, 12, 26, 0, 0, 0, 0, time.UTC)
	hs, err := holidays.On(dec26, holidays.Options{Regions: []string{"gb"}, Observed: true})
	if err != nil {
		t.Fatalf("On: %v", err)
	}
	if !hasNamedHoliday(hs, "Boxing Day") {
		t.Errorf("expected Boxing Day observed on 2022-12-26, got %+v", hs)
	}

	dec27 := time.Date(2022, 12, 27, 0, 0, 0, 0, time.UTC)
	hs, err = holidays.On(dec27, holidays.Options{Regions: []string{"gb"}, Observed: true})
	if err != nil {
		t.Fatalf("On: %v", err)
	}
	if !hasNamedHoliday(hs, "Christmas Day") {
		t.Errorf("expected Christmas Day observed on 2022-12-27 (Sun 25th moved past Boxing Day), got %+v", hs)
	}
}
