# Verification

Required verification:

```sh
node -c tools/vibit
node tools/vibit inspect rule runtime.pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate
node tools/vibit inspect next --json
node tools/vibit check change define-pitaya-aligned-session-binding-kick-disconnect-session-data-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
git diff --check
node tools/vibit check all --json
```

Status: Verified.
