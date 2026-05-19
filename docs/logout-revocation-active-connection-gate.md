# Logout Revocation Active Connection Gate

Status: Draft v0.1
Last updated: 2026-05-18
Scope: Gate-only boundary for future logout, token revocation, runtime session revocation, and active WebSocket connection invalidation behavior
Depends on: `docs/authentication-contract-error-permission-surfaces.md`, `docs/session-persistence-websocket-handshake-ratification.md`, `docs/first-message-connection-binding-implementation-gate.md`, `docs/runtime-session-validation-gate.md`, `docs/session-creation-composition-gate.md`, `docs/bound-identity-route-policy-gate.md`, `decisions/ADR-0070-bound-identity-route-policy-implementation.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0071`

The paired Simplified Chinese translation is `docs/logout-revocation-active-connection-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit now has enough authentication and session pieces that logout and revocation can no longer be treated as a storage-only detail:

- Request-level opaque access-token validation.
- A public device-credential login command.
- Durable access-token verifier records.
- Durable `runtime_sessions` rows created at successful login.
- First-message connection binding that can associate an observed WebSocket connection with a validated player identity.
- Explicit route policy families for request-token, bound-connection, session-validated, and bound-session routes.

The missing boundary is what a future logout or revocation operation should do to active WebSocket connections. Mature game servers expose the pressure clearly:

- Nakama treats authenticated session material as a lifecycle object with logout, refresh, expiration, and realtime socket implications, while still keeping token/session validity and active socket behavior as distinct implementation concerns.
- Pitaya separates acceptors, session context, handlers, and connection management; session binding and kick-like behavior are connection/session lifecycle concerns, not handler-owned credential parsing.

vibit should adapt those lessons by defining the policy boundary before any code can revoke tokens, revoke runtime sessions, close sockets, or add a connection registry. This standard is gate-only.

```yaml
logout_revocation_active_connection_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0154
decision: ADR-0071
check_rule: runtime.logout_revocation_active_connection_gate
future_policy_owner: runtime/internal/app
future_authentication_service_owner: runtime/internal/app/authentication
future_session_repository_owner: runtime/internal/app/session
future_connection_registry_owner_candidate: runtime/internal/app
future_connection_registry_package_candidate: runtime/internal/app/connection
future_transport_invalidation_interface_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
existing_logout_contract: contracts/runtime/authentication/commands/LogoutAccessToken.yaml
existing_logout_service_method: runtime/internal/app/authentication.Service.LogoutAccessToken
existing_logout_service_behavior: fail_closed_not_implemented
existing_refresh_service_behavior: refresh_not_supported
first_recommended_logout_scope: presented_access_token_only
token_revocation_and_session_revocation_separate: true
connection_registry_required_before_targeting_active_connections: true
logout_execution_added: false
token_revocation_execution_added: false
runtime_session_revocation_execution_added: false
active_connection_invalidation_added: false
connection_registry_added: false
websocket_close_policy_added: false
transport_credential_carrier_added: false
websocket_handshake_authentication_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
existing_protobuf_envelope_change_added: false
reconnect_epoch_behavior_added: false
cleanup_jobs_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

This standard does not implement logout or revocation behavior.

## 2. Ownership

Future logout and revocation policy must stay application-owned:

```yaml
future_policy_owner: runtime/internal/app
future_logout_execution_owner: runtime/internal/app/authentication
future_token_record_owner: runtime/internal/modules/authentication
future_session_repository_owner: runtime/internal/app/session
future_connection_registry_owner_candidate: runtime/internal/app
transport_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
domain_handler_owner: runtime/internal/modules/*
postgres_adapter_owner: runtime/internal/platform/persistence/postgres
```

Rules:

- The application layer owns the decision that a token, runtime session, or active connection must be invalidated.
- The authentication service may execute token revocation only after a later implementation slice authorizes it.
- The session repository may execute runtime session revocation only after a later implementation slice authorizes it.
- A future connection registry must expose a narrow application-owned invalidation target; it must not make WebSocket transport the owner of authentication state.
- WebSocket transport may close concrete sockets only through a narrow transport-owned capability requested by application policy.
- Protobuf adapters may map public logout or invalidation results only after a protocol route gate; they must not decide revocation policy.
- Domain modules must never parse raw tokens, session ids, connection ids, WebSocket close metadata, or transport credential carriers to decide logout or revocation behavior.

