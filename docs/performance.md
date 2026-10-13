# Performance

[Documentation index](README.md)

Values are small immutable structs. Normalized sets own copied slices.

| Operation | Complexity |
|---|---|
| period relation/membership | `O(1)` |
| set construction | `O(n log n)` |
| normalized set search | `O(log n)` |
| normalized union | `O((n+m) log(n+m))` |
| normalized intersection | `O(n+m)` scan plus `O(k log k)` result normalization |
| split/steps | `O(k)`, bounded by `Limits.Steps` |

`k` is the emitted result count. The linear scan is followed by the common
stable-sorting constructor; this table describes actual code paths, not a
performance measurement. Limits are operation-specific ceilings, not aggregate
request, memory, or deadline budgets. Formatters check their byte cap after
notation assembly; `SplitDays`, daily `Split`, and `Steps` use the step ceiling,
not a universal `OutputPeriods` check.

The shared `golib check` contract runs relation, normalization, parser, and
daily-set allocation benchmarks. Benchmark output is an artifact, not a
cross-machine hard threshold; regression decisions compare like-for-like
runners.

Continue with [API reference](api.md).
