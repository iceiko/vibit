# Authentication Command Protocol And Login Route Gate

状态：Draft v0.1
最后更新：2026-05-17
范围：为 public `AuthenticateWithDeviceCredential` route 定义未来 Protobuf command payloads、application route registration、protocol bridge 和 startup composition
依赖：`docs/authentication-contract-error-permission-surfaces.md`、`docs/device-credential-login-service-behavior-gate.md`、`docs/access-token-protocol-carrier-route-protection-gate.md`、`docs/runtime-authentication-startup-composition-gate.md`、`docs/runtime-protocol-adapter.md`、`docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0055`

配对英文源文件是 `docs/authentication-command-protocol-login-route-gate.md`。英文文件是权威版本。

## 1. 目的

Runtime 现在已经能在 application authentication service 中执行 device credential login，也能为 protected routes 验证 opaque access tokens，并且已经在 PostgreSQL runtime startup path 中接入 route protection。

剩下的 client-facing 缺口是 public login command route。Client 需要一个 protocol message 和已注册的 application route：

```text
runtime.authentication.AuthenticateWithDeviceCredential
```

本 gate 定义下一个 bounded implementation slice。它只允许添加 public device-credential login command protocol messages、generated Protobuf output、protocol bridge behavior、application route handler registration，以及把已有 service method 暴露出来所需的 PostgreSQL startup composition。

它不添加 session persistence、WebSocket handshake authentication、HTTP `Authorization` 或 Bearer carriers、cookies、query-string carriers、WebSocket subprotocol carriers、logout、refresh、cleanup、token rotation、repository interface changes、migration changes、new dependencies，或 direct Nakama/Pitaya public API compatibility。

## 2. 核心规则

Authentication command protocol and login route gate 是：

```yaml
authentication_command_protocol_login_route_gate: defined
implementation_authorized_by_this_standard: true
completed_gate_work_item: W-0119
future_implementation_work_item: W-0120
decision: ADR-0055
public_login_route: runtime.authentication.AuthenticateWithDeviceCredential
route_kind: command
route_public_policy_status: already_public_in_route_protection_policy
first_protocol_source: proto/vibit/authentication/v1/authentication.proto
first_generated_go_output: runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go
application_handler_owner: runtime/internal/app/bootstrap
authentication_service_owner: runtime/internal/app/authentication
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
startup_owner: runtime/cmd/vibit-server
first_composed_runtime_store: postgres
memory_store_login_route_status: unavailable_bootstrap
websocket_transport_credential_neutral: true
protobuf_envelope_change_status: unchanged
session_persistence_status: deferred
websocket_handshake_authentication_status: deferred
```

Future implementation 只能通过 public command route 暴露已有 device credential login service behavior。不得新增 authentication semantics，也不得把 authentication logic 移入 WebSocket transport、Protobuf envelope metadata、domain modules、repositories、migrations 或 generated files。

## 3. 选择的协议形状

第一条 authentication command protocol source 是：

```text
proto/vibit/authentication/v1/authentication.proto
```

第一批计划 messages 是：

```yaml
messages:
  AuthenticateWithDeviceCredentialRequest:
    fields:
      credential_proof: string
      requested_player_id: string
      client_instance_id: string
      account_creation_intent: AccountCreationIntent
  AuthenticateWithDeviceCredentialResponse:
    fields:
      authentication_status: string
      actor_kind: string
      player_id: string
      access_token: string
      token_type: string
      issued_at: string
      expires_at: string
      token_record_id: string
  AccountCreationIntent:
    values:
      - ACCOUNT_CREATION_INTENT_UNSPECIFIED
      - ACCOUNT_CREATION_INTENT_ALLOW_CREATE
      - ACCOUNT_CREATION_INTENT_AUTHENTICATE_EXISTING_ONLY
```

规则：

- 现有 `proto/vibit/protocol/v1/envelope.proto` 必须保持 unchanged。
- Envelope route 仍然是 `kind=command`、`module=runtime.authentication`、`name=AuthenticateWithDeviceCredential`。
- Payload type 是 `vibit.authentication.v1.AuthenticateWithDeviceCredentialRequest`。
- Response payload type 是 `vibit.authentication.v1.AuthenticateWithDeviceCredentialResponse`。
- `credential_proof` 和 `access_token` 是 secret values，不得出现在 errors、logs、events 或 repository records 中。
- `access_token` 是 existing service result 产生的一次性 client-visible response material；startup 或 protocol adapter code 不得存储它。
- `token_type` 必须反映 service-local posture `opaque_access_token`。
- Time values 必须使用 RFC3339 或 RFC3339Nano UTC text。

## 4. Route Registration Flow

Future implementation 必须保持这个层次顺序：

