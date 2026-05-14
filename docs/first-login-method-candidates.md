# First Login Method Candidate Comparison

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Candidate first login-method comparison for M-013
Depends on: `docs/login-method-token-format-ratification.md`

The paired Simplified Chinese translation is `docs/first-login-method-candidates.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This document compares the first login-method candidates for vibit before the repository ratifies the first production login-method set.

It does not implement authentication, add credential storage, add external identity linking, add token behavior, add session persistence, change the Protobuf envelope, change WebSocket handshake authentication, add runtime player handlers, or add WebSocket routes.

## 2. Evaluation Criteria

Each candidate is evaluated against:

- Production safety.
- Game onboarding ergonomics.
- Agent-readable implementation shape.
- Required public contracts.
- Required storage and migrations.
- Required dependencies.
- Abuse and recovery posture.
- Nakama capability alignment.
- Pitaya session and handler vocabulary alignment.
- Reversibility.

## 3. Candidate Summary

| Candidate | Recommendation | Reason |
| --- | --- | --- |
| `device_credential_login` | Recommend as first production player login method. | It gives a low-friction game login path without OAuth, OIDC, password hashing, provider SDKs, or WebSocket handshake changes. It can be made production-capable if the credential is high entropy, stored hashed, rate limited, and bound to account lifecycle through explicit contracts. |
| `guest_anonymous_login` | Defer. | Useful for onboarding, but easy to over-authorize. It requires strict anonymous actor rules, upgrade behavior, expiration, data ownership, and abuse posture before durable player state uses it. |
| `custom_id_login` | Defer except as a future trusted-issuer method. | Useful for studios with an existing identity service, but unsafe if arbitrary clients can mint IDs. It requires issuer trust, subject collision, account linking, and service-auth boundaries. |
| `email_password_login` | Defer. | Familiar but requires password hashing, recovery, rate limiting, reset flows, breach posture, and stronger secret-handling rules. |
| `external_provider_login` | Defer. | Needed later for platform accounts and social identity, but provider dependencies, issuer/audience validation, account linking, conflict handling, and availability make it too broad for the first slice. |
| `service_authentication` | Defer. | Important for future operations and server-to-server work, but it is not a player login method and should be separated from first player authentication. |

## 4. Recommended First Set

Recommended first login-method set:

```text
device_credential_login
```

This recommendation means:

- The first player login method should prove possession of a high-entropy device or installation credential.
- The credential must be treated as secret proof material, not as a public device identifier.
- The credential must not be placed in existing Protobuf `Session` metadata fields.
- The credential exchange should be modeled as an application/protocol login command payload unless a later protocol decision chooses another carrier.
- The method may create a new player account or authenticate an existing account only after account creation/linking policy is ratified.
- Credential storage, token issuance, token storage, session persistence, route registration, and runtime handlers remain future work.

This is not direct Nakama device-auth compatibility. It adapts Nakama's low-friction game login capability into a vibit-native model that is explicit about credential secrecy, storage, redaction, replay controls, and account lifecycle boundaries.

## 5. Candidate Details

### Device Credential Login

Position:

```text
recommend_first
```

Definition:

```text
A client proves possession of a high-entropy device or installation credential.
```

Benefits:

- Low-friction game onboarding.
- No external identity provider dependency.
- No password hashing dependency if the credential is generated as high-entropy secret material and stored through a ratified hash/lookup boundary.
- Works well with WebSocket-first gameplay because the credential can be exchanged through a login command before normal player routes.
- Fits a small first production slice better than email/password, social login, OAuth, OIDC, or provider SDKs.
- Gives agents a narrow path from contract to storage to validator to tests.

Risks:

- Raw OS device IDs are not sufficient proof.
- Weak identifiers are replayable.
- Device replacement, reinstall, account recovery, and account merge behavior need later designs.
- A leaked credential can impersonate the player until revocation or rotation exists.
- Rate limiting and abuse controls need explicit gates.

Required artifacts before implementation:

- Login command semantic contract.
- Login response semantic contract.
- Credential schema boundary.
- Credential hash and lookup rule.
- Account creation/linking policy.
- Token issuance boundary.
- Error and permission catalog entries.
- Redaction rules.
- Replay and collision tests.
- Repository checks preventing credential storage inside player account lifecycle tables.

Recommended status:

```yaml
candidate: device_credential_login
recommended_for_first_set: true
production_capable_after_required_gates: true
creates_player_account: allowed_after_policy_ratification
links_existing_account: deferred
authenticates_existing_account: allowed_after_policy_ratification
requires_major_dependency: false
requires_credential_storage: true
requires_external_identity_linking: false
requires_protobuf_envelope_change: false
requires_websocket_handshake_change: false
confidence: high
```

### Guest Or Anonymous Login

Position:

```text
defer
```

Benefits:

- Fastest onboarding path.
- Useful for try-before-register gameplay.
- Can be useful for local development and smoke tests.

Risks:

- Easy to confuse anonymous actor with a durable player.
- Abuse and spam controls matter immediately.
- Account upgrade behavior is product-sensitive.
- Anonymous durable state can create recovery and ownership disputes.

Required artifacts before implementation:

- Anonymous actor contract.
- Permission limits.
- Expiration and upgrade rules.
- Data ownership rules.
- Abuse and rate-limit posture.
- Tests proving anonymous identity cannot satisfy player-owned production permissions.

Recommendation:

```text
Do not include in the first production login-method set.
```

Guest or anonymous login may be useful later, but it should not be the first production proof path because vibit's current architecture is deliberately preventing metadata-only identity from becoming authority.

### Custom ID Login

Position:

```text
defer_trusted_issuer_variant
```

Benefits:

- Useful when a studio already has an identity service.
- Can be simple if called only by trusted services.
- Aligns with Nakama's custom identifier capability coverage.

Risks:

- Unsafe if arbitrary clients can choose their own custom IDs.
- Requires issuer trust boundaries.
- Requires subject collision rules.
- Requires account linking and recovery semantics.
- Often implies service authentication before player authentication.

Required artifacts before implementation:

- Trusted issuer model.
- Service-auth or caller-auth boundary.
- Subject namespace and collision rules.
- Account linking policy.
- Replay and audit behavior.

Recommendation:

```text
Defer until service-auth and issuer boundaries are ratified.
```

### Email Password Login

Position:

```text
defer
```

Benefits:

- Familiar to users.
- Supports cross-device account recovery.
- Common in general backend systems.

Risks:

- Requires password hashing dependency adoption.
- Requires password reset, recovery, breach posture, and rate limiting.
- Creates sensitive secret-handling and support workflows early.
- Larger public contract surface than the first slice needs.

Required artifacts before implementation:

- Password hash dependency adoption.
- Credential schema.
- Password policy.
- Reset and recovery contracts.
- Rate-limit and lockout policy.
- Redaction rules.
- Security tests.

Recommendation:

```text
Defer until the first credential/token/session slice is stable.
```

### External Provider Login

Position:

```text
defer
```

Benefits:

- Strong platform fit.
- Cross-device identity.
- Important for production games on platform stores and social platforms.
- Aligns with Nakama's broad provider coverage.

Risks:

- Provider SDKs and validation dependencies.
- Issuer, audience, key, and token validation complexity.
- Account linking conflicts.
- Provider outages and metadata retention rules.
- Different providers have different identity and token semantics.

Required artifacts before implementation:

- Provider namespace and subject schema.
- External identity link schema boundary.
- Dependency adoption records.
- Issuer and audience validation.
- Conflict, unlink, recovery, and merge behavior.
- Provider metadata redaction.

Recommendation:

```text
Defer until external identity linking is ratified.
```

### Service Authentication

Position:

```text
defer_separate_track
```

Benefits:

- Required for future operations, internal services, and server-to-server work.
- Can help trusted custom ID login later.

Risks:

- Not a player login method.
- Easy to over-privilege.
- Requires key management, rotation, permission scope, and audit behavior.
- May interact with future distributed runtime design.

Required artifacts before implementation:

- Service actor model.
- Service permission catalog.
- Key or proof material lifecycle.
- Rotation and revocation semantics.
- Audit events.

Recommendation:

```text
Defer to a separate service-auth milestone or sub-milestone.
```

## 6. Comparative Matrix

Scores are qualitative:

```text
high = favorable
medium = manageable with gates
low = unfavorable for first slice
```

| Candidate | Production safety | Game ergonomics | Agent clarity | Dependency load | Storage complexity | Abuse/recovery load | First-slice fit |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `device_credential_login` | Medium | High | High | High | Medium | Medium | High |
| `guest_anonymous_login` | Low | High | Medium | High | Medium | Low | Medium |
| `custom_id_login` | Medium only with trusted issuer | Medium | Medium | Medium | Medium | Medium | Medium |
| `email_password_login` | Medium | Medium | Medium | Low | Medium | Low | Low |
| `external_provider_login` | Medium | High | Low | Low | Medium | Medium | Low |
| `service_authentication` | Medium | Not player-facing | Medium | Medium | Medium | Medium | Low for player login |

Interpretation:

- `device_credential_login` is the best first candidate because its complexity is mostly local and contract-checkable.
- `guest_anonymous_login` is tempting but weakens the project's current discipline unless anonymous permissions and upgrade behavior are carefully designed first.
- `custom_id_login` is valuable later, but only after trusted issuer semantics exist.
- `email_password_login` and `external_provider_login` are important product capabilities but too broad for the first authentication slice.
- `service_authentication` should be separated from player login.

## 7. Reference Mapping

### Nakama

Nakama supports multiple authentication methods, including device, email, social/provider, and custom identifier methods. It also produces session tokens after authentication and supports refresh-token based session continuation.

vibit adapts this as:

- First capability target: low-friction player login.
- First recommended method: high-entropy device credential login.
- Deferred capability targets: email/password, social/provider login, custom identifier login, guest/anonymous login, and refresh-token lifecycle.
- Rejected for now: direct Nakama API compatibility.

### Pitaya

Pitaya's useful input here is vocabulary around sessions, request handlers, session binding, frontend/backend separation, route handling, and push.

vibit adapts this as:

- Future validated request identity is passed to handlers before domain dispatch.
- Future session binding must happen through application-owned validation results, not transport metadata.
- First login method selection must not put authentication inside WebSocket acceptors.
- Frontend/backend and distributed session routing remain deferred.

## 8. Recommended W-0066 Ratification Packet

W-0066 should ratify:

```yaml
first_login_method_set:
  - device_credential_login
