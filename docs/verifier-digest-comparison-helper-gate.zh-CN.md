# Verifier Digest Comparison Helper Gate 中文版

状态：Draft v0.1
最后更新：2026-05-15
范围：在添加 comparison code 之前，定义 verifier digest comparison helper implementation gate、future helper file boundaries、constant-time primitive posture、input redaction、failure posture、tests 和 deferrals
依赖：`docs/verifier-digest-computation-comparison-boundary.md`、`docs/verifier-digest-helper-implementation-gate.md`、`docs/token-credential-verifier-algorithm-redaction-boundary.md`、`docs/authentication-service-implementation-readiness-gate.md`
规范化决策：`ADR-0049`

本文件是 `docs/verifier-digest-comparison-helper-gate.md` 的简体中文译本。英文文件是权威版本。

## 1. 目的

本 gate 定义 verifier digest comparison helpers 的下一个 bounded implementation slice。

仓库已经在 `runtime/internal/app/authentication` 下拥有 helper-only verifier digest computation。下一步风险是 future agent 使用非 constant-time primitive 比较 verifier material，比较错误的 material，把 authentication service behavior 藏进 helper，或者把 comparison 移到 repositories、protocol adapters、transport handlers 或 generated output。

这是 implementation-gate standard。它不添加 Go code、imports、verifier comparison、authentication service behavior、login execution、token validation、logout execution、refresh behavior、cleanup jobs、Protobuf messages、WebSocket proof carriers、repository methods、SQL migrations、startup wiring、authentication dependencies、external cryptography services、KMS、cloud secret-manager integrations 或 production authentication behavior。

## 2. 核心规则

Verifier digest comparison helper gate 是：

```yaml
verifier_digest_comparison_helper_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0105
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
future_slice: constant_time_verifier_digest_comparison_helpers
future_source: runtime/internal/app/authentication/verifier_comparison.go
future_tests: runtime/internal/app/authentication/verifier_comparison_test.go
verifier_algorithm_family: vibit_hmac_sha256_v1
verifier_version: 1
digest_size_bytes: 32
input_shape: computed_verifier_digest_and_stored_verifier_digest_bytes
computed_digest_handoff: ComputedDigest
stored_digest_handoff: copied_repository_digest_bytes
preferred_constant_time_primitive: crypto/hmac.Equal
acceptable_constant_time_primitive: crypto/subtle.ConstantTimeCompare
comparison_result_shape: redacted_match_or_mismatch_result
raw_material_comparison: forbidden
lookup_digest_comparison_for_authentication: forbidden
database_only_verifier_comparison: forbidden
service_behavior_status: deferred
authentication_service_behavior_status: deferred
protobuf_authentication_messages_status: deferred
websocket_proof_carrier_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
startup_wiring_status: deferred
authentication_dependencies_status: deferred
production_authentication_behavior_status: deferred
```

Future implementation 必须是 helper-only slice。它可以用 constant-time primitive 比较 computed verifier digest 和 stored verifier digest bytes，并返回 redacted comparison result。它不得 compute digests、load keys、select records、parse proofs、issue login responses、validate tokens、revoke tokens、refresh tokens、call repositories、inspect lifecycle state 或 touch protocol carriers。

## 3. Future Helper Shape

Future implementation ownership：

```text
runtime/internal/app/authentication
```

后续 implementation work item 授权 code 后允许的 future files：

```text
runtime/internal/app/authentication/verifier_comparison.go
runtime/internal/app/authentication/verifier_comparison_test.go
```

`verifier_digest.go` 保持 computation-only。Comparison helpers 必须放在单独的 comparison file 中，这样后续 agent 可以分别理解 computation 和 comparison。

推荐的 future API shape：

