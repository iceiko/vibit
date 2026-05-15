# Environment Verifier Key Loader Gate

状态：Draft v0.1
最后更新：2026-05-15
范围：添加 loader code 之前，process environment verifier key loader gate、environment variable contract、decoding posture、validation handoff、redaction rules、file boundaries、tests 和 deferrals
依赖：`docs/secret-configuration-verifier-key-loading-boundary.md`, `docs/local-verifier-key-configuration-loading-gate.md`
Canonical decision: `ADR-0046`

配套英文源文档是 `docs/environment-verifier-key-loader-gate.md`。英文文件是权威版本。

## 1. Purpose

本 gate 在添加 loader code 前，定义未来的 process environment verifier key loader。

现在 explicit in-memory verifier key set validator 已经存在于 `runtime/internal/app/authentication`。下一类风险是未来 Agent 在添加 environment parsing 时复制 validation rules，或者通过 decoder errors、logs、test snapshots、change specs 或 conversation logs 泄露 verifier key values。

这只是 implementation-gate standard。它不添加 Go code、imports、process environment parsing、Base64 text decoding、local secret file reading、`.env` behavior、CLI flag behavior、startup wiring、token generation、credential generation、verifier digest computation、verifier comparison、authentication service behavior、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、repository methods、SQL migrations、external secret-manager integrations、authentication dependencies 或 production authentication behavior。

## 2. Core Rule

Environment verifier key loader gate 是：

```yaml
environment_verifier_key_loader_gate: defined
implementation_authorized_by_this_standard: false
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
future_loader_slice: process_environment_verifier_key_loader
required_validator_handoff: NewVerifierKeySet
explicit_in_memory_validator_required: true
explicit_in_memory_validator_source: runtime/internal/app/authentication/verifier_key_config.go
future_loader_source: runtime/internal/app/authentication/verifier_key_env.go
future_loader_tests: runtime/internal/app/authentication/verifier_key_env_test.go
input_source: process_environment
environment_variable_contract_declared: true
base64_text_decoding_posture_declared: true
startup_wiring_status: deferred
local_secret_file_status: deferred
dotenv_status: deferred
external_kms_secret_manager_status: deferred
token_generation_status: deferred
credential_generation_status: deferred
digest_helper_status: deferred
verifier_comparison_status: deferred
authentication_service_behavior_status: deferred
protobuf_authentication_messages_status: deferred
websocket_proof_carrier_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
authentication_dependencies_status: deferred
```

未来 loader 必须保持 adapter 身份。它可以从 process environment source 收集 text values，把 key text decode 成 bytes，并把结果交给 explicit in-memory validator。它不得变成第二套 validator、startup composition point、service implementation 或 authentication behavior entry point。

## 3. Environment Variable Contract

未来 process environment contract 是：

```yaml
environment_variables:
  VIBIT_AUTH_VERIFIER_KEY_SET_ID:
    required: true
    value_kind: key_set_identifier
    decoding: trim_space_string
    log_safe_value: false
  VIBIT_AUTH_CREDENTIAL_LOOKUP_KEY:
    required: true
    value_kind: verifier_key_bytes
    decoding: base64_text_to_bytes
    log_safe_value: false
  VIBIT_AUTH_CREDENTIAL_VERIFIER_KEY:
    required: true
    value_kind: verifier_key_bytes
    decoding: base64_text_to_bytes
    log_safe_value: false
  VIBIT_AUTH_TOKEN_LOOKUP_KEY:
    required: true
    value_kind: verifier_key_bytes
    decoding: base64_text_to_bytes
    log_safe_value: false
  VIBIT_AUTH_TOKEN_VERIFIER_KEY:
    required: true
    value_kind: verifier_key_bytes
    decoding: base64_text_to_bytes
    log_safe_value: false
```

规则：

- 第一版 process environment posture 要求全部五个变量存在。
- Environment variable names 是 log-safe。
- Environment variable values 不是 log-safe。
- Full concrete key set id value 默认不是 log-safe。
- Missing values 必须 fail closed。
- Empty 或 whitespace-only key set id 必须通过 explicit validator fail closed。
- Empty key text 必须在 decode 前后或通过 explicit validator fail closed。
- Decoded key bytes 必须传给 `NewVerifierKeySet`。
- Loader 不得削弱 validator 对 missing、short、duplicate、all-zero 或 repeated-byte keys 的要求。

## 4. Decoding Posture

未来 loader 只能使用 Go standard library 进行 decoding。

后续 implementation work item 授权 code 后允许的 future imports：

```yaml
future_standard_library_imports_allowed:
  - encoding/base64
  - os
  - strings
  - errors
  - fmt
```

Decoding policy：

```yaml
preferred_encoding: base64url_unpadded
compatibility_encoding: standard_base64_padded
invalid_base64_text: fail_closed
raw_unencoded_key_text: forbidden
hex_key_text: forbidden
json_key_blob: forbidden
partial_key_set: forbidden
```

未来第一段 implementation 应优先使用 URL-safe unpadded Base64，因为它适合 environment variables 和 shell configuration。为了 operator ergonomics，如果 tests 覆盖两种路径，也可以接受 standard padded Base64。它不得接受 raw key text、hex text、JSON blobs、comma-delimited values 或 partial key sets。

## 5. Package And File Boundary

未来 implementation ownership：

```text
runtime/internal/app/authentication
```

后续 implementation work item 授权 code 后允许的 future loader files：

```text
runtime/internal/app/authentication/verifier_key_env.go
runtime/internal/app/authentication/verifier_key_env_test.go
```

未来 loader 可以使用已有 validator：

