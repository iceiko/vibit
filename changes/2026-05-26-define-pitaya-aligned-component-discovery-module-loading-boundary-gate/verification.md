# Verification

Required verification commands:

```text
node -c tools/vibit
node tools/vibit inspect rule runtime.pitaya_aligned_component_discovery_module_loading_boundary_gate
node tools/vibit inspect next --json
node tools/vibit check change define-pitaya-aligned-component-discovery-module-loading-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Verification status is recorded by the agent run that completes this change.
