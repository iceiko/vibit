# Session Repository Boundary

Status: Draft v0.1
Last updated: 2026-05-17
Scope: Gate-only boundary for the future storage-neutral runtime session repository after the PostgreSQL `runtime_sessions` migration source
Depends on: `docs/session-persistence-websocket-handshake-ratification.md`, `docs/postgres-session-persistence-schema-gate.md`, `decisions/ADR-0060-runtime-sessions-migration-source.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0061`

The paired Simplified Chinese translation is `docs/session-repository-boundary.zh-CN.md`. The English file is authoritative.

## 1. Purpose

The repository now has a concrete PostgreSQL `runtime_sessions` migration source, but it still does not have an application session repository, PostgreSQL session adapter, runtime session creation, or runtime session validation.

The next useful boundary is the repository seam for future session lifecycle behavior. Mature game servers put pressure on this design:

- Nakama treats authenticated sessions as first-class lifecycle objects with expiration, refresh, logout, and operational management concerns.
- Pitaya makes session-like context available to handlers while keeping acceptors, routing, and handler execution separated.

vibit should adapt those lessons as a storage-neutral, application-owned repository boundary. This standard defines that future boundary only.

```yaml
session_repository_boundary: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0134
decision: ADR-0061
check_rule: runtime.session_repository_boundary
future_repository_owner_candidate: runtime/internal/app/session
future_repository_interface_candidate: runtime/internal/app/session.Repository
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
session_logical_table: runtime_sessions
repository_interface_added: false
postgresql_adapter_added: false
runtime_session_behavior_added: false
runtime_session_validation_added: false
runtime_session_creation_added: false
runtime_session_revocation_added: false
runtime_session_cleanup_added: false
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
protobuf_session_messages_added: false
existing_protobuf_envelope_change_added: false
route_policy_session_identity_added: false
logout_revocation_active_connection_added: false
reconnect_epoch_behavior_added: false
dependencies_added: false
memory_durable_session_behavior_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

This is a gate-only standard. It does not add Go session repository code, SQL adapter code, runtime session behavior, or protocol behavior.

## 2. Ownership

The future repository is runtime application-owned:

```yaml
future_repository_owner_candidate: runtime/internal/app/session
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
runtime_session_table_owner: runtime.session
authentication_module_owns_runtime_sessions: false
player_module_owns_runtime_sessions: false
websocket_transport_owns_runtime_sessions: false
protocol_adapter_owns_runtime_sessions: false
```

Rules:

- The repository interface must be storage-neutral and application-facing.
- The interface must not mention PostgreSQL, pgx, SQL rows, goose migrations, prepared statements, or database transaction implementation details.
- The PostgreSQL adapter may later implement the interface under `runtime/internal/platform/persistence/postgres`, but only after a separate adapter gate.
- Authentication may link to session records through `authentication_access_tokens(token_record_id)` only as a record reference; it does not own session lifecycle storage.
- Player account storage owns player lifecycle state, not runtime sessions.
- WebSocket transport owns connection plumbing, not durable session persistence.
- Protocol adapters own wire conversion, not session repository behavior.

## 3. Future Repository Capability Vocabulary

A later implementation gate may choose a smaller or renamed API, but the first candidate capability family is:

```yaml
candidate_repository_capabilities:
  - CreateRuntimeSession
  - GetRuntimeSession
  - FindActiveSessionByID
  - UpdateRuntimeSessionLastSeen
  - MarkRuntimeSessionExpired
  - RevokeRuntimeSession
  - ListActiveSessionsForPlayer
```

Capability rules:

- `CreateRuntimeSession` may create a lifecycle row only after a later runtime behavior gate defines the caller, transaction boundary, id generation, expiration, and token linkage.
- `GetRuntimeSession` is a storage lookup. It must not validate proof or create request identity.
- `FindActiveSessionByID` may later express active-status filtering, but it must not collapse token validation and session validation into one hidden operation.
- `UpdateRuntimeSessionLastSeen` is lifecycle metadata mutation only. It must not refresh access tokens, extend sessions, or imply presence.
- `MarkRuntimeSessionExpired` and `RevokeRuntimeSession` are durable state transitions only. They must not close active WebSocket connections without a later active-connection gate.
- `ListActiveSessionsForPlayer` is allowed only if a later gate defines permission, pagination, leakage boundaries, and administrative versus player self-inspection semantics.

The repository must return typed, application-owned records and errors. It must not return raw SQL rows or database driver errors directly to domain modules or protocol adapters.

## 4. Data Boundary

The repository may later use the existing `runtime_sessions` table:

```yaml
session_logical_table: runtime_sessions
allowed_linkage_table: authentication_access_tokens
allowed_player_table: player_accounts
forbidden_material:
  - raw_access_token
  - raw_credential
  - token_lookup_digest
  - token_verifier_digest
  - credential_lookup_digest
  - credential_verifier_digest
  - websocket_connection_state
  - websocket_connection_registry_rows
