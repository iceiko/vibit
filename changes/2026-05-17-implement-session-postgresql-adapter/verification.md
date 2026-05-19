# Verification

Status: Verified

Commands:

```bash
node -c tools/vibit
go test ./internal/platform/persistence/postgres
go test ./...
node tools/vibit check migrations --json
node tools/vibit check runtime --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check change confirm-next-direction-after-session-postgresql-adapter-gate --json
node tools/vibit check change implement-session-postgresql-adapter --json
node tools/vibit check all --json
git diff --check
```
