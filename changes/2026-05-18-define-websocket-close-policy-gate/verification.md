# Verification

Required verification:

```text
node -c tools/vibit
node tools/vibit check schemas --json
node tools/vibit check runtime --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check change confirm-next-direction-after-active-connection-registry-implementation --json
node tools/vibit check change define-websocket-close-policy-gate --json
node tools/vibit check change confirm-next-direction-after-websocket-close-policy-gate --json
node tools/vibit check all --json
git diff --check
```

Status: Verified.
