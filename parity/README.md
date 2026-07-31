# Parity suite: Ruby gem vs. Go port

This directory holds an in-repo comparison suite that checks the Go port
returns the same results as the upstream Ruby `holidays` gem for every
result-producing public function.

## Design: apples-to-apples

The challenge is comparing two engines fairly: the gem carries its own bundled
data set, which we do not want to depend on. We close that gap by separating
CODE from DATA:

- **CODE** is the installed gem (`holidays`, pinned to 11.3.0 via `Gemfile`): its
  real resolution logic.
- **DATA** is our `definitions/` submodule (v8.1.0 YAML), loaded into the running
  gem on startup via `Holidays.load_custom`, so the comparison never rides on
  whatever data the gem happens to bundle.

So both sides resolve the *same* holiday rules. Any difference in output is then
a real behavioural difference between the engines, not a difference in the data.

## How it works

```
parity_test.go  --NDJSON request-->  oracle.rb  (gem 11.3.0 + our v8.1.0 YAML)
   (Go side)     <--NDJSON result--             
```

- `oracle.rb` runs the gem, loads our v8.1.0 YAML, and answers one
  line-delimited JSON (NDJSON) request per line on stdin, one response per line
  on stdout (see its RUN CONTRACT header).
- `oracle_client.go` starts the oracle as a subprocess and speaks that contract.
- `parity_test.go` (build tag `parity`) runs a corpus of cases through both the
  Go API and the oracle and diffs the results. Each case runs under all four
  flag combinations: plain, observed, informal, informal+observed.
- `oracle_smoke.sh` is a standalone smoke check of the oracle.

## Running

```
make parity                      # or: go test -tags parity ./parity/...
```

The suite is behind the `parity` build tag, so plain `make test` never builds or
runs it and stays Ruby-free.

**Prerequisites:**
- Ruby with the pinned `holidays` gem installed. The oracle activates
  `parity/Gemfile` via `bundler/setup`, so `holidays 11.3.0` must be available
  (`gem install holidays -v 11.3.0`, or `bundle install` from `parity/`).
- The `definitions/` submodule checked out (`git submodule update --init
  definitions`). The oracle loads its v8.1.0 YAML; without it the gem falls back
  to its own bundled data and regions diverge.

## Scope

The eight result-producing public functions are compared:

| Ruby gem | Go port |
|----------|---------|
| `on` | `On` |
| `between` | `Between` |
| `cache_between` | `CacheBetween` |
| `next_holidays` | `NextHolidays` |
| `year_holidays` | `YearHolidaysFrom` |
| `any_holidays_during_work_week?` | `AnyHolidaysDuringWorkWeek` |
| `available_regions` | `AvailableRegions` |
| `load_custom` | `LoadCustom` |

`year_holidays` is compared against `YearHolidaysFrom(from, opts)` (the from-date
variant that matches the gem's semantics), not `YearHolidays(year)`.

### `load_all` is intentionally not ported (N/A)

The gem's `load_all` eagerly loads every region at runtime. The Go port already
loads all regions eagerly via `init()` side effects (a blank import of
`pkg/definitions`), so `load_all` would be a no-op with nothing to do. It is
documented here as not-applicable rather than implemented.

## Known limitations and intentional divergences

### `jp` is skipped (oracle limitation, not a Go bug)

The gem resolves some regions (notably Japan) through Ruby-coded custom methods
that live in a module the gem only defines for its own bundled data.
`load_custom` of our YAML cannot supply `Holidays::JP`, so any `jp` request
raises `uninitialized constant Holidays::JP`. The harness detects this
(`isOracleUnsupported`) and skips such cases, logging each skip. The Go port
itself resolves `jp` fine; only the oracle cannot serve it.

### `next_holidays` for `de` from 2024-03-01 (intentional Go divergence)

The gem's `next_holidays` windows its search through a `dates_driver` that
buckets each holiday by its *source month* (function/variable holidays such as
Easter offsets live in "month 0") and only reaches out to `from >> 12`. For this
case the 2025 bucket ends at month 4, so the gem **drops** the fixed-date
`Tag der Arbeit` (2025-05-01, month 5) while still **keeping** the Easter-based
`Christi Himmelfahrt` (2025-05-29, month 0) it computes for that year.

Go's `NextHolidays` instead expands year by year until it has `count` holidays,
so it correctly includes `Tag der Arbeit 2025`. We deliberately do **not**
replicate the gem's month-bucketing quirk (Go's behaviour is arguably more
correct), so this specific case is excluded from the asserted corpus with an
explanatory comment in `parity_test.go`.
