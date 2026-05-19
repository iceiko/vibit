# Session Persistence 与 WebSocket Handshake Ratification

状态：Draft v0.1
最后更新：2026-05-17
范围：public login route 暴露之后，request-level validation、future session persistence、future connection binding、WebSocket handshake authentication deferral，以及 session-related implementation gates 的第一版 ratified posture
依赖：`docs/session-persistence-websocket-handshake-decision-gates.md`、`docs/access-token-protocol-carrier-route-protection-gate.md`、`docs/authentication-command-protocol-login-route-gate.md`、`docs/runtime-protocol-adapter.md`、`docs/reference-game-server-alignment.md`
权威决策：`ADR-0056`

配套英文原文是 `docs/session-persistence-websocket-handshake-ratification.md`。英文文件是权威版本。

## 1. 目的

vibit 现在已经有最小 authenticate-then-gameplay loop：

- Public device credential login route：`runtime.authentication.AuthenticateWithDeviceCredential`。
- Application authentication service 中的 opaque access-token validation。
- 通过 `vibit.authentication.v1.AuthenticatedRequest` 实现 request-level route protection。
- PostgreSQL startup composition 中组合 authentication service 和 route protection。

下一个 game-server concern 是 session 与 connection lifecycle。Nakama 这类成熟服务器提供 authentication/session 模型，并把 realtime socket 与 authenticated client state 关联起来。Pitaya 则把 acceptors、sessions、route handlers、groups 和后续 cluster roles 分开。vibit 应该学习这些系统，但仍保留自己的 agent-native boundaries。

本标准 ratify public login 之后的第一版 posture：

```yaml
session_persistence_websocket_handshake_ratification: defined
current_validation_model: request_level_access_token_validation
current_proof_carrier: protobuf_authenticated_request_payload_wrapper
current_websocket_handshake_authentication: not_selected
current_websocket_transport_credential_neutral: true
current_protobuf_envelope_change: unchanged
future_connection_binding_preferred_gate: first_message_protocol_binding
future_session_store_preferred_first_durable_target: postgres
implementation_authorized_by_this_standard: false
completed_work_item: W-0122
decision: ADR-0056
check_rule: runtime.session_persistence_websocket_handshake_ratification
```

这是 ratification gate，不是 implementation gate。

## 2. 已选择的当前路径

当前 production-sensitive path 仍然是 request-level validation：

```yaml
login:
  route: runtime.authentication.AuthenticateWithDeviceCredential
  result: opaque_access_token
protected_request:
  carrier: vibit.authentication.v1.AuthenticatedRequest
  validation_owner: runtime/internal/app/authentication
  route_policy_owner: runtime/internal/app
  domain_handoff: normalized RequestIdentity
  session_validated: false
websocket_transport:
  credential_neutral: true
protobuf_envelope:
  changed: false
```

规则：

- 在后续 session persistence 或 binding implementation 明确验证 session 之前，`RequestIdentity.SessionValidated` 继续保持 false。
- Client-supplied envelope `Session.player_id`、`Session.session_id`、`Session.connection_id` 和 `Session.connection_epoch` 仍然只是 metadata。
- Protected domain routes 必须依赖 route protection 生成的 normalized request identity，而不是 metadata-only envelope fields。
- 当前 proof carrier 仍然是 Protobuf authenticated request payload wrapper，而不是 WebSocket handshake metadata。

## 3. WebSocket Handshake Authentication 的 Ratified Deferral

WebSocket handshake authentication 不被选择为下一条 implementation path。

WebSocket transport 必须继续忽略 credential carriers：

```yaml
forbidden_current_transport_carriers:
  - HTTP Authorization header
  - Bearer header value
  - Cookie access_token or session token
  - query-string access_token or session token
  - Sec-WebSocket-Protocol authentication token
```

规则：

- WebSocket transport 不得拥有 player account lookup。
- WebSocket transport 不得拥有 credential lookup。
- WebSocket transport 不得拥有 token verifier lookup。
- WebSocket transport 不得创建 authenticated `RequestIdentity`。
- 未来任何 transport-level extraction 都需要单独的 transport-auth boundary，并且必须 hand off 给 application-owned 或 authentication-owned validation contracts。

理由：

- Browser 与 non-browser clients 对 headers、cookies、query strings、subprotocols 的实际约束不同。
- 在 Protobuf envelope 存在前拒绝连接，会改变 error behavior 和 compatibility。
- Transport-owned authentication 会让未来 non-WebSocket transports 更难保持一致。

## 4. Future Connection Binding Preference

优先的下一条 connection-level gate 是 first-message protocol binding：

```yaml
future_connection_binding_preferred_gate: first_message_protocol_binding
candidate_owner_layers:
  protocol_adapter: decode system/authentication binding message
  application_runtime: validate proof and bind normalized identity to connection context
  websocket_transport: keep connection plumbing credential-neutral
```

这个 preference 不实现 first-message binding。后续 gate 必须定义：

