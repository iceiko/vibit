# Session Repository Boundary

状态：Draft v0.1
最后更新：2026-05-17
范围：在 PostgreSQL `runtime_sessions` migration source 之后，为未来 storage-neutral runtime session repository 定义 gate-only boundary
依赖：`docs/session-persistence-websocket-handshake-ratification.md`、`docs/postgres-session-persistence-schema-gate.md`、`decisions/ADR-0060-runtime-sessions-migration-source.md`、`docs/reference-game-server-alignment.md`
Canonical decision：`ADR-0061`

英文原文是 `docs/session-repository-boundary.md`。英文文件是架构、契约和治理的权威版本；本文是面向中文读者的对应翻译。

## 1. 目的

仓库现在已经有具体的 PostgreSQL `runtime_sessions` migration source，但还没有 application session repository、PostgreSQL session adapter、runtime session creation，也没有 runtime session validation。

下一个有价值的边界，是未来 session lifecycle behavior 所需的 repository seam。成熟游戏服务器对这里的设计有很强参考价值：

- Nakama 把 authenticated session 作为一等 lifecycle object，包含 expiration、refresh、logout 和运维管理问题。
- Pitaya 让 handler 能拿到 session-like context，同时保持 acceptor、routing 和 handler execution 分层。

vibit 应该吸收这些经验，但把它们转换成 storage-neutral、application-owned repository boundary。本文只定义未来边界。

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

这只是 gate-only standard。它不添加 Go session repository code、SQL adapter code、runtime session behavior 或 protocol behavior。

## 2. 所有权

未来 repository 由 runtime application 层拥有：

```yaml
future_repository_owner_candidate: runtime/internal/app/session
future_postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
runtime_session_table_owner: runtime.session
authentication_module_owns_runtime_sessions: false
player_module_owns_runtime_sessions: false
websocket_transport_owns_runtime_sessions: false
protocol_adapter_owns_runtime_sessions: false
```

规则：

- Repository interface 必须是 storage-neutral、application-facing。
- Interface 不得提到 PostgreSQL、pgx、SQL rows、goose migrations、prepared statements 或 database transaction 实现细节。
- PostgreSQL adapter 未来可以在 `runtime/internal/platform/persistence/postgres` 下实现该 interface，但必须等待单独 adapter gate。
- Authentication 可以只通过 `authentication_access_tokens(token_record_id)` 与 session records 保持 record reference；它不拥有 session lifecycle storage。
- Player account storage 拥有 player lifecycle state，不拥有 runtime sessions。
- WebSocket transport 拥有 connection plumbing，不拥有 durable session persistence。
- Protocol adapter 拥有 wire conversion，不拥有 session repository behavior。

## 3. 未来 Repository 能力词汇

后续 implementation gate 可以选择更小或重命名后的 API，但第一版候选能力族是：

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

能力规则：

- `CreateRuntimeSession` 只有在后续 runtime behavior gate 定义 caller、transaction boundary、id generation、expiration 和 token linkage 后，才可以创建 lifecycle row。
- `GetRuntimeSession` 是 storage lookup。它不得验证 proof，也不得创建 request identity。
- `FindActiveSessionByID` 未来可以表达 active-status filtering，但不得把 token validation 和 session validation 隐藏合并成一个操作。
- `UpdateRuntimeSessionLastSeen` 只能修改 lifecycle metadata。它不得 refresh access tokens、延长 sessions 或暗示 presence。
- `MarkRuntimeSessionExpired` 和 `RevokeRuntimeSession` 只是 durable state transitions。没有后续 active-connection gate 时，它们不得关闭 active WebSocket connections。
- `ListActiveSessionsForPlayer` 只有在后续 gate 定义 permission、pagination、information leakage boundary，以及 admin inspection 与 player self-inspection 语义后才允许。

Repository 必须返回 typed、application-owned records 和 errors。它不得把 raw SQL rows 或 database driver errors 直接返回给 domain modules 或 protocol adapters。

