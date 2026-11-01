package engine

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/holidays/go-holidays/internal/definition"
)

// Regression coverage for go-holidays-67s: two definitions that share every
// field compared for identity (name, wday, mday, week, function,
// function_modifier, type, observed, year_ranges) but fall in different
// months are distinct holidays. Definitions are bucketed per month, so month
// is implicitly part of a definition's identity. ruleSignature must therefore
// include the month, or ResolveYear's dedup silently drops the second one on an
// aggregate query (e.g. "au_" losing the ACT/NSW/SA October Labour Day because
// it collides with the WA March Labour Day).
var _ = Describe("ruleSignature month component", func() {
	AfterEach(func() {
		UnregisterCountry("zz")
	})

	It("keeps two same-named first-Monday rules that differ only by month", func() {
		RegisterCountry("zz", []definition.HolidayRule{
			{Name: "Labour Day", Regions: []string{"zz_west"}, Month: 3, Wday: 1, Week: 1},
			{Name: "Labour Day", Regions: []string{"zz_east"}, Month: 10, Wday: 1, Week: 1},
		})

		got, err := ResolveYear(2020, ResolveOptions{Regions: []string{"zz_"}})
		Expect(err).NotTo(HaveOccurred())

		var months []int
		for _, r := range got {
			if r.Name == "Labour Day" {
				months = append(months, int(r.Date.Month()))
			}
		}
		Expect(months).To(ConsistOf(3, 10),
			"expected both the March and October Labour Day, got %+v", got)
	})
})

// Rules that differ only by excluded_regions are distinct definitions (the
// gem's generator never merges them), so ResolveYear must not collapse them.
var _ = Describe("ruleSignature excluded regions component", func() {
	AfterEach(func() {
		UnregisterCountry("zz")
	})

	It("keeps two rules that differ only by their exclusions", func() {
		RegisterCountry("zz", []definition.HolidayRule{
			{Name: "Example Day", Regions: []string{"zz"}, ExcludedRegions: []string{"zz_a"}, Month: 6, Mday: 19, Wday: -1},
			{Name: "Example Day", Regions: []string{"zz_a"}, Month: 6, Mday: 19, Wday: -1},
		})

		got, err := ResolveYear(2020, ResolveOptions{Regions: []string{"zz", "zz_a"}})
		Expect(err).NotTo(HaveOccurred())

		var regions [][]string
		for _, r := range got {
			if r.Name == "Example Day" {
				regions = append(regions, r.Regions)
			}
		}
		Expect(regions).To(ConsistOf([]string{"zz"}, []string{"zz_a"}))
	})

	It("still merges rules whose exclusions match", func() {
		RegisterCountry("zz", []definition.HolidayRule{
			{Name: "Example Day", Regions: []string{"zz"}, ExcludedRegions: []string{"zz_b_x"}, Month: 6, Mday: 19, Wday: -1},
			{Name: "Example Day", Regions: []string{"zz_b"}, ExcludedRegions: []string{"zz_b_x"}, Month: 6, Mday: 19, Wday: -1},
		})

		got, err := ResolveYear(2020, ResolveOptions{Regions: []string{"zz_b"}})
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(HaveLen(1))
		Expect(got[0].Regions).To(Equal([]string{"zz"}))
	})
})
