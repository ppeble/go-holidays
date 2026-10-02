# go-holidays

A Go library for working with statutory and other holidays.

All holiday definitions are maintained in the
[holidays/definitions](https://github.com/holidays/definitions) repository
(vendored here as a submodule). By default this library returns statutory
(formally government-defined) holidays. Culturally recognized but non-statutory
holidays (such as Valentine's Day) are available via the `Informal` option. See
the [definitions syntax guide](https://github.com/holidays/definitions/blob/master/doc/SYNTAX.md#formalinformal)
for details on how holidays are classified.

## Installation

```bash
go get github.com/holidays/go-holidays
```

## Tested versions

This module requires **Go 1.25+**, the floor declared by the `go` directive in
`go.mod`.

CI runs the full test suite against the latest Go minor release and the
previous two:

  * 1.27
  * 1.26
  * 1.25

That list lives in the `test-matrix` job in `.github/workflows/ci.yml` and moves
forward as new Go minors are released, with the floor in `go.mod` moving up in
step so every matrix entry is actually exercised (rather than silently
upgraded by `GOTOOLCHAIN=auto`).

## Semver

This module follows [semantic versioning](http://semver.org/). The guarantee
specifically covers the exported surface of the root package
`github.com/holidays/go-holidays` (its functions, types, and their fields).

Please note that we consider definition changes to be "minor" bumps, meaning
they are backwards compatible with your code but might give different holiday
results.

## Time zones

Dates are always constructed in UTC and truncated to the calendar day
(`time.Date(y, m, d, 0, 0, 0, 0, time.UTC)`). Pass in whatever time zone you
like: comparisons are done on the UTC calendar day, not on wall-clock time.

## Usage

Every example below assumes this import, shown once here and omitted from the
rest for brevity:

```go
import (
	"time"

	holidays "github.com/holidays/go-holidays"
)
```

This library offers multiple ways to check for holidays for a variety of
scenarios.

Lookups return a `[]holidays.Holiday`. Each one has a `Date`, a `Name` and the
`Regions` it applies to, which you can print like this (with `fmt` and
`strings` imported):

```go
for _, h := range hs {
	fmt.Printf("%s  %s  [%s]\n", h.Date.Format("2006-01-02"), h.Name, strings.Join(h.Regions, ", "))
}
```

The results under each example below are shown in that form, with the names
padded so the columns line up, or as `(no holidays)` when the result is empty.
`Regions` holds the regions listed on the matching definition, which can be
more than the ones you asked for. A `Holiday` also has an `Informal` field,
which is `true` for informal observances and can only be set when the lookup
used the `Informal` option.

`On`, `Between` and `YearHolidays` return holidays in no guaranteed order, so
sort them if the order matters. The results below are listed by date, then
name.

#### Checking a specific date

Get all holidays on April 25, 2008 in Australia:

```go
d := time.Date(2008, time.April, 25, 0, 0, 0, 0, time.UTC)
hs, err := holidays.On(d, holidays.Options{Regions: []string{"au"}})
```
```
2008-04-25  ANZAC Day  [au]
```

You can check multiple regions in a single call:

```go
d := time.Date(2008, time.January, 1, 0, 0, 0, 0, time.UTC)
hs, err := holidays.On(d, holidays.Options{Regions: []string{"us", "fr"}})
```
```
2008-01-01  Jour de l'an    [fr]
2008-01-01  New Year's Day  [us]
```

You can leave `Regions` empty to get holidays for any registered region:

```go
d := time.Date(2007, time.April, 25, 0, 0, 0, 0, time.UTC)
hs, err := holidays.On(d, holidays.Options{})
```
```
2007-04-25  ANZAC Day                       [au]
2007-04-25  ANZAC Day                       [au_qld, au_nt, au_act, au_sa]
2007-04-25  ANZAC Day                       [nz]
2007-04-25  Dia da Liberdade                [pt]
2007-04-25  Festa della Liberazione         [it]
2007-04-25  Festa di San Marco Evangelista  [it_ve]
```

#### Wildcard regions

A region code ending in an underscore (e.g. `au_`, `ca_`) is a *wildcard*. It
matches the parent country region and all of its sub-regions in a single call:

```go
d := time.Date(2017, time.March, 13, 0, 0, 0, 0, time.UTC)
hs, err := holidays.On(d, holidays.Options{Regions: []string{"au_"}})
```
```
2017-03-13  Canberra Day          [au_act]
2017-03-13  Eight Hours Day       [au_tas]
2017-03-13  Labour Day            [au_vic]
2017-03-13  March Public Holiday  [au_sa]
```

The same date queried with the plain `au` region returns nothing, because none
of those holidays are observed nation-wide:

```go
hs, err := holidays.On(d, holidays.Options{Regions: []string{"au"}})
```
```
(no holidays)
```

Use a wildcard when you want "this country and every sub-region it defines"
without listing each sub-region explicitly.

Note that a wildcard always collapses to the top-level country region. The
portion between the country prefix and the trailing underscore is ignored, so
`au_vic_` behaves identically to `au_` (it loads every Australian sub-region,
not just Victoria's). There is currently no way to wildcard-match only the
children of a sub-region.

#### Checking a date range

Get all holidays during the month of July 2008 in Canada and the US:

```go
from := time.Date(2008, time.July, 1, 0, 0, 0, 0, time.UTC)
to := time.Date(2008, time.July, 31, 0, 0, 0, 0, time.UTC)
hs, err := holidays.Between(from, to, holidays.Options{Regions: []string{"ca", "us"}})
```
```
2008-07-01  Canada Day        [ca]
2008-07-04  Independence Day  [us]
```

#### Informal holidays

Set `Options.Informal` to include holidays specified as informal in your
results. See the [definitions syntax guide](https://github.com/holidays/definitions/blob/master/doc/SYNTAX.md#formalinformal)
for what constitutes "informal" vs "formal".

By default this option is `false`, meaning no informal holidays are returned.

Get Valentine's Day in the US:

```go
d := time.Date(2018, time.February, 14, 0, 0, 0, 0, time.UTC)
hs, err := holidays.On(d, holidays.Options{Regions: []string{"us"}, Informal: true})
```
```
2018-02-14  Valentine's Day  [us, ca]
```

Leaving `Informal` false means Valentine's Day is not returned:

```go
hs, err := holidays.On(d, holidays.Options{Regions: []string{"us"}})
```
```
(no holidays)
```

#### Observed holidays

Set `Options.Observed` to include holidays that are observed on different days
than they actually occur. See the [definitions syntax guide](https://github.com/holidays/definitions/blob/master/doc/SYNTAX.md#observed)
for further explanation of "observed".

By default this option is `false`, meaning no observed logic is applied.

Get holidays that are observed on Monday July 2, 2007 in British Columbia,
Canada:

```go
d := time.Date(2007, time.July, 2, 0, 0, 0, 0, time.UTC)
hs, err := holidays.On(d, holidays.Options{Regions: []string{"ca_bc"}, Observed: true})
```
```
2007-07-02  Canada Day  [ca]
```

Leaving `Observed` false means "Canada Day" is not returned on July 2, since it
actually falls on Sunday July 1:

```go
hs, err := holidays.On(d, holidays.Options{Regions: []string{"ca_bc"}})
```
```
(no holidays)
```
```go
d = time.Date(2007, time.July, 1, 0, 0, 0, 0, time.UTC)
hs, err = holidays.On(d, holidays.Options{Regions: []string{"ca_bc"}})
```
```
2007-07-01  Canada Day  [ca]
```

#### Any holidays during work week

`AnyHolidaysDuringWorkWeek` reports whether any holiday matching the options
falls during the Monday-Friday work week containing the given date.

Check whether a holiday falls during the first week of the year for any
region:

```go
d := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
any, err := holidays.AnyHolidaysDuringWorkWeek(d, holidays.Options{})
// any == true
```

`Informal` and `Observed` apply the same way they do everywhere else:

```go
// true: Valentine's Day falls on a Wednesday
d = time.Date(2018, time.February, 14, 0, 0, 0, 0, time.UTC)
any, _ = holidays.AnyHolidaysDuringWorkWeek(d, holidays.Options{Regions: []string{"us"}, Informal: true})

// false without Informal
any, _ = holidays.AnyHolidaysDuringWorkWeek(d, holidays.Options{Regions: []string{"us"}})

// true: Veterans Day is observed on Monday November 12, 2018
d = time.Date(2018, time.November, 12, 0, 0, 0, 0, time.UTC)
any, _ = holidays.AnyHolidaysDuringWorkWeek(d, holidays.Options{Regions: []string{"us"}, Observed: true})

// false without Observed: the actual holiday is on Sunday November 11
any, _ = holidays.AnyHolidaysDuringWorkWeek(d, holidays.Options{Regions: []string{"us"}})
```

#### Next holidays

`NextHolidays` returns the next `count` holidays on or after a given date,
sorted by date ascending:

```go
from := time.Date(2016, time.February, 23, 0, 0, 0, 0, time.UTC)
hs, err := holidays.NextHolidays(from, 3, holidays.Options{Regions: []string{"us"}, Informal: true})
```
```
2016-03-17  St. Patrick's Day  [us, ca]
2016-03-25  Good Friday        [us]
2016-03-27  Easter Sunday      [us]
```

#### Year holidays

`YearHolidays` returns every holiday matching the options in a given calendar
year:

```go
hs, err := holidays.YearHolidays(2016, holidays.Options{Regions: []string{"ca_on"}})
```
```
2016-01-01  New Year's Day  [ca]
2016-02-15  Family Day      [ca_on]
2016-03-25  Good Friday     [ca]
2016-05-23  Victoria Day    [ca_ab, ca_bc, ca_mb, ca_nt, ca_nu, ca_on, ca_sk, ca_yt]
2016-07-01  Canada Day      [ca]
2016-09-05  Labour Day      [ca]
2016-10-10  Thanksgiving    [ca_ab, ca_bc, ca_mb, ca_nt, ca_nu, ca_on, ca_qc, ca_sk, ca_yt]
2016-12-25  Christmas Day   [ca_on]
2016-12-26  Boxing Day      [ca_on]
```

`YearHolidaysFrom` returns every holiday from a given date through December 31
of that date's year, sorted ascending:

```go
from := time.Date(2016, time.February, 23, 0, 0, 0, 0, time.UTC)
hs, err := holidays.YearHolidaysFrom(from, holidays.Options{Regions: []string{"ca_on"}})
```
```
2016-03-25  Good Friday    [ca]
2016-05-23  Victoria Day   [ca_ab, ca_bc, ca_mb, ca_nt, ca_nu, ca_on, ca_sk, ca_yt]
2016-07-01  Canada Day     [ca]
2016-09-05  Labour Day     [ca]
2016-10-10  Thanksgiving   [ca_ab, ca_bc, ca_mb, ca_nt, ca_nu, ca_on, ca_qc, ca_sk, ca_yt]
2016-12-25  Christmas Day  [ca_on]
2016-12-26  Boxing Day     [ca_on]
```

#### Available regions

`AvailableRegions` returns every registered region code, sorted
lexicographically:

```go
regions := holidays.AvailableRegions()
fmt.Println(len(regions), regions[:5])
```
```
290 [ar at au au_act au_nsw]
```

#### Region display names

`RegionName` returns the display name for one region code, and whether it is
registered. `RegionNames` returns every registered region code mapped to its
display name:

```go
name, ok := holidays.RegionName("gb_sct")
fmt.Println(name, ok)

names := holidays.RegionNames()
fmt.Println(len(names), names["ch_ge"])
```
```
Scotland true
290 Genève
```

## Command-line interface

`make build` produces `bin/holidays` and `bin/gen-holidays`. `bin/holidays`
wraps the same public API from the command line:

```bash
bin/holidays on 2024-07-04 --regions us
```
```
2024-07-04  Independence Day  us
```
```bash
bin/holidays between 2024-12-20 2024-12-31 --regions us
```
```
2024-12-25  Christmas Day  us
```
```bash
bin/holidays year 2024 --regions us
```
```
2024-01-01  New Year's Day                        us
2024-01-15  Martin Luther King, Jr. Day           us
2024-02-19  Presidents' Day                       us
2024-05-27  Memorial Day                          us
2024-06-19  Juneteenth National Independence Day  us
2024-07-04  Independence Day                      us
2024-09-02  Labor Day                             us
2024-11-11  Veterans Day                          us
2024-11-28  Thanksgiving                          us
2024-12-25  Christmas Day                         us
```
```bash
bin/holidays next 5 2024-05-28 --regions us
```
```
2024-06-19  Juneteenth National Independence Day  us
2024-07-04  Independence Day                      us
2024-09-02  Labor Day                             us
2024-11-11  Veterans Day                          us
2024-11-28  Thanksgiving                          us
```
```bash
bin/holidays workweek 2024-11-25 --regions us
```
```
true
```
```bash
bin/holidays regions
```
```
ar
at
au
au_act
au_nsw
au_nt
au_qld
au_qld_brisbane
au_qld_cairns
au_sa
... (280 more)
```

Holidays print one per line as date, name and comma-separated regions,
separated by tabs (shown aligned here) and sorted by date, then name.

Every subcommand except `regions` also accepts `--informal` and `--observed`.
Flags may appear before or after the positional arguments.

## Loading custom definitions on the fly

In addition to the [provided definitions](https://github.com/holidays/definitions)
you can load a custom definitions file on the fly and use it immediately.

To load a custom "Company Founding" holiday on June 1st, put this in
`custom_holidays.yaml`:

```yaml
months:
  6:
  - name: Company Founding
    regions: [my_custom_region]
    mday: 1
```

Then load it and query it by the region code from its `regions:` list:

```go
err := holidays.LoadCustom("/home/user/holiday_definitions/custom_holidays.yaml")
hs, err := holidays.On(time.Date(2013, time.June, 1, 0, 0, 0, 0, time.UTC),
    holidays.Options{Regions: []string{"my_custom_region"}})
```
```
2013-06-01  Company Founding  [my_custom_region]
```

Custom definition files must match the [syntax of the existing definition files](https://github.com/holidays/definitions/blob/master/doc/SYNTAX.md).
Region codes are lowercased and trimmed, the same as `Options.Regions`.

Multiple files can be loaded at the same time:

```go
err := holidays.LoadCustom(
    "/home/user/holidays/custom_holidays1.yaml",
    "/home/user/holidays/custom_holidays2.yaml",
)
```

Each file is keyed by its base name without the directory or extension. Loading
the same path again replaces its prior load, and so does loading any other file
with the same base name (`/a/holidays.yaml` and `/b/holidays.yaml`, or
`holidays.yaml` and `holidays.yml`), so give every file a distinct base name.
Files with different base names add rules without overwriting one another.

`LoadCustom` is all-or-nothing: if any file fails to read, parse, or validate,
it returns an error and none of the files are registered. `UnloadCustom`
removes rules previously loaded from the given paths, matched by base name.
Both call `ResetCache` when they succeed, since the rule set changed.

### Custom methods

Custom rules can't embed executable logic in YAML: a rule's `function:` or
`observed:` reference must point at a method already registered in Go before
you call `LoadCustom`, or `LoadCustom` returns an error:

```go
holidays.RegisterMethod("my_method", func(a holidays.MethodArgs) (time.Time, error) {
    return time.Date(a.Year, time.March, 15, 0, 0, 0, 0, time.UTC), nil
})
err := holidays.LoadCustom("my_team.yaml") // YAML can now use function: my_method(year)
```

The method receives only a `MethodArgs`. The parentheses in the YAML call are
required, but any arguments inside them (the `year` above) are not passed to
your function.

- `function:` methods return the holiday's date for `a.Year`. Only the month and
  day of the result are kept, after the rule's `function_modifier:` days are
  added: the year is forced to the year being resolved and the time zone to UTC.
  Returning the zero `time.Time` means the holiday does not occur that year.
  `a.Month` and `a.Day` are the rule's month and mday, `a.Date` is the date
  computed from them (or the rule's wday/week), and `a.Region` is the first
  region listed on the rule.
- `observed:` methods run only when `Options.Observed` is true. `a.Date` is the
  holiday's date, `a.Region` is the first region in `Options.Regions` (empty if
  none), and the date you return replaces the holiday's date.
- `RegisterMethod` panics if the name is already registered, including any
  built-in method name, and a registered method can't be removed. Register each
  name once, for example from an `init` function.

## Caching holiday lookups

If you are checking holidays regularly you can cache your results for improved
performance. Run this before looking up a holiday (e.g. during startup):

```go
year := 365 * 24 * time.Hour
err := holidays.CacheBetween(time.Now(), time.Now().Add(2*year), holidays.Options{
    Regions:  []string{"ca", "us"},
    Observed: true,
})
```

Holidays for the regions and options specified within the dates specified will
be pre-computed and stored in memory. Subsequent `On`/`Between` calls with the
same options and a range fully contained in the cached window are answered
from the cache. `ResetCache` clears every cached range.

## How to contribute

See [CONTRIBUTING.md](CONTRIBUTING.md) for information on how to help out.

## Credits

* Holiday definitions come from [holidays/definitions](https://github.com/holidays/definitions),
  along with all of its wonderful contributors.
