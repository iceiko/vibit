# Verification

Verified:

- `buf generate`
- `node -c tools/vibit`
- `node tools/vibit check change implement-protocol-logout-route --json`
- `node tools/vibit check runtime --json`
- `cd runtime && go test ./...`

Interim results:

- Change check passed.
- Runtime check passed with the existing `runtime.identity_boundary` warning for `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
- Runtime Go tests passed.

Broader continuation verification:

- `node tools/vibit check schemas --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check protocol --json`
- `node tools/vibit check work --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not applicable:

- Live PostgreSQL verification is not required for the default repository check path.
- No migration, repository interface, WebSocket transport close, runtime session revocation, reconnect, protocol session carrier, dependency, SDK, cluster, or direct Nakama/Pitaya API compatibility change was made.
