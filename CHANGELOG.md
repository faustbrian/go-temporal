# Changelog

This project follows Semantic Versioning. Dates use ISO 8601.

## 3.0.0 - 2026-10-22

### Changed

- Move the module to `github.com/faustbrian/go-temporal/v3` and adopt
  `go-validation/v3`; `adapters/validation` returns Validation v3 validators.
  Tests use `go-config/v3`; `adapters/config` does not import go-config.

### Removed

- Remove the deprecated `postgres`, `temporalwire`, `temporalconfig` and
  `temporalvalidation` facades. Import the matching `adapters/*` package;
  `adapters/postgres` is package `temporalpostgres`. The removal precedes the
  announced support window; published v2 releases keep the facades and no
  further v2 releases are planned.

## Earlier unversioned notes (shipped by 2.0.3)

### Documentation

- Clarify parser versus formatter admission, zero values, callback and adapter
  ownership, applicable resource limits, and retained adapter semantics. Correct
  intersection complexity and document source-visible precision and degenerate
  relation caveats without changing the Temporal v2 API or runtime behavior.

## 2.0.3 - 2026-10-22

### Changed

- Publish the current source and documentation with a fresh immutable tag.
- Build with Go 1.27.3 to address GO-2026-6809 in `crypto/x509`.
  Consumers should use Go 1.27.3 or newer.

## 2.0.2 - 2026-10-14

### Fixed

- Classify coincident non-empty closed singleton endpoints as Equal in instant
  and civil-date relations without changing structural or represented-set equality.
- Restore TEMPORAL-DEC-001 endpoint admission: reject one-digit hours, comma
  fractions, offsets of 24 hours or 60 minutes, and excess fractional precision.
  Canonical and retained adapters share the parser: configuration retains prior
  values and wire returns zero on refusal. Valid numeric offsets and nanosecond
  fractions remain.
- Detach bounded diagnostic text from input backing storage while preserving safe
  classification and caller-owned cause objects. Declare the existing internal
  diagnostic package in the native inventories.

### Source accountability

Prepared from 9 non-merge commits after public v2.0.1 source
`9464371d0d72094c2096799de24e37389d1a5a99` through
`0c7a2d18c5b452c467173dc31f1a2b78f1aa31f9`. The restoration is delivered
separately after this frozen range; other Unreleased notes remain.

- `b0e0111a52e42e3bf6f7fb403842469000a66c61` — ci: use patched documentation tooling
- `fe964990960faf83f65147d77838ba59a813c945` — ci: adopt bounded published verification
- `01aa3dc435877e0c4f564c0df90ff5166c8d7b94` — docs: standardize repository development policy
- `78741c274d231addd6f323d115f649f419bf3869` — chore: cover non-Go dependency updates
- `62d7193e16e2f87785f5757e46674bfc73afcbe5` — chore(ci): pin shared tooling to v2.3.0
- `0e407212e5b10f9fcbeaf0ecd4b9b0db852b2302` — chore(ci): align manual tooling source pins
- `c7c36fd8865af6d1b7d52a776a07cfc926579c19` — docs: clarify temporal admission and ownership
- `762de8ebeddc393a320e705091c9914d5a7859b2` — chore: enforce source documentation contracts
- `0c7a2d18c5b452c467173dc31f1a2b78f1aa31f9` — docs: align adoption and documentation navigation

## 2.0.1 - 2026-10-07

### Changed

- Advance the test-only JSON codec composition dependency to published Wire
  v3.0.0, superseding the earlier Unreleased v1.0.1 pin. Canonical and retained
  adapters keep the same `temporal/v1` documents and Temporal v2 public API.

- Use Wire v1.0.1 for JSON codec composition in canonical and retained
  wire adapters, preserving the `temporal/v1` document representation.

## 2.0.0 - 2026-10-05

### Fixed

- Bound daily set Union input segments before copying or sorting, allowing
  twice the receiver's input-period limit to preserve circular expansion.
- Reject duplicate JSON fields and enforce configured wire parser depth before
  decoding scalar and collection documents.
- Admit configuration and SQL text bytes before conversion, reject nil text
  receivers, and enforce set and collection counts before owned copies.
- Count admitted JSON output bytes without allocating an encoded string before
  the format budget check.

### Changed

- Adopt public Validation v2.0.0 for canonical and retained validation
  adapters before the first Temporal v2 release. Returned validators compose
  with Validation `/v2` contexts, reports, and sentinels; nonempty and inclusive
  range rules retain their behavior. Historical Temporal v1 graphs are unchanged.
- Use the published Config v2.0.0 decoder for canonical and retained config
  adapter integration coverage. Temporal wrapper types and transactional
  decoding behavior remain unchanged.
- Move the root module and all owned imports to
  `github.com/faustbrian/go-temporal/v2`. Major-version types and sentinels have
  distinct identities; see the migration guide. Retained adapter facades and
  the `temporal/v1` wire schema remain supported within the v2 module.
- Adopt Calendar v2.0.0 for civil dates and timezone composition. Exposed
  Calendar types and sentinels require matching `/v2` imports; historical
  Temporal v1 dependency graphs are unchanged.
