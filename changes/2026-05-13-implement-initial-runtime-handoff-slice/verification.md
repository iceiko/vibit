# Verification

Verified:

- `buf generate`
- `go mod tidy`
- `node --check tools/vibit`
- `go test ./...` from `runtime/`
- `node tools/vibit check runtime --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check protocol --json`
- `node tools/vibit check change implement-initial-runtime-handoff-slice --json`
- `node tools/vibit check all --json`
- `buf lint`
- `git diff --check`
- Secret scan for GitHub tokens outside ignored local files:
  `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' --glob '!node_modules/**' .`

Not verified:

- Full WebSocket transport behavior is not verified because transport implementation is not part of this slice.
- PostgreSQL repository behavior is not verified because persistence implementation is not part of this slice.
- Migration behavior is not verified because no migrations are added.

Not applicable:

- Public contract compatibility testing is not applicable because no public command, query, event, error, or permission contract changes are made.
