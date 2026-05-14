# 玩家账号与 Session 契约标准

状态：Draft v0.2
最后更新：2026-05-14
范围：玩家账号生命周期契约与运行时 session validation 契约
依赖：`docs/player-identity-session-boundary.md`

配套英文文档是 `docs/player-account-session-contracts.md`。英文文件是权威版本。

## 1. 目的

本标准定义 vibit 在实现之前，如何 ratify 玩家账号与 session 契约。

前一个边界标准已经把 player identity、player account、authentication、runtime session validation、transport connection metadata、envelope session metadata 和 request identity context 分开。本标准把这个边界推进成契约 ratification 路径。

目标不是现在实现生产登录系统。目标是让下一批契约足够明确，使 Agent 后续添加它们时，不会在错误层里发明 authentication、token、persistence 或 protocol 行为。

Milestone 状态：

```text
M-005 已于 2026-05-14 完成。
```

这个里程碑完成了第一批 player account 和 runtime session validation surface 的契约 ratification。它没有选择或实现 authentication scheme、token behavior、credential storage、player account persistence、session persistence、Protobuf envelope changes、WebSocket handshake authentication、runtime player account handlers、WebSocket routes、直接 Nakama/Pitaya API 兼容，或 major external framework dependency。

## 2. 参考对齐

Nakama 是 account、user、authentication 和 session 能力覆盖面的主要参考。

Pitaya 是 session binding、route handler、frontend/backend server role 和 realtime server 词汇的主要参考。

vibit 使用这些项目作为参考，而不是治理性 API 形状。

规则：

- 把每个参考 pattern 记录为 adopted、adapted、deferred 或 rejected。
- 没有兼容性 ADR，不得复制 Nakama 或 Pitaya 的 public API。
- 保持 vibit 的 module ownership、contract-first behavior、generated output rules 和 repository checks。
- 不要让 transport handler 拥有 account、authentication、token、credential 或 session persistence 行为。

## 3. 契约族

`M-005` 有两个契约族。

### Player Account Contracts

拥有者：

```text
modules/player
```

目的：

```text
定义稳定 player identity 和持久 player account lifecycle 行为。
```

候选未来契约词汇：

- `CreatePlayerAccount`
- `GetPlayerAccount`
- `LinkPlayerAccountIdentity`
- `DisablePlayerAccount`
- `DeletePlayerAccount`
- `PlayerAccountCreated`
- `PlayerAccountLinked`
- `PlayerAccountDisabled`
- `PlayerAccountDeleted`

这些名称只是候选词汇。只有后续 work item ratify 具体 public contracts 后，它们才成为正式契约。

### Runtime Session Validation Contracts

拥有者：

```text
runtime/internal/app
```

目的：

```text
定义 decoded request 如何在 domain handler 运行前变成 validated request identity。
```

候选未来契约词汇：

- `ValidateSession`
- `BindSessionToConnection`
- `RefreshSession`
- `InvalidateSession`
- `SessionValidated`
- `SessionInvalidated`

这些名称只是候选词汇。只有后续 work item ratify 具体 public contracts 后，它们才成为正式契约。

## 4. 所有权规则

player module 拥有：

- 稳定的 `player_id` 语义。
- Player account lifecycle。
- ratification 之后的 player account public commands、queries、events、errors 和 permissions。
- persistence 被 ratify 后的未来 player account repository interfaces。

player module 不拥有：

- WebSocket connections。
- Protobuf envelope framing。
- Token parsing 或 token signing。
- Credential storage 或 password hashing。
- 单独 ratify 前的 runtime session persistence。
- Inventory、currency、rewards、quests、matches、rooms、chat 或 presence。

Application dispatch 拥有：

- Request identity context。
- Session validation handoff interfaces。
- metadata-only identity 被 validated identity 替换的位置。

Application dispatch 不拥有：