deferred_login_method_families:
  - guest_anonymous_login
  - custom_id_login
  - email_password_login
  - external_provider_login
  - service_authentication
recommended_first_carrier_posture: login_command_payload_before_normal_gameplay_routes
requires_before_implementation:
  - semantic_login_contract
  - credential_schema_boundary
  - credential_hash_lookup_rule
  - account_creation_or_lookup_policy
  - token_issuance_boundary
  - error_catalog_entries
  - permission_catalog_entries
  - redaction_rules
  - repository_checks
  - focused_tests
does_not_authorize:
  - runtime_authentication_code
  - token_parsing
  - credential_tables
  - external_identity_tables
  - token_tables
  - session_tables
  - protobuf_envelope_change
  - websocket_handshake_authentication
  - runtime_player_handlers
  - websocket_routes
```

## 9. Open Questions For Later Work

- Whether the initial credential is generated by the client, issued by the server during a bootstrap exchange, or both.
- Whether the login command creates a player account by default or requires an explicit create intent.
- Whether credential rotation is part of the first implementation or deferred behind token lifecycle work.
- Whether account recovery exists in the first implementation or is explicitly deferred.
- Which token model W-0067 and W-0068 will select after the login-method set is ratified.

These are not blockers for W-0065. They should be answered by W-0066 through W-0071.
