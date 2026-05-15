# Request

Implement the verifier digest comparison helpers as the next bounded work item after the verifier digest comparison helper gate.

The change must stay helper-only. It may compare computed verifier digest bytes to stored verifier digest bytes with a constant-time primitive, but it must not compute digests, load keys, select records, call repositories, inspect lifecycle state, parse proofs, issue login responses, validate tokens, execute logout, refresh tokens, expose Protobuf messages, expose WebSocket proof carriers, change repositories, change migrations, wire startup, add dependencies, or add production authentication behavior.

## Acceptance Criteria

- Implement `runtime/internal/app/authentication/verifier_comparison.go`.
- Add focused tests in `runtime/internal/app/authentication/verifier_comparison_test.go`.
- Use `crypto/hmac.Equal` for constant-time verifier digest comparison.
- Compare only `ComputedDigest` verifier bytes against stored verifier digest bytes.
- Provide class-specific helpers for credential verifier digests and token verifier digests.
- Return a redacted match or mismatch result.
- Reject lookup digest classes, wrong verifier classes, missing computed input, malformed computed input, missing stored input, and malformed stored input.
- Keep all errors redacted.
- Update architecture manifests, authentication module metadata, agent guides, repository checks, work state, and conversation memory.
