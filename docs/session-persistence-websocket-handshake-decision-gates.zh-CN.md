# Session Persistence 与 WebSocket Handshake 决策门

状态：Draft v0.1
最后更新：2026-05-14
范围：Session persistence decision gates、WebSocket handshake authentication decision gates、request-level validation options、reconnect gates、connection epoch gates、Protobuf envelope interaction gates，以及未来 implementation artifacts
依赖：`docs/authentication-token-session-validation.md`

配套英文原文是 `docs/session-persistence-websocket-handshake-decision-gates.md`。英文文件是权威版本。

## 1. 目的

本标准定义 vibit 在实现 session persistence 或 WebSocket handshake authentication 之前必须具备的 decision gates。

目标是让未来 session work 足够显式，Agent 可以扩展它，而不是把 authentication 或 binding behavior 隐藏在 transport、protocol 或 domain modules 中。

本标准不选择：

- Request-level validation 作为生产模型。
- First-message validation 作为生产模型。
- WebSocket handshake-level validation 作为生产模型。
- Every-request validation 作为生产模型。
- Hybrid validation model。
- Session store。
- Session tables 或 migrations。
- Token carrier behavior。
- Protobuf envelope changes。
- WebSocket handshake/system messages。
- WebSocket handshake authentication behavior。
- Route-level authentication implementation。

## 2. 必读材料

本标准需要与以下材料一起阅读：

- `docs/authentication-token-session-validation.md`
- `docs/authentication-proof-token-session-contract-dimensions.md`
- `docs/credential-storage-external-identity-linking-boundaries.md`
- `docs/player-identity-session-boundary.md`
- `docs/game-protocol.md`
- `docs/runtime-protocol-adapter.md`
- `docs/postgresql-persistence-boundary.md`
- `.arch/runtime.yaml`
- `.arch/protocol.yaml`
- `.arch/reference.yaml`
- `runtime/AGENTS.md`
- `ADR-0015`
- `ADR-0018`
- `ADR-0021`
- `ADR-0023`

参考阅读：

- Nakama authentication 和 sessions concepts：`https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama realtime socket concepts：`https://heroiclabs.com/docs/nakama/concepts/multiplayer/`
- Pitaya session、handler 和 frontend/backend vocabulary：`https://pitaya.readthedocs.io/`

Nakama 和 Pitaya 是参考对象。它们不支配 vibit 的 public API shape、validation model、session persistence model、envelope behavior、WebSocket handshake behavior 或 agent workflow。

## 3. 核心词汇

### Session Persistence

Session persistence 是 logical session state 在单个 request 或 in-memory validation step 之外的 server-side storage。

未来可能的 state 包括 `session_id`、`actor_kind`、`actor_id`、`player_id`、claims、expiration、revocation status、refresh linkage、connection binding metadata 和 audit metadata。

规则：

- 本标准不实现 session persistence。
- Persisted sessions 不是 player account lifecycle rows。
- Persisted sessions 不是 WebSocket connections。
- Persisted sessions 不是 Protobuf envelope metadata。
- Session storage 不得为了方便添加到 player、inventory、WebSocket transport 或 Protobuf adapter code 中。

### WebSocket Handshake Authentication

WebSocket handshake authentication 是在 WebSocket connection establishment 之前或期间 validation proof。

规则：

- 本标准不实现 handshake authentication。
- 当前 WebSocket transport 保持 credential-neutral。
- 除非未来 handshake standard 明确授予狭窄 transport responsibility，否则 WebSocket transport 不得检查 `Authorization`、cookies、subprotocol authentication carriers、token query parameters 或 credential headers。
- 即使未来采用 handshake authentication，application dispatch 仍必须收到 normalized request identity。

### Request-Level Validation

Request-level validation 是在 envelope decode 之后、domain dispatch 之前，对每个 request validation proof 或 session identity。

规则：

- 这是未来选项，不是已选择的生产模型。
- 它与当前 application-owned `SessionValidator` hook 对齐。
- 它未来可以与其他 validation gates 组合。

### First-Message Validation

