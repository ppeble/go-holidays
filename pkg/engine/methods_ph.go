package engine

import (
	"time"

	"github.com/ppeble/go-holidays/pkg/calc"
)

func init() {
	// Philippines National Heroes Day: last Monday of August.
	RegisterMethod("ph_heroes_day", func(a MethodArgs) (time.Time, error) {
		return calc.DayOfMonth(a.Year, time.August, -1, time.Monday), nil
	})
}
