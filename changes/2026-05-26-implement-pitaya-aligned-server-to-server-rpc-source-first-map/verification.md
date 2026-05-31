# Verification

## RED Checks

The initial checks failed because the W-0251 command, rule, and change artifacts were absent:

```text
node tools/vibit inspect pitaya-rpc --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_server_to_server_rpc_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_server_to_server_rpc_source_first_map

node tools/vibit check change implement-pitaya-aligned-server-to-server-rpc-source-first-map --json
# change directory does not exist
```

## Required Verification

- `node -c tools/vibit`
- `node tools/vibit inspect pitaya-rpc --json`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.pitaya_aligned_server_to_server_rpc_source_first_map`
- `node tools/vibit check change implement-pitaya-aligned-server-to-server-rpc-source-first-map --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Result

Verification passed after the W-0251 source-first map artifacts were updated. `node tools/vibit check runtime --json` reported the repository's existing `runtime.identity_boundary` warning for `runtime/internal/platform/persistence/postgres/authentication_repository.go` and no failures.