First-message validation 是 client 在 WebSocket connection establishment 之后发送 protocol message，用于在普通 gameplay routes 被接受前 authenticate 或 bind connection。

规则：

- 这是未来选项，不是已选择的生产模型。
- 它在实现前需要显式 protocol/system-message contracts。
- 它不得通过 ad hoc domain routes 被临时发明出来。

### Every-Request Validation

Every-request validation 是每个 command 或 query 都携带足够的 proof 或 session metadata，以便在 domain dispatch 前独立 validation。

规则：

- 这是未来选项，不是已选择的生产模型。
- 它可以减少对 connection state 的依赖，但可能增加 validation cost 和 protocol verbosity。
- Token/session carrier behavior 必须先 ratify，之后才能实现。

### Hybrid Validation

Hybrid validation 组合多个 validation gates。

例子：

- Handshake validation 加 request-level permission checks。
- First-message binding 加 periodic session revalidation。
- Sensitive routes 使用 every-request validation，low-risk routes 使用 cached session context。

规则：

- 这是未来选项，不是已选择的生产模型。
- Hybrid model 必须声明每一步由哪个 layer 拥有，以及 failures 如何处理。

## 4. Validation Model 决策门

选择生产 validation model 前，未来工作至少必须比较这些选项：

| Option | Layer touched | Benefits | Risks | Required before implementation |
| --- | --- | --- | --- | --- |
| Request-level validation | Application dispatch after protocol decode | Domain handoff 清晰，适用于 multiple transports，匹配当前 `SessionValidator` hook | 如果没有 cache/session rules，可能重复 validation | Validation contract、proof/session carrier、failure behavior、route policy mapping |
| First-message validation | Protocol/application after connection open | 保持 transport credential-neutral，并支持显式 game protocol negotiation | 需要 connection-bound state 和 system-message contracts | System-message contract、connection binding model、timeout/error behavior、reconnect rules |
| Handshake-level validation | WebSocket transport/process boundary | 早期拒绝 unauthenticated connections，并可减少 unauthenticated connection load | 容易把 auth 放进 transport，且会让非 WebSocket transports 更复杂 | Transport auth boundary、header/subprotocol/cookie/query carrier decision、normalized identity handoff |
| Every-request validation | Protocol/application on each request | Stateless-friendly，per-request proof semantics 清晰 | Overhead 更高，envelopes 或 payloads 更大，carrier leakage 风险更高 | Token/session carrier contract、replay handling、failure/retryability rules |
| Hybrid validation | Multiple layers | 可以平衡 early rejection、connection binding 和 route sensitivity | 如果 manifests 不精确，很容易模糊 ownership | Layer ownership matrix、cache invalidation、failure precedence、verification plan |

选择这些选项中的任何一种都是 ask-first decision。

在模型被选择前：

- 当前 metadata-only validation 仍是 non-authenticated。
- `Session.connection_id`、`Session.session_id`、`Session.player_id` 和 `Session.connection_epoch` 仍然只是 metadata carriers。
- Domain modules 不得把 metadata-only identity 当作 production proof。
- WebSocket transport 保持 credential-neutral。

## 5. Session Store 决策门

添加 session store 前，未来工作必须定义：

- Sessions 是否 server-side persisted。
- Store 是 PostgreSQL、memory、Redis-like、其他被 ratify 的 store，还是没有 store。
- Session ID generation semantics。
- Lookup semantics。
- Expiration semantics。
- Refresh semantics。
- Revocation 和 logout semantics。
- 如果涉及 tokens，则定义 rotation 和 replay behavior。
- Connection binding 和 connection epoch behavior。
- Cleanup 和 migration behavior。
- 与 account lifecycle 或 token changes 的 transaction boundaries。
- Opt-in live verification requirements。

规则：

- 本标准不选择 session store。
- 本标准不得添加 session table。
- PostgreSQL 是 module state 的第一版 authoritative durable store，但这不会自动让 PostgreSQL 成为 session store。
- Memory 可用于 tests 或 local bootstrap，但这不会自动让 memory 成为 production session store。
- Redis-like storage 仍然只是未来选项，只有在 dependency adoption 和 architecture ratification 后才可能选择。

