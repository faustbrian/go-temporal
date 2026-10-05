# Security policy

## Reporting

Report vulnerabilities privately through GitHub Security Advisories. Do not
include secrets, production data, or exploit payloads in a public issue.

Include the affected module and version, impact, reproduction, preconditions,
and any suggested mitigation. Maintainers follow the pinned ecosystem
[vulnerability-management process](https://github.com/faustbrian/go-library-tools/blob/77bfd78c12a853f0d490bb27a3fbcb5f34330772/docs/ecosystem/security/vulnerability-management.md)
for severity assignment, acknowledgement and remediation targets, private
embargo handling, advisory ranges, and coordinated affected-module releases.
Targets begin when sufficient private evidence exists to reproduce or bound
the report; the owner explains any changed target to the reporter.

## Supported versions

The latest minor of each supported released major receives security fixes.
The v1.1 line remains supported; the v2 security model applies to the `/v2`
module. Published tags and GitHub releases establish version availability.

## Threat model

The package treats notation, JSON, SQL range text, sequence inputs, and split
parameters as untrusted. `temporal.Limits` bounds parse bytes, precision, error
and output bytes, parser depth, input/output period counts, and steps. Parsers
require valid UTF-8, ASCII grammar tokens, complete consumption, duplicate-free
JSON objects, unique ordered notation components, and checked arithmetic.

Set operations reject cardinality expansion before returning partial output.
Iterators reject zero and negative steps and prove progress. Returned slices
are copies. PostgreSQL adapters reject unbounded, empty, NULL-as-value, and
microsecond-loss cases unless represented explicitly by nullable wrappers.

DST resolution is outside core algebra. Callers must supply a location and a
`calendar/timezone.Resolution`; no implicit local timezone is consulted.

Known residual risks are documented in [docs/security.md](docs/security.md).
