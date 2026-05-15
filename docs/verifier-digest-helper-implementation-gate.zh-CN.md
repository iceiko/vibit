# 验证摘要辅助函数实现门限

状态：草案 v0.1
最后更新：2026-05-15
范围：验证摘要辅助函数实现门限、未来辅助函数文件边界、规范化输入构造、验证密钥交接、摘要类映射、脱敏规则、测试以及在添加摘要代码之前的延迟项
依赖：`docs/verifier-digest-computation-comparison-boundary.md`、`docs/token-credential-verifier-algorithm-redaction-boundary.md`、`docs/token-credential-material-generation-implementation-gate.md`、`docs/authentication-service-implementation-readiness-gate.md`
规范化决策：`ADR-0048`

配对的英文版本为 `docs/verifier-digest-helper-implementation-gate.md`。英文文件为权威来源。

## 1. 目的

此门限定义了查找摘要和验证摘要计算辅助函数的下一个有界实现切片。

仓库已定义验证摘要计算与比较边界、验证算法与脱敏边界以及材料生成辅助函数。下一个风险是未来代理在服务编排、仓储适配器、传输处理器、测试夹具或协议代码中计算摘要，或将摘要计算与验证比较和认证行为混合。

这是一个实现门限标准。它不会添加 Go 代码、导入、HMAC 计算、摘要辅助函数、验证比较、认证服务行为、登录执行、令牌验证、注销执行、清理任务、Protobuf 消息、WebSocket 证明载体、仓储方法、SQL 迁移、启动接线、认证依赖、外部加密服务、KMS、云密钥管理集成或生产认证行为。

## 2. 核心规则

验证摘要辅助函数实现门限为：

```yaml
verifier_digest_helper_implementation_gate: defined
implementation_authorized_by_this_standard: false
future_implementation_work_item: W-0103
planned_owner: runtime/internal/app
planned_package: runtime/internal/app/authentication
future_slice: lookup_and_verifier_digest_computation_helpers
future_source: runtime/internal/app/authentication/verifier_digest.go
future_tests: runtime/internal/app/authentication/verifier_digest_test.go
verifier_algorithm_family: vibit_hmac_sha256_v1
verifier_version: 1
digest_size_bytes: 32
canonical_input_version: vibit.auth.verifier.input.v1
canonical_input_encoding: ascii_header_length_prefixed_purpose_label_length_prefixed_raw_material
hmac_hash: crypto/sha256
constant_time_comparison_primitive: crypto/hmac.Equal
verifier_key_handoff: VerifierKeySet
raw_material_handoff: GeneratedSecretMaterial_or_raw_bytes
digest_output_shape: raw_32_byte_digest
digest_copying_required: true
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

未来实现必须是仅辅助函数的切片。它可以构建规范化摘要输入、计算 HMAC-SHA-256 查找摘要、计算 HMAC-SHA-256 验证摘要，并返回已复制的摘要字节切片。它不得比较验证摘要、写入仓储、选择账户、发出登录响应、验证令牌、解析持有者证明或触及协议载体。

## 3. 未来辅助函数形态

未来实现归属：

```text
runtime/internal/app/authentication
```

在未来实现工作项授权代码后允许的文件：

```text
runtime/internal/app/authentication/verifier_digest.go
runtime/internal/app/authentication/verifier_digest_test.go
```

首选的未来 API 形态：

```yaml
future_types:
  DigestClass:
    values:
      - credential_lookup
      - credential_verifier
      - token_lookup
      - token_verifier
  ComputedDigest:
    owns: copied_digest_bytes_and_class
    methods:
      - Class() DigestClass
      - Bytes() []byte
    constraints:
      - Bytes 返回副本。
      - 错误文本和字符串格式化不得暴露摘要字节、原始材料或密钥值。

future_constants:
  CanonicalInputVersion: "vibit.auth.verifier.input.v1"
  PurposeLabelCredentialLookup: "vibit.credential.lookup.v1"
  PurposeLabelCredentialVerifier: "vibit.credential.verifier.v1"
  PurposeLabelTokenLookup: "vibit.access_token.lookup.v1"
  PurposeLabelTokenVerifier: "vibit.access_token.verifier.v1"

