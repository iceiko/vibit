# Login Method And Token Format Ratification Standard 中文版

状态：Draft v0.1
最后更新：2026-05-14
范围：第一批 production login-method selection、token model selection、proof carrier boundaries、lifecycle semantics、storage implications，以及实现前 gates
依赖：`docs/authentication-token-session-validation.md`
权威决策：`ADR-0024`

本文件是 `docs/login-method-token-format-ratification.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

## 1. 目的

本标准定义 vibit 如何在实现之前 ratify 第一批 production login methods 和 token model。

目标不是快速加入 authentication code。目标是让第一批 authentication choices 足够明确，使后续 agents 能实现它们，而不会把 security behavior 偷偷放进 transport handlers、Protobuf adapters、player account persistence 或 domain modules。

本标准可用于 ratify：

- 第一批 login-method set。
- 第一版 token model。
- Token 与 proof carrier boundaries。
- Token lifecycle semantics。
- 必需的 contracts、schemas、dependencies、checks 和 tests。
- 有边界的 implementation queue。

本标准本身不实现：

- Runtime authentication。
- Token parsing、signing、validation、refresh、revocation、rotation、replay handling 或 storage。
- Credential storage。
- External identity linking。
- Session persistence。
- Protobuf envelope changes。
- WebSocket handshake authentication。
- Runtime player account handlers。
- WebSocket routes。

## 2. 必读材料

本标准应与以下材料一起阅读：

- `docs/authentication-token-session-validation.md`
- `docs/authentication-proof-token-session-contract-dimensions.md`
- `docs/credential-storage-external-identity-linking-boundaries.md`
- `docs/session-persistence-websocket-handshake-decision-gates.md`
- `docs/player-identity-session-boundary.md`
- `docs/player-account-session-contracts.md`
- `docs/game-protocol.md`
- `docs/runtime-protocol-adapter.md`
- `docs/postgresql-persistence-boundary.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/player/module.yaml`
- `runtime/AGENTS.md`
- `ADR-0019`
- `ADR-0021`
- `ADR-0022`
- `ADR-0023`
- `ADR-0024`

参考阅读：

- Nakama authentication concepts: `https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama session concepts: `https://heroiclabs.com/docs/nakama/concepts/session/`
- Pitaya API and handler vocabulary: `https://pitaya.readthedocs.io/en/latest/API.html`
- Pitaya session and server-role features: `https://pitaya.readthedocs.io/en/stable/features.html`

Nakama 和 Pitaya 是参考。它们不支配 vibit 的 public API shape、credential schema、token format、generated file conventions 或 agent workflow。

## 3. Ratification 姿态

Authentication choices 必须作为持久 architecture decisions 来做，而不能成为偶然的 implementation details。

规则：

- Login methods 是 public behavior，必须在 handlers 存在前 ratify。
- Token format 是 security 与 operations decision，必须在 parser 或 validator code 存在前 ratify。
- Token carrier behavior 必须在修改 Protobuf envelope fields、WebSocket handshake behavior 或 per-request payload contracts 前 ratify。
- Credential storage 与 external identity linking 仍是独立边界。
- Session persistence 仍是独立边界。
- 当前 metadata-only `player_id`、`session_id` 和 `connection_id` 仍然不是 authenticated proof。
- 被选择的 login method 不会自动授权该方法需要的所有 dependencies、tables、routes 或 protocol changes。

## 4. Reference Alignment

### Nakama

Nakama 表明成熟 game backend 通常支持多种 authentication methods、account linking、session tokens、refresh tokens、session expiration、logout 和 realtime socket behavior。

vibit 采纳它作为 capability coverage，但不采纳它的 API shape。

参考定位：

