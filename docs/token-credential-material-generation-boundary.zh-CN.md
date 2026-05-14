# Token And Credential Material Generation Boundary

状态：Draft v0.1
最后更新：2026-05-15
范围：High-entropy device credential 与 opaque access-token material 的 future generation ownership、entropy enforcement、text encoding、one-time presentation、non-storage、repository handoff、redaction、dependency posture 和 test expectations
依赖：`docs/token-credential-verifier-algorithm-redaction-boundary.md`, `docs/secret-configuration-verifier-key-loading-boundary.md`, `docs/token-lifecycle-storage-implications.md`, `docs/first-login-method-set.md`
Canonical decision: `ADR-0042`

配套英文源文档是 `docs/token-credential-material-generation-boundary.md`。英文文件是权威版本。

## 1. Purpose

本标准定义 future application-owned authentication code 如何生成 raw credential material 和 raw access-token material。

它位于 token generation code、credential generation code、secret loading、verifier digest computation、verifier comparison、login execution、token validation、logout execution、cleanup jobs、Protobuf messages、WebSocket proof carriers、authentication dependencies、repository interface changes、migration schema changes 或 production authentication behavior 之前。

这只是 boundary-only standard。它不添加 Go code、imports、runtime services、token generation、credential generation、verifier digest computation、verifier comparison、secret loading、repository methods、SQL migrations、Protobuf messages、WebSocket carriers、routes 或 production authentication behavior。

## 2. Core Rule

第一版 material generation posture 是：

```yaml
material_generation_boundary: defined
implementation_authorized_by_this_standard: false
planned_owner: runtime/internal/app
planned_optional_child_package: runtime/internal/app/authentication
first_device_credential_source: server_issued_application_generated
first_access_token_source: server_issued_application_generated
minimum_raw_device_credential_entropy_bits: 256
minimum_raw_access_token_entropy_bits: 256
raw_material_size_bytes: 32
text_encoding: url_safe_unpadded_base64_or_equivalent
one_time_client_visible_presentation_required: true
raw_credential_storage: forbidden
raw_token_storage: forbidden
raw_material_in_repository: forbidden
raw_material_in_transport_logs: forbidden
external_randomness_dependency_required_for_first_posture: false
future_allowed_standard_library_packages:
  - crypto/rand
  - encoding/base64
```

后续 code gate 授权 implementation 后，future first-posture generation helpers 可以使用 Go standard library `crypto/rand` 和 `encoding/base64`。该 generation posture 不需要 external randomness、cryptography、JWT、OAuth、OIDC、provider、KMS、cloud secret-manager 或 operations dependency。

## 3. Ownership

Future material generation 由 application 拥有：

```text
runtime/internal/app
```

Future helper code 可以放在 application-owned child package，例如：

```text
runtime/internal/app/authentication
```

Ownership rules：

- Application authentication code 在后续 implementation gate 授权后拥有 raw token 和 credential material generation。
- `authentication.Repository` 绝不生成、为存储接收或返回 raw credential 或 raw token material。
- PostgreSQL adapters 绝不生成、encode、log、persist 或返回 raw credential 或 raw token material。
- Protobuf adapters 和 WebSocket transports 可在后续 protocol gates 后承载 already-decoded proof fields，但它们不生成 secrets。
- Generated authentication contract shapes 保持 metadata-only 且 immutable。

Material generation 不得放入 transport handlers、protocol adapters、domain modules、repositories、generated output、migrations、SQL fixtures，或未明确限定为 future generation helpers 的 tests。

## 4. Device Credential Material

第一版 device credential posture 是 server-issued 且 application-generated。

```yaml
device_credential_material:
  source: server_issued_application_generated
  credential_kind: high_entropy_installation_credential
  minimum_entropy_bits: 256
  raw_size_bytes: 32
  text_encoding: url_safe_unpadded_base64_or_equivalent
  one_time_client_visible_presentation: required
  raw_server_storage: forbidden
  public_metadata_source: forbidden
```

规则：

- Credential 是 secret proof material，不是 raw operating-system device ID、advertising ID、hardware serial number、account email、player name、provider subject、player id、session id、connection id 或 metadata-only value。
- Server 只能通过 future explicitly authorized authentication response 或 credential bootstrap response 向 client 展示 credential。
- Client 负责存储该 future response 中 one-time presented credential。
- Future server code 必须先计算 credential lookup 和 verifier digests，再进行 repository storage，但 exact digest helper implementation 仍属于后续 gate。
- Future code 只能存储 digest material、metadata、lifecycle state、timestamps 和 `verifier_key_id`。
- 第一版姿态不接受 client-generated installation credentials。后续 boundary 可以添加 client-generated credential enrollment，但必须定义 entropy、replay、collision、proof-carrier 和 abuse controls。

本标准不定义 account creation policy、credential rotation command、recovery flow 或 Protobuf response shape。

## 5. Access Token Material

第一版 access-token posture 是 server-issued 且 application-generated。

```yaml
access_token_material:
  source: server_issued_application_generated
  token_format: opaque_high_entropy_token
  minimum_entropy_bits: 256
  raw_size_bytes: 32
  text_encoding: url_safe_unpadded_base64_or_equivalent
  one_time_client_visible_presentation: required
  raw_server_storage: forbidden
  claim_container: false
```

规则：

