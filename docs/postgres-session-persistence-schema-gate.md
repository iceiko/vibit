# PostgreSQL Session Persistence Schema Gate

Status: Draft v0.1
Last updated: 2026-05-17
Scope: Gate-only boundary for the first durable runtime session persistence schema after first-message connection binding
Depends on: `docs/session-persistence-websocket-handshake-ratification.md`, `docs/session-persistence-websocket-handshake-decision-gates.md`, `docs/first-message-connection-binding-implementation-gate.md`, `docs/postgresql-persistence-boundary.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0059`

The paired Simplified Chinese translation is `docs/postgres-session-persistence-schema-gate.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit now has the minimum authentication and realtime-connection foundation:

- Public device credential login route.
- Opaque access-token validation.
- Request-level protected route validation through `vibit.authentication.v1.AuthenticatedRequest`.
- PostgreSQL startup composition for authentication.
- First-message `runtime.authentication.BindConnection` system route.

The next durable boundary is session persistence. Mature game servers such as Nakama treat authenticated sessions as a first-class lifecycle concept before richer realtime features. Pitaya exposes useful vocabulary around session-like connection context and binding, but keeps acceptor and handler concerns separated. vibit should adapt those lessons without copying public APIs or moving authentication into transport.

This standard defines a gate for the future PostgreSQL session schema.

```yaml
postgres_session_persistence_schema_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0130
decision: ADR-0059
check_rule: runtime.postgres_session_persistence_schema_gate
future_session_logical_table: runtime_sessions
future_migration_source_candidate: runtime/migrations/postgres/000005_create_runtime_sessions.sql
future_repository_owner_candidate: runtime/internal/app/session
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
session_persistence_implementation_added: false
session_table_added: false
migration_source_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
websocket_handshake_authentication_added: false
route_policy_bound_identity_added: false
logout_revocation_active_connection_added: false
reconnect_epoch_behavior_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

This is a schema gate. It does not add SQL migration source, Go repository interfaces, PostgreSQL adapters, runtime validation behavior, or connection lifecycle behavior.

## 2. Selected First Durable Session Target

The first durable session persistence target is PostgreSQL:

```yaml
selected_first_session_store: postgres
future_session_logical_table: runtime_sessions
future_migration_source_candidate: runtime/migrations/postgres/000005_create_runtime_sessions.sql
future_session_repository_boundary: separate_future_work_item
```

Rationale:

- PostgreSQL is already vibit's first accepted authoritative durable store.
- Session persistence needs transactionally clear relationships with player account state and token verifier state.
- A SQL-first schema gate is easier for agents to inspect than ad hoc in-memory state.
- Deferring Redis-like stores avoids a dependency and operations decision before it is needed.

## 3. Future `runtime_sessions` Table Candidate

The future first session migration may define one logical table:

```yaml
runtime_sessions:
  primary_key:
    - session_id
  required_columns:
    - session_id
    - actor_kind
    - actor_id
    - player_id
    - session_status
    - issued_at
    - expires_at
    - last_seen_at
    - created_at
    - updated_at
  nullable_columns:
    - revoked_at
    - revocation_reason
    - access_token_record_id
  forbidden_columns:
    - raw_access_token
    - raw_credential
    - token_verifier_digest
    - token_lookup_digest
    - credential_verifier_digest
    - credential_lookup_digest
```

Rules:

- `session_id` is a server-generated opaque identifier. It is not an access token.
- `actor_kind` and `actor_id` use the ratified authentication/session vocabulary.
- `player_id` is required for the first player session posture.
- `session_status` must be a closed vocabulary such as `active`, `revoked`, or `expired`.
- `issued_at`, `expires_at`, `last_seen_at`, `created_at`, and `updated_at` must use UTC timestamps.
- `revoked_at` and `revocation_reason` support future logout/revocation work but do not implement active connection invalidation.
- `access_token_record_id` may be a future opaque linkage to token verifier storage, but raw token text and digest bytes must not be copied into the session table.

## 4. Explicit Deferral Of Connection Records

The first session schema gate does not choose a durable connection table.

```yaml
runtime_session_connections_table_selected_now: false
connection_binding_registry_added: false
active_connection_invalidation_added: false
reconnect_epoch_behavior_added: false
```

A future gate must define connection lifecycle storage before adding any of:

- `runtime_session_connections` table.
- Durable `connection_id` to `session_id` binding rows.
- Duplicate connection replacement.
- Reconnect/resume state.
- Connection epoch policy beyond metadata handoff.
- Logout-triggered active WebSocket close or invalidation.
- Presence, rooms, groups, parties, or match attachment.

This keeps the first schema small and prevents session persistence from silently becoming a full realtime connection manager.

## 5. Ownership

The future session persistence boundary is runtime-owned:

