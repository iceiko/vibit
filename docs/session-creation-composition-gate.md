# Session Creation Composition Gate

Status: Draft v0.1
Last updated: 2026-05-17
Scope: Gate-only boundary for future durable runtime session creation composition after persistent runtime session validation exists
Depends on: `docs/runtime-session-validation-gate.md`, `decisions/ADR-0066-runtime-session-validation-implementation.md`, `docs/session-repository-boundary.md`, `docs/session-persistence-websocket-handshake-ratification.md`, `docs/authentication-command-protocol-login-route-gate.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0067`

The paired Simplified Chinese translation is `docs/session-creation-composition-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit now has the durable prerequisites for runtime sessions:

- A PostgreSQL `runtime_sessions` migration source.
- A storage-neutral `runtime/internal/app/session.Repository` interface.
- A PostgreSQL adapter for that repository.
- An application-owned persistent runtime session validator.
- Existing device-credential login and access-token validation behavior.

The missing boundary is creation composition. `AuthenticateWithDeviceCredential` can issue an opaque access token, and `PersistentSessionValidator` can validate a persisted session, but no production path creates a durable runtime session row. The next useful step is to define where future session creation belongs and how it composes with login and token issuance.

Mature game servers shape this boundary:

- Nakama treats authentication as a session lifecycle entry point with token material, expiration, refresh, logout, and management pressure.
- Nakama also separates session validity from active socket disconnect behavior.
- Pitaya keeps acceptors, handler routing, and session context separated, so durable session creation should be application composition rather than transport behavior.

vibit should adapt those lessons by making durable session creation explicit, transactional, and application-owned. This standard defines the gate only.

```yaml
session_creation_composition_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0146
decision: ADR-0067
check_rule: runtime.session_creation_composition_gate
future_composition_owner: runtime/internal/app
future_authentication_service_owner: runtime/internal/app/authentication
session_repository_owner: runtime/internal/app/session
session_repository_interface: runtime/internal/app/session.Repository
session_repository_create_method: CreateRuntimeSession
future_login_composition_candidate: AuthenticateWithDeviceCredential
future_session_id_generation_owner: runtime/internal/app
session_creation_behavior_added: false
authentication_service_code_changed: false
runtime_session_validation_changed: false
request_identity_session_validated_policy_changed: false
route_policy_session_identity_added: false
route_policy_bound_identity_added: false
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
protobuf_session_messages_added: false
existing_protobuf_envelope_change_added: false
generated_output_added: false
logout_revocation_active_connection_added: false
reconnect_epoch_behavior_added: false
cleanup_jobs_added: false
dependencies_added: false
memory_durable_session_behavior_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

This is a gate-only standard. It does not create runtime sessions in code.

## 2. Ownership

Future session creation composition is application-owned:

```yaml
future_composition_owner: runtime/internal/app
authentication_service_owner: runtime/internal/app/authentication
session_record_owner: runtime/internal/app/session
postgresql_session_adapter_owner: runtime/internal/platform/persistence/postgres
transport_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
```

Rules:

- Future durable session creation may call `session.Repository.CreateRuntimeSession` only through application unit-of-work capabilities.
- The session repository and PostgreSQL adapter remain storage-oriented. They must not decide when login should create a session.
- WebSocket transport must not create durable runtime session rows.
- Protobuf adapters must not create durable runtime session rows.
- Domain modules must not create authentication runtime sessions as a side effect of domain commands.
- Authentication may compose session creation only because the first selected login path owns successful token issuance, not because the authentication module owns runtime session persistence.

## 3. Future Composition Semantics

A later implementation slice may make successful device-credential login create a durable runtime session in the same unit of work as access-token storage.

Candidate future order:

1. Reject missing or malformed device credential proof before unit-of-work.
2. Validate the device credential and active player account.
3. Generate opaque access-token material.
4. Compute and store token lookup and verifier digests.
5. Generate a server-owned runtime `session_id`.
6. Create a `runtime_sessions` row linked to the stored `access_token_record_id`.
7. Commit the unit of work.
8. Return client-visible token material and any future authorized session material only after commit.

