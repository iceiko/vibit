# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit check schemas --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check protocol --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change confirm-next-direction-after-authentication-login-route --json`
- `node tools/vibit check change ratify-session-persistence-websocket-handshake-authentication --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- `go test ./...` is not required for this standards-only gate because no Go runtime behavior is added or changed.
