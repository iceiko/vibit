# Access Token Validation 服务行为 Gate

状态: Draft v0.1
最后更新: 2026-05-16
范围: 在添加真实 validation execution 之前，定义未来 access-token validation service behavior sequence、proof input shape、repository handoff、token lifecycle checks、request identity handoff、public failure collapse、redaction、tests 和 deferrals
依赖: `docs/authentication-service-behavior-implementation-gate.md`, `docs/application-authentication-service-interface-boundary.md`, `docs/verifier-digest-helper-implementation-gate.md`, `docs/verifier-digest-comparison-helper-gate.md`, `docs/device-credential-login-service-behavior-gate.md`
规范决策: `ADR-0052`

对应的英文源文档是 `docs/access-token-validation-service-behavior-gate.md`。英文文件是权威版本。

## 1. 目的

Device credential login 现在已经可以通过 application authentication service 签发 opaque access tokens。下一类风险是未来 Agent 直接在 runtime dispatch 或 transport code 里临时加入 bearer parsing、repository calls、route protection、session persistence 或公开错误披露。

这个 gate 在代码允许执行真实 access-token validation 之前，定义精确的未来验证行为。

这只是 gate-only standard。它不实现 access-token validation、不改变 service method signatures、不暴露 protocol carriers、不改变 repositories、不改变 migrations、不 wire startup、不添加 session persistence、不添加 dependencies，也不添加 production authentication behavior。

## 2. 核心规则

Access-token validation 服务行为 gate 是：

```yaml
access_token_validation_service_behavior_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0111
completed_gate_work_item: W-0110
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
planned_source: runtime/internal/app/authentication/service.go
planned_tests: runtime/internal/app/authentication/service_test.go
service_method: ValidateAccessToken
token_kind: access_token
token_type: opaque_access_token
proof_carrier_status: already_decoded_service_request_only
protocol_carrier_status: deferred
startup_wiring_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
session_persistence_status: deferred
authentication_dependencies_status: deferred
production_authentication_behavior_status: deferred
```

只有后续 work item 明确授权 token validation 时，未来实现才可以执行该流程。

## 3. 未来 Dependency Shape

未来实现应复用为 login service slice 创建的现有 `ServiceDependencies` shape。

未来必需 dependency 类别：

```yaml
future_service_dependencies:
  unit_of_work_runner: already_present
  verifier_key_set: already_present
  clock: already_present
  token_audience: already_present
  access_token_entropy_reader: not_used_by_validation
  token_record_id_generator: not_used_by_validation
  access_token_lifetime: not_used_by_validation
```

规则：

- `UnitOfWorkRunner` 仍然是唯一 transaction entry point。
- 未来 validation behavior 必须在传入的 `tx.UnitOfWork` 上使用局部 capability interface 获取 `NewAuthenticationRepository()` 和 `NewPlayerAccountRepository()`。
- 不得仅为了这个 validation slice 扩展全局 `tx.UnitOfWork` interface。
- Validation 不得生成新的 token material 或 token record id。
- Injected clock 只用于检查 token time windows。
- Configured token audience 必须匹配 stored token audience。
- 本 gate 不 wire startup configuration。

## 4. Proof 输入形态

未来验证方法通过 `AccessTokenValidationRequest.AccessToken` 接收 proof。

第一版 proof 规则：

```yaml
access_token_proof:
  source: already_decoded_service_request
  text_encoding: base64url_unpadded
  encoded_length_chars: 43
  raw_length_bytes: 32
  raw_entropy_floor_bits: 256
  bearer_prefix: forbidden
  authorization_header_parsing: forbidden
  cookie_parsing: forbidden
  query_string_parsing: forbidden
  session_id_as_token: forbidden
```

规则：

- Access token 是由 service login flow 签发的 opaque high-entropy material。
- 它不是 JWT、signed claim token、session id、WebSocket connection id、device id、player id、route field 或 transport metadata。
- Missing、whitespace-only、padded、wrongly sized、non-Base64URL 或不是 32-byte decoded proof 的输入，必须在任何 unit-of-work 或 repository call 之前失败。
- 未来实现只应在计算 digest 的短暂阶段把文本 decode 为 raw material，并且必须避免它出现在 logs、public errors、test snapshots、conversation logs 和 docs examples 中。
- Application service 不得解析 HTTP `Authorization` header、`Bearer` string、cookie、query parameter、WebSocket handshake field 或 Protobuf authentication carrier。Protocol carriers 仍然延期。

## 5. 必需验证顺序

当 W-0111 或后续 work item 授权行为时，`ValidateAccessToken` 必须按以下顺序执行：

