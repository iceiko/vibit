# 访问令牌退出行为 Gate

状态：Draft v0.1
最后更新：2026-05-18
范围：未来 presented opaque access-token logout 行为的 gate-only 边界
依赖：`docs/authentication-service-behavior-implementation-gate.md`、`docs/access-token-validation-service-behavior-gate.md`、`docs/logout-revocation-active-connection-gate.md`、`docs/session-creation-composition-gate.md`、`docs/bound-identity-route-policy-gate.md`、`decisions/ADR-0071-logout-revocation-active-connection-gate.md`、`docs/reference-game-server-alignment.md`
规范决策：`ADR-0072`

配对英文原文是 `docs/logout-access-token-behavior-gate.md`。英文文件是权威版本。

## 1. 目的

`LogoutAccessToken` 目前已经作为语义 contract 和 fail-closed service method 存在，但 vibit 还没有选择执行 logout 的行为边界。

上一轮 gate 已经确定 logout、runtime session revocation 和 active WebSocket connection invalidation 是分离的生命周期决策。本 gate 把下一步 future behavior 缩小到第一版安全姿态：只撤销 presented opaque access token，由 application authentication service 执行，不改变 runtime session state，也不关闭 active WebSocket connections。

Nakama 提供关键产品压力：authenticated session material 需要包含 logout、refresh、expiration，以及 revoked material 不得继续授权 gameplay requests 的生命周期。Pitaya 提供关键分层压力：handlers 应接收 session/connection infrastructure 提供的 context，不应解析 credentials，也不应拥有 connection lifecycle side effects。

vibit 吸收这些经验，但保持自己的边界：第一版 logout behavior 是 token-record scoped、事务明确、redacted，并且与 socket close、connection registry、reconnect、protocol route exposure、refresh 和 logout-all behavior 分离。

