# Changelog

This project follows Semantic Versioning. Dates use ISO 8601.

## Unreleased

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
- Reconcile the `go-calendar`, `go-config`, `go-validation`, and `go-wire`
  v1.0.0 checksums with their immutable public module archives.

### Documentation

- Revalidate the pinned Temporal compatibility profile and advance the
  loss-checked range authority and supported deployment baseline to PostgreSQL
  18.6 without changing the selected mapping behavior.
  TEMPORAL-DEC-004 sha256:116d63c5d537d80713f36b10a1a992c6ad104b561abd53f5a3fae3e0ab7f5d1b.
- Remove the archived monorepo documentation link; package guidance remains in
  the repository-owned documentation.
- Link the module to the immutable v1.4.0 Golib ecosystem guidance.
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
