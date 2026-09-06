package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("reorderFlagsFirst", func() {
	It("moves a flag placed after positionals to the front", func() {
		got := reorderFlagsFirst([]string{"5", "2026-01-01", "--regions", "us"})
		Expect(got).To(Equal([]string{"--regions", "us", "--", "5", "2026-01-01"}))
	})

	It("keeps a flag placed before positionals in front", func() {
		got := reorderFlagsFirst([]string{"--regions", "us", "5", "2026-01-01"})
		Expect(got).To(Equal([]string{"--regions", "us", "--", "5", "2026-01-01"}))
	})

	It("keeps a bool flag after positionals without consuming the next token", func() {
		got := reorderFlagsFirst([]string{"2018-02-14", "--regions", "us", "--informal"})
		Expect(got).To(Equal([]string{"--regions", "us", "--informal", "--", "2018-02-14"}))
	})

	It("returns flags unchanged when there are no positionals", func() {
		got := reorderFlagsFirst([]string{"--regions", "us"})
		Expect(got).To(Equal([]string{"--regions", "us"}))
	})

	It("treats a bare negative integer as positional, not a flag", func() {
		got := reorderFlagsFirst([]string{"-3", "2026-01-01", "--regions", "us"})
		Expect(got).To(Equal([]string{"--regions", "us", "--", "-3", "2026-01-01"}))
	})
})

var _ = Describe("cmdNext", func() {
	It("surfaces the library's count-must-be-positive error for a negative count", func() {
		err := cmdNext([]string{"-3", "2026-01-01", "--regions", "us"})
		Expect(err).To(MatchError("holidays.NextHolidays: count must be positive, got -3"))
	})

	It("still rejects a zero count with the library's error", func() {
		err := cmdNext([]string{"0", "2026-01-01", "--regions", "us"})
		Expect(err).To(MatchError("holidays.NextHolidays: count must be positive, got 0"))
	})

	It("still resolves a positive count with flags after the positionals", func() {
		err := cmdNext([]string{"5", "2026-01-01", "--regions", "us"})
		Expect(err).NotTo(HaveOccurred())
	})

	It("still resolves a positive count with flags before the positionals", func() {
		err := cmdNext([]string{"--regions", "us", "5", "2026-01-01"})
		Expect(err).NotTo(HaveOccurred())
	})
})
