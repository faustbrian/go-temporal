# Temporal specification conformance

This repository implements notation and database adaptation; it does not
implement the Temporal service protocol. The [decision register](../docs/specification-decisions.md)
separates normative source behavior from package-owned profiles and defensive
loss rejection. Source identities and integrity pins are in
[`sources.tsv`](sources.tsv); maintained-peer observations are in
[`interoperability.tsv`](interoperability.tsv).

| Decision | Source | Executable boundary | Disposition |
| --- | --- | --- | --- |
| TEMPORAL-DEC-001 | RFC 3339 | instant notation tests and fuzzing | resolved bounded Go profile |
| TEMPORAL-DEC-002 | pinned compatibility profile | interval and duration tests and fuzzing | resolved explicit subset |
| TEMPORAL-DEC-003 | pinned compatibility profile | bound notation tests and fuzzing | resolved plus named Bourbaki extension |
| TEMPORAL-DEC-004 | PostgreSQL 18.3 | pgx, SQL, fuzz, and live PostgreSQL integration | resolved finite lossless mapping |

Peer output is evidence, not authority. A disagreement remains classified and
cannot silently change a decision.
