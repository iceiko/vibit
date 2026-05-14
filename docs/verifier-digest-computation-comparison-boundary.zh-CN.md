# Verifier Digest Computation And Comparison Boundary

状态：Draft v0.1
最后更新：2026-05-15
范围：第一版 device-credential 与 opaque access-token 姿态下 future verifier digest computation ownership、canonical byte input construction、purpose-label use、logical key use、key-id selection、lookup digest handoff、verifier digest comparison、failure redaction、dependency posture 和 test expectations
依赖：`docs/token-credential-verifier-algorithm-redaction-boundary.md`, `docs/secret-configuration-verifier-key-loading-boundary.md`, `docs/token-credential-material-generation-boundary.md`
Canonical decision: `ADR-0043`

配套英文源文档是 `docs/verifier-digest-computation-comparison-boundary.md`。英文文件是权威版本。

## 1. Purpose

本标准定义 future application-owned authentication code 如何计算 lookup digests、计算 verifier digests，并比较 verifier digests。

它位于 verifier digest helper code、verifier comparison code、token generation code、credential generation code、secret loading、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、authentication dependencies、repository interface changes、migration schema changes 或 production authentication behavior 之前。

这只是 boundary-only standard。它不添加 Go code、imports、runtime services、HMAC helpers、digest helpers、comparison helpers、token generation、credential generation、secret loading、repository methods、SQL migrations、Protobuf messages、WebSocket carriers、routes 或 production authentication behavior。

## 2. Core Rule

第一版 verifier digest computation and comparison posture 是：

```yaml
verifier_digest_computation_comparison_boundary: defined
implementation_authorized_by_this_standard: false
planned_owner: runtime/internal/app
planned_optional_child_package: runtime/internal/app/authentication
verifier_algorithm_family: vibit_hmac_sha256_v1
verifier_version: 1
digest_size_bytes: 32
digest_output_shape: raw_32_byte_digest
canonical_input_version: vibit.auth.verifier.input.v1
canonical_input_encoding: ascii_header_length_prefixed_purpose_label_length_prefixed_raw_material
lookup_digest_database_equality_allowed_for_record_selection: true
lookup_digest_database_equality_sufficient_for_authentication: false
constant_time_verifier_comparison_required: true
missing_lookup_public_failure: invalid_authentication_proof
invalid_verifier_public_failure: invalid_authentication_proof
unknown_key_id_public_failure: invalid_authentication_proof
external_cryptography_dependency_required_for_first_posture: false
future_allowed_standard_library_packages:
  - crypto/hmac
  - crypto/sha256
  - crypto/subtle
```

后续 code gate 授权 implementation 后，future first-posture digest helpers 可以使用 Go standard library `crypto/hmac`、`crypto/sha256` 和 `crypto/subtle`。该 digest computation and comparison posture 不需要 external cryptography、password-hashing、JWT、OAuth、OIDC、provider、KMS、cloud secret-manager 或 operations dependency。

## 3. Ownership

Future verifier digest computation and comparison 由 application 拥有：

```text
runtime/internal/app
```

Future helper code 可以放在 application-owned child package，例如：

```text
runtime/internal/app/authentication
```

Ownership rules：

- Application authentication code 在后续 implementation gate 授权后拥有 digest computation 和 verifier comparison。
- `authentication.Repository` 只能存取 already-computed digest material。
- PostgreSQL adapters 可使用 lookup digest equality 进行 record selection，但它们不 compute digests、compare verifier digests、choose keys 或决定 authentication outcomes。
- Protobuf adapters 和 WebSocket transports 可在后续 protocol gates 后承载 already-decoded proof fields，但它们不 compute digests 或 compare verifiers。
- Generated authentication contract shapes 保持 metadata-only 且 immutable。

Digest computation 和 comparison 不得放入 transport handlers、protocol adapters、domain modules、repositories、generated output、migrations、SQL fixtures，或未明确限定为 future application-owned digest helpers 的 tests。

## 4. Canonical Digest Input

Future digest helpers 必须先构建 deterministic canonical byte input，再进行 HMAC computation。

第一版 planned canonical input 是：

```yaml
canonical_digest_input:
  version_header_ascii: vibit.auth.verifier.input.v1
  header_separator: 0x00
  purpose_label_length: uint16_big_endian_byte_length
  purpose_label: ascii_bytes
  raw_material_length: uint16_big_endian_byte_length
  raw_material: generated_secret_material_bytes
```

Byte sequence 是：

```text
ascii("vibit.auth.verifier.input.v1")
|| 0x00
|| uint16be(len(purpose_label))
|| ascii(purpose_label)
|| uint16be(len(raw_material))
|| raw_material
```

规则：

