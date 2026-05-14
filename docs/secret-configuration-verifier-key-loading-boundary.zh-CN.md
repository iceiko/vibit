# Secret Configuration And Verifier Key Loading Boundary

状态：Draft v0.1
最后更新：2026-05-15
范围：第一版 device-credential 和 opaque access-token 姿态的 secret configuration ownership、verifier key separation、future key loading posture、key identifier handling、rotation expectations、development/test posture、production failure behavior 和 redaction requirements
依赖：`docs/token-credential-verifier-algorithm-redaction-boundary.md`, `docs/application-authentication-service-interface-boundary.md`, `docs/runtime-authentication-implementation-boundary.md`
Canonical decision: `ADR-0041`

配套英文源文档是 `docs/secret-configuration-verifier-key-loading-boundary.md`。英文文件是权威版本。

## 1. Purpose

本标准定义未来 application-owned authentication service code 必须遵守的 secret configuration 和 verifier key loading boundary。

它位于 secret loading、token material generation、credential material generation、verifier comparison、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、authentication dependencies、repository interface changes、migration schema changes 或 production authentication behavior 之前。

这只是 boundary-only standard。它不添加 Go code、imports、runtime services、environment parsing、secret loading、token generation、credential generation、digest comparison、repository methods、SQL migrations、Protobuf messages、WebSocket carriers、routes、KMS integration、secret-manager integration 或 production authentication behavior。

## 2. Core Rule

第一版 secret configuration posture 是：

```yaml
secret_configuration_boundary: defined
implementation_authorized_by_this_standard: false
planned_owner: runtime/internal/app
planned_optional_child_package: runtime/internal/app/authentication
first_local_key_source: process_environment_or_explicit_runtime_secret_input
external_kms_secret_manager_required_for_first_local_posture: false
external_secret_manager_dependency_adoption_record_required: true
verifier_key_set_required: true
verifier_key_id_required_for_stored_records: true
minimum_key_entropy_bits: 256
decoded_key_minimum_length_bytes: 32
key_text_encoding: base64url_unpadded_or_standard_base64_decoded_to_bytes
production_missing_key_behavior: fail_closed
development_default_production_keys: forbidden
automatic_persistent_random_key_generation: forbidden
runtime_secret_values_in_committed_artifacts: forbidden
```

在后续 code gate 授权 implementation 后，第一版 local implementation 可以使用 process environment configuration 或 explicit runtime secret input。它不需要 KMS 或 cloud secret-manager dependency。

Production secret storage、external secret managers、KMS providers、cloud provider SDKs、operational rotation systems 或 container orchestration secret integrations，都必须先通过后续 dependency adoption record 和 operations boundary 才能实现。

## 3. Ownership

Future secret configuration loading 由 application 拥有：

```text
runtime/internal/app
```

Future code 可以使用 application-owned child package，例如：

```text
runtime/internal/app/authentication
```

Ownership rules：

- Application authentication code 在后续 implementation gate 授权后，可以 load、validate 和 hold verifier key material。
- `authentication.Repository` 只能存取 already-computed digests、status fields、timestamps 和 `verifier_key_id` values。
- PostgreSQL adapters 只持久化 key identifiers 和 digest bytes；它们不 load keys、不 decode secret configuration、不 compute HMAC、不 compare verifiers、不 rotate keys，也不决定 authentication outcomes。
- WebSocket transport 和 Protobuf protocol adapters 不读取 secret configuration，也不拥有 key loading。
- Generated authentication contract shapes 保持 metadata-only 且 immutable。

不得为了方便把 secret configuration 移入 transport、protocol、domain modules、generated output、migration files、SQL fixtures 或 repository adapters。

## 4. Logical Key Set

第一版 verifier posture 需要四个独立的 logical server-side keys：

```yaml
key_set_logical_keys:
  credential_lookup_key:
    purpose_label: vibit.credential.lookup.v1
    digest_class: credential_lookup_digest
    reuse_with_other_digest_classes: forbidden
  credential_verifier_key:
    purpose_label: vibit.credential.verifier.v1
    digest_class: credential_verifier_digest
    reuse_with_other_digest_classes: forbidden
  token_lookup_key:
    purpose_label: vibit.access_token.lookup.v1
    digest_class: token_lookup_digest
    reuse_with_other_digest_classes: forbidden
  token_verifier_key:
    purpose_label: vibit.access_token.verifier.v1
    digest_class: token_verifier_digest
    reuse_with_other_digest_classes: forbidden
```

规则：

