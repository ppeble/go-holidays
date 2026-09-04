package holidays_test

import (
	holidays "github.com/ppeble/go-holidays"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegionName / RegionNames", func() {
	It("returns the display name for a known region", func() {
		name, ok := holidays.RegionName("gb")
		Expect(ok).To(BeTrue())
		Expect(name).To(Equal("United Kingdom"))
	})

	It("returns comma-ok false for an unknown region", func() {
		name, ok := holidays.RegionName("does_not_exist")
		Expect(ok).To(BeFalse())
		Expect(name).To(Equal(""))
	})

	It("round-trips a non-ASCII display name", func() {
		name, ok := holidays.RegionName("ch_ge")
		Expect(ok).To(BeTrue())
		Expect(name).To(Equal("Genève"))
	})

	It("returns every registered region with the expected count", func() {
		names := holidays.RegionNames()
		Expect(names).To(HaveLen(290))
		Expect(names).To(HaveKeyWithValue("gb", "United Kingdom"))
		Expect(names).To(HaveKeyWithValue("ch_ge", "Genève"))
	})
})
