# Verification

Verified:

- `cd runtime && go test ./...`
- `node tools/vibit check architecture --json`
- `node tools/vibit check protocol --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change define-application-request-identity-handoff-boundary --json`

Not verified:

- `node tools/vibit check all --json` pending after work queue update.
- `git diff --check` pending after work queue update.

Not applicable:

- Live PostgreSQL verification is not required for this application handoff type change.
- Authentication, token, credential, session store, player account migration, Protobuf envelope, and WebSocket handshake verification are not applicable because this change does not implement or change those surfaces.
