# Plan

1. Implement digest classes, computed digest value object, typed errors, canonical input construction, and helper functions under `runtime/internal/app/authentication/verifier_digest.go`.
2. Add focused unit tests under `runtime/internal/app/authentication/verifier_digest_test.go`.
3. Update `tools/vibit` to allow only this helper slice while preserving deferrals in other authentication boundaries.
4. Update architecture manifests, authentication module metadata, AGENTS guides, work items, reference planning, and conversation memory.
5. Run Go tests and repository verification.

## Non-Goals

- No verifier comparison.
- No authentication service behavior.
- No login, token validation, logout, refresh, or cleanup behavior.
- No Protobuf or WebSocket authentication carriers.
- No repository or migration changes.
- No startup wiring.
- No dependency changes.
- No production authentication behavior.
