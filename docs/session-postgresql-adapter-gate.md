# Session PostgreSQL Adapter Gate

Status: Draft v0.1
Last updated: 2026-05-17
Scope: Gate-only boundary for the future PostgreSQL adapter that implements the storage-neutral runtime session repository interface
Depends on: `docs/session-repository-boundary.md`, `decisions/ADR-0062-session-repository-interface-implementation.md`, `runtime/migrations/postgres/000005_create_runtime_sessions.sql`, `docs/postgresql-persistence-boundary.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0063`

The paired Simplified Chinese translation is `docs/session-postgresql-adapter-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit now has:

- A PostgreSQL `runtime_sessions` migration source.
- A storage-neutral `runtime/internal/app/session.Repository` interface.
- No PostgreSQL adapter for that interface.
- No runtime session creation, lookup, validation, revocation execution, cleanup job, or route-policy use of persisted sessions.

The next useful step is a gate for the future PostgreSQL adapter. Mature game servers shape this boundary:

- Nakama treats sessions as first-class lifecycle records that need durable lookup, expiration, revocation, logout, and management-ready listing.
- Pitaya keeps session/context handoff separate from acceptors and routing, which means transport should not own durable session persistence.

vibit should adapt those lessons by making the future adapter persistence-only and transaction-bound. This standard defines the adapter gate only.

```yaml
session_postgresql_adapter_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0138
decision: ADR-0063
check_rule: runtime.session_postgresql_adapter_gate
repository_interface: runtime/internal/app/session.Repository
repository_owner: runtime/internal/app/session
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
future_adapter_source_candidate: runtime/internal/platform/persistence/postgres/session_repository.go
future_adapter_test_candidate: runtime/internal/platform/persistence/postgres/session_repository_test.go
session_logical_table: runtime_sessions
unit_of_work_factory_candidate: NewSessionRepository
postgresql_adapter_added: false
unit_of_work_factory_added: false
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

This is a gate-only standard. It does not add Go adapter code or runtime behavior.

## 2. Ownership

The future adapter is platform persistence-owned:

```yaml
repository_interface_owner: runtime/internal/app/session
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
transaction_boundary_owner: runtime/internal/platform/tx
application_runtime_behavior_owner: runtime/internal/app
authentication_module_owns_session_adapter: false
websocket_transport_owns_session_adapter: false
protocol_adapter_owns_session_adapter: false
```

Rules:

- The adapter may later implement `runtime/internal/app/session.Repository`.
- The adapter may later import `runtime/internal/app/session` and `github.com/jackc/pgx/v5` only inside `runtime/internal/platform/persistence/postgres`.
- The storage-neutral `runtime/internal/app/session` package must not import PostgreSQL, pgx, SQL rows, WebSocket transport, or generated Protobuf packages.
- The adapter must not become the owner of runtime session validation policy. It maps repository calls to SQL and maps SQL rows back to application records.
- The adapter must not parse access-token proof, compare verifier digests, create `RequestIdentity`, or set `RequestIdentity.SessionValidated`.

## 3. Future Adapter Surface

A later implementation slice may add:

```yaml
future_adapter_files:
  - runtime/internal/platform/persistence/postgres/session_repository.go
  - runtime/internal/platform/persistence/postgres/session_repository_test.go
future_constructor:
  - NewSessionRepositoryForUnitOfWork(executor Executor)
future_unit_of_work_factory:
  - NewSessionRepository() (session.Repository, error)
```

The adapter should implement only the existing repository methods:

```yaml
repository_methods:
  - CreateRuntimeSession
  - GetRuntimeSession
  - FindActiveSessionByID
  - UpdateRuntimeSessionLastSeen
  - MarkRuntimeSessionExpired
  - RevokeRuntimeSession
  - ListActiveSessionsForPlayer
```

Rules:

- The adapter implementation must stay behind the already-ratified `Executor` and unit-of-work patterns.
- Repository methods must not open, commit, or roll back transactions directly.
- The future unit-of-work factory may expose `NewSessionRepository` only after the adapter implementation slice authorizes it.
- Query and mutation inputs must be normalized through `runtime/internal/app/session` helpers before SQL execution.
- Returned rows must be normalized through `session.NormalizeRuntimeSessionRecord`.
- Adapter errors must be typed and redacted. They may include action labels and constraint labels, but must not include raw token material, digest bytes, verifier key ids, SQL argument values, or access-token proof.

## 4. SQL Shape Boundary

The future adapter may use the existing `runtime_sessions` table only:

```yaml
allowed_table:
  - runtime_sessions
allowed_reference_tables:
  - player_accounts
  - authentication_access_tokens
forbidden_tables:
  - runtime_session_connections
forbidden_material:
  - raw_access_token
  - raw_credential
  - token_lookup_digest
  - token_verifier_digest
  - credential_lookup_digest
  - credential_verifier_digest
  - verifier_key_id
  - websocket_connection_state
  - websocket_connection_registry_rows
```

Future SQL shape rules:

