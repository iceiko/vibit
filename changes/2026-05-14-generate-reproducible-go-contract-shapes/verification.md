# Verification

Verified:

- `node tools/vibit generate contract-shapes all`
- `node tools/vibit inspect generated --json`
- `node tools/vibit check generated --json`
- `cd runtime && go test ./...`

Not verified:

- None.

Not applicable:

- Live PostgreSQL verification is not required because no persistence behavior changed.
- WebSocket, authentication, token, credential, and session validation behavior are not changed by generated contract shape metadata.
