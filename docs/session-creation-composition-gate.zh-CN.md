# Session Creation Composition Gate

Status: Draft v0.1
Last updated: 2026-05-17
Scope: persistent runtime session validation 已存在之后，future durable runtime session creation composition 的 gate-only boundary
Depends on: `docs/runtime-session-validation-gate.md`, `decisions/ADR-0066-runtime-session-validation-implementation.md`, `docs/session-repository-boundary.md`, `docs/session-persistence-websocket-handshake-ratification.md`, `docs/authentication-command-protocol-login-route-gate.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0067`

配套英文原文是 `docs/session-creation-composition-gate.md`。英文文件是权威版本。

## 1. 目的

vibit 现在已经有 runtime sessions 的 durable prerequisites：

- PostgreSQL `runtime_sessions` migration source。
- Storage-neutral `runtime/internal/app/session.Repository` interface。
- 该 repository 的 PostgreSQL adapter。
- Application-owned persistent runtime session validator。
- 已有 device-credential login 和 access-token validation behavior。

缺失的边界是 creation composition。`AuthenticateWithDeviceCredential` 可以签发 opaque access token，`PersistentSessionValidator` 可以验证 persisted session，但还没有任何 production path 创建 durable runtime session row。下一个有价值的步骤，是定义 future session creation 应属于哪里，以及它如何与 login 和 token issuance 组合。

成熟 game server 会塑造这个边界：

- Nakama 把 authentication 当作 session lifecycle 的入口，涉及 token material、expiration、refresh、logout 和 management pressure。
- Nakama 也把 session validity 与 active socket disconnect behavior 分开处理。
- Pitaya 保持 acceptors、handler routing 和 session context 分层，所以 durable session creation 应该是 application composition，而不是 transport behavior。

vibit 应该吸收这些经验，把 durable session creation 做成显式、transactional、application-owned。本标准只定义 gate。

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

这是 gate-only standard。它不在代码中创建 runtime sessions。

## 2. Ownership

Future session creation composition 由 application 拥有：

```yaml
future_composition_owner: runtime/internal/app
authentication_service_owner: runtime/internal/app/authentication
session_record_owner: runtime/internal/app/session
postgresql_session_adapter_owner: runtime/internal/platform/persistence/postgres
transport_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
```

规则：

- Future durable session creation 只能通过 application unit-of-work capabilities 调用 `session.Repository.CreateRuntimeSession`。
- Session repository 和 PostgreSQL adapter 保持 storage-oriented。它们不得决定 login 应何时创建 session。
- WebSocket transport 不得创建 durable runtime session rows。
- Protobuf adapters 不得创建 durable runtime session rows。
- Domain modules 不得把创建 authentication runtime sessions 作为 domain commands 的副作用。
- Authentication 可以组合 session creation，仅仅因为第一条已选择的 login path 拥有 successful token issuance，而不是因为 authentication module 拥有 runtime session persistence。

## 3. Future Composition Semantics

后续 implementation slice 可以让 successful device-credential login 在存储 access-token 的同一个 unit of work 内创建 durable runtime session。

Candidate future order：

1. 在 unit-of-work 前拒绝 missing 或 malformed device credential proof。
2. 验证 device credential 和 active player account。
3. 生成 opaque access-token material。
4. 计算并存储 token lookup 和 verifier digests。
5. 生成 server-owned runtime `session_id`。
6. 创建一条链接到已存储 `access_token_record_id` 的 `runtime_sessions` row。
7. Commit unit of work。
8. 只有 commit 后，才返回 client-visible token material 以及 future authorized session material。

规则：

- Session creation 必须发生在 credential proof validation 之后、commit 之前。
- Raw credential material 和 raw access-token material 不得存入 `runtime_sessions`。
- Token lookup digests、token verifier digests、credential lookup digests、credential verifier digests 和 verifier key ids 不得复制到 runtime session records。
- `access_token_record_id` linkage 是 private server metadata。它不是 proof，也不得作为 client credential 暴露。
- Session creation 本身不得让 request routes 变成 session-validated。
- Session creation 不得设置 `RequestIdentity.SessionValidated = true`；validation 仍然是独立的 runtime-session validation concern。
- Session creation 不得关闭、替换或恢复 WebSocket connections。

## 4. Session ID And Lifetime Posture

Future session creation 需要显式 session id 与 lifetime posture：

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

规则：

- Future session id generator 在添加代码前必须单独授权。
- 如果后续 implementation gate 选择该 posture，第一版 session id 可以是 opaque high-entropy text。
- 创建 session 时不得接受 client-supplied session id。
- Session id uniqueness collisions 必须 fail closed，并保持 redacted。
- Session id values 默认是 operationally sensitive，除非后续 observability gate 定义 log-safe redaction format。
- 第一版 session lifetime 可以与 access-token lifetime 对齐，但该选择必须在 implementation gate 中重新说明。
- Refresh、renewal、extension、rotation 和 token-session rekeying 仍然 deferred。

