# API baseline

`v1.export` preserves the historical v1 module export data. The v2 module uses
`v2.export`, generated with the repository's pinned `apidiff` policy after its
dependency identities are finalized:

```sh
golib api update
```

`golib api check` rejects incompatible changes relative to the selected major's
baseline. The v2 baseline must be generated before release; the v1 baseline
must not be rewritten as v2 evidence.
