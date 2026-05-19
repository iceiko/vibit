# First Message Connection Binding Gate

状态：Draft v0.1
最后更新：2026-05-17
范围：Session persistence 与 WebSocket handshake ratification 之后，future protocol/application connection binding posture
依赖：`docs/session-persistence-websocket-handshake-ratification.md`、`docs/access-token-protocol-carrier-route-protection-gate.md`、`docs/authentication-command-protocol-login-route-gate.md`、`docs/runtime-protocol-adapter.md`、`docs/reference-game-server-alignment.md`
权威决策：`ADR-0057`

配套英文原文是 `docs/first-message-connection-binding-gate.md`。英文文件是权威版本。

## 1. 目的

vibit 现在已经有：

- Public device credential login route：`runtime.authentication.AuthenticateWithDeviceCredential`。
- 通过 `vibit.authentication.v1.AuthenticatedRequest` 实现 request-level opaque access-token validation。
- Authentication service 和 route protection 的 PostgreSQL startup composition。
- `ADR-0056` 中 ratified session 与 WebSocket handshake posture。

下一个 connection-lifecycle 问题是：已经 authenticated 的 player 如何关联到一个打开的 WebSocket connection，同时不把 credential parsing 放进 WebSocket transport，也不修改现有 Protobuf envelope。

本 gate 定义 future first-message connection binding posture：

```yaml
first_message_connection_binding_gate: defined
implementation_authorized_by_this_standard: false
completed_gate_work_item: W-0124
decision: ADR-0057
check_rule: runtime.first_message_connection_binding_gate
selected_binding_message_kind: system
selected_binding_route: runtime.authentication.BindConnection
selected_binding_payload_candidate: vibit.authentication.v1.BindConnectionRequest
selected_binding_response_candidate: vibit.authentication.v1.BindConnectionResponse
selected_proof_carrier: protobuf_system_payload_access_token
websocket_transport_credential_neutral: true
websocket_handshake_authentication_added: false
existing_protobuf_envelope_change_added: false
protocol_source_added: false
generated_output_added: false
connection_binding_registry_added: false
route_policy_bound_identity_added: false
session_persistence_added: false
```

这是 gate-only standard，不实现 connection binding。

## 2. 已选择的 Future Shape

已选择的 future route 是：

```yaml
route:
  kind: system
  module: runtime.authentication
  name: BindConnection
semantic_route_key: runtime.authentication.BindConnection
payload_candidate: vibit.authentication.v1.BindConnectionRequest
response_candidate: vibit.authentication.v1.BindConnectionResponse
```

这条 message 是 protocol/application binding message。它不是 domain command，也不得通过 inventory、player 或其他 gameplay modules 路由。

计划中的 request shape：

```yaml
BindConnectionRequest:
  access_token: opaque Base64URL unpadded 32-byte access-token proof
  client_instance_id: optional client installation or runtime instance hint
```

计划中的 response shape：

```yaml
BindConnectionResponse:
  binding_status: bound | rejected
  actor_kind: player
  player_id: validated player id when bound
  connection_id: server-assigned connection id
  connection_epoch: server-assigned or server-observed connection epoch
  session_validated: false until durable session persistence is implemented
  bound_at: server timestamp when binding succeeds
```

规则：

- Access token 仍然是 opaque proof，必须由 existing application authentication service boundary 验证。
- WebSocket transport 不得从 HTTP headers、cookies、query strings、bearer strings 或 subprotocols 解析 access token。
- First-message binding 不得修改现有 `proto/vibit/protocol/v1/envelope.proto` shape。
- Future bind route 必须在 normal domain dispatch 前被处理。
- Binding 成功必须生成 normalized application identity，而不是信任 metadata-only envelope fields。

## 3. Future Layer Ownership

Future implementation 必须保持这个 layer split：

```yaml
websocket_transport:
  owns:
    - connection accept
    - server-assigned connection id
    - binary frame read and write
  must_not_own:
    - access-token parsing
    - authentication service calls
    - player account lookup
    - request identity construction

protobuf_protocol_adapter:
  owns:
    - existing envelope decode
    - recognition of runtime.authentication.BindConnection system route
    - bind request payload decode
    - application binder handoff
    - bind response or public error envelope encode
  must_not_own:
    - token verifier lookup
    - credential lookup
    - durable session storage
    - logout or reconnect policy

application_runtime:
  owns:
    - connection binding policy
    - access-token validation handoff
    - normalized identity creation
    - connection-bound identity registry if implemented
    - route policy use of bound identity if implemented
  must_not_own:
    - WebSocket handshake credential extraction
    - direct Nakama or Pitaya API compatibility
```

