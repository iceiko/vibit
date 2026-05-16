# 认证服务行为实现 Gate

状态: Draft v0.1
最后更新: 2026-05-16
范围: 在添加服务行为代码之前，定义未来应用层认证服务行为的所有权、文件边界、仓储交接、helper 组合流程、公开错误折叠、脱敏、测试与延期项
依赖: `docs/runtime-authentication-implementation-boundary.md`, `docs/application-authentication-service-interface-boundary.md`, `docs/token-credential-verifier-algorithm-redaction-boundary.md`, `docs/secret-configuration-verifier-key-loading-boundary.md`, `docs/token-credential-material-generation-boundary.md`, `docs/verifier-digest-computation-comparison-boundary.md`, `docs/authentication-service-implementation-readiness-gate.md`, `docs/verifier-digest-comparison-helper-gate.md`
规范决策: `ADR-0050`

对应的英文源文档是 `docs/authentication-service-behavior-implementation-gate.md`。英文文件是权威版本。

## 1. 目的

这个 gate 定义 verifier key、材料生成、digest 计算、verifier 比较 helper 已存在之后，未来认证服务行为的实现边界。

仓库现在已经具备第一种已选认证姿态所需的小型应用层 helper 链：

```text
VerifierKeySet
-> raw credential/token material generation
-> lookup and verifier digest computation
-> constant-time verifier digest comparison
```

下一类风险是未来 Agent 把这些 helper 直接接进 transport、Protobuf adapter、repository、startup 或 generated files，或者在实现登录和 token 校验时顺手发明错误映射、仓储调用顺序、proof 脱敏、key 选择或 request identity 交接方式。

这只是 implementation-gate standard。它不添加 Go service code、登录执行、access-token 校验、logout 执行、refresh 行为、cleanup 任务、协议承载、仓储接口变更、SQL 迁移、启动 wiring、外部依赖或生产认证行为。

## 2. 核心规则

认证服务行为实现 gate 是：

```yaml
authentication_service_behavior_implementation_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0107
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
future_slice: authentication_service_behavior_skeleton
future_source: runtime/internal/app/authentication/service.go
future_tests: runtime/internal/app/authentication/service_test.go
repository_handoff: application_unit_of_work
helper_composition_flow_defined: true
public_error_collapse_defined: true
request_identity_handoff_defined: true
service_behavior_status: gated
login_execution_status: deferred
token_validation_status: deferred
logout_execution_status: deferred
refresh_behavior_status: deferred
cleanup_execution_status: deferred
protocol_carrier_status: deferred
protobuf_authentication_messages_status: deferred
websocket_proof_carrier_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
startup_wiring_status: deferred
authentication_dependencies_status: deferred
production_authentication_behavior_status: deferred
```

未来认证服务行为必须归应用层所有。只有后续有边界清晰的 work item 明确授权某个切片时，它才可以编排现有 helper 输出、repository interface 和 request identity 交接。

## 3. 未来服务所有权

未来服务行为所有者：

```text
runtime/internal/app/authentication
```

后续实现 work item 授权代码之后，允许的未来文件：

```text
runtime/internal/app/authentication/service.go
runtime/internal/app/authentication/service_test.go
```

第一个未来服务代码切片应该只是 skeleton，除非后续 work item 明确授权真实登录或 token 校验行为。这个 skeleton 可以定义类型化 dependency 边界、request/result 词汇、脱敏后的内部错误类别，以及 fail-closed 的 `not implemented` 行为。它不得调用仓储、签发 token、校验 token、撤销 token、refresh token、清理 token、暴露协议承载或 wiring 启动。

## 4. 仓储交接

认证服务行为只能通过应用层 unit-of-work 边界使用 storage-neutral 的 authentication repository。

未来状态变更行为的必需流程：

```text
application service method
-> application-owned unit-of-work runner
-> UnitOfWork.NewAuthenticationRepository(...)
-> authentication.Repository
-> persistence-only PostgreSQL adapter
```

规则：

- service 负责编排。
- repository 只负责 storage-neutral 的记录查找与变更。
- PostgreSQL adapter 只负责 SQL 持久化。
- service 不得导入 PostgreSQL driver package。
- service 不得绕过 unit-of-work 边界执行会改变状态的认证行为。
- repository 不得生成材料、计算 digest、比较 verifier、解析 proof、映射公开失败或构造 `RequestIdentity`。
- PostgreSQL adapter 不得决定认证结果。

仓储接口变更仍然延期。如果 service 需要当前不存在的仓储方法，后续 work item 必须先更新模块边界和 adapter 测试，再让行为代码消费它。

## 5. Helper 组合流程

未来 device credential login 行为在被单独授权时，必须按这个顺序组合 helper：

