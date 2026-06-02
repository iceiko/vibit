# Verification

Status: verified during implementation.

RED evidence:

- `node tools/vibit inspect rule runtime.pitaya_aligned_startup_shutdown_hook_boundary_gate` failed with `Unknown rule_id: runtime.pitaya_aligned_startup_shutdown_hook_boundary_gate`.
- `node tools/vibit check change define-pitaya-aligned-startup-shutdown-hook-boundary-gate --json` failed because the change directory did not exist.

Expected verification:

- `node -c tools/vibit`
- `node tools/vibit inspect rule runtime.pitaya_aligned_startup_shutdown_hook_boundary_gate`
- `node tools/vibit check change define-pitaya-aligned-startup-shutdown-hook-boundary-gate --json`
- `node tools/vibit inspect next --json`
- `node tools/vibit check work --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check all --json`
- `git diff --check`

Observed verification will be recorded after the completing agent run.
