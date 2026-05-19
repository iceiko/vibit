# Plan

1. Replace fail-closed `LogoutAccessToken` service behavior with the gate-defined presented-token revocation sequence.
2. Reuse existing access-token proof decoding, lookup digest, verifier digest, comparison, token posture checks, unit-of-work, and repository vocabulary.
3. Add a redacted rejected logout result helper and a fixed revocation reason.
4. Extend authentication service tests for missing, malformed, lookup miss, token posture failures, verifier mismatch, dependency failures, commit failure, successful revocation, no raw token leakage, and no player/session repository side effects.
5. Record ADR-0073, change specs, conversation memory, manifests, AGENTS guides, rule catalog, and repository checks.
6. Run Go and repository verification.
