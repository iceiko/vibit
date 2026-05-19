# First Message Connection Binding Implementation Gate

状态：Draft v0.1
最后更新：2026-05-17
范围：`runtime.authentication.BindConnection` 的 future bounded implementation slice
依赖：`docs/first-message-connection-binding-gate.md`、`docs/session-persistence-websocket-handshake-ratification.md`、`docs/access-token-protocol-carrier-route-protection-gate.md`、`docs/authentication-command-protocol-login-route-gate.md`、`docs/runtime-authentication-startup-composition-gate.md`、`docs/runtime-protocol-adapter.md`、`docs/generated-output.md`、`docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0058`

英文版 `docs/first-message-connection-binding-implementation-gate.md` 是权威版本。本文件是简体中文翻译。

## 1. 目的

`ADR-0057` 已选择 first-message protocol/application binding 作为 future connection-level identity posture。选定 route 是：

```text
runtime.authentication.BindConnection
```

本 gate 定义后续 bounded implementation slice 可以添加什么，以及仍然不能添加什么。

Implementation gate 是：

```yaml
first_message_connection_binding_implementation_gate: defined
implementation_authorized_by_this_standard: false
completed_gate_work_item: W-0126
future_implementation_work_item: W-0128
decision: ADR-0058
check_rule: runtime.first_message_connection_binding_implementation_gate
future_route: runtime.authentication.BindConnection
future_route_kind: system
future_protocol_source: proto/vibit/authentication/v1/authentication.proto
future_generated_go_output: runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go
future_protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
future_application_binding_owner: runtime/internal/app
future_startup_owner: runtime/cmd/vibit-server
future_request: vibit.authentication.v1.BindConnectionRequest
future_response: vibit.authentication.v1.BindConnectionResponse
future_status_enum: vibit.authentication.v1.ConnectionBindingStatus
first_composed_runtime_store: postgres
memory_store_binding_status: unavailable_bootstrap
websocket_transport_credential_neutral: true
websocket_handshake_authentication_added: false
existing_protobuf_envelope_change_added: false
session_persistence_added: false
route_policy_bound_identity_added_by_this_gate: false
```

这是 gate-only standard。本文件不实现 connection binding。

## 2. Future Implementation Boundary

后续 implementation slice 只可以添加这些行为族：

```yaml
future_allowed_implementation:
  protocol_source:
    - extend proto/vibit/authentication/v1/authentication.proto with BindConnection messages
  generated_output:
    - regenerate runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go through Buf
  protocol_adapter:
    - recognize runtime.authentication.BindConnection as a system route
    - decode BindConnectionRequest payload
    - map public binding result to BindConnectionResponse
    - map public binding errors to error envelopes
  application_binding:
    - validate access-token proof through existing application authentication service boundary
    - create normalized RequestIdentity for the server-observed connection id
    - keep an in-memory process-local connection binding registry if explicitly implemented in the slice
  startup_composition:
    - wire the binder only when the PostgreSQL authentication service is composed
  tests:
    - protocol adapter tests
    - application binder tests
    - startup composition tests
    - WebSocket transport neutrality regression tests
```

Future implementation 不得把 authentication 移进 transport、generated code、repositories、domain modules 或 Protobuf envelope metadata。

## 3. Future Protocol Shape

Future implementation 可以扩展：

```text
proto/vibit/authentication/v1/authentication.proto
```

计划 messages：

```yaml
BindConnectionRequest:
  fields:
    access_token:
      type: string
      secret: true
      semantics: opaque Base64URL unpadded 32-byte access-token proof
    client_instance_id:
      type: string
      secret: false
      semantics: optional client installation or runtime instance hint

BindConnectionResponse:
  fields:
    binding_status:
      type: ConnectionBindingStatus
      semantics: bound or rejected
    actor_kind:
      type: string
      semantics: player when bound
    player_id:
      type: string
      semantics: validated player id when bound
    connection_id:
      type: string
      semantics: server-observed connection id
    connection_epoch:
      type: uint64
      semantics: server-observed connection epoch
    session_validated:
      type: bool
      semantics: false until durable session persistence is implemented
    bound_at:
      type: string
      semantics: server timestamp in RFC3339Nano UTC text

ConnectionBindingStatus:
  values:
    - CONNECTION_BINDING_STATUS_UNSPECIFIED
    - CONNECTION_BINDING_STATUS_BOUND
    - CONNECTION_BINDING_STATUS_REJECTED
```

规则：

- Route kind 是 `system`。
- Route key 是 `runtime.authentication.BindConnection`。
- Request payload type 是 `vibit.authentication.v1.BindConnectionRequest`。
- Response payload type 是 `vibit.authentication.v1.BindConnectionResponse`。
- 现有 `proto/vibit/protocol/v1/envelope.proto` 必须保持不变。
- Generated Go Protobuf output 必须由 Buf 生成，不能手工编辑。
- `access_token` 是 secret material，不得出现在 errors、logs、events、debug details 或 repository records 中。

## 4. Future Binding Flow