Candidate future files：

```yaml
protocol_source: proto/vibit/authentication/v1/authentication.proto
generated_go_output: runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go
protocol_adapter_source: runtime/internal/platform/protocol/protobuf/connection_binding*.go
application_binding_source: runtime/internal/app/connection_binding*.go
startup_composition_owner: runtime/cmd/vibit-server
```

Generated Go Protobuf output 必须由 Buf 生成，不能手工编辑。

## 4. Binding Flow

计划中的 future flow：

```yaml
connection_binding_flow:
  - websocket_transport_accepts_connection_without_reading_credentials
  - websocket_transport_assigns_connection_id
  - client_may_call_public_login_route_if_it_needs_an_access_token
  - client_sends_runtime.authentication.BindConnection_system_message
  - protobuf_adapter_decodes_bind_payload
  - application_binding_policy_validates_access_token
  - authentication_service_returns_validated_player_identity
  - application_runtime_binds_identity_to_connection_id
  - bind_response_confirms_bound_identity_without_returning_secret_material
```

Anonymous pre-bind messages：

```yaml
allowed_before_binding:
  - runtime.authentication.AuthenticateWithDeviceCredential
  - runtime.authentication.BindConnection
  - heartbeat_or_ack_frames_after_a_later_protocol_rule_allows_them
forbidden_before_binding:
  - bound_identity_only_gameplay_routes
  - future_presence_room_party_match_membership_routes
  - future_group_broadcast_membership_routes
```

规则：

- Public login 必须仍可在 connection binding 前使用，因为新 client 可能还没有 access token。
- First-message binding 指连接上的第一条 identity-binding message，不必是该连接发送的第一帧。
- Metadata-only envelope session fields 永远不能满足 binding。
- 成功 bind 必须使用 server-observed connection id，而不是 client-supplied connection id。

## 5. 与 Route Policy 的关系

当前 route-protection behavior 仍然是通过 `vibit.authentication.v1.AuthenticatedRequest` 做 request-level access-token validation。

Future binding implementation 只有在显式更新 application route policy 时，才能把 connection-bound identity 作为 proof source。计划顺序是：

```yaml
future_protected_route_identity_sources:
  first: request_level_access_token_wrapper
  second: bound_connection_identity_after_explicit_route_policy_update
```

规则：

- 现有 protected routes 不得静默变成 public。
- Bound connection identity 仍然必须表现为 normalized `RequestIdentity`。
- 在单独的 durable session persistence implementation 验证 session 之前，`RequestIdentity.SessionValidated` 保持 false。
- Domain handlers 不得直接查询 connection binding registries。
- Inventory、player、authentication repository interfaces 和 PostgreSQL adapters 不得被这个 gate 修改。

## 6. Failure Behavior

Future implementation 必须折叠 public failures：

```yaml
future_public_errors:
  missing_bind_proof: CONNECTION_BINDING_TOKEN_MISSING
  malformed_bind_payload_or_token: CONNECTION_BINDING_TOKEN_MALFORMED
  invalid_or_expired_or_revoked_token: CONNECTION_BINDING_TOKEN_INVALID
  validation_dependency_unavailable: CONNECTION_BINDING_UNAVAILABLE
  protected_bound_route_without_binding: CONNECTION_BINDING_REQUIRED
```

规则：

- Public errors 不得泄露 token lookup hit/miss、token lifecycle status、verifier key id、player account state、verifier mismatch 或 internal binding registry state。
- Invalid binding proof 不得把任何 identity bind 到 connection。
- Rebinding already bound connection、duplicate player connections、active-connection invalidation 和 kick/replace behavior 需要后续 gates。
- 本 gate 不选择 bind timeout，也不选择 bind failure 后自动关闭连接。后续 implementation 可以返回 protocol error 而不关闭 WebSocket。
- Repeated failures 的 transport-level close policy 属于 operations 或 abuse-control concern，继续 deferred。

## 7. Session Persistence 与 Reconnect Deferrals

First-message connection binding 不是 durable session persistence。

规则：

- 不授权 `runtime_sessions` table。
- 不授权 session repository。
- 不授权 session migration。
- 不选择 Redis-like 或 external session-store dependency。
- 不授权 reconnect resume behavior。
- 不授权 connection epoch replacement behavior。
- 不授权 logout-triggered active connection invalidation。
- 不授权 token rotation 或 refresh behavior。

