# Verification

## Verified

- `node --check tools/vibit`
- `node tools/vibit check architecture`
- `node tools/vibit check schemas`
- `node tools/vibit check memory`
- `node tools/vibit check contracts`
- `node tools/vibit check generated`
- `node tools/vibit check runtime`
- `node tools/vibit check all --json`
- `git diff --check`
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' --glob '!node_modules/**' .`

## Not Verified

- PostgreSQL runtime behavior is not verified because no Go runtime persistence implementation exists yet.
- Migrations are not verified because migration tooling has not been selected.
- S3-compatible object storage behavior is not verified because no object-storage adapter exists yet.
- MinIO deployment is not verified because MinIO has not been adopted as a dependency.

## Not Applicable

- Contract tests are not applicable because no public command, query, event, error, or permission contract changed.
