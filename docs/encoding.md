# Notation and encoding

[Documentation index](README.md)

Parsers check UTF-8, complete notation structure, precision, arithmetic overflow
and selected formats. Instant endpoints enforce TEMPORAL-DEC-001: two-digit
hours, dot-only fractions within the configured precision (at most nine digits),
and Z or numeric offsets with hours below 24 and minutes below 60. The complete
Go RFC3339Nano parser validates the remaining date/time fields; local checks
reject its permissive fallback extensions. Parsed numeric offsets remain
observable. Bracketed forms separately check their separator structure.

- ISO 8601 `start/end` implies `ClosedOpen`; formatting another bound mode is
  rejected as lossy.
- ISO 80000 uses `[a,b]`, `[a,b)`, `(a,b]`, and `(a,b)`.
- Bourbaki uses inward/outward square brackets.
- Fixed ISO durations support elapsed weeks, days, hours, minutes, seconds, and
  nanoseconds. Years and months are rejected because they require a reference
  calendar date.

`adapters/wire.Document` uses `version`, `kind`, and canonical `value` fields for
scalar values. Scalar factories format values but do not independently reparse
them; `Marshal` and typed accessors perform separate admission. Formatters apply
`FormatBytes` after assembling notation and do not enforce parse `Precision` or
all endpoint/date-domain restrictions. Successful formatting is not a universal
round-trip promise. `CollectionDocument` uses `version`, `kind`, and canonical
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
`adapters/wire` or transactional `adapters/config` wrappers. Transactional means
failed parsing preserves the prior value, not lock-free synchronization.
JSON structure validation rejects repeated exact key strings; field binding still
uses `encoding/json` matching rules. Applications must synchronize mutable wrapper
and document fields. Collection accessors borrow `Values` during the call; returned
sets own their storage. Boundary diagnostics select safe causes explicitly; the
internal diagnostic helper owns its bounded display text, but is not a redactor
for arbitrary text or caller-owned cause objects.

Continue with [API reference](api.md).
