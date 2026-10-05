# Migration from faustbrian/temporal

## Go major-version migration

Use `github.com/faustbrian/go-temporal/v2` and add `/v2` before the package
suffix in every Temporal import. Update directly composed dependencies and
consumers together: a v1 `instant.Period`, `temporal.Limits`, or facade value is
not assignable to its v2 counterpart. Sentinel names and `errors.Is` semantics
remain, but v1 and v2 sentinel values are distinct; classify with the sentinel
from the major that produced the error. Calendar types exposed by civil-date
and timezone APIs must match the Calendar major selected by the Temporal
module's `go.mod`. Temporal v2 uses the published
`github.com/faustbrian/go-calendar/v2` v2.0.0 contract. Update Calendar imports
to `/v2`, including `/v2/timezone`; Calendar v1 dates, local-time values,
policies, and sentinels do not have the same identity as Calendar v2 values.
Historical Temporal v1 consumers remain on their published dependency graph.

Temporal v2's canonical `adapters/validation` and retained `temporalvalidation`
constructors return validators from the published
`github.com/faustbrian/go-validation/v2` v2.0.0 module. Compose them with that
major's `Validator`, `Context`, `Report`, and sentinels; Validation v1 types are
not interchangeable. Nonempty and inclusive-range behavior, including
Temporal's `ErrReversed`, remains unchanged. Retaining an adapter package path
within Temporal v2 does not retain the Validation v1 dependency identity.

Direct Opening and Rule Engine adapter adoption requires separate consumer
migration and verification; this source change does not certify those routes.

The `temporal/v1` wire schema is independent of the Go module major and remains
unchanged, including empty collection normalization to JSON `null`. Retained
facade packages remain distinct nominal APIs within `/v2`; prefer canonical
adapters for new code. Public tags and releases establish availability.

## PHP migration

The Go API is deliberately not a transliteration. Start from represented-set
behavior, then select the type whose semantics match.

| PHP concept | Go replacement | Deliberate change |
|---|---|---|
| `Period\Bounds` | `temporal.Bounds` | immutable enum; closed-open is zero/default |
| `DatePoint` | `time.Time` or `calendar.Date` | instant and civil date are distinct |
| `Period\Duration` | `time.Duration` / `timeofday.Duration` | no implicit calendar months |
| `Period` | `instant.Period` or `dateperiod.Period` | typed immutable values |
| `Sequence` | normalized `instant.Set` / `dateperiod.Set` | no mutation collection API |
| `Time\Time` | `timeofday.Time` | nanoseconds; strict ASCII ISO text |
| `Time\Duration` | `timeofday.Duration` | checked `time.Duration` interoperability |
| `Time\Interval` | `timeofday.Interval` | all four bounds; explicit collapsed/full |
| `Time\IntervalSet` | `timeofday.IntervalSet` | normalized disjoint immutable segments |

PHP accepts whitespace in mathematical notation; Go's strict codecs do not.
Trim only at a trusted application boundary if the protocol permits it. PHP's
local-time precision is microseconds; Go supports nanoseconds and rejects
lossy PostgreSQL writes.

PHP `Time::endOfDay()` denotes the last microsecond
(`23:59:59.999999`). Go `EndOfDay()` is the distinct boundary `24:00`. Use the
PHP last-microsecond value only when reproducing legacy sampled membership; use
`24:00` for interval boundaries.

PHP represents collapsed and full-day intervals with equal formatted endpoints
and distinguishes them through duration/type. Go rejects equal `Between`
endpoints and requires `Collapsed(anchor)` or `FullDay()`.

Date factories move to `dateperiod` and delegate civil arithmetic to
`calendar`. End-of-day conversion becomes a next-boundary exclusive instant
range, preserving DST behavior.

PHP snap helpers become `instant.Snap` or `Period.SnapOutward`; callers provide
the unit, direction, location, and `calendar` gap/fold policy. PHP
`Time::applyTo` becomes `Time.Apply`, and circular `Interval::toNative` becomes
`Interval.ToInstant`; both require the same explicit civil context.

Unversioned `jsonSerialize` payloads become `temporalwire.Document` for scalar
values and `CollectionDocument` for normalized sets. Decoders reject unknown
fields and trailing data instead of accepting partially understood payloads.

All variable-output operations accept `temporal.Limits` and may return
`LimitError`. All arithmetic and parsing errors are typed and compatible with
`errors.Is`/`errors.As`.

## Adapter import migration

| Retained compatibility import | Canonical import |
|---|---|
| `go-temporal/v2/temporalconfig` | `go-temporal/v2/adapters/config` |
| `go-temporal/v2/postgres` | `go-temporal/v2/adapters/postgres` |
| `go-temporal/v2/temporalvalidation` | `go-temporal/v2/adapters/validation` |
| `go-temporal/v2/temporalwire` | `go-temporal/v2/adapters/wire` |

Change imports when convenient; compatibility signatures and named-type identities
remain supported within the same module major. The compatibility paths remain for the
longer of 180 days after successor public availability and two subsequently
published stable root-module minor releases. Removal additionally requires an
authorized next major release and consumer verification.

## Unsupported charting gap

`Period\Chart` has no core Go implementation. This includes `Chart`, `Data`,
`Dataset`, `GanttChart`, `GanttChartConfig`, `Output`, `StreamOutput`, terminal
capabilities, colors, alignments, affix/reverse/generated labels, decimal,
Latin-letter and Roman-number labels, chart errors, and every rendering
fixture. The detailed inventory is in `docs/compatibility.md`.

Do not mark a migration fully compatible if it uses PHP charting. Keep core
period/set data and introduce an application renderer, or wait for a separately
versioned future `temporalchart` package.
