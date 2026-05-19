# Runtime Session Validation Gate

Status: Draft v0.1
Last updated: 2026-05-17
Scope: Gate-only boundary for future runtime session validation behavior after the PostgreSQL session adapter
Depends on: `docs/session-postgresql-adapter-gate.md`, `decisions/ADR-0064-session-postgresql-adapter-implementation.md`, `docs/session-repository-boundary.md`, `docs/session-persistence-websocket-handshake-ratification.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0065`

The paired Simplified Chinese translation is `docs/runtime-session-validation-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit now has the durable pieces needed before session validation can be designed:

- A PostgreSQL `runtime_sessions` migration source.
- A storage-neutral `runtime/internal/app/session.Repository` interface.
- A PostgreSQL adapter for that repository.
- Existing request-level access-token validation and first-message connection binding, both of which keep `SessionValidated` false in their current production paths.

The next useful step is not to silently enable session validation. The next step is to define the future validation gate. Mature game servers shape this boundary:

- Nakama treats authenticated sessions as lifecycle objects with expiration, refresh, logout, and socket-related implications.
- Nakama also separates token/session logout from active socket disconnect behavior.
- Pitaya exposes session context to handlers while keeping acceptors, session binding, and handler execution separated.

vibit should adapt those lessons by making runtime session validation application-owned and explicit. This standard defines the gate only.

```yaml
runtime_session_validation_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0142
decision: ADR-0065
check_rule: runtime.runtime_session_validation_gate
future_validation_owner: runtime/internal/app
future_session_repository_owner: runtime/internal/app/session
postgresql_session_adapter_owner: runtime/internal/platform/persistence/postgres
future_validator_source_candidate: runtime/internal/app/runtime_session_validator.go
future_validator_test_candidate: runtime/internal/app/runtime_session_validator_test.go
runtime_session_validation_added: false
request_identity_session_validated_true_added: false
session_creation_composition_added: false
route_policy_session_identity_added: false
route_policy_bound_identity_added: false
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

This is a gate-only standard. It does not add runtime validation code.

## 2. Ownership

Future runtime session validation is application-owned:

```yaml
future_validation_owner: runtime/internal/app
session_record_owner: runtime/internal/app/session
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
authentication_service_owner: runtime/internal/app/authentication
transport_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
```

Rules:

- The future validator may depend on an application-facing session repository capability, not on PostgreSQL details.
- The future validator may consume an already-normalized request identity, access-token validation result, or connection binding result only through application-owned types.
- The future validator must not import WebSocket transport packages, generated Protobuf packages, PostgreSQL adapters, pgx, SQL rows, or migration packages.
- The PostgreSQL session adapter must remain persistence-only and must not create `RequestIdentity`.
- WebSocket transport must remain credential-neutral for this gate.
- Protocol adapters must not decide whether a session is valid. They may later carry proof or metadata only after a separate protocol gate.

## 3. Future Validation Semantics

A later implementation slice may define a validator with these candidate inputs:

```yaml
candidate_validation_input:
  - route_request_identity
  - route_request_session_metadata
  - server_observed_connection_id
  - server_observed_connection_epoch
  - observed_at
  - session_repository
```

Candidate future validation order:

1. Reject missing or malformed required validation input before repository lookup.
2. Require an already-validated actor identity before trusting a persisted session row.
3. Look up the runtime session by server-owned `session_id` through `session.Repository`.
4. Require `session_status = active`.
5. Require `expires_at > observed_at`.
6. Require the session actor to match the already-validated actor identity.
7. Require player account identity handoff to remain normalized and application-owned.
8. Optionally update `last_seen_at` only if a later implementation gate authorizes that mutation.
9. Return `RequestIdentity.SessionValidated = true` only after all selected checks pass.
10. Collapse public failure output to a stable invalid-session error class.

Rules:

- A persisted `session_id` alone is not proof.
- `access_token_record_id` linkage is private metadata. It must not replace access-token proof validation unless a later ADR explicitly defines that composition.
- A client-supplied envelope `Session.session_id` remains metadata until future validation checks it against durable state and already-validated actor identity.
- Session validation must not implicitly create sessions.
- Session validation must not refresh access tokens, extend session expiration, revoke tokens, or logout users.
- Session validation must not close WebSocket connections.

## 4. Request Identity Handoff

Future validation may set:

```yaml
future_request_identity:
  status: validated
  actor_kind: player
  player_id_validated: true
  session_validated: true
```

Only a later implementation slice may set `SessionValidated` true in production behavior.

Rules:

- `SessionValidated` true requires both an already-validated actor identity and a valid persisted runtime session row.
- Metadata-only identity must never become session-validated.
- First-message bound identity must not satisfy ordinary protected route policy through this gate.
- Route policy may use session-validated identity only after a separate route-policy gate defines which routes require it and how failures map to public errors.
- Domain modules must continue to receive normalized `RequestIdentity`, not session repository records.

## 5. Error And Redaction Boundary

Future public session validation errors should collapse to:

```yaml
public_error_class:
  - SESSION_INVALID
```

Internal failure reasons may be more specific for tests and private control flow:

