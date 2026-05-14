# Token And Credential Verifier Algorithm Redaction Boundary

Status: Draft v0.1
Last updated: 2026-05-15
Scope: first device-credential 和 opaque access-token posture 的 verifier algorithm posture、digest classification、key identifier treatment、constant-time comparison expectations、dependency posture 和 redaction test expectations
Depends on: `docs/credential-record-schema-boundary.md`, `docs/token-verifier-record-schema-boundary.md`, `docs/application-authentication-service-interface-boundary.md`
Canonical decision: `ADR-0040`

配套英文源文档是 `docs/token-credential-verifier-algorithm-redaction-boundary.md`。英文文件是权威版本。

## 1. Purpose

本标准定义未来 application-owned authentication service code 必须遵守的 verifier algorithm 和 redaction boundary。

它位于 token material generation、verifier comparison、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、authentication dependencies、repository interface changes、migration schema changes 或 production authentication behavior 之前。

这只是 boundary-only standard。它不添加 Go code、imports、runtime services、token generation、digest comparison、key loading、configuration、repository methods、SQL migrations、Protobuf messages、WebSocket carriers、routes 或 production authentication behavior。

## 2. Core Rule

第一版 verifier posture 是：

```yaml
verifier_algorithm_family: vibit_hmac_sha256_v1
minimum_raw_token_entropy_bits: 256
minimum_raw_device_credential_entropy_bits: 256
token_text_encoding: url_safe_unpadded_base64_or_equivalent
credential_kind: high_entropy_installation_credential
credential_lookup_digest_required: true
credential_verifier_digest_required: true
token_lookup_digest_required: true
token_verifier_digest_required: true
raw_credential_storage: forbidden
raw_token_storage: forbidden
constant_time_verifier_comparison_required: true
password_hashing_dependency_required: false
external_cryptography_dependency_required: false
jwt_or_claim_parsing_required: false
oauth_oidc_dependency_required: false
kms_dependency_required_for_first_posture: false
implementation_authorized_by_this_standard: false
```

在后续 implementation work item 授权 code 之后，第一版 planned verifier algorithm 可以使用 Go standard library 实现：

```yaml
future_allowed_standard_library_packages:
  - crypto/hmac
  - crypto/sha256
  - crypto/subtle
  - crypto/rand
  - encoding/base64
external_dependency_adoption_record_required_for_first_posture: false
```

这里提到 `crypto/rand` 和 `encoding/base64`，只是因为未来 token 或 credential material generation 需要 cryptographically secure randomness 和安全文本编码。本标准不授权 generation code。

## 3. Algorithm Identifiers

第一版 planned algorithm identifier 是：

```yaml
verifier_algorithm: vibit_hmac_sha256_v1
verifier_version: 1
digest_size_bytes: 32
digest_storage_format: raw_32_byte_digest
test_fixture_text_encoding: base64url_unpadded_when_text_is_required
```

该 algorithm identifier 覆盖 lookup digests 和 verifier digests。Purpose labels 和独立 server-side keys 用于区分每一种 digest class。

Required purpose labels:

```yaml
purpose_labels:
  credential_lookup_digest: vibit.credential.lookup.v1
  credential_verifier_digest: vibit.credential.verifier.v1
  token_lookup_digest: vibit.access_token.lookup.v1
  token_verifier_digest: vibit.access_token.verifier.v1
```

规则：

- Purpose labels 是 digest input domain separation 的一部分。
- Lookup 和 verifier digests 不得复用同一个 purpose label。
- Credential 和 token digests 不得复用同一个 purpose label。
- Future algorithm identifiers 必须先 versioned，才能被 stored records 使用。
- Stored records 必须保留 `verifier_algorithm`、`verifier_version` 和 `verifier_key_id`。

## 4. Digest Construction Posture

Planned digest construction 是对 raw proof material 使用 server-side secret keys 和固定 purpose labels 计算 HMAC-SHA-256。

Future code 必须对 canonical byte input 计算 digest：

```yaml
digest_input:
  purpose_label: required
  separator: required_non_ambiguous_byte_separator
  raw_secret_material: credential_or_access_token_bytes
```

Planned digest classes:

