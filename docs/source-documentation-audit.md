# Source documentation audit

## Pre-edit acceptance, owner and selected checks

Base `0e407212e5b10f9fcbeaf0ecd4b9b0db852b2302`, clean detached main-source
worktree. This worker owns the standalone Temporal source documentation pass;
the original unrelated toolchain branch is preserved. Parent owns enforcement,
remote delivery and aggregate evidence. Tier A truthful source prose only: preserve
all executable AST/tokens/literals/directives, adapter contracts, dependency graphs,
wire/SQL fixtures, pinned PHP compatibility/specification corpora and typed history.
No runtime repair, service/campaign qualification, normative rewrite or release.

Outcome: distinguish dense instant, discrete civil-date and circular daily set
semantics; zero/default/empty/singleton states; elapsed versus calendar arithmetic,
DST resolution and monotonic stripping; variable-size admission and bounded errors;
representation versus semantic equality; wire/config/validation/PostgreSQL adapter
ownership and failure/partial-result limits. Consumed dependencies are published Calendar
v2, Config v2, Validation v2, Wire v3 and pgx/v5 contracts. Direct surfaces include
configuration decoders, SQL/pgx callers, external wire codecs and validators, alongside retained
compatibility adapter facades and temporaltest helper interfaces.
Named risks: off-by-one boundary normalization, date/calendar overflow, midnight
versus 24:00 identity, DST gaps/folds versus elapsed steps, callback/predicate
cooperation, nil receivers/state-preservation on decode, database canonicalization,
source-data disclosure through retained causes, count/byte limits versus actual
allocation and false cross-adapter equivalence claims.
Smallest checks: full 122-source AST/significant-token/literal/directive proof,
useful parsed/HTML/native GoDoc for all actual packages, configured native apidiff
(no text snapshot yet discovered), actual docs snippet/example/generation recipe,
source/local links/offline specification/pinned changed-content+staged secrets and
ordinary complete final-diff truth review. No new runtime assertions/probes, active
fuzz, race, performance, PHP/native full conformance or PostgreSQL service campaign
for unchanged behavior. Every Go/native workflow awaits explicit parent grant with
three owned caches and immediate all-exit resource cleanup.

## Reproducible complete inventory

`git ls-files "*.go"` identifies 122 tracked Go files, 80 test files and 15
source directories in one root v2 module. modules.json/packages.json declare 14
packages; the actual imported internal/diagnostic directory is not catalogued and
is nevertheless included in this source audit/render plan. Parent owns tooling/
metadata disposition. Authenticated structural input reports 160 exported-comment
omissions, no missing package comment or marker/generator review. Existing terse
canonical package comments still need source-semantic review; counts are not prose
quality or proof of temporal/specification conformance.

- `.`: `bounded_diagnostic_test.go`, `bounds.go`, `bounds_test.go`, `errors.go`, `example_test.go`, `limits.go`, `limits_test.go`, `ordinary_boundary_test.go`, `relation.go`, `relation_test.go`, `security_boundary_test.go`, `successor_contract_test.go`.
- `adapters/config`: `admission_boundary_test.go`, `config.go`, `config_test.go`.
- `adapters/postgres`: `fuzz_test.go`, `integration_test.go`, `mutation_regression_test.go`, `range.go`, `range_test.go`, `sql.go`, `sql_boundary_test.go`, `sql_test.go`.
- `adapters/validation`: `validation.go`, `validation_test.go`.
- `adapters/wire`: `boundary_test.go`, `collection.go`, `collection_test.go`, `fuzz_test.go`, `marshal_boundary_test.go`, `size_internal_test.go`, `wire.go`, `wire_test.go`.
- `dateperiod`: `algebra.go`, `algebra_test.go`, `boundary_test.go`, `factories.go`, `fuzz_test.go`, `invariants_test.go`, `operations.go`, `operations_test.go`, `period.go`, `period_test.go`, `property_test.go`, `sequence.go`, `sequence_test.go`, `set.go`, `split.go`, `split_test.go`, `union_admission_test.go`.
- `instant`: `calendar.go`, `calendar_internal_test.go`, `calendar_test.go`, `concurrency_test.go`, `constructors.go`, `convenience_test.go`, `fuzz_test.go`, `mutation_regression_test.go`, `operations.go`, `operations_test.go`, `period.go`, `period_test.go`, `predicates.go`, `property_test.go`, `sequence.go`, `sequence_test.go`, `set.go`, `set_test.go`, `split.go`, `split_test.go`, `union_admission_test.go`.
- `internal/diagnostic`: `error.go`, `error_test.go`.
- `notation`: `date.go`, `date_test.go`, `duration.go`, `duration_test.go`, `fuzz_test.go`, `instant.go`, `instant_test.go`, `mutation_regression_test.go`, `timeofday.go`, `timeofday_test.go`.
- `postgres`: `fuzz_test.go`, `integration_test.go`, `postgres.go`, `range_test.go`, `sql_test.go`.
- `temporalconfig`: `admission_boundary_test.go`, `config.go`, `config_test.go`, `security_boundary_test.go`.
- `temporaltest`: `compatibility_fixture_test.go`, `temporaltest.go`, `temporaltest_test.go`.
- `temporalvalidation`: `validation.go`, `validation_test.go`.
- `temporalwire`: `admission_test.go`, `boundary_test.go`, `collection_test.go`, `fuzz_test.go`, `marshal_boundary_test.go`, `security_boundary_test.go`, `wire.go`, `wire_test.go`.
- `timeofday`: `adapters_test.go`, `calendar.go`, `calendar_internal_test.go`, `calendar_test.go`, `constructors.go`, `duration.go`, `duration_test.go`, `fuzz_test.go`, `interval.go`, `interval_calendar.go`, `interval_set.go`, `interval_test.go`, `mutation_regression_test.go`, `property_test.go`, `time.go`, `time_test.go`, `union_admission_test.go`.

