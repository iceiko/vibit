# Plan

1. Close `M-089/W-0161` by selecting `implement_active_connection_registry_single_process`.
2. Create `M-090/W-0162` as the bounded implementation slice.
3. Create `M-091/W-0163` as the next confirmation gate.
4. Add `runtime/internal/app/connection/registry.go`.
5. Add `runtime/internal/app/connection/registry_test.go`.
6. Implement register, bind, close, invalidate, find, and list capabilities.
7. Add focused tests for validation, duplicate active records, metadata-only rejection, terminal-state exclusion, copy semantics, and record redaction.
8. Add `ADR-0075`, change specs, and conversation memory.
9. Update architecture manifests, module manifest, AGENTS guides, rules, and checks.
10. Run Go and repository verification.
