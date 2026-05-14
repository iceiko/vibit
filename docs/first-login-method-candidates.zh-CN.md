# First Login Method Candidate Comparison 中文版

状态：Draft v0.1
最后更新：2026-05-14
范围：M-013 的第一批 login-method candidates comparison
依赖：`docs/login-method-token-format-ratification.md`

本文件是 `docs/first-login-method-candidates.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

## 1. 目的

本文在 repository ratify 第一批 production login-method set 前，比较 vibit 的第一批 login-method candidates。

它不实现 authentication，不添加 credential storage，不添加 external identity linking，不添加 token behavior，不添加 session persistence，不修改 Protobuf envelope，不修改 WebSocket handshake authentication，不添加 runtime player handlers，也不添加 WebSocket routes。

## 2. Evaluation Criteria

每个 candidate 按以下维度评估：

- Production safety。
- Game onboarding ergonomics。
- Agent-readable implementation shape。
- Required public contracts。
- Required storage and migrations。
- Required dependencies。
- Abuse and recovery posture。
- Nakama capability alignment。
- Pitaya session and handler vocabulary alignment。
- Reversibility。

## 3. Candidate Summary

| Candidate | Recommendation | Reason |
| --- | --- | --- |
| `device_credential_login` | 推荐作为第一版 production player login method。 | 它提供低摩擦 game login path，同时不需要 OAuth、OIDC、password hashing、provider SDKs 或 WebSocket handshake changes。如果 credential 是 high entropy、hashed storage、rate limited，并通过显式 contracts 绑定 account lifecycle，它可以生产化。 |
| `guest_anonymous_login` | Deferred。 | 对 onboarding 有用，但很容易过度授权。它需要严格 anonymous actor rules、upgrade behavior、expiration、data ownership 和 abuse posture，才能用于 durable player state。 |
| `custom_id_login` | Deferred，未来可作为 trusted-issuer method。 | 适合已有 identity service 的 studio，但如果任意 client 都能 mint IDs 就不安全。它需要 issuer trust、subject collision、account linking 和 service-auth boundaries。 |
| `email_password_login` | Deferred。 | 常见且熟悉，但需要 password hashing、recovery、rate limiting、reset flows、breach posture 和更强 secret-handling rules。 |
| `external_provider_login` | Deferred。 | 后续 platform accounts 和 social identity 会需要，但 provider dependencies、issuer/audience validation、account linking、conflict handling 和 availability 对第一步过宽。 |
| `service_authentication` | Deferred。 | 对未来 operations 和 server-to-server work 很重要，但它不是 player login method，应该与第一版 player authentication 分离。 |

## 4. Recommended First Set

推荐的第一批 login-method set：

```text
device_credential_login
```

这个建议意味着：

- 第一版 player login method 应证明 client 持有 high-entropy device 或 installation credential。
- Credential 必须被视为 secret proof material，而不是 public device identifier。
- Credential 不得放进现有 Protobuf `Session` metadata fields。
- Credential exchange 应建模为 application/protocol login command payload，除非后续 protocol decision 选择其他 carrier。
- 该方法只有在 account creation/linking policy 被 ratify 后，才能创建新 player account 或认证 existing account。
- Credential storage、token issuance、token storage、session persistence、route registration 和 runtime handlers 仍是未来工作。

这不是直接 Nakama device-auth compatibility。它是把 Nakama 的低摩擦 game login capability 改造成 vibit-native model，并显式处理 credential secrecy、storage、redaction、replay controls 和 account lifecycle boundaries。

## 5. Candidate Details

### Device Credential Login

Position:

```text
recommend_first
```

定义：

```text
Client 证明它持有 high-entropy device 或 installation credential。
```

Benefits:

- 低摩擦 game onboarding。
- 不需要 external identity provider dependency。
- 如果 credential 是 high-entropy secret material，并通过已 ratify 的 hash/lookup boundary 存储，则不需要 password hashing dependency。
- 适合 WebSocket-first gameplay，因为 credential 可通过 login command 在 normal player routes 前交换。
- 比 email/password、social login、OAuth、OIDC 或 provider SDKs 更适合小型第一版 production slice。
- 为 agents 提供从 contract 到 storage 到 validator 到 tests 的窄路径。

Risks:

- Raw OS device IDs 不足以作为 proof。
- Weak identifiers 可 replay。
- Device replacement、reinstall、account recovery 和 account merge behavior 需要后续设计。
- 泄露的 credential 在 revocation 或 rotation 存在前可冒充 player。
- Rate limiting 和 abuse controls 需要显式 gates。

实现前 required artifacts：

- Login command semantic contract。
- Login response semantic contract。
- Credential schema boundary。
- Credential hash and lookup rule。
- Account creation/linking policy。
- Token issuance boundary。
- Error and permission catalog entries。
- Redaction rules。
- Replay and collision tests。
- 防止 credential storage 进入 player account lifecycle tables 的 repository checks。

推荐状态：

```yaml
candidate: device_credential_login
recommended_for_first_set: true
production_capable_after_required_gates: true
creates_player_account: allowed_after_policy_ratification
links_existing_account: deferred
authenticates_existing_account: allowed_after_policy_ratification
requires_major_dependency: false
requires_credential_storage: true
requires_external_identity_linking: false
requires_protobuf_envelope_change: false
requires_websocket_handshake_change: false
confidence: high
```

### Guest Or Anonymous Login

Position:

```text
defer
```

Benefits:

- 最快 onboarding path。
- 适合 try-before-register gameplay。
- 对 local development 和 smoke tests 有用。

Risks:

- 容易混淆 anonymous actor 和 durable player。
- Abuse 和 spam controls 立刻变得重要。
- Account upgrade behavior 是 product-sensitive。
- Anonymous durable state 会制造 recovery 和 ownership disputes。

实现前 required artifacts：

- Anonymous actor contract。
- Permission limits。
- Expiration and upgrade rules。
- Data ownership rules。
- Abuse and rate-limit posture。
- 证明 anonymous identity 不能满足 player-owned production permissions 的 tests。

Recommendation:

```text
Do not include in the first production login-method set.
```

Guest 或 anonymous login 后续可能有用，但不应作为第一版 production proof path，因为 vibit 当前 architecture 正在刻意防止 metadata-only identity 变成 authority。

### Custom ID Login

Position:

```text
defer_trusted_issuer_variant
```

Benefits:

- 当 studio 已有 identity service 时有用。
- 如果只由 trusted services 调用，可以相对简单。
- 对齐 Nakama 的 custom identifier capability coverage。

Risks:

- 如果任意 client 都能选择自己的 custom IDs，则不安全。
- 需要 issuer trust boundaries。
- 需要 subject collision rules。
- 需要 account linking 和 recovery semantics。
- 通常意味着 service authentication 要先于 player authentication。

实现前 required artifacts：

- Trusted issuer model。
- Service-auth 或 caller-auth boundary。
- Subject namespace and collision rules。
- Account linking policy。
- Replay and audit behavior。

Recommendation:

```text
Defer until service-auth and issuer boundaries are ratified.
```

### Email Password Login

Position:

```text
defer
```

Benefits:

- 用户熟悉。
- 支持 cross-device account recovery。
- 在通用 backend systems 中常见。

Risks:

- 需要 password hashing dependency adoption。
- 需要 password reset、recovery、breach posture 和 rate limiting。
- 会过早引入 sensitive secret-handling 和 support workflows。
- Public contract surface 比第一步所需更大。

实现前 required artifacts：

- Password hash dependency adoption。
- Credential schema。
- Password policy。
- Reset and recovery contracts。
- Rate-limit and lockout policy。
- Redaction rules。
- Security tests。

Recommendation:

```text
Defer until the first credential/token/session slice is stable.
```

### External Provider Login

Position:

```text
defer
```

Benefits:

- 强 platform fit。
- Cross-device identity。
- 对 platform stores 和 social platforms 上的 production games 很重要。
- 对齐 Nakama 的 broad provider coverage。

Risks:

- Provider SDKs 和 validation dependencies。
- Issuer、audience、key 和 token validation complexity。
- Account linking conflicts。
- Provider outages 和 metadata retention rules。
- 不同 providers 有不同 identity 和 token semantics。

实现前 required artifacts：

- Provider namespace and subject schema。
- External identity link schema boundary。
- Dependency adoption records。
- Issuer and audience validation。
- Conflict、unlink、recovery 和 merge behavior。
- Provider metadata redaction。

Recommendation:

```text
Defer until external identity linking is ratified.
```

### Service Authentication

Position:

```text
defer_separate_track
```

Benefits:

- 未来 operations、internal services 和 server-to-server work 需要。
- 后续可帮助 trusted custom ID login。

Risks:

- 不是 player login method。
- 容易过度授权。
- 需要 key management、rotation、permission scope 和 audit behavior。
- 可能与未来 distributed runtime design 交互。

实现前 required artifacts：

- Service actor model。
- Service permission catalog。
- Key or proof material lifecycle。
- Rotation and revocation semantics。
- Audit events。

Recommendation:

```text
Defer to a separate service-auth milestone or sub-milestone.
```

## 6. Comparative Matrix

Scores 是 qualitative：

```text
high = favorable
medium = manageable with gates
low = unfavorable for first slice
```

| Candidate | Production safety | Game ergonomics | Agent clarity | Dependency load | Storage complexity | Abuse/recovery load | First-slice fit |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `device_credential_login` | Medium | High | High | High | Medium | Medium | High |
| `guest_anonymous_login` | Low | High | Medium | High | Medium | Low | Medium |
| `custom_id_login` | Medium only with trusted issuer | Medium | Medium | Medium | Medium | Medium | Medium |
| `email_password_login` | Medium | Medium | Medium | Low | Medium | Low | Low |
| `external_provider_login` | Medium | High | Low | Low | Medium | Medium | Low |
| `service_authentication` | Medium | Not player-facing | Medium | Medium | Medium | Medium | Low for player login |

Interpretation:

- `device_credential_login` 是最佳第一候选，因为它的复杂度主要是 local 且 contract-checkable。
- `guest_anonymous_login` 很诱人，但除非先仔细设计 anonymous permissions 和 upgrade behavior，否则会削弱项目当前纪律。
- `custom_id_login` 后续很有价值，但要等 trusted issuer semantics 存在。
- `email_password_login` 和 `external_provider_login` 是重要产品能力，但对第一版 authentication slice 过宽。
- `service_authentication` 应与 player login 分离。

## 7. Reference Mapping

### Nakama

Nakama 支持多种 authentication methods，包括 device、email、social/provider 和 custom identifier methods。它也会在 authentication 后产生 session tokens，并支持基于 refresh-token 的 session continuation。

vibit 将其改造为：

- 第一 capability target：low-friction player login。
- 第一推荐方法：high-entropy device credential login。
- Deferred capability targets：email/password、social/provider login、custom identifier login、guest/anonymous login 和 refresh-token lifecycle。
- 当前 rejected：direct Nakama API compatibility。

### Pitaya

Pitaya 在这里的有用输入是 sessions、request handlers、session binding、frontend/backend separation、route handling 和 push 相关词汇。

vibit 将其改造为：

- 未来 validated request identity 会在 domain dispatch 前传给 handlers。
- 未来 session binding 必须通过 application-owned validation results 发生，而不是通过 transport metadata。
- 第一版 login method selection 不得把 authentication 放进 WebSocket acceptors。
- Frontend/backend 与 distributed session routing 继续 deferred。

## 8. Recommended W-0066 Ratification Packet

W-0066 应 ratify：

```yaml
first_login_method_set:
  - device_credential_login