- Lookup keys 和 verifier keys 必须分离。
- Credential keys 和 token keys 必须分离。
- 一个 logical key 泄露时，不应直接暴露所有 digest class。
- Key set 由内部 `verifier_key_id` 选择。
- 新的 credential 和 token verifier records 必须存储用于计算 digest 的 `verifier_key_id`。
- Future verifier code 必须拒绝不完整 key set，而不是静默复用同一个 key 到多个 digest class。

## 5. Planned Local Configuration Surface

后续 code gate 授权后，第一版 local implementation 可以使用这些 planned environment variable names：

```yaml
planned_environment_variables:
  verifier_key_set_id: VIBIT_AUTH_VERIFIER_KEY_SET_ID
  credential_lookup_key: VIBIT_AUTH_CREDENTIAL_LOOKUP_KEY
  credential_verifier_key: VIBIT_AUTH_CREDENTIAL_VERIFIER_KEY
  token_lookup_key: VIBIT_AUTH_TOKEN_LOOKUP_KEY
  token_verifier_key: VIBIT_AUTH_TOKEN_VERIFIER_KEY
```

这些名称是 planned configuration contract，不是 implementation。本标准不授权 `os.Getenv`、environment parsing、config structs、CLI flags、process wiring 或 runtime key loading。

Future implementation rules：

- Environment variable values 是 secret material；key set id 是 secret-adjacent，默认也不是 log-safe。
- Key values decode 后必须至少有 32 bytes。
- Production key values 必须来自 cryptographically secure randomness。
- Key values 不得是 human-readable passwords、passphrases、repeated-pattern strings、all-zero values、user identifiers、device identifiers、provider subjects 或 copied tokens。
- Production key configuration 缺失、格式错误、过短、重复或不完整时，必须 fail closed。
- Developers 可以用 local ignored files 填充 environment variables，但 committed file 不得包含 production-like values。
- 本仓库 ignored `.vibit.local.env` 约定只用于 local machine configuration，不得被当作 committed secret source。

## 6. Key Identifier Rules

`verifier_key_id` 标识 stored credential 或 token verifier record 使用的 logical verifier key set。它不是 secret key value。

Classification：

```yaml
verifier_key_id:
  secret_value: false
  public_api_field: false
  database_record_field: true
  internal_selection_input: true
  log_safe_by_default: false
  public_error_safe: false
  documentation_example_safe: placeholder_only
  conversation_log_safe: placeholder_only
  change_spec_example_safe: placeholder_only
```

规则：

- `verifier_key_id` 不得包含 key bytes、encoded key values、credentials、tokens、account identifiers、tenant identifiers、cloud secret paths、provider secret names、access keys、environment variable values、hostnames、deployment names 或 operator names。
- Public errors 不得泄露 key id、key-set miss、key rotation state、key decoding error 或 key length failure。
- Logs、traces、metrics labels、audit-safe facts、ADRs、change specs、docs 和 conversation logs 只能使用 `<verifier-key-id>` 这样的 placeholder。
- Future operations standard 可以定义 short redacted fingerprint format，但本标准不 ratify 任何 fingerprint format。

## 7. Rotation Posture

Future key rotation 必须以 key-set 为单位。

```yaml
rotation_model:
  active_key_set: required
  accepted_previous_key_sets: allowed
  new_writes_use_active_key_set: required
  stored_records_keep_verifier_key_id: required
  verification_may_select_key_by_record_key_id: required
  automatic_rotation_implemented_by_this_standard: false
```

Expected rotation phases：

1. 引入新的完整 key set。
2. 将新 key set 标记为 active，供新的 credential 和 token verifier records 使用。
3. 在 access-token TTL、credential replacement window 或其他后续 ratified retention window 内保留 previous key sets。
4. 按后续 authentication behavior gate reissue、rotate、revoke 或 replace records。
5. 只有在没有 valid record 需要 old key set 后，才能 retire old key sets。

Rotation rules：

- Rotation 不得要求手工编辑 stored digests。
- Rotation 不得把 database-only equality 当作 verifier proof。
- Rotation failures 必须使用与其他 invalid proof failures 相同的 public failure class，除非后续 semantic error standard 明确允许更多细节。
- Key retirement 必须能在内部被观察，但不得在 public artifacts 中暴露 raw keys 或 full key identifiers。

本标准不实现 rotation behavior。

## 8. Development And Test Posture

Development 和 test configuration 必须保持 explicit。

后续 implementation gates 授权后允许：

