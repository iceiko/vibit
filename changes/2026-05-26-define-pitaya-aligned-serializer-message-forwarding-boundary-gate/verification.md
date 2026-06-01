# Verification

Required commands:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_serializer_message_forwarding_boundary_gate
node tools/vibit check change define-pitaya-aligned-serializer-message-forwarding-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Status: Verified.