deferred_login_method_families:
  - guest_anonymous_login
  - custom_id_login
  - email_password_login
  - external_provider_login
  - service_authentication
recommended_first_carrier_posture: login_command_payload_before_normal_gameplay_routes
requires_before_implementation:
  - semantic_login_contract
  - credential_schema_boundary
  - credential_hash_lookup_rule
  - account_creation_or_lookup_policy
  - token_issuance_boundary
  - error_catalog_entries
  - permission_catalog_entries
  - redaction_rules
  - repository_checks
  - focused_tests
does_not_authorize:
  - runtime_authentication_code
  - token_parsing
  - credential_tables
  - external_identity_tables
  - token_tables
  - session_tables
  - protobuf_envelope_change
  - websocket_handshake_authentication
  - runtime_player_handlers
  - websocket_routes
```

## 9. Open Questions For Later Work

- Initial credential 是 client 生成、server 在 bootstrap exchange 中签发，还是二者都允许。
- Login command 是默认创建 player account，还是需要 explicit create intent。
- Credential rotation 是否属于第一版 implementation，还是 deferred 到 token lifecycle work。
- Account recovery 是否存在于第一版 implementation，还是明确 deferred。
- W-0067 和 W-0068 在 login-method set ratify 后会选择哪种 token model。

这些问题不是 W-0065 的 blockers。它们应由 W-0066 到 W-0071 回答。
