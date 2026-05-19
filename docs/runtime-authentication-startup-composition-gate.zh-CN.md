# Runtime Authentication Startup Composition Gate

状态：Draft v0.1
最后更新：2026-05-17
范围：为已有 application authentication service、verifier key environment loader、request-level route protection 和 Protobuf frame handler 定义 runtime process startup composition
依赖：`docs/environment-verifier-key-loader-gate.md`、`docs/access-token-validation-service-behavior-gate.md`、`docs/access-token-protocol-carrier-route-protection-gate.md`、`docs/runtime-protocol-adapter.md`、`docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0054`

配对英文源文件是 `docs/runtime-authentication-startup-composition-gate.md`。英文文件是权威版本。

## 1. 目的

Runtime 现在已经能在 application service 层验证 opaque access-token proof，并且 Protobuf adapter 在注入 `RouteProtector` 后可以执行 request-level route protection。

剩下的缺口是 process startup composition：`runtime/cmd/vibit-server` 必须把已有 authentication service、verifier key configuration、PostgreSQL unit-of-work runner、route access-token validator、route protector 和 Protobuf frame handler 组装起来。

这是 startup-composition gate。它不新增认证语义，不添加 WebSocket handshake authentication、session persistence、authentication command Protobuf messages、login route registration、repository interface changes、PostgreSQL adapter changes、migrations、dependencies、logout、refresh、cleanup、token rotation、token validation audit mutation 或 broader production authentication behavior。

## 2. 核心规则

Runtime authentication startup composition gate 是：

```yaml
runtime_authentication_startup_composition_gate: defined
implementation_authorized_by_this_standard: true
completed_gate_work_item: W-0116
future_implementation_work_item: W-0117
decision: ADR-0054
startup_owner: runtime/cmd/vibit-server
application_authentication_owner: runtime/internal/app/authentication
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
transport_owner: runtime/internal/platform/transport/ws
first_composed_runtime_store: postgres
memory_store_authentication_status: bootstrap_metadata_only
websocket_transport_credential_neutral: true
protobuf_envelope_change_status: unchanged
session_persistence_status: deferred
websocket_handshake_authentication_status: deferred
```

Future implementation 只能通过 composition 把已有 authentication behavior 接入 startup。不得把 authentication logic 移到 WebSocket transport、Protobuf envelope metadata、domain handlers、repositories 或 generated files。

## 3. 选择的启动路径

第一条 startup composition path 限定为：

```yaml
runtime_store: postgres
selector: VIBIT_RUNTIME_STORE=postgres
startup_file: runtime/cmd/vibit-server/main.go
required_store_capabilities:
  - postgres.UnitOfWork.NewAuthenticationRepository
  - postgres.UnitOfWork.NewPlayerAccountRepository
  - postgres.NewPoolRunner
```

规则：

- PostgreSQL startup 在 required authentication verifier key configuration 缺失或无效时必须 fail closed。
- 默认 in-memory runtime store 仍是 bootstrap path，不得假装拥有 durable authentication repository capability。
- Startup 不得自动 apply migrations。
- Startup 不得自行创建 credential 或 token records。
- Startup 不得注册 authentication command routes，除非后续 protocol/route work item 授权。

## 4. 组合流程

Future implementation 必须遵循这个层次顺序：

```yaml
startup_composition_flow:
  - read_runtime_store_selection
  - open_postgres_pool_for_postgres_store
  - build_postgres_unit_of_work_runner
  - load_verifier_key_set_from_explicit_environment_lookup
  - create_authentication_service_with_existing_dependencies
  - create_route_access_token_validator_from_service
  - create_application_route_protector
  - inject_route_protector_into_protobuf_frame_handler
  - mount_websocket_transport_with_opaque_frame_handler
```

规则：