## 6. WebSocket Handshake 决策门

添加 WebSocket handshake authentication 前，未来工作必须定义：

- Handshake 前或期间使用哪种 proof carrier。
- Carrier 是 header、cookie、query parameter、WebSocket subprotocol、mTLS-like service proof，还是其他被 ratify 的 input。
- 如何支持 browser 和 non-browser clients。
- Proof 失败时是 reject handshake，还是接受 connection 进入 anonymous/non-authenticated state。
- 哪个 component 从 handshake validation 创建 normalized request identity。
- Handshake 后如何 revalidate identity。
- Logout、revocation、expiration、refresh、reconnect 和 connection migration 如何影响 active WebSocket connections。
- 当连接在 Protobuf envelope 存在前被拒绝时，errors 如何暴露。
- 哪些 tests 证明 transport code 保持狭窄。

规则：

- WebSocket transport 不得拥有 player account lookup。
- WebSocket transport 不得拥有 credential lookup。
- WebSocket transport 不得拥有 long-lived session persistence。
- 未来任何 transport-level proof extraction 都必须 hand off 给 application-owned 或 auth-owned validation contracts。
- 除非 application dispatch 传递 normalized request identity，否则 route-level domain handlers 不得假设 handshake identity。

## 7. Protobuf Envelope 交互门

当前 envelope 包含：

```text
Session.connection_id
Session.session_id
Session.player_id
Session.connection_epoch
```

本标准不改变这些字段。

改变 envelope behavior 前，未来工作必须定义：

- Token/session proof 是通过 envelope metadata、system message、module payload、WebSocket handshake carrier，还是其他设计携带。
- `session_id` 是 candidate identifier、validated identifier，还是在不同 states 中两者皆可。
- `player_id` 是否可以在 validation 前出现，以及如何标记为 untrusted。
- `connection_epoch` 是 server-issued、client-presented，还是两者皆可。
- 新字段是否需要 Protobuf package version change。
- Generated output impact。
- Backward compatibility behavior。
- Unsupported、missing、malformed、expired 或 revoked proof 的 error behavior。

规则：

- Envelope fields 在 validation 存在前仍然只是 metadata carriers。
- 添加 token fields、credential fields、authentication result fields、handshake fields 或新的 system messages，需要 protocol change spec 和 ADR。
- Generated Protobuf output 不得手工编辑。

## 8. Reconnect 与 Connection Epoch 门槛

实现 reconnect 或 connection epoch behavior 前，未来工作必须定义：

- Reconnect 是 tied to session persistence、token proof、connection IDs，还是其他 proof。
- 谁签发 `connection_epoch`。
- Epoch 是每次 connection、reconnect、migration 增加，还是只在 successful validation 后增加。
- Stale epochs 是 rejected、ignored，还是 treated as hints。
- 如何处理同一 session 的 duplicate active connections。
- Previous connections 是 closed、replaced，还是 allowed to coexist。
- 未来 room、match、party、presence 和 stream membership 如何恢复。
- 需要什么 replay protection。

规则：

- 当前 `connection_epoch` 仍然只是 metadata。
- 当前 WebSocket transport 不拥有 reconnect semantics。
- 未来 reconnect behavior 必须兼容后续 multiplayer 和 presence modules。

## 9. 参考模式映射

### Nakama Patterns

| Pattern | Vibit position | Reason |
| --- | --- | --- |
| Session token after authentication | Deferred | Token format、issuance、validation、expiration、revocation 和 carrier behavior 仍是独立选择。 |
| Refresh token concept | Deferred | Refresh、rotation、replay、logout 和 storage behavior 需要未来 ratification。 |
| Realtime socket associated with an authenticated session | Adapted | vibit 要求 production-sensitive dispatch 前有 normalized request identity，但尚未选择 handshake-level binding。 |
| Session expiration and revocation vocabulary | Adopted as design dimensions | 未来 session work 必须显式处理这些维度。 |
| Direct Nakama API compatibility | Rejected for now | vibit 定义 agent-native contracts，只有未来 ADR 才能采用兼容性。 |

### Pitaya Patterns

