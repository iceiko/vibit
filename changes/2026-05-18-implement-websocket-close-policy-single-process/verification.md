# Verification

Required verification:

```text
go test ./internal/app/connection
go test ./...
node -c tools/vibit
node tools/vibit check schemas --json
node tools/vibit check runtime --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check change confirm-next-direction-after-websocket-close-policy-gate --json
node tools/vibit check change implement-websocket-close-policy-single-process --json
node tools/vibit check change confirm-next-direction-after-websocket-close-policy-implementation --json
node tools/vibit check all --json
git diff --check
node tools/vibit inspect next --json
```

Status: Verified.