```yaml
access_token_validation_sequence:
  - reject_missing_or_malformed_access_token_before_unit_of_work
  - decode_access_token_text_to_raw_32_byte_material
  - compute_token_lookup_digest_with_active_VerifierKeySet
  - enter_application_unit_of_work
  - obtain_authentication_repository_from_unit_of_work_capability
  - obtain_player_account_repository_from_unit_of_work_capability
  - find_token_by_lookup_digest
  - collapse_lookup_miss_to_public_invalid_token
  - require_token_kind_access_token
  - require_token_status_active
  - require_supported_verifier_algorithm_vibit_hmac_sha256_v1
  - require_supported_verifier_version_1
  - require_active_verifier_key_set_id_match_for_first_posture
  - require_configured_token_audience_match
  - require_token_issued_at_not_in_future_beyond_clock_tolerance
  - require_token_not_expired_by_injected_clock
  - compute_token_verifier_digest
  - compare_token_verifier_digest_with_CompareTokenVerifierDigest
  - collapse_verifier_mismatch_to_public_invalid_token
  - get_player_account_by_token_player_id
  - require_player_account_active
  - construct_validated_player_RequestIdentity
  - exit_unit_of_work_successfully
  - return_validated_identity_after_unit_of_work_success
```

规则：

- 缺失或 malformed proof 被拒绝之前，不得发生 repository call。
- Token lookup 成功之前，不得发生 verifier digest comparison。
- Token verifier comparison、token lifecycle checks、audience check 和 player account state checks 成功之前，不得把 request identity 标记为 validated。
- Lookup digest equality 本身不是 proof。
- Database equality 本身不是 final proof。
- 第一版 validation behavior 要求 player account active。
- 第一姿态只检查 active `VerifierKeySet.KeySetID()`。Previous-key validation 仍然延期。
- Token disabled、expired、revoked、replaced、wrong kind、wrong algorithm、wrong version、wrong key id、wrong audience、future-issued token、player missing、player disabled、player deleted、lookup miss 和 verifier mismatch 都折叠到同一个公开 invalid-token family，除非后续披露决策改变这一点。

## 6. 仓储交接

未来 validation behavior 必须使用现有 repository interfaces。本 gate 不授权 repository interface changes。

未来必需 repository access：

```yaml
repository_handoff:
  transaction_boundary: UnitOfWorkRunner.WithinUnitOfWork
  authentication_repository_source: unit_of_work.NewAuthenticationRepository
  player_repository_source: unit_of_work.NewPlayerAccountRepository
  authentication_lookup_method: FindTokenByLookupDigest
  token_mutation_method: deferred
  player_lookup_method: GetPlayerAccount
  direct_postgres_import: forbidden
  repository_interface_change: forbidden_by_this_gate
```

Service 不得为 validation behavior import PostgreSQL adapters、pgx、goose、WebSocket transport packages、Protobuf adapter packages、generated Protobuf packages 或 generated contract-shape packages。

Authentication module 保持 storage-neutral。它只存储和读取 already-computed digest records。它不得 parse access-token carriers、compute digests、compare verifiers、decide proof validity、collapse public failures 或 construct application responses。

更新 `LastValidatedAt` 或 `LastFailedValidationAt` 仍然延期，除非后续 work item 授权 repository mutation shape 和 storage behavior。

## 7. Request Identity Handoff

成功验证必须在 production-sensitive domain dispatch 之前创建 application-owned request identity。

```yaml
request_identity_handoff:
  owner: runtime/internal/app
  source: AccessTokenValidationResult
  target_type: RequestIdentity
  success_validation_status: validated
  success_proof_status: valid
  identity_status: validated
  actor_kind: player
  actor_id: token.player_id
  player_id: token.player_id
  player_id_validated: true
  session_validated: false_until_session_persistence_gate
  session_id_source: request_metadata_only_if_present
  connection_id_source: request_metadata_only_if_present
  connection_epoch_source: request_metadata_only_if_present
  metadata_only_allowed_as_proof: false
```

规则：

- 未来代码必须 keep SessionValidated false until session persistence is ratified；不得仅凭 access-token validation 就把 `SessionValidated` 标记为 true。
- 如果使用 application identity helpers，必须覆盖或避开任何把 metadata-only session id 视为 validated proof 的 helper 行为。
- `ConnectionID`、`ConnectionEpoch`、route fields 和 session metadata 可以作为 context 复制，但它们不是 proof。
- Domain modules 必须消费 `RequestIdentity`；它们不得 parse access tokens、compare verifier digests、select token records 或 decide authentication proof validity。

## 8. 公开错误折叠

未来行为必须保留 redacted internal classes，同时折叠公开 proof failures。

