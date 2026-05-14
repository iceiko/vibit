# 首个 Token 格式与证明承载姿态

状态：草案 v0.1
最后更新：2026-05-14
范围：M-013 的首个 token 格式与证明承载姿态裁定
依赖：`docs/token-format-carrier-options.md`
权威决策：`ADR-0026`

英文文件 `docs/first-token-format-proof-carrier-posture.md` 是权威版本。本文是面向中文读者的人类可读翻译。

## 1. 目的

本文在 `device_credential_login` 被选为首个登录方法集合之后，裁定 vibit 的首个 token 格式与证明承载姿态。

它选择首个访问 token 格式、颁发承载方式、请求证明承载方式，以及未被授权的边界。

本文不实现 token 生成、解析、签名、验证、刷新、吊销、轮换、重放处理、存储、会话持久化、Protobuf envelope 变更、WebSocket 握手认证、运行时 player handler 或 WebSocket route。

## 2. 裁定姿态

首个 token 姿态为：

```yaml
first_access_token_format: opaque_high_entropy_token
token_issuance_carrier: login_command_response_token
request_proof_carrier: explicit_request_proof_payload
refresh_token: deferred
session_token_vocabulary: deferred_until_session_persistence
protobuf_envelope_change: false
websocket_handshake_authentication_change: false
current_session_metadata_as_proof: false
first_system_message_binding: deferred
implementation_authorized: false
```

被选中的 token 是访问 token。它不是 refresh token、持久化运行时 session 标识符、WebSocket 连接标识符、Protobuf `Session.session_id`、player ID 或外部 provider token。

## 3. Access Token 格式

### `opaque_high_entropy_token`

定义：

```text
客户端呈递的高熵 bearer secret，不包含客户端可读取的 claims。
```

裁定姿态：

```yaml
format: opaque_high_entropy_token
token_kind: access_token
issuer: future_application_owned_authentication_boundary
verifier: future_application_owned_token_validator_before_domain_dispatch
subject: player_account_id_after_credential_and_account_policy_success
audience: vibit_gameplay_runtime_requests
expiration: required_finite_expiration_exact_policy_deferred_to_W_0069
refresh: deferred
revocation: required_capability_policy_deferred_to_W_0069
rotation: required_for_new_issuance_policy_deferred_to_W_0069
replay_posture: bearer_token_risk_must_be_controlled_by_lifecycle_policy
redaction: raw_token_secret_redacted_everywhere_except_client_presentation
storage: lookup_safe_hash_or_equivalent_non_plaintext_verifier_required_before_implementation
requires_signing_dependency: false
requires_key_management: false
requires_protobuf_envelope_change: false
requires_websocket_handshake_authentication: false
confidence: high
```

该 token 不包含客户端可检查的 claims。授权事实不能隐藏在 token 内容中。模块自己的权限逻辑仍然必须与 token 验证分离。

## 4. 颁发者与验证者

未来的颁发者是 application-owned authentication boundary：它位于协议解码之后，并且只在 `device_credential_login` 凭证验证成功之后颁发 token。

未来的验证者是 application-owned token validator：它在生产敏感的 domain dispatch 之前运行。

以下层不能颁发或验证 token：

- WebSocket transport adapter。
- Protobuf envelope adapter。
- inventory 等 domain module。
- player account persistence repository。
- 生成的 Protobuf 文件或生成的 contract shape 文件。

这保持了 vibit 已采用的运行时边界：transport 和 protocol 负责字节与 envelope，application dispatch 拥有验证交接，domain module 接收已经规范化的 request identity。

## 5. Subject 与 Audience

token subject 只能是在所选登录方法与账号创建或查找策略成功之后得到的 player account identifier。

初始 audience 为：

```text
vibit gameplay runtime requests
```

首个 token 姿态不授权：

- 服务间权限。
- 管理员权限。
- 外部 provider session。
- 跨游戏或跨项目 audience 共享。
- 长期离线凭证。

这些都需要后续合同与决策。

## 6. 过期、刷新、吊销、轮换与重放

首个 access token 必须有有限过期时间。具体过期时长延后到 W-0069。

refresh token 延后。未来若加入 refresh token，必须先定义轮换、吊销、重放、存储、清理、脱敏和错误语义。

由于 opaque token 是 storage-backed 或 verifier-backed，吊销能力是必需的。具体吊销模型延后到 W-0069。

新颁发与替换策略需要轮换规则，但具体触发条件延后到 W-0069。

opaque access token 是 bearer secret。除非后续裁定绑定策略，被盗 token 在过期或吊销前可能被重放。W-0069 必须在实现前定义重放控制。

本裁定不把 token 绑定到 WebSocket 连接、first system message、设备指纹、IP 地址或当前 Protobuf session metadata。

## 7. 脱敏与存储

原始 token 值是 secret material。

规则：

