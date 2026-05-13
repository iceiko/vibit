# Verification

Verified:

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-minimal-websocket-process-wiring --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- Long-running manual server operation beyond the documented `go run ./cmd/vibit-server` path.

Not applicable:

- PostgreSQL persistence and migrations; deferred to `M-002`.
- Authentication/session validation; deferred until an auth or player module exists.
- Generated route registration; deferred until the generator standard exists.
