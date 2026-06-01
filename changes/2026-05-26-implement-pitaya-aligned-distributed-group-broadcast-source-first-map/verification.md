# Verification

## RED Checks

The initial checks failed because the W-0255 command, rule, and change artifacts were absent:

```text
node tools/vibit inspect pitaya-groups --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_distributed_group_broadcast_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_distributed_group_broadcast_source_first_map

node tools/vibit check change implement-pitaya-aligned-distributed-group-broadcast-source-first-map --json
# change directory does not exist
```

## Required Verification

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect pitaya-groups --json`
- `node tools/vibit inspect rule runtime.pitaya_aligned_distributed_group_broadcast_source_first_map`
- `node tools/vibit check change implement-pitaya-aligned-distributed-group-broadcast-source-first-map --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Result

Verified for this commit:

- `node -c tools/vibit`: passed.
- `node tools/vibit inspect next --json`: next-ready `W-0256 Define Pitaya-aligned cluster-safe session routing boundary gate`.
- `node tools/vibit inspect pitaya-groups --json`: emitted `kind: pitaya_groups_inspection` with next-ready `W-0256`.
- `node tools/vibit inspect rule runtime.pitaya_aligned_distributed_group_broadcast_source_first_map`: passed.
- `node tools/vibit check change implement-pitaya-aligned-distributed-group-broadcast-source-first-map --json`: passed 13, warnings 0, failures 0.
- `node tools/vibit check work --json`: passed 1548, warnings 0, failures 0.
- `node tools/vibit check runtime --json`: passed 21304, warnings 1, failures 0. The warning is the existing `runtime.identity_boundary` warning on `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
- `node tools/vibit check memory --json`: passed 4364, warnings 0, failures 0.
- `node tools/vibit check schemas --json`: passed 4765, warnings 0, failures 0.
- `git diff --check`: passed.
