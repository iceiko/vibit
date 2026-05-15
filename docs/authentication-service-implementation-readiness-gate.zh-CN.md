# Authentication Service Implementation Readiness Gate

状态：Draft v0.1
最后更新：2026-05-15
范围：第一段 application authentication service implementation work 之前的 entry criteria、package ownership、file boundaries、test expectations、sequencing、reference mapping 和 deferrals
依赖：`docs/runtime-authentication-implementation-boundary.md`, `docs/application-authentication-service-interface-boundary.md`, `docs/token-credential-verifier-algorithm-redaction-boundary.md`, `docs/secret-configuration-verifier-key-loading-boundary.md`, `docs/token-credential-material-generation-boundary.md`, `docs/verifier-digest-computation-comparison-boundary.md`
Canonical decision: `ADR-0044`

配套英文源文档是 `docs/authentication-service-implementation-readiness-gate.md`。英文文件是权威版本。

## 1. Purpose

本 gate 定义 vibit 在开始第一段 application authentication service implementation code 前必须满足什么条件。

Authentication 是 security-sensitive 且 cross-cutting 的。如果没有 implementation readiness gate，future agent 可能一边写 service code，一边悄悄选择 package ownership、secret loading、token generation、digest helper names、repository call shape、protocol carriers、failure behavior 或 test posture。

这只是 readiness-only standard。它不添加 Go code、imports、runtime services、secret loading、token generation、credential generation、verifier digest computation、verifier comparison、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket carriers、routes、repository methods、SQL migrations、authentication dependencies 或 production authentication behavior。

## 2. Core Rule

第一版 authentication service implementation readiness posture 是：

```yaml
authentication_service_implementation_readiness_gate: defined
implementation_authorized_by_this_standard: false
planned_owner: runtime/internal/app
planned_package_candidate: runtime/internal/app/authentication
first_code_slice_must_be_separately_authorized: true
service_code_status: deferred
secret_loading_code_status: deferred
material_generation_code_status: deferred
digest_helper_code_status: deferred
verifier_comparison_code_status: deferred
login_execution_status: deferred
token_validation_status: deferred
logout_execution_status: deferred
cleanup_execution_status: deferred
protobuf_authentication_messages_status: deferred
websocket_proof_carrier_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
major_dependency_status: deferred
```

这个 gate 只有在 agents 能检查一个标准并知道以下内容时才算完成：

- 必须阅读哪些 prior boundaries。
- 第一段 implementation slice 可以编辑哪些 packages。
- 哪些 behaviors 仍是 separate gates。
- 接受 behavior 前需要哪些 tests。
- 正在适配 Nakama 和 Pitaya 的哪些 capabilities。
- 哪些 public/protocol surfaces 仍不在范围内。

## 3. Required Prior Boundaries

Implementation code 开始前，implementing agent 必须阅读并保持：

```yaml
required_boundaries:
  runtime_authentication_implementation_boundary: docs/runtime-authentication-implementation-boundary.md
  application_authentication_service_interface_boundary: docs/application-authentication-service-interface-boundary.md
  verifier_algorithm_redaction_boundary: docs/token-credential-verifier-algorithm-redaction-boundary.md
  secret_configuration_key_loading_boundary: docs/secret-configuration-verifier-key-loading-boundary.md
  material_generation_boundary: docs/token-credential-material-generation-boundary.md
  verifier_digest_computation_comparison_boundary: docs/verifier-digest-computation-comparison-boundary.md
```

第一段 implementation slice 不得重新解释这些 boundaries。如果 implementation 需要改变其中之一，必须先打开单独的 standard 或 ADR，再改 code。

## 4. Package And File Ownership

第一段 authentication service implementation 应保持 application-owned。

后续 implementation work item 授权后允许的 future write area：

```text
runtime/internal/app/authentication/
```

后续 implementation work item 授权后允许的 application integration points：

```text
runtime/internal/app/
runtime/internal/app/bootstrap/
```

允许的 test area：

```text
runtime/internal/app/authentication/*_test.go
runtime/internal/app/*_test.go
```

除非后续 work item 明确命名，否则第一段 slice 禁止写入：

- `runtime/internal/platform/transport/ws/`
- `runtime/internal/platform/protocol/protobuf/`
- `runtime/internal/generated/`
- `runtime/internal/modules/authentication/`
- `runtime/internal/platform/persistence/postgres/`
- `runtime/migrations/postgres/`
- `proto/`
- `contracts/runtime/authentication/`

第一段 implementation slice 可以使用 existing generated metadata 和 existing repository interfaces，但不得顺手 hand-edit generated files 或改变 public contracts。

## 5. First Implementation Queue

Readiness gate 推荐这个 implementation order：

```yaml
recommended_queue:
  - define_local_verifier_key_configuration_code_gate
  - implement_application_secret_configuration_loader
  - implement_token_credential_material_generation_helpers
  - implement_verifier_digest_helpers_and_comparison
  - implement_authentication_service_in_memory_unit_tests
  - implement_device_credential_authentication_service_flow
  - implement_access_token_validation_service_flow
  - define_protocol_authentication_message_gate
  - define_websocket_request_proof_carrier_gate
```

这个 queue 是 guidance，不是 implicit authorization。每一项仍需要 bounded work item、change spec、tests、verification 和 documentation update。

## 6. Minimum First Code Slice

最小可接受的 first code slice 应证明一个狭窄 behavior，且不打开 protocol 或 transport surfaces。

