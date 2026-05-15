# Plan

1. Implement the material generation value object, material kind type, typed errors, and helper functions under `runtime/internal/app/authentication/material_generation.go`.
2. Add focused unit tests under `runtime/internal/app/authentication/material_generation_test.go`.
3. Update `tools/vibit` to allow only this helper slice while preserving deferrals in other authentication boundaries.
4. Update architecture manifests, authentication module metadata, AGENTS guides, work items, reference planning, and conversation memory.
5. Run Go tests and repository verification.

## Non-Goals

- No digest computation.
- No verifier comparison.
- No authentication service behavior.
- No login, token validation, logout, refresh, or cleanup behavior.
- No Protobuf or WebSocket authentication carriers.
- No repository or migration changes.
- No startup wiring.
- No dependency changes.
- No production authentication behavior.