```yaml
credential_lookup_digest:
  classification: secret_adjacent_index_material
  construction: HMAC-SHA-256(credential_lookup_key, purpose_label || separator || raw_credential)
  storage: authentication_device_credentials.credential_lookup_digest
  index_use: allowed
  log_safe: false
credential_verifier_digest:
  classification: secret_verifier_material
  construction: HMAC-SHA-256(credential_verifier_key, purpose_label || separator || raw_credential)
  storage: authentication_device_credentials.credential_verifier_digest
  constant_time_comparison_required: true
  log_safe: false
token_lookup_digest:
  classification: secret_adjacent_index_material
  construction: HMAC-SHA-256(token_lookup_key, purpose_label || separator || raw_access_token)
  storage: authentication_access_tokens.token_lookup_digest
  index_use: allowed
  log_safe: false
token_verifier_digest:
  classification: secret_verifier_material
  construction: HMAC-SHA-256(token_verifier_key, purpose_label || separator || raw_access_token)
  storage: authentication_access_tokens.token_verifier_digest
  constant_time_comparison_required: true
  log_safe: false
```

确切 byte separator 和 helper function names 留给后续 code implementation gate。Implementation 必须让 byte construction deterministic、unambiguous、有测试覆盖，并且独立于 transport 或 Protobuf message shape。

通过 database index 进行 lookup digest equality 不能单独作为 proof。Future validation 仍必须执行 status checks、相关 expiration 或 revocation checks，以及 constant-time verifier digest comparison。

## 5. Entropy And Encoding

Raw access-token material 必须至少包含 256 bits entropy。

第一版 `device_credential_login` posture 的 raw device credential material 也必须至少包含 256 bits entropy。它必须是为 authentication 生成的 installation credential，而不是 raw operating-system device identifier、advertising identifier、hardware serial number、account email、player name 或 provider subject。

Token 和 credential text presentation rules:

```yaml
minimum_entropy_bits: 256
first_text_encoding: url_safe_unpadded_base64_or_equivalent
case_sensitive: true
allowed_in_url_query: false
allowed_in_route_name: false
allowed_in_session_metadata: false
allowed_in_logs: false
```

Future text encoding 必须避免 whitespace、control characters、path separators、query delimiters 和 visually ambiguous formatting。该 encoding 不是 claim container，不得嵌入 player、credential、route、permission、timestamp 或 account lifecycle data。

## 6. Key Identifier And Secret Configuration

`verifier_key_id` 标识 stored record 使用的 server-side verifier key set。它不是 secret key value。

Classification:

```yaml
verifier_key_id:
  secret_value: false
  public_api_field: false
  log_safe_by_default: false
  allowed_in_database_record: true
  allowed_in_internal_rotation_plan: true
  allowed_in_public_error: false
  allowed_in_conversation_log: false
  allowed_in_change_spec_examples: false
```

规则：

- Server-side verifier keys、peppers 和 raw secret configuration values 是 secret material。
- `verifier_key_id` 不得包含 key value、environment variable value、credential、token、会暴露 tenancy 的 cloud secret path 或 provider secret。
- Future key loading、environment variables、rotation、fallback 和 operational secret storage 需要单独的 implementation/configuration gate。
- 第一版 local posture 不需要 KMS 或 external secret-management integration。未来添加此类能力需要 dependency adoption record 和 operations boundary。

## 7. Constant-Time Comparison

Future verifier comparison 必须使用 constant-time equality 比较 verifier digests。

Acceptable future Go primitives:

```yaml
preferred_go_comparison: crypto/hmac.Equal
acceptable_go_comparison: crypto/subtle.ConstantTimeCompare
plaintext_comparison: forbidden
bytes_equal_for_verifier_digest: forbidden
string_equal_for_verifier_digest: forbidden
```

规则：

- 先计算 verifier digests，再进行 comparison。
- 比较 verifier digest bytes，而不是 raw credential 或 token strings。
- 不得用 `==`、`bytes.Equal`、string conversion、map lookup 或 database-only equality 比较 verifier digests。
- Missing lookup record 必须产生与 invalid verifier 相同的 public failure class。
- Public failures 不得泄露 failure 是由 lookup digest、verifier digest、player account、credential record、token record、key id、algorithm version 或 expiration state 导致的，除非相应 semantic error class 明确允许。

Missing records、disabled accounts、expired tokens 和 revoked tokens 的 timing equalization 可能需要后续 implementation tests。本标准要求 verifier material 的 constant-time digest comparison，但不实现 validation behavior。

## 8. Redaction Requirements

以下内容禁止出现在 logs、traces、metrics labels、public errors、panic output、audit-safe facts、test snapshots、fixtures、ADRs、change specs、documentation examples 和 conversation logs 中：

