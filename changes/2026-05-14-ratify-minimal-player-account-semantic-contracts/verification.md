# Verification

Verified:

- `node tools/vibit check contracts --json`
- `node tools/vibit check protocol --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change ratify-minimal-player-account-semantic-contracts --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Go tests; this contract ratification step does not change Go source.
- Live PostgreSQL verification; this step does not run PostgreSQL.
- Authentication, token, credential, player account persistence, session persistence, Protobuf envelope, and WebSocket handshake verification are not applicable because this change does not implement or change those surfaces.