| Nakama concept | vibit position |
| --- | --- |
| Multiple authentication methods | 在选择第一批方法前需要比较的 capability coverage。 |
| Device-style authentication | 候选 login family；只有 credential 和 abuse semantics 明确后才能改造采纳。 |
| Email/password authentication | 候选 login family；除非 password hashing、reset、rate limit 和 recovery gates 被 ratify，否则 deferred。 |
| Social/provider authentication | 候选 login family；直到 provider subject 与 external identity linking 被 ratify 前 deferred。 |
| Custom identifier authentication | 候选 login family；只有 issuer/trust boundaries 明确后才允许。 |
| Session token | 候选 token concept，但不会自动 Nakama-compatible。 |
| Refresh token | 候选 lifecycle concept，需要 storage、rotation、revocation、replay 和 redaction decisions。 |
| Session variables | Deferred；不得变成 opaque token claims 中的隐藏 authority。 |
| Logout invalidating tokens | 任何 token model 的必需比较维度。 |
| Realtime socket bound to an authenticated session | 通过 vibit 的 application-owned request identity 和未来 session validation model 改造。 |

### Pitaya

Pitaya 展示了 Go game server 在 handlers、request context、session binding、frontend/backend roles、message routing、push 和 session lifecycle 方面的词汇。

vibit 采纳有用词汇，但不采纳 Pitaya 的 API shape 或 cluster assumptions。

参考定位：

| Pitaya concept | vibit position |
| --- | --- |
| Handler receives request context containing session | 改造为 domain dispatch 前的 `RequestIdentity` 和 `SessionValidationResult`。 |
| Session binding to user ID | 候选未来 session-binding concept，不是 token format。 |
| Session data storage between requests | Deferred，直到 session persistence 被 ratify。 |
| Frontend/backend session distinction | Deferred，直到 distributed runtime planning。 |
| Route-aware routing function | 作为未来 route policy mapping 的参考；不得绕过 module contracts。 |
| Push to bound users | 未来 realtime capability；除非 session lifecycle tests 需要，否则不属于第一批 authentication ratification。 |

## 5. Login-Method Candidate Families

未来工作必须在选择第一批 production set 前比较候选 login methods。

| Family | Description | Main benefits | Main risks | Required gates before implementation |
| --- | --- | --- | --- | --- |
| Device credential login | Client 证明它持有 device-scoped high-entropy secret 或 identifier。 | 适合游戏低摩擦进入，可作为第一类账号入口。 | Device identifiers 可能变化；弱 identifiers 可 replay；account recovery 和 abuse controls 很重要。 | Credential semantics、storage/hashing、account creation/linking policy、replay controls、rate limits、tests。 |
| Guest or anonymous login | Server 创建 temporary 或 low-assurance actor，不依赖 durable external proof。 | 快速 onboarding，适合 early game entry。 | Abuse、data ownership、upgrade behavior 和 production permission scope 必须严格。 | Anonymous actor rules、upgrade path、expiration、anti-abuse posture、permission limits。 |
| Custom ID login | Trusted issuer 把 custom subject 映射到 player。 | 适合集成已有 game platform 或 studio identity。 | 如果任意 client 都能 mint IDs 会不安全；issuer trust 必须明确。 | Issuer model、trusted caller boundary、subject collision rules、linking policy。 |
| Email/password login | User 使用 email 或 username 加 password 认证。 | 熟悉且可恢复。 | 需要 password hashing、reset flows、rate limits、breach posture 和 secret handling。 | Password dependency adoption、credential schema、reset/recovery flow、rate limiting、redaction。 |
| External provider login | User 通过 platform、social、OAuth、OIDC 或 game-provider account 认证。 | 适配平台和跨设备 identity。 | Provider dependencies、issuer/audience validation、account linking conflicts、token validation、availability。 | Provider subject schema、dependency adoption、issuer/audience validation、linking/conflict/recovery behavior。 |
| Service authentication | Non-player service 证明 authority。 | 未来 internal operations 和 server-to-server work 会需要。 | 容易过度授权；必须与 player login 分离。 | Service actor model、permissions、key management、rotation、audit events。 |

选择规则：

- 第一批方法应刻意保持小。
- 被选择的 login method 必须有明确 confidence level 和 known gaps。
- 被选择的 login method 必须说明它是创建 player accounts、链接 existing accounts，还是只认证 existing accounts。
- 被选择的 login method 必须说明它是 production-capable、bootstrap-only，还是 local-development-only。
- 任何被选择的方法都不得把 credential parsing 放进 WebSocket transport 或 Protobuf adapter code。