Future durable sessions 需要 PostgreSQL-first schema/repository/migration gate。Future reconnect 和 duplicate connection behavior 需要 reconnect/epoch gate。

## 8. Nakama 与 Pitaya Reference Mapping

Nakama reference mapping：

```yaml
adopted_concepts:
  - authenticate_before_realtime_socket_features
  - socket_connection_associated_with_authenticated_state
  - session_and_socket_lifecycle_are_related_but_not_identical
adapted_concepts:
  - vibit_uses_opaque_access_token_proof_not_nakama_jwt_session_compatibility
  - vibit_binds_connection_through_protocol_application_message_not_websocket_handshake
deferred_concepts:
  - refresh_token_model
  - single_socket_single_session_policy
  - server_side_session_disconnect
  - reconnect_restore_behavior
rejected_for_now:
  - direct_nakama_api_compatibility
```

Pitaya reference mapping：

```yaml
adopted_concepts:
  - acceptor_session_route_handler_separation
  - bind_user_identity_to_session_like_connection_context
  - groups_rooms_and_broadcast_need_bound_identity_later
adapted_concepts:
  - vibit_binding_is_application_owned_and_protocol_explicit
  - vibit_single_process_registry_precedes_cluster_session_broadcast
deferred_concepts:
  - frontend_backend_split
  - server_to_server_rpc
  - remote_session_bind_broadcast
  - direct_group_membership_api
rejected_for_now:
  - direct_pitaya_api_compatibility
```

Nakama 和 Pitaya 指导 capability planning。它们不覆盖 vibit 的 constitution、manifests、ADRs、generated boundaries 或 verification commands。

## 9. 必需的 Future Gates

在 first-message binding behavior 存在之前，后续 implementation gate 必须定义：

- Exact Protobuf message fields 和 field numbers。
- 是扩展 `authentication.proto` 还是创建新的 authentication protocol source。
- Generated output update path。
- Protocol adapter recognition 和 response mapping。
- Application connection binding registry shape。
- Binder 的 startup composition。
- Route policy 如何使用 bound identity。
- Missing、malformed、invalid、expired、revoked 和 unavailable proof 的 focused tests。
- Raw tokens、token digests、key ids、connection ids 和 internal binding state 的 redaction rules。
- Failed binding attempts 的 close/error behavior。

在 durable 或 multi-connection behavior 存在之前，单独的 future gates 必须定义：

- PostgreSQL session schema。
- Logout、revocation 和 active-connection invalidation。
- Reconnect 与 connection epoch behavior。
- Single-socket、duplicate connection、kick 和 replacement behavior。
- Presence、room、party、match 和 group membership attachment。
- Operations、metrics 和 abuse-control posture。

## 10. Deferrals

本 gate 不授权：

- 新 `.proto` messages。
- Generated Protobuf output。
- 修改 `proto/vibit/protocol/v1/envelope.proto`。
- WebSocket handshake authentication。
- HTTP `Authorization`、Bearer、cookie、query-string 或 subprotocol proof carriers。
- Connection-bound identity registry code。
- Route policy 使用 bound identity。
- Session persistence。
- Session tables、repositories、PostgreSQL adapters、migrations 或 cleanup jobs。
- Logout、refresh、token rotation、token validation audit mutation 或 active-connection invalidation。
- Reconnect、resume、connection epoch replacement 或 duplicate-connection policy。
- 新 dependencies。
- Memory durable authentication behavior。
- Direct Nakama 或 Pitaya public API compatibility。

## 11. Verification

本 gate 的 repository check rule 是：

```text
runtime.first_message_connection_binding_gate
```

该 check 必须验证：

- Standard、translation、ADR、change specs 和 conversation log 存在。
- Runtime、conventions、contracts、reference、work-item、module 和 AGENTS markers 存在。
- WebSocket transport non-test Go files 保持 credential-neutral。
- 现有 Protobuf envelope 不包含 authentication proof 或 binding fields。
- Implementation gate 之前没有新增 `BindConnection` Protobuf source 或 generated output。
- 不存在 session 或 connection-binding migration source。

## 12. References

- Nakama sockets：`https://heroiclabs.com/docs/nakama/concepts/sockets/`
- Nakama sessions：`https://heroiclabs.com/docs/nakama/concepts/session/`
- Nakama configuration session limits：`https://heroiclabs.com/docs/nakama/getting-started/configuration/`
- Pitaya session package documentation：`https://pkg.go.dev/github.com/topfreegames/pitaya/v3/pkg/session`
