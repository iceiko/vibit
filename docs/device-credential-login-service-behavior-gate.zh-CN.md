# 设备凭证登录服务行为 Gate

状态: Draft v0.1
最后更新: 2026-05-16
范围: 在添加真实登录执行之前，定义未来设备凭证登录服务行为顺序、dependency shape、仓储交接、token 签发姿态、公开失败折叠、脱敏、测试和延期项
依赖: `docs/authentication-service-behavior-implementation-gate.md`, `docs/application-authentication-service-interface-boundary.md`, `docs/token-credential-material-generation-implementation-gate.md`, `docs/verifier-digest-helper-implementation-gate.md`, `docs/verifier-digest-comparison-helper-gate.md`
规范决策: `ADR-0051`

对应的英文源文档是 `docs/device-credential-login-service-behavior-gate.md`。英文文件是权威版本。

## 1. 目的

认证服务 skeleton 已经存在，并且 fail closed。下一类风险是未来 Agent 直接在 skeleton method 中填入临时仓储调用、token 生成、公开错误映射或 transport 假设。

这个 gate 在代码允许执行真实登录之前，定义精确的未来登录行为。

这只是 gate-only standard。它不实现 device credential login、不签发 access token、不 validate access tokens、不改变 service method signatures、不暴露 protocol carriers、不改变 repositories、不改变 migrations、不 wire startup、不添加 dependencies，也不添加 production authentication behavior。

## 2. 核心规则

设备凭证登录服务行为 gate 是：

```yaml
device_credential_login_service_behavior_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0109
completed_gate_work_item: W-0108
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
planned_source: runtime/internal/app/authentication/service.go
planned_tests: runtime/internal/app/authentication/service_test.go
service_method: AuthenticateWithDeviceCredential
login_method: device_credential_login
credential_kind: device_credential_login
token_kind: access_token
token_type: opaque_access_token
proof_carrier_status: already_decoded_service_request_only
protocol_carrier_status: deferred
startup_wiring_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
authentication_dependencies_status: deferred
production_authentication_behavior_status: deferred
```

只有后续 work item 明确授权登录流程执行时，未来实现才可以执行该流程。

## 3. 未来 Dependency Shape

未来实现只能在 application service 边界内扩展现有 `ServiceDependencies` shape。

未来必需 dependency 类别：

```yaml
future_service_dependencies:
  unit_of_work_runner: already_present
  verifier_key_set: required
  access_token_entropy_reader: required
  clock: required
  token_record_id_generator: required
  access_token_lifetime: required_positive_duration
  token_audience: required_non_empty_string
```

规则：

- `UnitOfWorkRunner` 仍然是唯一 transaction entry point。
- 未来登录行为必须在传入的 `tx.UnitOfWork` 上使用局部 capability interface 获取 `NewAuthenticationRepository()` 和 `NewPlayerAccountRepository()`。
- 不得仅为了这个登录切片扩展全局 `tx.UnitOfWork` interface。
- Token record id generation 必须通过 dependency 注入。本 gate 不选择 UUID、ULID、KSUID、database-generated ids 或外部 id package。
- Access-token lifetime 必须通过 service dependencies 配置，并且必须为正数。
- 本 gate 不批准任何 production default token lifetime。
- 本 gate 不 wire startup configuration。

## 4. Proof 输入形态

未来登录方法通过 `DeviceCredentialAuthenticationRequest.CredentialProof` 接收 proof。

第一版 proof 规则：

```yaml
device_credential_proof:
  source: already_decoded_service_request
  text_encoding: base64url_unpadded
  encoded_length_chars: 43
  raw_length_bytes: 32
  raw_entropy_floor_bits: 256
  bearer_prefix: forbidden
  raw_device_identifier: forbidden
  client_generated_low_entropy_value: forbidden
```

规则：

- Proof 是 server-issued high-entropy device credential material。
- 它不是 raw operating-system device id、advertising id、hardware serial number、user name、email address、provider subject、session id 或 transport metadata。
- Missing、whitespace-only、padded、wrongly sized、non-Base64URL 或不是 32-byte decoded proof 的输入，必须在任何 unit-of-work 或 repository call 之前失败。
- 未来实现只应在计算 digest 的短暂阶段把文本 decode 为 raw material，并且必须避免它出现在 logs、public errors、test snapshots、conversation logs 和 docs examples 中。
- Application service 不得为这个 login method 解析 HTTP `Authorization` header 或 `Bearer` string。Protocol carriers 仍然延期。

## 5. 必需登录顺序

当 W-0109 或后续 work item 授权行为时，`AuthenticateWithDeviceCredential` 必须按以下顺序执行：

