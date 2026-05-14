# 认证、Token 与会话校验设计标准

状态：Draft v0.1
最后更新：2026-05-14
范围：认证证明、token 行为、runtime session validation、credential 边界、external identity 边界、session persistence 边界、request identity trust、Protobuf envelope 交互和 WebSocket handshake 交互
权威决策：`ADR-0023`

本文件是 `docs/authentication-token-session-validation.md` 的简体中文译本。英文版本是权威版本。

## 1. 目的

vibit 现在已经有 player account lifecycle persistence 和 application-owned session validation hook，但它仍然没有 production authentication。

本标准在未来 Agent 实现 login、token validation、credential storage、external identity linking、runtime session persistence、WebSocket handshake authentication、player account handlers 或 player WebSocket routes 之前定义设计边界。

目标是让 authentication 和 session 行为对 Agent 可读、可验证、可安全扩展。未来 Agent 应该能明确识别：proof 在哪里产生，在哪里校验，在哪里存储，在哪里变成 request identity，以及 domain permissions 在哪里消费它。

本标准不选择：

- 具体 login method。
- JWT、opaque tokens、refresh tokens、signing、expiration、revocation 或 token storage 行为。
- Credential storage、password hashing、OAuth、OIDC、social login、device login、guest login 或 custom ID 行为。
- External identity linking tables 或 provider dependencies。
- Runtime session persistence。
- Protobuf envelope changes。
- WebSocket handshake authentication behavior。
- Runtime player account handlers 或 WebSocket routes。

## 2. 必读内容

本标准应与以下文件一起阅读：

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
- `modules/player/module.yaml`
- `runtime/AGENTS.md`
- `ADR-0019`
- `ADR-0021`
- `ADR-0022`
- `ADR-0023`

参考阅读：

- Nakama documentation：`https://heroiclabs.com/docs/nakama/`
- Nakama authentication and sessions concepts：`https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama realtime multiplayer and socket concepts：`https://heroiclabs.com/docs/nakama/concepts/multiplayer/`
- Pitaya documentation：`https://pitaya.readthedocs.io/`

Nakama 和 Pitaya 是参考，不支配 vibit 的 public API shape、runtime module layout、generated file conventions 或 agent workflow。

## 3. 核心词汇

以下术语必须保持分离。

### Authentication Proof

Authentication proof 是验证某个 actor 可以绑定到 player identity、service identity 或未来 administrative identity 后得到的结果。

规则：

- Authentication proof 不等同于 token string。
- Authentication proof 不等同于 player account row。
- Authentication proof 必须在 domain handlers 接收请求前变成 machine-readable validation result。
- Authentication proof 不得从 client-supplied `player_id`、`session_id` 或 `connection_id` metadata 推断出来。

### Login Method

Login method 是 client 获得 authentication proof 或 session credential 的方式。

例子包括 guest login、device login、email/password login、custom ID login、social login、OAuth、OIDC 和 external identity-provider login。

规则：

- 本标准不 ratify 任何 login method。
- 每个 login method 在实现前都必须有未来 contract、storage boundary、error model 和 verification path。
- Login method 不得实现在 WebSocket transport、Protobuf adapters、inventory handlers 或通用 player account persistence 中。

### Token

Token 是 client 或 service 提交给服务器进行 validation 的 credential-like artifact。

例子包括 access tokens、session tokens、refresh tokens、opaque tokens 和 signed tokens。

规则：

- 本标准不 ratify 任何 token format。
- Token parsing、signing、verification、refresh、revocation、rotation、expiration、storage 和 replay handling 都是独立设计维度。
- Token behavior 必须先声明再实现。
- Token validation 不得隐藏在 domain modules、Protobuf payload bridges 或 WebSocket frame handlers 中。

### Credential

Credential 是 login method 或 identity provider 使用的 secret 或 proof material。

例子包括 passwords、password hashes、device secrets、provider secrets、OAuth credentials、OIDC subjects 和 provider-issued identity material。

规则：

- 本标准不 ratify credential storage。
- Credentials 不得存储在 `player_accounts` 或 `player_account_events` 中。
- Credential storage 必须有单独 schema boundary、secret-handling rule、必要时的 dependency adoption record 和 verification path。

### External Identity Link

External identity link 把未来 vibit player account 映射到 identity provider subject。

规则：

- 本标准不 ratify external identity linking。
- Provider names、provider subject IDs、provider metadata、linking conflicts、unlinking、account recovery 和 merge behavior 都需要未来设计。
- External identity links 不得为了方便而加入 player account lifecycle tables。

