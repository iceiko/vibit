# Session PostgreSQL Adapter Gate

Status: Draft v0.1
Last updated: 2026-05-17
Scope: future PostgreSQL adapter 的 gate-only boundary，该 adapter 未来实现 storage-neutral runtime session repository interface
Depends on: `docs/session-repository-boundary.md`, `decisions/ADR-0062-session-repository-interface-implementation.md`, `runtime/migrations/postgres/000005_create_runtime_sessions.sql`, `docs/postgresql-persistence-boundary.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0063`

配套英文原文是 `docs/session-postgresql-adapter-gate.md`。英文文件是权威版本。

## 1. 目的

vibit 现在已经具备：

- PostgreSQL `runtime_sessions` migration source。
- Storage-neutral `runtime/internal/app/session.Repository` interface。
- 还没有该 interface 的 PostgreSQL adapter。
- 还没有 runtime session creation、lookup、validation、revocation execution、cleanup job 或 route-policy use of persisted sessions。

下一个有价值的步骤是 future PostgreSQL adapter 的 gate。成熟 game server 会塑造这个边界：

- Nakama 把 sessions 当作 first-class lifecycle records，需要 durable lookup、expiration、revocation、logout 和 management-ready listing。
- Pitaya 保持 session/context handoff 与 acceptors、routing 分离，这意味着 transport 不应该拥有 durable session persistence。

vibit 应该吸收这些经验，把 future adapter 限定为 persistence-only 且 transaction-bound。本标准只定义 adapter gate。

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

这是 gate-only standard。它不添加 Go adapter code 或 runtime behavior。

## 2. Ownership

Future adapter 由 platform persistence 拥有：

```yaml
repository_interface_owner: runtime/internal/app/session
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
transaction_boundary_owner: runtime/internal/platform/tx
application_runtime_behavior_owner: runtime/internal/app
authentication_module_owns_session_adapter: false
websocket_transport_owns_session_adapter: false
protocol_adapter_owns_session_adapter: false
```

规则：

- Adapter 后续可以实现 `runtime/internal/app/session.Repository`。
- Adapter 后续只能在 `runtime/internal/platform/persistence/postgres` 内 import `runtime/internal/app/session` 和 `github.com/jackc/pgx/v5`。
- Storage-neutral `runtime/internal/app/session` package 不得 import PostgreSQL、pgx、SQL rows、WebSocket transport 或 generated Protobuf packages。
- Adapter 不得成为 runtime session validation policy 的所有者。它只把 repository calls 映射到 SQL，并把 SQL rows 映射回 application records。
- Adapter 不得解析 access-token proof、比较 verifier digests、创建 `RequestIdentity`，或设置 `RequestIdentity.SessionValidated`。

## 3. Future Adapter Surface

后续 implementation slice 可以添加：

```yaml
future_adapter_files:
  - runtime/internal/platform/persistence/postgres/session_repository.go
  - runtime/internal/platform/persistence/postgres/session_repository_test.go
future_constructor:
  - NewSessionRepositoryForUnitOfWork(executor Executor)
future_unit_of_work_factory:
  - NewSessionRepository() (session.Repository, error)
```

Adapter 应只实现现有 repository methods：

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

规则：

- Adapter implementation 必须留在已经 ratified 的 `Executor` 和 unit-of-work patterns 后面。
- Repository methods 不得直接 open、commit 或 rollback transactions。
- Future unit-of-work factory 只有在 adapter implementation slice 授权后才可以暴露 `NewSessionRepository`。
- Query 和 mutation inputs 必须先通过 `runtime/internal/app/session` helpers normalize，再执行 SQL。
- 返回 rows 必须通过 `session.NormalizeRuntimeSessionRecord` normalize。
- Adapter errors 必须 typed 且 redacted。它们可以包含 action labels 和 constraint labels，但不得包含 raw token material、digest bytes、verifier key ids、SQL argument values 或 access-token proof。

## 4. SQL Shape Boundary

Future adapter 只能使用已有 `runtime_sessions` table：

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

Future SQL shape 规则：