- Purpose labels 必须是 ASCII，并且必须匹配 digest class 注册的 labels。
- Raw material 必须是 decoded raw secret material bytes，不是 normalized text、metadata、player id、session id、route name 或 provider subject。
- 第一版姿态下 raw material length 是 32 bytes。
- Length prefixes 让 input 即使在 future raw material shape 改变时也保持 unambiguous。
- Version header 是 HMAC input 的一部分；如果 canonical input shape 变化，必须改变 version header。
- Future tests 必须包含 canonical byte construction 的 deterministic fixture vectors。

本标准定义 future helpers 的 input shape。它不实现这些 helpers。

## 5. Digest Classes

第一版姿态有四个 digest classes：

```yaml
credential_lookup_digest:
  purpose_label: vibit.credential.lookup.v1
  logical_key: credential_lookup_key
  output_bytes: 32
  storage_column: authentication_device_credentials.credential_lookup_digest
  database_equality_for_selection: allowed
  log_safe: false
credential_verifier_digest:
  purpose_label: vibit.credential.verifier.v1
  logical_key: credential_verifier_key
  output_bytes: 32
  storage_column: authentication_device_credentials.credential_verifier_digest
  constant_time_comparison_required: true
  log_safe: false
token_lookup_digest:
  purpose_label: vibit.access_token.lookup.v1
  logical_key: token_lookup_key
  output_bytes: 32
  storage_column: authentication_access_tokens.token_lookup_digest
  database_equality_for_selection: allowed
  log_safe: false
token_verifier_digest:
  purpose_label: vibit.access_token.verifier.v1
  logical_key: token_verifier_key
  output_bytes: 32
  storage_column: authentication_access_tokens.token_verifier_digest
  constant_time_comparison_required: true
  log_safe: false
```

规则：

- Lookup digests 和 verifier digests 必须使用不同 purpose labels。
- Credential digests 和 token digests 必须使用不同 purpose labels。
- 每个 digest class 必须使用 selected verifier key set 中匹配的 logical key。
- Digest bytes 不得为 storage、lookup、comparison、logs、public errors、metrics labels、conversation logs、ADRs 或 change specs 而被截断。
- 任何 future digest prefix 或 fingerprint format 都需要单独 redaction/fingerprint standard。

## 6. Key Selection

New credential and token writes 必须使用 future secret configuration code 选出的 active verifier key set。

Future presented credential 或 access-token proof validation 必须考虑 accepted key sets，因为 opaque proof text 不携带 `verifier_key_id`。

Planned validation posture：

```yaml
new_write_key_selection:
  key_set: active_key_set
  stored_record_keeps_verifier_key_id: true
validation_lookup_key_selection:
  candidate_key_sets: active_and_accepted_previous_key_sets
  compute_lookup_digest_per_candidate_key_set: required
  repository_lookup: by_lookup_digest
validation_verifier_key_selection:
  selected_by_stored_record_verifier_key_id: required
  unknown_key_id_public_failure: invalid_authentication_proof
  retired_key_id_public_failure: invalid_authentication_proof
```

规则：

- Future validation 可以为 accepted key sets 计算多个 lookup digests。
- Repository lookup by lookup digest 只选择 candidate record；它不是 authentication proof。
- 选中 record 后，future code 必须使用 stored `verifier_key_id` 标识的 key set 计算 matching verifier digest。
- 如果 stored key id unknown、retired、malformed 或 unavailable，public behavior 必须使用同一个 invalid proof failure class。
- Public error 不得暴露 active key set id、previous key set ids、candidate count、lookup miss、key miss、rotation state 或 digest value。

本标准不改变 repository interface。如果后续 implementation 需要 batch lookup by multiple digest candidates，那个 work item 必须显式更新 repository boundary。

## 7. Comparison Boundary

Future verifier comparison 必须使用 constant-time equality 比较 computed verifier digest bytes 和 stored verifier digest bytes。

可接受的 future Go comparison primitives：

```yaml
preferred_go_comparison: crypto/hmac.Equal
acceptable_go_comparison: crypto/subtle.ConstantTimeCompare
forbidden_comparison:
  - bytes.Equal
  - string equality
  - byte_slice_string_conversion
  - reflect.DeepEqual
  - database_only_equality
  - map_lookup_equality
```

规则：

- 比较 verifier digest bytes，而不是 raw credential text、raw access-token text、encoded material、lookup digest bytes、key ids 或 public identifiers。
- Database equality on lookup digest 只能用于 candidate record selection。
- Lookup hit 仍必须通过 lifecycle checks 和 constant-time verifier digest comparison，才能接受 proof。
- Lookup miss、missing stored verifier digest、mismatched verifier digest、unknown key id、unsupported algorithm version、expired token、revoked token、disabled credential 或 disabled account 不得暴露更详细的 public failure reason，除非后续 semantic error standard 明确允许。
- Future code 应避免 missing records 和 invalid verifier digests 之间明显的 timing differences；exact equalization strategy 属于后续 implementation gate。

