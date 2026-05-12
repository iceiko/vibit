# Plan

## Files To Change

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `decisions/ADR-0011-postgresql-and-object-storage-persistence.md`
- `conversations/2026-05-12-persistence-direction.md`

## Steps

1. Record the maintainer's persistence proposal in a change spec.
2. Add ADR-0011 for PostgreSQL, S3-compatible object storage, and MinIO candidate status.
3. Update runtime and convention manifests.
4. Update repository and module agent rules so future agents do not import persistence clients directly inside domain modules.
5. Update README files and paired Chinese translations.
6. Run repository verification and secret scan.

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

If this direction is rejected, supersede `ADR-0011` with a new ADR and update `.arch/runtime.yaml` and `.arch/conventions.yaml`. Do not silently delete the decision history.
