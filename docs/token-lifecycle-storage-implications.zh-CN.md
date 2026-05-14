# Token 生命周期与存储影响

状态：草案 v0.1
最后更新：2026-05-14
范围：首个 token 姿态的颁发、过期、刷新、吊销、轮换、重放、登出、清理、脱敏、审计和存储影响
依赖：`docs/first-token-format-proof-carrier-posture.md`
权威决策：`ADR-0027`

英文文件 `docs/token-lifecycle-storage-implications.md` 是权威版本。本文是面向中文读者的人类可读翻译。

## 1. 目的

本文定义 vibit 首个已裁定 token 姿态的生命周期与存储影响：

```yaml
first_access_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
```

它把 W-0068 的 token 姿态转化为未来实现前必须满足的 gate。

本文不实现 token 生成、解析、验证、刷新、吊销、轮换、重放处理、存储、清理任务、审计事件、migration、Protobuf 变更、WebSocket 握手认证、runtime player handler 或 WebSocket route。

## 2. 生命周期概要

首个生命周期姿态为：

```yaml
token_kind: access_token
format: opaque_high_entropy_token
minimum_entropy_bits: 256
token_text_encoding: url_safe_unpadded_base64_or_equivalent
token_ttl: 1h
refresh_token: not_in_first_implementation
renewal_method: reauthenticate_with_selected_login_method
revocation_required: true
rotation_required: true
replay_control_required: true
logout_required: true
cleanup_required: true
audit_required: true
raw_token_storage: forbidden
verifier_storage_required: true
token_storage_default_target: postgresql_schema_gate
session_storage_required_for_first_posture: false
external_identity_storage_required_for_first_posture: false
credential_storage_required_for_device_credential_login: true
implementation_authorized: false
```

`1h` access-token TTL 是首个实现队列的生产取向默认值。它足够短，可以限制被盗 token 的有效时间；也足够长，不会让早期玩法循环频繁重新登录。后续运维决策可以通过配置调整该值，但默认值必须保持有限。

## 3. 颁发

未来成功的 `device_credential_login` command 只能在以下条件都成立后颁发 access token：

- credential proof 有效。
- player account 创建或查找策略成功。
- account lifecycle state 允许登录。
- token verifier storage boundary 已存在。
- 脱敏与审计规则已存在。
- login command contract 与 response contract 已存在。

颁发影响：

```yaml
issuer_owner: future_application_owned_authentication_boundary
issuer_layer: application_after_protocol_decode
transport_issuer: forbidden
protobuf_adapter_issuer: forbidden
domain_module_issuer: forbidden
player_account_repository_issuer: forbidden
generated_file_issuer: forbidden
```

token response 是一次性的客户端可见 secret presentation。服务端日志、trace、错误、测试快照、conversation log 和 change spec 都不能保存原始 token。

## 4. Token 形态

原始 access token 必须由加密安全随机数生成，并至少具备 256 bit 熵。

首个可接受的文本编码是 URL-safe unpadded Base64，或等价的编码；该编码必须避免控制字符、空白字符、路径分隔符、查询分隔符和视觉上容易混淆的格式。

规则：

- token 值区分大小写。
- token 值是 bearer secret。
- token 值不能包含客户端可读取 claims。
- token 值不能嵌入 `player_id`、`session_id`、provider subject、route name、timestamp、permission 或 account lifecycle state。
- 首个姿态不能从 URL query parameter 接收 token 值。
- token 值不能复制到 Protobuf `Session` metadata 字段。

## 5. 过期

首个 access-token TTL 为：

```yaml
access_token_ttl: 1h
```

规则：

- 未来 verifier schema 必须有 `issued_at` 与 `expires_at` 语义。
- 过期由未来 application-owned token validator 判断。
- 过期 token 产生独立的 expired-proof failure class。
- 过期 token 不能被静默当作 missing proof。
- 过期 token record 可以为了审计、重放检测或滥用分析临时保留。

过期 token record 的具体保留期仍然是 W-0071 的 schema gate。首个生命周期建议为：

```yaml
expired_token_retention_recommendation: 7d
```

这只是建议，不是 migration。

## 6. 刷新与续期

refresh token 不属于首个实现姿态。

首个续期方法是：

```yaml
renewal_method: reauthenticate_with_device_credential_login
```

也就是说，客户端通过再次执行所选登录方法来获得新的 access token。

规则：

- W-0069 不添加 refresh-token contract。
- W-0069 不添加 refresh-token storage。
- 不要把首个 access token 称为 session token。
- 不要把当前 `Session.session_id` 用作 refresh 或 renewal proof。
- 未来 refresh token 必须有自己的轮换、吊销、重放、存储、清理、脱敏、错误、权限和测试 gate。