```yaml
future_types:
  VerifierComparisonResult:
    owns: digest_class_and_match_status_only
    methods:
      - Class() DigestClass
      - Matched() bool
    constraints:
      - Result text 和 string formatting 不得暴露 digest bytes、raw material、key ids、account ids、token ids 或 lookup miss details。
      - False match 本身不是 public authentication error；service behavior 以后再映射。

future_functions:
  CompareCredentialVerifierDigest:
    signature: "func CompareCredentialVerifierDigest(computed ComputedDigest, storedDigest []byte) (VerifierComparisonResult, error)"
    expected_computed_class: credential_verifier
    behavior: constant_time_compare_computed_verifier_digest_to_stored_credential_verifier_digest
  CompareTokenVerifierDigest:
    signature: "func CompareTokenVerifierDigest(computed ComputedDigest, storedDigest []byte) (VerifierComparisonResult, error)"
    expected_computed_class: token_verifier
    behavior: constant_time_compare_computed_verifier_digest_to_stored_token_verifier_digest
```

允许使用 unexported shared helper，只要 exported functions 仍然是 class-specific，并且能避免意外比较 lookup digest。

## 4. 输入边界

允许的输入：

- 由 W-0103 digest helpers 产生的 `ComputedDigest`。
- 从 repository records 返回的 stored verifier digest bytes。

禁止的输入：

- Raw device credential material。
- Raw access-token material。
- Encoded credential 或 token text。
- Lookup digest bytes。
- `credential_lookup_digest`。
- `token_lookup_digest`。
- `verifier_key_id`。
- Player account id。
- Token record id。
- Credential record id。
- Provider subject。
- Session id。
- WebSocket connection metadata。
- Route name 或 protocol metadata。

Comparison helper 只能比较 verifier digest bytes。Lookup digest equality 仍然只能用于 record selection，不能作为 authentication proof。

## 5. Constant-Time Comparison Posture

首选 primitive：

```yaml
preferred_go_comparison: crypto/hmac.Equal
```

可接受 primitive：

```yaml
acceptable_go_comparison: crypto/subtle.ConstantTimeCompare
```

Verifier digest comparison 禁止使用：

- `bytes.Equal`
- `reflect.DeepEqual`
- 对由 digest bytes 派生出的 strings 或 arrays 使用 `==`
- 为 equality 把 byte slice 转换成 string
- map lookup equality
- database-only equality
- 把 SQL comparison of verifier digest bytes 作为最终 proof
- 比较 encoded digest text
- 比较 raw credential 或 token material

Length validation 必须 fail closed。第一版姿态要求 32-byte verifier digests。Helper 可以对 missing 或 malformed digest input 返回 redacted error，但不得让 future public response 暴露 authentication attempt 的哪一侧失败。

## 6. Failure Posture

Helper 必须为测试保留小的 internal distinction，同时支持后续 public failure collapse：

```yaml
internal_failure_classes:
  verifier_digest_mismatch: redacted
  missing_computed_digest: redacted
  wrong_computed_digest_class: redacted
  missing_stored_digest: redacted
  malformed_stored_digest: redacted
  invalid_comparison_input: redacted
future_public_failure_class: invalid_authentication_proof
```

规则：

- Mismatch 不得暴露 lookup 是否成功。
- Missing stored verifier digest 不得通过 public service behavior 暴露 record existence 或 corruption。
- Wrong digest class 必须 fail closed。
- Malformed stored digest 必须 fail closed。
- Error text 不得包含 digest bytes、raw material、key ids、account ids、credential ids、token ids、lookup values 或 candidate counts。
- Helper 不得判断 credential 是否 active、expired、revoked、可用于某个 player，或是否 bound to a session。

## 7. Package And File Boundary

允许的 future implementation area：

```text
runtime/internal/app/authentication/verifier_comparison.go
runtime/internal/app/authentication/verifier_comparison_test.go
```

Future verifier comparison helper slice 禁止写入：

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

Future implementation 不得把 comparison wire into login execution、account creation、token issuance、token validation、startup、WebSocket transport、Protobuf protocol、PostgreSQL persistence、generated contract output、migrations 或 domain repository code。

## 8. Error And Redaction Requirements

允许的 error classification examples：

```yaml
error_classes:
  verifier_digest_mismatch
  missing_computed_digest
  wrong_computed_digest_class
  missing_stored_digest
  malformed_stored_digest
  invalid_verifier_comparison_input
```

Error text 允许包含：

- Error classes。
- Digest class names。
- 非 secret 数字期望值，例如 `32` bytes。

