# Impact

## Changed

- Adds a verifier algorithm and redaction-test boundary standard.
- Adds ADR-0040.
- Adds a conversation log for the work item.
- Updates architecture manifests and agent guides so future agents can discover the boundary.
- Adds a runtime check rule for the boundary.
- Completes `W-0091` and prepares the next bounded work item.

## Security Impact

This change improves security posture by defining:

- `vibit_hmac_sha256_v1` as the first planned high-entropy verifier algorithm family.
- Separate lookup and verifier digest classes.
- Purpose-label domain separation.
- Minimum 256-bit entropy for raw access tokens and raw device credentials.
- Constant-time verifier digest comparison expectations.
- Secret and digest redaction requirements.
- Redaction test expectations for future implementation work.

No runtime authentication behavior is added.

## Dependency Impact

No dependency is added.

The first planned verifier posture can use Go standard library packages after a later implementation gate:

- `crypto/hmac`
- `crypto/sha256`
- `crypto/subtle`
- `crypto/rand`
- `encoding/base64`

No external cryptography, password-hashing, JWT, OAuth, OIDC, KMS, provider, or Redis-like token/session dependency is required for the first high-entropy posture.

## Runtime Impact

No Go runtime code is added or modified.

The following remain deferred:

- token generation
- credential generation
- verifier comparison
- application authentication service code
- login execution
- token validation
- logout execution
- cleanup jobs
- Protobuf authentication messages
- WebSocket proof carriers
- authentication dependencies
- repository interface changes
- migration schema changes
- secret configuration loading

## Documentation Impact

English and Simplified Chinese public documentation are updated together.
