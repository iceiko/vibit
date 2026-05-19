# Verification

Status: Verified

Commands:

```bash
go test ./internal/app/connection
go test ./...
node -c tools/vibit
node tools/vibit check schemas --json
node tools/vibit check runtime --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check change confirm-next-direction-after-active-connection-registry-gate --json
node tools/vibit check change implement-active-connection-registry-single-process --json
node tools/vibit check all --json
git diff --check
```
