# Bound Identity Route Policy Gate

Status: Draft v0.1
Last updated: 2026-05-18
Scope: Gate-only boundary for future route-policy use of request-level proof, durable runtime session validation, and first-message bound connection identity
Depends on: `docs/access-token-protocol-carrier-route-protection-gate.md`, `docs/first-message-connection-binding-implementation-gate.md`, `docs/runtime-session-validation-gate.md`, `docs/session-creation-composition-gate.md`, `decisions/ADR-0068-session-creation-composition-implementation.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0069`

The paired Simplified Chinese translation is `docs/bound-identity-route-policy-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit now has several identity-related runtime pieces:

- Request-level access-token route protection through an explicit Protobuf payload wrapper.
- First-message connection binding through `runtime.authentication.BindConnection`.
- Durable `runtime_sessions` persistence and a session repository.
- An application-owned persistent runtime session validator.
- Login-time durable runtime session creation composed with device-credential login.

Those pieces are intentionally not yet one route policy. The current protected-route posture still validates access-token proof per protected request. The bound connection identity and persisted session validation exist as separate application capabilities, but normal domain routes do not yet rely on them.

The next useful step is to define the future route-policy boundary before implementation. Mature game servers shape this boundary:

- Nakama treats authenticated session material as the basis for gameplay API access, while keeping session lifetime, logout, refresh, and socket behavior as distinct lifecycle decisions.
- Pitaya keeps acceptors, sessions, and route handlers separated so handlers receive context rather than parsing transport credentials.

vibit should adapt those lessons by making route identity policy application-owned, explicit, testable, and staged. This standard defines the gate only.

```yaml
bound_identity_route_policy_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0150
decision: ADR-0069
check_rule: runtime.bound_identity_route_policy_gate
future_policy_owner: runtime/internal/app
future_policy_source_candidate: runtime/internal/app/route_authentication.go
future_policy_test_candidate: runtime/internal/app/route_authentication_test.go
existing_access_token_proof_path: request_level_authenticated_request_wrapper
existing_connection_binding_route: runtime.authentication.BindConnection
existing_session_validator: runtime/internal/app.PersistentSessionValidator
route_policy_session_identity_added: false
route_policy_bound_identity_added: false
ordinary_protected_routes_use_bound_identity: false
ordinary_protected_routes_use_session_validated_identity: false
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
protobuf_session_messages_added: false
existing_protobuf_envelope_change_added: false
logout_revocation_active_connection_added: false
reconnect_epoch_behavior_added: false
cleanup_jobs_added: false
dependencies_added: false
memory_durable_session_behavior_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

This is a gate-only standard. It does not change production route authorization behavior.

## 2. Ownership

Future bound identity route policy is application-owned:

```yaml
future_route_policy_owner: runtime/internal/app
authentication_service_owner: runtime/internal/app/authentication
session_validator_owner: runtime/internal/app
session_repository_owner: runtime/internal/app/session
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
transport_owner: runtime/internal/platform/transport/ws
domain_handler_owner: runtime/internal/modules/*
```

Rules:

- Route policy must stay under `runtime/internal/app`.
- Route policy may consume normalized `RequestIdentity`, access-token validation results, connection binding identity, and runtime session validation results only through application-owned types.
- Route policy must not import WebSocket transport packages, generated Protobuf packages, PostgreSQL adapters, SQL rows, migration packages, or provider SDKs.
- WebSocket transport must remain credential-neutral.
- Protobuf adapters may carry or unwrap proof only through already-authorized protocol carrier boundaries; they must not decide policy.
- Domain modules must receive normalized identity context and must not parse access tokens, session ids, connection ids, or transport credential carriers.

## 3. Future Policy Families

A later implementation slice may introduce explicit route policy families instead of a single implicit protected-route default:

```yaml
candidate_policy_families:
  public:
    examples:
      - runtime.authentication.AuthenticateWithDeviceCredential
    requirement: no authenticated proof required
  request_token_required:
    requirement: fresh request-level access-token proof validates for the route
  bound_connection_required:
    requirement: connection has a server-observed bound identity matching the request identity
  session_validated_required:
    requirement: already-validated actor identity also validates against an active durable runtime session row
  bound_session_required:
    requirement: bound connection identity and active durable session validation both match the request identity
```

Rules:

- Public routes must remain explicit.
- There must be no implicit public gameplay route default.
- Metadata-only identity must not satisfy any protected policy family.
- A bound connection identity alone must not become a universal substitute for request-level proof unless a later implementation slice explicitly chooses that posture for named routes.
- A persisted `session_id` alone must not become proof.
- Session-validated identity requires an already-validated actor identity and an active durable runtime session row.
- Policy family selection must be route-scoped and testable without a live PostgreSQL server.