- `CreateRuntimeSession` 可以 insert 一条 `runtime_sessions` row 并返回 inserted row。
- `GetRuntimeSession` 可以按 `session_id` select 一条 row，不做 status 或 expiration 过滤。
- `FindActiveSessionByID` 可以过滤 `session_status = 'active'` 和 `expires_at > observed_at`；它不得验证 token proof。
- `UpdateRuntimeSessionLastSeen` 可以为一条 row 更新 `last_seen_at` 和 `updated_at`。
- `MarkRuntimeSessionExpired` 可以把一个 active session transition 为 `expired`。
- `RevokeRuntimeSession` 可以把一个 session transition 为 `revoked`，并写入 `revoked_at` 和 `revocation_reason`。
- `ListActiveSessionsForPlayer` 可以按 explicit bounded limit 列出某 player 的 active、unexpired sessions。

Adapter 不得添加 cleanup jobs、connection registries、logout-triggered socket invalidation、reconnect state、presence、rooms、parties、match runtime 或 social graph behavior 的 SQL。

## 5. Error Mapping

Future adapter errors 应跟随现有 PostgreSQL repository 风格：

```yaml
future_error_classes:
  - ErrSessionRecordNotFound
  - ErrSessionRecordConflict
  - ErrSessionRecordConstraint
  - ErrSessionRecordStale
```

规则：

- `pgx.ErrNoRows` 映射到 not-found class。
- Duplicate key conflicts 映射到 conflict class。
- not-null、foreign-key 和 check constraint failures 映射到 constraint class。
- Lifecycle update 中 zero affected rows 映射到 not found 或 stale state，具体由 implementation gate 定义。
- Public runtime validation failures 不得泄露失败来自 session id lookup、expiration、revocation、actor mismatch、token mismatch 或 player state。这个 public collapse 属于后续 runtime validation gate。

## 6. Relationship To Authentication And Request Identity

该 gate 不改变 authentication behavior：

```yaml
token_validation_owner: runtime/internal/app/authentication
session_adapter_validates_token_proof: false
session_adapter_sets_request_identity: false
session_validated_status_added: false
```

规则：

- Authentication token validation 仍由 application authentication service 和 route access-token validator 拥有。
- Adapter 后续可以存储或返回 `access_token_record_id` linkage，但只能作为 opaque metadata。
- Adapter 不得读取 token lookup digests、token verifier digests、credential digests、raw token text 或 raw credential material。
- `RequestIdentity.SessionValidated` 保持 false，直到后续 runtime session validation gate 定义行为。
- Bound connection identity 仍然不会通过该 gate 满足普通 protected route policy。

## 7. Relationship To WebSocket And Protocol

该 gate 不改变 WebSocket 或 Protobuf behavior：

```yaml
websocket_transport_credential_neutral: true
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
protobuf_session_messages_added: false
existing_protobuf_envelope_change_added: false
connection_registry_added: false
```

规则：

- WebSocket transport 不得因为该 gate 解析 access tokens、cookies、query-string tokens、Authorization headers 或 session tokens。
- Existing Protobuf envelope 保持不变。
- 这里不授权 session protocol messages 或 generated output。
- 这里不授权 logout-triggered WebSocket close、duplicate connection replacement、reconnect、resume、durable connection epoch、presence、rooms、parties 或 match attachment behavior。

## 8. Future Implementation 的测试要求

后续 adapter implementation 必须包含 focused tests：

- Interface conformance to `session.Repository`。
- Constructor/executor nil rejection。
- 每个 repository method 的 SQL shape。
- Argument normalization 和 UTC timestamp handoff。
- Row scanning into `session.RuntimeSession`，包括 nullable revocation 和 token linkage fields。
- `pgx.ErrNoRows`、conflict、foreign-key、check、not-null error mapping。
- Update/revoke/expire paths 的 zero affected row handling。
- Bounded listing limit behavior。
- Transaction neutrality：repository 内没有直接 `BEGIN`、`COMMIT` 或 `ROLLBACK` SQL。
- Redaction：errors 或 tests 中没有 raw token、raw credential、digest bytes、verifier keys、Authorization headers、cookies 或 WebSocket credential carriers。

Live PostgreSQL verification 可以保持 opt-in，除非后续 work item 在 disposable PostgreSQL verification environment 下把它设为 mandatory。

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

Nakama 和 Pitaya 只指导 capability planning。它们不覆盖 vibit constitution、ADRs、manifests、generated boundaries 或 verification commands。

## 10. Future Implementation Queue

该 gate 之后，future work 应继续拆分：

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

没有新 ADR，不要把这些合并成一个宽泛 session subsystem slice。

## 11. Verification

该 gate 的 repository verification 是：

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-session-postgresql-adapter-gate --json
node tools/vibit check all --json
```

Repository check rule 是：

```yaml
runtime.session_postgresql_adapter_gate
```
