# go-holidays

A Go port of the top-level public API of the Ruby [`holidays`](https://github.com/holidays/holidays) gem, driven by the same upstream YAML definitions from [`holidays/definitions`](https://github.com/holidays/definitions) (vendored as a submodule pinned to a specific tag).

## Status

| Item | State |
|------|-------|
| Top-level API | `On`, `Between`, `YearHolidays`, `NextHolidays`, `AnyHolidaysDuringWorkWeek`, `CacheBetween`, `ResetCache`, `LoadCustom`, `UnloadCustom`, `RegisterMethod`, `AvailableRegions` |
| Upstream definitions | Vendored as a submodule, pinned to the tag in `definitions/VERSION.txt` |

## Layout

```
cmd/
  holidays/        # CLI wrapper around the public API
  gen-holidays/    # YAML -> Go code generator
pkg/
  holidays.go      # package holidays — public On/Between/YearHolidays
  calc/            # Easter, day-of-month, observance helpers, lunar stub
  definition/      # internal types (HolidayRule, YearRange)
  engine/          # method + region registries, year-range eval, resolver,
                   #   builtins.go and one methods_<country>.go per country
                   #   that uses custom Ruby methods we've hand-ported
  generator/       # YAML parser + Go emitter + test emitter
  definitions/     # GENERATED: one <country>.go + <country>_test.go per region
definitions/       # git submodule -> holidays/definitions (pinned tag in VERSION.txt)
```

## Install

Requires Go 1.22+. Clone with the submodule:

```bash
git clone --recurse-submodules https://github.com/ppeble/go-holidays
# or, in an existing clone:
git submodule update --init
```

## Library use

```go
import (
    "time"
    holidays "github.com/ppeble/go-holidays/pkg"
    _ "github.com/ppeble/go-holidays/pkg/definitions" // wire up region registrations
)

func example() {
    d, _ := time.Parse("2006-01-02", "2024-07-04")
    hs, err := holidays.On(d, holidays.Options{Regions: []string{"us"}})
    // hs[0].Name == "Independence Day"
    _ = err
}
```

`Options` fields:
- `Regions []string` — empty means "all registered regions". Trailing `_` is a subregion wildcard (`gb_` matches `gb_eng`, `gb_sct`, `gb_wls`, ...).
- `Informal bool` — include `type: informal` holidays.
- `Observed bool` — apply the rule's `observed:` transformation (e.g. roll weekend dates to Monday).

## CLI

```bash
make build  # produces bin/holidays and bin/gen-holidays

bin/holidays year 2024 --regions us
bin/holidays on 2024-07-04 --regions us
bin/holidays between 2024-12-20 2024-12-31 --regions us
bin/holidays next 5 2024-05-28 --regions us
bin/holidays workweek 2024-11-25 --regions us
bin/holidays regions
```

Flags can appear before or after positional arguments.

## Custom holiday definitions

`LoadCustom` parses YAML files in the same format as upstream `holidays/definitions` and registers their rules at runtime alongside the built-in ones:

```go
holidays.LoadCustom("my_team.yaml")
hs, _ := holidays.On(date, holidays.Options{Regions: []string{"my_team"}})
```

The YAML's `methods:` block (Ruby source) is parsed but ignored — we can't interpret Ruby. Any `function:` or `observed:` reference in a custom rule must point to a method registered in Go first:

```go
holidays.RegisterMethod("my_method", func(a holidays.MethodArgs) (time.Time, error) {
    return time.Date(a.Year, time.March, 15, 0, 0, 0, 0, time.UTC), nil
})
holidays.LoadCustom("my_team.yaml") // YAML can now use function: my_method(year)
```

Loading the same path again replaces the prior load; loading distinct paths adds rules. `UnloadCustom(path)` removes them.

## Regenerating definitions

The generated `pkg/definitions/*.go` files are checked in. Regenerate after bumping the submodule or after adding new per-region methods:

```bash
make generate                                       # all regions, fail on unported methods
bin/gen-holidays --allow-unported                   # all regions, skip those with missing methods
bin/gen-holidays -regions us,gb                     # only specific regions
```

The generator hard-fails on any `function:` or `observed:` YAML reference that lacks a registered Go implementation. To add a missing method:

1. Write the Go function in `pkg/engine/methods_<country>.go`.
2. Register it in that file's `init()` with `engine.RegisterMethod("<name>", func(a MethodArgs) (time.Time, error) { ... })`.
3. Re-run `make generate`.

## Bumping the definitions tag

```bash
make update-definitions DEFS_TAG=v6.2.0
git add definitions
git commit -s -m "vendor holidays/definitions v6.2.0"
make generate
```

## Tests

```bash
make test          # go vet + go test ./...
go test ./pkg/definitions/... -v
```

Each region's YAML `tests:` block is emitted as a corresponding `_test.go` so failures point at the specific holiday.