## 7. 吊销与登出

首个 opaque-token 实现必须支持吊销。

未来所需最小状态：

```yaml
token_statuses:
  - active
  - expired
  - revoked
```

未来 schema 工作可以添加更多状态，但这三个是最小集合。

登出语义：

```yaml
logout_scope_first_posture: presented_access_token
logout_all_sessions: deferred
admin_revocation: deferred_to_permission_surface
forced_account_revocation: deferred_to_account_policy_and_audit_surface
```

首个 logout 行为应该吊销当前呈递的 access token。它不能吊销某个 player、credential、device 或 account 的所有 token，除非后续 contract 与 permission 决策明确授权。

吊销影响：

- 已吊销 token 的验证失败必须与 malformed 或 expired proof 区分。
- 吊销必须在生产敏感的 domain dispatch 之前生效。
- 当 audit storage 被裁定后，吊销必须能被审计工具看见。
- 对生产 opaque-token 姿态而言，仅运行时内存吊销是不够的，除非整个实现被明确标为 local-only。

## 8. 轮换

新颁发必须支持轮换。

首个姿态：

```yaml
rotation_on_successful_login: required
previous_token_for_same_credential: revoke_when_schema_supports_credential_token_linkage
automatic_background_rotation: deferred
refresh_rotation: deferred
```

首个实现应在成功登录时轮换 access token。一旦 schema gate 定义 credential proof 与 token verifier record 的关系，成功登录应吊销同一 credential installation 之前的 active access token，除非后续决策明确支持同一 credential 的多个并发 active token。

这不需要 session persistence。

## 9. 重放控制

opaque access token 是 bearer secret。token 格式本身不能消除重放。

首个必需控制：

- 高熵 token 生成。
- 有限 TTL。
- 非明文 verifier 存储。
- 吊销与登出。
- 成功登录时轮换。
- token 脱敏。
- token carrier 不得进入 route name、request ID、target ID、日志、URL query parameter 或 Protobuf `Session` metadata。
- 后续测试要覆盖重放敏感 failure class 和被盗 token 在既定模型内的行为。

延后控制：

- 每请求 nonce。
- token 绑定到 WebSocket connection。
- token 绑定到设备指纹、IP 地址、TLS session 或 first system message。
- 分布式 replay cache。
- Redis-like token/session store。

这些延后控制需要未来架构决策，因为它们影响协议形态、状态模型、分布式运行时和运维依赖。

## 10. 清理

启用生产 token storage 之前必须有清理策略。

首个清理姿态：

```yaml
cleanup_required: true
cleanup_owner: future_authentication_or_token_storage_boundary
cleanup_target: expired_and_revoked_token_verifier_records
cleanup_trigger_first_posture: explicit_maintenance_command_or_scheduled_runtime_job_deferred
default_retention_recommendation: 7d
```

本文不添加 cleanup job。

未来清理工作必须定义：

- 清理通过 CLI command、scheduled process、admin operation 还是 maintenance worker 执行。
- 清理是否可并发安全运行。
- 清理如何被审计。
- 清理是否进入默认 verification。
- 本地开发如何避免破坏性意外。

## 11. 脱敏

原始 token 值是 secret material。

必需脱敏规则：

```yaml
redact_raw_tokens_in:
  - logs
  - errors
  - traces
  - metrics_labels
  - test_snapshots
  - migration_fixtures
  - conversation_logs
  - change_specs
  - documentation_examples
  - panic_or_recovery_output
```

允许引用：

- 不由原始 token 文本派生的稳定 token record identifier。
- 若未来标准定义 fingerprint 算法，可以使用短 redacted fingerprint。
- 仅在受控数据库测试中，并且没有原始 token 且 fixture policy 明确时，才可使用 hash/verifier 值。

禁止引用：

- 原始 token 值。
- 长到对暴力破解有帮助的 token prefix。
- 嵌入 URL 的 token。
- 被复制到 `player_id`、`session_id`、`connection_id`、`request_id`、`target_id`、route name 或 error message 的 token。

## 12. 审计影响

审计是未来必需能力，但本文不添加 audit event。

未来 authentication audit surface 应覆盖：

- Token issued。
- Token validation failed。
- Token expired。
- Token revoked。
- Token logout requested。
- Token rotated。
- Token cleanup executed。
- Credential mismatch 或 account state 阻止颁发。

audit event 不能包含原始 token 值。