future_functions:
  ComputeCredentialLookupDigest:
    signature: "func ComputeCredentialLookupDigest(keySet VerifierKeySet, rawMaterial []byte) (ComputedDigest, error)"
    behavior: build_canonical_input_and_hmac_with_credential_lookup_key
  ComputeCredentialVerifierDigest:
    signature: "func ComputeCredentialVerifierDigest(keySet VerifierKeySet, rawMaterial []byte) (ComputedDigest, error)"
    behavior: build_canonical_input_and_hmac_with_credential_verifier_key
  ComputeTokenLookupDigest:
    signature: "func ComputeTokenLookupDigest(keySet VerifierKeySet, rawMaterial []byte) (ComputedDigest, error)"
    behavior: build_canonical_input_and_hmac_with_token_lookup_key
  ComputeTokenVerifierDigest:
    signature: "func ComputeTokenVerifierDigest(keySet VerifierKeySet, rawMaterial []byte) (ComputedDigest, error)"
    behavior: build_canonical_input_and_hmac_with_token_verifier_key
```

`VerifierKeySet` 交接是必需的。辅助函数必须接受已验证的密钥集；它不得加载密钥、解析环境变量或选择密钥集。原始材料交接是来自 `GeneratedSecretMaterial.RawBytes()` 或等效的已解码原始字节；它不得接受编码文本、玩家 ID、会话 ID、路由名称或元数据。

## 4. 规范化输入构造

未来辅助函数必须在 HMAC 计算之前构建确定性规范化字节输入。

```yaml
canonical_digest_input:
  version_header_ascii: "vibit.auth.verifier.input.v1"
  header_separator: 0x00
  purpose_label_length: uint16_big_endian_byte_length
  purpose_label: ascii_bytes
  raw_material_length: uint16_big_endian_byte_length
  raw_material: generated_secret_material_bytes
```

字节序列为：

```text
ascii("vibit.auth.verifier.input.v1")
|| 0x00
|| uint16be(len(purpose_label))
|| ascii(purpose_label)
|| uint16be(len(raw_material))
|| raw_material
```

规则：

- 版本头必须是字面 ASCII 字节。如果规范化输入形状变化则必须改变。
- 目的标签必须是每个摘要类的已注册常量。
- 原始材料必须是已解码的原始秘密材料字节，不是归一化文本、元数据、玩家 ID、会话 ID、路由名称或提供者主体。
- 第一姿态的原始材料长度为 32 字节。辅助函数必须拒绝长度为零的原始材料。
- 长度前缀使用大端无符号 16 位整数，使输入即使在未来原始材料形状变化时也无歧义。
- 未来测试必须包含确定性夹具向量，通过比较整个输入字节序列来验证规范化字节构造。

## 5. 摘要类与密钥映射

每个摘要类使用已注册的目的标签和验证密钥集中匹配的逻辑密钥：

```yaml
credential_lookup_digest:
  purpose_label: vibit.credential.lookup.v1
  key_accessor: VerifierKeySet.CredentialLookupKey
  output_bytes: 32
credential_verifier_digest:
  purpose_label: vibit.credential.verifier.v1
  key_accessor: VerifierKeySet.CredentialVerifierKey
  output_bytes: 32
token_lookup_digest:
  purpose_label: vibit.access_token.lookup.v1
  key_accessor: VerifierKeySet.TokenLookupKey
  output_bytes: 32
token_verifier_digest:
  purpose_label: vibit.access_token.verifier.v1
  key_accessor: VerifierKeySet.TokenVerifierKey
  output_bytes: 32
```

规则：

- 查找摘要和验证摘要必须使用不同的目的标签。
- 凭证摘要和令牌摘要必须使用不同的目的标签。
- 每个计算函数必须恰好使用匹配的密钥；通过函数签名不可能传递错误的密钥。
- 摘要字节不得在存储、查找、比较或日志中被截断。
- 辅助函数必须在返回时复制摘要字节。

## 6. 包和文件边界

允许的未来实现区域：

```text
runtime/internal/app/authentication/verifier_digest.go
runtime/internal/app/authentication/verifier_digest_test.go
```

未来验证摘要辅助函数切片的禁止写入区域：

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

未来实现不得将摘要计算接线到登录执行、账户创建、令牌发行、令牌验证、启动、WebSocket 传输、Protobuf 协议、PostgreSQL 持久化、生成的合约输出、迁移或域仓储代码中。

## 7. 错误与脱敏要求

未来辅助函数应暴露类型化或哨兵错误，对测试有用而不暴露摘要材料或密钥值。

允许的错误分类示例：

```yaml
error_classes:
  missing_key_set
  missing_raw_material
  invalid_digest_computation
