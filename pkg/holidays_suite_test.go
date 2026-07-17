package holidays_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestHolidays is the single Ginkgo bootstrap for the whole pkg/ test binary.
// The binary compiles two test packages, holidays (white-box, cache_internal_test.go)
// and holidays_test (black-box, the rest), that share one process-global Ginkgo
// spec tree. Exactly one RunSpecs call may exist for the binary; a second one
// would panic or double-run the suite. Do not add another RunSpecs anywhere
// else in pkg/.
func TestHolidays(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "holidays suite")
}
