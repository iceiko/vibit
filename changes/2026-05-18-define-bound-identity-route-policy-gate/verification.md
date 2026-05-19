# Verification

Status: Verified

Commands:

```bash
node -c tools/vibit
node tools/vibit check runtime --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check change confirm-next-direction-after-session-creation-composition-implementation --json
node tools/vibit check change define-bound-identity-route-policy-gate --json
node tools/vibit check all --json
git diff --check
```
