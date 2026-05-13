package engine

import "time"

func init() {
	// Epiphany (Trzech Króli) became a Polish public holiday in 2011.
	// The rule is placed under month 1 in pl.yaml; we use rule.Month from MethodArgs.
	RegisterMethod("pl_trzech_kroli", func(a MethodArgs) (time.Time, error) {
		if a.Year < 2011 {
			return time.Time{}, nil
		}
		return time.Date(a.Year, time.Month(a.Month), 6, 0, 0, 0, 0, time.UTC), nil
	})
	// Same date treated as informal before 2011.
	RegisterMethod("pl_trzech_kroli_informal", func(a MethodArgs) (time.Time, error) {
		if a.Year >= 2011 {
			return time.Time{}, nil
		}
		return time.Date(a.Year, time.Month(a.Month), 6, 0, 0, 0, 0, time.UTC), nil
	})
}
