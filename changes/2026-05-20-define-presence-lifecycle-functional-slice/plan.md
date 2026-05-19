# Plan

1. Keep W-0177 as a Tier 2 functional slice under `ADR-0082`.
2. Define the embedded presence lifecycle boundary in `spec.yaml`.
3. Add registry-backed player presence snapshots under `runtime/internal/app/connection`.
4. Add a credential-neutral WebSocket connection lifecycle observer interface.
5. Wire PostgreSQL startup composition from WebSocket open/close and successful connection binding into the active connection registry.
6. Add focused tests for registry presence, lifecycle observation, and startup binding.
7. Update manifests, guides, ADR, conversation memory, check catalog, and `tools/vibit`.
8. Verify focused tests, full runtime tests, and repository checks.
9. Commit and push.

