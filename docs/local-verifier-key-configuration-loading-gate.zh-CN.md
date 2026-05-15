# Local Verifier Key Configuration Loading Gate

状态：Draft v0.1
最后更新：2026-05-15
范围：添加 secret-loading code 之前，第一段 local verifier key configuration loading implementation gate、explicit input posture、validation rules、redaction rules、package boundaries、test requirements 和 deferrals
依赖：`docs/secret-configuration-verifier-key-loading-boundary.md`, `docs/verifier-digest-computation-comparison-boundary.md`, `docs/authentication-service-implementation-readiness-gate.md`
Canonical decision: `ADR-0045`

配套英文源文档是 `docs/local-verifier-key-configuration-loading-gate.md`。英文文件是权威版本。

## 1. Purpose

本 gate 在添加 code 前，定义第一段 bounded local verifier key configuration loading implementation slice。

它存在的原因是：key configuration 是 security-sensitive 的，但第一段 implementation 仍然必须小到 agents 能实现和验证。未来实现应该先证明一个完整的 in-memory verifier key set 可以被 validate、redact 和 test，而不把核心 validator 耦合到 process environment parsing、KMS providers、cloud secret managers、protocol carriers、repositories 或 login behavior。

这只是 implementation-gate standard。它不添加 Go code、imports、secret loading、environment parsing、token generation、credential generation、verifier digest computation、verifier comparison、authentication service behavior、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、repository methods、SQL migrations、external secret-manager integrations、authentication dependencies 或 production authentication behavior。

## 2. Core Rule

第一版 local verifier key configuration loading gate 是：

```yaml
local_verifier_key_configuration_loading_gate: defined
implementation_authorized_by_this_standard: false
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
first_implementation_slice: explicit_in_memory_verifier_key_set_validation
first_environment_loader_status: deferred_to_follow_up_gate
process_environment_parsing_status: deferred
base64_text_decoding_status: deferred
external_kms_secret_manager_status: deferred
token_generation_status: deferred
credential_generation_status: deferred
digest_helper_status: deferred
verifier_comparison_status: deferred
authentication_service_behavior_status: deferred
login_execution_status: deferred
token_validation_status: deferred
protobuf_authentication_messages_status: deferred
websocket_proof_carrier_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
authentication_dependencies_status: deferred
```

第一段 implementation 应该只 validate explicit in-memory input。Environment variable loading 被故意排在 validator 存在之后。

原因很直接：validator 是承载 invariants 的核心，而 environment parsing 是 adapter。Agents 应该能在没有 process environment、local files、shell quoting 或 deployment assumptions 的情况下测试这个 invariant-bearing core。

## 3. Package And File Boundary

未来 implementation ownership：

```text
runtime/internal/app/authentication
```

后续 implementation work item 授权 code 后，允许的 first-slice files：

```text
runtime/internal/app/authentication/verifier_key_config.go
runtime/internal/app/authentication/verifier_key_config_test.go
```

如果第一段 implementation 仍保持很小，允许的 helper file：

```text
runtime/internal/app/authentication/errors.go
```

Forbidden first-slice write areas：

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

第一段 implementation 不得把 configuration 接入 process startup。它应该定义一个很小的 application-owned validation primitive，供后续 loaders 调用。

## 4. First Input Shape

第一段 code slice 应接收已经 decode 为 bytes 的 explicit in-memory data：

```yaml
VerifierKeySetConfig:
  key_set_id: string
  credential_lookup_key: []byte
  credential_verifier_key: []byte
  token_lookup_key: []byte
  token_verifier_key: []byte
```

接受后的 output 应是 immutable validated key-set value：

```yaml
VerifierKeySet:
  key_set_id: internal_non_log_safe_identifier
  credential_lookup_key: private_bytes
  credential_verifier_key: private_bytes
  token_lookup_key: private_bytes
  token_verifier_key: private_bytes
```

规则：

- Input byte slices 必须在存储前 copy。
- Output accessors 不得暴露 mutable internal slices。
- Key set id 是必填的，但默认不是 log-safe。
- Key set id 不是 key、credential、token，也不是 public API value。
- 该 value 不得通过 `fmt.Stringer` 暴露 secrets。
- Error values 不得包含 key bytes、encoded key values、full concrete key ids、environment variable values 或 deployment identifiers。

## 5. Validation Rules

第一版 validator 必须拒绝：

```yaml
validation_failures:
  missing_key_set_id: fail_closed
  missing_credential_lookup_key: fail_closed
  missing_credential_verifier_key: fail_closed
  missing_token_lookup_key: fail_closed
  missing_token_verifier_key: fail_closed
  decoded_key_shorter_than_32_bytes: fail_closed
  duplicate_logical_key_bytes: fail_closed
  all_zero_key_bytes: fail_closed
  obvious_repeated_single_byte_key: fail_closed
```

