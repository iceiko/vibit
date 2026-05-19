# Protocol Logout Route Gate

状态：Draft v0.1
最后更新：2026-05-18
范围：未来 client-facing `runtime.authentication.LogoutAccessToken` protocol route behavior 的 gate-only 边界
依赖：`docs/nakama-pitaya-product-parity-roadmap.md`、`docs/logout-access-token-behavior-gate.md`、`decisions/ADR-0073-logout-access-token-behavior-implementation.md`、`docs/authentication-command-protocol-login-route-gate.md`、`docs/websocket-close-policy-gate.md`、`docs/runtime-protocol-adapter.md`、`docs/reference-game-server-alignment.md`
权威决策：`ADR-0079`

英文文件 `docs/protocol-logout-route-gate.md` 是权威版本。本文是简体中文译本。

## 1. 目的

vibit 现在已经可以通过 `authentication.Service.LogoutAccessToken` revoke 一个已经验证过的 presented opaque access token，但 client 还没有 protocol route 来请求这个行为。

Nakama/Pitaya 产品同级路线图把 logout 放进近期 runtime lifecycle closure。Nakama server runtime 把 session logout 和 session disconnect 作为不同操作暴露出来，这说明 token/session invalidation 和 socket disconnection 是不同生命周期动作。Pitaya 文档强调 sessions、handler routing、acceptors、kick/disconnect 风格能力是不同 framework surfaces。vibit 吸收这些经验，把 protocol logout route 定义成一个明确边界：暴露 token logout，但不偷偷关闭 socket，也不把 token 逻辑移动到 transport。

本标准定义下一步有边界的 protocol route implementation slice。本 work item 只定义 gate。

```yaml
protocol_logout_route_gate: defined
implementation_authorized_by_this_standard: true
completed_gate_work_item: W-0169
future_implementation_work_item: W-0170
decision: ADR-0079
check_rule: runtime.protocol_logout_route_gate
public_logout_route: runtime.authentication.LogoutAccessToken
route_kind: command
route_protection_posture: explicit_service_validated_token_lifecycle_route
proof_carrier_posture: access_token_in_logout_request_payload
first_protocol_source: proto/vibit/authentication/v1/authentication.proto
first_generated_go_output: runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go
request_message_candidate: vibit.authentication.v1.LogoutAccessTokenRequest
response_message_candidate: vibit.authentication.v1.LogoutAccessTokenResponse
application_handler_owner: runtime/internal/app/bootstrap
authentication_service_owner: runtime/internal/app/authentication
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
startup_owner: runtime/cmd/vibit-server
first_composed_runtime_store: postgres
memory_store_logout_route_status: unavailable_bootstrap
transaction_bypass_required: true
protobuf_envelope_change_status: unchanged
websocket_transport_credential_neutral: true
first_logout_scope: presented_access_token_only
already_revoked_token_public_behavior: invalid_token
expired_token_public_behavior: invalid_token
logout_closes_socket: false
runtime_session_revocation_added: false
active_connection_invalidation_added: false
concrete_socket_close_added: false
reconnect_epoch_behavior_added: false
protocol_session_carrier_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. 所有权

未来 protocol logout route implementation 必须保持这些 ownership boundaries：

```yaml
authentication_service_owner: runtime/internal/app/authentication
application_handler_owner: runtime/internal/app/bootstrap
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
startup_owner: runtime/cmd/vibit-server
websocket_transport_owner: runtime/internal/platform/transport/ws
close_policy_owner: runtime/internal/app/connection
```

规则：

- Authentication service 仍然是 raw access-token proof validation、digest computation、verifier comparison、token lookup、token posture checks 和 token revocation mutation 的唯一 owner。
- 未来 route handler 只能调用已有 `LogoutAccessToken` service method。它不能 compute digests、compare verifiers、直接调用 repositories、自己开启 transaction、关闭 sockets 或 mutate runtime sessions。
- Protobuf adapter 只能映射 request 和 response payload。它不能决定 logout validity、route authorization、session revocation、active connection invalidation 或 transport close behavior。
- WebSocket transport 仍然 credential-neutral 和 policy-neutral。它不能从 headers、cookies、query strings、bearer values 或 subprotocols 解析 logout proof。
- Close policy 和 active connection registry 仍与 protocol logout route exposure 分离。

## 3. Route 与 Carrier Posture

第一版 logout route 是：

```yaml
route:
  kind: command
  module: runtime.authentication
  name: LogoutAccessToken
  semantic_contract: contracts/runtime/authentication/commands/LogoutAccessToken.yaml
