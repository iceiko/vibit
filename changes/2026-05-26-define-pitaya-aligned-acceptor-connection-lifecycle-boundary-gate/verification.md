# Verification

Required before completion:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_acceptor_connection_lifecycle_boundary_gate
node tools/vibit check change define-pitaya-aligned-acceptor-connection-lifecycle-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Status: Verified after the commands above pass.
