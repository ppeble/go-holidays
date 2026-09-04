package generator

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmitDefinitions region names", func() {
	Context("when the region file has region_names", func() {
		It("emits a sorted map var and registers it in init()", func() {
			rf := &RegionFile{
				Country: "xx",
				RegionNames: map[string]string{
					"xx_reg": "Example Region",
					"xx":     "Example Country",
				},
			}

			src, err := EmitDefinitions(rf)
			Expect(err).NotTo(HaveOccurred())
			out := string(src)

			Expect(out).To(ContainSubstring(`var xxRegionNames = map[string]string{`))
			// sorted by key: "xx" before "xx_reg"
			xxIdx := indexOf(out, `"Example Country"`)
			xxRegIdx := indexOf(out, `"Example Region"`)
			Expect(xxIdx).To(BeNumerically(">=", 0))
			Expect(xxRegIdx).To(BeNumerically(">", xxIdx))

			Expect(out).To(ContainSubstring(`engine.RegisterCountry("xx", xxRules)`))
			Expect(out).To(ContainSubstring(`engine.RegisterRegionNames("xx", xxRegionNames)`))
		})
	})

	Context("when the region file has no region_names", func() {
		It("does not emit a RegionNames var or RegisterRegionNames call", func() {
			rf := &RegionFile{Country: "xx"}

			src, err := EmitDefinitions(rf)
			Expect(err).NotTo(HaveOccurred())
			out := string(src)

			Expect(out).NotTo(ContainSubstring("RegionNames"))
			Expect(out).NotTo(ContainSubstring("RegisterRegionNames"))
		})
	})
})

func indexOf(haystack, needle string) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}
