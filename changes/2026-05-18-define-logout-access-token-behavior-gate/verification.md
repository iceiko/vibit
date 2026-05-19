# Verification

Status: Verified

Commands:

```bash
node -c tools/vibit
node tools/vibit check schemas --json
node tools/vibit check runtime --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check change confirm-next-direction-after-logout-revocation-active-connection-gate --json
node tools/vibit check change define-logout-access-token-behavior-gate --json
node tools/vibit check all --json
git diff --check
```
