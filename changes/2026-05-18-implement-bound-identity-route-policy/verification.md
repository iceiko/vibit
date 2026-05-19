# Verification

Status: Verified

Commands:

```bash
node -c tools/vibit
go test ./internal/app
go test ./...
node tools/vibit check runtime --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check change confirm-next-direction-after-bound-identity-route-policy-gate --json
node tools/vibit check change implement-bound-identity-route-policy --json
node tools/vibit check all --json
git diff --check
```
