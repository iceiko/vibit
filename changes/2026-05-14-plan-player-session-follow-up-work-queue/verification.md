# Verification

Verified:

- `node tools/vibit check work --json`
- `node tools/vibit inspect work --json`
- `node tools/vibit check change plan-player-session-follow-up-work-queue --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Go tests; this planning step does not change Go source.
- Live PostgreSQL verification; this planning step does not run PostgreSQL.
- Authentication, token, credential, player account persistence, Protobuf envelope, and WebSocket handshake verification are not applicable because this change does not implement or change those surfaces.