```yaml
device_credential_login_sequence:
  - reject_missing_or_malformed_device_credential_proof_before_unit_of_work
  - decode_device_credential_proof_text_to_raw_32_byte_material
  - compute_credential_lookup_digest_with_active_VerifierKeySet
  - enter_application_unit_of_work
  - obtain_authentication_repository_from_unit_of_work_capability
  - obtain_player_account_repository_from_unit_of_work_capability
  - find_credential_by_lookup_digest
  - collapse_lookup_miss_to_public_invalid_credential
  - require_credential_kind_device_credential_login
  - require_credential_status_active
  - require_supported_verifier_algorithm_vibit_hmac_sha256_v1
  - require_supported_verifier_version_1
  - require_active_verifier_key_set_id_match_for_first_posture
  - compute_credential_verifier_digest
  - compare_credential_verifier_digest_with_CompareCredentialVerifierDigest
  - collapse_verifier_mismatch_to_public_invalid_credential
  - get_player_account_by_credential_player_id
  - require_player_account_active
  - generate_access_token_material_with_explicit_entropy_reader
  - compute_token_lookup_digest_with_active_VerifierKeySet
  - compute_token_verifier_digest_with_active_VerifierKeySet
  - create_token_record_id_with_injected_generator
  - store_access_token_record_through_authentication_repository
  - exit_unit_of_work_successfully
  - return_raw_access_token_text_once_after_unit_of_work_success
```

规则：

- 缺失或 malformed proof 被拒绝之前，不得发生 repository call。
- Credential lookup 成功之前，不得发生 verifier digest comparison。
- Credential proof 和 player account state 被接受之前，不得生成 access-token material。
- Token storage、transaction commit 或任何 dependency 失败时，不得返回 raw access token。
- Raw access token text 只能在 unit of work 成功之后，通过 `AuthenticationResult.AccessToken` 一次性返回。
- Lookup digest equality 本身不是 proof。
- Database equality 本身不是 final proof。
- 第一版登录行为要求 player account active。
- Credential disabled、revoked、replaced、wrong algorithm、wrong version、wrong key id、player missing、player disabled、player deleted、lookup miss 和 verifier mismatch，都折叠到同一个公开 invalid-credential family，除非后续披露决策改变这一点。

## 6. 仓储交接

未来登录行为必须使用现有 repository interfaces。本 gate 不授权 repository interface changes。

未来必需 repository access：

```yaml
repository_handoff:
  transaction_boundary: UnitOfWorkRunner.WithinUnitOfWork
  authentication_repository_source: unit_of_work.NewAuthenticationRepository
  player_repository_source: unit_of_work.NewPlayerAccountRepository
  authentication_lookup_method: FindCredentialByLookupDigest
  token_store_method: StoreToken
  player_lookup_method: GetPlayerAccount
  direct_postgres_import: forbidden
  repository_interface_change: forbidden_by_this_gate
```

Service 不得为登录行为 import PostgreSQL adapters、pgx、goose、WebSocket transport packages、Protobuf adapter packages、generated Protobuf packages 或 generated contract-shape packages。

Authentication module 保持 storage-neutral。它只存储和读取 already-computed digest records。它不得 generate access tokens、compute digests、compare verifiers、decide proof validity、collapse public failures 或 construct application responses。

## 7. Token 签发姿态

第一版未来登录行为只能签发 opaque access tokens。

```yaml
token_issuance:
  token_kind: access_token
  token_type: opaque_access_token
  actor_kind: player
  token_material_generation: GenerateAccessTokenMaterial
  token_lookup_digest_helper: ComputeTokenLookupDigest
  token_verifier_digest_helper: ComputeTokenVerifierDigest
  token_record_store_method: StoreToken
  refresh_token_issued: false
  jwt_or_signed_claim_token: forbidden
  previous_token_rotation: deferred
  cleanup_job: deferred
```

规则：

- `StoreTokenMutation` 只能收到 digest bytes，绝不能收到 raw token material。
- `VerifierAlgorithm` 必须是 `vibit_hmac_sha256_v1`。
- 第一版实现的 `VerifierVersion` 必须是 `1`。
- `VerifierKeyID` 必须来自 active `VerifierKeySet.KeySetID()`。
- `IssuedAt` 和 `ExpiresAt` 必须来自 injected clock 和 positive access-token lifetime。
- `TokenRecordID` 必须来自 injected token record id generator。
- `Audience` 必须来自 service configuration。
- Refresh-token behavior 仍然 unsupported。
- 除非后续 work item 授权 exact behavior，否则同一 credential 下 previous tokens 的 token rotation 或 revocation 继续延期。

## 8. 公开错误折叠

未来行为必须保留 redacted internal classes，同时折叠公开 proof failures。

