package holidays_test

import (
	"time"

	holidays "github.com/ppeble/go-holidays"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Seollal eve (설날 연휴) is the day before Seollal, the first day of the first
// lunar month. Upstream kr.yaml models it as kr_seollal_eve(year, region),
// which is lunar_to_solar(year, 1, 1, region) minus one day.
var _ = Describe("KR Seollal eve", func() {
	DescribeTable("falls on the day before Seollal",
		func(dateStr string) {
			d, err := time.Parse("2006-01-02", dateStr)
			Expect(err).NotTo(HaveOccurred())
			hs, err := holidays.On(d, holidays.Options{Regions: []string{"kr"}, Informal: true})
			Expect(err).NotTo(HaveOccurred())
			Expect(hasNamedHoliday(hs, "설날 연휴")).To(BeTrue(),
				"expected Seollal eve on %s", dateStr)
		},
		Entry("2017", "2017-01-27"),
		Entry("2020", "2020-01-24"),
		Entry("2022", "2022-01-31"),
		Entry("2025", "2025-01-28"),
	)
})