## 4. 数据边界

未来 repository 可以使用已有 `runtime_sessions` table：

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

规则：

- Repository 不得存储或返回 raw access-token text、raw credential material、lookup digests、verifier digests、verifier keys 或 HMAC inputs。
- `access_token_record_id` 是可选 linkage metadata，不是 proof。
- `session_id` 是 server-owned runtime session identifier，不是 access token。
- Session status 必须保持 closed vocabulary，例如 active、revoked 或 expired。
- Expiration 和 revocation 是 session lifecycle states，不是 proof-verifier states。
- Runtime session records 不得变成 WebSocket connection registry rows。

## 5. 与 Authentication 和 Request Identity 的关系

未来 session repository 不拥有 token validation。

```yaml
token_validation_owner: runtime/internal/app/authentication
session_repository_token_validation_owner: false
request_identity_owner: runtime/internal/app
session_repository_sets_request_identity: false
session_validated_status_added: false
```

规则：

- Access-token validation 仍由 application authentication service 和 route access-token validator 拥有。
- Repository 未来可以在 token validation 之后被使用，但不得 parse token proof 或 compare token verifier digests。
- `RequestIdentity.SessionValidated` 保持 false，直到后续 runtime session validation gate 明确定义 persisted session 如何验证并绑定到 request identity。
- `runtime.authentication.BindConnection` 的 bound connection identity 不会通过这个 repository boundary 满足普通 protected route policy。
- Domain modules 仍然消费 normalized `RequestIdentity`，而不是 session repository records。

## 6. 与 WebSocket 和 Protocol 的关系

该边界不改变 WebSocket 或 Protobuf 行为：

```yaml
websocket_transport_credential_neutral: true
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
protobuf_session_messages_added: false
existing_protobuf_envelope_change_added: false
connection_registry_added: false
```

规则：

- WebSocket transport 不得为了本边界解析 access tokens、cookies、query-string tokens、Authorization headers 或 session tokens。
- 现有 Protobuf envelope 保持不变。
- 这里不授权 session protocol messages 或 generated output。
- 这里不授权 logout-triggered WebSocket close、duplicate connection replacement、reconnect、resume 或 durable connection epoch behavior。

## 7. Nakama 和 Pitaya 参考映射

Nakama reference mapping：

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

Pitaya reference mapping：

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

Nakama 和 Pitaya 用来指导 capability planning。它们不覆盖 vibit constitution、ADRs、manifests、generated boundaries 或 verification commands。

## 8. 后续实现队列

这个边界之后，后续工作仍应拆开：

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

不要在没有新 ADR 的情况下把这些合并成一个宽泛 session subsystem slice。

## 9. 验证

该边界的 repository verification 是：

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-session-repository-boundary --json
node tools/vibit check all --json
```

Repository check rule 是：

```yaml
runtime.session_repository_boundary
```

该检查必须确认 standard、translation、ADR、conversation log、change specs、manifests、AGENTS guides 和 rule registry 存在，同时 Go session repository code、PostgreSQL session adapter code、runtime behavior、protocol changes、WebSocket credential carriers、new migrations、dependencies 和 direct Nakama/Pitaya API compatibility 仍然不存在。

## 10. Agent 指引

添加后续 session work 时：

- 阅读本文、`ADR-0061`、`ADR-0060`、`docs/postgres-session-persistence-schema-gate.md` 和 `docs/session-persistence-websocket-handshake-ratification.md`。
- 保持 session repository behavior application-owned 且 storage-neutral。
- 保持 PostgreSQL adapter behavior 在单独 gate 后面。
- 保持 token validation 由 authentication service behavior 拥有。
- 在 runtime session validation 被单独 ratify 前，保持 `RequestIdentity.SessionValidated` 为 false。
- 保持 WebSocket transport credential-neutral。
- 不要复制 Nakama 或 Pitaya public APIs。
