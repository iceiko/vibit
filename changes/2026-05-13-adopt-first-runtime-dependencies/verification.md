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

- Runtime dependency behavior is not verified because no Go runtime implementation or `go.mod` exists yet.
- Live WebSocket behavior is not verified because no transport adapter exists yet.
- Protobuf generation is not verified because no `proto/`, `buf.yaml`, or generated protocol package exists yet.
- PostgreSQL persistence behavior is not verified because no repository implementation exists yet.
- Goose migration behavior is not verified because no migration files exist yet.

## Not Applicable

- Database migrations are not verified because no schema or migration files are added in this change.
- Public contract tests are not applicable because no command, query, event, error, or permission contract changed.
