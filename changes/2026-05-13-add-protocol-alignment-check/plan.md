# Plan

## Files To Edit

- `tools/vibit`
- `rules/check-rules.json`
- `.arch/contracts.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `proto/README.md`
- `proto/README.zh-CN.md`

## Files To Create

- `changes/2026-05-13-add-protocol-alignment-check/`
- `conversations/2026-05-13-protocol-alignment-check.md`

## Generated Artifacts

None.

## Handwritten Logic

Add CLI logic to:

- Discover registered command, query, and event contracts.
- Derive expected Protobuf source paths.
- Derive expected message names.
- Parse shallow field names from contract input/output/payload sections.
- Validate `.proto` files when they exist.
- Report planned missing `.proto` files as passing while protocol generation has not started.

## Tests To Add Or Update

No separate test framework exists for `tools/vibit` yet.

The change is verified through existing CLI check commands.

## Verification Plan

Run:

```bash
node --check tools/vibit
node tools/vibit check protocol
node tools/vibit check protocol --json
node tools/vibit check architecture
node tools/vibit check schemas
node tools/vibit check memory
node tools/vibit check contracts
node tools/vibit check generated
node tools/vibit check runtime
node tools/vibit check all --json
node tools/vibit check change add-protocol-alignment-check
git diff --check
rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' --glob '!node_modules/**' .
```

## Rollback

Remove `check protocol` and return `.arch/contracts.yaml` alignment status to `planned` if the first alignment rule proves misleading before the first `.proto` files exist.
