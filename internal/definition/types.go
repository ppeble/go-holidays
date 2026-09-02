package definition

type HolidayType int

const (
	Formal HolidayType = iota
	Informal
)

type YearRangeKind int

const (
	YearRangeAny YearRangeKind = iota
	YearRangeUntil
	YearRangeFrom
	YearRangeLimited
	YearRangeBetween
)

type YearRange struct {
	Kind  YearRangeKind
	Years []int
}

type HolidayRule struct {
	Name             string
	Regions          []string
	Month            int
	Mday             int
	Wday             int
	Week             int
	Function         string
	FunctionArgs     []string
	FunctionModifier int
	Observed         string
	Type             HolidayType
	YearRanges       []YearRange
}

func (r HolidayRule) HasWday() bool     { return r.Week != 0 }
func (r HolidayRule) HasMday() bool     { return r.Mday != 0 }
func (r HolidayRule) HasFunction() bool { return r.Function != "" }
func (r HolidayRule) HasObserved() bool { return r.Observed != "" }

func (r HolidayRule) AppliesIn(year int) bool {
	if len(r.YearRanges) == 0 {
		return true
	}
	for _, yr := range r.YearRanges {
		if yr.Matches(year) {
			return true
		}
	}
	return false
}

func (yr YearRange) Matches(year int) bool {
	switch yr.Kind {
	case YearRangeUntil:
		return len(yr.Years) == 1 && year <= yr.Years[0]
	case YearRangeFrom:
		return len(yr.Years) == 1 && year >= yr.Years[0]
	case YearRangeLimited:
		for _, y := range yr.Years {
			if y == year {
				return true
			}
		}
		return false
	case YearRangeBetween:
		return len(yr.Years) == 2 && year >= yr.Years[0] && year <= yr.Years[1]
	default:
		return true
	}
}