- System 或 authentication message name。
- Protobuf payload shape。
- Timeout 和 failure behavior。
- Anonymous pre-bind messages 是否允许。
- Connection context storage owner。
- Revalidation behavior。
- Logout、revocation、expiration 和 reconnect behavior。
- Route protection 如何使用 bound identity 或 per-request proof。
- Duplicate active connections 如何处理。
- 未来 room、party、match、presence 和 group membership 如何挂到 binding 上。

规则：

- First-message binding 不得被发明成 ad hoc domain command。
- Binding state 不得放在 inventory、player 或其他 domain modules 中。
- Binding identity 仍然必须通过 normalized application request identity 进入 domain dispatch。

## 5. Future Session Persistence Posture

未来 session persistence 的第一版 durable target 优先选择 PostgreSQL，因为 PostgreSQL 已经是 module state 的第一版 accepted authoritative durable store。

这个 preference 不添加 session persistence。

未来 session persistence 需要单独的 schema gate，定义：

- Sessions 是否持久化。
- Session record ownership。
- Session id generation。
- Actor kind 和 actor id fields。
- Player id linkage。
- Issue、expiration、revocation 和 last-seen semantics。
- Connection binding fields，如需要。
- Token record linkage，如需要。
- Logout 和 revocation transaction boundaries。
- Cleanup strategy。
- Migration source。
- Repository interface 和 PostgreSQL adapter boundaries。
- Live verification requirements。

规则：

- 本标准不授权 `runtime_sessions` table。
- 本标准不授权 repository method。
- 本标准不授权 PostgreSQL adapter behavior。
- 本标准不选择 Redis-like 或 external session-store dependency。
- Memory session storage 不被选择为 production behavior。

## 6. Nakama 与 Pitaya Reference Mapping

Nakama reference mapping：

```yaml
adopted_concepts:
  - authenticate_before_protected_gameplay
  - session_lifetime_dimensions
  - realtime_socket_associated_with_authenticated_state
adapted_concepts:
  - token_or_session_material_is_vibit_opaque_access_token_first
  - realtime_socket_binding_is_deferred_to_protocol_application_gate
deferred_concepts:
  - refresh_token_behavior
  - session_persistence_schema
  - socket_reconnect_restore_behavior
rejected_for_now:
  - direct_nakama_api_compatibility
```

Pitaya reference mapping：

```yaml
adopted_concepts:
  - transport_acceptor_separated_from_session_binding
  - route_handler_separated_from_connection_plumbing
  - groups_rooms_broadcast_require_later_binding_context
adapted_concepts:
  - session_binding_becomes_application_protocol_first_message_gate
deferred_concepts:
  - frontend_backend_split
  - server_to_server_rpc
  - cluster_service_discovery
rejected_for_now:
  - direct_pitaya_api_compatibility
```

Nakama 和 Pitaya 指导 capability planning。它们不覆盖 vibit 的 constitution、manifests、ADRs、generated boundaries 或 verification commands。

## 7. 必需的 Future Gates

实现 session 或 connection behavior 之前，未来工作必须创建这些 gate 之一：

```yaml
future_gates:
  first_message_connection_binding_gate:
    required_before:
      - system/authentication binding protobuf messages
      - connection-bound identity cache
      - route policy use of bound identity
      - reconnect or epoch behavior
  postgres_session_persistence_schema_gate:
    required_before:
      - session tables
      - session repository interfaces
      - PostgreSQL session adapters
      - session cleanup jobs
  websocket_handshake_authentication_gate:
    required_before:
      - Authorization/Bearer/cookie/query/subprotocol credential parsing
      - pre-envelope handshake rejection behavior
      - transport-level proof extraction
  logout_revocation_active_connection_gate:
    required_before:
      - logout execution
      - active session revocation
      - closing or invalidating bound WebSocket connections
  reconnect_epoch_gate:
    required_before:
      - connection_epoch semantics
      - duplicate connection replacement
      - session resume behavior
```

## 8. Deferrals

本标准不授权：

- Session persistence implementation。
- Session schema、migrations、repositories 或 PostgreSQL adapters。
- WebSocket handshake authentication。
- HTTP `Authorization`、Bearer、cookie、query-string 或 WebSocket subprotocol proof carriers。
- Existing Protobuf envelope changes。
- New Protobuf system messages。
- Generated Protobuf output。
- Connection-bound identity cache。
- Reconnect、resume 或 connection epoch behavior。
- Logout、refresh、cleanup、token rotation 或 token validation audit mutation。
- New dependencies。
- Memory durable authentication behavior。
- Direct Nakama 或 Pitaya public API compatibility。

## 9. Verification

该 ratification 的 repository check rule 是：

```text
runtime.session_persistence_websocket_handshake_ratification
```

该检查必须验证：

- Standard、translation、ADR、change specs 和 conversation log 存在。
- Manifests 和 AGENTS guides 记录 ratified posture。
- WebSocket transport 仍然 credential-neutral。
- Existing Protobuf envelope 没有 token、credential、login 或 handshake proof fields。
- 本 gate 不存在 session migration source。
- 本 gate 不要求新的 generated Protobuf output。
