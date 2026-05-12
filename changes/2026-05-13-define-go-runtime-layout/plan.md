# Plan

## Files To Change

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `decisions/ADR-0014-go-runtime-layout-and-boundaries.md`
- `conversations/2026-05-13-go-runtime-layout-and-boundaries.md`
- `changes/2026-05-13-define-go-runtime-layout/`

## Steps

1. Record the continuation request and clarified scope.
2. Add ADR-0014 for Go runtime package layout and boundary rules.
3. Update `.arch/runtime.yaml` with machine-readable layout, protocol, migration, and transaction boundary sections.
4. Update `.arch/conventions.yaml` to link ADR-0014 and summarize layout.
5. Update AGENTS guides and README translations.
6. Run verification and secret scan.

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
node tools/vibit check all --json
git diff --check
rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' --glob '!node_modules/**' .
```

## Rollback

If the package layout proves unsuitable, supersede ADR-0014 with a new runtime layout ADR and update `.arch/runtime.yaml` before creating or moving Go implementation code.
