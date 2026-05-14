# First Login Method Set 中文版

状态：Draft v0.1
最后更新：2026-05-14
范围：M-013 的第一版 production login-method set ratification
依赖：`docs/first-login-method-candidates.md`
Canonical decision：`ADR-0025`

本文件是 `docs/first-login-method-set.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

## 1. 目的

本文在 runtime authentication implementation 开始前，ratify vibit 的第一版 login-method set。

它选择第一版面向 production 的 player login method，记录被拒绝的 alternatives，并定义在添加 handler、credential table、token parser、session store、Protobuf change、WebSocket handshake authentication、runtime player route 或 WebSocket route 前必须存在的 gates。

本文不实现 authentication。

## 2. Ratified Set

第一版 login-method set 是：

```yaml
first_login_method_set:
  - device_credential_login
```

没有其他 login method 属于第一版 set。

## 3. Selected Method

### `device_credential_login`

定义：

```text
Player client 证明它持有 high-entropy installation credential。
```

Credential 是 secret proof material。它不是 raw operating-system device ID、advertising ID、model identifier、connection ID、player ID、session ID 或其他 public metadata。

Ratified posture：

```yaml
method: device_credential_login
actor_kind_after_success: player
production_classification: production_capable_after_required_gates
bootstrap_only: false
local_development_only: false
creates_player_account: allowed_after_account_creation_policy
authenticates_existing_account: allowed_after_credential_lookup_policy
links_existing_account: deferred
recovers_account: deferred
upgrades_anonymous_account: deferred
requires_major_dependency_before_contracts: false
requires_credential_storage_before_implementation: true
requires_external_identity_linking: false
requires_password_hashing: false
requires_oauth_or_oidc: false
requires_provider_sdk: false
requires_protobuf_envelope_change: false
requires_websocket_handshake_authentication: false
confidence: high
```

第一版 implementation 必须保持 WebSocket transport credential-neutral。Credential proof 必须由未来 application-owned authentication boundary 在 protocol decoding 之后、production-sensitive domain dispatch 之前处理。

Token format、access-token behavior、refresh behavior、token carrier behavior、runtime session persistence 和 WebSocket connection binding 仍由 W-0067 到 W-0071 分别决定。

## 4. Public Rationale

`device_credential_login` 是最小的、面向 production 的第一版 login method，符合 vibit 当前目标。

它提供类似 Nakama 这类成熟 game backend 中 device-style capability 的低摩擦游戏入口，但会改造成 vibit terms：

- Proof 是 high-entropy secret material，不是 public metadata。
- Player account lifecycle 与 credential records 分离。
- WebSocket transport 保持 credential-neutral。
- Protobuf `Session` metadata 保持 metadata-only。
- Runtime handlers 只有在 validator 产生 machine-readable result 后才接收 authenticated identity。
- 后续 login families 仍然可能，不会把 provider dependencies 或 password workflows 强行塞进第一版 slice。

这也保留了 Pitaya 风格的 realtime connection/session vocabulary 与 handler request context 分离。vibit 后续应该把 validated identity 绑定到 request context，但绑定必须来自 application-owned validation，而不是 WebSocket acceptors 或 route handlers。

## 5. Rejected Alternatives

| Alternative | Decision | Reason |
| --- | --- | --- |
| `guest_anonymous_login` | Deferred。 | 它对 onboarding 有用，但第一版 production player-owned state 不应依赖 anonymous authority，除非 expiration、upgrade、abuse、recovery 和 permission limits 已 ratify。 |
| `custom_id_login` | Deferred。 | 它只有在 trusted issuer 下才安全。这需要先有 service-auth、subject namespace、collision、linking、replay 和 audit boundaries。 |
| `email_password_login` | Deferred。 | 它需要 password hashing、recovery、reset flows、rate limiting、breach posture 和 support workflows，对第一版 authentication slice 过宽。 |
| `external_provider_login` | Deferred。 | 它对后续 platform identity 很重要，但 provider validation、issuer/audience rules、account linking、conflicts、outages 和 dependencies 现在过宽。 |
| `service_authentication` | Deferred 到独立 track。 | 它不是 player login method，不应混入第一版 player-authentication slice。 |
| Metadata-only `player_id` 或 `session_id` | Rejected。 | Metadata-only values 不是 proof，不得满足 production permissions。 |
| Direct Nakama API compatibility | 当前 rejected。 | Nakama 仍是 capability reference，不是 governing API shape。 |
| Pitaya session binding as implementation API | 当前 rejected。 | Pitaya 仍是 vocabulary input；vibit 拥有自己的 contracts、manifests 和 runtime handoff shape。 |

## 6. Decision Weights

```yaml
decision_weights:
  production_safety: medium
  game_onboarding_ergonomics: high
  agent_context: high
  dependency_load: low
  storage_complexity: medium
  abuse_and_recovery_load: medium
  reversibility: high
  long_term_maintainability: high
