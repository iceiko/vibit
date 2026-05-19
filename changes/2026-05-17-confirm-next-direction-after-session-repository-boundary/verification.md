# Verification

Verified:

- `go test ./internal/app/session`

Planned final verification:

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

- None yet.

Not applicable:

- Live PostgreSQL verification, because this direction change adds no adapter or SQL behavior.
