# Request

Implement the token and credential material generation helpers as the next bounded work item after the implementation gate.

The change must stay helper-only. It must not compute lookup digests, compute verifier digests, compare verifier material, implement authentication service behavior, execute login, validate tokens, execute logout, add refresh behavior, add cleanup jobs, expose Protobuf messages, expose WebSocket proof carriers, change repositories, change migrations, wire startup, add dependencies, or add production authentication behavior.

## Acceptance Criteria

- Implement `runtime/internal/app/authentication/material_generation.go`.
- Add focused tests in `runtime/internal/app/authentication/material_generation_test.go`.
- Use explicit `io.Reader` entropy-source handoff.
- Read 32 raw bytes.
- Encode URL-safe unpadded Base64 presentation text.
- Return copied raw bytes.
- Preserve distinct material kind values for device credentials and access tokens.
- Fail closed for nil reader, read error, short read, all-zero material, and repeated single-byte material.
- Keep errors redacted.
- Update architecture manifests, agent guides, repository checks, work state, and conversation memory.