### Runtime Session

Runtime session 是 authentication/session validation 被 ratify 后使用的服务端 validation 和 binding state。

规则：

- Runtime sessions 不是 player accounts。
- Runtime sessions 不是 WebSocket connection IDs。
- Runtime sessions 未来可以绑定 `session_id`、`player_id`、actor kind、claims、expiration、connection metadata 和 revocation state，但这里不 ratify 任何 storage model。
- Session validation 归 application 所有，并在 production-sensitive domain dispatch 之前运行。

### Request Identity

Request identity 是传给 command 和 query handlers 的 application-facing identity context。

当前 owner：

```text
runtime/internal/app
```

规则：

- `RequestIdentity` 是 domain permissions 的 handoff type。
- Domain modules 消费 request identity；它们不校验 tokens 或 credentials。
- 当前 metadata-only request identity 不是 authenticated proof。

### WebSocket Handshake Authentication

WebSocket handshake authentication 指在 connection establishment 之前或过程中校验 actor。

规则：

- 本标准不 ratify WebSocket handshake authentication。
- 在未来 handshake standard 存在前，WebSocket transport 必须保持 credential-neutral。
- 任何未来 handshake 设计都必须保持 transport mechanics 与 authentication proof、session binding semantics 的分离。

## 4. 信任状态

vibit 使用以下 trust states 作为设计词汇。

| State | 含义 | Production authority |
| --- | --- | --- |
| `anonymous` | 没有 actor proof。 | 不得授予 player-owned 或 privileged permissions。 |
| `metadata_only` | Identity text 来自 envelope 或 transport metadata，并经当前 metadata-only path normalize。 | 本身不得授予 production permissions。 |
| `authentication_proven` | 未来已 ratify 的 authenticator 校验了 login method、token 或 trusted service proof。 | 只有 validating component 和 contracts 被 ratify 后才能使用。 |
| `session_validated` | 未来已 ratify 的 session validator 把 request identity 绑定到有效 logical session。 | Validation 成功后可被 module permission policies 使用。 |
| `service_validated` | 未来已 ratify 的 service-auth path 校验了 non-player service authority。 | 只能被显式建模的 service permissions 使用。 |

当前 runtime 只有 metadata-only 行为。`MetadataOnlySessionValidator` 保留现有 proof-slice request flow；它不认证 clients。

## 5. Ownership Model

### Authentication Boundary

Owner 状态：

```text
planned, not implemented
```

Ratify 后的职责：

- 校验 login-specific proof 或 token/session credentials。
- 产生 machine-readable authentication results。
- 将 authentication failures 映射到 registered errors。
- 避免把 credential implementation details 泄漏到 domain modules。

不得：

- 除非 contract 要求，否则不得把存储 account lifecycle rows 作为副作用。
- 不得把 token parsing 隐藏在 transport、protocol、player account repository 或 inventory code 中。
- 不得把 metadata-only identity 当作 proof。

### Application Session Validation

Owner：

```text
runtime/internal/app
```

职责：

- 在 protocol decoding 后、domain dispatch 前调用 session validation。
- 把 validation results 转换成 `RequestIdentity`。
- 仅把当前 metadata-only behavior 作为 non-authenticated bootstrap path 保留。
- 实现时把 authentication provider details 放在显式 interfaces 后面。

不得：

- Import WebSocket transport libraries。
- Import generated Protobuf packages。
- 存储 credential records。
- 拥有 player account lifecycle persistence。

### Player Module

Owner：

```text
modules/player
runtime/internal/modules/player
```

职责：

- 拥有稳定 `player_id` 和 player account lifecycle。
- 拥有 player account contracts 和 lifecycle persistence。
- 作为未来 authentication 和 linking flows 使用的 durable account owner。

不得：

- 在 account lifecycle tables 中存储 credentials、token state、refresh tokens、runtime sessions、WebSocket connection state 或 request validation results。
- 在 account repositories 中 validate tokens。
- 未经单独 ratification 添加 runtime player handlers 或 WebSocket routes。

### WebSocket Transport

Owner：

```text
runtime/internal/platform/transport/ws
```

职责：

- 接受 WebSocket connections。
- 读取和写入 binary frames。
- 提供 transport-local metadata，例如 connection IDs。
- 把 opaque frame bytes 交给 protocol/application composition。

不得：

- 解析 credentials 或 tokens。
- 认证 players。
- 把 player accounts 绑定到 connections。
- 拥有 session persistence。
- 执行 domain permissions。

### Protobuf Protocol Adapter

Owner：