- Player account lifecycle。
- Authentication provider implementation。
- Credential storage。
- WebSocket connection mechanics。
- Generated Protobuf packages。

Transport 拥有：

- 已接受的 WebSocket connections。
- Binary frame IO。
- Transport-local connection metadata。

Transport 不拥有：

- Authentication。
- Player accounts。
- Tokens。
- Credentials。
- Domain permissions。
- Durable session state。

## 5. 参考 Pattern 映射

### Nakama Patterns

| Pattern | Vibit position | Reason |
| --- | --- | --- |
| User/account as first-class backend capability | Adopted | vibit 需要 stable player identity 与 account lifecycle 的持久拥有者。 |
| Multiple authentication methods | Deferred | 支持哪些登录方式会影响安全、存储和 public API shape，必须单独选择。 |
| Session token and refresh token concepts | Deferred | Token format、signing、refresh、expiration 和 revocation 需要单独决策。 |
| Realtime socket bound to authenticated session | Adapted | vibit 会在 domain dispatch 之前验证 request identity，但现在不改变 WebSocket handshake。 |
| Broad social, storage, leaderboard, matchmaker capability surface | Deferred | 对 roadmap 有价值，但不属于 account/session contract ratification。 |

### Pitaya Patterns

| Pattern | Vibit position | Reason |
| --- | --- | --- |
| Session object separate from transport acceptor | Adopted | vibit 保持 WebSocket transport metadata 与 application request identity 分离。 |
| Route handler receives session context | Adapted | vibit handler 通过 application dispatch 接收 `RequestIdentity`，而不是采用 Pitaya API shape。 |
| Frontend/backend server role split | Deferred | 在 modular monolith boundary 稳定前，distributed topology 继续延后。 |
| Groups and push vocabulary | Deferred | 对未来 room、presence 和 broadcast 有用，但不属于 account/session contract ratification。 |
| Direct Pitaya API compatibility | Rejected for now | vibit 定义 agent-readable contracts；兼容性只能通过未来 ADR 引入。 |

## 6. Ratification 顺序

推荐的契约 ratification 顺序是：

1. 标准化 account/session contract rules 和 reference mapping。
2. Ratify 最小 player account semantic contract set。
3. Ratify runtime session validation semantic contracts。
4. 为需要 client/server wire shape 的 account contracts ratify Protobuf request/response messages。
5. 决定是否需要 player account persistence 和 session persistence。
6. 只有相关 contracts 与 ownership rules 注册之后才实现。

这个顺序是有意保守的。它允许 Agent 继续推进，但不会隐式选择 authentication 或 storage 决策。

如果满足以下条件，实际执行顺序可以在 runtime session semantic contracts 之前先 ratify player account Protobuf wire messages：

- Player account semantic contracts 已经 ratify。
- 不添加 runtime player handlers 或 WebSocket routes。
- Protobuf envelope 和 WebSocket handshake 保持不变。
- Authentication、token behavior、credential storage、session persistence 和 player account persistence 仍然 deferred。

第一版 player account wire shape 使用了这个例外，以便在 runtime session contracts 完成前先检查 generated Protobuf output。

## 7. 第一个最小契约集

第一个最小 player account contract set 应优先覆盖 account lifecycle，而不选择登录方式。

已 ratify 的第一批 semantic contracts：

- `CreatePlayerAccount`
- `GetPlayerAccount`
- `PlayerAccountCreated`
- `player_account_errors`
- `player_account_permissions`

建议的 contract-level account fields：

- `player_id`
- `display_name`
- `account_state`
- `created_at`
- `updated_at`

第一个最小 session validation contract set 应优先覆盖 runtime handoff 词汇，而不选择 token format。

已 ratify 的第一批 runtime session validation semantic contracts：

- `ValidateSession`
- `SessionValidated`
- `session_errors`
- `session_permissions`

建议的 contract-level session validation fields：

- `session_id`
- `player_id`
- `connection_id`
- `validation_status`
- `actor_kind`
- `validated_at`