```yaml
device_credential_login_flow:
  input: already_decoded_service_request
  proof_shape: raw_device_credential_material
  steps:
    - reject_missing_or_malformed_proof_before_repository_lookup
    - compute_credential_lookup_digest_with_VerifierKeySet
    - find_credential_by_lookup_digest_through_unit_of_work_repository
    - collapse_lookup_miss_or_unusable_record_to_public_invalid_credential
    - compute_credential_verifier_digest_with_record_verifier_key_context
    - compare_credential_verifier_digest_with_CompareCredentialVerifierDigest
    - collapse_mismatch_to_public_invalid_credential
    - require_active_credential_and_allowed_player_account_state
    - generate_access_token_material_with_explicit_entropy_reader
    - compute_token_lookup_and_verifier_digests
    - store_token_record_through_unit_of_work_repository
    - return_raw_access_token_text_once_to_the_authorized_response_carrier
  output: redacted_authentication_result
```

未来 access-token validation 行为在被单独授权时，必须按这个顺序组合 helper：

```yaml
access_token_validation_flow:
  input: already_decoded_explicit_request_proof_payload
  proof_shape: raw_access_token_material
  steps:
    - reject_missing_or_malformed_proof_before_repository_lookup
    - compute_token_lookup_digest_with_VerifierKeySet
    - find_token_by_lookup_digest_through_unit_of_work_repository
    - collapse_lookup_miss_or_unusable_record_to_public_invalid_token
    - compute_token_verifier_digest_with_record_verifier_key_context
    - compare_token_verifier_digest_with_CompareTokenVerifierDigest
    - collapse_mismatch_to_public_invalid_token
    - require_active_token_lifecycle_state_and_unexpired_window
    - convert_validated_actor_to_RequestIdentity
  output: application_owned_request_identity
```

这个 gate 只记录预期组合方式。它不授权任何一个流程执行。

## 6. 公开错误折叠

服务行为必须保留对测试有用的内部 proof 区分，同时折叠公开 proof 失败。

内部失败类别可以包括：

```yaml
internal_failure_classes:
  missing_proof: redacted
  malformed_proof: redacted
  lookup_miss: redacted
  wrong_verifier_algorithm: redacted
  unknown_verifier_key_id: redacted
  unsupported_verifier_version: redacted
  verifier_digest_mismatch: redacted
  credential_not_active: redacted
  token_not_active: redacted
  token_expired: redacted
  token_revoked: redacted
  repository_unavailable: redacted_dependency
```

第一姿态要求的公开折叠：

```yaml
public_error_collapse:
  missing_device_credential_proof: AUTHENTICATION_PROOF_MISSING
  malformed_device_credential_proof: AUTHENTICATION_PROOF_MALFORMED
  invalid_device_credential_proof_family: AUTHENTICATION_CREDENTIAL_INVALID
  missing_access_token_proof: AUTHENTICATION_TOKEN_MISSING
  malformed_access_token_proof: AUTHENTICATION_TOKEN_MALFORMED
  invalid_access_token_proof_family: AUTHENTICATION_TOKEN_INVALID
  credential_store_unavailable: AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE
  token_store_unavailable: AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
  unsupported_refresh: AUTHENTICATION_REFRESH_NOT_SUPPORTED
  not_implemented: AUTHENTICATION_NOT_IMPLEMENTED
```

`AUTHENTICATION_TOKEN_EXPIRED`、`AUTHENTICATION_TOKEN_REVOKED` 和 account-disabled 的公开区分已经存在于语义 catalog 中，但第一版行为是否公开这些区分仍然是后续显式决策。未来 work item 只有确认披露姿态之后，才可以授权更具体的公开映射。

## 7. Request Identity 交接

Access-token validation 最终必须在生产敏感的 domain dispatch 之前，把 proof 转换成应用层拥有的 request identity。

目标交接：

```yaml
request_identity_handoff:
  owner: runtime/internal/app
  source: access_token_validation_result
  target_type: RequestIdentity
  success_status: authentication_proven
  actor_kind: player
  player_id_validated: true
  session_validated: false_until_session_persistence_gate
  metadata_only_allowed_as_proof: false
```

Domain modules 必须接收 `RequestIdentity`。它们不得解析 access token、比较 verifier digest、选择 credential record 或决定认证 proof 是否有效。

## 8. 脱敏要求

以下内容禁止出现在 errors、logs、traces、metrics labels、test snapshots、ADRs、change specs、documentation examples、conversation logs 以及 public responses 中，除非是一次性授权 token response carrier：

