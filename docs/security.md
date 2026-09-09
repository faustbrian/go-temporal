# Hostile input and resource safety

Always lower defaults at a trust boundary when the application knows tighter
limits. A parser for a single header may reasonably use `ParseBytes: 256`; a
bulk API should cap both `InputPeriods` and `OutputPeriods` per request.

Precision is zero through nine decimal digits. Split and step operations reject
non-positive steps and stop at `Limits.Steps`. Arithmetic checks overflow before
allocation or returned mutation. Formatting checks output bytes.

Unicode punctuation that resembles ASCII brackets, commas, slashes, signs, or
digits is rejected. There is no natural-language or locale-dependent parser.

At hostile text boundaries, `ErrorBytes` is the maximum byte length of the
outer error string after valid limits resolve. Fixed safe text is truncated at
a UTF-8 boundary without an ellipsis. Returned error graphs preserve Temporal
sentinels and safe structured `LimitError` values, but deliberately do not
retain rejected text, JSON decoder errors, `time.ParseError`, numeric parser errors,
or PostgreSQL parser details. Classify failures with `errors.Is` or
`errors.As`, never string matching.

Do not resolve a local time onto a date without an explicit location and DST
policy. A local time alone is not an instant. Do not treat a calendar day as
`24*time.Hour`; DST transitions disprove that conversion.
