package holidays

import "time"

type Holiday struct {
	Date    time.Time
	Name    string
	Regions []string
}