| Pattern | Vibit position | Reason |
| --- | --- | --- |
| Session object separate from acceptor | Adopted as vocabulary | Transport connection state 必须与 application request identity 分离。 |
| Handler receives session context | Adapted | vibit handlers 通过 application dispatch 接收 `RequestIdentity`，不是 transport-owned session object。 |
| Session binding | Adopted as vocabulary | Binding 是有用词汇，但 vibit 要求 ratified validation results 和 manifests。 |
| Frontend/backend split | Deferred | Distributed routing 和 cluster topology 仍是未来工作。 |
| Direct Pitaya API compatibility | Rejected for now | Architecture vocabulary 可以影响 vibit，但不复制 public APIs。 |

## 10. 未来产物门槛

实现 session persistence 前，未来工作必须新增或更新：

- `changes/` 下 change spec。
- 当选择 validation model、store、reconnect behavior 或 transport responsibilities 时新增 ADR。
- Session persistence standard。
- Session contract sources。
- Error catalog 和 retryability rules。
- Permission catalog。
- Store ownership manifest。
- 如果存储数据，则新增 schema boundary 和 migrations。
- Session cleanup 和 expiration plan。
- Creation、lookup、expiration、revocation、logout、refresh、reconnect 和 failure behavior 的 tests。
- 当规则可以静态执行时新增 repository checks。
- 英文文档和简体中文译本。

实现 WebSocket handshake authentication 前，未来工作必须新增或更新：

- `changes/` 下 change spec。
- Handshake carrier 和 validation model 的 ADR。
- WebSocket handshake authentication standard。
- Transport/auth handoff contract。
- Protocol 或 non-protocol error surface definition。
- Browser 和 non-browser client compatibility notes。
- 证明 transport responsibilities 保持狭窄的 tests。
- 证明 application dispatch 收到 normalized request identity 的 tests。
- 当规则可以静态执行时新增 repository checks。
- 英文文档和简体中文译本。

## 11. Ask-First 边界

在以下操作前询问 maintainer：

- 选择 request-level、first-message、handshake-level、every-request 或 hybrid validation 作为生产模型。
- 选择 session persistence store。
- 添加 session tables 或 migrations。
- 选择 session expiration、refresh、revocation、reconnect 或 connection epoch behavior。
- 选择 token/session carrier behavior。
- 改变 Protobuf envelope fields 或添加 handshake/system messages。
- 添加 WebSocket handshake authentication behavior。
- 添加 route-level authentication implementation。
- 声明 metadata-only `player_id`、`session_id` 或 `connection_id` 足以获得 production permissions。

## 12. 验证

本标准默认 repository verification：

```bash
node tools/vibit check architecture --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check change define-session-persistence-websocket-handshake-decision-gates --json
node tools/vibit check all --json
```

只有 runtime Go code 发生变化时才需要 Go tests。本 decision-gate standard 不要求修改 Go runtime code。

## 13. Agent 规则

Agents 必须：

- 在添加 session persistence、WebSocket handshake authentication、reconnect behavior、connection epoch behavior、token/session carriers 或 session-related protocol changes 前阅读本标准。
- 在未来 handshake standard 被 ratify 前，保持当前 WebSocket transport credential-neutral。
- 在 validation 存在前，保持当前 envelope session fields 为 metadata-only。
- 保持 domain modules 依赖 normalized request identity，而不是 token、session、credential 或 transport internals。
- 当使用 Nakama 和 Pitaya reference patterns 进行规划时，记录 adopted、adapted、deferred 或 rejected。
- 如实记录 verification。

Agents 不得：

- 隐式选择 validation model。
- 隐式添加 session tables。
- 未经 ratification，把 PostgreSQL、memory、Redis-like storage 或任何其他 store 当作已选择的 session store。
- 未经未来 handshake standard，在 WebSocket transport 中解析 tokens、credentials、cookies、`Authorization` headers、WebSocket subprotocols 或 query-parameter proof。
- 仅凭本标准添加 token、credential、authentication result、handshake 或 system-message envelope fields。
- 把 metadata-only `player_id`、`session_id`、`connection_id` 或 `connection_epoch` 当作 authenticated proof。
