# 访问令牌协议载体与路由保护门禁

状态：草案 v0.1
最后更新：2026-05-17
范围：在添加协议或路由保护实现之前，定义未来访问令牌 proof 载体、请求级路由保护、应用验证交接、公开路由策略、脱敏、测试和延期项
依赖：`docs/access-token-validation-service-behavior-gate.md`、`docs/runtime-protocol-adapter.md`、`docs/game-protocol.md`、`docs/session-persistence-websocket-handshake-decision-gates.md`
规范决策：`ADR-0053`

英文文件 `docs/access-token-protocol-carrier-route-protection-gate.md` 是权威版本。本文是配套简体中文翻译。

## 1. 目的

应用认证服务现在已经能在服务本地边界内签发和验证 opaque 访问令牌。下一个风险是为了让客户端携带令牌，而把 Bearer 解析、cookie、查询参数、WebSocket 握手认证、Protobuf envelope 变更或路由保护直接塞进 transport、protocol 或 domain 包。

本门禁选择 service-local authentication 之后的下一个里程碑方向：先定义请求级访问令牌载体和路由保护边界，再进入实现。

这是一个仅门禁标准。它不添加 `.proto` 源文件、生成文件、协议适配器代码、路由保护代码、启动 wiring、session persistence、WebSocket 握手认证、仓库变更、迁移、logout、refresh、cleanup、依赖或更宽泛的生产认证行为。

## 2. 核心规则

访问令牌协议载体与路由保护门禁为：

```yaml
access_token_protocol_carrier_route_protection_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0114
completed_gate_work_item: W-0113
planned_owner: runtime/internal/app
planned_protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
planned_transport_owner: runtime/internal/platform/transport/ws
planned_authentication_service_owner: runtime/internal/app/authentication
validation_model: request_level_validation
first_proof_carrier_posture: protobuf_payload_wrapper
first_wrapper_message_candidate: vibit.authentication.v1.AuthenticatedRequest
protobuf_envelope_change_status: deferred
websocket_handshake_authentication_status: deferred
session_persistence_status: deferred
startup_wiring_status: deferred
route_policy_status: defined_not_implemented
```

只有后续工作项明确授权具体协议源、生成输出、适配器交接、应用路由策略和启动组合切片时，未来实现才可以把 token proof 暴露给普通 gameplay 路由。

## 3. 已选择的载体姿态

第一个路由保护载体姿态是 protected request 的显式 Protobuf payload wrapper。

计划语义形态：

```yaml
authenticated_request_payload_wrapper:
  message_candidate: vibit.authentication.v1.AuthenticatedRequest
  owner: proto/vibit/authentication/v1
  outer_envelope_kind: original_command_or_query
  outer_envelope_module: original_domain_module
  outer_envelope_name: original_domain_route_name
  outer_payload_type: vibit.authentication.v1.AuthenticatedRequest
  inner_payload_type: original_payload_type
  inner_payload: original_payload_bytes
  proof_field: access_token
  proof_encoding: base64url_unpadded_32_byte_opaque_access_token
```

规则：

- 第一个载体姿态不得改变现有 Protobuf envelope 形状。
- `Envelope.session.player_id`、`Envelope.session.session_id`、`Envelope.session.connection_id` 和 `Envelope.session.connection_epoch` 仍只是 metadata。
- 访问令牌不得放进 `Session`、`Target`、路由字段、`request_id`、错误详情、日志或连接 metadata。
- 在此姿态中，不得从 HTTP `Authorization`、`Bearer` 字符串、cookie、查询字符串或 WebSocket subprotocol 解析访问令牌。
- wrapper 只是协议 payload 载体。它不验证 proof、不选择 token 记录、不比较 digest、不决定玩家账号状态，也不构造领域响应。
- 只有 wrapper 结构被接受并且 proof 验证成功后，才可以解码 inner payload type。

## 4. 请求级验证流程

未来实现必须保持以下层次顺序：

```yaml
request_level_validation_flow:
  - websocket_transport_receives_binary_frame_without_reading_authentication_proof
  - protobuf_adapter_decodes_existing_envelope
  - protobuf_adapter_recognizes_authenticated_payload_wrapper_for_protected_route
  - protobuf_adapter_extracts_access_token_text_as_secret
  - protobuf_adapter_keeps_inner_payload_bytes_undispatched
  - application_route_protection_policy_requires_authenticated_identity_for_protected_route
  - application_authentication_validator_calls_ValidateAccessToken
  - authentication_service_validates_token_and_returns_RequestIdentity
  - application_handoff_sets_validated_identity_with_SessionValidated_false
  - protocol_adapter_decodes_inner_payload_for_original_route
  - application_dispatch_calls_domain_handler_with_validated_identity
```

规则：

- Domain handler 只能接收规范化后的 `RequestIdentity`；不得接收或解析访问令牌 proof。
- 当 proof 缺失、畸形、无效、不可用，或没有用已选择 wrapper 正确携带时，路由保护层必须在 domain dispatch 前拒绝 protected route。
- Metadata-only identity 永远不能满足 protected route policy。
- 在后续 session persistence gate 被选择和实现之前，`SessionValidated` 仍保持 false。
- WebSocket transport 保持 credential-neutral。
- Protocol adapter 只能把 wrapper 字段作为窄 proof handoff 交给应用验证。

