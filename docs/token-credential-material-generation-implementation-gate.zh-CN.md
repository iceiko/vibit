# Token And Credential Material Generation Implementation Gate

状态：Draft v0.1
最后更新：2026-05-15
范围：添加 generation code 之前，token and credential material generation helper implementation gate、future helper file boundaries、entropy source、encoding posture、redaction rules、tests 和 deferrals
依赖：`docs/token-credential-material-generation-boundary.md`, `docs/verifier-digest-computation-comparison-boundary.md`, `docs/authentication-service-implementation-readiness-gate.md`, `docs/environment-verifier-key-loader-gate.md`
Canonical decision: `ADR-0047`

配套英文源文档是 `docs/token-credential-material-generation-implementation-gate.md`。英文文件是权威版本。

## 1. Purpose

本 gate 定义未来 raw device credential 与 opaque access-token material generation helpers 的下一段 bounded implementation slice。

仓库已经定义了 material generation boundary。下一类风险是未来 Agent 把 secrets 生成放进 service orchestration、repository adapters、transport handlers、test fixtures 或 protocol code，或者把 generation 与 verifier digest computation、authentication behavior 混在一起。

这只是 implementation-gate standard。它不添加 Go code、imports、token generation、credential generation、verifier digest computation、verifier comparison、authentication service behavior、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、repository methods、SQL migrations、startup wiring、authentication dependencies、external randomness services、KMS、cloud secret-manager integrations 或 production authentication behavior。

## 2. Core Rule

Token and credential material generation implementation gate 是：

```yaml
token_credential_material_generation_implementation_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0101
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
future_slice: raw_device_credential_and_access_token_material_helpers
future_source: runtime/internal/app/authentication/material_generation.go
future_tests: runtime/internal/app/authentication/material_generation_test.go
production_entropy_source: crypto/rand.Reader
testable_entropy_handoff: io.Reader
random_read_primitive: io.ReadFull
raw_material_size_bytes: 32
minimum_entropy_bits: 256
text_encoding: base64.RawURLEncoding
encoded_text_shape: url_safe_unpadded_base64
encoded_text_length_chars: 43
raw_material_copying_required: true
all_zero_material_fails_closed: true
repeated_single_byte_material_fails_closed: true
one_time_client_visible_presentation_required: true
raw_material_storage: forbidden
raw_material_repository_handoff: forbidden
digest_computation_status: deferred
verifier_comparison_status: deferred
authentication_service_behavior_status: deferred
protobuf_authentication_messages_status: deferred
websocket_proof_carrier_status: deferred
repository_interface_change_status: deferred
migration_schema_change_status: deferred
startup_wiring_status: deferred
authentication_dependencies_status: deferred
production_authentication_behavior_status: deferred
```

未来 implementation 必须是 helper-only slice。它可以创建 random raw bytes、拒绝 supplied reader 产生的 malformed 或 weak generated material、把 material encode 成 one-time presentation text，并返回给后续 application-owned digest helpers 使用的小型 value object。它不得 compute lookup digests、compute verifier digests、compare verifiers、写 repositories、选择 accounts、发 login responses、validate tokens、parse bearer proofs 或接触 protocol carriers。

## 3. Future Helper Shape

未来 implementation ownership：

```text
runtime/internal/app/authentication
```

后续 implementation work item 授权 code 后允许的 future files：

```text
runtime/internal/app/authentication/material_generation.go
runtime/internal/app/authentication/material_generation_test.go
```

首选 future API shape：

```yaml
future_types:
  GeneratedSecretMaterial:
    owns: copied_raw_bytes_and_encoded_text
    methods:
      - Kind() MaterialKind
      - RawBytes() []byte
      - Text() string
    constraints:
      - RawBytes returns a copy.
      - Text returns URL-safe unpadded Base64 presentation text.
      - Error text and string formatting must not expose raw bytes or text.
  MaterialKind:
    values:
      - device_credential
      - access_token

future_functions:
  GenerateDeviceCredentialMaterial:
    signature: "func GenerateDeviceCredentialMaterial(random io.Reader) (GeneratedSecretMaterial, error)"
    behavior: read_32_random_bytes_validate_encode_as_device_credential_material
  GenerateAccessTokenMaterial:
    signature: "func GenerateAccessTokenMaterial(random io.Reader) (GeneratedSecretMaterial, error)"
    behavior: read_32_random_bytes_validate_encode_as_access_token_material
```

明确的 `io.Reader` handoff 是为了 testability。未来 production service code 可以在自己的 service behavior gate 授权调用路径后传入 `crypto/rand.Reader`。Helper implementation 本身只有在保留 explicit reader seam 且不 wire process startup 或 service behavior 时，才可以 import `crypto/rand`。

## 4. Entropy And Encoding

未来 helper 必须先生成 canonical raw material，再进行 encoding。

```yaml
raw_material:
  bytes: 32
  entropy_bits: at_least_256
  production_source: crypto/rand.Reader
  read_primitive: io.ReadFull
  nil_reader: fail_closed
  short_read: fail_closed
  read_error: fail_closed
  all_zero_bytes: fail_closed
  repeated_single_byte_bytes: fail_closed

encoded_material:
  encoding: base64.RawURLEncoding
  alphabet: url_safe_base64
  padding: forbidden
  expected_length_chars_for_32_bytes: 43
  whitespace: forbidden
  control_characters: forbidden
  path_separators: forbidden
  query_delimiters: forbidden
  claims_or_metadata: forbidden
```