上面列出的 player account contracts 已经有 ratified semantic contracts 和第一版 Protobuf wire messages。

上面列出的 runtime session validation contracts 只是已 ratify 的 semantic contracts。它们由 `runtime/internal/app` 拥有，source 位于：

```text
contracts/runtime/session/
```

它们描述现有 application-owned session validation handoff，并保留当前 metadata-only validator behavior。它们不实现 real authentication、token validation、credential lookup、session persistence、player account lookup、Protobuf envelope changes 或 WebSocket handshake authentication。

已 ratify 的第一批 player account Protobuf messages：

- `CreatePlayerAccountRequest`
- `CreatePlayerAccountResponse`
- `GetPlayerAccountRequest`
- `GetPlayerAccountResponse`
- `PlayerAccountCreated`

已 ratify 的 player account Protobuf source 是：

```text
proto/vibit/player/v1/player.proto
```

生成的 Go Protobuf output 是：

```text
runtime/internal/generated/proto/vibit/player/v1/player.pb.go
```

Database columns、indexes、token claims、generated dispatch shapes、runtime handlers 和 WebSocket handshake fields 仍然必须单独 ratify。

## 8. 必需契约产物

当 player public contract 被 ratify 时，更新：

- `contracts/player/...`
- `.arch/contracts.yaml`
- `modules/player/module.yaml`
- `docs/player-account-session-contracts.md`
- `docs/player-account-session-contracts.zh-CN.md`
- `changes/` 下相关 change spec

当 runtime session validation public contracts 被 ratify 时，更新：

- `contracts/runtime/session/...`
- `.arch/contracts.yaml`
- `.arch/runtime.yaml`
- `docs/player-account-session-contracts.md`
- `docs/player-account-session-contracts.zh-CN.md`
- `changes/` 下相关 change spec

如果添加 Protobuf messages，更新：

- `proto/vibit/<module>/v1/...`
- 只有 generation roots 改变时才更新 `buf.yaml` 或 `buf.gen.yaml`。
- Generated output declarations。
- Protocol alignment checks。

对于第一版 player account wire shape，generation roots 没有改变。`buf generate` 只从新的 player Protobuf source 生成 Go Protobuf output。

## 9. Ask-First 边界

以下情况必须先询问维护者：

- 选择具体登录方式，例如 guest、device、email/password、social login、custom ID 或 external identity providers。
- 选择 JWT、opaque tokens、refresh tokens、signing、expiration 或 revocation behavior。
- 添加 credential storage、password hashing、cryptography、OAuth、OIDC 或 external auth dependencies。
- 添加 player account database schema、migrations 或 session persistence。
- 改变 Protobuf envelope。
- 改变 WebSocket handshake authentication behavior。
- 让 metadata-only `player_id` 或 `session_id` 足以进行生产权限授予。
- 复制 Nakama 或 Pitaya public API shape。

## 10. 验证

实现本标准中的任何契约之前，运行：

```bash
node tools/vibit check contracts --json
node tools/vibit check protocol --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```

添加 public contracts 后，运行：

```bash
node tools/vibit check all --json
```

如果 Go runtime code 改变，也运行：

```bash
cd runtime && go test ./...
```

## 11. 下一方向门

`M-005` 之后，Agent 必须在 `M-006/W-0037` 停下，不能自行选择下一项重大实现方向。

下一方向可能是以下之一，但本标准不替维护者选择：

- Authentication 和 token/session validation design。
- Player account PostgreSQL schema 与 persistence。
- Runtime player account handlers 和 WebSocket route wiring。
- 参考 Nakama 与 Pitaya 能力覆盖后扩展额外 core game backend modules。
- 在继续添加 runtime features 前扩展 generator 和 contract tooling。

这些路径会影响长期 security、data ownership、protocol、runtime 或 roadmap shape，因此必须由 maintainer 决策。
