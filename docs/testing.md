# Testing evidence

Production statement coverage is exactly 100%, enforced by the shared `golib`
repository contract.
Confidence additionally comes from:

- all 13 Allen relations across all four bound modes;
- randomized union/intersection/difference conservation;
- exhaustive hourly circular complements across all bounds;
- fixed regression tests for surviving singleton subtraction boundaries;
- strict notation round trips and generated PHP fixtures for bounds, local
  values, duration arithmetic, time arithmetic, interval kinds, predicates,
  algebra, complement, splitting, and stepping;
- fuzz targets for instant/date/daily/duration/time notation, split progress,
  set normalization, both wire paths, and both PostgreSQL range paths;
- race tests with concurrent reads over shared immutable values;
- PostgreSQL 18 range and multirange integration;
- Gremlins arithmetic and conditional mutation operators;
- allocation-reporting benchmarks for relations, parsing, 1,000-period
  normalization, splitting, daily algebra, and early limit rejection.

The reusable CI workflow runs the proportional local contract on pull requests
and pushes. `make check` runs the broader all-tier contract, including configured
service-backed and expensive gates. PostgreSQL integration tests participate
when `TEMPORAL_POSTGRES_DSN` is configured. NilAway is advisory because the
upstream analyzer explicitly permits false positives; other gates in the active
contract block that run.

The requirement-by-requirement truth tables, algebra laws, resource budgets,
interoperability runs, and mutation classifications are recorded in the
[hardening report](hardening.md).
