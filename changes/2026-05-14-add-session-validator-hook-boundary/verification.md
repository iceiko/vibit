# Verification

Verified:

- `cd runtime && go test ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-session-validator-hook-boundary --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Live PostgreSQL verification is not required for this application hook change.
- Authentication, token, credential, session store, player account migration, Protobuf envelope, and WebSocket handshake verification are not applicable because this change does not implement or change those surfaces.