Rules:

- Session creation must happen after credential proof validation and before commit.
- Raw credential material and raw access-token material must not be stored in `runtime_sessions`.
- Token lookup digests, token verifier digests, credential lookup digests, credential verifier digests, and verifier key ids must not be copied into runtime session records.
- `access_token_record_id` linkage is private server metadata. It is not proof and must not be exposed as a client credential.
- Session creation must not make request routes session-validated by itself.
- Session creation must not set `RequestIdentity.SessionValidated = true`; validation remains a separate runtime-session validation concern.
- Session creation must not close, replace, or resume WebSocket connections.

## 4. Session ID And Lifetime Posture

Future session creation needs an explicit session id and lifetime posture:

```yaml
candidate_session_id_posture:
  generated_by: application_owned_secure_material_generator
  client_supplied_session_id_allowed: false
  stored_raw_session_id_allowed: true
  session_id_is_proof: false
candidate_lifetime_posture:
  issued_at_source: application_clock
  expires_at_source: selected_access_token_expiration_or_later_session_policy
  last_seen_at_initial_value: issued_at
  initial_status: active
```

Rules:

- A future session id generator must be authorized separately before code is added.
- The first session id may be opaque high-entropy text if a later implementation gate chooses that posture.
- A client-supplied session id must not be accepted during session creation.
- Session id uniqueness collisions must fail closed and remain redacted.
- Session id values are operationally sensitive unless a later observability gate defines a log-safe redaction format.
- The first session lifetime may align to the access-token lifetime, but that choice must be restated in the implementation gate.
- Refresh, renewal, extension, rotation, and token-session rekeying remain deferred.

## 5. Unit-Of-Work Composition Boundary

Future login-created sessions should be transactionally composed:

```yaml
future_unit_of_work_capabilities:
  - NewAuthenticationRepository
  - NewPlayerAccountRepository
  - NewSessionRepository
future_session_mutation:
  - session.CreateRuntimeSessionMutation
future_token_linkage:
  field: access_token_record_id
  role: private_server_metadata
```

Rules:

- Token storage and session creation should commit or roll back together.
- A token must not be returned to the client if session creation is required by the selected posture and session creation fails.
- If a later implementation keeps session creation optional, the public result must make that posture explicit and tests must cover it.
- Repository acquisition failures and session creation failures must collapse to redacted dependency or authentication-unavailable public errors.
- The future implementation must test commit failure behavior so raw token or session material is not returned as a successful result before commit.

## 6. Relationship To Validation And Route Policy

This gate does not change validation or route policy:

```yaml
runtime_session_validation_changed: false
request_identity_session_validated_policy_changed: false
route_policy_session_identity_added: false
route_policy_bound_identity_added: false
```

Rules:

- `PersistentSessionValidator` remains lookup-only and is not wired into route policy by this gate.
- Access-token validation remains the current protected-route proof path.
- Future session creation may produce a durable row that a later validation or route-policy slice can use, but this gate does not select that policy.
- First-message bound identity still does not satisfy ordinary protected route policy.
- Metadata-only player id, session id, and connection metadata are still not authenticated proof.

## 7. Relationship To WebSocket And Protocol

This gate does not change WebSocket or Protobuf behavior:

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
- No session creation command, response field, session carrier, system message, generated Protobuf output, or generated contract shape is authorized here.
- A future protocol gate must authorize any client-visible session id or session result shape before the login response exposes it.
- No reconnect, resume, duplicate replacement, durable connection epoch, logout disconnect, presence, rooms, parties, groups, or match attachment behavior is authorized here.

## 8. Error And Redaction Boundary

Future session creation failures should not leak sensitive internals:

```yaml
candidate_public_failure_classes:
  - AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE
  - AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
  - AUTHENTICATION_CREDENTIAL_INVALID
candidate_internal_failure_reasons:
  - session_id_generation_failed
  - session_repository_unavailable
  - session_record_conflict
  - session_record_constraint_failure
  - unit_of_work_commit_failed
```