```

第一版 route-protection posture 是：

```yaml
route_protection_posture: explicit_service_validated_token_lifecycle_route
authenticated_request_wrapper_required: false
access_token_source: LogoutAccessTokenRequest.access_token
service_validation_required: true
```

原因：

- Logout 需要 exact presented access token，这样 service 才能在 verifier comparison 后 revoke 同一个 token record。
- 这个 route 不是普通 gameplay protected route。它是 token lifecycle route，payload 携带要被 revoke 的 proof。
- 把 route 标记为 explicit 且 service-validated，可以避免 route protector 在 logout service revoke 之前消费掉 token。
- Route 仍必须显式注册；不存在隐式 public authentication route family。

## 4. Protocol Shape

未来 implementation 可以扩展现有 authentication Protobuf source：

```text
proto/vibit/authentication/v1/authentication.proto
```

计划中的 request：

```yaml
LogoutAccessTokenRequest:
  access_token:
    type: string
    secret: true
    source: presented_opaque_access_token
  logout_reason:
    type: string
    required: false
    public_safe: true
```

计划中的 response：

```yaml
LogoutAccessTokenResponse:
  logout_status:
    type: string
    values:
      - revoked
      - rejected
  revoked:
    type: bool
  logout_scope:
    type: string
    first_value: presented_access_token
  revoked_at:
    type: string
    format: rfc3339nano_utc
  token_record_id:
    type: string
    visibility: audit_safe
```

规则：

- 现有 `proto/vibit/protocol/v1/envelope.proto` 必须保持不变。
- `access_token` 是 secret input，不能出现在 errors、logs、events、close intents、connection records 或 test names 中。
- `logout_reason` 不能被当作可信 server policy input，也不能携带 secrets。
- `logout_scope` 必须把第一版 posture 暴露为 `presented_access_token`，即使内部 service vocabulary 使用更窄的 implementation enum。
- Time values 必须使用 RFC3339Nano UTC text。
- `token_record_id` 如果暴露，只是 audit-safe。它不是 proof、不是 session id，也不是 connection target。

## 5. Future Route Flow

未来 implementation 必须保持这个顺序：

```yaml
logout_route_flow:
  - websocket_transport_receives_binary_frame_without_reading_credentials
  - protobuf_adapter_decodes_existing_envelope
  - route_policy_allows_explicit_service_validated_logout_route_without_authenticated_wrapper
  - protobuf_adapter_decodes LogoutAccessTokenRequest
  - protocol_bridge_maps_request_to authentication.LogoutAccessTokenRequest
  - application_bootstrap_handler_calls authentication.Service.LogoutAccessToken
  - authentication_service_owns_unit_of_work_and token revocation
  - protocol_bridge_maps LogoutAccessTokenResult to LogoutAccessTokenResponse
  - protobuf_adapter_returns success or existing error envelope
```

规则：

- Authentication route 必须 bypass 外层 `TransactionalDispatcher` unit of work，因为 authentication service 拥有自己的 unit-of-work boundary。
- Failed logout 不能返回 `revoked: true`、`revoked_at` 或 success status。
- Successful logout 只能在 service 报告 unit-of-work commit 成功后返回。
- 未来 route 不能为 logout 解码或 dispatch `AuthenticatedRequest` wrapper。

## 6. Public Error Mapping

未来 protocol behavior 必须映射 service public errors，且不泄露 proof details：

```yaml
service_public_errors:
  AUTHENTICATION_TOKEN_MISSING: application_error_same_code
  AUTHENTICATION_TOKEN_MALFORMED: application_error_same_code
  AUTHENTICATION_TOKEN_INVALID: application_error_same_code
  AUTHENTICATION_TOKEN_STORE_UNAVAILABLE: application_error_same_code
  AUTHENTICATION_NOT_IMPLEMENTED: application_error_same_code