- Raw device credential material。
- Raw access-token material。
- Encoded generated credential material。
- 一次性 response carrier 之外的 encoded generated access-token material。
- Lookup digest bytes。
- Verifier digest bytes。
- HMAC input bytes。
- HMAC output bytes。
- Verifier key values。
- Encoded verifier key values。
- 完整具体的 `verifier_key_id` 值。
- Credential lookup hit 或 miss 细节。
- Token lookup hit 或 miss 细节。
- Candidate key-set counts。
- failed proof 的 repository 实现细节。

允许：

- 已注册的公开 error codes。
- 在流程已经证明 actor 且目标 carrier 已授权时的非秘密 record ids。
- 诸如 `<raw-access-token>` 和 `<verifier-key-id>` 的脱敏占位符。
- 不向公开面暴露 secret 或 existence detail 的内部测试失败类别名。

## 9. 文件边界

允许的未来实现区域：

```text
runtime/internal/app/authentication/service.go
runtime/internal/app/authentication/service_test.go
```

后续 service code work item 授权之后，可以依赖的现有 helper：

```text
runtime/internal/app/authentication/verifier_key_config.go
runtime/internal/app/authentication/verifier_key_env.go
runtime/internal/app/authentication/material_generation.go
runtime/internal/app/authentication/verifier_digest.go
runtime/internal/app/authentication/verifier_comparison.go
```

除非后续 work item 明确命名，否则第一个 service behavior slice 禁止写入：

- `runtime/cmd/vibit-server/`
- `runtime/internal/app/bootstrap/`
- `runtime/internal/platform/transport/ws/`
- `runtime/internal/platform/protocol/protobuf/`
- `runtime/internal/platform/persistence/postgres/`
- `runtime/internal/platform/migrations/`
- `runtime/internal/modules/authentication/`
- `runtime/internal/generated/`
- `runtime/migrations/postgres/`
- `proto/`
- `contracts/runtime/authentication/`

未来服务行为不得手工编辑 generated files，也不得把服务行为藏在 transport、protocol、repository、migration 或 startup code 中。

## 10. 未来服务行为所需测试

未来实现必须在以下位置添加聚焦测试：

```text
runtime/internal/app/authentication/service_test.go
```

真实行为被接受之前，至少需要这些测试类别：

```yaml
required_tests:
  service_skeleton_fails_closed_without_behavior_authorization
  service_dependencies_reject_nil_or_missing_unit_of_work
  service_does_not_log_or_return_raw_proof_material
  credential_login_composes_lookup_digest_before_repository_lookup
  credential_login_compares_verifier_digest_before_token_issuance
  credential_login_collapses_lookup_miss_and_mismatch_to_same_public_error
  access_token_validation_composes_lookup_digest_before_repository_lookup
  access_token_validation_compares_verifier_digest_before_request_identity
  access_token_validation_collapses_lookup_miss_and_mismatch_to_same_public_error
  request_identity_is_populated_only_after_valid_access_token_proof
  repository_is_used_only_through_unit_of_work_boundary
  protocol_and_transport_packages_do_not_import_service_behavior
```

除非后续 live integration work item 明确选择 PostgreSQL，否则测试必须使用 fake repository 和 fake unit-of-work。普通 service behavior 测试不得要求运行中的 PostgreSQL server。

## 11. Nakama 与 Pitaya 映射

需要吸收的 Nakama 能力：

- 服务端账户认证。
- session/access-token 签发。
- session/access-token 校验。
- token 过期和撤销检查。
- 账号 disabled 检查。

需要吸收的 Pitaya 能力：

- Frontend acceptor 与 backend handler logic 分离。
- Realtime handler 在校验之后接收 identity/session context。
- Session binding 是 runtime context，本身不是 proof。

vibit 规则：

```text
transport accepts frames
-> protocol decodes messages
-> application authentication validates proof
-> application request identity carries validated actor context
-> domain modules consume request identity
```

不要直接复制 Nakama 或 Pitaya 的公开 API。把它们作为能力和词汇参考，同时保留 vibit 的 agent-native 边界。

## 12. 延期项

这个 gate 不授权：

- Go authentication service behavior code。
- 登录执行。
- Access-token validation 执行。
- Logout 执行。
- Refresh 执行。
- Token cleanup jobs。
- Credential bootstrap 或 account creation policy。
- Session persistence。
- Protocol authentication messages。
- WebSocket proof carriers。
- WebSocket handshake authentication。
- Repository interface changes。
- PostgreSQL adapter changes。
- SQL migration changes。
- Startup wiring。
- External dependencies。
- KMS 或 cloud secret-manager integration。
- Redis-like token/session stores。
- Production authentication behavior。

## 13. 验证

这个 gate 的仓库检查规则是：

```text
runtime.authentication_service_behavior_implementation_gate
```

触碰这个 gate 的变更需要运行：

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

检查必须保持：除非后续 work item 明确授权 service code，否则 `service.go` 和 `service_test.go` 不存在。