## 4. Recommended First Implementation Posture

The recommended first implementation after this gate is conservative:

```yaml
recommended_first_implementation_posture:
  default_domain_route_policy: request_token_required
  public_routes:
    - runtime.authentication.AuthenticateWithDeviceCredential
  system_binding_route:
    route: runtime.authentication.BindConnection
    policy: request_token_required_for_binding_request
  bound_connection_policy:
    status: available_for_explicit_future_routes_only
  session_validated_policy:
    status: available_for_explicit_future_routes_only
  bound_session_policy:
    status: deferred
```

Rules:

- Ordinary protected domain routes should continue to require request-level access-token proof in the first implementation.
- `BindConnection` may establish process-local bound identity, but that does not automatically authorize ordinary domain routes.
- The first route-policy implementation may add policy vocabulary and tests for future route classes without changing existing domain route behavior.
- Any route that stops requiring per-request proof must be named explicitly and must have tests showing which identity source is trusted.
- Inventory route behavior must not change accidentally as a side effect of this gate.

## 5. Composition Order

Future route policy must compose identity checks in a deterministic order:

```yaml
candidate_composition_order:
  - normalize_route_key
  - classify_route_policy_family
  - accept_explicit_public_route_without_authenticated_identity
  - require_structural_proof_or_bound_identity_for_protected_family
  - validate_access_token_proof_when_required
  - validate_bound_connection_identity_when required_by_route_family
  - validate_runtime_session_when_required_by_route_family
  - build_one_normalized_request_identity
  - dispatch_domain_handler_only_after_policy_success
```

Rules:

- Domain dispatch must happen only after policy success.
- If multiple identity sources are used, they must agree on actor kind, actor id, player id, session id when required, connection id when required, and connection epoch when required.
- Mismatched identity sources must fail closed.
- Policy failure must not mutate token state, session state, connection binding state, inventory state, or domain state.
- `SessionValidated = true` may be accepted only when produced by the runtime session validator; route policy must not set it by assertion.
- Future last-seen updates remain outside this gate unless a later implementation explicitly authorizes them.

## 6. Error And Redaction Boundary

Future implementation must keep public failures stable and redacted:

```yaml
candidate_public_errors:
  request_token_missing: AUTHENTICATION_TOKEN_MISSING
  request_token_malformed: AUTHENTICATION_TOKEN_MALFORMED
  request_token_invalid: AUTHENTICATION_TOKEN_INVALID
  request_token_unavailable: AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
  bound_connection_required: CONNECTION_BINDING_REQUIRED
  bound_connection_invalid: CONNECTION_BINDING_TOKEN_INVALID
  session_invalid: SESSION_INVALID
```

Rules:

- Public failures must not reveal whether token lookup, token mismatch, account inactive, session missing, session expired, session revoked, connection unbound, epoch mismatch, actor mismatch, player mismatch, repository failure, or policy-family mismatch caused denial.
- Errors, logs, events, and test output must not include raw access-token text, raw credential material, lookup digests, verifier digests, verifier key ids, session ids, token record ids, Authorization headers, cookies, query strings, WebSocket subprotocol values, or inner payload bytes.
- Route policy should return route-specific public errors only at the application error boundary.
- More precise internal reasons may exist in tests only if they stay redacted from public output.

## 7. Relationship To Existing Pieces

This gate does not change existing behavior:

```yaml
access_token_validation_changed: false
connection_binding_changed: false
runtime_session_validation_changed: false
session_creation_changed: false
request_identity_session_validated_policy_changed: false
route_policy_session_identity_added: false
route_policy_bound_identity_added: false
```

Rules:

- `RouteProtector` may continue to require request-level access-token proof for ordinary protected routes until a later implementation work item changes it.
- `ConnectionBinder` may continue to set `SessionValidated = false`.
- `PersistentSessionValidator` may continue to validate only when called directly by application code.
- `AuthenticateWithDeviceCredential` may continue creating durable runtime sessions without exposing those session ids through protocol responses.
- The PostgreSQL session adapter remains persistence-only.
- The authentication service remains proof/token/session-creation composition; it must not own route policy.

## 8. Relationship To WebSocket And Protocol

This gate does not authorize WebSocket or Protobuf changes:

```yaml
websocket_transport_credential_neutral: true
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
protobuf_session_messages_added: false
existing_protobuf_envelope_change_added: false
generated_output_added: false
```

