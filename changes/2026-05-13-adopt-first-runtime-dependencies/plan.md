# Plan

## Files To Change

- `.arch/dependencies.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `decisions/ADR-0012-agent-technical-decision-authority.md`
- `decisions/ADR-0013-first-go-runtime-dependencies.md`
- `conversations/2026-05-13-technical-decision-authority-and-runtime-dependencies.md`
- `changes/2026-05-13-adopt-first-runtime-dependencies/`

## Steps

1. Record the maintainer authorization in a conversation log.
2. Add ADR-0012 for the technical decision authority boundary.
3. Add ADR-0013 for the first accepted Go runtime dependency choices.
4. Update `.arch/dependencies.yaml` for accepted and deferred slots.
5. Update `.arch/runtime.yaml` and `.arch/conventions.yaml` so agents can inspect the accepted choices.
6. Update AGENTS guides and architecture README translations.
7. Run repository verification and secret scan.

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

If one accepted dependency becomes unsuitable, supersede ADR-0013 or create a focused replacement ADR. Update `.arch/dependencies.yaml` before changing implementation imports.
