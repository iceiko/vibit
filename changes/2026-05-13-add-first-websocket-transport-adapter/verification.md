# Verification

Verified:

- `cd runtime && go mod tidy`
- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check runtime --json`
- `node tools/vibit check change add-first-websocket-transport-adapter --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`
- Secret scan for unignored files:

```bash
git ls-files --cached --others --exclude-standard -z | node -e '...'
```

Not verified:

- None.

Not applicable:

- Protobuf generation is not applicable because this change does not add or change `.proto` sources.
- Database migration verification is not applicable because no persistence adapter or migration is added.
- Domain module tests beyond existing runtime tests are not applicable because this change adds only transport code.

## Notes

An initial `node tools/vibit check all --json` run failed because this change spec still declared `verification.status: Not run`, which is not a valid final schema value. After updating the change spec to `Verified`, `check schemas` and `check all` passed.
