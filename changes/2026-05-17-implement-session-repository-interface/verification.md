# Verification

Verified:

- `go test ./internal/app/session`
- `node -c tools/vibit`
- `go test ./...`
- `node tools/vibit check migrations --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check work --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check change confirm-next-direction-after-session-repository-boundary --json`
- `node tools/vibit check change implement-session-repository-interface --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Live PostgreSQL verification, because this slice adds no PostgreSQL adapter or SQL execution.
- Protobuf generation, because no Protobuf source changed.

Notes:

- `node tools/vibit check runtime --json` passed with one existing warning: `runtime.identity_boundary` warns that `runtime/internal/platform/persistence/postgres/authentication_repository.go` mentions a credential dependency and should remain behind an explicit ratified boundary.