Recommended first code slice：

```yaml
first_code_slice:
  name: local_verifier_key_configuration_loader
  owner: runtime/internal/app/authentication
  behavior: load_explicit_in_memory_or_environment_supplied_key_set
  production_behavior: fail_closed_for_invalid_config
  external_dependencies: none
  protocol_changes: none
  repository_changes: none
  migration_changes: none
```

为什么它排在第一：

- 没有 key material，token 和 credential generation 无法验证。
- Digest helpers 应接收 validated logical keys，而不是自己 parse configuration。
- Service code 不应在 login logic 里发明 secret-loading behavior。

本标准不授权这个 code slice。它只是命名 recommended next implementation gate。

## 7. Service Behavior Entry Criteria

Login 或 token validation behavior 被实现前，必须满足：

```yaml
service_behavior_entry_criteria:
  explicit_key_configuration_loader_exists: required
  material_generation_helpers_exist: required_for_token_or_credential_issuance
  digest_helpers_exist: required
  constant_time_verifier_comparison_exists: required
  redaction_tests_exist: required
  repository_usage_through_unit_of_work: required
  generated_authentication_shapes_remain_metadata_only: required
  protobuf_authentication_messages_defined: required_before_wire_exposure
  websocket_proof_carrier_defined: required_before_realtime_exposure
```

Login execution 不得通过 WebSocket 或 Protobuf 暴露，直到 protocol carrier gates 存在。

Token validation 不得成为 production-sensitive domain authorization，直到 explicit request proof carrier behavior 存在且 request identity handoff 已测试。

## 8. Required Tests For First Behavior

Future authentication implementation work 必须按其引入的 behavior 添加 tests。

Minimum test classes：

```yaml
required_test_classes:
  configuration_tests:
    - missing_key_fails_closed
    - malformed_key_fails_closed
    - duplicate_logical_keys_fail_closed
    - secret_values_absent_from_errors
  generation_tests:
    - generated_material_has_32_raw_bytes
    - generated_material_round_trips_through_text_encoding
    - raw_material_not_stored
  digest_tests:
    - canonical_input_is_stable
    - purpose_labels_are_separated
    - digest_output_is_32_bytes
    - lookup_digest_is_not_authentication_proof
  comparison_tests:
    - verifier_comparison_uses_constant_time_primitive
    - mismatch_and_missing_record_share_public_failure
  service_tests:
    - repository_is_used_through_unit_of_work
    - request_identity_is_populated_only_after_valid_proof
    - public_errors_are_registered_and_redacted
```

Tests 不得要求运行中的 PostgreSQL server，除非明确标记为 live 且通过 existing PostgreSQL verification standard 选择性启用。

## 9. Redaction And Observability

Future service code 必须通过结构保持 redaction。

禁止出现在 logs、public errors、traces、metrics labels、audit-safe facts、tests snapshots、ADRs、change specs、documentation examples 和 conversation logs 中：

- Raw credentials。
- Raw access tokens。
- Encoded generated material。
- Lookup digests。
- Verifier digests。
- Verifier key values。
- Encoded verifier key values。
- Full concrete `verifier_key_id` values。
- HMAC input bytes。
- HMAC output bytes。
- Environment variable values。

允许：

- Registered error codes。
- Non-secret record identifiers。
- 明确 safe 的 lifecycle state names。
- 用于说明 redaction rules 的 `<raw-access-token>` 和 `<verifier-key-id>` placeholders。

Observability 必须有用，但不能暴露 proof material。后续 observability standard 可以定义 redacted fingerprints，但这里没有 ratify 任何 fingerprint。

## 10. Nakama And Pitaya Mapping

需要适配的 Nakama capabilities：

- Device/custom authentication entry point。
- Session token issuance。
- Session token validation。
- Token expiration and revocation。
- User/account status checks。

需要适配的 Pitaya capabilities：

- Frontend transport acceptor 与 backend handler logic 分离。
- Route handler context 接收 validated identity。
- Session binding 是 application/runtime context，不是 proof 本身。

vibit implementation rule：

- Transport accepts frames。
- Protocol decodes messages。
- Application authentication validates proof。
- Application request identity carries validated actor context。
- Domain modules consume request identity and do not parse proof。

不要直接复制 Nakama 或 Pitaya public APIs。只把它们作为 capability coverage 和 vocabulary references。

## 11. Deferrals

这个 readiness gate 不授权：

- Authentication service code。
- Secret loading code。
- Token generation。
- Credential generation。
- Verifier digest computation。
- Verifier comparison。
- Login execution。
- Access-token validation。
- Logout execution。
- Refresh behavior。
- Cleanup jobs。
- Protobuf authentication messages。
- WebSocket proof carriers。
- WebSocket handshake authentication。
- Authentication dependencies。
- `authentication.Repository` changes。
- PostgreSQL migration schema changes。
- Production authentication behavior。

## 12. Verification Path

该 gate 的 repository check rule 是：

```text
runtime.authentication_service_implementation_readiness_gate
```

触碰该 gate 的 change 应运行：

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

本标准不添加 runtime Go authentication behavior，因此不验证该 behavior。

## 13. Completion Criteria

这个 gate 完成条件：

- Readiness standard 和 translation 存在。
- ADR-0044 记录 readiness decision。
- Manifests 和 agent guides 引用该 gate。
- Repository checks enforce readiness markers。
- `W-0095` completed。
- Next ready work item 是 bounded implementation 或 preparation gate。
