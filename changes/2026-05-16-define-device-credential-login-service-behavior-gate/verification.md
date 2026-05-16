# Verification

Verified:

- `node -c tools/vibit`
- `cd runtime && go test ./...`
- `node tools/vibit check schemas --json`
  - passed: 2054
  - warnings: 0
  - failures: 0
- `node tools/vibit check memory --json`
  - passed: 1515
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
  - passed: 2833
  - warnings: 1
  - failures: 0
  - known warning: `runtime.identity_boundary` for `runtime/internal/platform/persistence/postgres/authentication_repository.go`
- `node tools/vibit check module authentication --json`
  - passed: 23
  - warnings: 0
  - failures: 0
- `node tools/vibit check work --json`
  - passed: 665
  - warnings: 0
  - failures: 0
- `node tools/vibit check change define-device-credential-login-service-behavior-gate --json`
  - passed: 13
  - warnings: 0
  - failures: 0
- `node tools/vibit check all --json`
  - subchecks: 161
  - passed: 161
  - warnings: 1
  - failures: 0
  - known warning: `runtime.identity_boundary` for `runtime/internal/platform/persistence/postgres/authentication_repository.go`
- `git diff --check`
- Secret scan for GitHub tokens excluding `.git` and `.vibit.local.env`

Not verified:

- None.

Not applicable:

- Live PostgreSQL verification is not required because this gate does not change persistence behavior or SQL migrations.
- WebSocket/protocol integration verification is not required because this gate does not expose authentication protocol carriers or wire runtime authentication behavior.
- Authentication service behavior is not verified because this gate intentionally does not add login execution code.