## Variants, native ownership and investigation questions

The actual integration build tags live in canonical adapters/postgres and retained
postgres tests; untagged library docs/Examples do not certify a live database.
Retained postgres, temporalconfig, temporalvalidation and temporalwire packages
have canonical Deprecated paragraphs directing callers to adapters; retain v2
compatibility facades and their intentionally distinct symbols. No tracked Go
source generator or generated header is declared. Root .golib uses apidiff at
api/v2.export. verification/package.mk docs delegates scripts/check-docs.sh;
inspect its real snippets/selectors before choosing bounded checks. PHP generation
is a separate pinned external-source compatibility recipe, not source-doc output.

Missing exported fields/interfaces/aliases require semantic documentation.
Unanswered contracts require actual source/test/pinned dependency evidence:
constructor/zero state validity, normalization and collection copy boundaries,
custom limit propagation across operations, parser precision/bytes/depth effects,
callback or decoder ownership/cancellation, SQL NULL/scanner state and timestamp
precision, format validation versus exact represented membership, diagnostics
retention/privacy, overflow/partial-result semantics and actual DST delegate policy.
Record source-observed limitations separately from unanswered desired rationale.

## Semantic coverage and source-only limitations

Every tracked directory above was examined for declarations, API members, private
helpers, assurance selectors and resource/variant boundaries. The 80 test files
cover exact endpoint ties, singleton/empty states, copied slices, receiver limits,
normalization properties, overflow, synchronous callbacks, calendar/DST delegates,
parser/diagnostic ceilings, config/SQL assignment and wire structure. Fuzz files
have saved seed entrypoints for notation, splitting, daily algebra, SQL and JSON;
retained adapters intentionally duplicate canonical admission contracts. Six
benchmarks cover relation, splitting, rejection, normalization, notation and daily
algebra. Three root Examples compose instant/date/daily values; their bodies and
outputs remain unchanged. Tagged PostgreSQL tests own real service qualification.
PHP compatibility fixtures and the externally supplied PHP producer remain frozen,
not a generation recipe to execute in this source pass. No build/platform source
variant, generated Go source, go:generate route or unfinished source marker was
found. Existing lint directives remain executable tooling input and are preserved.