## 6. Token Model Dimensions

未来工作必须按这些维度比较 token models。

### Token Kinds

候选 token kinds：

- Access token。
- Refresh token。
- Session token。
- One-time proof token。
- Service token。
- External provider token。

规则：

- Token kind 必须说明它是 client-presented、server-stored、derived from external provider，还是只在 trusted server boundary 内使用。
- Access 与 refresh token behavior 必须分离，即使二者都是 opaque strings。
- Refresh token 不得在没有 rotation、revocation、storage、replay 和 redaction semantics 的情况下引入。
- External provider tokens 不得因为方便就成为 vibit session tokens。

### Token Formats

候选 formats：

| Format | Benefits | Risks | Required gates |
| --- | --- | --- | --- |
| Opaque high-entropy token | 容易 redact，server-side revoke 简单，implementation 可保持 storage-backed 且显式。 | 需要 lookup storage 和谨慎 hashing。 | Token generation、hash storage、lookup index、expiration、rotation、cleanup。 |
| Signed structured token | 可无 storage lookup 验证，并可携带 claims。 | Key management、revocation、replay、claim drift 和 secret rotation 更难。 | Signing dependency、key management、issuer/audience/claims、revocation story。 |
| External provider token | 复用 provider proof。 | Provider-specific validation 和 availability；audience 与 issuer 错误风险很高。 | Provider validation、trust boundary、subject mapping、failure handling。 |
| Plain session ID as token | 词汇简单。 | 容易混淆 identifier 与 proof；除非 high entropy 且作为 secret 保护，否则不安全。 | 必须当作 secret token，而不是 metadata-only `session_id`。 |

选择规则：

- Token format 必须说明 issuer、verifier、subject、audience、entropy 或 signing model、expiration、refresh、revocation、rotation、replay、redaction 和 storage implications。
- 除非未来 decision 明确授予例外并记录理由，否则 token strings 绝不能 plaintext 存储。
- Logs、errors、traces、conversation logs 和 change specs 必须 redact token values。
- Token validation errors 必须在实现前映射到 registered error dimensions。

## 7. Proof Carrier Boundaries

Proof carrier 是 credentials、tokens 或 session proof 在 wire 或 request 内出现的位置。

候选 carriers：

| Carrier | Layer touched | Position |
| --- | --- | --- |
| Login command payload | Application/protocol payload after decode | 候选第一 carrier，因为它保持 WebSocket transport credential-neutral。 |
| Per-request payload field | Application/protocol payload after decode | 可作为显式 validation 候选，但 verbose 且容易重复。 |
| Protobuf envelope extension | Wire envelope | 需要 protocol decision；不得复用当前 metadata fields。 |
| Current `Session.session_id` field | Existing envelope metadata | 当前 metadata-only；未经 protocol ratification 不得成为 proof。 |
| WebSocket handshake header | Transport/process boundary | 需要 handshake-auth decision 和 normalized identity handoff。 |
| WebSocket subprotocol、cookie 或 query parameter | Transport/process boundary | 需要明确 risk analysis；不能因为方便就使用。 |
| First system message | Protocol/application after connection open | 需要 system-message contract 和 connection binding rules。 |

规则：

- 当前 Protobuf `Session` message 在未来 protocol decision 前仍是 metadata-only。
- 不要把 access tokens 或 refresh tokens 放入 `player_id`、`session_id`、`connection_id`、`connection_epoch`、route names、request IDs 或 target IDs。
- WebSocket transport 在未来 handshake decision 赋予 narrow role 前保持 credential-neutral。
- Carrier selection 必须说明 proof 如何在 domain handlers 运行前变成 `RequestIdentity`。

## 8. Ratification Packet

每个未来选择 login method、token model、token carrier 或 lifecycle rule 的 work item，都必须在 change spec 或 documentation 中留下 ratification packet。

该 packet 必须包括：

- Selected option。
- Rejected plausible options。
- 与 Nakama 和 Pitaya 的 reference alignment。
- Public rationale。
- Decision weights。
- Security and abuse notes。
- Public contract impact。
- Protobuf impact。
- WebSocket handshake impact。
- Persistence and migration impact。
- Dependency adoption impact。
- Error、permission 和 audit-event impact。
- Test and repository-check impact。
- Redaction rules。
- Reversal conditions。
- Next implementation gate。