Future implementation 必须保持这个 flow：

```yaml
bind_connection_flow:
  - websocket_transport_accepts_connection_without_reading_credentials
  - websocket_transport_assigns_or_observes_connection_id
  - client_may_call_public_login_route_if_access_token_is_needed
  - client_sends_runtime.authentication.BindConnection_system_message
  - protobuf_adapter_decodes_existing_envelope_and_bind_payload
  - application_binding_service_validates_access_token_through_authentication_service
  - application_binding_service_builds_validated_player_identity_for_connection
  - application_binding_service_records_process_local_binding_if_registry_is_in_slice
  - protocol_adapter_returns_bind_response_without_secret_material
```

规则：

- Binding 通过现有 application authentication service path 验证 access-token proof。
- Binding 必须使用 server-observed `connection_id` 和 `connection_epoch`，不得信任 client-supplied connection identity。
- Binding 成功可以创建 connection-bound identity，但 route policy 只有在 implementation slice 明确更新时才可以使用这个 identity。
- `RequestIdentity.SessionValidated` 在 durable session persistence 单独实现之前必须保持 false。
- Domain handlers 不得直接查询 connection binding registries。

## 5. Future Application Ownership

Future application boundary 应使用 `runtime/internal/app` 下的 app-owned types。

Candidate files：

```yaml
application_binding_source:
  - runtime/internal/app/connection_binding.go
  - runtime/internal/app/connection_binding_test.go
optional_authentication_adapter_source:
  - runtime/internal/app/authentication/connection_binding_validator.go
  - runtime/internal/app/authentication/connection_binding_validator_test.go
```

Candidate application vocabulary：

```yaml
BindConnectionRequest:
  access_token: string
  route: RouteKey
  connection_id: string
  connection_epoch: uint64
  client_instance_id: string

BindConnectionResult:
  bound: bool
  identity: RequestIdentity
  binding_status: bound | rejected
  public_error_code: ErrorCode
  connection_id: string
  connection_epoch: uint64
  bound_at: time.Time
```

规则：

- Binder 可以依赖一个验证 access tokens 的接口；它不得 import platform persistence 或 protocol packages。
- Binder 必须把 public token failures 折叠成 binding-specific public error codes。
- Binder 不得生成 token、撤销 token、刷新 token、mutate token audit state 或创建 session records。
- 如果添加 registry，第一版必须只是 process-local，不能宣称 durability。

## 6. Future Protocol Adapter Ownership

Future Protobuf adapter boundary 应使用：

```text
runtime/internal/platform/protocol/protobuf
```

Candidate files：

```yaml
protocol_adapter_source:
  - runtime/internal/platform/protocol/protobuf/connection_binding.go
  - runtime/internal/platform/protocol/protobuf/connection_binding_test.go
```

规则：

- Protocol adapter 可以 decode `BindConnectionRequest`。
- Protocol adapter 可以把 binding results 映射到 `BindConnectionResponse`。
- Protocol adapter 可以把 public binding errors 映射到现有 error envelopes。
- Protocol adapter 不得验证 token digests、查询 repositories、加载 verifier keys、创建 sessions 或决定 token lifecycle state。
- Adapter 不得改变正常 request-level `vibit.authentication.v1.AuthenticatedRequest` route protection behavior。

## 7. Startup Composition

Future startup composition 限制为：

```yaml
first_runtime_store: postgres
startup_owner: runtime/cmd/vibit-server
memory_store_binding_status: unavailable_bootstrap
```

规则：

- PostgreSQL runtime startup 只有在现有 authentication service 和 route protector 已组合时，才可以 wire connection binding。
- Memory runtime startup 可以让 `BindConnection` 保持 unavailable。
- Startup 不得从 process arguments、environment variables、HTTP headers、cookies、query strings 或 WebSocket subprotocols 解析 access tokens。
- Startup 不得自动 apply migrations。
- Startup 不得注册 durable session stores，除非后续 session persistence gate 授权。

## 8. Error Mapping

Future implementation 必须使用 binding-specific public errors：

```yaml
future_public_errors:
  missing_bind_proof: CONNECTION_BINDING_TOKEN_MISSING
  malformed_bind_payload_or_token: CONNECTION_BINDING_TOKEN_MALFORMED
  invalid_or_expired_or_revoked_token: CONNECTION_BINDING_TOKEN_INVALID
  validation_dependency_unavailable: CONNECTION_BINDING_UNAVAILABLE
  protected_bound_route_without_binding: CONNECTION_BINDING_REQUIRED
```

规则：

- Public errors 不得泄露 lookup hit/miss、verifier key id、verifier mismatch、token status、player account state、audience mismatch、connection registry state 或 internal dependency class。
- Failed bind 不得把 identity 绑定到 connection。
- Failed bind 不得返回 token record ids、credential record ids、lookup digests、verifier digests 或 raw proof。
- 已绑定连接的 rebinding、duplicate player connections、kick/replace behavior 和 repeated failure close policy 仍然需要独立 gate。