- `CreateRuntimeSession` may insert one `runtime_sessions` row and return the inserted row.
- `GetRuntimeSession` may select a row by `session_id` without status or expiration filtering.
- `FindActiveSessionByID` may filter `session_status = 'active'` and `expires_at > observed_at`; it must not validate token proof.
- `UpdateRuntimeSessionLastSeen` may update `last_seen_at` and `updated_at` for one row.
- `MarkRuntimeSessionExpired` may transition an active session to `expired` for one row.
- `RevokeRuntimeSession` may transition a session to `revoked` with `revoked_at` and `revocation_reason` for one row.
- `ListActiveSessionsForPlayer` may list active, unexpired sessions for a player with an explicit bounded limit.

The adapter must not add SQL for cleanup jobs, connection registries, logout-triggered socket invalidation, reconnect state, presence, rooms, parties, match runtime, or social graph behavior.

## 5. Error Mapping

Future adapter errors should follow the existing PostgreSQL repository style:

```yaml
future_error_classes:
  - ErrSessionRecordNotFound
  - ErrSessionRecordConflict
  - ErrSessionRecordConstraint
  - ErrSessionRecordStale
```

Rules:

- `pgx.ErrNoRows` maps to a not-found class.
- Duplicate key conflicts map to a conflict class.
- not-null, foreign-key, and check constraint failures map to a constraint class.
- Zero affected rows on lifecycle updates map to not found or stale state, as the implementation gate defines.
- Public runtime validation failures must not leak whether session id lookup, expiration, revocation, actor mismatch, token mismatch, or player state failed. That public collapse belongs to a later runtime validation gate.

## 6. Relationship To Authentication And Request Identity

This gate does not change authentication behavior:

```yaml
token_validation_owner: runtime/internal/app/authentication
session_adapter_validates_token_proof: false
session_adapter_sets_request_identity: false
session_validated_status_added: false
```

Rules:

- Authentication token validation remains owned by the application authentication service and route access-token validator.
- The adapter may later store or return `access_token_record_id` linkage as opaque metadata only.
- The adapter must not read token lookup digests, token verifier digests, credential digests, raw token text, or raw credential material.
- `RequestIdentity.SessionValidated` remains false until a later runtime session validation gate defines behavior.
- Bound connection identity still does not satisfy ordinary protected route policy through this gate.

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

- WebSocket transport must not parse access tokens, cookies, query-string tokens, Authorization headers, or session tokens for this gate.
- The existing Protobuf envelope remains unchanged.
- No session protocol messages or generated output are authorized here.
- No logout-triggered WebSocket close, duplicate connection replacement, reconnect, resume, durable connection epoch, presence, rooms, parties, or match attachment behavior is authorized here.

## 8. Test Requirements For Future Implementation

A later adapter implementation must include focused tests for:

- Interface conformance to `session.Repository`.
- Constructor/executor nil rejection.
- SQL shape for each repository method.
- Argument normalization and UTC timestamp handoff.
- Row scanning into `session.RuntimeSession`, including nullable revocation and token linkage fields.
- `pgx.ErrNoRows`, conflict, foreign-key, check, and not-null error mapping.
- Zero affected row handling for update/revoke/expire paths.
- Bounded listing limit behavior.
- Transaction neutrality: no direct `BEGIN`, `COMMIT`, or `ROLLBACK` SQL inside the repository.
- Redaction: no raw token, raw credential, digest bytes, verifier keys, Authorization headers, cookies, or WebSocket credential carriers in errors or tests.

Live PostgreSQL verification may remain opt-in unless a later work item makes it mandatory under the disposable PostgreSQL verification environment.

## 9. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

```yaml
adopted_concepts:
  - sessions_are_first_class_durable_lifecycle_records
  - session_lookup_expiration_revocation_logout_and_listing_need_adapter_ready_storage
  - operational_session_management_pressure_should_shape_repository_queries
adapted_concepts:
  - adapter_is_vibit_postgresql_persistence_boundary
  - opaque_access_token_and_session_id_remain_separate
  - token_record_linkage_is_private_metadata
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
  - session_context_is_separate_from_acceptor
  - handlers_should_receive_application_context_not_transport_storage
  - transport_and_routing_must_not_own_durable_session_persistence
adapted_concepts:
  - postgresql_session_adapter_is_platform_persistence_only
  - request_identity_handoff_remains_application_owned
deferred_concepts:
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
  session_postgresql_adapter_implementation:
    may_add:
      - runtime/internal/platform/persistence/postgres/session_repository.go
      - runtime/internal/platform/persistence/postgres/session_repository_test.go
      - unit-of-work factory wiring for NewSessionRepository
    must_not_add:
      - runtime session creation at login or BindConnection
      - runtime session validation
      - RequestIdentity.SessionValidated true
      - WebSocket credential carriers
      - Protobuf session messages
      - logout-triggered active connection invalidation
      - reconnect or epoch behavior
  runtime_session_validation_gate:
    requires_later_gate: true
  logout_revocation_active_connection_gate:
    requires_later_gate: true
  reconnect_connection_epoch_gate:
    requires_later_gate: true
  bound_identity_route_policy_gate:
    requires_later_gate: true
```

Do not combine these into one broad session subsystem slice without a new ADR.

## 11. Verification

Repository verification for this gate is:

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-session-postgresql-adapter-gate --json
node tools/vibit check all --json
```

The repository check rule is:

```yaml
runtime.session_postgresql_adapter_gate
```
