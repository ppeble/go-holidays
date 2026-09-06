package holidays_test

import (
	"time"

	holidays "github.com/ppeble/go-holidays"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// go-holidays-coi: the two AU year-based Boxing/Proclamation observance methods
// and their expected observed dates across every Dec-26 weekday.
//
//   Boxing Day (au_tas, au_nt) uses to_weekday_if_boxing_weekend_from_year,
//   defined as to_tuesday_if_sunday_or_monday_if_saturday(Dec26):
//   Sat->+2, Sun->+2, Mon UNCHANGED, otherwise unchanged.
//
//   Proclamation Day (au_sa) uses ..._or_to_tuesday_if_monday,
//   defined as to_weekday_if_boxing_weekend(Dec26):
//   Sat->+2, Sun->+2, Mon->+1, otherwise unchanged.

func mustOn(dateStr, region string) []holidays.Holiday {
	d, err := time.Parse("2006-01-02", dateStr)
	Expect(err).NotTo(HaveOccurred())
	hs, err := holidays.On(d, holidays.Options{Regions: []string{region}, Observed: true})
	Expect(err).NotTo(HaveOccurred())
	return hs
}

var _ = Describe("AU Boxing/Proclamation Day observance", func() {
	It("matches the gem for Boxing Day (au_tas, au_nt) across every Dec-26 weekday", func() {
		// date string -> expected observed Boxing Day date, keyed by Dec-26 weekday.
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
				hs := mustOn(c.want, region)
				Expect(hasNamedHoliday(hs, "Boxing Day")).To(BeTrue(),
					"%s Dec26=%s: Boxing Day not observed on %s (want it here)", region, c.wday, c.want)
			}
		}
	})

	It("matches the gem for Proclamation Day (au_sa) across every Dec-26 weekday", func() {
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
			hs := mustOn(c.want, "au_sa")
			Expect(hasNamedHoliday(hs, "Proclamation Day")).To(BeTrue(),
				"au_sa Dec26=%s: Proclamation Day not observed on %s (want it here)", c.wday, c.want)
		}
	})
})
