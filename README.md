# temporal

[![CI](https://github.com/faustbrian/go-temporal/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-temporal/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-temporal/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-temporal.svg)](https://pkg.go.dev/github.com/faustbrian/go-temporal)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-temporal?sort=semver)](https://github.com/faustbrian/go-temporal/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`temporal` is immutable temporal algebra for Go. It provides bounded instant
and civil-date periods, explicit endpoint bounds, Allen relations, normalized
period sets, fixed elapsed durations, local times, circular daily intervals,
strict notation, versioned scalar/set encoding, PostgreSQL range adapters, and
bounded parsers and iterators.

The module is a Go-native successor to `github.com/faustbrian/go-temporal`. It
preserves useful mathematics without copying mutable or PHP-specific APIs.

## Install

```sh
go get github.com/faustbrian/go-temporal
```

The minimum supported toolchain is Go 1.26.6. Civil-date features use
`github.com/faustbrian/go-calendar`; clocks and timers deliberately remain in
`clock`.

## Five-minute instant period

```go
start := time.Date(2026, 7, 16, 9, 0, 0, 0, time.UTC)
period, err := instant.Range(start, start.Add(90*time.Minute))
if err != nil { log.Fatal(err) }

parts, err := period.SplitForward(30*time.Minute, temporal.Limits{Steps: 10})
if err != nil { log.Fatal(err) }

set, err := instant.NewSet(temporal.Limits{}, parts...)
if err != nil { log.Fatal(err) }

dayStart, err := instant.Snap(
    start, instant.Day, instant.Floor, time.UTC, timezone.Reject,
)
if err != nil { log.Fatal(err) }
fmt.Println(set.Includes(start), set.Len()) // true 1: normalization merged parts
```

`instant.Range` is closed-open. Use `instant.New` for another bound mode.
Instant duration is elapsed duration; calendar months are never inferred.

## Five-minute civil-date period

```go
month, err := dateperiod.Month(2026, time.July)
if err != nil { log.Fatal(err) }

weeks, err := month.SplitDays(7, temporal.Limits{Steps: 10})
if err != nil { log.Fatal(err) }

helsinki, _ := time.LoadLocation("Europe/Helsinki")
instants, err := month.ToInstant(helsinki, timezone.Reject)
if err != nil { log.Fatal(err) }
```

Date arithmetic delegates to `calendar`. Conversion uses the next civil
boundary as an exclusive instant end, so DST-short and DST-long days remain
correct.

## Five-minute daily interval

```go
start, _ := timeofday.Parse("22:00", temporal.Limits{})
end, _ := timeofday.Parse("02:00", temporal.Limits{})
night, err := timeofday.Between(start, end, temporal.ClosedOpen)
if err != nil { log.Fatal(err) }

set, _ := timeofday.NewIntervalSet(temporal.Limits{}, night)
offHours, err := set.Complement()
if err != nil { log.Fatal(err) }

helsinki, _ := time.LoadLocation("Europe/Helsinki")
date := calendar.MustDate(2026, time.July, 16)
instantNight, err := night.ToInstant(date, helsinki, timezone.Reject)
if err != nil { log.Fatal(err) }

fmt.Println(night.Kind(), night.Duration(), offHours.Len(), instantNight.Start())
// Circular 4h0m0s 1
```

`24:00` is a distinct end boundary. Equal endpoints must be constructed as
`Collapsed` or `FullDay`; they are never guessed.

## Packages

- `temporal`: bounds, Allen relations, typed errors, and resource limits.
- `instant`: bounded `time.Time` periods and normalized sets.
- `dateperiod`: bounded `calendar.Date` periods and normalized sets.
- `timeofday`: local times, fixed durations, circular intervals, and sets.
- `notation`: strict ISO 8601, ISO 80000, and Bourbaki codecs.
- `postgres`: loss-checked PostgreSQL range and multirange adapters.
- `temporalwire`: versioned format-neutral scalar and set documents.
- `temporalvalidation`: deterministic `validation` rules.
- `temporalconfig`: atomic `config` text wrappers.
- `temporaltest`: exhaustive relation fixtures and set assertions.

## Compatibility status

The audited PHP source is pinned at
`469603239dbe700739c29b4c532a90382b6cbedf`. The complete behavior inventory
has a machine-checked classification and evidence pointer for every non-chart
public symbol. Deliberate divergences are in
[docs/compatibility.md](docs/compatibility.md) and
[docs/migration.md](docs/migration.md). Generated compatibility evidence is
stored under `compat/fixtures`.

Charting is intentionally unsupported in v1. There is no `temporalchart`
package and no terminal/Gantt renderer in core. Every PHP chart type and option
is inventoried, and the core period/set values preserve a future renderer seam.
Full PHP-package compatibility is not claimed while this gap remains.

## Quality and local gates

```sh
make check
golib mutation --module .
golib docs check --module .
make -f verification/package.mk php-compat PHP_TEMPORAL_SOURCE=/path/to/php-temporal
```

See [docs/testing.md](docs/testing.md) for the evidence model and
[docs/hardening.md](docs/hardening.md) for the current algebra audit. See
[SECURITY.md](SECURITY.md) for hostile-input and disclosure guidance.
The [specification decision register](docs/specification-decisions.md) separates
standard-backed codec and PostgreSQL behavior from package API policy. This
library does not implement the Temporal service protocol.

## License

MIT. See [LICENSE](LICENSE).