Errors、logs、test snapshots、docs、ADRs、change specs 和 conversation logs 中禁止出现：

- Raw credential material 或 bytes。
- Raw access-token material 或 bytes。
- Encoded generated material。
- Lookup digest bytes。
- Verifier digest bytes。
- HMAC input bytes。
- HMAC output bytes。
- Stored verifier digest bytes。
- Verifier key values。
- Encoded verifier key values。
- Full concrete `verifier_key_id` values。
- 来自真实数据的 account ids、credential record ids 或 token record ids。
- Candidate key-set counts。
- Lookup hit 或 miss details。

## 9. Future Helper Required Tests

Future implementation 必须在以下路径添加 focused unit tests：

```text
runtime/internal/app/authentication/verifier_comparison_test.go
```

最低测试集合：

```yaml
required_tests:
  credential_verifier_digest_match_returns_matched
  credential_verifier_digest_mismatch_returns_not_matched
  token_verifier_digest_match_returns_matched
  token_verifier_digest_mismatch_returns_not_matched
  comparison_uses_crypto_hmac_equal
  comparison_rejects_lookup_digest_classes
  comparison_rejects_wrong_verifier_digest_class
  comparison_rejects_missing_computed_digest
  comparison_rejects_missing_stored_digest
  comparison_rejects_malformed_stored_digest_length
  comparison_does_not_compare_raw_material
  comparison_does_not_call_repositories
  comparison_result_does_not_expose_digest_bytes
  comparison_errors_are_redacted
  comparison_helpers_do_not_implement_authentication_service_behavior
```

Tests 只能在测试内部使用 deterministic synthetic digest bytes。不得记录 production-like secrets、真实 credential material、真实 access tokens、真实 key ids、account ids、credential record ids、token record ids 或 lookup digests。

## 10. Dependency Posture

本 gate 不允许新增 external dependency。

Future implementation work item 授权 code 后允许使用的 Go standard library packages：

```yaml
future_standard_library_imports_allowed:
  - crypto/hmac
  - crypto/subtle
  - errors
  - fmt
```

第一版 helper implementation 不得添加 JWT、JWK、OAuth、OIDC、provider SDKs、password-hashing dependencies、Redis-like stores、KMS SDKs、cloud secret-manager SDKs、operations libraries 或 external cryptography services。

## 11. Nakama And Pitaya Mapping

Nakama capability reference：

- Server-side credential and token validation 在 lookup 和 digest computation 之后需要最终 verifier comparison step。vibit 采用 capability need，不采用 Nakama 的 public API shape。

Pitaya capability reference：

- Realtime route handlers 应该在 validation 之后接收 identity context。vibit 把 verifier digest comparison 保持在 application-owned helpers 中，而不是放到 frontend acceptors、route handlers 或 session binding code 中。

本 gate 把这些 reference 映射为一个窄 helper slice：只做 comparison，service behavior 后续再做。

## 12. Non-Goals

本 gate 不会：

- 添加 verifier comparison code。
- 添加 digest computation code。
- 添加 token generation orchestration。
- 添加 credential generation orchestration。
- 实现 authentication service behavior。
- 执行 login。
- Validate access tokens。
- 执行 logout。
- 添加 refresh behavior。
- 添加 cleanup jobs。
- 添加 Protobuf authentication messages。
- 添加 WebSocket proof carriers。
- 改变 WebSocket handshake behavior。
- 改变 `authentication.Repository`。
- 改变 PostgreSQL migration schemas。
- 添加 startup wiring。
- 添加 authentication dependencies。
- 添加 production authentication behavior。

## 13. Verification Path

本 gate 的 repository check rule 是：

```text
runtime.verifier_digest_comparison_helper_gate
```

修改本 gate 时运行：

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

本 gate 不验证 runtime Go verifier digest comparison behavior，因为它不添加 comparison behavior。

## 14. Follow-Up Gates

推荐 follow-up gates：

- Implement verifier digest comparison helpers。
- Implement application authentication service behavior。
- Add Protobuf authentication messages。
- Add WebSocket request proof carriers。
