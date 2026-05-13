# Verification

Verified:

- `cd runtime && go test ./...`
  - Result: passed
- `cd runtime && go vet ./...`
  - Result: passed
- `node tools/vibit check runtime --json`
  - Result: passed
  - Summary: 83 passed, 0 warnings, 0 failures
- `node tools/vibit check change add-application-error-envelope-mapping --json`
  - Result: passed
  - Summary: 13 passed, 0 warnings, 0 failures
- `node tools/vibit check all --json`
  - Result: passed
  - Summary: 49 subchecks, 49 passed, 0 warnings, 0 failures
- `git diff --check`
  - Result: passed
- Secret scan for GitHub tokens in tracked and unignored files
  - Result: passed

Not verified:

- None.

Not applicable:

- Protobuf generation is not applicable because this change does not add or change `.proto` sources.
- Database migration verification is not applicable because no persistence adapter or migration is added.
- WebSocket transport verification is not applicable because this change does not add transport code.
