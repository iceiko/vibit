# Plan

1. Capture RED checks for the missing startup/shutdown hook source-first rule, command, and change directory.
2. Add `node tools/vibit inspect pitaya-startup-shutdown --json`.
3. Register `runtime.pitaya_aligned_startup_shutdown_hook_source_first_map`.
4. Add runtime check coverage for the command output, required artifacts, manifest references, redaction posture, and deferrals.
5. Add `ADR-0198`, conversation memory, and change artifacts.
6. Update `.arch/work-items.yaml` from `W-0290` to `W-0291`.
7. Update `.arch/runtime.yaml`, `.arch/reference.yaml`, `.arch/conventions.yaml`, `.arch/contracts.yaml`, `.arch/modules.yaml`, module manifests, README/AGENTS files, alpha docs, product maturity docs, and roadmap docs.
8. Run targeted and full verification.
