package engine

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ppeble/go-holidays/pkg/definition"
)

// MethodArgs is the uniform invocation payload for every registered method.
// The resolver populates every applicable field; each method picks what it needs.
type MethodArgs struct {
	Year   int
	Month  int
	Day    int
	Date   time.Time
	Region string
}

// Method computes or transforms a date.
type Method func(args MethodArgs) (time.Time, error)

var (
	methodMu sync.RWMutex
	methods  = map[string]Method{}

	regionMu  sync.RWMutex
	countries = map[string][]definition.HolidayRule{} // keyed by country code (us, gb, ...)
)

// RegisterMethod registers a named method. Panics if name is already registered.
func RegisterMethod(name string, fn Method) {
	methodMu.Lock()
	defer methodMu.Unlock()
	if _, exists := methods[name]; exists {
		panic(fmt.Sprintf("engine: method %q already registered", name))
	}
	methods[name] = fn
}

// LookupMethod returns the registered method for name.
func LookupMethod(name string) (Method, bool) {
	methodMu.RLock()
	defer methodMu.RUnlock()
	fn, ok := methods[name]
	return fn, ok
}

// IsMethodRegistered reports whether name has been registered.
func IsMethodRegistered(name string) bool {
	_, ok := LookupMethod(name)
	return ok
}

// RegisterCountry installs the rules for a country YAML (e.g. "us" from us.yaml).
// Subregions like us_ga are addressable via these same rules — each rule carries
// its own Regions slice.
func RegisterCountry(country string, rules []definition.HolidayRule) {
	regionMu.Lock()
	defer regionMu.Unlock()
	countries[country] = rules
}

// AvailableRegions returns every region code mentioned across all registered rules,
// sorted lexicographically.
func AvailableRegions() []string {
	regionMu.RLock()
	defer regionMu.RUnlock()
	seen := map[string]struct{}{}
	for _, rules := range countries {
		for _, r := range rules {
			for _, reg := range r.Regions {
				seen[reg] = struct{}{}
			}
		}
	}
	out := make([]string, 0, len(seen))
	for r := range seen {
		out = append(out, r)
	}
	sort.Strings(out)
	return out
}

// rulesFor returns the rules to scan given requested regions. Since registered
// country keys do not always align with region prefixes (e.g. be_fr.yaml
// registers under "be_fr", not "be"), we walk every registered country's
// rules and let ruleMatchesRequested decide.
func rulesFor(requested []string) []definition.HolidayRule {
	regionMu.RLock()
	defer regionMu.RUnlock()
	var out []definition.HolidayRule
	for _, rs := range countries {
		out = append(out, rs...)
	}
	return out
}

// rulesForCountry returns the rules registered for the named country file.
// Method implementations call this at resolve time, after all init() functions
// have run, to inspect peer rules (e.g. JP substitute-holiday logic needs to
// know which days are already holidays).
func rulesForCountry(country string) []definition.HolidayRule {
	regionMu.RLock()
	defer regionMu.RUnlock()
	return countries[country]
}

// ruleMatchesRequested reports whether a rule applies to any of the requested regions.
// Matching rules:
//   - Exact equality: rule region == requested region.
//   - Parent: rule region "us" applies when the user requests "us_ga".
//   - Wildcard: requested region ending with "_" (e.g. "gb_") matches any rule
//     region that equals the bare prefix ("gb") or starts with the prefix
//     plus "_" ("gb_sct", "gb_wls").
func ruleMatchesRequested(rule definition.HolidayRule, requested []string) bool {
	if len(requested) == 0 {
		return true
	}
	for _, req := range requested {
		if strings.HasSuffix(req, "_") {
			bare := strings.TrimSuffix(req, "_")
			for _, rr := range rule.Regions {
				if rr == bare || strings.HasPrefix(rr, req) {
					return true
				}
			}
			continue
		}
		for _, rr := range rule.Regions {
			if rr == req || isParentOf(rr, req) {
				return true
			}
		}
	}
	return false
}

func isParentOf(parent, child string) bool {
	return len(child) > len(parent) && strings.HasPrefix(child, parent+"_")
}
