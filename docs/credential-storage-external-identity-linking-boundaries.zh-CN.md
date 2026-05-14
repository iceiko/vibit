# Credential Storage 与 External Identity Linking 边界

状态：Draft v0.1
最后更新：2026-05-14
范围：Credential storage 边界、external identity linking 边界、player account lifecycle 分离、login-method family 延后、provider subject 延后，以及未来 implementation gates
依赖：`docs/authentication-token-session-validation.md`

配套英文原文是 `docs/credential-storage-external-identity-linking-boundaries.md`。英文文件是权威版本。

## 1. 目的

本标准在 credential storage 和 external identity linking 被实现前，定义 vibit 对它们的边界理解。

目的是让未来 authentication work 对 Agent 可读且边界清晰。未来 Agent 必须能够看出：player account lifecycle storage、credential storage、provider identity linking、token/session behavior 和 runtime request identity validation 是不同职责。

本标准不选择：

- 支持哪一种 login method。
- Password model。
- Credential schema。
- External identity provider。
- OAuth、OIDC、social login、device login、guest login 或 custom ID 行为。
- Password hashing、encryption、signing、key management 或 provider dependencies。
- Credential tables、external identity tables、token tables、session tables 或 migrations。
- Runtime credential lookup、account linking handlers、recovery flows、merge behavior 或 WebSocket routes。

## 2. 必读材料

本标准需要与以下材料一起阅读：

- `docs/authentication-token-session-validation.md`
- `docs/authentication-proof-token-session-contract-dimensions.md`
- `docs/player-identity-session-boundary.md`
- `docs/player-account-session-contracts.md`
- `docs/postgresql-persistence-boundary.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `modules/player/module.yaml`
- `runtime/AGENTS.md`
- `ADR-0019`
- `ADR-0021`
- `ADR-0022`
- `ADR-0023`

参考阅读：

- Nakama authentication concepts：`https://heroiclabs.com/docs/nakama/concepts/authentication/`
- Nakama account/session capability surface：`https://heroiclabs.com/docs/nakama/`
- Pitaya session 与 route-handler vocabulary：`https://pitaya.readthedocs.io/`

Nakama 和 Pitaya 是参考对象。它们不支配 vibit 的 public API shape、credential schema、linking behavior、generated file conventions 或 agent workflow。

## 3. 核心词汇

### Credential

Credential 是 login method 或 identity provider 使用的 secret 或 proof material。

例子包括 passwords、password hashes、device secrets、custom ID secrets、provider secrets、OAuth credentials、OIDC subjects、provider-issued identity material，或未来 service credentials。

规则：

- Credentials 不是 player account lifecycle rows。
- Credentials 不是 runtime sessions。
- Credentials 不是 WebSocket connection metadata。
- Credentials 不是 Protobuf envelope session metadata。
- Credentials 不得存储在 `player_accounts` 或 `player_account_events` 中。
- 在未来标准 ratify schema、dependencies、redaction、tests 和 ownership 前，credential storage 仍未实现。

### External Identity Link

External identity link 把 vibit-owned account identity 映射到 provider namespace 和 provider subject。

例子包括未来 platform accounts、device identities、social identity providers、OIDC issuers、custom identity providers 或 game-platform accounts 的 provider subjects。

规则：

- External identity links 不是 player account lifecycle rows。
- External identity links 本身不等同于 credentials，虽然它们可能与某个 credential 或 login method 相关。
- Provider subject semantics 必须先定义，之后才允许存在 storage。
- Link、unlink、conflict、recovery 和 merge behavior 仍未实现。
- External identity links 不得为了方便添加到 `player_accounts` 或 `player_account_events`。

### Player Account Lifecycle

Player account lifecycle storage 拥有稳定的 account identity state。

当前 lifecycle tables：

```text
player_accounts
player_account_events
```

规则：

- 这些 tables 仍然只存储 account lifecycle。
- 它们可以记录 player account PostgreSQL schema boundary 已 ratify 的 lifecycle state，例如 created、disabled 或 deleted account status。
- 它们不得存储 credential material、provider subjects、token state、runtime sessions、WebSocket connection state 或 request validation results。

### Login Method Family

Login method family 是未来产生 authentication proof 的方式类别。

Deferred families 包括：

- guest 或 anonymous-login family
- device-login family
- email/password family
- custom ID family
- social-login family
- OAuth family
- OIDC family
- external identity-provider family
- service-auth family

规则：

- 本标准不选择任何 family。
- 列出 family 只是 capability coverage，不是 implementation permission。
- 每个被选择的 family 在实现前都必须有未来 contract、schema boundary、dependency review、error model、permission model 和 verification path。

## 4. Ownership 边界

### Player Module

Owner：

```text
modules/player
runtime/internal/modules/player
```

拥有：