## 5. Unit-Of-Work Composition Boundary

Future login-created sessions 应该以 transaction 方式组合：

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

规则：

- Token storage 和 session creation 应该一起 commit 或 rollback。
- 如果 selected posture 要求 session creation，而 session creation 失败，则 token 不得返回给 client。
- 如果后续 implementation 保持 session creation optional，public result 必须显式说明该 posture，并用 tests 覆盖。
- Repository acquisition failures 和 session creation failures 必须折叠成 redacted dependency 或 authentication-unavailable public errors。
- Future implementation 必须测试 commit failure behavior，确保 raw token 或 session material 不会在 commit 前作为成功结果返回。

## 6. Relationship To Validation And Route Policy

该 gate 不改变 validation 或 route policy：

```yaml
runtime_session_validation_changed: false
request_identity_session_validated_policy_changed: false
route_policy_session_identity_added: false
route_policy_bound_identity_added: false
```

规则：

- `PersistentSessionValidator` 保持 lookup-only，并且不会因为该 gate 接入 route policy。
- Access-token validation 仍然是当前 protected-route proof path。
- Future session creation 可以产生 durable row，后续 validation 或 route-policy slice 可以使用它，但本 gate 不选择该 policy。
- First-message bound identity 仍然不能满足 ordinary protected route policy。
- Metadata-only player id、session id 和 connection metadata 仍然不是 authenticated proof。

## 7. Relationship To WebSocket And Protocol

该 gate 不改变 WebSocket 或 Protobuf behavior：

```yaml
websocket_transport_credential_neutral: true
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
protobuf_session_messages_added: false
existing_protobuf_envelope_change_added: false
generated_output_added: false
```

规则：

- WebSocket transport 不得因为该 gate 解析 Authorization headers、bearer values、cookies、query-string tokens、session tokens 或 `Sec-WebSocket-Protocol` authentication material。
- Existing Protobuf envelope 保持不变。
- 这里不授权 session creation command、response field、session carrier、system message、generated Protobuf output 或 generated contract shape。
- Future protocol gate 必须先授权任何 client-visible session id 或 session result shape，login response 才能暴露它。
- 这里不授权 reconnect、resume、duplicate replacement、durable connection epoch、logout disconnect、presence、rooms、parties、groups 或 match attachment behavior。

## 8. Error And Redaction Boundary

Future session creation failures 不应泄露敏感内部细节：

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

规则：

- Public errors 不得泄露 raw token material、raw credential material、session ids、token record ids、lookup digests、verifier digests、verifier key ids、SQL argument values、Authorization headers、cookies 或 WebSocket credential carriers。
- Internal test-only failure classes 可以更具体，但 logged 和 public output 必须保持 redacted。
- Session id generation failures 和 uniqueness conflicts 不得回显 candidate session ids。
- Session creation 不得为 `session_id`、`token_record_id` 或 `player_id` 引入新的 log-safe status。

## 9. Future Implementation 的测试要求

后续 session creation implementation 必须包含 focused tests：

- Successful login 通过 `session.Repository.CreateRuntimeSession` 创建 exactly one active runtime session row。
- Session creation 只发生在 credential proof validation、player account validation、token generation、digest computation 和 token storage 之后。
- Runtime session 链接 `access_token_record_id`，但不存储 raw proof 或 digest material。
- Session id generation 根据后续选择的 generator posture 拒绝 missing、malformed、duplicate 或 low-entropy values。
- Session creation failure 会阻止 successful token/session material return。
- Unit-of-work commit failure 会阻止 successful token/session material return。
- Repository acquisition errors 和 session creation errors 保持 redacted。
- Access-token validation behavior 保持不变。
- `RequestIdentity.SessionValidated` 保持 false，直到独立 validation/route-policy path 设置它。
- WebSocket transport 和 Protobuf envelope 保持不变。

Live PostgreSQL verification 可以保持 opt-in，除非后续 implementation work item 要求它。

## 10. Nakama And Pitaya Reference Mapping

Nakama reference mapping：

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

Pitaya reference mapping：

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

Nakama 和 Pitaya guide capability planning。它们不会覆盖 vibit constitution、ADRs、manifests、generated boundaries 或 verification commands。

## 11. Future Implementation Queue

此 gate 之后，future work 应继续拆分：

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

不要在没有新 ADR 的情况下把这些合并成一个宽泛的 session subsystem slice。

## 12. Verification

本 gate 的 repository verification 是：

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check change define-session-creation-composition-gate --json
node tools/vibit check all --json
```

Repository check rule 是：

```yaml
runtime.session_creation_composition_gate
```
