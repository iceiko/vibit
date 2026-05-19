# 绑定身份路由策略关口

状态：Draft v0.1
最后更新：2026-05-18
范围：为未来路由策略使用请求级 proof、持久 runtime session 验证，以及 first-message 连接绑定身份定义 gate-only 边界
依赖：`docs/access-token-protocol-carrier-route-protection-gate.md`、`docs/first-message-connection-binding-implementation-gate.md`、`docs/runtime-session-validation-gate.md`、`docs/session-creation-composition-gate.md`、`decisions/ADR-0068-session-creation-composition-implementation.md`、`docs/reference-game-server-alignment.md`
Canonical decision：`ADR-0069`

配对英文源文档是 `docs/bound-identity-route-policy-gate.md`。英文文件是权威版本。

## 1. 目的

vibit 现在已经有几块和身份相关的 runtime 能力：

- 通过显式 Protobuf payload wrapper 实现的请求级 access-token 路由保护。
- 通过 `runtime.authentication.BindConnection` 实现的 first-message 连接绑定。
- 持久化 `runtime_sessions` 和 session repository。
- 应用层拥有的持久 runtime session validator。
- 与 device-credential login 组合的登录时 durable runtime session 创建。

这些能力故意还没有合并成一个路由策略。当前 protected route 姿态仍然是每个受保护请求验证 access-token proof。绑定连接身份和持久 session 验证已经是独立应用能力，但普通 domain routes 还没有依赖它们。

下一步有价值的工作，是先定义未来 route-policy 边界，再实现。成熟游戏服务器对这个边界有直接启发：

- Nakama 把认证后的 session material 作为 gameplay API 访问基础，同时把 session lifetime、logout、refresh、socket 行为当作不同 lifecycle 决策。
- Pitaya 把 acceptor、session、route handler 分开，使 handler 接收上下文，而不是解析 transport credential。

vibit 应该吸收这些经验，把 route identity policy 做成应用层拥有、显式、可测试、分阶段的边界。本标准只定义 gate。