- Test code 内的 test-only deterministic fixture keys，前提是它们明显非 production，且只用于 repeatable unit tests。
- 为 developer machine 设置 environment variables 的 local ignored files。
- Tests 提供的 explicit in-memory test configuration。
- 使用明显 sentinel strings 证明 secret values 不泄露的 redaction tests。

禁止：

- Committed production-like key values。
- Default production keys。
- Automatic persistent random key generation。
- Silent fallback 到 shared hard-coded development key。
- 把 key set id 当作 secret key。
- 使用 raw credential material、raw access-token material、provider tokens、user identifiers、device identifiers 或 player identifiers 作为 key material。

Automatic random keys 只能在后续 gate 明确授权的 ephemeral in-memory development mode 中允许。它们不得用于 durable local authentication state，因为重启会让 stored verifier records 无法验证。

## 9. Failure And Fallback Behavior

Future production behavior 在 required secret configuration 不可用或无效时必须 fail closed。

```yaml
failure_behavior:
  missing_required_key: fail_closed
  malformed_key_encoding: fail_closed
  decoded_key_too_short: fail_closed
  duplicate_logical_keys: fail_closed
  incomplete_key_set: fail_closed
  unknown_record_key_id: invalid_proof_public_failure
  retired_record_key_id: invalid_proof_public_failure
  public_error_discloses_secret_config_problem: forbidden
```

规则：

- Key configuration 无效时，startup 可以在 serving requests 前失败。
- Request-time validation 不得泄露 failure 是由 key configuration、missing records、invalid proof、expiration、revocation、account state、algorithm version 或 key id 引起，除非后续 semantic error standard 明确允许。
- Development behavior 可以使用更清晰的 local diagnostics，但 diagnostics 必须 redact secret values，且不得进入 public client responses。
- Metrics labels 不得包含 key values、encoded key values、full key ids、credentials、tokens 或 cloud secret paths。

## 10. Redaction Requirements

以下内容禁止出现在 logs、traces、metrics labels、public errors、panic output、audit-safe facts、test snapshots、fixtures、ADRs、change specs、documentation examples 和 conversation logs 中：

- Verifier key values。
- Encoded verifier key values。
- Decoded verifier key bytes。
- Environment variable values。
- Secret-manager response payloads。
- 会暴露 tenancy 或 deployment structure 的 cloud secret paths。
- Provider credentials 或 provider secret names。
- Raw credential proof。
- Raw access-token text。
- Credential lookup digest。
- Credential verifier digest。
- Token lookup digest。
- Token verifier digest。
- Full concrete `verifier_key_id` values。

Allowed with care：

- Environment variable names。
- `<verifier-key-id>` 这样的 placeholder key identifiers。
- Non-secret configuration field names。
- Registered error codes。
- Internal lifecycle state names。

Secret configuration values 不能因为“只是本地”就被视为安全。Agent-facing documents 和 conversation logs 必须把 local secrets 当作真实 secrets 处理。

## 11. Dependency Posture

第一版 local secret configuration posture 不需要 external dependency。

```yaml
external_kms_secret_manager_required_for_first_local_posture: false
process_environment_allowed_after_code_gate: true
explicit_runtime_secret_input_allowed_after_code_gate: true
dependency_adoption_record_required_for_external_secret_manager: true
operations_boundary_required_for_external_secret_manager: true
```

本标准 deferred 且未选择：

- KMS providers。
- Cloud secret-manager SDKs。
- Vault-like systems。
- Container orchestration secret APIs。
- Password-hashing dependencies。
- JWT、JWK、OAuth 或 OIDC dependencies。
- Provider SDKs。

Future production deployment guide 可以要求 external secret manager，但该决定属于单独的 dependency 和 operations gate。

## 12. Non-Goals

本标准不：

- 添加 secret loading。
- 添加 environment parsing。
- 添加 config structs。
- 添加 CLI flags。
- 添加 KMS 或 cloud secret-manager integration。
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
- 添加 production authentication behavior。

## 13. Verification Path

该 boundary 的 repository check rule 是：

```text
runtime.secret_configuration_verifier_key_loading_boundary
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

本标准不添加 runtime Go authentication behavior，因此不验证 runtime Go authentication behavior。

## 14. Follow-Up Gates

Recommended follow-up gates：

- Token and credential material generation boundary。
- Verifier digest computation and constant-time comparison implementation boundary。
- Application authentication service implementation gate。
- Authentication redaction test implementation gate。
- Protobuf authentication message gate。
- WebSocket request proof carrier gate。
- Operations secret-management adoption gate。