```yaml
public_error_collapse:
  missing_token: AUTHENTICATION_TOKEN_MISSING
  malformed_token: AUTHENTICATION_TOKEN_MALFORMED
  lookup_miss: AUTHENTICATION_TOKEN_INVALID
  wrong_token_kind: AUTHENTICATION_TOKEN_INVALID
  inactive_token: AUTHENTICATION_TOKEN_INVALID
  expired_token: AUTHENTICATION_TOKEN_INVALID
  revoked_token: AUTHENTICATION_TOKEN_INVALID
  unsupported_algorithm: AUTHENTICATION_TOKEN_INVALID
  unsupported_version: AUTHENTICATION_TOKEN_INVALID
  verifier_key_id_mismatch: AUTHENTICATION_TOKEN_INVALID
  audience_mismatch: AUTHENTICATION_TOKEN_INVALID
  future_issued_token: AUTHENTICATION_TOKEN_INVALID
  verifier_mismatch: AUTHENTICATION_TOKEN_INVALID
  player_missing_or_inactive: AUTHENTICATION_TOKEN_INVALID
  repository_or_unit_of_work_unavailable: AUTHENTICATION_TOKEN_STORE_UNAVAILABLE
```

`AUTHENTICATION_TOKEN_EXPIRED`、`AUTHENTICATION_TOKEN_REVOKED` 和 `AUTHENTICATION_ACCOUNT_DISABLED` 已存在于 semantic catalog。第一版 validation behavior 不应公开区分它们，除非后续披露决策明确授权这种姿态。

公开面不得暴露 token record 是否存在、token 是否 expired 或 revoked、player 是否存在、stored verifier key id 是什么、audience 是否匹配，或者 verifier mismatch 是否发生。

## 9. 脱敏要求

未来行为不得把以下值放进 errors、logs、docs examples、test failure messages、conversation logs、ADRs、change specs 或 public responses：

- Raw access-token text。
- Raw access-token bytes。
- Lookup digest bytes。
- Verifier digest bytes。
- HMAC input 或 output bytes。
- Verifier key bytes。
- 完整具体的 `verifier_key_id` 值。
- Token lookup hit 或 miss details。
- Failed proof 的 token lifecycle details。
- Failed proof 的 audience mismatch details。
- Player lookup hit 或 miss details。
- Failed proof 的 repository implementation details。

允许：

- 已注册的 public error codes。
- `<access-token>` 和 `<verifier-key-id>` 这类 redacted placeholders。
- 只有 proof 成功后，才能在授权 result fields 或 internal tests 中使用 non-secret record ids。

## 10. 必需测试

未来实现必须在 `runtime/internal/app/authentication/service_test.go` 中添加或更新 focused tests。

必需 test classes：

```yaml
required_tests:
  validation_rejects_missing_token_without_unit_of_work
  validation_rejects_malformed_token_without_unit_of_work
  validation_computes_lookup_digest_before_repository_lookup
  validation_uses_authentication_repository_from_unit_of_work_only
  validation_uses_player_repository_from_unit_of_work_only
  validation_collapses_lookup_miss_to_invalid_token
  validation_rejects_inactive_expired_or_revoked_token
  validation_rejects_wrong_kind_algorithm_version_key_or_audience
  validation_compares_verifier_digest_before_request_identity
  validation_collapses_verifier_mismatch_to_invalid_token
  validation_requires_active_player_account
  validation_returns_request_identity_only_after_unit_of_work_success
  validation_keeps_session_validated_false_without_session_persistence
  validation_errors_do_not_leak_raw_token_or_digest_material
  validation_does_not_touch_protocol_carriers_or_route_protection
```

常规测试必须使用 fake unit-of-work 和 repository implementations。它们不得要求 live PostgreSQL server。

## 11. Nakama And Pitaya 映射

Nakama capability mapping：

- Session/access-token validation 被适配为 application-owned proof validation。
- Token expiration 和 revocation 在 service validation flow 中检查。
- Server-side account state 仍然决定 token 是否可以产生 player identity。

Pitaya capability mapping：

- Realtime handlers 应该在 validation 后接收 identity context。
- Frontend acceptors 和 backend handlers 继续与 proof validation 分离。
- Session binding 是 runtime context，本身不是 proof。

这些参考只指导 capability coverage。它们不授权复制 Nakama 或 Pitaya public APIs。

## 12. 延期项

本 gate 不授权：

- Access-token validation execution。
- WebSocket proof carriers。
- WebSocket handshake authentication。
- Route protection。
- Session persistence。
- Logout execution。
- Refresh execution。
- Token cleanup jobs。
- Token validation audit mutation。
- Protocol authentication messages。
- Startup wiring。
- Repository interface changes。
- PostgreSQL adapter changes。
- SQL migration changes。
- Generated file changes。
- External cryptography、JWT、OAuth、OIDC、provider、KMS、cloud secret-manager、password-hashing、Redis-like、queue 或 session-store dependencies。
- Production authentication behavior。

## 13. Verification

此 gate 的 repository check rule 是：

```text
runtime.access_token_validation_service_behavior_gate
```