```yaml
login_route_flow:
  - websocket_transport_receives_binary_frame_without_reading_credentials
  - protobuf_adapter_decodes_existing_envelope
  - route_protector_allows_public_authentication_route_without_access_token
  - protobuf_adapter_decodes AuthenticateWithDeviceCredentialRequest
  - protocol_bridge_maps_request_to authentication.DeviceCredentialAuthenticationRequest
  - application_bootstrap_handler_calls authentication.Service.AuthenticateWithDeviceCredential
  - authentication_service_owns_unit_of_work_and token issuance
  - protocol_bridge_maps AuthenticationResult to AuthenticateWithDeviceCredentialResponse
  - protobuf_adapter_returns success or existing error envelope
```

规则：

- Route 必须显式注册。不得创建隐式 public route family。
- Application route registration 属于 `runtime/internal/app/bootstrap` 或等价 application-composition package，不属于 authentication module。
- Handler 只调用已有 `AuthenticateWithDeviceCredential` service method。
- Handler 不得 compute digests、compare verifiers、generate tokens、直接调用 repositories、open transactions，或把 transport metadata 当作 proof。
- 因为 authentication service 拥有自己的 unit-of-work boundary，application transaction wrapper 不得为 public authentication route 再打开一层 inventory-style command transaction。
- Memory runtime path 可以保持 login route unavailable，因为它没有 durable authentication repository capability。

## 5. Error Mapping

Future implementation 必须通过 existing application error envelope boundary 映射 service public errors。

第一版 mapping：

```yaml
service_public_errors:
  AUTHENTICATION_PROOF_MISSING: application_error_same_code
  AUTHENTICATION_PROOF_MALFORMED: application_error_same_code
  AUTHENTICATION_CREDENTIAL_INVALID: application_error_same_code
  AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE: application_error_same_code
  AUTHENTICATION_TOKEN_STORE_UNAVAILABLE: application_error_same_code
  AUTHENTICATION_NOT_IMPLEMENTED: application_error_same_code
```

规则：

- Public errors 不得披露 credential lookup hit/miss、player account existence、verifier key id、digest class、HMAC input/output，或 success output 中 audit-safe ids 之外的 token record internals。
- Error messages 不得包含 `credential_proof`、`access_token`、lookup digests、verifier digests、raw HMAC input、verifier keys 或 full concrete key ids。
- 失败登录不得返回 `access_token` 或 `token_record_id`。

## 6. Nakama 和 Pitaya 参考映射

Nakama 指导能力顺序：client 先 authenticate 并收到 session/token material，然后再使用 normal gameplay 或 realtime features。vibit 把这个思路适配为：先暴露 public login command，再处理 session persistence 或 handshake authentication。

Pitaya 指导分层：acceptors/connections、session context、routes 和 handlers 保持分离。vibit 把这个思路适配为：WebSocket transport 保持 credential-neutral，使用 existing envelope route，在 protocol adapter 中 bridge Protobuf payloads，并调用 application-owned handlers。

两个参考都不会覆盖 vibit 自身边界。本 gate 不采用 direct Nakama 或 Pitaya public API compatibility。

## 7. 必需实现测试

Future implementation slice 必须添加或更新 focused tests：

```yaml
required_tests:
  proto_source_and_generated_output_exist
  login_route_is_registered_only_when_authentication_service_is_composed
  login_route_is_public_and_does_not_require_access_token_wrapper
  login_route_bypasses_outer_transactional_dispatcher_unit_of_work
  login_request_maps_to_service_request_without_treating_metadata_as_proof
  login_success_maps_service_result_to_response_payload
  login_failure_maps_public_service_error_to_error_envelope
  login_errors_do_not_leak_credential_proof_or_access_token
  protected_routes_still_require_authenticated_wrapper
  websocket_transport_remains_credential_neutral
  existing_protobuf_envelope_remains_unchanged
```

Live PostgreSQL login verification 仍然 optional，不得成为默认 repository checks 的要求。

## 8. 延期项

本 gate 不授权：

- Session persistence。
- Runtime session tables 或 migrations。
- WebSocket handshake authentication。
- HTTP `Authorization`、Bearer、cookie、query string 或 WebSocket subprotocol credential carriers。
- Logout、refresh、cleanup、token rotation、token validation audit mutation 或 previous-token revocation。
- Repository interface changes。
- PostgreSQL adapter changes。
- Migration changes。
- Automatic startup migrations。
- New external dependencies。
- Memory-store durable authentication behavior。
- Direct Nakama 或 Pitaya API compatibility。

## 9. Verification

本 gate 的 repository check rule 是：

```text
runtime.authentication_command_protocol_login_route_gate
```