```

错误文本中允许的内容：

- 错误类。
- 摘要类名称。
- 非秘密数值预期，如 `32` 字节。

在错误、日志、测试快照、文档、ADR、变更规范和对话日志中禁止的内容：

- 原始凭证材料或字节。
- 原始访问令牌材料或字节。
- 编码的生成材料。
- 查找摘要字节。
- 验证摘要字节。
- HMAC 输入字节。
- HMAC 输出字节。
- 验证密钥值。
- 编码的验证密钥值。
- 完整的具体 `verifier_key_id` 值。
- 候选密钥集计数。

摘要材料不安全仅因为它是确定性的。摘要字节仍然是秘密相邻的，必须保留在公共工件之外。

## 8. 未来辅助函数的必需测试

未来实现必须在以下位置添加聚焦的单元测试：

```text
runtime/internal/app/authentication/verifier_digest_test.go
```

最低测试用例：

```yaml
required_tests:
  canonical_input_is_deterministic
  canonical_input_uses_version_header
  canonical_input_null_separator_present
  canonical_input_length_prefixes_purpose_label
  canonical_input_length_prefixes_raw_material
  lookup_and_verifier_purpose_labels_differ
  credential_and_token_purpose_labels_differ
  digest_output_is_32_bytes
  credential_lookup_digest_uses_credential_lookup_key
  credential_verifier_digest_uses_credential_verifier_key
  token_lookup_digest_uses_token_lookup_key
  token_verifier_digest_uses_token_verifier_key
  different_keys_produce_different_digests
  different_raw_material_produces_different_digests
  digest_bytes_are_copied_on_return
  empty_raw_material_fails_closed
  errors_do_not_include_digest_or_key_material
  helper_does_not_compare_verifiers
```

测试可以仅在测试内部使用确定性合成密钥和原始材料。它们不得成为生产默认值、真实秘密形状的文档示例或已提交的类生产秘密值。

## 9. 依赖姿态

此门限不允许新的外部依赖。

在未来实现工作项授权代码后允许的 Go 标准库包：

```yaml
future_standard_library_imports_allowed:
  - crypto/hmac
  - crypto/sha256
  - encoding/binary
  - errors
  - fmt
```

第一个辅助函数实现不得添加 JWT、JWK、OAuth、OIDC、提供者 SDK、密码哈希依赖、类 Redis 存储、KMS SDK、云密钥管理 SDK、运营库或外部加密服务。

## 10. Nakama 和 Pitaya 映射

Nakama 能力参考：

- 服务端认证通常需要基于 HMAC 的凭证和令牌验证。vibit 采用能力需求，而非精确的实现形状。

Pitaya 能力参考：

- 实时路由处理器应在框架/应用验证后接收身份上下文。vibit 将摘要计算保留在应用拥有的辅助函数中，而非传输接受器或路由调度中。

此门限将这些参考映射为一个窄辅助函数切片：先摘要计算，后验证比较和认证行为。

## 11. 非目标

此门限不：

- 添加验证摘要计算代码。
- 添加验证比较代码。
- 添加令牌生成编排。
- 添加凭证生成编排。
- 实现认证服务行为。
- 执行登录。
- 验证访问令牌。
- 执行注销。
- 添加刷新行为。
- 添加清理任务。
- 添加 Protobuf 认证消息。
- 添加 WebSocket 证明载体。
- 更改 WebSocket 握手行为。
- 更改 `authentication.Repository`。
- 更改 PostgreSQL 迁移模式。
- 添加启动接线。
- 添加认证依赖。
- 添加生产认证行为。

## 12. 验证路径

此门限的仓库检查规则为：

```text
runtime.verifier_digest_helper_implementation_gate
```

对于涉及此门限的变更，运行：

```bash
node tools/vibit check runtime --json
node tools/vibit check contracts --json
node tools/vibit check module authentication --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

此门限不验证运行时 Go 验证摘要计算行为，因为未添加计算行为。

## 13. 后续门限

推荐的后续门限：

- 实现验证摘要计算辅助函数。
- 实现验证摘要比较辅助函数。
- 实现应用认证服务行为。
- 添加 Protobuf 认证消息。
- 添加 WebSocket 请求证明载体。