## 5. 路由策略门禁

在实现路由保护前，未来工作必须声明应用拥有的路由策略。

第一版策略要求：

```yaml
route_policy:
  owner: runtime/internal/app
  public_routes:
    - runtime.authentication.AuthenticateWithDeviceCredential
  protected_route_default: authenticated_player_required
  protected_identity_requirement:
    identity_status: validated
    actor_kind: player
    player_id_validated: true
    session_validated: false_allowed_until_session_persistence
  domain_module_token_parsing: forbidden
```

规则：

- public route 必须显式列出。普通 gameplay route 没有隐式 public 默认值。
- protected route 要求 `IdentityValidationValidated`、`ActorKindPlayer` 和 `PlayerIDValidated`。
- 只有因为 access-token validation 是请求 proof 且 session persistence 仍延期，route policy 才可以允许 `SessionValidated: false`。
- Inventory permission policy 可以继续执行领域权限，但它不能替代路由级认证策略。
- 路由策略必须可以在不依赖 live PostgreSQL 的情况下测试。

## 6. 错误映射

未来实现必须在 domain dispatch 前映射认证验证失败。

第一版公开姿态：

```yaml
route_protection_error_mapping:
  missing_wrapper_or_missing_token: AUTHENTICATION_TOKEN_MISSING
  malformed_wrapper_or_malformed_token: AUTHENTICATION_TOKEN_MALFORMED
  invalid_token_family: AUTHENTICATION_TOKEN_INVALID
  validation_dependency_unavailable: AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
  protected_route_without_validated_identity: AUTHENTICATION_TOKEN_INVALID
```

规则：

- 公开协议错误不得透露 token lookup 命中/未命中、token 生命周期状态、verifier key id、audience mismatch、玩家存在性、玩家账号状态或 verifier mismatch。
- 路由未授权失败不得包含 raw token text、inner payload bytes、lookup digest、verifier digest、HMAC input/output、verifier key 或完整具体 key id。
- 到 Protobuf error envelope 的映射必须复用现有 application-error 边界，或等待后续显式映射 gate。

## 7. 未来必需工件

后续实现切片必须在路由保护生效前定义或更新这些工件：

```yaml
required_future_artifacts:
  protocol_source: proto/vibit/authentication/v1/authenticated_request.proto
  generated_go_proto_output: runtime/internal/generated/proto/vibit/authentication/v1/authenticated_request.pb.go
  protocol_adapter_tests: runtime/internal/platform/protocol/protobuf/*authentication*_test.go
  application_route_policy_source: runtime/internal/app/*route*_auth*.go
  application_route_policy_tests: runtime/internal/app/*route*_auth*_test.go
  authentication_validator_adapter_source: runtime/internal/app/authentication/*validator*.go
  process_startup_wiring: deferred_until_separate_startup_work_item
```

本门禁不创建这些工件。

## 8. 必需测试

未来实现必须添加聚焦测试，覆盖：

- protected route 缺少 wrapper。
- wrapper 中缺少 access token。
- wrapper 中 access token 畸形。
- invalid token 的公开失败折叠。
- store unavailable 的公开失败映射。
- metadata-only identity 被 protected route 拒绝。
- valid token 在 domain dispatch 前生成 validated player identity。
- access-token validation 之后 `SessionValidated` 仍保持 false。
- public authentication route 显式公开且不要求 token proof。
- WebSocket transport 不解析 `Authorization`、cookie、query string、subprotocol 或 bearer 值。
- Protobuf envelope 字段仍只是 metadata。
- raw access-token text 和 inner payload bytes 不出现在错误中。

## 9. 延期项

本门禁不授权：

- 添加 `.proto` 文件。
- 生成 Protobuf 输出。
- 修改 `proto/vibit/protocol/v1/envelope.proto`。
- 修改 WebSocket 握手认证。
- 解析 HTTP header、cookie、query string、bearer string 或 subprotocol。
- 把认证服务接入进程启动。
- 添加 session persistence。
- 添加 runtime session 表或迁移。
- 修改认证仓库接口。
- 修改 PostgreSQL adapter。
- 实现 logout、refresh、cleanup、token rotation 或 token validation audit mutation。
- 添加依赖。
- 声明直接 Nakama 或 Pitaya API 兼容。

## 10. 参考映射

Nakama 说明了 gameplay request 前验证 token 的能力需求。vibit 通过请求级验证和显式路由策略适配该能力，而不采用直接 API 兼容。

Pitaya 说明了 frontend connection handling 与 backend route handler context 的分离。vibit 通过保持 WebSocket transport credential-neutral，并在 domain handler 运行前向 application dispatch 传递 validated identity 来适配该模式。

## 11. 验证

本门禁的仓库检查规则是：

```text
runtime.access_token_protocol_carrier_route_protection_gate
```