```yaml
logout_access_token_behavior_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0156
decision: ADR-0072
check_rule: runtime.logout_access_token_behavior_gate
future_service_owner: runtime/internal/app/authentication
future_repository_capability_owner: runtime/internal/modules/authentication
future_transaction_owner: runtime/internal/app
existing_contract: contracts/runtime/authentication/commands/LogoutAccessToken.yaml
existing_service_method: runtime/internal/app/authentication.Service.LogoutAccessToken
existing_service_behavior: fail_closed_not_implemented
first_logout_scope: presented_access_token_only
proof_shape: opaque_base64url_unpadded_32_byte_access_token
proof_carrier_status: already_decoded_service_request_only
token_lookup_before_revocation_required: true
token_verifier_comparison_before_revocation_required: true
token_status_must_be_active_before_revocation: true
revoked_token_public_behavior: invalid_token
already_revoked_token_public_behavior: invalid_token
expired_token_public_behavior: invalid_token
runtime_session_revocation_added: false
active_connection_invalidation_added: false
connection_registry_added: false
websocket_close_policy_added: false
protobuf_logout_route_added: false
protobuf_session_carrier_added: false
existing_protobuf_envelope_change_added: false
refresh_behavior_added: false
logout_all_sessions_added: false
admin_revocation_added: false
cleanup_jobs_added: false
dependencies_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

本标准不实现 logout。

## 2. 所有权

未来 logout execution 必须由 application 拥有：

```yaml
future_logout_service_owner: runtime/internal/app/authentication
future_unit_of_work_owner: runtime/internal/app
repository_interface_owner: runtime/internal/modules/authentication
postgres_adapter_owner: runtime/internal/platform/persistence/postgres
transport_owner: runtime/internal/platform/transport/ws
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
domain_handler_owner: runtime/internal/modules/*
```

规则：

- Authentication service 拥有 proof validation、token lookup、verifier comparison、public error collapse，以及发起 token revoke 的请求。
- Authentication repository 只拥有 storage-neutral token record mutation vocabulary。它不得解析 raw access-token text、计算 digest、比较 verifier、决定 logout proof validity，或关闭 connections。
- Unit-of-work boundary 必须包住 token lookup、verifier comparison decision inputs、token revocation mutation 和 commit outcome。
- WebSocket transport 不得解析 logout credentials，也不得决定 logout side effects。
- Protobuf adapters 只有在后续 protocol route gate 授权后，才可以暴露 logout。
- Domain modules 必须接收 already-authenticated request identity 或 logout result；它们不得解析 access tokens 或直接调用 token repositories。

## 3. 未来行为顺序

后续 implementation 只有按以下顺序才可以执行 logout：

```yaml
future_logout_sequence:
  - reject_missing_or_malformed_access_token_before_unit_of_work
  - compute_access_token_lookup_digest
  - begin_application_unit_of_work
  - obtain_authentication_repository_from_unit_of_work
  - find_token_record_by_lookup_digest
  - require_token_kind_access_token
  - require_token_status_active
  - require_not_expired_at_service_clock_now
  - require_expected_audience
  - require_supported_verifier_algorithm_and_version
  - compute_access_token_verifier_digest_using_record_key_id
  - compare_verifier_digest_constant_time
  - revoke_presented_token_record_with_reason_logout_presented_access_token
  - commit_unit_of_work
  - return_revoked_result_after_commit_only
```

规则：

- Missing 或 malformed token proof 必须在打开 unit of work 前被拒绝。
- Token lookup 必须使用现有 lookup digest helper；raw token text 不得进入 repository。
- Verifier digest comparison 必须发生在 revoke 之前，避免 lookup collision 或错误 raw token 撤销 token record。
- 第一版 implementation posture 只能撤销 verified presented token record。
- Revoked、expired、wrong-audience、wrong-kind、unsupported、unknown-key、lookup-missing 或 mismatched token 必须折叠为 public invalid-token behavior。
- Repository unavailable 必须折叠为 public token-store-unavailable behavior。
- Raw token text 不得出现在结果、错误、日志、事件或测试名称中。

## 4. 事务边界

未来第一版 logout behavior 必须对 token state 保持事务化：

```yaml
transaction_boundary:
  includes:
    - token_record_lookup_by_lookup_digest
    - token_record_status_and_expiry_check
    - verifier_digest_comparison_decision_inputs
    - token_record_revocation_mutation
  excludes:
    - runtime_session_revocation
    - active_connection_invalidation
    - websocket_close
    - protocol_response_mapping
    - cleanup_jobs
```

规则：

- Service 只有在 unit of work commit 后才可以返回 `revoked`。
- 如果 commit 失败，public result 不得声称 logout 成功。
- Logout 必须明确 idempotence posture：本 gate 为 already revoked、expired 或 inactive tokens 选择 fail-closed invalid-token behavior。后续 ADR 可以选择 idempotent success，但必须显式修改标准和测试。
- Runtime session revocation 不属于第一版 transaction。`runtime_sessions` row 可以保持 active，直到后续 session revocation policy 改变它。
- Active WebSocket connection invalidation 不在 SQL transaction 控制内，并且本 gate 不授权。

## 5. Public Result And Error 边界

未来 public logout behavior 必须保持最小且 redacted：

```yaml
candidate_public_success:
  status: revoked
  revoked: true
  logout_scope: presented_access_token
  token_type: opaque_access_token

candidate_public_errors:
  missing: AUTHENTICATION_TOKEN_MISSING
  malformed: AUTHENTICATION_TOKEN_MALFORMED
  invalid: AUTHENTICATION_TOKEN_INVALID
  unavailable: AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
```

规则：

- Public invalid-token behavior 不得区分 lookup miss、already revoked、expired、wrong audience、wrong token kind、unsupported verifier metadata、unknown verifier key id、verifier mismatch 或 missing player account。
- Success 不得包含 raw token text、lookup digest、verifier digest、verifier key id、Authorization header、cookie、query string、WebSocket subprotocol value、session id、connection id 或 remote address。
- Internal errors 必须保留足够类型供测试和 redacted observability 使用，但 `Error()` 不得泄露 secrets。
- 如果未来 result 包含 `TokenRecordID`，它必须只在 observability standard 分类之后作为 internal/audit-safe identifier 使用。

## 6. 与 Runtime Session 和 Active Connections 的关系

本 gate 有意把 runtime session 和 active connection behavior 分离：

```yaml
runtime_session_revocation_added: false
active_connection_invalidation_added: false
connection_registry_added: false
websocket_close_policy_added: false
session_last_seen_update_added: false
duplicate_connection_replacement_added: false
reconnect_epoch_behavior_added: false
```

规则：

- Presented access token logout 不得隐式 revoke player 的所有 runtime sessions。
- Presented access token logout 不得隐式 revoke player、credential、device 或 account 的所有 tokens。
- Presented access token logout 不得隐式关闭 active WebSocket connections。
- 当后续 protected requests 执行 request-level access-token validation 时，既有 validation 会拒绝 revoked token。
- Token revocation 后 bound-connection 和 bound-session route behavior 仍是后续 policy question。
- Reconnect、resume、duplicate connection replacement、connection epoch 和 targeted kick/disconnect 继续 deferred。

## 7. 与 Protocol 的关系

本 gate 不通过 Protobuf 或 WebSocket 暴露 logout：

```yaml
protobuf_logout_route_added: false
protobuf_authentication_message_added: false
protobuf_session_carrier_added: false
existing_protobuf_envelope_change_added: false
websocket_handshake_authentication_added: false
transport_credential_carrier_added: false
generated_output_added: false
```

规则：

- `proto/vibit/authentication/v1/authentication.proto` 在本 gate 下不得添加 logout request 或 response messages。
- Generated Go Protobuf output 不得手工编辑。
- WebSocket transport 不得从 headers、cookies、query strings 或 subprotocol values 解析 logout credentials。
- Protocol adapters 不得调用 `LogoutAccessToken`，直到后续 protocol logout route gate 定义 carrier、route name、response mapping 和 error mapping。

## 8. 与既有 Service Skeleton 的关系

当前 service 可以保持 fail-closed：

```yaml
current_logout_service_method: LogoutAccessToken
current_status: not_implemented
current_public_error: AUTHENTICATION_NOT_IMPLEMENTED
behavior_change_authorized_by_this_gate: false
```

规则：

- 本 gate 可以指导未来 implementation slice，但不得改变 `runtime/internal/app/authentication/service.go`。
- `RefreshAccessToken` 继续 unsupported。
- Access-token validation behavior 不变。
- Device credential login behavior 不变。
- Session creation behavior 不变。
- Route policy 不变。

## 9. 未来测试期望

后续 implementation slice 必须添加聚焦测试：

- Missing token proof 在 unit of work 之前被拒绝。
- Malformed token proof 在 unit of work 之前被拒绝。
- Lookup miss 折叠为 public invalid-token behavior。
- Already revoked token 折叠为 public invalid-token behavior。
- Expired token 折叠为 public invalid-token behavior。
- Wrong token kind、audience、algorithm、version 或 key id 折叠为 public invalid-token behavior。
- Verifier mismatch 不会 revoke token record。
- Repository lookup/revocation/commit failures 不得声称成功。
- Successful logout 只调用一次 `RevokeToken`，reason 为 `logout_presented_access_token`。
- Success 只在 commit 后返回。
- Raw token text、lookup digest、verifier digest、verifier key id、session id、connection id、Authorization header、cookie、query string 和 WebSocket subprotocol value 不得出现在 public result 和 error strings。
- Runtime session repository、connection registry、WebSocket transport 和 Protobuf adapter 不被调用。

## 10. Nakama 和 Pitaya 参考映射

Nakama 参考映射：

- 吸收 session material 具有 lifecycle state，且 revoked material 不得授权未来 gameplay requests 的经验。
- 把 logout 适配到 vibit 的 opaque server-issued access-token record model，而不是 Nakama-compatible JWT/session APIs。
- Refresh、logout-all、session management APIs、dashboard/admin revocation 和 realtime socket invalidation 都保留为后续显式 surfaces。

Pitaya 参考映射：

- 吸收 connection/session infrastructure 与 handler logic 分离的经验。
- 把 logout 适配为 application service behavior，未来可以通过窄接口与 connection/session infrastructure 协作。
- Frontend/backend cluster routing、distributed kick/disconnect 和 server-to-server invalidation 都保留为后续 gates。

Nakama 和 Pitaya 只指导 capability pressure。本 gate 不添加与任何一个项目的 direct public API compatibility。

## 11. 非目标

本 gate 不授权：

- `LogoutAccessToken` implementation。
- Token revocation execution。
- Runtime session revocation。
- Active WebSocket connection invalidation。
- Connection registry。
- WebSocket close policy。
- Kick/disconnect behavior。
- Reconnect、resume、duplicate replacement 或 epoch behavior。
- Protobuf logout route。
- Protocol session carrier。
- 现有 Protobuf envelope changes。
- Refresh token behavior。
- Logout-all-sessions behavior。
- Admin revocation。
- Cleanup jobs。
- 新 dependencies。
- Memory durable session behavior。
- Direct Nakama 或 Pitaya API compatibility。

## 12. 必需 Follow-Up

如果后续选择 implementation slice，建议是：

```text
implement_logout_access_token_behavior
```

该 slice 必须把行为保持在 `runtime/internal/app/authentication` 内，使用既有 repository 和 verifier helper boundaries，添加聚焦测试，并继续保留所有 protocol、transport、session revocation、active connection、reconnect、dependency 和 direct compatibility deferrals，除非有独立 ADR 授权。