| Boundary | Clarified actual contract | Unresolved rationale / qualification |
|---|---|---|
| Shared enums/limits/errors | Bounds zero is half-open; invalid enum handling varies by helper; limits resolve fields independently and apply per operation; mutable LimitError is not a redactor. | No request-wide memory, CPU, deadline or diagnostic-graph budget is owned. |
| Instant/date algebra | Location-sensitive structural equality differs from represented equality; date sets canonicalize included dates; relation seam checks precede full endpoint equality. | Degenerate Allen laws are not established by the nondegenerate fixture table. Source-derived caveat only; desired semantic laws remain unchanged. |
| Sets, split and callbacks | Owned slices, receiver-limit propagation, linear intersection scan followed by sorting; callback effects are not rolled back; date/daily splitting uses Steps, not universal OutputPeriods. | No callback preemption, cancellation, panic recovery or aggregate budget. |
| Daily/calendar | Explicit midnight/24:00 and interval kinds; zero ordinary representation differs from Between construction; early duration branches precede bounds checks; Calendar v2 resolves explicit DST policy. | No fresh timezone/civil-range conformance qualification. DailyNonEmpty only checks CollapsedKind; intent for other degenerate states remains unknown. |
| Notation | Format limits checked after assembly; parse Precision is not a formatting contract; instant endpoints delegate to pinned Go RFC3339Nano parsing. | Local digit counting recognizes dot; Go 1.27.2 source admits comma fractions and truncates excess fractions in fallback parsing. No endpoint/fraction probe executed; strict DEC-001 remains desired profile, not a policy extension. |
| Diagnostics | Outer text and cause slice are bounded/copied selectively; cause objects and string backing can remain shared. Parser/wire callers select fixed text and safe causes explicitly. | Internal helper is not arbitrary-input redaction, erasure, deep copying or retained-memory qualification. |
| Config/SQL | Transactional assignment, nil rejection, borrowed text admission, zero SQL NULL, explicit precision/domain conversion constraints. | Not concurrent atomics or a running database/server-domain certification. |
| Wire/validation | Caller-owned envelopes; factories versus subsequent parsing admission; borrowed collection input; exact JSON key duplicate check plus standard field binding; synchronous validators/report context. | External codecs own their own structural admission. No blanket factory round-trip or normalization-free representation promise. |
| Retained adapters/helpers | Preserve distinct named wrappers, v2 successor delegation, temporal/v1 envelope and deprecation support window; Allen helpers require representable nondegenerate fixtures. | No successor removal or broader PHP/chart compatibility claim. |

Source findings above are not newly executed runtime results. Normative semantic
expectations, DEC-001, typed decision history, source authorities, corpora, immutable
API exports and dependencies remain unchanged. The actual imported diagnostic
package missing from both catalogs remains a parent-owned metadata question; the
source audit and planned native render include it.

## Selected local results and review boundary

Go 1.27.2 bounded native session completed: `gofmt` on changed Go paths;
stdlib parser/scanner comparison across all 122 files established exact
comment-free AST, raw significant tokens, literal spellings and lint/build
directives. All 580 comment groups parsed/rendered as GoDoc HTML; all 15 actual
packages rendered with native `go doc -all`, including the uncatalogued internal
package. No optional-semicolon normalization or executable exception was needed.
Assurance comments describe actual fuzz assertions and benchmark setup; equal
step/output settings in existing date/daily fuzz targets are not independent
proof of those two limits. Their input, assertions, seeds and measured bodies
are unchanged.

`go test . -run '^Example' -count=1 -timeout=2m -json` passed the three existing
root Examples with ordinary default vet enabled. `make -f verification/package.mk
docs` executed the real required-file/phrase and all-package GoDoc recipe. The
checksum-authenticated declared Tools v2.3.0 passed `api check --module .`,
`docs check --local`, and offline `specification check`. API mode is apidiff;
`api/v2.export` is unchanged and there is no full-prose snapshot producer to
regenerate. Qualified Tools v2.4.1 separately checks source documentation
structure; this does not replace v2.3.0 admission for native repository gates.

No behavioral tests, active fuzzing, benchmarks, race, services, PHP generator,
endpoint/fraction probe or release qualification ran. Existing assertions were
read as source support, not reported as newly passed runtime proof. Go cache,
module cache, temporary Go directory, helper, baseline copies, render output and
logs were task-owned and removed immediately after the bounded session, including
readonly downloaded modules. Original checkout remains untouched.

The ordinary Tier A whole-diff review checks truth, usefulness, idiomatic
attachment, duplication and compatibility: comments/current factual prose only;
strict semantic expectations, decision records/history, API export, compatibility
fixtures/corpora, producer inputs, dependency manifests, CI and catalog metadata
remain exact. Cheap affected local links, source structure, changed-content and
staged-content secret checks and explicit Intent commit hygiene complete the
local handoff; hosted CI and remote-main delivery are coordinator-owned.

## Objective Enforcement

The checksum-pinned public Tools v2.4.0 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.