```yaml
runtime_session_owner: runtime/internal/app
future_repository_owner_candidate: runtime/internal/app/session
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
player_module_owns_session_table: false
authentication_module_owns_session_table: false
websocket_transport_owns_session_table: false
```

Rules:

- The player module owns player account lifecycle, not runtime sessions.
- The authentication module owns credential and token verifier record boundaries, not runtime sessions.
- WebSocket transport owns connection plumbing and frame metadata, not durable session records.
- The Protobuf adapter owns wire conversion, not session persistence.
- Application dispatch must receive normalized request identity; domain modules must not query session storage directly.

## 6. Relationship To Current BindConnection

The current `runtime.authentication.BindConnection` implementation does not create durable session rows.

Future work may later connect BindConnection to session persistence only after a separate implementation gate defines:

- How access-token validation maps to session creation or lookup.
- Whether a session id is returned to clients.
- Whether BindConnection requires an existing session, creates one, or only binds a token-validated actor to a transient connection.
- Whether session validation can satisfy route policy.
- How session expiration or revocation affects existing connections.
- How errors map to public protocol failures.

Until then:

- `RequestIdentity.SessionValidated` remains false.
- Ordinary protected routes continue to use request-level `AuthenticatedRequest`.
- Bound connection identity does not satisfy route policy.

## 7. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

```yaml
adopted_concepts:
  - sessions_are_first_class_lifecycle_records
  - session_expiration_and_revocation_are_distinct_from_login
  - realtime_socket_identity_depends_on_authenticated_session_material
adapted_concepts:
  - first_session_store_is_vibit_postgresql_runtime_sessions
  - opaque_access_token_remains_separate_from_session_id
deferred_concepts:
  - refresh_token_session_extension
  - logout_disconnect_active_socket
  - reconnect_restore_socket_state
rejected_for_now:
  - direct_nakama_session_api_compatibility
```

Pitaya reference mapping:

```yaml
adopted_concepts:
  - session_like_context_is_separate_from_acceptor
  - user_identity_can_be_bound_to_session_context
adapted_concepts:
  - durable_runtime_session_schema_is_vibit_owned
  - handler_context_handoff_remains_application_owned
deferred_concepts:
  - durable_connection_registry
  - group_broadcast_session_membership
  - frontend_backend_cluster_session_routing
rejected_for_now:
  - direct_pitaya_session_api_compatibility
```

Nakama and Pitaya guide capability planning. They do not override vibit's constitution, ADRs, manifests, generated boundaries, or verification commands.

## 8. Future Implementation Queue

After this gate, future work should remain split:

```yaml
future_work_items:
  session_migration_source:
    may_add:
      - runtime/migrations/postgres/000005_create_runtime_sessions.sql
    must_not_add:
      - Go repository code
      - runtime validation behavior
      - WebSocket handshake authentication
  session_repository_boundary:
    may_add:
      - storage-neutral session repository interface
      - focused repository boundary tests
    must_not_add:
      - PostgreSQL adapter behavior
      - route policy use of sessions
  session_postgresql_adapter:
    may_add:
      - PostgreSQL adapter for ratified session repository
      - fake-executor adapter tests
    must_not_add:
      - runtime session validation behavior
      - active WebSocket invalidation
  session_validation_runtime_behavior:
    requires_later_gate: true
  logout_revocation_active_connection_behavior:
    requires_later_gate: true
  reconnect_epoch_behavior:
    requires_later_gate: true
  bound_identity_route_policy:
    requires_later_gate: true
```

## 9. Deferrals

This standard does not authorize:

- SQL migration source.
- Creating `runtime_sessions` table.
- Creating `runtime_session_connections` table.
- Repository interfaces.
- PostgreSQL adapters.
- Runtime session creation, validation, revocation, or cleanup behavior.
- WebSocket handshake authentication.
- HTTP `Authorization`, Bearer, cookie, query-string, or WebSocket subprotocol credential carriers.
- Existing Protobuf envelope changes.
- New Protobuf messages or generated output.
- Route-policy use of session or bound connection identity.
- Logout-triggered active connection invalidation.
- Reconnect, resume, duplicate replacement, or durable epoch behavior.
- New dependencies.
- Memory durable session behavior.
- Direct Nakama or Pitaya public API compatibility.

## 10. Verification

The repository check rule for this gate is:

```text
runtime.postgres_session_persistence_schema_gate
```

The check must verify:

- The standard, translation, ADR, change specs, and conversation log exist.
- Runtime, convention, contract, reference, work, module, and AGENTS artifacts reference the gate.
- No session migration source exists before a future implementation work item.
- WebSocket transport remains credential-neutral.
- The existing Protobuf envelope remains free of proof or session persistence fields.
- No repository, PostgreSQL adapter, generated output, dependency, route-policy, logout/revocation, reconnect, or direct Nakama/Pitaya compatibility behavior is added by this gate.