## 9. 实现前必需 Artifacts

在 selected login method 或 token model 的 implementation code 存在前，repository 必须具备：

- Login commands、token refresh、logout、validation 以及 scope 内 account-linking behavior 的 semantic contracts。
- 对 invalid proof、expired proof、revoked proof、replayed proof、malformed proof、credential mismatch、account disabled、rate limited 和 dependency unavailable 等适用场景的 error catalog entries。
- 对 player login、token refresh、logout、account linking、service-auth validation 和 administrative revocation 等适用场景的 permission catalog entries。
- 如果需要 storage，则有 credential records、token/session records 或 external identity links 的 schema boundaries。
- Migration sources 只能在 schema boundaries ratify 后添加。
- Hashing、signing、OAuth/OIDC、provider SDKs、Redis-like stores 或 key-management libraries 的 dependency adoption records。
- Credential 和 token material 的 redaction 与 logging rules。
- Validator interfaces 和 storage adapters 的 runtime ownership rules。
- 防止 shortcut implementations 的 repository checks。
- 覆盖 success、invalid proof、expired proof、revoked proof、replay/collision behavior、redaction 和 boundary ownership 的 focused tests。
- Bilingual documentation updates。

## 10. M-013 Work Queue

M-013 按有边界的步骤推进：

1. 定义本 ratification standard。
2. 比较第一批 login-method candidates。
3. Ratify 第一批 login-method set。
4. 比较 token format 和 token carrier options。
5. Ratify 第一版 token format 和 proof carrier posture。
6. 定义 token lifecycle 和 storage implications。
7. 定义 authentication contract、error 和 permission surfaces。
8. 定义 credential、token 和 session schema gates。
9. 为 selected login/token boundaries 添加 repository checks。
10. 关闭 milestone 并创建下一个 implementation gate。

如果 verification 或 reference analysis 证明某一步过宽或过窄，未来 work item 可以调整该队列。不得因为方便而把它折叠成 implementation work。

## 11. Repository Rules

Agents 不得：

- 仅凭本标准实现 login handlers。
- 仅凭本标准 parse 或 validate tokens。
- 仅凭本标准添加 token、credential、external identity 或 session tables。
- 仅凭本标准添加 password hashing、JWT、OAuth、OIDC、provider SDK、cryptography、key-management、Redis-like 或其他 major authentication dependencies。
- 仅凭本标准修改 Protobuf envelope behavior。
- 仅凭本标准修改 WebSocket handshake authentication behavior。
- 仅凭本标准添加 runtime player account command/query handlers 或 WebSocket routes。
- 把 metadata-only identity 当作 production proof。
- 未经 compatibility ADR 复制 Nakama 或 Pitaya public APIs。

Agents 可以：

- 添加 comparison docs。
- 添加 decision records。
- 添加 change specs。
- 更新 manifests。
- 规划 contracts、schemas、checks 和 tests。
- 强化防止 shortcut implementations 的 repository checks。

## 12. Verification

本标准下的 changes 默认验证包括：

```bash
node tools/vibit inspect next --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check architecture --json
node tools/vibit check memory --json
node tools/vibit check change <change-id> --json
node tools/vibit check all --json
git diff --check
```

只有 runtime code 变化时才需要 runtime Go tests。

只有 persistence behavior 或 migrations 变化时才需要 live PostgreSQL verification；否则记录为 not applicable。

## 13. Exit Criteria

M-013 只有在满足以下条件时才能关闭：

- 第一批 login-method set 已选择，或已带 next gate 明确 deferred。
- 第一版 token model 已选择，或已带 next gate 明确 deferred。
- Proof carrier posture 已选择，或已明确 deferred，且没有隐式 Protobuf 或 WebSocket behavior changes。
- Token lifecycle semantics 已记录。
- 必需的 contracts、schemas、dependency records、checks 和 tests 已规划。
- 没有意外跨过 implementation boundary。
- `node tools/vibit check all --json` 通过。
