# Verification

Verified:

- `gofmt -w runtime/internal/app/authentication/material_generation.go runtime/internal/app/authentication/material_generation_test.go`
- `cd runtime && go test ./internal/app/authentication`
- `cd runtime && go test ./...`
- `node -c tools/vibit`
- `node tools/vibit check schemas --json`
  - passed: 1943
  - warnings: 0
  - failures: 0
- `node tools/vibit check memory --json`
  - passed: 1393
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
  - passed: 2117
  - warnings: 1
  - failures: 0
  - known warning: `runtime.identity_boundary` for `runtime/internal/platform/persistence/postgres/authentication_repository.go`
- `node tools/vibit check module authentication --json`
  - passed: 23
  - warnings: 0
  - failures: 0
- `node tools/vibit check work --json`
  - passed: 623
  - warnings: 0
  - failures: 0
- `node tools/vibit check change implement-token-credential-material-generation-helpers --json`
  - passed: 13
  - warnings: 0
  - failures: 0
- `node tools/vibit check all --json`
  - subchecks: 153
  - passed: 153
  - warnings: 1
  - failures: 0
  - known warning: `runtime.identity_boundary` for `runtime/internal/platform/persistence/postgres/authentication_repository.go`
- `git diff --check`
- Secret scan for GitHub tokens excluding `.git` and `.vibit.local.env`

Not verified:

- None.

Not applicable:

- Live PostgreSQL verification is not required because this change does not touch persistence behavior or SQL migrations.
- WebSocket/protocol integration verification is not required because this change does not expose authentication protocol carriers or wire runtime authentication behavior.
