package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestHolidaysCLI is the single Ginkgo bootstrap for the cmd/holidays test
// binary. It is white-box (package main) so specs can exercise unexported
// helpers such as reorderFlagsFirst directly.
func TestHolidaysCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cmd/holidays suite")
}
