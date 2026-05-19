# Verification

Status: Verified

Commands:

```bash
node -c tools/vibit
go test ./...
node tools/vibit check runtime --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check change confirm-next-direction-after-runtime-session-validation-implementation --json
node tools/vibit check change define-session-creation-composition-gate --json
node tools/vibit check all --json
git diff --check
```