```yaml
bound_identity_route_policy_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0150
decision: ADR-0069
check_rule: runtime.bound_identity_route_policy_gate
future_policy_owner: runtime/internal/app
future_policy_source_candidate: runtime/internal/app/route_authentication.go
future_policy_test_candidate: runtime/internal/app/route_authentication_test.go
existing_access_token_proof_path: request_level_authenticated_request_wrapper
existing_connection_binding_route: runtime.authentication.BindConnection
existing_session_validator: runtime/internal/app.PersistentSessionValidator
route_policy_session_identity_added: false
route_policy_bound_identity_added: false
ordinary_protected_routes_use_bound_identity: false
ordinary_protected_routes_use_session_validated_identity: false
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

这是 gate-only 标准。它不改变生产路由授权行为。

## 2. 所有权

未来绑定身份路由策略由应用层拥有：

```yaml
future_route_policy_owner: runtime/internal/app
authentication_service_owner: runtime/internal/app/authentication
session_validator_owner: runtime/internal/app
session_repository_owner: runtime/internal/app/session
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
transport_owner: runtime/internal/platform/transport/ws
domain_handler_owner: runtime/internal/modules/*
```

规则：

- Route policy 必须留在 `runtime/internal/app`。
- Route policy 只能通过应用层类型消费标准化后的 `RequestIdentity`、access-token validation result、connection binding identity、runtime session validation result。
- Route policy 不得 import WebSocket transport 包、generated Protobuf 包、PostgreSQL adapter、SQL row、migration 包或 provider SDK。
- WebSocket transport 必须继续保持 credential-neutral。
- Protobuf adapter 只能通过已经授权的 protocol carrier 边界携带或解包 proof；它们不得决定 policy。
- Domain module 必须接收标准化后的 identity context，不得解析 access token、session id、connection id 或 transport credential carrier。

## 3. 未来策略族

后续 implementation slice 可以引入显式 route policy families，而不是一个隐式 protected-route default：

```yaml
candidate_policy_families:
  public:
    examples:
      - runtime.authentication.AuthenticateWithDeviceCredential
    requirement: no authenticated proof required
  request_token_required:
    requirement: fresh request-level access-token proof validates for the route
  bound_connection_required:
    requirement: connection has a server-observed bound identity matching the request identity
  session_validated_required:
    requirement: already-validated actor identity also validates against an active durable runtime session row
  bound_session_required:
    requirement: bound connection identity and active durable session validation both match the request identity
```

规则：

- Public routes 必须显式列出。
- 不得有隐式 public gameplay route default。
- Metadata-only identity 不能满足任何 protected policy family。
- 除非后续 implementation slice 明确为具名 route 选择这种姿态，否则绑定连接身份不能成为 request-level proof 的通用替代品。
- 持久化 `session_id` 本身不能成为 proof。
- Session-validated identity 需要已经验证的 actor identity 和 active durable runtime session row。
- Policy family selection 必须按 route scope 声明，并且不依赖 live PostgreSQL 也能测试。

## 4. 推荐的首个实现姿态

本 gate 之后推荐的首个实现姿态应该保守：

```yaml
recommended_first_implementation_posture:
  default_domain_route_policy: request_token_required
  public_routes:
    - runtime.authentication.AuthenticateWithDeviceCredential
  system_binding_route:
    route: runtime.authentication.BindConnection
    policy: request_token_required_for_binding_request
  bound_connection_policy:
    status: available_for_explicit_future_routes_only
  session_validated_policy:
    status: available_for_explicit_future_routes_only
  bound_session_policy:
    status: deferred
```

规则：

- 首个实现中，普通受保护 domain routes 应继续要求 request-level access-token proof。
- `BindConnection` 可以建立 process-local bound identity，但这不会自动授权普通 domain routes。
- 首个 route-policy implementation 可以添加未来 route classes 的 policy vocabulary 和 tests，而不改变现有 domain route 行为。
- 任何不再要求 per-request proof 的 route 都必须显式具名，并且测试必须说明信任哪个 identity source。
- Inventory route behavior 不得因为本 gate 意外改变。

## 5. 组合顺序

未来 route policy 必须用确定顺序组合 identity checks：

```yaml
candidate_composition_order:
  - normalize_route_key
  - classify_route_policy_family
  - accept_explicit_public_route_without_authenticated_identity
  - require_structural_proof_or_bound_identity_for_protected_family
  - validate_access_token_proof_when_required
  - validate_bound_connection_identity_when_required_by_route_family
  - validate_runtime_session_when_required_by_route_family
  - build_one_normalized_request_identity
  - dispatch_domain_handler_only_after_policy_success
```

规则：

- Domain dispatch 只能在 policy success 之后发生。
- 如果使用多个 identity source，它们必须在 actor kind、actor id、player id、需要时的 session id、需要时的 connection id、需要时的 connection epoch 上一致。
- 身份来源不一致必须 fail closed。
- Policy failure 不得 mutate token state、session state、connection binding state、inventory state 或 domain state。
- `SessionValidated = true` 只能由 runtime session validator 产生后被接受；route policy 不能靠断言自行设置它。
- 未来 last-seen updates 不属于本 gate，除非后续实现显式授权。

## 6. 错误和脱敏边界

未来实现必须保持 public failure 稳定且脱敏：

```yaml
candidate_public_errors:
  request_token_missing: AUTHENTICATION_TOKEN_MISSING
  request_token_malformed: AUTHENTICATION_TOKEN_MALFORMED
  request_token_invalid: AUTHENTICATION_TOKEN_INVALID
  request_token_unavailable: AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
  bound_connection_required: CONNECTION_BINDING_REQUIRED
  bound_connection_invalid: CONNECTION_BINDING_TOKEN_INVALID
  session_invalid: SESSION_INVALID
```

规则：

- Public failures 不得暴露拒绝原因是 token lookup、token mismatch、account inactive、session missing、session expired、session revoked、connection unbound、epoch mismatch、actor mismatch、player mismatch、repository failure，还是 policy-family mismatch。
- Errors、logs、events 和 test output 不得包含 raw access-token text、raw credential material、lookup digest、verifier digest、verifier key id、session id、token record id、Authorization header、cookie、query string、WebSocket subprotocol value 或 inner payload bytes。
- Route policy 只应在 application error boundary 返回 route-specific public errors。
- 更精确的 internal reasons 可以只存在于测试，但不能泄漏到 public output。

## 7. 与现有能力的关系

本 gate 不改变现有行为：

```yaml
access_token_validation_changed: false
connection_binding_changed: false
runtime_session_validation_changed: false
session_creation_changed: false
request_identity_session_validated_policy_changed: false
route_policy_session_identity_added: false
route_policy_bound_identity_added: false
```

规则：

- 在后续 implementation work item 改变它之前，`RouteProtector` 可以继续要求普通 protected routes 使用 request-level access-token proof。
- `ConnectionBinder` 可以继续设置 `SessionValidated = false`。
- `PersistentSessionValidator` 可以继续只在应用代码直接调用时验证。
- `AuthenticateWithDeviceCredential` 可以继续创建 durable runtime sessions，但不通过 protocol responses 暴露这些 session id。
- PostgreSQL session adapter 仍然只做 persistence。
- Authentication service 仍然负责 proof/token/session-creation composition；它不拥有 route policy。

## 8. 与 WebSocket 和 Protocol 的关系

本 gate 不授权 WebSocket 或 Protobuf 变更：

```yaml
websocket_transport_credential_neutral: true
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
protobuf_session_messages_added: false
existing_protobuf_envelope_change_added: false
generated_output_added: false
```

规则：

- 本 gate 下 WebSocket transport 不得解析 Authorization headers、bearer values、cookies、query-string tokens、session tokens 或 `Sec-WebSocket-Protocol` authentication material。
- 现有 Protobuf envelope 保持不变。
- 本 gate 不授权 session carrier field、login response session id、handshake message、reconnect message、route policy Protobuf message、generated Protobuf output 或 generated contract shape。
- 在 route policy 要求客户端输入任何 session id 或 session proof 之前，未来 protocol session carrier gate 必须先授权它。

## 9. 延后事项

本 gate 不授权：

- 生产 route-policy 使用 bound identity。
- 生产 route-policy 使用 session-validated identity。
- 从普通 protected routes 移除 per-request access-token proof。
- WebSocket handshake authentication。
- Transport credential carriers。
- Protobuf envelope changes。
- Protobuf session messages 或 generated output。
- Logout、refresh、cleanup、token rotation、token-session rekeying 或 token validation audit mutation。
- Logout 或 revocation 时让 active connection invalidation。
- Reconnect、resume、duplicate connection replacement 或 durable connection epoch policy。
- Presence、rooms、parties、groups、chat、matchmaking、match runtime、broadcast groups 或 social modules。
- Metrics、dashboards、admin APIs、session management APIs 或 operations posture。
- Memory durable session behavior。
- Direct Nakama 或 Pitaya API compatibility。

## 10. 未来实现的测试要求

后续 bound identity route-policy implementation 必须包含聚焦测试：

- Public routes 仍然显式，且不要求 proof。
- 除非显式重新分类，普通 protected routes 仍然要求 request-level proof。
- Metadata-only identity 对所有 protected policy families 都被拒绝。
- Bound identity 只能满足显式分类为 bound identity 的 routes。
- Session-validated identity 只能满足显式分类为 session validation 的 routes。
- 当 route 同时要求 bound 和 session identity 时，actor/player/session/connection 字段必须一致。
- 不一致必须在 domain dispatch 前 fail closed。
- Policy failures 保持脱敏。
- WebSocket transport 保持 credential-neutral。
- Protobuf envelope 和 authentication response shape 保持不变。
- 除非显式重新分类，inventory route behavior 不变。

除非后续 implementation work item 要求真实持久化，live PostgreSQL verification 可以继续保持 opt-in。

## 11. Nakama 和 Pitaya 参考映射

Nakama 参考映射：

```yaml
adopted_concepts:
  - authenticated_session_material_controls_gameplay_access
  - session_lifetime_and_route_authorization_are_related
  - logout_refresh_and_session_management_have_future_policy_pressure
adapted_concepts:
  - vibit_keeps_request_token_bound_connection_and_durable_session_identity_as_explicit_policy_families
  - vibit_does_not_copy_nakama_session_api_or_jwt_shape
  - vibit_requires_route_scoped_policy_before_trusting_bound_or_session_identity
deferred_concepts:
  - refresh_token_flow
  - session_management_api
  - logout_disconnect_active_socket
  - single_session_or_single_socket_policy
rejected_for_now:
  - direct_nakama_api_compatibility
```

Pitaya 参考映射：

```yaml
adopted_concepts:
  - acceptor_transport_and_handler_logic_are_separate
  - sessions_can_bind_user_identity_to_connection_context
  - handlers_should_receive_context_not_parse_credentials
adapted_concepts:
  - vibit_route_policy_is_application_owned_and_contract_first
  - vibit_bound_identity_is_not_implicitly_all_routes_authenticated
  - vibit_keeps_cluster_session_routing_deferred
deferred_concepts:
  - frontend_backend_cluster_session_routing
  - remote_session_binding_broadcast
  - groups_rooms_and_presence_attachment
  - server_to_server_rpc
rejected_for_now:
  - direct_pitaya_api_compatibility
```

Nakama 和 Pitaya 指导 capability planning。它们不覆盖 vibit 的 constitution、ADR、manifests、generated boundaries 或 verification commands。

## 12. 未来实现队列

本 gate 之后，未来工作仍应拆分：

```yaml
future_work_items:
  bound_identity_route_policy_implementation:
    may_add:
      - explicit application route policy families
      - route-scoped classification for public, request-token, bound-connection, and session-validated policies
      - tests for fail-closed identity agreement
  logout_revocation_active_connection_gate:
    may_define:
      - whether token/session revocation closes active WebSocket connections
  reconnect_connection_epoch_gate:
    may_define:
      - reconnect, resume, duplicate replacement, and epoch mismatch behavior
  protocol_session_carrier_gate:
    may_define:
      - whether and how clients receive or carry session ids
```

没有新的 ADR 时，不要把这些合成一个宽泛的 session subsystem slice。

## 13. 验证

本 gate 的仓库检查规则是：

```text
runtime.bound_identity_route_policy_gate
```