- 稳定 `player_id`
- player account lifecycle contracts
- player account lifecycle repository interfaces
- 已 ratify 的 player account lifecycle persistence boundaries

不拥有：

- credential storage
- password hashing
- external identity provider subjects
- token issuance 或 validation
- runtime session persistence
- WebSocket connection binding
- request validation results

### 未来 Credential Boundary

Owner 状态：

```text
planned, not implemented
```

未来工作必须定义：

- Credential owner module 或 runtime subsystem。
- Credential record lifecycle。
- Secret material 和 non-secret metadata。
- Hashing、encryption、signing 或 provider dependency adoption。
- Redaction 和 logging rules。
- Access permissions。
- Failure classes 和 retryability。
- Migration ownership。
- Tests 和 repository checks。

### 未来 External Identity Boundary

Owner 状态：

```text
planned, not implemented
```

未来工作必须定义：

- Provider namespace semantics。
- Provider subject semantics。
- Link 和 unlink authority。
- Duplicate provider-subject behavior。
- Provider subject 映射到已有 account 时的 conflict behavior。
- Account recovery behavior。
- Account merge behavior，如果支持。
- Provider metadata retention 和 redaction rules。
- Audit events 和 client-visible events。

### Runtime Session Validation

Owner：

```text
runtime/internal/app
```

规则：

- Runtime session validation 只能在未来 contracts ratify 后消费 authentication proof 或 session validation results。
- 除非未来边界明确授予该职责，否则它不得直接查询 credential stores 或 external identity stores。
- Domain modules 接收 normalized request identity；它们不解析 credentials 或 provider subjects。

## 5. Deferred Login-Method Coverage

Nakama 展示了生产级 game backend 通常会支持多个 authentication methods 和 account/session concepts。

vibit 采纳这一点作为 capability coverage，但延后具体 login-method selection。

| Capability family | Vibit position | Reason |
| --- | --- | --- |
| Device-style login | Deferred | 需要 device identifier semantics、replay controls、secret treatment、account recovery behavior 和 abuse controls。 |
| Email/password login | Deferred | 需要 password hashing、password reset、rate limiting、credential storage、secret redaction 和 recovery flows。 |
| Custom ID login | Deferred | 需要 issuer semantics、collision behavior、account linking rules 和 trusted caller boundaries。 |
| Social/provider login | Deferred | 需要 provider namespace、subject semantics、provider metadata、conflict behavior 和 dependency adoption。 |
| OAuth/OIDC-style login | Deferred | 需要 provider trust、issuer/audience validation、key management、token validation、refresh behavior 和 dependency adoption。 |
| Guest/anonymous login | Deferred | 需要显式 anonymous actor rules、account upgrade behavior、data ownership 和 anti-abuse posture。 |
| Session token 和 refresh token concepts | Deferred | 需要 token format、issuer、verifier、expiration、revocation、rotation、storage 和 replay behavior。 |
| Direct Nakama API compatibility | Rejected for now | vibit 定义 agent-native contract surface；兼容性需要未来 ADR。 |

Pitaya 展示了 connection acceptors、sessions、route handlers 和 server roles 之间有价值的分离。

vibit 这样适配该 vocabulary：

| Pitaya pattern | Vibit position | Reason |
| --- | --- | --- |
| Session object separate from acceptor | Adopted as vocabulary | Transport connections 必须与 application request identity 和 future runtime sessions 分离。 |
| Handler receives session context | Adapted | vibit handlers 接收 application-owned `RequestIdentity`，不是 transport-owned session object。 |
| Session binding | Adopted as vocabulary | Binding 是有用词汇，但必须来自 ratified validation results，而不是 raw metadata。 |
| Frontend/backend split | Deferred | Distributed topology 不属于当前 credential/linking boundary。 |
| Direct Pitaya API compatibility | Rejected for now | vibit 可以学习 architecture vocabulary，但不复制 public APIs。 |

## 6. 禁止的捷径

Agents 不得：

- 向 `player_accounts` 添加 credential columns。
- 向 `player_accounts` 添加 provider subject columns。
- 向 `player_accounts` 添加 token 或 session columns。
- 向 `player_account_events` 添加 credential、provider subject、token、session、WebSocket state 或 request-validation rows。
- 未经未来 schema boundary 添加 credential tables 或 external identity tables。
- 仅凭本标准添加 password hashing、OAuth、OIDC、JWT、provider SDK、cryptography 或 key-management dependencies。
- 仅凭本标准添加 runtime credential lookup 或 external identity lookup。
- 仅凭本标准添加 account linking、unlinking、recovery 或 merge behavior。
- 从 reference project list 推断 login method。
- 把 direct Nakama 或 Pitaya API compatibility 当成未说明的目标。

## 7. 未来 Credential Storage 产物门槛

实现 credential storage 前，未来工作必须新增或更新：