规则：

- Encoding 只是 presentation。它不得 embed player id、credential id、token id、session id、timestamp、route name、permissions、account state、provider subject 或 claims。
- Decode encoded text 必须能精确恢复 generated raw bytes。
- Device credential material 和 access-token material 共享 byte shape，但必须携带不同的 `MaterialKind` values。
- Helper 不得让 generated values 在非 test builds 中变成 stable。
- 如果 randomness 失败，helper 不得 silent retry forever。

## 5. Package And File Boundary

允许的未来 implementation area：

```text
runtime/internal/app/authentication/material_generation.go
runtime/internal/app/authentication/material_generation_test.go
```

Future material generation helper slice 的 forbidden write areas：

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

未来 implementation 不得把 generation wire into login execution、account creation、token issuance、startup、WebSocket transport、Protobuf protocol、PostgreSQL persistence、generated contract output、migrations 或 domain repository code。

## 6. Error And Redaction Requirements

未来 helper 应暴露 typed 或 sentinel errors，让 tests 能使用，同时不披露 generated material。

允许的 error classification examples：

```yaml
error_classes:
  missing_random_source
  random_read_failed
  invalid_generated_material
```

允许出现在 error text 中：

- Error classes。
- Material kind names。
- 非 secret 的 numeric expectations，例如 `32` bytes。

禁止出现在 errors、logs、test snapshots、docs、ADRs、change specs 和 conversation logs 中：

- Raw device credential text。
- Raw device credential bytes。
- Raw access-token text。
- Raw access-token bytes。
- Encoded generated material。
- Token 或 credential prefixes。
- Randomness seeds。
- Lookup digests。
- Verifier digests。
- Verifier keys。
- Encoded verifier keys。
- Full concrete `verifier_key_id` values。

Generated material 不会因为 short-lived 就安全。One-time presentation 的含义是未来只能有一个 client-visible delivery path，不是“可以 log 一次”。

## 7. Required Tests For The Future Helper

未来 implementation 必须在这里添加 focused unit tests：

```text
runtime/internal/app/authentication/material_generation_test.go
```

Minimum test cases：

```yaml
required_tests:
  device_credential_material_uses_32_random_bytes
  access_token_material_uses_32_random_bytes
  encoded_material_is_base64url_unpadded
  encoded_material_length_is_43_characters
  encoded_material_round_trips_to_raw_bytes
  generated_material_kind_is_preserved
  raw_bytes_are_copied_on_return
  nil_random_source_fails_closed
  random_source_error_fails_closed
  short_random_read_fails_closed
  all_zero_generated_material_fails_closed
  repeated_single_byte_generated_material_fails_closed
  generated_values_are_not_constant_with_progressing_source
  errors_do_not_include_raw_or_encoded_material
  helper_does_not_compute_digests_or_compare_verifiers
```

Tests 可以只通过 explicit test seam 使用 deterministic readers。它们不得引入 committed production-like secrets、stable generated tokens、stable generated credentials、包含 material 的 logs、包含 material 的 repository fixtures 或包含 material 的 protocol fixtures。

## 8. Dependency Posture

本 gate 不允许添加任何 new external dependency。

后续 implementation work item 授权 code 后允许的 Go standard library packages：

```yaml
future_standard_library_imports_allowed:
  - crypto/rand
  - encoding/base64
  - errors
  - fmt
  - io
```

第一段 helper implementation 不得添加 JWT、JWK、OAuth、OIDC、provider SDKs、password-hashing dependencies、Redis-like stores、KMS SDKs、cloud secret-manager SDKs、operations libraries 或 external randomness services。

## 9. Nakama And Pitaya Mapping

Nakama capability reference：

- Server-side account authentication 通常需要 server-issued secret material 和 session token issuance。
- vibit 采纳这种 capability need，而不是照搬其 exact implementation shape。

Pitaya capability reference：

- Realtime route handlers 应在 framework/application validation 之后接收 identity context。
- vibit 把 secret generation 保持在 application-owned helpers 中，而不是放在 transport acceptors 或 route dispatch 中。

本 gate 把这些参考映射成一段 narrow helper slice：先做 raw material creation，digest helpers 和 authentication behavior 后续再做。

## 10. Non-Goals

本 gate 不会：

- 添加 token generation code。
- 添加 credential generation code。
- Compute lookup digests。
- Compute verifier digests。
- Compare verifiers。
- Implement authentication service behavior。
- Execute login。
- Validate access tokens。
- Execute logout。
- Add refresh behavior。
- Add cleanup jobs。
- Add Protobuf authentication messages。
- Add WebSocket proof carriers。
- Change WebSocket handshake behavior。
- Change `authentication.Repository`。
- Change PostgreSQL migration schemas。
- Add startup wiring。
- Add authentication dependencies。
- Add production authentication behavior。

## 11. Verification Path

该 gate 的 repository check rule 是：

```text
runtime.token_credential_material_generation_implementation_gate
```

触及该 gate 的 change 需要运行：

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

本 gate 不验证 runtime Go token 或 credential material generation behavior，因为本 gate 不添加 generation behavior。

## 12. Follow-Up Gates

建议的后续 gates：

- Implement token and credential material generation helpers。
- Implement verifier digest computation helpers。
- Implement verifier digest comparison helpers。
- Implement application authentication service behavior。
- Add Protobuf authentication messages。
- Add WebSocket request proof carriers。