## 3. Future Policy Questions

A later implementation gate must answer these questions before code changes:

```yaml
future_policy_questions:
  logout_scope:
    - Does logout revoke only the presented access token?
    - Does logout also revoke the linked runtime session?
    - Does logout support all sessions for a player or credential family?
  revocation_source:
    - Is revocation user-requested logout, admin action, expiration cleanup, credential compromise, or account disablement?
  active_connection_targeting:
    - Can the server find active connections by token record id?
    - Can the server find active connections by runtime session id?
    - Can the server find active connections by player id?
  invalidation_effect:
    - Should the active connection be closed immediately?
    - Should the connection remain open but fail protected routes?
    - Should the server emit a system message before close?
  transport_boundary:
    - Which layer chooses WebSocket close reason and code?
    - Which layer performs the concrete close?
  route_policy_boundary:
    - Which route families are affected after revocation?
    - Does request-token validation alone catch the revocation before dispatch?
    - Do bound-connection routes require an active registry check?
  reconnect_epoch_boundary:
    - Does revocation advance an epoch?
    - Does revocation prevent resume?
  observability_boundary:
    - Which counters, audit events, and redacted logs are required later?
```

Rules:

- These questions must not be answered implicitly inside WebSocket handlers.
- Metadata-only `player_id`, `session_id`, `connection_id`, or `connection_epoch` must not be enough to target or authorize revocation.
- A connection may be targeted only from server-observed binding/session/token state, never from client-asserted metadata alone.
- Public errors must not reveal whether token lookup, token revocation, session revocation, active connection lookup, or connection close failed.

## 4. Recommended Future First Posture

The recommended first future posture is conservative:

```yaml
recommended_future_first_posture:
  logout_scope: presented_access_token_only
  runtime_session_revocation: linked_session_policy_deferred_until_implementation_gate
  active_connection_invalidation: policy_defined_before_implementation
  connection_registry_required_before_targeting_active_connections: true
  websocket_transport_auth_state_owner: false
  close_active_socket_on_logout_default: not_selected_by_this_gate
  bound_route_policy_reclassification: not_changed_by_this_gate
  reconnect_epoch_interaction: deferred
  public_logout_route_protocol: deferred
  direct_nakama_pitaya_api_compatibility: false
```

Rules:

- The first future logout implementation should revoke only the presented access token unless a later ADR explicitly broadens the scope.
- Runtime session revocation should be a distinct policy choice even when the session row links to the access-token record.
- Active connection invalidation must not be best-effort hidden behavior. A later implementation must state whether failure to invalidate active connections fails logout, succeeds with warning, or is retried asynchronously.
- A future connection registry is required before the server can target already-open sockets by player, session, token, connection id, or epoch.
- Ordinary protected routes remain request-token protected by default; this gate does not reclassify routes.

## 5. Future Event And State Vocabulary

Future implementation may use this vocabulary only after a later implementation gate:

```yaml
candidate_internal_state_transitions:
  token:
    active_to_revoked: authentication_access_tokens.status
  runtime_session:
    active_to_revoked: runtime_sessions.session_status
  active_connection:
    bound_to_invalidated: future_connection_registry_state
    bound_to_closed: future_transport_close_result

candidate_internal_reasons:
  - logout_presented_access_token
  - token_revoked_by_policy
  - runtime_session_revoked_by_policy
  - player_account_disabled
  - credential_compromised
  - admin_session_revocation
  - duplicate_connection_replacement
```

Rules:

- Token revocation, session revocation, and active-connection invalidation are separate state transitions.
- A later implementation must choose transaction boundaries explicitly. In particular, it must state whether token revocation and runtime session revocation occur in the same unit of work.
- Active connection close is outside SQL transaction control and must be modeled as an application/transport side effect with redacted outcome handling.
- Future events must not include raw token text, raw credential material, lookup digests, verifier digests, verifier key ids, Authorization headers, cookies, query strings, WebSocket subprotocol values, or inner payload bytes.

