# Verification

Verified:

- `gofmt -w runtime/internal/app/authentication/service.go runtime/internal/app/authentication/service_test.go`
- `cd runtime && go test ./internal/app/authentication`
- `cd runtime && go test ./...`
- `node -c tools/vibit`
- `node tools/vibit check schemas --json`
  - passed: 2032
  - warnings: 0
  - failures: 0
- `node tools/vibit check memory --json`
  - passed: 1492
  - warnings: 0
  - failures: 0
- `node tools/vibit check contracts --json`
  - passed: 291
  - warnings: 0
  - failures: 0
- `node tools/vibit check generated --json`
  - passed: 205
  - warnings: 0
  - failures: 0
- `node tools/vibit check runtime --json`
  - passed: 2657
  - warnings: 1
  - failures: 0
  - known warning: `runtime.identity_boundary` for `runtime/internal/platform/persistence/postgres/authentication_repository.go`
- `node tools/vibit check module authentication --json`
  - passed: 23
  - warnings: 0
  - failures: 0
- `node tools/vibit check work --json`
  - passed: 659
  - warnings: 0
  - failures: 0
- `node tools/vibit check change implement-authentication-service-behavior-skeleton --json`
  - passed: 13
  - warnings: 0
  - failures: 0
- `node tools/vibit check all --json`
  - subchecks: 160
  - passed: 160
  - warnings: 1
  - failures: 0
  - known warning: `runtime.identity_boundary` for `runtime/internal/platform/persistence/postgres/authentication_repository.go`
- `git diff --check`
- Secret scan for GitHub tokens excluding `.git` and `.vibit.local.env`.

Not verified:

- None.

Not applicable:

- Live PostgreSQL verification is not required because this skeleton does not change persistence behavior, SQL migrations, or repository implementations.
- WebSocket/protocol integration verification is not required because this skeleton does not expose authentication protocol carriers or wire runtime authentication behavior.
