package holidays_test

import (
	"testing"
	"time"

	holidays "github.com/ppeble/go-holidays/pkg"
	_ "github.com/ppeble/go-holidays/pkg/definitions"
)

// go-holidays-coi: the two AU year-based Boxing/Proclamation observance methods
// must match holidays gem 11.0.0 exactly. Ground truth captured from the gem
// (Holidays.on(..., :observed)) across every Dec-26 weekday.
//
//   Boxing Day (au_tas, au_nt) uses to_weekday_if_boxing_weekend_from_year, which
//   the gem defines as to_tuesday_if_sunday_or_monday_if_saturday(Dec26):
//   Sat->+2, Sun->+2, Mon UNCHANGED, otherwise unchanged.
//
//   Proclamation Day (au_sa) uses ..._or_to_tuesday_if_monday, which the gem
//   defines as to_weekday_if_boxing_weekend(Dec26):
//   Sat->+2, Sun->+2, Mon->+1, otherwise unchanged.

func mustOn(t *testing.T, dateStr, region string) []holidays.Holiday {
	t.Helper()
	d, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		t.Fatalf("parse %s: %v", dateStr, err)
	}
	hs, err := holidays.On(d, holidays.Options{Regions: []string{region}, Observed: true})
	if err != nil {
		t.Fatalf("On(%s, %s): %v", dateStr, region, err)
	}
	return hs
}

func TestAU_BoxingDayFromYear_MatchesGem(t *testing.T) {
	// date string -> observed Boxing Day date per the gem, keyed by Dec-26 weekday.
	cases := []struct {
		year int
		wday string
		want string // observed Boxing Day date
	}{
		{1970, "Sat", "1970-12-28"},
		{1971, "Sun", "1971-12-28"},
		{1972, "Tue", "1972-12-26"},
		{1973, "Wed", "1973-12-26"},
		{1974, "Thu", "1974-12-26"},
		{1975, "Fri", "1975-12-26"},
		{1977, "Mon", "1977-12-26"}, // the bug: Go currently emits 1977-12-27
	}
	for _, region := range []string{"au_tas", "au_nt"} {
		for _, c := range cases {
			hs := mustOn(t, c.want, region)
			if !hasNamedHoliday(hs, "Boxing Day") {
				t.Errorf("%s Dec26=%s: Boxing Day not observed on %s (want it here)", region, c.wday, c.want)
			}
		}
	}
}

func TestAU_ProclamationDayFromYear_MatchesGem(t *testing.T) {
	cases := []struct {
		year int
		wday string
		want string // observed Proclamation Day date
	}{
		{1970, "Sat", "1970-12-28"},
		{1971, "Sun", "1971-12-28"}, // the bug: Go currently emits 1971-12-27
		{1972, "Tue", "1972-12-26"},
		{1973, "Wed", "1973-12-26"},
		{1974, "Thu", "1974-12-26"},
		{1975, "Fri", "1975-12-26"},
		{1977, "Mon", "1977-12-27"},
	}
	for _, c := range cases {
		hs := mustOn(t, c.want, "au_sa")
		if !hasNamedHoliday(hs, "Proclamation Day") {
			t.Errorf("au_sa Dec26=%s: Proclamation Day not observed on %s (want it here)", c.wday, c.want)
		}
	}
}
