package holidays

import (
	"time"

	"github.com/ppeble/go-holidays/internal/engine"
)

type Holiday struct {
	Date    time.Time
	Name    string
	Regions []string
}

// Options controls a holiday lookup. An empty Regions slice means "all registered regions".
type Options struct {
	Regions  []string
	Informal bool
	Observed bool
}

// MethodArgs re-exports engine.MethodArgs so callers using RegisterMethod do
// not need to import internal/engine.
type MethodArgs = engine.MethodArgs
