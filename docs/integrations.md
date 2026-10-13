# Calendar, clock, validation, config, and wire integration

`temporal` does not own a clock. Obtain current instants and timers from
`clock`, then pass concrete `time.Time` values into `instant`.

Civil arithmetic and DST resolution belong to `calendar`. Date-to-instant,
local-time application, daily-interval conversion, and civil snapping require
both `*time.Location` and a timezone resolution policy. This prevents
accidental dependence on the process-local timezone or the standard library's
implicit gap/fold choice.

`adapters/validation` returns immutable `validation` reports with stable
codes: `temporal_empty`, `time_of_day_range`, and `fixed_duration_range`.

`adapters/config` wrappers implement transactional text decoding recognized by
`config`: assignment occurs only after parsing succeeds, but shared mutable
wrappers require caller synchronization. Keep the wrapper in configuration structs
and call `Value()` after successful loading.

`adapters/wire.Document` and `CollectionDocument` are format-neutral and carry
JSON/YAML/TOML tags. They may be passed through the corresponding `wire`
encoder. Their JSON helpers apply the admission documented in [encoding](encoding.md);
other codecs do not automatically inherit those structural checks.

`adapters/postgres` maps finite Temporal periods and normalized sets to pgx
ranges and multiranges and supplies nullable `database/sql` wrappers. It owns
no connection or transaction. The four legacy integration imports remain
supported compatibility facades during the documented migration interval.
