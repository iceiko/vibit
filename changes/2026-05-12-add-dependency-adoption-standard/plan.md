# Plan

## Files To Change

- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `.arch/conventions.yaml`
- `.arch/dependencies.yaml`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `docs/dependency-adoption.md`
- `docs/dependency-adoption.zh-CN.md`
- `docs/_templates/dependency-adoption.md`
- `tools/vibit`
- `conversations/2026-05-12-dependency-adoption-standard.md`

## Steps

1. Add the dependency adoption standard and Simplified Chinese translation.
2. Add the reusable dependency adoption template.
3. Add `.arch/dependencies.yaml` with dependency slots and statuses.
4. Update architecture and repository docs to reference the registry.
5. Update `tools/vibit check architecture` to require the new artifacts.
6. Record conversation memory and verification.

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

If this standard is rejected, remove `.arch/dependencies.yaml`, the dependency adoption docs/template, and the architecture check requirements. Keep the change spec as history.