## 9. Nakama And Pitaya Reference Mapping

Nakama reference mapping：

```yaml
adopted_concepts:
  - authenticated_session_material_precedes_authenticated_realtime_socket_use
  - realtime_socket_lifecycle_can_have_authenticated_player_context
  - session_token_lifetime_and_socket_connection_lifetime_are_related_but_distinct
adapted_concepts:
  - vibit_uses_opaque_access_token_not_nakama_session_api_compatibility
  - vibit_binds_connection_through_existing_protobuf_system_route_not_transport_handshake
  - vibit_keeps_request_level_validation_as_current_protected_route_path
deferred_concepts:
  - refresh_token_flow
  - single_session_or_single_socket_enforcement
  - disconnect_on_session_revocation
  - reconnect_restore_behavior
rejected_for_now:
  - direct_nakama_api_compatibility
```

Pitaya reference mapping：

```yaml
adopted_concepts:
  - transport_acceptor_separated_from_session_binding
  - user_identity_can_bind_to_session_like_connection_context
  - route_handlers_should_receive_context_not_parse_transport_credentials
adapted_concepts:
  - vibit_connection_binding_is_application_owned_and_protocol_explicit
  - vibit_first_registry_is_process_local_before_cluster_session_broadcast
  - vibit_routes_remain_contract_first_instead_of_pitaya_route_string_compatibility
deferred_concepts:
  - frontend_backend_cluster_split
  - remote_session_binding_broadcast
  - groups_rooms_and_presence_attachment
  - server_to_server_rpc
rejected_for_now:
  - direct_pitaya_api_compatibility
```

Nakama 和 Pitaya 用于指导 capability planning。它们不覆盖 vibit 的 constitution、manifests、ADRs、generated boundaries 或 verification commands。

## 10. Required Future Tests

Future implementation slice 必须添加 focused tests：

```yaml
required_tests:
  protocol_source_includes_bind_connection_messages
  generated_output_is_buf_generated_and_traced_to_proto_source
  bind_connection_payload_decodes_to_application_request
  bind_connection_success_maps_to_response_without_secret_material
  bind_connection_missing_token_maps_to_CONNECTION_BINDING_TOKEN_MISSING
  bind_connection_malformed_token_maps_to_CONNECTION_BINDING_TOKEN_MALFORMED
  bind_connection_invalid_token_maps_to_CONNECTION_BINDING_TOKEN_INVALID
  bind_connection_unavailable_dependency_maps_to_CONNECTION_BINDING_UNAVAILABLE
  bind_connection_uses_server_observed_connection_id
  bind_connection_keeps_session_validated_false
  failed_bind_does_not_create_bound_identity
  public_login_route_remains_available_before_binding
  request_level_authenticated_wrapper_remains_current_protected_route_path
  websocket_transport_remains_credential_neutral
  existing_protobuf_envelope_remains_unchanged
```

Live PostgreSQL verification 仍然是 optional，除非后续 change 明确要求 opt-in live check。

## 11. Deferrals

本 gate 不授权：

- 实现 `BindConnection`。
- 添加 `.proto` messages。
- 运行 `buf generate`。
- 编辑 generated `.pb.go` files。
- 修改 `proto/vibit/protocol/v1/envelope.proto`。
- WebSocket handshake authentication。
- HTTP `Authorization`、Bearer、cookie、query-string 或 subprotocol proof carriers。
- Durable session persistence。
- Session tables、session repositories、session migrations 或 session cleanup jobs。
- Route-policy use of bound identity，除非 future implementation slice 明确修改 route policy。
- Logout-triggered active-connection invalidation。
- Refresh behavior、token rotation 或 token validation audit mutation。
- Reconnect、resume、duplicate connection replacement 或 connection epoch policy。
- Presence、rooms、parties、match runtime、group membership 或 broadcast behavior。
- 新 dependencies。
- Memory durable authentication behavior。
- Direct Nakama 或 Pitaya public API compatibility。

## 12. Verification

本 gate 的 repository check rule 是：

```text
runtime.first_message_connection_binding_implementation_gate
```

检查必须验证：

- Standard、translation、ADR、change specs 和 conversation log 存在。
- Runtime、conventions、contracts、reference、work-item、module 和 AGENTS markers 存在。
- Gate 声明 future implementation artifacts 和 deferrals。
- WebSocket transport non-test Go files 保持 credential-neutral。
- 现有 Protobuf envelope 不含 authentication proof 和 binding fields。
- Gate-only change 没有添加 `BindConnection` Protobuf source 或 generated output。
- 没有 session 或 connection-binding migration source。

## 13. References

- Nakama sockets: `https://heroiclabs.com/docs/nakama/concepts/sockets/`
- Nakama sessions: `https://heroiclabs.com/docs/nakama/concepts/session/`
- Nakama configuration session limits: `https://heroiclabs.com/docs/nakama/getting-started/configuration/`
- Pitaya session package documentation: `https://pkg.go.dev/github.com/topfreegames/pitaya/v3/pkg/session`