```text
runtime/internal/platform/protocol/protobuf
```

职责：

- Decode 和 encode 当前 envelope。
- 保留现有 session metadata fields。
- 把 envelope session metadata 转换成 application handoff types。
- 当相关 errors 被注册后，把 validation 或 protocol errors 映射成 error envelopes。

不得：

- 未经 protocol change spec 和 ADR 向 envelope 添加 credential 或 token fields。
- 把 envelope `player_id` 或 `session_id` 当作 proof。
- 拥有 long-lived session state。
- 执行 permission decisions。

### Domain Modules

例子：

```text
inventory, currency, reward, quest, match
```

职责：

- 使用 `RequestIdentity` 执行 module-owned permissions 和 invariants。
- 对 production-sensitive operations，把 `metadata_only` 视为 unauthenticated。

不得：

- Validate tokens。
- Parse credentials。
- 直接查询 credential 或 session stores。
- 创建或链接 player accounts，除非它们是 player module 且行为已 ratify。

## 6. Request Flow

当前 non-authenticated bootstrap flow：

```text
websocket frame
-> protobuf envelope
-> route request
-> metadata-only session validator
-> metadata-only request identity
-> application dispatch
-> domain handler
```

未来 production-sensitive flow：

```text
websocket frame
-> protobuf envelope
-> route request
-> ratified authentication/session validation boundary
-> validated request identity
-> application dispatch
-> module permission policy
-> domain handler
```

规则：

- Production-sensitive domain handlers 必须在信任 player-owned 或 privileged actions 前收到 validated request identity。
- Request-level validation 是 domain dispatch 前所必需的 application handoff。
- Handshake-level validation、first-message validation、every-request validation 和 hybrid validation 都是未来设计选择。
- 未来 handshake model 不得移除 application-owned validation handoff。

## 7. Token 设计维度

未来 token work 在实现前必须决定并记录：

- Token format：opaque、signed、structured 或其他形式。
- Token issuer 和 verifier ownership。
- Subject semantics：player、service、admin、guest 或其他 actor kind。
- Audience 和 route scope。
- Session binding 和 connection binding。
- Expiration 和 clock-skew 行为。
- Refresh 行为。
- Revocation 行为。
- Rotation 行为。
- Replay detection 行为。
- Storage requirements。
- Secret 和 key management。
- Error codes 和 retryability。
- Logging 和 redaction。
- Test fixtures 和 negative tests。

本标准有意不决定这些维度。

这些未来选择所需的 semantic vocabulary 已在 `docs/authentication-proof-token-session-contract-dimensions.md` 中 ratify。未来 token/session contracts 在实现前，必须把 actor kinds、validation statuses、proof statuses、failure classes、retryability 和 request identity handoff semantics 映射回该标准。

## 8. Credential 与 External Identity 边界

Credential 和 external identity storage 需要单独边界，之后才允许出现 schema 或 code。

该边界现在由 `docs/credential-storage-external-identity-linking-boundaries.md` 定义。

该边界标准定义 ownership separation、deferred login-method families、provider-subject deferral、future artifact gates 和 forbidden shortcuts。它不实现 credential storage 或 external identity linking。

未来 credential work 必须定义：

- 支持哪些 login methods。
- 存储哪些 secret material。
- 采用哪些 hashing、encryption 或 provider dependencies。
- 哪些 tables 拥有 credential records。
- Credential rows 如何关联 player accounts。
- Disabled 或 deleted accounts 如何影响 login。
- 如何记录 failures 且不泄漏 secrets。
- Tests 如何证明常见 failure modes。

未来 external identity work 必须定义：

- Provider namespace 和 subject semantics。
- Link 和 unlink 行为。
- Duplicate 和 conflict 行为。
- Account recovery 和 merge 行为。
- Provider metadata retention rules。
- Audit/event requirements。

在这些标准存在前：

- `player_accounts` 仍然只存储 account lifecycle。
- `player_account_events` 仍然只存储 lifecycle event。
- 不得为了方便添加 credential、provider subject、access token、refresh token 或 session row。

在该边界标准存在后，同一规则仍然适用，直到未来 implementation standard ratify 具体 schema、dependencies、contracts 和 verification：

- `player_accounts` 仍然不包含 credential 和 provider subject。
- `player_account_events` 仍然不包含 credential、provider subject、token、session、WebSocket state 和 request-validation。
- Login-method family coverage 仍然只是 deferred capability coverage，不是 implementation permission。

## 9. Session Persistence 边界

Session persistence 继续 deferred。