Minimum key length：

```yaml
decoded_key_minimum_length_bytes: 32
minimum_key_entropy_bits: 256
```

第一版 validator 无法证明 caller-supplied bytes 的 entropy，但它必须拒绝表示误用的常见 invalid shapes。Production key generation requirements 仍由 secret configuration boundary 定义。

## 6. Error And Redaction Requirements

未来第一段 implementation 应暴露 typed 或 sentinel errors，让 tests 能使用，同时不披露 secret material。

允许的 error classification examples：

```yaml
error_classes:
  missing_key_set_id
  missing_required_key
  key_too_short
  duplicate_logical_key
  weak_repeated_key
```

禁止出现在 errors、logs、test snapshots、docs、ADRs、change specs 和 conversation logs 中：

- Verifier key bytes。
- Encoded verifier key values。
- Environment variable values。
- Full concrete `verifier_key_id` values。
- Credentials。
- Access tokens。
- Lookup digests。
- Verifier digests。
- HMAC inputs。
- HMAC outputs。

Tests 可以使用明显的 sentinel strings 或 short byte slices 来证明 redaction，但不得提交 production-like keys。

## 7. Environment Loading Sequence

Environment loading 不属于第一段 implementation slice。

后续 environment loader 应调用 explicit validator，而不是复制 validation rules。后续 gate 可以授权：

```yaml
future_environment_loader:
  variable_names:
    - VIBIT_AUTH_VERIFIER_KEY_SET_ID
    - VIBIT_AUTH_CREDENTIAL_LOOKUP_KEY
    - VIBIT_AUTH_CREDENTIAL_VERIFIER_KEY
    - VIBIT_AUTH_TOKEN_LOOKUP_KEY
    - VIBIT_AUTH_TOKEN_VERIFIER_KEY
  decoding: base64url_unpadded_or_standard_base64_to_bytes
  input_source: process_environment
  redaction_required: true
```

本 gate 不授权 `os.Getenv`、`os.LookupEnv`、`.env` parsing、local file reading、CLI flags、startup wiring 或 environment-driven runtime authentication。

## 8. Required Tests For The Future First Code Slice

未来 implementation 必须在这里添加 focused unit tests：

```text
runtime/internal/app/authentication/verifier_key_config_test.go
```

Minimum test cases：

```yaml
required_tests:
  accepts_complete_distinct_32_byte_key_set
  copies_input_key_material
  accessors_do_not_expose_mutable_internal_slices
  missing_key_set_id_fails_closed
  missing_each_logical_key_fails_closed
  short_each_logical_key_fails_closed
  duplicate_logical_keys_fail_closed
  all_zero_key_fails_closed
  repeated_single_byte_key_fails_closed
  errors_do_not_include_secret_bytes
  errors_do_not_include_full_key_set_id
```

Tests 不得要求 PostgreSQL、MinIO、WebSocket transport、Protobuf generation、process environment variables、KMS、cloud SDKs 或 external services。

## 9. Dependency Posture

本 gate 不允许添加任何 new external dependency。

未来第一段 implementation 应只使用 Go standard library。Major external KMS、cloud secret-manager、operations、cryptography、password-hashing、JWT、OAuth、OIDC、provider 或 dotenv dependency 都需要单独的 adoption record 和 operations boundary。

## 10. Nakama And Pitaya Mapping

Nakama capability reference：

- Account authentication 和 session token validation 依赖可信的 server-side secret material。
- Secret handling 应保持 server-owned 且 redacted。

Pitaya capability reference：

- Handler identity context 应在 proof validation 后接收 validated identity。
- Transport/session state 不应拥有 verifier key material。

vibit adaptation：

- Verifier key configuration 由 application 拥有。
- Key validation 先于 digest helpers 和 service behavior。
- Transport、protocol、repositories、generated code 和 domain modules 不 parse 或持有 verifier keys。

## 11. Deferrals

本 gate 不授权：

- Go code by itself。
- Environment variable parsing。
- Base64 text decoding。
- `.env` parsing。
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

```text
runtime.local_verifier_key_configuration_loading_gate
```

该 check 应验证：

- 本 standard、translation 和 ADR 存在。
- Architecture manifests 和 agent guides 引用了这个 gate。
- Required markers 将 explicit in-memory validation 标识为第一段 implementation slice。
- Runtime code 没有基于本 gate 实现 environment parsing、secret loading、token generation、digest helpers、verifier comparison、authentication behavior、protocol carriers、repository changes、migrations、KMS、secret-manager integration 或 new authentication dependencies。

Future code verification 必须包含在 `runtime/` 下运行的 `go test ./...`，以及 `node tools/vibit check runtime --json`。
