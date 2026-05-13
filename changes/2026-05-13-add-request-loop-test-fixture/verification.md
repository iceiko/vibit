# Verification

Verified:

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change add-request-loop-test-fixture --json`
- `node tools/vibit check all --json`
- `git diff --check`

Not verified:

- None.

Not applicable:

- Runtime behavior changes; this change is test-only.
- PostgreSQL persistence and migrations.
- Authentication/session validation.
- Generated route registration.
