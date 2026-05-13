# Verification

Verified:

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-websocket-protobuf-application-composition --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- `/v1/ws` process wiring; deferred to `W-0008`.
- PostgreSQL persistence and migrations; deferred to `M-002`.
- Generated Protobuf updates; this change does not edit protocol sources or generated files.