- Require Go 1.27.0 and use it for development, CI, and interoperability
  verification.
  TEMPORAL-DEC-001 sha256:aa3f62b689f2a90a208f0f50e57d4a00b096b67e2e1f4866f5cff785393427b2.

## 1.1.0 - 2026-09-09

### Added

- Add canonical `adapters/config`, `adapters/postgres`,
  `adapters/validation`, and `adapters/wire` package paths while retaining the
  released integration imports as compatibility facades.

### Fixed

- Bound hostile parse and decode diagnostics by `Limits.ErrorBytes`, remove
  rejected input and unsafe parser causes from returned errors, and preserve
  stable Temporal and safe structured classifications.

### Changed

- Replace copied repository tooling with the released `go-library-tools`
  v1.2.0 specification-governance contract while retaining repository-owned
  fixtures, API baselines, and content-addressed mutation evidence.
- Adopt the checksum-verified `go-library-tools` v1.3.0 CLI, schema-v2 cohesion
  metadata, and repository-local cohesion gate while retaining package-owned
  source and evidence.
- Adopt the checksum-verified `go-library-tools` v1.4.0 CLI and immutable
  shared workflow so authority monitoring uses the stabilized request profile
  and public-first module resolution.
- Pin the shared workflow and CLI to `go-library-tools` v1.6.1 so ordinary CI
  runs the proportional local contract while release rehearsals remain
  explicit.
- Reconcile the `go-calendar`, `go-config`, `go-validation`, and `go-wire`
  v1.0.0 checksums with their immutable public module archives.
- Require `go-validation` v1.1.0 so Temporal findings compose with context
  terminal reports without losing either error identity.

### Documentation

- Revalidate the pinned Temporal compatibility profile and advance the
  loss-checked range authority and supported deployment baseline to PostgreSQL
  18.6 without changing the selected mapping behavior.
  TEMPORAL-DEC-004 sha256:116d63c5d537d80713f36b10a1a992c6ad104b561abd53f5a3fae3e0ab7f5d1b.
- Remove the archived monorepo documentation link; package guidance remains in
  the repository-owned documentation.
- Link the module to the immutable v1.5.5 Golib ecosystem guidance.
- Publish the auditable [temporal specification register](docs/specification-decisions.md)
  and conformance map: TEMPORAL-DEC-001 sha256:7858a8bedd4143c177fb9670508d564d30049b8c099c66d19ffc3b83a7f3da14,
  TEMPORAL-DEC-002 sha256:01b85dcf0eac47ccfa1f69262044dd1206e261a72009a24832f112a2bb8555a1,
  TEMPORAL-DEC-003 sha256:2caaf9ce53ec7c2c383c73925f4a672a8f78e4724370bba981fbfdc8c82b2028,
  and TEMPORAL-DEC-004 sha256:aa94f39bba10c27ad85d3c2a400ce1190f0695bbba30d0ed526dcd5ec234755c.

## 1.0.0 - 2026-08-25

### Changed

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Link the package README to the repository-wide Golib documentation portal.

### Fixed

- Resolve API compatibility checks against the active monorepo workspace
  instead of downloading unpublished internal module versions.
- Parse duration component counts directly at the platform integer width before
  multiplication, eliminating narrowing conversions from hostile input.
- Run API compatibility through the isolated versioned-tool path so clean
  verification snapshots do not leak module flags into `apidiff`.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-temporal` identity while preserving its documented API and behavior.
- Refresh local `v0.0.0` owned-module checksums after dependency manifests and
  release notes were normalized; runtime behavior and public APIs are
  unchanged.
- Delegate local mutation checks to the canonical exact-100 repository runner
  and remove the superseded package-local Gremlins toolchain and configuration.
- Require owned sibling modules at local `v0.0.0`; clean external consumers
  pin each module to an exact main pseudo-version.

- Validate unsigned rounding modes with their meaningful upper bound only.
- Refresh owned-module checksums against the final consolidated archives.
- Use deterministic execution counts for default fuzz smoke campaigns while
  allowing explicit duration overrides for extended fuzzing.
- Normalized standalone module metadata against the canonical owned dependency
  graph, including complete checksums for clean consumer resolution.

### Added

- Exact boundary handling for parsing, limits, interval algebra, local-time
  arithmetic, splitting, rounding, and PostgreSQL range decoding.
- Explicit four-mode bounds and exhaustive Allen relations.
- Immutable instant and civil-date periods and normalized sets.
- Fixed durations, local times, circular daily intervals, and complements.
- Strict ISO 8601, ISO 80000, Bourbaki, JSON, SQL, and pgx adapters.
- `calendar`, `config`, `validation`, and format-neutral wire seams.
- Explicit civil snapping, local-time/daily application, and versioned set
  documents.
- Differential PHP fixtures, property/fuzz/race/mutation/benchmark gates.
- Exhaustive convenience-predicate tables, associative algebra properties, and
  a reproducible hardening evidence report.
- A generated, pinned inventory and behavior classification for every
  non-chart public PHP symbol.

### Compatibility

- PHP terminal and Gantt chart rendering is deferred. Full PHP-package
  compatibility is not claimed until an optional future renderer closes it.
