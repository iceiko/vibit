# Plan

## Files To Create

- `runtime/go.mod`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `runtime/cmd/vibit-server/.gitkeep`
- `runtime/internal/app/.gitkeep`
- `runtime/internal/platform/transport/ws/.gitkeep`
- `runtime/internal/platform/protocol/protobuf/.gitkeep`
- `runtime/internal/platform/persistence/postgres/.gitkeep`
- `runtime/internal/platform/migrations/.gitkeep`
- `runtime/internal/platform/events/.gitkeep`
- `runtime/internal/platform/tx/.gitkeep`
- `runtime/internal/modules/inventory/.gitkeep`
- `runtime/internal/generated/contracts/inventory/.gitkeep`
- `runtime/internal/generated/proto/.gitkeep`
- `runtime/migrations/postgres/.gitkeep`
- `proto/README.md`
- `proto/README.zh-CN.md`
- `proto/vibit/inventory/v1/.gitkeep`
- `conversations/2026-05-13-go-runtime-skeleton.md`
- `changes/2026-05-13-bootstrap-go-runtime-skeleton/`

## Files To Edit

- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Generated Artifacts

None.

## Handwritten Logic

No server business logic.

Tooling logic changes are limited to runtime skeleton detection and future Go test discovery.

## Tests To Add Or Update

No Go tests are added.

Runtime CLI verification is updated and exercised through existing `node tools/vibit check ...` commands.

## Verification Plan

Run:

```bash
node --check tools/vibit
node tools/vibit check architecture
node tools/vibit check schemas
node tools/vibit check memory
node tools/vibit check contracts
node tools/vibit check generated
node tools/vibit check runtime
node tools/vibit check runtime --json
node tools/vibit check all --json
node tools/vibit check change bootstrap-go-runtime-skeleton
git diff --check
rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' --glob '!node_modules/**' .
```

## Rollback

If this skeleton proves too early, remove the new runtime/proto skeleton files and return `.arch/runtime.yaml` to `implementation_code_status: not_started`.
