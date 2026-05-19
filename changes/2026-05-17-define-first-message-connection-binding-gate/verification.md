# Verification

Verified:

- `node -c tools/vibit`
- `node tools/vibit check schemas --json`
- `node tools/vibit check protocol --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check module authentication --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change confirm-next-direction-after-session-handshake-ratification --json`
- `node tools/vibit check change define-first-message-connection-binding-gate --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- `go test ./...` is not required for this standards-only gate because no Go runtime behavior is added or changed.
