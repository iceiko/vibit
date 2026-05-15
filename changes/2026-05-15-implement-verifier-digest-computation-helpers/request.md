# Request

Implement the verifier digest computation helpers as the next bounded work item after the verifier digest helper implementation gate.

The change must stay helper-only. It may compute lookup and verifier digests from already-validated verifier keys and raw material, but it must not compare verifier digests, implement authentication service behavior, execute login, validate tokens, execute logout, add refresh behavior, add cleanup jobs, expose Protobuf messages, expose WebSocket proof carriers, change repositories, change migrations, wire startup, add dependencies, or add production authentication behavior.

## Acceptance Criteria

- Implement `runtime/internal/app/authentication/verifier_digest.go`.
- Add focused tests in `runtime/internal/app/authentication/verifier_digest_test.go`.
- Use deterministic canonical input with the `vibit.auth.verifier.input.v1` version header, null separator, length-prefixed purpose label, and length-prefixed raw material.
- Compute HMAC-SHA-256 with the matching logical key from `VerifierKeySet`.
- Provide separate helper functions for credential lookup, credential verifier, token lookup, and token verifier digests.
- Return copied 32-byte digest bytes.
- Keep digest classes and errors useful for tests but redacted.
- Fail closed for missing verifier key set and missing raw material.
- Update architecture manifests, authentication module metadata, agent guides, repository checks, work state, and conversation memory.
