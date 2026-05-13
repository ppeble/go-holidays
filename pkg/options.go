package holidays

// Options controls a holiday lookup. An empty Regions slice means "all registered regions".
type Options struct {
	Regions  []string
	Informal bool
	Observed bool
}