```yaml
public_error_collapse:
  missing_proof: AUTHENTICATION_PROOF_MISSING
  malformed_proof: AUTHENTICATION_PROOF_MALFORMED
  lookup_miss: AUTHENTICATION_CREDENTIAL_INVALID
  wrong_credential_kind: AUTHENTICATION_CREDENTIAL_INVALID
  inactive_credential: AUTHENTICATION_CREDENTIAL_INVALID
  unsupported_algorithm: AUTHENTICATION_CREDENTIAL_INVALID
  unsupported_version: AUTHENTICATION_CREDENTIAL_INVALID
  verifier_key_id_mismatch: AUTHENTICATION_CREDENTIAL_INVALID
  verifier_mismatch: AUTHENTICATION_CREDENTIAL_INVALID
  player_missing_or_inactive: AUTHENTICATION_CREDENTIAL_INVALID
  repository_or_unit_of_work_unavailable: AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE
  token_generation_or_storage_unavailable: AUTHENTICATION_CREDENTIAL_STORE_UNAVAILABLE
```

公开面不得暴露 credential record 是否存在、player 是否存在、stored verifier key id 是什么、proof accepted 之后 token store 是否失败，或者 verifier mismatch 是否发生。

## 9. 脱敏要求

未来行为不得把以下值放进 errors、logs、docs examples、test failure messages、conversation logs、ADRs、change specs 或 public responses：

- Raw device credential text。
- Raw device credential bytes。
- Raw access-token text，唯一例外是成功时一次性 `AuthenticationResult.AccessToken` carrier。
- Raw access-token bytes。
- Lookup digest bytes。
- Verifier digest bytes。
- HMAC input 或 output bytes。
- Verifier key bytes。
- 完整具体的 `verifier_key_id` 值。
- Credential lookup hit 或 miss details。
- Player lookup hit 或 miss details。
- failed proof 的 token store internals。

允许：

- 已注册的 public error codes。
- `<device-credential-proof>` 和 `<access-token>` 这类 redacted placeholders。
- 只有 proof 成功后，才能在授权 result fields 或 internal tests 中使用 non-secret record ids。

## 10. 必需测试

未来实现必须在 `runtime/internal/app/authentication/service_test.go` 中添加或更新 focused tests。

必需 test classes：

```yaml
required_tests:
  login_rejects_missing_proof_without_unit_of_work
  login_rejects_malformed_proof_without_unit_of_work
  login_computes_lookup_digest_before_repository_lookup
  login_uses_authentication_repository_from_unit_of_work_only
  login_uses_player_repository_from_unit_of_work_only
  login_collapses_lookup_miss_to_invalid_credential
  login_rejects_inactive_or_wrong_kind_credential
  login_rejects_wrong_algorithm_version_or_key_id
  login_compares_verifier_digest_before_token_generation
  login_collapses_verifier_mismatch_to_invalid_credential
  login_requires_active_player_account
  login_generates_access_token_only_after_proof_and_player_acceptance
  login_stores_token_digest_only
  login_does_not_return_access_token_when_store_or_commit_fails
  login_returns_access_token_once_after_unit_of_work_success
  login_errors_do_not_leak_raw_proof_or_token_material
  login_does_not_validate_access_tokens_or_touch_protocol_carriers
```

普通测试必须使用 fake unit-of-work 和 repository implementations。它们不得要求 live PostgreSQL server。

## 11. Nakama 与 Pitaya 映射

Nakama capability mapping：

- Server-side account authentication 被适配为 application-owned proof validation。
- Session/access-token issuance 被适配为 opaque access-token generation 和 digest storage。
- Token expiration 和 revocation 仍属于 token validation behavior，不属于 login behavior。

Pitaya capability mapping：

- Realtime handlers 应该在 validation 之后接收 identity context。
- Frontend 和 backend handler separation 映射到 vibit 的 transport/protocol/application separation。
- Login proof validation 不属于 route/domain handlers。

这些参考只指导 capability coverage。它们不授权复制 Nakama 或 Pitaya public APIs。

## 12. 延期项

本 gate 不授权：

- Device credential login execution。
- Access-token validation execution。
- Logout execution。
- Refresh execution。
- Token cleanup jobs。
- Protocol authentication messages。
- WebSocket proof carriers。
- WebSocket handshake authentication。
- Startup wiring。
- Repository interface changes。
- PostgreSQL adapter changes。
- SQL migration changes。
- Generated file changes。
- External cryptography、JWT、OAuth、OIDC、provider、KMS、cloud secret-manager、password-hashing、Redis-like、queue 或 session-store dependencies。
- Production authentication behavior。

## 13. 验证

该 gate 的 repository check rule 是：

```text
runtime.device_credential_login_service_behavior_gate
```