```text
runtime/internal/app/authentication/verifier_key_config.go
```

Environment loader implementation slice 的 forbidden write areas：

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

未来 implementation 不得把 loader 接入 process startup。Startup wiring 是单独的 composition decision。

## 6. Future Loader Shape

首选 future shape 是：

```yaml
EnvLookup:
  type: function
  signature: "func(name string) (string, bool)"

future_functions:
  LoadVerifierKeySetFromEnvironment:
    input: EnvLookup
    output: VerifierKeySet
    behavior: decode_required_environment_values_then_call_NewVerifierKeySet
  LoadVerifierKeySetFromProcessEnvironment:
    input: none
    output: VerifierKeySet
    behavior: small_os_LookupEnv_adapter_only
```

可测试的 loader 应接收 explicit lookup function，这样 tests 在更简单时可以不修改 global process environment。后续 implementation gate 授权后，一个很小的 process adapter 可以调用 `os.LookupEnv`，但它仍不得 wire into server startup。

## 7. Error And Redaction Requirements

未来 loader 应暴露 typed 或 sentinel errors，让 tests 能使用，同时不披露 secret material。

允许的 error classification examples：

```yaml
error_classes:
  missing_environment_variable
  invalid_environment_key_encoding
  invalid_environment_key_set
```

允许出现在 public error text 中：

- Environment variable names。
- Error classes。
- Logical key purposes。

禁止出现在 errors、logs、test snapshots、docs、ADRs、change specs 和 conversation logs 中：

- Environment variable values。
- Encoded verifier key values。
- Decoded verifier key bytes。
- Full concrete key set id values。
- Credentials。
- Access tokens。
- Lookup digests。
- Verifier digests。
- HMAC inputs。
- HMAC outputs。
- Deployment-specific identifiers。

未来 loader 必须 wrap 或 map validation errors，但不得在错误文本中添加 secret values。

## 8. Required Tests For The Future Loader

未来 implementation 必须在这里添加 focused unit tests：

```text
runtime/internal/app/authentication/verifier_key_env_test.go
```

Minimum test cases：

```yaml
required_tests:
  accepts_complete_base64url_unpadded_environment_key_set
  accepts_complete_standard_base64_padded_environment_key_set
  missing_each_environment_variable_fails_closed
  invalid_each_encoded_key_fails_closed
  decoded_short_key_fails_through_validator
  duplicate_decoded_keys_fail_through_validator
  all_zero_decoded_key_fails_through_validator
  repeated_single_byte_decoded_key_fails_through_validator
  loader_calls_explicit_in_memory_validator
  errors_include_environment_variable_name_when_safe
  errors_do_not_include_environment_variable_values
  errors_do_not_include_full_key_set_id
  process_environment_adapter_is_small_and_unwired
```

Tests 不得要求 PostgreSQL、MinIO、WebSocket transport、Protobuf generation、KMS、cloud SDKs、local files、`.env` files 或 external services。

## 9. Dependency Posture

本 gate 不允许添加任何 new external dependency。

未来 loader 应只使用 Go standard library。Major external KMS、cloud secret-manager、operations、cryptography、password-hashing、JWT、OAuth、OIDC、provider 或 dotenv dependency 都需要单独的 adoption record 和 operations boundary。

## 10. Nakama And Pitaya Mapping

Nakama capability reference：

- Server-side authentication secret material 必须在 authentication token behavior 有 production 意义前可靠配置。
- Session token validation 依赖 stable server-owned verifier material。

Pitaya capability reference：

- Realtime route handlers 应在 proof validation 后接收 identity context。
- Server/session infrastructure 不应拥有 verifier key material。

vibit adaptation：

- Environment configuration 是 application-owned adapter。
- 承载 invariants 的 verifier key set validator 仍是 validation truth source。
- Transport、protocol、repositories、generated code 和 domain modules 不 parse 或持有 verifier keys。
- Startup composition 与 loader implementation 保持分离。

## 11. Deferrals

本 gate 不授权：

- Go code by itself。
- 本 change 中的 process environment parsing。
- 本 change 中的 Base64 text decoding。
- Local secret file reading。
- `.env` parsing。
- CLI flag secret input。
- Startup wiring。
- Token generation。
- Credential generation。
- Verifier digest computation。
- Verifier comparison。
- Authentication service behavior。
- Login execution。
- Access-token validation。
- Logout execution。
- Refresh execution。
- Cleanup jobs。
- Protobuf authentication messages。
- WebSocket proof carriers。
- WebSocket handshake authentication。
- Repository interface changes。
- PostgreSQL adapter changes。
- SQL migration changes。
- External KMS 或 secret-manager integrations。
- Authentication dependencies。
- Production authentication behavior。

## 12. Verification

本 gate 的 repository check rule 是：

```yaml
runtime.environment_verifier_key_loader_gate
```

Verification for this gate 必须证明：

- English standard 和 Simplified Chinese translation 存在。
- ADR-0046 存在。
- Architecture manifests 和 agent guides reference this gate。
- Environment variable contract 已声明。
- Future loader 必须 hand off to `NewVerifierKeySet`。
- 本 gate 没有添加 process environment parsing 或 Base64 decoding code。
- Deferrals 在 manifests 中保持可见。

## 13. Migration Path

1. 以 standards and manifest work only 完成本 gate。
2. 在后续 code slice 中实现 environment verifier key loader。
3. Startup wiring 保持在单独 composition gate 之后。
4. KMS、cloud secret-manager、`.env`、local secret files、CLI secret input 和 production operations posture 保持在单独 decisions 之后。
