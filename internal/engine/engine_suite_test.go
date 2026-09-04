package engine

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestEngine is the single Ginkgo bootstrap for the internal/engine test
// binary. It drives specs from both this package (white-box, so specs can
// reach unexported helpers like ruleMatchesRequested directly) and the
// engine_test package (black-box, for the exported public surface) in the
// same directory: Ginkgo's Describe/It registration is global regardless of
// which of the two source packages declares it, so one RunSpecs call here
// is enough to run everything.
func TestEngine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "engine suite")
}
