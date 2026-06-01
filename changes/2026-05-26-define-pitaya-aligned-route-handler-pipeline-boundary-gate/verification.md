# Verification

Status: Verified

Commands:

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.pitaya_aligned_route_handler_pipeline_boundary_gate`
- `node tools/vibit check change define-pitaya-aligned-route-handler-pipeline-boundary-gate --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`