```yaml
candidate_internal_failure_reasons:
  - missing_session_id
  - malformed_session_id
  - actor_identity_not_validated
  - session_not_found
  - session_inactive
  - session_expired
  - session_actor_mismatch
  - session_player_mismatch
  - session_repository_unavailable
```

Rules:

- Public errors must not reveal whether lookup, expiration, revocation, actor mismatch, player mismatch, token linkage, or repository failure caused validation failure.
- Errors, logs, events, and tests must not include raw access-token text, raw credential material, lookup digests, verifier digests, verifier key ids, Authorization headers, cookies, query-string tokens, or WebSocket subprotocol token material.
- Session ids and player ids should be treated as operationally sensitive unless a later observability gate classifies specific redacted forms as log-safe.

## 6. Relationship To Authentication

Runtime session validation does not replace token proof validation.

```yaml
access_token_validation_owner: runtime/internal/app/authentication
session_validation_owner: runtime/internal/app
session_validation_replaces_token_validation: false
session_validator_reads_token_digests: false
session_validator_reads_raw_tokens: false
```

Rules:

- Access-token validation remains the proof-validation path for the current protected routes.
- The future session validator may consume validated actor identity produced by authentication, but it must not compute token digests, compare token verifiers, issue tokens, refresh tokens, or revoke tokens.
- Session validation and access-token validation composition must be defined explicitly before route policy depends on both.
- Logout and revocation active-connection behavior remain behind a separate gate.

## 7. Relationship To WebSocket And Protocol

This gate does not change WebSocket or Protobuf behavior:

```yaml
websocket_transport_credential_neutral: true
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
protobuf_session_messages_added: false
existing_protobuf_envelope_change_added: false
connection_registry_added: false
```

Rules:

- WebSocket transport must not parse Authorization headers, bearer values, cookies, query-string tokens, session tokens, or `Sec-WebSocket-Protocol` authentication material for this gate.
- The existing Protobuf envelope remains unchanged.
- No session validation protocol messages or generated output are authorized here.
- No reconnect, resume, duplicate replacement, durable connection epoch, logout disconnect, presence, rooms, parties, groups, or match attachment behavior is authorized here.

## 8. Test Requirements For Future Implementation

A later runtime session validation implementation must include focused tests for:

- Missing or malformed session id rejection before repository lookup.
- Metadata-only identity rejection.
- Already-validated actor identity requirement.
- Active session lookup through `session.Repository`.
- Expired, revoked, not-found, and actor-mismatch collapse to the same public invalid-session class.
- Successful validation sets `SessionValidated` true only after all checks pass.
- Repository errors remain redacted and do not leak raw proof or digest material.
- No WebSocket transport credential parsing.
- No Protobuf envelope shape change.
- No route-policy use of session-validated identity unless a later route-policy gate authorizes it.

Live PostgreSQL verification may remain opt-in unless a later implementation work item requires it.

## 9. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

```yaml
adopted_concepts:
  - sessions_have_lifetime_and_invalid_state
  - logout_and_revocation_affect_session_validity
  - realtime_socket_lifecycle_is_related_but_not_identical_to_session_validity
adapted_concepts:
  - vibit_uses_opaque_access_token_proof_plus_durable_runtime_session_records
  - public_invalid_session_failures_are_collapsed
  - session_validation_is_application_owned
deferred_concepts:
  - refresh_token_session_extension
  - logout_disconnect_active_socket
  - admin_session_management_api
  - single_socket_or_single_session_policy
rejected_for_now:
  - direct_nakama_session_api_compatibility
```

Pitaya reference mapping:

```yaml
adopted_concepts:
  - session_context_is_available_to_handlers_through_context
  - acceptor_transport_and_handler_logic_are_separate
  - session_lifecycle_callbacks_are_distinct_from_request_handler_logic
adapted_concepts:
  - durable_session_validation_builds_application_request_identity
  - transport_acceptor_remains_credential_neutral
  - route_handlers_receive_normalized_identity_not_storage_rows
deferred_concepts:
  - unique_session_enforcement
  - frontend_backend_cluster_session_routing
  - durable_connection_registry
  - group_or_room_session_membership
rejected_for_now:
  - direct_pitaya_session_api_compatibility
```

Nakama and Pitaya guide capability planning. They do not override vibit's constitution, ADRs, manifests, generated boundaries, or verification commands.

## 10. Future Implementation Queue

After this gate, future work should remain split:

```yaml
future_work_items:
  runtime_session_validation_implementation:
    may_add:
      - application-owned runtime session validator
      - tests with fake session repository
      - redacted error mapping
  session_creation_composition_gate:
    may_define:
      - whether login or BindConnection creates durable sessions
      - session id generation
      - token record linkage
      - expiration and last_seen semantics
  bound_identity_route_policy_gate:
    may_define:
      - which routes can use bound or session-validated identity
      - failure behavior
  logout_revocation_active_connection_gate:
    may_define:
      - whether revocation closes active WebSocket connections
  reconnect_connection_epoch_gate:
    may_define:
      - duplicate replacement
      - reconnect and resume behavior
```

Do not combine these into one broad session subsystem slice without a new ADR.

## 11. Verification

Repository verification for this boundary is:

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-runtime-session-validation-gate --json
node tools/vibit check all --json
```

The repository check rule is:

```yaml
runtime.runtime_session_validation_gate
```
