# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected package's docs.
2. Run `go test ./...` or the affected package's narrower tests.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected stable decision entries, and complete the Specification Decisions
section of the pull request template. An unresolved interpretation or stale
source pin is release-blocking; peer behavior cannot silently select policy.
The current [specification decision register](docs/specification-decisions.md)
and conformance matrix are mandatory review inputs.

Mutation is optional local or scheduled assurance, not a required PR status.
Confirmed behavioral and security defects still require a focused regression
and repair; equivalent/invalid mutants retain reviewed narrow dispositions.

Do not add package-local workflows, permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Standalone `go test ./...` runs deterministic tests, fuzz seeds and saved corpus
regressions without sibling checkouts or ecosystem ledgers. Tagged PostgreSQL
tests are separate from this untagged route.

With the checksum-pinned public verifier selected by `.golib.yaml`, use:

```bash
make inventory
make check
make ci
```

`check` and `ci` use bounded local checks; ordinary CI consumes the public binary
without a verifier source build. The manifest requires PostgreSQL and the
`integration` tag, so service-aware Make checks may require Docker and differ
from untagged Go tests. No new ordinary race campaign is added for immutable
algebra. Specification, persistence, API and adapter changes retain relevant
conformance, compatibility, live integration and direct-consumer checks.

Short active exploration selects one target. Optional local mutation selects
one package; neither command is an all-change prerequisite:

```bash
GOMAXPROCS=2 go test -run '^$' -fuzz '^FuzzInstantNotation$' -fuzztime=2s -parallel=1 -timeout=2m ./notation
golib mutation --module . --package ./notation
```

Relevant PostgreSQL changes use a disposable database and exercise canonical
and retained adapters. Never use shared or production data:

```bash
TEMPORAL_POSTGRES_DSN='postgres://temporal:temporal@127.0.0.1:5432/temporal_test?sslmode=disable' \
go test -tags=integration -run '^TestPostgreSQLRangeAndMultirangeRoundTrips$' -count=1 -timeout=2m ./adapters/postgres ./postgres
```

`make check-full` retains the complete full-assurance contract, all fourteen
native fuzz target budgets, conformance, interoperability and benchmarks.
Preserve pinned PHP compatibility fixtures; relevant PHP compatibility changes
also use `make -f verification/package.mk php-compat` with the documented
`PHP_TEMPORAL_SOURCE`. Release qualification remains
`golib release dry-run --module .`, with the selected manual release workflow
retaining its PostgreSQL and native qualification. The separate verifier
diagnostic remains opt-in development evidence, not ordinary PR assurance.
Report unavailable, skipped and failed boundaries honestly; a local pass is not
full assurance, hosted-CI or public-release qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
