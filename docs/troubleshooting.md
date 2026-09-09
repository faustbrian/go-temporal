# Troubleshooting

- `ErrBounds`: the bound enum is unknown or notation brackets do not match the
  selected format.
- `ErrEmpty`: the requested relation/gap requires a non-empty result.
- `ErrReversed`: the end precedes the start.
- `ErrStep`: a split, iterator, divisor, or rounding unit cannot progress.
- `ErrOverflow`: fixed arithmetic or PostgreSQL canonicalization exceeded its
  representable range.
- `ErrPrecision`: input has more fractional digits than permitted.
- `ErrLimit`: inspect `LimitError.Field`, `Value`, and `Max`.
- Truncated parse diagnostic: increase `Limits.ErrorBytes` only at a trusted
  boundary; use `errors.Is` or `errors.As` because the text is intentionally
  bounded and may be shorter than a sentinel's message.
- DST conversion failure: choose and pass an explicit `calendar` resolution
  policy; do not retry with an implicit timezone.
- PostgreSQL loss rejection: remove sub-microsecond precision explicitly at the
  application boundary if that loss is acceptable.
