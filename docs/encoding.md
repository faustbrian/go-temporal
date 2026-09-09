# Notation and encoding

Parsers accept exactly one complete valid UTF-8 value. They reject trailing
data, duplicate separators/components, unsupported precision, overflow, and
unknown formats.

- ISO 8601 `start/end` implies `ClosedOpen`; formatting another bound mode is
  rejected as lossy.
- ISO 80000 uses `[a,b]`, `[a,b)`, `(a,b]`, and `(a,b)`.
- Bourbaki uses inward/outward square brackets.
- Fixed ISO durations support elapsed weeks, days, hours, minutes, seconds, and
  nanoseconds. Years and months are rejected because they require a reference
  calendar date.

`adapters/wire.Document` uses `version`, `kind`, and canonical `value` fields for
scalar values. `CollectionDocument` uses `version`, `kind`, and canonical
`values` for normalized instant, date, and daily sets. Elements retain stable
normalized order and use ISO 80000 notation. The only current version is
`temporal/v1`; unknown versions, kinds, fields, trailing JSON, cardinality
excess, and byte-limit violations fail. New compatible readers may be added,
but v1 output will not silently change.

Hostile parse and decode operations bound the outer diagnostic by
`Limits.ErrorBytes`. Their messages use fixed package-owned stage text and do
not retain rejected input or unsafe parser causes. Use `errors.Is` and
`errors.As` for stable classification; error strings are not classification
contracts.

Core values do not expose mutable JSON unmarshal hooks. Decode through
`adapters/wire` or atomic `adapters/config` wrappers.