## 6. Error And Redaction Boundary

Future public behavior must be stable and redacted:

```yaml
candidate_public_errors:
  logout_token_missing: AUTHENTICATION_TOKEN_MISSING
  logout_token_malformed: AUTHENTICATION_TOKEN_MALFORMED
  logout_token_invalid: AUTHENTICATION_TOKEN_INVALID
  logout_unavailable: AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
  session_invalid: SESSION_INVALID
  active_connection_invalidation_unavailable: CONNECTION_INVALIDATION_UNAVAILABLE
```

Rules:

- Public logout and revocation failures must not reveal whether the token was unknown, already revoked, expired, linked to a missing session, linked to a closed connection, or blocked by repository failure.
- Active connection invalidation errors must not leak connection ids, session ids, token record ids, player ids, remote addresses, close codes, or transport-specific internals.
- Logs and events must use redacted ids or intentionally classified internal ids only after an observability standard defines what is log-safe.
- `verifier_key_id`, lookup digest, verifier digest, raw token text, and raw credential material remain not log-safe.
- WebSocket close reason text must be generic unless a future protocol/transport close policy standard authorizes a more specific client-visible value.

## 7. Relationship To Existing Runtime Pieces

This gate does not change existing behavior:

```yaml
authentication_service_logout_changed: false
authentication_service_refresh_changed: false
access_token_validation_changed: false
route_policy_changed: false
connection_binding_changed: false
runtime_session_validation_changed: false
session_creation_changed: false
session_repository_changed: false
postgres_adapter_changed: false
websocket_transport_changed: false
protobuf_protocol_changed: false
```

Rules:

- `LogoutAccessToken` may remain fail-closed until a future implementation gate.
- `RefreshAccessToken` may remain unsupported.
- Access-token validation may continue rejecting revoked tokens through existing record status checks when invoked per request.
- Bound connection identity does not become an active connection registry.
- Runtime session repository revocation methods remain storage-neutral capabilities until a future behavior slice uses them.
- Route policy continues to classify ordinary protected routes as request-token required by default.
- No session last-seen, cleanup, audit mutation, or connection replacement behavior is authorized here.

## 8. Relationship To WebSocket And Protocol

This gate does not authorize WebSocket or Protobuf changes:

```yaml
websocket_transport_credential_neutral: true
websocket_close_policy_added: false
transport_credential_carrier_added: false
websocket_handshake_authentication_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
existing_protobuf_envelope_change_added: false
generated_output_added: false
```

Rules:

- WebSocket transport must not parse Authorization headers, bearer values, cookies, query-string tokens, session tokens, or `Sec-WebSocket-Protocol` authentication material for this gate.
- No logout command Protobuf source, generated output, authentication response session fields, revocation system message, close reason enum, reconnect message, or envelope field is authorized here.
- A future protocol logout route gate must authorize any client-visible logout command before it is exposed.
- A future transport close policy gate must authorize any custom WebSocket close code or close reason.
- A future protocol session carrier gate must authorize any client-visible session id, session proof, resume token, or connection epoch carrier.

## 9. Deferrals

This gate does not authorize:

- Implementing `LogoutAccessToken`.
- Implementing token revocation execution.
- Implementing runtime session revocation execution.
- Implementing logout-all-sessions, account-wide revocation, credential-wide revocation, admin revocation, or credential compromise workflows.
- Closing WebSocket connections on logout or revocation.
- Adding an active connection registry.
- Wiring a connection registry into frame handling.
- Adding kick, disconnect, duplicate replacement, reconnect, resume, or epoch behavior.
- Adding WebSocket close codes, close reason policy, or close system messages.
- Adding Protobuf logout routes, session carriers, resume carriers, generated output, or existing envelope changes.
- Adding WebSocket handshake authentication or transport credential carriers.
- Adding cleanup jobs, async invalidation workers, queues, metrics, dashboards, admin APIs, or observability dependencies.
- Reclassifying ordinary protected routes away from request-level token proof.
- Adding memory durable session behavior.
- Adding direct Nakama or Pitaya public API compatibility.

