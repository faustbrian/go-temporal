# Testing evidence

[Documentation index](README.md)

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
and pushes. `make check` and `make ci` select bounded local checks with the
checksum-pinned public verifier; `make check-full` preserves the explicit
broader all-tier contract, including configured service-backed and expensive
gates. Standalone `go test ./...` runs deterministic tests and fuzz corpus
regressions without sibling checkouts. The manifest requires PostgreSQL and
the `integration` test tag; service-aware Make checks may require Docker, unlike
untagged Go tests. Live PostgreSQL tests additionally require
`TEMPORAL_POSTGRES_DSN` pointing at disposable data. NilAway is advisory because the
upstream analyzer explicitly permits false positives; other gates in the active
contract block that run.

For active fuzz exploration select one target and a 2s budget; mutation is
optional local/scheduled assurance and may select one package. Neither is a
universal PR prerequisite. No new race campaign is added for immutable algebra.
Public API, codec, specification and PostgreSQL changes retain affected
compatibility, conformance, integration and direct-consumer evidence. See
[contributor verification](../CONTRIBUTING.md#verification) for copyable commands
and explicit PHP/release routes; [`.golib.yaml`](../.golib.yaml) retains the
fourteen-target full fuzz mapping.
The optional source-built verifier diagnostic and full release rehearsal retain
their original selection, refusal, evidence and cleanup contracts; they do not
substitute for ordinary hosted CI or actual public publication.

The requirement-by-requirement truth tables, algebra laws, resource budgets,
interoperability runs, and mutation classifications are recorded in the
[hardening report](hardening.md).