- WebSocket transport 继续只接收和发出 opaque frames。
- Protobuf frame handler 继续 unwrap 已 ratify 的 `vibit.authentication.v1.AuthenticatedRequest` payload wrapper。
- Application route protector 仍然拥有 protected-route authentication policy。
- Domain modules 接收 `RequestIdentity`；不得接收 raw access-token proof。
- `SessionValidated` 在 session persistence 被 ratify 前保持 false。

## 5. 启动依赖

Future implementation 只可以使用已经 ratify 的 helpers 和标准库包：

```yaml
authentication_service_dependencies:
  unit_of_work_runner: postgres.NewPoolRunner(pool)
  verifier_key_set: authentication.LoadVerifierKeySetFromEnvironment(lookup)
  access_token_random: crypto/rand.Reader
  clock: startup_owned_system_clock
  token_record_id_generator: startup_owned_standard_library_generator
  access_token_lifetime:
    default: 1h
    optional_environment: VIBIT_AUTH_ACCESS_TOKEN_TTL
  token_audience:
    default: vibit_gameplay_runtime_requests
    optional_environment: VIBIT_AUTH_TOKEN_AUDIENCE
```

规则：

- `VIBIT_AUTH_ACCESS_TOKEN_TTL` 存在时必须能 parse 为 positive Go duration。
- `VIBIT_AUTH_TOKEN_AUDIENCE` 必须 trim whitespace；缺失或为空时回退到 `vibit_gameplay_runtime_requests`。
- Startup 生成的 token record ids 必须是非 secret identifiers。
- Raw verifier keys、access tokens、credential proof、lookup digests、verifier digests 和 full concrete verifier key ids 不得出现在日志或 startup errors 中。
- 本 gate 不授权外部 UUID、ULID、KSUID、JWT、OAuth、OIDC、password-hashing、Redis-like、KMS、cloud secret-manager 或 session-store dependency。

## 6. Nakama 和 Pitaya 参考映射

Nakama 指导能力需求：client 先 authenticate 获取 session/token，然后再使用 server 和 realtime features。vibit 把这个思路适配为：在 PostgreSQL runtime path 中把 token validation 组合进 request-level route protection。

Pitaya 指导架构分层：acceptors/connections、sessions、routing 和 handlers 是不同关注点。vibit 把这个思路适配为：WebSocket transport 保持 credential-neutral，在 domain handlers 执行前注入 validated identity。

两个参考都不会覆盖 vibit 自身边界。本 gate 不采用直接 Nakama 或 Pitaya public API compatibility。

## 7. 必需实现测试

Implementation slice 必须在以下文件下添加或更新 focused tests：

```text
runtime/cmd/vibit-server/main_test.go
```

必需测试类别：

```yaml
required_tests:
  memory_startup_remains_bootstrap_without_route_protector
  explicit_route_protector_can_be_injected_into_frame_handler
  postgres_auth_startup_requires_verifier_key_configuration
  auth_startup_accepts_default_lifetime_and_audience
  auth_startup_accepts_configured_lifetime_and_audience
  auth_startup_rejects_invalid_lifetime
  token_record_id_generator_returns_non_secret_stable_shape
  startup_errors_do_not_include_verifier_key_values
```

Live PostgreSQL startup verification 仍然是 optional，不得成为默认 repository checks 的要求。

## 8. 延期项

本 gate 不授权：

- WebSocket handshake authentication。
- Session persistence。
- Runtime session tables 或 migrations。
- Existing payload wrapper 之外的 authentication command Protobuf messages。
- Login route registration。
- HTTP `Authorization`、Bearer、cookie、query string 或 WebSocket subprotocol credential carriers。
- Repository interface changes。
- PostgreSQL adapter changes。
- Migration changes。
- Automatic startup migrations。
- Logout、refresh、cleanup、token rotation、token validation audit mutation 或 previous-token revocation。
- External dependencies。
- Direct Nakama 或 Pitaya API compatibility。

## 9. Verification

本 gate 的 repository check rule 是：

```text
runtime.authentication_startup_composition_gate
```
