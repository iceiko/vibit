# Plan

1. Implement class-specific verifier digest comparison helpers under `runtime/internal/app/authentication/verifier_comparison.go`.
2. Add focused unit tests under `runtime/internal/app/authentication/verifier_comparison_test.go`.
3. Update `tools/vibit` so the W-0105 helper file is authorized while verifier comparison remains forbidden elsewhere.
4. Update architecture manifests, authentication module metadata, AGENTS guides, work items, reference planning, and conversation memory.
5. Run Go tests and repository verification.

## Non-Goals

- No authentication service behavior.
- No login, token validation, logout, refresh, or cleanup behavior.
- No Protobuf or WebSocket authentication carriers.
- No repository or migration changes.
- No startup wiring.
- No dependency changes.
- No production authentication behavior.