Rules:

- WebSocket transport must not parse Authorization headers, bearer values, cookies, query-string tokens, session tokens, or `Sec-WebSocket-Protocol` authentication material for this gate.
- The existing Protobuf envelope remains unchanged.
- No session carrier field, login response session id, handshake message, reconnect message, route policy Protobuf message, generated Protobuf output, or generated contract shape is authorized here.
- A future protocol session carrier gate must authorize any client-visible session id or session proof before route policy requires that client input.

## 9. Deferrals

This gate does not authorize:

- Production route-policy use of bound identity.
- Production route-policy use of session-validated identity.
- Removing per-request access-token proof from ordinary protected routes.
- WebSocket handshake authentication.
- Transport credential carriers.
- Protobuf envelope changes.
- Protobuf session messages or generated output.
- Logout, refresh, cleanup, token rotation, token-session rekeying, or token validation audit mutation.
- Active-connection invalidation on logout or revocation.
- Reconnect, resume, duplicate connection replacement, or durable connection epoch policy.
- Presence, rooms, parties, groups, chat, matchmaking, match runtime, broadcast groups, or social modules.
- Metrics, dashboards, admin APIs, session management APIs, or operations posture.
- Memory durable session behavior.
- Direct Nakama or Pitaya API compatibility.

## 10. Test Requirements For Future Implementation

A later bound identity route-policy implementation must include focused tests for:

- Public routes remain explicit and require no proof.
- Ordinary protected routes still require request-level proof unless explicitly reclassified.
- Metadata-only identity is rejected for all protected policy families.
- Bound identity can satisfy only routes explicitly classified for bound identity.
- Session-validated identity can satisfy only routes explicitly classified for session validation.
- Bound and session identity sources must match actor/player/session/connection fields when a route requires both.
- Mismatches fail closed before domain dispatch.
- Policy failures remain redacted.
- WebSocket transport remains credential-neutral.
- Protobuf envelope and authentication response shapes remain unchanged.
- Inventory route behavior does not change unless explicitly reclassified.

Live PostgreSQL verification may remain opt-in unless the later implementation work item requires real persistence.

## 11. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

```yaml
adopted_concepts:
  - authenticated_session_material_controls_gameplay_access
  - session_lifetime_and_route_authorization_are_related
  - logout_refresh_and_session_management_have_future_policy_pressure
adapted_concepts:
  - vibit_keeps_request_token_bound_connection_and_durable_session_identity_as_explicit_policy_families
  - vibit_does_not_copy_nakama_session_api_or_jwt_shape
  - vibit_requires_route_scoped_policy_before_trusting_bound_or_session_identity
deferred_concepts:
  - refresh_token_flow
  - session_management_api
  - logout_disconnect_active_socket
  - single_session_or_single_socket_policy
rejected_for_now:
  - direct_nakama_api_compatibility
```

Pitaya reference mapping:

```yaml
adopted_concepts:
  - acceptor_transport_and_handler_logic_are_separate
  - sessions_can_bind_user_identity_to_connection_context
  - handlers_should_receive_context_not_parse_credentials
adapted_concepts:
  - vibit_route_policy_is_application_owned_and_contract_first
  - vibit_bound_identity_is_not_implicitly_all_routes_authenticated
  - vibit_keeps_cluster_session_routing_deferred
deferred_concepts:
  - frontend_backend_cluster_session_routing
  - remote_session_binding_broadcast
  - groups_rooms_and_presence_attachment
  - server_to_server_rpc
rejected_for_now:
  - direct_pitaya_api_compatibility
```

Nakama and Pitaya guide capability planning. They do not override vibit's constitution, ADRs, manifests, generated boundaries, or verification commands.

## 12. Future Implementation Queue

After this gate, future work should remain split:

```yaml
future_work_items:
  bound_identity_route_policy_implementation:
    may_add:
      - explicit application route policy families
      - route-scoped classification for public, request-token, bound-connection, and session-validated policies
      - tests for fail-closed identity agreement
  logout_revocation_active_connection_gate:
    may_define:
      - whether token/session revocation closes active WebSocket connections
  reconnect_connection_epoch_gate:
    may_define:
      - reconnect, resume, duplicate replacement, and epoch mismatch behavior
  protocol_session_carrier_gate:
    may_define:
      - whether and how clients receive or carry session ids
```

Do not combine these into one broad session subsystem slice without a new ADR.

## 13. Verification

The repository check rule for this gate is:

```text
runtime.bound_identity_route_policy_gate
```
