# Verification

Required verification:

```text
go test ./...
node -c tools/vibit
node tools/vibit check schemas --json
node tools/vibit check runtime --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check change confirm-next-direction-after-logout-access-token-behavior-gate --json
node tools/vibit check change implement-logout-access-token-behavior --json
node tools/vibit check all --json
git diff --check
```

Status: Verified.