- Access tokens 是 bearer secrets。
- Token text 不得嵌入 player id、credential id、token record id、session id、route name、timestamp、permission、account state、provider subject 或 claims。
- Future server code 必须先计算 token lookup 和 verifier digests，再进行 repository storage，但 exact digest helper implementation 仍属于后续 gate。
- Future code 只能存储 digest material、token metadata、lifecycle state、timestamps 和 `verifier_key_id`。
- Token text 只能在 future explicitly authorized authentication response 中向 client 返回一次。
- Token text 不得被 URL query parameters、route names、当前 Protobuf `Session` metadata fields、logs、traces、metrics labels、public errors、audit-safe facts、conversation logs 或 change specs 接收或保存。

本标准不定义 login execution、token issuance command behavior、token validation、logout、cleanup、refresh、Protobuf messages 或 WebSocket proof carriers。

## 6. Encoding And Byte Shape

Future first-posture generation 必须先产生 canonical raw bytes，再进行 text encoding。

```yaml
raw_material:
  source: cryptographically_secure_random_bytes
  bytes: 32
  entropy_bits: at_least_256
  all_zero_value: forbidden
  repeated_pattern_value: forbidden
  human_readable_password_like_value: forbidden
text_material:
  encoding: url_safe_unpadded_base64_or_equivalent
  whitespace: forbidden
  control_characters: forbidden
  path_separators: forbidden
  query_delimiters: forbidden
  case_sensitive: true
```

规则：

- Encoding 只是 presentation。它不得成为 claim container。
- Decoding 必须能无损恢复原始 generated bytes。
- Future tests 应覆盖 entropy length、text alphabet、如果选择 base64url unpadded 则无 padding，以及 round-trip decoding。
- Future code 如果无法取得 cryptographically secure random bytes，必须 fail closed。

## 7. Repository Handoff

Future application authentication code 必须保持 repository handoff non-raw。

后续 gates 授权后允许的 repository inputs：

- `credential_lookup_digest`
- `credential_verifier_digest`
- `token_lookup_digest`
- `token_verifier_digest`
- `verifier_algorithm`
- `verifier_version`
- `verifier_key_id`
- lifecycle state
- timestamps
- stable record identifiers

禁止的 repository inputs：

- Raw credential text。
- Raw credential bytes。
- Raw access-token text。
- Raw access-token bytes。
- Token prefixes。
- Credential prefixes。
- Randomness seeds。
- Verifier keys。
- Encoded verifier keys。
- Secret-manager payloads。

Repository 可以持久化 already-computed digest bytes。它不得 generate raw material、encode raw material、compute verifier digests、compare verifiers 或决定 authentication outcomes。

## 8. Redaction Requirements

以下内容禁止出现在 logs、traces、metrics labels、public errors、panic output、audit-safe facts、test snapshots、fixtures、ADRs、change specs、documentation examples 和 conversation logs 中：

- Raw credential text。
- Raw credential bytes。
- Raw access-token text。
- Raw access-token bytes。
- Token 或 credential prefixes。
- Randomness seeds。
- Encoding 前的 generated material。
- Encoded generated material。
- Lookup digests。
- Verifier digests。
- Verifier keys。
- Encoded verifier keys。
- Full concrete `verifier_key_id` values。

Allowed with care：

- Non-secret record identifiers。
- Lifecycle state names。
- Registered error codes。
- `<raw-access-token>` 这样的 placeholder names，但只用于说明 redaction rules。

Generated material 不能因为 short-lived 就被视为安全。One-time presentation 意味着一次 client-visible delivery，不是“可以 log 一次”。

## 9. Test Expectations

未来只要 implementation work 生成 raw credential 或 access-token material，就必须添加 focused tests。

Minimum expectations：

```yaml
generation_tests:
  random_source_error_fails_closed: required
  raw_material_length_is_32_bytes: required
  encoded_material_round_trips_to_raw_bytes: required
  encoded_material_uses_allowed_text_alphabet: required
  generated_values_are_not_constant: required
  raw_material_absent_from_logs: required_when_logging_path_exists
  raw_material_absent_from_public_errors: required
  raw_material_absent_from_repository_records: required
  token_text_contains_no_claims: required
  credential_text_contains_no_metadata: required
```

Tests 可以使用 deterministic fakes，但只能通过 future implementation code 中的 test seam。它们不得替代 production randomness、泄露 fixture secrets，或让 generated values 在 non-test builds 中稳定。

## 10. Dependency Posture

第一版 material generation posture 不需要 external dependency。

后续 explicit code gate 授权后可接受：

```yaml
go_standard_library:
  randomness: crypto/rand
  base64url: encoding/base64
external_dependency_adoption_record_required_for_first_posture: false
```

第一版 posture deferred 且不需要：

- External randomness services。
- KMS 或 cloud secret-manager SDKs。
- JWT、JWK 或 signing libraries。
- OAuth 或 OIDC provider SDKs。
- Password-hashing dependencies。
- Redis-like token/session stores。

如果未来 login method 接受 password-like 或 low-entropy material，必须定义单独的 credential boundary 和 dependency adoption record。

## 11. Non-Goals

本标准不：

- 添加 token generation。
- 添加 credential generation。
- 添加 secret loading。
- 添加 verifier digest computation。
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

## 12. Verification Path

该 boundary 的 repository check rule 是：

```text
runtime.token_credential_material_generation_boundary
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

本标准不添加 runtime Go token 或 credential generation behavior，因此不验证该 behavior。

## 13. Follow-Up Gates

Recommended follow-up gates：

- Verifier digest computation and constant-time comparison boundary。
- Application authentication service implementation gate。
- Authentication redaction test implementation gate。
- Protobuf authentication message gate。
- WebSocket request proof carrier gate。