## 8. Failure Redaction

Future credential 与 token verifier problems 的 public failures 必须收敛到 registered invalid-proof class，除非后续 semantic standard 明确允许更具体的 public class。

Required public failure posture：

```yaml
public_failure_class:
  missing_lookup_record: invalid_authentication_proof
  verifier_digest_mismatch: invalid_authentication_proof
  unknown_verifier_key_id: invalid_authentication_proof
  unsupported_verifier_algorithm: invalid_authentication_proof
  malformed_presented_proof: invalid_authentication_proof
  expired_or_revoked_token: invalid_authentication_proof
```

禁止出现在 public errors、logs intended for client support、audit-safe facts、traces、metrics labels、ADRs、change specs、documentation examples 和 conversation logs 中：

- Raw credential material。
- Raw access-token material。
- Encoded credential 或 token material。
- Lookup digest bytes。
- Verifier digest bytes。
- Verifier key values。
- Encoded verifier key values。
- Full concrete `verifier_key_id` values。
- Candidate key-set counts。
- HMAC input bytes。
- HMAC output bytes。
- 能区分 lookup miss 与 verifier mismatch 的 reason。

Allowed with care：

- Registered error codes。
- 周围 artifact 已允许时的 non-secret lifecycle state names。
- 用于说明 redaction rules 的 `<lookup-digest>` 或 `<verifier-key-id>` placeholders。

## 9. Test Expectations

未来只要 implementation work 计算 lookup digests、计算 verifier digests 或比较 verifier digests，就必须添加 focused tests。

Minimum expectations：

```yaml
digest_tests:
  canonical_input_is_deterministic: required
  canonical_input_uses_version_header: required
  canonical_input_length_prefixes_purpose_label: required
  canonical_input_length_prefixes_raw_material: required
  lookup_and_verifier_purpose_labels_differ: required
  credential_and_token_purpose_labels_differ: required
  digest_output_is_32_bytes: required
  lookup_digest_uses_lookup_key: required
  verifier_digest_uses_verifier_key: required
  lookup_digest_not_used_as_authentication_proof: required
  verifier_comparison_uses_constant_time_primitive: required
  missing_record_and_mismatch_share_public_failure: required
  raw_material_absent_from_outputs: required
  digest_material_absent_from_outputs: required
```

Test fixtures 只能在 tests 内使用 deterministic synthetic keys 和 raw material。它们不得成为 production defaults、real secret shape 的 documentation examples，或 committed production-like secret values。

## 10. Dependency Posture

第一版 digest computation and comparison posture 不需要 external dependency。

后续 explicit code gate 授权后可接受：

```yaml
go_standard_library:
  hmac: crypto/hmac
  hash: crypto/sha256
  constant_time: crypto/subtle
external_dependency_adoption_record_required_for_first_posture: false
```

第一版 posture deferred 且不需要：

- External cryptography libraries。
- Password hashing dependencies。
- JWT、JWK 或 signing libraries。
- OAuth 或 OIDC provider SDKs。
- KMS 或 cloud secret-manager SDKs。
- Redis-like token/session stores。

如果未来 login method 接受 password-like 或 low-entropy material，必须定义单独的 credential boundary 和 dependency adoption record。

## 11. Non-Goals

本标准不：

- 添加 verifier digest computation。
- 添加 verifier comparison。
- 添加 token generation。
- 添加 credential generation。
- 添加 secret loading。
- 添加 application authentication service code。
- 添加 login execution。
- 添加 access-token validation。
- 添加 logout execution。
- 添加 refresh behavior。
- 添加 cleanup jobs。
- 添加 Protobuf authentication messages。
- 添加 WebSocket proof carriers。
- 改变 WebSocket handshake。
- 添加 authentication dependencies。
- 改变 `authentication.Repository`。
- 改变 PostgreSQL migration schemas。
- 添加 production authentication behavior。

## 12. Verification Path

该 boundary 的 repository check rule 是：

```text
runtime.verifier_digest_computation_comparison_boundary
```

触碰该 boundary 的 change 应运行：

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

本标准不添加 runtime Go verifier digest computation 或 comparison behavior，因此不验证该 behavior。

## 13. Follow-Up Gates

Recommended follow-up gates：

- Authentication service implementation readiness gate。
- Local verifier key configuration implementation gate。
- Token and credential material generation implementation gate。
- Verifier digest helper implementation gate。
- Authentication redaction test implementation gate。
- Protobuf authentication message gate。
- WebSocket request proof carrier gate。