```

规则：

- Lookup miss、already revoked token、expired token、wrong audience、wrong kind、unsupported verifier metadata、unknown verifier key id 和 verifier mismatch 必须继续 collapse 为 public invalid-token behavior。
- Error messages 不能包含 raw token text、lookup digest、verifier digest、HMAC input、verifier key id、headers、cookies、query strings、subprotocol values、session ids、connection ids、remote addresses 或 database errors。
- Memory runtime 可以继续让 durable logout route behavior unavailable，直到未来 ratify memory authentication posture。

## 7. 与 Socket Close 和 Sessions 的关系

本 gate 不改变 socket 或 session lifecycle behavior：

```yaml
logout_closes_socket: false
close_policy_called_by_logout_route: false
active_connection_invalidation_added: false
runtime_session_revocation_added: false
bound_connection_identity_after_logout_policy: deferred
bound_session_identity_after_logout_policy: deferred
reconnect_epoch_behavior_added: false
protocol_session_carrier_added: false
```

规则：

- Successful protocol logout 只 revoke 已验证的 presented access-token record。
- Successful protocol logout 不得隐式关闭 WebSocket connection。
- Successful protocol logout 不得隐式 revoke linked runtime session。
- Successful protocol logout 不得 invalidate active connection registry records。
- 现有 request-level token validation 在后续 protected requests 运行时会拒绝 revoked token。
- Logout 后 bound connection 是否还能使用已经绑定的 identity，仍是之后 bound-identity/session/reconnect policy 问题。

## 8. Required Future Tests

未来 implementation slice 必须增加 focused tests：

```yaml
required_tests:
  proto_source_and_generated_output_include_logout_messages
  logout_route_is_registered_only_when_authentication_service_is_composed
  logout_route_is_explicit_service_validated_token_lifecycle_route
  logout_route_bypasses_outer_transactional_dispatcher_unit_of_work
  logout_request_maps_access_token_to_service_request_without_logging_it
  logout_success_maps_service_result_to_response_payload
  logout_failure_maps_public_service_error_to_error_envelope
  logout_errors_do_not_leak_access_token
  websocket_transport_remains_credential_neutral
  existing_protobuf_envelope_remains_unchanged
  logout_route_does_not_call_close_policy_or_close_socket
  protected_gameplay_routes_still_require_authenticated_wrapper
```

Live PostgreSQL verification 仍然 optional，且不得成为默认 repository checks 的要求。

## 9. Nakama And Pitaya Reference Mapping

Nakama reference mapping：

- 采纳 logout/token invalidation 与 explicit session disconnect 之间的区别。
- 改造其 product expectation：client 可以通过稳定 server API 请求 logout。
- 推迟 Nakama session token compatibility、refresh token compatibility、logout-all semantics、realtime socket disconnect compatibility 和 dashboard/admin compatibility。

Pitaya reference mapping：

- 采纳 client connection acceptors、sessions、routes、handlers 和 connection management 的分离。
- 改造其 session/kick lifecycle 经验：logout proof validation 留在 application service behavior，socket close 留给未来独立 policy/handoff。
- 推迟 Pitaya route naming compatibility、frontend/backend routing、distributed kick/disconnect、groups integration 和 RPC/session propagation。

## 10. 非目标

本 gate 不授权：

- 在本 work item 添加 Protobuf logout messages。
- 在本 work item 添加 generated Go Protobuf output。
- 在本 work item 注册 `runtime.authentication.LogoutAccessToken`。
- 改变现有 Protobuf envelope。
- 改变 WebSocket handshake authentication。
- 解析 HTTP headers、bearer values、cookies、query strings 或 WebSocket subprotocols。
- 关闭 sockets、选择 close codes、选择 close reason text 或发送 protocol close messages。
- Revoke runtime sessions。
- Invalidate active connection registry records。
- 添加 reconnect、resume、duplicate replacement 或 connection epoch behavior。
- 添加 protocol session carriers。
- 添加 refresh、logout-all、admin revocation、cleanup jobs、presence、chat、social modules、matchmaking、match runtime、SDKs、cluster、RPC、service discovery、dependencies 或 direct Nakama/Pitaya API compatibility。

## 11. Verification

本 gate 的 repository check rule 是：

```text
runtime.protocol_logout_route_gate
```