- 原始 token 值只能由客户端呈递。
- 原始 token 值不能明文存储。
- 如果服务器端需要存储，只能存储 lookup-safe hash 或等价的非明文 verifier。
- 日志、错误、trace、conversation log、change spec、测试和文档必须脱敏 token 值。
- token 值不能出现在 route name、request ID、target ID、player ID、session ID、connection ID 或 migration fixture 中。

W-0069 和 W-0071 必须在存储存在前定义具体生命周期与 schema gate。

## 8. 证明承载

### 颁发承载

首个 token 应由未来成功的 login command response 颁发：

```yaml
token_issuance_carrier: login_command_response_token
```

这意味着未来语义化的 `device_credential_login` command 可以在凭证验证和账号策略成功后返回 access token。

这还没有授权 command contract。W-0070 必须在运行时实现前定义语义 command、response、error、permission 和 audit surface。

### 请求证明承载

首个 authenticated route 的请求证明承载为：

```yaml
request_proof_carrier: explicit_request_proof_payload
```

authenticated command 或 query 可以在显式的 contract-owned payload 字段中携带 access-token proof，直到后续协议决策裁定更干净的 carrier。

这有意保持冗长。好处是后续 Agent 能在语义合同中直接看到 proof 要求，而不需要改变当前 Protobuf envelope 或 WebSocket handshake。

## 9. 不变的协议与传输行为

本裁定不改变 Protobuf envelope。

当前 Protobuf session metadata 仍然只是 metadata：

- `Session.session_id`
- `Session.player_id`
- `Session.connection_id`
- `Session.connection_epoch`

这些字段不是 authenticated proof，不能通过重新解释来满足生产权限。

本裁定不改变 WebSocket 握手认证。WebSocket transport 仍然保持 credential-neutral。

本裁定不加入 first system-message binding。连接绑定、重连行为和 session persistence 都仍然是未来决策。

## 10. 拒绝与延后选项

| 选项 | 状态 | 原因 |
| --- | --- | --- |
| signed structured token | 延后 | 需要签名依赖、密钥管理、claim ownership、密钥轮换、吊销、重放和时钟偏差决策。 |
| external provider token 作为 vibit access token | 延后 | 属于 provider login 和 external identity linking，不属于首个 device credential login。 |
| plain session ID as secret | 首个姿态拒绝 | 混淆 identifier 与 proof，容易重新解释 metadata-only session 字段。 |
| 当前 Protobuf `Session` metadata 作为 proof | 拒绝 | 现有字段只是 metadata，不能因为方便而变成 authority。 |
| Protobuf envelope extension | 延后 | 这是兼容性敏感的 wire decision。 |
| WebSocket handshake carrier | 延后 | 容易把认证放入 transport，并且需要浏览器/非浏览器 carrier 分析。 |
| first system-message binding | 延后 | 需要 system-message contract、超时行为、重连规则和连接状态。 |
| refresh token | 延后 | 需要生命周期、轮换、吊销、重放、脱敏、存储和错误语义。 |

## 11. 参考对齐

### Nakama

Nakama 仍然是 session token、refresh token、expiration 和 logout 词汇的能力参考。

vibit 采用“认证交换成功后可以颁发客户端呈递 token”的思路，但不复制 Nakama 的 token 格式、公开 API、refresh 行为或 realtime socket binding 语义。

### Pitaya

Pitaya 仍然是 session binding 与 handler context 的词汇参考。

vibit 通过要求未来 application-owned validator 在 domain dispatch 之前生成 request identity 来吸收 session-context 思路。它不把 token 验证放进 WebSocket acceptor 或 route handler。

## 12. 实现前所需 Gate

实现该姿态之前，仓库必须具备：

- W-0069 定义 token 生命周期与存储影响。
- W-0070 定义认证 contract、error、permission 与 audit surface。
- W-0071 定义 credential、token、session schema gate。
- W-0072 增加所选 login/token 边界的仓库检查。
- 语义化 login command 与 response contract。
- 语义化 authenticated-request proof contract shape。
- credential 与 token 的脱敏规则。
- 如果验证需要持久查找，必须有 token verifier storage 或等价 lookup boundary。
- 针对成功、缺失 proof、格式错误 proof、无效 proof、过期 proof、已吊销 proof、重放/碰撞行为、脱敏和层所有权的聚焦测试。

## 13. 非授权事项

本裁定不授权：

- 运行时认证代码。
- login handler。
- token 生成、解析、签名、验证、刷新、吊销、轮换、重放处理或存储。
- credential table。
- external identity table。
- token table。
- session table。
- migration。
- password hashing、JWT、OAuth、OIDC、provider SDK、Redis-like、cryptography、key-management 或重大认证依赖。
- Protobuf envelope 变更。
- WebSocket handshake authentication。
- first system-message authentication。
- runtime player account handler。
- WebSocket route。
- 把 metadata-only 的 `player_id`、`session_id`、`connection_id` 或 `connection_epoch` 当作 proof。

## 14. 后续

下一项工作：

```text
W-0069 Define token lifecycle and storage implications
```

W-0069 必须在实现前定义颁发、过期、吊销、轮换、重放、清理、登出、存储、审计和脱敏影响。