- `changes/` 下 change spec。
- 当选择影响长期架构、provider dependency posture、generated file conventions 或 security boundaries 时新增 ADR。
- Manifests 中的 credential owner declaration。
- Login method contract。
- Credential schema boundary。
- 如果存储数据，则新增 migration source。
- Secret-handling rules。
- Redaction 和 logging rules。
- Error catalog 与 failure-class mapping。
- Permission catalog。
- 常见 success 和 failure modes 的 tests。
- 证明 credentials 不存储在 player lifecycle tables 中的 negative tests。
- 当规则可以静态执行时新增 repository checks。
- 英文文档和简体中文译本。

未来 schema boundary 必须明确回答：

- 哪些 records 包含 secret material。
- 哪些 records 包含 non-secret lookup metadata。
- 哪些 identifiers 可以安全记录到日志。
- 哪些 fields 是 unique。
- 哪些 fields 是 mutable。
- 哪些 account lifecycle states 会阻止 login。
- 哪些 operations 需要与 player account lifecycle changes 一起事务执行。

## 8. 未来 External Identity Linking 产物门槛

实现 external identity linking 前，未来工作必须新增或更新：

- `changes/` 下 change spec。
- 当 provider semantics、merge behavior、recovery behavior 或 direct API compatibility 影响长期架构时新增 ADR。
- Manifests 中的 external identity owner declaration。
- Provider namespace contract。
- Provider subject contract。
- Link command contract。
- 如果支持 unlinking，则新增 unlink command contract。
- Conflict、duplicate、recovery 和 merge semantics。
- External identity schema boundary。
- 如果存储数据，则新增 migration source。
- Audit event catalog。
- Error catalog 和 retryability rules。
- Permission catalog。
- Link、unlink、duplicate、conflict、disabled-account、deleted-account 和 recovery behavior 的 tests。
- 当规则可以静态执行时新增 repository checks。
- 英文文档和简体中文译本。

未来 linking boundary 必须明确回答：

- Provider subject IDs 是 globally unique 还是 provider-scoped。
- Provider metadata 是 retained、normalized、redacted 还是 discarded。
- 一个 account 是否可以持有多个 provider links。
- 一个 provider subject 是否可能映射到多个 vibit accounts。
- Linked accounts 是否可以 merge。
- Links 是否可以在 accounts 之间移动。
- 哪些 events 只属于 security/audit，哪些可以 client-visible。

## 9. 与 Token 和 Session Work 的关系

Credential storage 和 external identity linking 本身不定义 runtime session behavior。

规则：

- Credential 只能通过未来 ratified login method 产生 authentication proof。
- External identity link 只能通过未来 ratified provider validation path 识别 account。
- Token 只能在 token format、issuance、expiration、revocation、storage 和 validation behavior 被 ratify 后签发。
- Runtime session 只能在 session persistence boundary 被 ratify 后持久化。
- WebSocket connection 只能在 request-level、first-message、handshake-level、every-request 或 hybrid validation behavior 被 ratify 后绑定到 identity。

## 10. Ask-First 边界

在以下操作前询问 maintainer：

- 选择支持的 login method。
- 选择 password model。
- 选择 credential table shape。
- 选择 external identity table shape。
- 选择 provider namespace 或 provider subject semantics。
- 选择 OAuth、OIDC、social login、provider SDK、password hashing、encryption、signing 或 key-management dependencies。
- 添加 credential storage、external identity storage、token storage 或 session storage。
- 添加 runtime credential lookup、account linking handlers、account recovery flows 或 account merge behavior。
- 改变 player account lifecycle table shape。
- 复制 Nakama 或 Pitaya public API shape。

## 11. 验证

本标准默认 repository verification：

```bash
node tools/vibit check architecture --json
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check memory --json
node tools/vibit check module player --json
node tools/vibit check change define-credential-storage-external-identity-linking-boundaries --json
node tools/vibit check all --json
```

只有 runtime Go code 发生变化时才需要 Go tests。本边界标准不要求修改 Go runtime code。

## 12. Agent 规则

Agents 必须：

- 在添加 credential storage、external identity linking、login methods 或 provider-related behavior 前阅读本标准。
- 保持 player account lifecycle tables 仅作为 account lifecycle storage。
- 在未来选择被 ratify 前，把 login-method family lists 保持为 deferred capability coverage。
- 使用 Nakama 和 Pitaya reference patterns 规划前，记录它们是 adopted、adapted、deferred 还是 rejected。
- 如实记录 verification。

Agents 不得：

- 在 player account lifecycle tables 中存储 credentials、provider subjects、tokens、runtime sessions、WebSocket connection state 或 request validation results。
- 隐式选择 login method。
- 隐式添加 credential 或 provider dependencies。
- 仅凭本边界标准实现 account linking、unlinking、recovery 或 merge behavior。
- 把 metadata-only `player_id` 或 `session_id` 当作 authenticated proof。
