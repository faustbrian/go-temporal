# API baseline

`v1.export` and `v2.export` preserve historical module export data. The v3
module uses `v3.export`, generated with the repository's pinned `apidiff`
policy after its dependency identities are finalized:

```sh
golib api update
```

`golib api check` rejects incompatible changes relative to the selected major's
baseline. The v3 baseline must be generated before release; earlier baselines
must not be rewritten as v3 evidence.
