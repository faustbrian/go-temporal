# PostgreSQL

`adapters/postgres` maps instant periods to `tstzrange`, date periods to
`daterange`, and normalized sets to corresponding multiranges. The released
`postgres` path remains a compatibility facade.

PostgreSQL 18.6 is the supported deployment baseline. Its range definition is
identical to the reviewed 18.3 definition, including timestamp precision and
discrete-range canonicalization.

PostgreSQL timestamps are microsecond precision. Instant values with non-zero
sub-microsecond data are rejected; they are never rounded silently. Unbounded
and PostgreSQL-empty ranges are rejected because core periods are bounded and
carry their own empty representation. SQL NULL uses `InstantRange` or
`DateRange` wrappers and is distinct from an empty range.

PostgreSQL canonicalizes discrete dateranges to closed-open. Conversion back to
`dateperiod` retains represented dates, so structural bounds may differ while
`SetEqual` remains true. Maximum-date exclusive-end overflow is rejected.

The release dispatch (`release_dry_run: true`) additionally requires the two
canonical and retained adapter tests against digest-pinned PostgreSQL 18.6.
Missing or skipped expected tests fail that release gate. Connections and
queries share a 30-second test context; cleanup uses an independent 10-second
context. Ordinary local runs without a DSN may still skip these tagged tests.

Run the disposable integration suite locally with:

```sh
TEMPORAL_POSTGRES_DSN='postgres://temporal:temporal@localhost/temporal_test' \
  go test -tags=integration ./adapters/postgres ./postgres
```
