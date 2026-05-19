# Verification

Status: Verified

## Commands

```bash
cd runtime && go test ./cmd/vibit-server ./internal/app/connection ./internal/platform/transport/ws
cd runtime && go test ./...
node -c tools/vibit
node tools/vibit check change define-presence-lifecycle-functional-slice --json
node tools/vibit inspect next
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Results

- Passed `cd runtime && go test ./cmd/vibit-server ./internal/app/connection ./internal/platform/transport/ws`.
- Passed `cd runtime && go test ./...`.
- Passed `node -c tools/vibit`.
- Passed `node tools/vibit check change define-presence-lifecycle-functional-slice --json`.
- Passed `node tools/vibit inspect next`.
- Passed `node tools/vibit check work --json`.
- Passed `node tools/vibit check runtime --json`.
- Passed `node tools/vibit check memory --json`.
- Passed `node tools/vibit check schemas --json`.
- Passed `node tools/vibit check all --json`.
- Passed `git diff --check`.

## Notes

No live PostgreSQL check is required for this slice because the new startup wiring is covered by unit tests and the persistent data model is unchanged.

`node tools/vibit check runtime --json` still reports the pre-existing `runtime.identity_boundary` warning for `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
