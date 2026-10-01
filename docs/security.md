# Security model

## Threat model revision 2 (2026-10-01)

This model covers the Temporal `/v2` module and its retained compatibility
packages; it does not certify the earlier v1.1.0 release. Public tags and
releases establish version availability. Calendar major adoption requires
explicit nominal-type and consumer reconciliation. This model must be reviewed
when a new input format, persistence adapter,
callback surface, dependency, or resource-producing operation is added.

### Assets and security properties

- Process availability: hostile inputs must not trigger unbounded allocation,
  iteration, recursion, or diagnostics.
- Temporal integrity: parsing and persistence conversion must not silently
  change bounds, precision, dates, instants, or DST resolution.
- Caller state: failed decoding must not partially update receiver values or
  expose mutable package-owned slices.
- Diagnostic confidentiality: rejected payloads and dependency parser details
  must not be retained in returned errors.

The module does not own credentials, authorization policy, network listeners,
files, environment variables, subprocesses, or database connections.

### Trust boundaries and attacker capabilities

The following values are untrusted: notation and configuration text, JSON wire
payloads, `database/sql` range values, pgx range and multirange values, caller
supplied collections and limits, split and iteration parameters, time zones,
and values returned by caller callbacks. An attacker may choose malformed,
ambiguous, deeply nested, oversized, high-cardinality, Unicode-confusable,
precision-losing, reversed, empty, or arithmetic-boundary inputs and may repeat
requests concurrently.

### Controls

- `temporal.Limits` validates hard ceilings for parse, format, and diagnostic
  bytes; fractional precision; parser depth; input and output cardinality; and
  iterative steps. Applications should lower these ceilings at each request
  boundary.
- Text and SQL byte inputs are length-checked before conversion or parsing.
  JSON is valid UTF-8, fully consumed, depth-limited, duplicate-key rejecting,
  and decoded with unknown fields disabled.
- Notation accepts ASCII grammar tokens only, requires unique ordered
  components, and rejects unsupported precision and checked-arithmetic
  overflow.
- Set, split, and iteration work has explicit cardinality or step ceilings and
  returns no partial result when a ceiling is exceeded. Returned slices are
  copies.
- Daily set Union admits combined normalized segments before copying or
  sorting, with an inclusive ceiling of twice the receiver's resolved
  `InputPeriods` (at most 200,000 segments). Circular raw intervals may expand
  into two segments and normalization does not preserve raw input counts;
  constructors still enforce raw `InputPeriods`, and normalized output retains
  its separate `OutputPeriods` ceiling.
- Text and SQL scanners parse into temporary values before assignment and
  reject nil receivers instead of panicking.
- PostgreSQL conversions reject infinite, invalid, empty-as-value, and
  precision-losing values unless nullability is represented explicitly.
- Local-to-instant conversion requires an explicit location and DST resolution
  policy; the process-local time zone is never an implicit input.
- Boundary diagnostics contain fixed safe text and sentinel or structured
  limit causes. They do not retain rejected text or third-party parser errors.

Hostile-boundary callers must classify failures with `errors.Is` or
`errors.As`, never by matching diagnostic text.

## Accepted residual risks

| Risk | Owner | Rationale | Mitigation | Review condition |
| --- | --- | --- | --- | --- |
| A caller callback passed to `Transform` or `Reduce` can block, panic, or perform side effects. | Integrating application | The library cannot safely cancel or recover arbitrary application code without changing callback semantics. | Require trusted callbacks that return promptly. Run untrusted callback code in process or subprocess isolation, or behind an explicitly leak-bounded application wrapper; a deadline helps only when the callback cooperates with it. Invocation count remains bounded by set cardinality. | Revisit if callbacks accept contexts, execute concurrently, or cross a plugin or RPC boundary. |
| Limits apply per operation, not across a complete request containing multiple operations. | Integrating application | The package has no request lifecycle or shared budget ownership. | Set tighter per-boundary limits and enforce aggregate request quotas, concurrency limits, and deadlines in the application. | Revisit if the module gains request-scoped execution or batch orchestration. |
| Time-zone behavior depends on the Go runtime's supplied IANA data and explicit caller location. | Integrating application | Civil-time rules are external, versioned data and may change legislatively. | Pin and update the runtime or bundled tzdata deliberately; test business-critical transitions with the chosen resolution policy. | Revisit on Go or tzdata upgrades, or when supporting a new jurisdiction-sensitive workflow. |
| Algorithms and errors are not constant-time. | Integrating application | Inputs are temporal values, not cryptographic secrets, and work remains bounded. | Do not encode secrets in temporal payloads; isolate callers that introduce secret-dependent behavior. | Revisit if temporal operations enter an authentication or cryptographic decision path. |
| Database authentication, authorization, query timeouts, and transport security are outside the PostgreSQL value adapters. | Integrating application | The adapters do not open connections or issue queries. | Configure least privilege, TLS, statement timeouts, and cancellation in the owning database client. | Revisit if an adapter begins owning connections or executing SQL. |
