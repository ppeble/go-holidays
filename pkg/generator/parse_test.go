package generator

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Verifies that a YAML containing a top-level region_names: block (added in
// v8.0.0 of the upstream definitions) does not cause ParseRegionFile to error
// when strict decoding is enabled.
var _ = Describe("ParseRegionFile", func() {
	Context("when the YAML has a top-level region_names block", func() {
		It("parses without error and keeps the rule intact", func() {
			yaml := []byte(`
months:
  1:
    - name: "New Year's Day"
      regions:
        - xx
      mday: 1
      type: formal

region_names:
  xx: "Example Country"
  xx_reg: "Example Region"

tests:
  - given:
      date: "2024-01-01"
      regions:
        - xx
    expect:
      name: "New Year's Day"
      holiday: true
`)

			rf, err := ParseRegionFile("xx", yaml)
			Expect(err).NotTo(HaveOccurred())
			Expect(rf.Rules).To(HaveLen(1))
			Expect(rf.Rules[0].Name).To(Equal("New Year's Day"))
		})
	})
})
