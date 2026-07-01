package generator

import (
	"testing"
)

// TestParseRegionFile_RegionNames verifies that a YAML containing a top-level
// region_names: block (added in v8.0.0 of the upstream definitions) does not
// cause ParseRegionFile to error when strict decoding is enabled.
func TestParseRegionFile_RegionNames(t *testing.T) {
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
	if err != nil {
		t.Fatalf("ParseRegionFile returned unexpected error: %v", err)
	}
	if len(rf.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rf.Rules))
	}
	if rf.Rules[0].Name != "New Year's Day" {
		t.Errorf("expected rule name %q, got %q", "New Year's Day", rf.Rules[0].Name)
	}
}