Session persistence 与 WebSocket handshake decision gates 现在由 `docs/session-persistence-websocket-handshake-decision-gates.md` 定义。

该标准把 request-level validation、first-message validation、handshake-level validation、every-request validation 和 hybrid validation 分离为未来选择。它不选择生产模型、session store、token/session carrier、Protobuf envelope change、handshake/system message 或 route-level authentication behavior。

未来 session persistence work 必须决定：

- Sessions 是否 server-side persisted。
- Session store 是 PostgreSQL、memory、Redis-like 还是其他被 ratify 的 store。
- Session ID generation 和 lookup semantics。
- Expiration 和 refresh semantics。
- Revocation 和 logout semantics。
- Reconnect 和 connection epoch semantics。
- 与 WebSocket connection lifecycle 的关系。
- 与 token validation 的关系。
- Cleanup 和 migration behavior。
- Opt-in live verification requirements。

在 ratify 前，runtime 必须保持 `session_id` 和 `player_id` 为 metadata-only，除非未来 validator 显式校验它们。

在 decision-gate standard 存在后，同一规则仍然适用，直到未来 implementation standard ratify 具体 validation model、storage model、protocol impact 和 verification path：

- WebSocket transport 保持 credential-neutral。
- 当前 envelope session fields 仍然只是 metadata carriers。
- `connection_epoch` 仍然只是 metadata。
- Session persistence 仍未实现。
- WebSocket handshake authentication 仍未实现。

## 10. Protobuf Envelope 与 WebSocket Handshake 交互

当前 Protobuf envelope 已经有：

```text
Session.connection_id
Session.session_id
Session.player_id
Session.connection_epoch
```

本标准不改变这些字段。

规则：

- 现有 envelope session fields 在 validation 存在前仍是 metadata carriers。
- 添加 credential fields、token fields、authentication result fields 或新的 handshake/system messages 需要 protocol change spec 和 ADR。
- 未来 token carrier 可以是 envelope metadata、system message、first request payload、WebSocket subprotocol/header pattern 或其他被 ratify 的设计。本标准不选择其中任何一种。
- 除非未来 handshake standard 赋予 WebSocket transport 狭窄的 transport-level responsibility，否则 WebSocket transport 不得 parse 或 validate credentials。
- 即使未来采用 handshake authentication，application dispatch 仍必须收到可供 domain permissions 检查的 normalized request identity。
- 选择 request-level、first-message、handshake-level、every-request 或 hybrid validation 需要未来 ask-first decision。

## 11. 参考模式映射

### Nakama Patterns

| Pattern | Vibit position | Reason |
| --- | --- | --- |
| Account/user as first-class backend capability | Adopted | vibit 已经把 player identity 和 account lifecycle 视为 player module 的 first-class concerns。 |
| Multiple authentication methods | Deferred | Login methods 会影响 contracts、storage、dependencies、security posture 和 public API shape。 |
| Session token concept | Deferred | Token format、issuer、verifier、expiration 和 validation behavior 需要单独 ratification。 |
| Refresh token concept | Deferred | Refresh、rotation、revocation、storage 和 replay handling 需要单独 ratification。 |
| Realtime socket bound to authenticated session | Adapted | vibit 要求 production-sensitive domain dispatch 前有 request identity validation，但 WebSocket handshake binding 尚未决定。 |
| Session expiration and revocation vocabulary | Adopted as design dimensions | 未来 token/session work 必须考虑这些维度，但这里不实现行为。 |
| Direct Nakama public API compatibility | Rejected for now | vibit 定义 agent-native contract surface；兼容性需要未来 ADR。 |

### Pitaya Patterns

| Pattern | Vibit position | Reason |
| --- | --- | --- |
| Session object separate from acceptor/transport | Adopted | vibit 将 transport connection metadata 与 application request identity、future runtime sessions 分开。 |
| Handler receives session context | Adapted | vibit domain handlers 接收 application-owned `RequestIdentity`，不是 Pitaya API session objects。 |
| Session binding vocabulary | Adopted as vocabulary | Binding 是有用词汇，但 vibit 会通过 ratified validation results 和 manifests 进行绑定。 |
| Frontend/backend server role split | Deferred | Distributed topology 在 single-process boundaries 和 checks 稳定前继续 deferred。 |
| Realtime session context | Adapted | vibit 保持 realtime connection metadata、session validation 和 domain permissions 分离。 |
| Direct Pitaya public API compatibility | Rejected for now | vibit 可以学习 Pitaya architecture vocabulary，但不复制其 API surface。 |

## 12. 未来必需产物

