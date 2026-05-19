# PostgreSQL Session Persistence Schema Gate

Status: Draft v0.1
Last updated: 2026-05-17
Scope: first-message connection binding 之后 first durable runtime session persistence schema 的 gate-only boundary
Depends on: `docs/session-persistence-websocket-handshake-ratification.md`, `docs/session-persistence-websocket-handshake-decision-gates.md`, `docs/first-message-connection-binding-implementation-gate.md`, `docs/postgresql-persistence-boundary.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0059`

配套英文原文是 `docs/postgres-session-persistence-schema-gate.md`。英文文件是权威版本。

## 1. 目的

vibit 现在已经具备最小 authentication 和 realtime-connection foundation：

- Public device credential login route。
- Opaque access-token validation。
- 通过 `vibit.authentication.v1.AuthenticatedRequest` 做 request-level protected route validation。
- PostgreSQL startup composition for authentication。
- First-message `runtime.authentication.BindConnection` system route。

下一个 durable boundary 是 session persistence。Nakama 这样的成熟 game server 把 authenticated sessions 当作 richer realtime features 之前的一等生命周期概念。Pitaya 在 session-like connection context 和 binding 上提供了有用 vocabulary，同时把 acceptor 和 handler concerns 分开。vibit 应该吸收这些经验，但不复制 public APIs，也不把 authentication 放进 transport。

本标准定义 future PostgreSQL session schema 的 gate。

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

这是 schema gate。它不添加 SQL migration source、Go repository interfaces、PostgreSQL adapters、runtime validation behavior 或 connection lifecycle behavior。

## 2. 选择的第一 durable session target

第一 durable session persistence target 是 PostgreSQL：

```yaml
selected_first_session_store: postgres
future_session_logical_table: runtime_sessions
future_migration_source_candidate: runtime/migrations/postgres/000005_create_runtime_sessions.sql
future_session_repository_boundary: separate_future_work_item
```

理由：

- PostgreSQL 已经是 vibit 第一个 accepted authoritative durable store。
- Session persistence 需要和 player account state、token verifier state 有清晰 transaction relationship。
- SQL-first schema gate 比临时 in-memory state 更容易让 agent 检查。
- 推迟 Redis-like store，避免在真正需要前引入 dependency 和 operations decision。

## 3. Future `runtime_sessions` table candidate

Future first session migration 可以定义一个 logical table：

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

规则：

- `session_id` 是 server-generated opaque identifier，不是 access token。
- `actor_kind` 和 `actor_id` 使用 ratified authentication/session vocabulary。
- 第一 player session posture 要求 `player_id`。
- `session_status` 必须是 closed vocabulary，例如 `active`、`revoked` 或 `expired`。
- `issued_at`、`expires_at`、`last_seen_at`、`created_at` 和 `updated_at` 必须使用 UTC timestamps。
- `revoked_at` 和 `revocation_reason` 支持 future logout/revocation work，但不实现 active connection invalidation。
- `access_token_record_id` 可以作为 future opaque linkage 指向 token verifier storage，但 raw token text 和 digest bytes 不得复制进 session table。

## 4. 明确推迟 connection records

第一 session schema gate 不选择 durable connection table。

```yaml
runtime_session_connections_table_selected_now: false
connection_binding_registry_added: false
active_connection_invalidation_added: false
reconnect_epoch_behavior_added: false
```

添加以下任何内容前，future gate 必须先定义 connection lifecycle storage：

- `runtime_session_connections` table。
- Durable `connection_id` to `session_id` binding rows。
- Duplicate connection replacement。
- Reconnect/resume state。
- Metadata handoff 之外的 connection epoch policy。
- Logout-triggered active WebSocket close or invalidation。
- Presence、rooms、groups、parties 或 match attachment。

这会保持第一 schema 小而明确，避免 session persistence 静默变成完整 realtime connection manager。

## 5. Ownership

Future session persistence boundary 由 runtime 拥有：

```yaml
runtime_session_owner: runtime/internal/app
future_repository_owner_candidate: runtime/internal/app/session
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
player_module_owns_session_table: false
authentication_module_owns_session_table: false
websocket_transport_owns_session_table: false
```

规则：

- Player module 拥有 player account lifecycle，不拥有 runtime sessions。
- Authentication module 拥有 credential 和 token verifier record boundaries，不拥有 runtime sessions。
- WebSocket transport 拥有 connection plumbing 和 frame metadata，不拥有 durable session records。
- Protobuf adapter 拥有 wire conversion，不拥有 session persistence。
- Application dispatch 必须接收 normalized request identity；domain modules 不得直接查询 session storage。

## 6. 与当前 BindConnection 的关系

当前 `runtime.authentication.BindConnection` implementation 不创建 durable session rows。

只有后续单独 implementation gate 定义以下内容后，future work 才能把 BindConnection 接入 session persistence：

- Access-token validation 如何映射到 session creation 或 lookup。
- 是否向 client 返回 session id。
- BindConnection 是要求已有 session、创建 session，还是只把 token-validated actor 绑定到 transient connection。
- Session validation 是否能满足 route policy。
- Session expiration 或 revocation 如何影响 existing connections。
- Errors 如何映射到 public protocol failures。

在此之前：

- `RequestIdentity.SessionValidated` 保持 false。
- 普通 protected routes 继续使用 request-level `AuthenticatedRequest`。
- Bound connection identity 不满足 route policy。

## 7. Nakama 和 Pitaya reference mapping

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

Nakama 和 Pitaya 只指导 capability planning。它们不覆盖 vibit constitution、ADRs、manifests、generated boundaries 或 verification commands。

## 8. Future implementation queue

该 gate 之后，future work 应继续拆分：

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

本标准不授权：

- SQL migration source。
- 创建 `runtime_sessions` table。
- 创建 `runtime_session_connections` table。
- Repository interfaces。
- PostgreSQL adapters。
- Runtime session creation、validation、revocation 或 cleanup behavior。
- WebSocket handshake authentication。
- HTTP `Authorization`、Bearer、cookie、query-string 或 WebSocket subprotocol credential carriers。
- Existing Protobuf envelope changes。
- New Protobuf messages 或 generated output。
- Route-policy use of session 或 bound connection identity。
- Logout-triggered active connection invalidation。
- Reconnect、resume、duplicate replacement 或 durable epoch behavior。
- New dependencies。
- Memory durable session behavior。
- Direct Nakama 或 Pitaya public API compatibility。

## 10. Verification

该 gate 的 repository check rule 是：

```text
runtime.postgres_session_persistence_schema_gate
```

检查必须验证：

- Standard、translation、ADR、change specs 和 conversation log 存在。
- Runtime、convention、contract、reference、work、module 和 AGENTS artifacts 引用该 gate。
- Future implementation work item 前不存在 session migration source。
- WebSocket transport 保持 credential-neutral。
- Existing Protobuf envelope 没有 proof 或 session persistence fields。
- 该 gate 没有新增 repository、PostgreSQL adapter、generated output、dependency、route-policy、logout/revocation、reconnect 或 direct Nakama/Pitaya compatibility behavior。