confidence: high
```

接受 medium storage 和 abuse complexity 的主要原因是：这些复杂度保持 local、contract-checkable，并且可以延后到显式 gates 之后。

## 7. Required Gates Before Implementation

在实现 `device_credential_login` 前，repository 必须具备：

- W-0067 与 W-0068 ratify 的 token format 和 proof carrier posture。
- W-0069 定义的 token lifecycle 和 storage implications。
- W-0070 定义的 authentication command、response、error、permission 和 audit surfaces。
- W-0071 定义的 credential、token 和 session schema gates。
- W-0072 添加的 selected login/token boundary repository checks。
- Selected method 的 semantic login contract。
- Credential storage boundary，确保 credentials 不进入 `player_accounts` 和 `player_account_events`。
- Credential hash and lookup rule。
- Player account creation or lookup policy。
- 对 credentials、tokens、logs、traces、errors、change specs、conversation logs 和 tests 的 redaction rules。
- 针对 success、missing proof、malformed proof、invalid proof、replay 或 collision behavior、redaction 和 boundary ownership 的 focused tests。

Implementation 继续 deferred，直到这些 gates 完成，或被后续 accepted decision 明确 supersede。

## 8. Non-Authorization

本 ratification 不授权：

- Runtime authentication code。
- Login handlers。
- Token parsing、signing、validation、refresh、revocation、rotation、replay handling 或 storage。
- Credential tables。
- External identity tables。
- Token tables。
- Session tables。
- Migrations。
- Password hashing、JWT、OAuth、OIDC、provider SDK、Redis-like、cryptography、key-management 或 major authentication dependencies。
- Protobuf envelope changes。
- WebSocket handshake authentication。
- Runtime player account handlers。
- WebSocket routes。
- 将 metadata-only `player_id`、`session_id` 或 `connection_id` 视为 proof。

## 9. Known Gaps

以下问题有意留给后续 M-013 work：

- Initial installation credential 是 client-generated、server-issued，还是二者都允许。
- First login 是否默认创建 player account，还是要求 explicit create intent。
- Credential rotation 是否包含在第一版 implementation 中。
- Account recovery 是否包含在第一版 implementation 中。
- Rate limiting 是否需要独立 store，还是可以先从 process-local limits 开始。
- 选择哪种 token format 和 token carrier。
- 第一版 production slice 是否包含 refresh tokens。
- Runtime sessions 是 persisted，还是从 token validation 派生。

这些 gaps 不是 ratify login-method set 的 blockers。它们是 implementation 的 blockers。

## 10. Reference Alignment

### Nakama

Nakama 仍是 broad game authentication coverage 的 capability reference，包括 device-style、email、custom identifier、provider login、sessions、refresh 和 logout。

vibit 采纳 low-friction device-style capability 作为第一版 player login direction，但不复制 Nakama 的 API shape 或 token/session semantics。

### Pitaya

Pitaya 仍是 session binding、handler context、routing 和 realtime server structure 的 vocabulary reference。

vibit 将 session-context idea 改造成 application-owned request identity。被选中的 login method 不得把 authentication 放进 WebSocket acceptors、transport handlers、protocol adapters 或 domain handlers。

## 11. Follow-Up

下一项工作：

```text
W-0067 Compare token format and carrier options
```

W-0067 必须比较 token formats 和 carrier postures，不得假设选择 `device_credential_login` 就自动选择 JWT、opaque tokens、refresh tokens、Protobuf envelope changes、WebSocket handshake authentication 或 session persistence。