在实现 production authentication 前，未来工作必须按相关范围新增或更新：

- `changes/` 下 change spec。
- 当决策影响长期架构时新增 ADR。
- `contracts/` 下 contract source files。
- 当 public commands、queries、events、errors 或 permissions 变化时更新 `.arch/contracts.yaml`。
- `.arch/runtime.yaml` 中的 runtime ownership 和 implementation state。
- 当新增 foundational dependencies 时更新 `.arch/dependencies.yaml` 和 dependency adoption records。
- Module manifests 和 module guides。
- 当新增数据时定义 database schema boundary 和 migrations。
- 只有 protocol impact 被 ratify 后才添加 Protobuf source 和 generated output。
- Runtime tests 和 repository checks。
- 英文文档和简体中文译本。

未来 authentication proof 或 token/session validation contract artifacts 在定义 actor kinds、validation statuses、proof statuses、failure classes、retryability、request identity handoff、error semantics、permission semantics 或 events 时，也必须引用 `docs/authentication-proof-token-session-contract-dimensions.md`。

## 13. Ask-First 边界

在以下操作前询问 maintainer：

- 选择 guest、device、email/password、custom ID、social login、OAuth、OIDC 或其他 external identity provider。
- 选择 JWT、opaque tokens、refresh tokens、signing、expiration、revocation、rotation 或 token storage behavior。
- 添加 credential storage、password hashing、cryptography、OAuth、OIDC、external identity 或 session-store dependencies。
- 添加 authentication runtime code、token parsing、credential lookup、external identity linking 或 session persistence。
- 改变 Protobuf envelope behavior。
- 改变 WebSocket handshake authentication behavior。
- 添加 runtime player account handlers 或 WebSocket routes。
- 声明 metadata-only `player_id` 或 `session_id` 足以获得 production permissions。
- 复制 Nakama 或 Pitaya public API shape。

## 14. M-011 剩余工作队列

Authentication design milestone 剩余部分应按 bounded steps 推进：

1. 添加 architecture checks，用于强制 authentication/token/session design boundary。
2. Ratify authentication proof 与 token/session validation 的 semantic contract dimensions，但不选择具体 token format。已由 `changes/2026-05-14-ratify-authentication-proof-token-session-contract-dimensions/` 完成。
3. 定义 credential storage 和 external identity linking 边界，但不添加 schema 或 dependencies。已由 `changes/2026-05-14-define-credential-storage-external-identity-linking-boundaries/` 完成。
4. 定义 session persistence 和 WebSocket handshake decision gates，但不实现任何路径。已由 `changes/2026-05-14-define-session-persistence-websocket-handshake-decision-gates/` 完成。
5. 关闭该 milestone，或为第一个实现方向创建 confirmation gate。

每一步都必须保持 metadata-only identity 为 non-authenticated，直到 real validator 被单独 ratify 并实现。

## 15. 验证

本标准默认 repository verification：

```bash
node tools/vibit check architecture --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check all --json
```

`node tools/vibit check runtime --json` 包含 `runtime.authentication_token_session_boundary`。该规则会静态检查 standard/ADR references、implementation-status markers、metadata-only validator markers、Protobuf source 和 generated-output boundaries、WebSocket transport boundary、player repository boundary 以及 player account migration boundary，并且不要求 live external services。

只有 runtime Go code 发生变化时才需要 Go tests。本设计标准不要求修改 Go runtime code。

## 16. Agent 规则

Agents 必须：

- 在添加 authentication、token、credential、external identity、session persistence、request identity、WebSocket handshake 或 runtime player route 行为前阅读本标准。
- 保持 authentication proof、login methods、tokens、credentials、external identity links、runtime sessions、request identity、transport connections 和 player account lifecycle 的分离。
- 在实现 public behavior 前使用 change specs、contracts、manifests、ADRs 和 checks。
- 把当前 metadata-only identity 视为 unauthenticated。
- 当使用 Nakama 或 Pitaya pattern 进行规划时，记录 adopted、adapted、deferred 或 rejected。

Agents 不得：

- 在 WebSocket transport、Protobuf bridges、player repositories 或 domain handlers 中添加 authentication shortcuts。
- 把 envelope `player_id` 或 `session_id` 当作 proof。
- 在 player account lifecycle tables 中存储 credentials、tokens、provider subjects 或 sessions。
- 未经 ratification 添加 token、credential、OAuth、OIDC、password hashing 或 session-store dependencies。
- 在 design-only work 中改变 Protobuf envelope 或 WebSocket handshake behavior。
