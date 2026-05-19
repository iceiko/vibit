# Verification

Verified:

- `node tools/vibit check change confirm-next-direction-after-session-handshake-ratification --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check all --json`

Not verified:

- None.

Not applicable:

- `go test ./...` is not required for this direction confirmation because no Go runtime behavior is added.