## 10. Test Requirements For Future Implementation

A later implementation must include focused tests for:

- Missing, malformed, unknown, expired, already-revoked, and active tokens collapse to the selected public outcomes.
- Presented-token logout revokes only the selected scope.
- Runtime session revocation behavior is either explicitly deferred or implemented with tests showing the transaction boundary.
- Active connection invalidation is either explicitly not selected or implemented through a narrow registry/transport interface.
- Connection targeting never relies on metadata-only identity.
- Bound/session route policies continue failing closed after revocation.
- Failure to close an already-closed or missing connection follows the selected policy.
- Logout side effects do not happen before the token proof is validated.
- Public errors, logs, and events are redacted.
- WebSocket transport remains credential-neutral.
- Protobuf envelope and authentication response shapes remain unchanged unless a separate protocol gate authorizes changes.

Live PostgreSQL verification may remain opt-in unless the later implementation changes persistence behavior in a way that needs real database checks.

## 11. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

```yaml
adopted_concepts:
  - authenticated_session_material_has_logout_revocation_and_expiration_lifecycle_pressure
  - revoked_or_expired_credentials_must_not_authorize_gameplay_requests
  - realtime_socket_behavior_must_be_considered_when session_lifecycle_changes
adapted_concepts:
  - vibit_keeps_presented_token_logout_runtime_session_revocation_and_active_socket_invalidation_as_separate_policy_decisions
  - vibit_requires_connection_registry_policy_before_targeting_open_sockets
  - vibit_does_not_copy_nakama_session_api_jwt_shape_or_realtime_socket_contract
deferred_concepts:
  - refresh_token_flow
  - logout_all_sessions
  - session_management_api
  - admin_session_revocation_surface
  - dashboard_session_operations
rejected_for_now:
  - direct_nakama_api_compatibility
```

Pitaya reference mapping:

```yaml
adopted_concepts:
  - acceptors_sessions_and_handlers_are_separate_lifecycle_surfaces
  - session_binding_can_support_targeted_connection_management
  - handler_logic_should_receive_context_not_parse_transport_credentials
adapted_concepts:
  - vibit_future_connection_registry_is_application_runtime_owned_with_narrow_transport_close_handoff
  - vibit_keeps_socket_close_policy_out_of_authentication_repositories
  - vibit_defers_cluster_session_routing_until_single_process_boundaries_are_stable
deferred_concepts:
  - frontend_backend_cluster_session_invalidation
  - distributed_kick_or_disconnect_routing
  - groups_broadcast_integration
  - server_to_server_rpc_invalidation
rejected_for_now:
  - direct_pitaya_api_compatibility
```

Nakama and Pitaya guide capability planning. They do not override vibit's constitution, ADRs, manifests, generated boundaries, or verification commands.

## 12. Future Implementation Queue

After this gate, future work should remain split:

```yaml
future_work_items:
  logout_access_token_behavior_gate:
    may_define:
      - presented-token logout execution
      - token revocation transaction boundary
      - whether linked runtime session revocation is included
  active_connection_registry_gate:
    may_define:
      - server-observed connection registry ownership
      - lookup keys and redaction
      - narrow transport invalidation handoff
  logout_revocation_active_connection_implementation:
    may_add:
      - token or session revocation execution
      - active connection invalidation only if registry and close policy are already selected
  reconnect_connection_epoch_gate:
    may_define:
      - reconnect, resume, duplicate replacement, and epoch mismatch behavior
  protocol_logout_route_gate:
    may_define:
      - client-visible logout command route and Protobuf messages
  protocol_session_carrier_gate:
    may_define:
      - whether and how clients receive or carry session ids, resume tokens, or connection epochs
```

Do not combine these into one broad connection/session subsystem slice without a new ADR.

## 13. Verification

The repository check rule for this gate is:

```text
runtime.logout_revocation_active_connection_gate
```
