# Plan

1. Implement a narrow environment verifier key loader under `runtime/internal/app/authentication`.
2. Keep the testable loader behind an explicit lookup function and add a tiny process environment adapter.
3. Decode only Base64URL unpadded and standard padded Base64 key text.
4. Call `NewVerifierKeySet` for key-set validation instead of duplicating invariant logic.
5. Add focused tests for success, missing input, invalid encoding, validator handoff failures, redaction, and process environment adapter behavior.
6. Update architecture manifests, agent guides, runtime checks, work queue state, and verification records.
