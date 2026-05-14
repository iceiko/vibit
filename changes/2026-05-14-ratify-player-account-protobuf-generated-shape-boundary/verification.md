# Verification

Verified:

- `buf lint`
- `buf generate`
- `node tools/vibit check protocol --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change ratify-player-account-protobuf-generated-shape-boundary --json`
- `node tools/vibit check all --json`
- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Live PostgreSQL verification; this change does not add database behavior.
- Authentication, token, credential, player account persistence, session persistence, Protobuf envelope, WebSocket handshake, and runtime player handler verification are not applicable because this change does not implement or change those surfaces.