Rules:

- Public errors must not reveal raw token material, raw credential material, session ids, token record ids, lookup digests, verifier digests, verifier key ids, SQL argument values, Authorization headers, cookies, or WebSocket credential carriers.
- Internal test-only failure classes may be more specific, but logged and public output must remain redacted.
- Session id generation failures and uniqueness conflicts must not echo candidate session ids.
- Session creation must not introduce new log-safe status for `session_id`, `token_record_id`, or `player_id`.

## 9. Test Requirements For Future Implementation

A later session creation implementation must include focused tests for:

- Successful login creates exactly one active runtime session row through `session.Repository.CreateRuntimeSession`.
- Session creation happens only after credential proof validation, player account validation, token generation, digest computation, and token storage.
- The runtime session links to `access_token_record_id` but stores no raw proof or digest material.
- Session id generation rejects missing, malformed, duplicate, or low-entropy values according to the later selected generator posture.
- Session creation failure prevents successful token/session material return.
- Unit-of-work commit failure prevents successful token/session material return.
- Repository acquisition errors and session creation errors remain redacted.
- Access-token validation behavior remains unchanged.
- `RequestIdentity.SessionValidated` remains false until a separate validation/route-policy path sets it.
- WebSocket transport and Protobuf envelope remain unchanged.

Live PostgreSQL verification may remain opt-in unless a later implementation work item requires it.

## 10. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

```yaml
adopted_concepts:
  - login_creates_session_lifecycle_pressure
  - sessions_have_expiration_logout_refresh_and_management_implications
  - token_or_session_response_material_is_client_visible_only_after_success
adapted_concepts:
  - vibit_keeps_opaque_access_token_proof_and_runtime_session_records_separate
  - runtime_session_creation_is_application_unit_of_work_composition
  - access_token_record_linkage_is_private_server_metadata
deferred_concepts:
  - refresh_token_session_extension
  - logout_disconnect_active_socket
  - session_management_api
  - single_session_or_single_socket_policy
rejected_for_now:
  - direct_nakama_session_api_compatibility
```

Pitaya reference mapping:

```yaml
adopted_concepts:
  - acceptor_transport_and_handler_logic_are_separate
  - session_context_can_be_handler_facing_without_transport_owning_persistence
  - backend_session_state_changes_need_explicit_composition
adapted_concepts:
  - durable_session_creation_belongs_to_application_composition
  - websocket_transport_remains_credential_neutral
  - route_handlers_receive_normalized_identity_not_storage_rows
deferred_concepts:
  - unique_session_enforcement
  - frontend_backend_cluster_session_routing
  - durable_connection_registry
rejected_for_now:
  - direct_pitaya_session_api_compatibility
```

Nakama and Pitaya guide capability planning. They do not override vibit's constitution, ADRs, manifests, generated boundaries, or verification commands.

## 11. Future Implementation Queue

After this gate, future work should remain split:

```yaml
future_work_items:
  session_creation_composition_implementation:
    requires_later_gate: true
    may_add:
      - application-owned session id generation
      - login-time session repository capability use
      - CreateRuntimeSession composition inside authentication service behavior
    must_not_add:
      - route-policy use of session-validated identity
      - WebSocket handshake authentication
      - Protobuf session carriers without protocol gate
      - logout-triggered active connection invalidation
      - reconnect or epoch behavior
  bound_identity_route_policy_gate:
    requires_later_gate: true
  logout_revocation_active_connection_gate:
    requires_later_gate: true
  reconnect_connection_epoch_gate:
    requires_later_gate: true
  operations_observability_and_admin_tooling:
    requires_later_gate: true
```

Do not combine these into one broad session subsystem slice without a new ADR.

## 12. Verification

Repository verification for this gate is:

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-session-creation-composition-gate --json
node tools/vibit check all --json
```

The repository check rule is:

```yaml
runtime.session_creation_composition_gate
```
