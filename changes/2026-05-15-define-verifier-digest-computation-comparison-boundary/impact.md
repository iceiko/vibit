# Impact

## Changed

- Adds a verifier digest computation and comparison boundary standard.
- Adds ADR-0043.
- Adds a conversation log for the work item.
- Updates architecture manifests and agent guides so future agents can discover the boundary.
- Adds a runtime check rule for the boundary.
- Completes `W-0094` and prepares the next bounded work item.

## Security Impact

This change improves security posture by defining:

- Application ownership for future digest computation and verifier comparison.
- Canonical HMAC input framing with version header and length prefixes.
- Required separation of lookup and verifier digest classes.
- Required separation of credential and token digest classes.
- Active and accepted-previous key set posture for opaque proof validation.
- Lookup digest equality as candidate record selection only.
- Constant-time comparison for verifier digest bytes.
- Public failure collapse for lookup miss, mismatch, unknown key id, unsupported algorithm, malformed proof, and expired or revoked proof.
- Redaction requirements for raw material, encoded material, digest material, key material, key identifiers, HMAC inputs, and HMAC outputs.
- Future digest and comparison test expectations.

No runtime authentication behavior is added.

## Dependency Impact

No dependency is added.

Future first-posture digest helpers may use Go standard library `crypto/hmac`, `crypto/sha256`, and `crypto/subtle` after a later code gate.

No external cryptography, password-hashing, JWT, OAuth, OIDC, provider, KMS, cloud secret-manager, or operations dependency is required for the first digest computation and comparison posture.

## Runtime Impact

No Go runtime code is added or modified.

The following remain deferred:

- verifier digest computation
- verifier comparison
- token generation
- credential generation
- secret loading
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
- production authentication behavior

## Documentation Impact

English and Simplified Chinese public documentation are updated together.