W-0070 必须在运行时实现前定义 public contract、error、permission 和 audit surface。W-0071 必须在任何 audit persistence 存在前定义 storage gate。

## 13. 存储影响

### Credential Storage

实现 `device_credential_login` 之前必须有 credential storage。

状态：

```yaml
credential_storage_required: true
credential_storage_added_now: false
credential_storage_schema_gate: W-0071
```

credential storage 仍然与 player account lifecycle storage 分离。

### Token Verifier Storage

实现 opaque access-token validation 之前必须有 token verifier storage。

状态：

```yaml
token_verifier_storage_required: true
token_verifier_storage_added_now: false
token_verifier_schema_gate: W-0071
default_store_target: PostgreSQL
redis_like_store_selected: false
```

首个 durable storage target 应该是 PostgreSQL，因为 PostgreSQL 已经是已裁定的权威持久存储。Redis-like store 仍然延后，直到依赖采纳和分布式运行时需求证明它有必要。

未来 token verifier schema 必须支持：

- 非明文 token verifier。
- token status。
- subject actor。
- audience。
- issued-at timestamp。
- expires-at timestamp。
- revoked-at timestamp，如果已吊销。
- rotation lineage 或 replacement relationship，如果被裁定。
- 如果强制同一 credential 只有一个 active token，需要 credential-token linkage。
- audit-safe record identifier。

现在不添加 table、migration、repository 或 adapter。

### External Identity Storage

首个 `device_credential_login` 姿态不需要 external identity storage。

状态：

```yaml
external_identity_storage_required_for_first_posture: false
external_identity_storage_added_now: false
```

provider login 与 identity linking 仍然延后。

### Session Storage

首个姿态不需要 session storage。

状态：

```yaml
session_storage_required_for_first_posture: false
session_storage_added_now: false
session_token_vocabulary: deferred_until_session_persistence
websocket_connection_binding: deferred
```

首个姿态可以在每个 authenticated request 上验证 access-token proof，而不需要持久化 runtime session。后续 session persistence milestone 可以选择 session store、connection binding、first-message authentication、handshake authentication 或 hybrid model。

### Player Account Lifecycle Storage

player account lifecycle storage 必须继续保持 credential-free、token-free、external-identity-free 和 session-free。

当前 lifecycle tables 仍然是：

```text
player_accounts
player_account_events
```

本生命周期标准禁止：

- 向 `player_accounts` 添加 credential column。
- 向 `player_accounts` 添加 token column。
- 向 `player_accounts` 添加 provider subject column。
- 向 `player_accounts` 添加 session 或 WebSocket state column。
- 向 `player_account_events` 添加 raw token、token verifier、credential、provider subject、session 或 WebSocket state row。

## 14. 参考对齐

### Nakama

Nakama 仍然是 access/session token 颁发、refresh、expiration、logout 和 revocation 词汇的能力参考。

vibit 采用这些生命周期维度，但继续延后 refresh token 和 session token 词汇。首个 vibit 姿态是通过所选登录方法续期的 storage-backed opaque access token，不直接复制 Nakama 的 token/session 行为。

### Pitaya

Pitaya 仍然是 session context 与 connection/session binding 的词汇参考。

vibit 延后 connection-bound session 行为。首个生命周期姿态保持 token validation 由 application 拥有，不把 token state 放进 WebSocket acceptor 或 route handler。

## 15. 未来所需 Gate

实现前，未来工作必须提供：

- W-0070 的 contract、error、permission 和 audit surface。
- W-0071 的 credential、token verifier 和可选 session schema gate。
- W-0072 的禁止 shortcut 的仓库检查。
- token verifier repository boundary。
- credential lookup boundary。
- 脱敏测试。
- 过期测试。
- 吊销测试。
- 登出测试。
- 轮换测试。
- 重放敏感测试。
- 清理测试，或明确的清理延后与生产分类。

## 16. 非授权事项

本文不授权：

- token 生成代码。
- token 解析或验证代码。
- token storage table。
- credential storage table。
- external identity table。
- session table。
- migration。
- cleanup job。
- audit event 实现。
- refresh token。
- JWT、signing、key-management、OAuth、OIDC、provider SDK、Redis-like 或 password-hashing 依赖。
- Protobuf envelope 变更。
- WebSocket handshake authentication。
- first system-message binding。
- runtime player handler。
- WebSocket route。

## 17. 后续

下一项工作：

```text
W-0070 Define authentication contract error permission surfaces
```

W-0070 必须在实现前定义所选 login 与 token 姿态所需的语义 contract、error、permission 和 audit surface。