- Raw credential proof。
- Raw access-token text。
- 用作 proof 或 proof input 的 raw device identifiers。
- Credential lookup digest。
- Credential verifier digest。
- Token lookup digest。
- Token verifier digest。
- Server-side HMAC keys、peppers、seed material 和 secret configuration values。
- Full `verifier_key_id` values，除非后续 operations standard 声明了某种 specific safe representation。
- Provider secrets、OAuth tokens、JWTs、JWKs、passwords 或 password-like credentials。

Allowed with care:

- `credential_record_id`
- `token_record_id`
- `player_id`，前提是上下文 artifact 已允许 player identifiers
- lifecycle state names
- registered error codes
- non-secret reason catalog values
- future short redacted fingerprint，但必须先有 fingerprint standard

Digest values 不是 raw tokens，并不代表它们安全。它们仍然是 secret-adjacent 或 verifier material，必须避免进入 public artifacts。

## 9. Redaction Test Expectations

Future implementation work 只要处理 credential proof、access-token proof、verifier inputs、verifier digests 或 application authentication errors，就必须添加 focused redaction tests。

Minimum test expectations:

```yaml
redaction_tests:
  raw_credential_absent_from_public_error: required
  raw_access_token_absent_from_public_error: required
  raw_credential_absent_from_logs: required_when_logging_path_exists
  raw_access_token_absent_from_logs: required_when_logging_path_exists
  raw_credential_absent_from_audit_safe_facts: required_when_audit_fact_path_exists
  raw_access_token_absent_from_audit_safe_facts: required_when_audit_fact_path_exists
  verifier_digest_absent_from_public_error: required
  verifier_key_value_absent_from_all_outputs: required
  key_identifier_absent_from_public_error: required
  registered_error_code_present: required
```

Test fixtures 应使用明显的 sentinel values，例如：

```text
vibit-test-raw-credential-do-not-log
vibit-test-raw-access-token-do-not-log
vibit-test-verifier-key-do-not-log
```

这些 sentinel values 只是 synthetic test material。它们不得成为 production credential 或 token shape 示例。

## 10. Dependency Posture

第一版 high-entropy verifier posture 不需要外部 dependency。

后续 explicit code gate 授权 implementation 后可使用：

```yaml
go_standard_library:
  hmac: crypto/hmac
  sha256: crypto/sha256
  constant_time: crypto/hmac.Equal_or_crypto/subtle
  randomness: crypto/rand
  base64url: encoding/base64
```

第一版 posture 不需要且继续 deferred：

- bcrypt。
- Argon2 或 Argon2id。
- JWT、JWK 或 signing libraries。
- OAuth 或 OIDC provider SDKs。
- KMS 或 cloud secret-manager SDKs。
- Redis-like token/session stores。
- Password policy libraries。

如果 future login method 接收 password-like 或 low-entropy human input，它必须定义单独的 credential boundary 和 dependency adoption record。HMAC-SHA-256 over low-entropy passwords 不由本标准 ratify。

## 11. Ownership

Future verifier algorithm code 由 application 拥有：

```text
runtime/internal/app
```

Code 后续可以位于 application-owned child package，例如：

```text
runtime/internal/app/authentication
```

Ownership rules:

- 后续 gates 授权后，application authentication service code 拥有 token generation orchestration 和 verifier comparison。
- `authentication.Repository` 只存取 already-computed digest material；它不计算或比较 verifiers。
- PostgreSQL adapters 只持久化 digest bytes 和 key identifiers；它们不做 authentication decisions。
- Protobuf adapters 和 WebSocket transports 在后续 protocol gates 后携带 already-decoded proof fields；它们不计算 digests 或比较 verifiers。
- Generated authentication contract shapes 仍是 metadata-only 且 immutable。

## 12. Non-Goals

本标准不做以下事情：

- 添加 token generation。
- 添加 credential generation。
- 添加 verifier comparison。
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
- 添加 secret configuration loading。
- 添加 KMS、OAuth、OIDC、JWT、bcrypt 或 Argon2 dependencies。

## 13. Verification Path

该 boundary 的 repository check rule 是：

```text
runtime.token_credential_verifier_algorithm_redaction_boundary
```

触及该 boundary 的 changes 应运行：

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

本标准不验证 runtime Go authentication behavior，因为没有添加 behavior。

## 14. Follow-Up Gates

Recommended follow-up gates:

- Secret configuration and verifier key loading boundary。
- Token and credential material generation implementation gate。
- Verifier digest comparison implementation gate。
- Application authentication service implementation gate。
- Protobuf authentication message gate。
- WebSocket request proof carrier gate。
- Authentication redaction test implementation gate。