```

Rules:

- The repository must not store or return raw access-token text, raw credential material, lookup digests, verifier digests, verifier keys, or HMAC inputs.
- `access_token_record_id` is optional linkage metadata, not proof.
- `session_id` is a server-owned runtime session identifier, not an access token.
- Session status must remain a closed vocabulary such as active, revoked, or expired.
- Expiration and revocation are session lifecycle states. They are not proof-verifier states.
- Runtime session records must not become WebSocket connection registry rows.

## 5. Relationship To Authentication And Request Identity

The future session repository does not own token validation.

```yaml
token_validation_owner: runtime/internal/app/authentication
session_repository_token_validation_owner: false
request_identity_owner: runtime/internal/app
session_repository_sets_request_identity: false
session_validated_status_added: false
```

Rules:

- Access-token validation remains owned by the application authentication service and route access-token validator.
- The repository may later be used after token validation, but it must not parse token proof or compare token verifier digests.
- `RequestIdentity.SessionValidated` remains false until a later runtime session validation gate defines exactly how a persisted session is validated and bound to request identity.
- Bound connection identity from `runtime.authentication.BindConnection` does not satisfy ordinary protected route policy through this repository boundary.
- Domain modules must continue to consume normalized `RequestIdentity`, not session repository records.

## 6. Relationship To WebSocket And Protocol

This boundary does not change WebSocket or Protobuf behavior:

```yaml
websocket_transport_credential_neutral: true
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
protobuf_session_messages_added: false
existing_protobuf_envelope_change_added: false
connection_registry_added: false
```

Rules:

- The WebSocket transport must not parse access tokens, cookies, query-string tokens, Authorization headers, or session tokens for this boundary.
- The existing Protobuf envelope remains unchanged.
- No session protocol messages or generated output are authorized here.
- No logout-triggered WebSocket close, duplicate connection replacement, reconnect, resume, or durable connection epoch behavior is authorized here.

## 7. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

```yaml
adopted_concepts:
  - sessions_are_first_class_lifecycle_records
  - session_expiration_revocation_and_logout_need_repository_ready_boundaries
  - active_session_lookup_is_needed_before_richer_realtime_features
adapted_concepts:
  - repository_is_vibit_storage_neutral_application_boundary
  - access_token_record_linkage_is_optional_and_not_public_session_api
deferred_concepts:
  - refresh_token_session_extension
  - logout_disconnect_active_socket
  - single_socket_or_single_session_policy
  - admin_session_management_api
rejected_for_now:
  - direct_nakama_session_api_compatibility
```

Pitaya reference mapping:

```yaml
adopted_concepts:
  - session_context_is_separate_from_acceptor
  - handlers_should_receive_application_context_not_transport_storage
  - routing_and_session_context_need_clean_boundaries
adapted_concepts:
  - durable_session_repository_is_not_transport_owned
  - handler_context_handoff_remains_request_identity_owned
deferred_concepts:
  - frontend_backend_cluster_session_routing
  - durable_connection_registry
  - group_or_room_session_membership
rejected_for_now:
  - direct_pitaya_session_api_compatibility
```

Nakama and Pitaya guide capability planning. They do not override vibit's constitution, ADRs, manifests, generated boundaries, or verification commands.

## 8. Future Implementation Queue

After this boundary, future work should remain split:

```yaml
future_work_items:
  session_postgresql_adapter_gate:
    may_define:
      - adapter ownership
      - transaction handoff
      - SQL query shape
      - adapter tests
  session_repository_interface_implementation:
    may_add:
      - runtime/internal/app/session
      - storage-neutral repository types
      - focused unit tests
  runtime_session_validation_gate:
    may_define:
      - when persisted sessions are looked up
      - how session status maps to request identity
      - when SessionValidated may become true
  logout_revocation_active_connection_gate:
    may_define:
      - whether session revocation closes active WebSocket connections
  reconnect_connection_epoch_gate:
    may_define:
      - duplicate replacement
      - reconnect and resume behavior
```

Do not combine these into one broad session subsystem slice without a new ADR.

## 9. Verification

Repository verification for this boundary is:

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-session-repository-boundary --json
node tools/vibit check all --json
```

The repository check rule is:

```yaml
runtime.session_repository_boundary
```

The check must verify that this standard, translation, ADR, conversation log, change specs, manifests, AGENTS guides, and rule registry are present, while Go session repository code, PostgreSQL session adapter code, runtime behavior, protocol changes, WebSocket credential carriers, new migrations, dependencies, and direct Nakama/Pitaya API compatibility remain absent.

## 10. Agent Guidance

When adding future session work:

- Read this standard, `ADR-0061`, `ADR-0060`, `docs/postgres-session-persistence-schema-gate.md`, and `docs/session-persistence-websocket-handshake-ratification.md`.
- Keep session repository behavior application-owned and storage-neutral.
- Keep PostgreSQL adapter behavior behind its own gate.
- Keep token validation owned by authentication service behavior.
- Keep `RequestIdentity.SessionValidated` false until runtime session validation is separately ratified.
- Keep WebSocket transport credential-neutral.
- Do not copy Nakama or Pitaya public APIs.
