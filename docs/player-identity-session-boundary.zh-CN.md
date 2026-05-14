# 玩家身份与会话边界标准

状态：Draft v0.1
最后更新：2026-05-14
范围：玩家身份、账号、认证、WebSocket 会话和请求身份上下文边界
权威决策：`ADR-0021`

本文件是 `docs/player-identity-session-boundary.md` 的简体中文译本。英文版本是权威版本。

## 1. 目的

vibit 的第一版 durable inventory runtime 已经能够处理 player-scoped commands 和 queries，但当前的 `player_id` 仍然只是请求数据，而不是经过认证的 runtime identity。

本标准在更多模块依赖临时 player context 之前定义边界。它防止未来 agent 把 identity、authentication、session validation 或 player account 行为放到最容易改的层里。

本步骤的目标是边界清晰，而不是生产级登录系统。

## 2. 外部参考对齐

Nakama 是账号、认证、用户和 session 相关产品能力面的主要参考。

Pitaya 是 Go 游戏服务器中 acceptor、session binding、route handler 和 server push 相关架构词汇的主要参考。

vibit 采用的是关注点拆分，而不是它们的公开 API：

- Account 和 authentication 概念不得隐藏在 transport handler 中。
- Connection session 和 route dispatch 必须保持分离。
- Realtime connection metadata 不等同于 durable player account ownership。
- 未来如果要兼容 Nakama、Pitaya 或其他项目的 API，需要单独 ADR。

## 3. 身份词汇

以下术语彼此独立。

### Player Identity

`player_id` 是游戏模块用来定位 player-owned aggregate 的稳定 domain identity。

Owner：

```text
player module
```

规则：

- `player_id` 是 domain identity，不是 WebSocket connection id。
- inventory 等 domain modules 可以引用 `player_id`。
- 除 `player` 外，其他 domain modules 不得创建、认证、合并、删除或重命名 players。
- 在 player module 存在之前，`player_id` 仍是带有明确边界说明的 external reference。

### Player Account

Player account 是玩家的 durable account record 和 lifecycle。

Owner：

```text
player module
```

规则：

- Account creation、lookup、linking、disabling、deletion 和 recovery 属于 player module。
- Account data 不得存储在 inventory tables 或 transport/session structures 中。
- 添加 player account database migrations 需要单独 change spec 和 maintainer confirmation。

### Authentication

Authentication 证明某个 actor 可以绑定到 player identity 或 privileged service identity。

Owner：

```text
authentication boundary, implemented through the player/session capability when ratified
```

规则：

- 本标准不选择具体认证方案。
- 未经单独决策，不得引入 JWT、OAuth、OIDC、password hashing、password storage、guest login、device login、social login 或 external identity providers。
- Authentication 必须为下游 session binding 产生 machine-readable result，而不是在 domain handlers 中解释。

### Runtime Session

Runtime session 是 session validation 存在后，服务端与已接受 client connection 关联的 logical session。

Owner：

```text
runtime session boundary
```

规则：

- Runtime session 可以把 connection 绑定到 `session_id`、`player_id` 和未来 authorization claims。
- Runtime session 不是 player account。
- Runtime session 不归 WebSocket transport adapter 所有。
- Session storage、expiration、reconnect behavior 和 token validation 推迟到单独确认后实现。

### Transport Connection

`connection_id` 是 transport-local connection metadata。

Owner：

```text
runtime/internal/platform/transport/ws
```

规则：

- `connection_id` 可以在 reconnect 后变化。
- Transport 可以暴露 connection metadata，但不得把它认证为 player。
- Transport 不得解析 credentials、domain payloads 或 player account data。

### Request Identity Context

Request identity context 是附加到 dispatched command 或 query 的 application-facing identity data。

Owner：

```text
runtime/internal/app
```

规则：

- Request identity context 在 protocol decoding 和 session validation 之后构建。
- Domain handlers 通过 vibit-owned application request types 接收已经 normalize 的 identity context。
- 模块可以在执行权限检查时比较 request identity 和 requested target data。
- 模块不得验证 transport credentials 或 token formats。

## 4. 分层所有权

### WebSocket Transport

Owner：

```text
runtime/internal/platform/transport/ws
```

职责：

- 接受 WebSocket connections。
- 读取和写入 binary frames。
- 提供 transport-local connection metadata。
- 把 opaque frame bytes 交给注入的 protocol/application composition。

不得：

- 认证 players。
- 解析 tokens、credentials 或 account payloads。
- 创建或验证 player accounts。
- 把 `connection_id` 当作 durable identity。
- 执行 domain permissions。

### Protobuf Protocol Adapter

Owner：

```text
runtime/internal/platform/protocol/protobuf
```

职责：

- Decode 和 encode 现有 Protobuf envelope。
- 保留 envelope 中已经存在的 session metadata fields。
- 把 envelope metadata 转换为 application handoff types。
- 当 validation rules 存在时，把 malformed identity/session metadata 映射为 protocol 或 application errors。

不得：

- 选择或实现认证方案。
- 拥有 long-lived session state。
- 未经 protocol change spec 和 ADR 修改 envelope shape。
- 执行 module-level permissions。

### Application Dispatch And Session Validation

Owner：

```text
runtime/internal/app
```

职责：

- 接收 decoded route requests。
- 当 session boundary 实现后调用 session validation。
- 把 normalized request identity context 传给 command 和 query handlers。
- 保持 route dispatch 与 authentication provider details 分离。

不得：

- 存储 player account lifecycle data。
- import WebSocket libraries。
- import generated Protobuf packages。
- 隐藏 module business rules。

### Player Module

未来 owner：

```text
modules/player
runtime/internal/modules/player
```

职责：

- 拥有 player identity 和 player account lifecycle。
- 实现前先定义公开 account/session commands、queries、events、permissions 和 errors。
- 当 migrations 被确认后拥有 persistent player account state。
- 在其他模块需要时发布 account lifecycle events。

不得：

- 拥有 inventory state。
- 拥有 WebSocket connection mechanics。
- 隐藏 transport 或 protocol behavior。

### Inventory Module

Owner：

```text
modules/inventory
runtime/internal/modules/inventory
```

职责：

- 拥有以 `player_id` 为 key 的 inventory records 和 item quantities。
- 执行 inventory permissions 和 invariants。
- 在 player module 存在之前，把 `player_id` 视为 external reference。

不得：

- 创建、认证、关联、禁用、删除或迁移 player accounts。
- 拥有 session state 或 token validation。
- 在未来 change 未明确允许前，直接依赖 player module。

## 5. 当前 Envelope 位置

现有 protocol envelope 已经包含：

```text
Session.connection_id
Session.session_id
Session.player_id
Session.connection_epoch
```

本标准不改变这些字段。

当前含义：

- `connection_id`：可用时的 transport-local metadata。
- `session_id`：保留的 authenticated logical session identifier。
- `player_id`：可用时保留的 authenticated player identity。
- `connection_epoch`：保留的 reconnect/lifecycle version。

在 session validation 存在之前，这些字段只是 envelope metadata。Domain logic 不得把它们当作 identity proof。

## 6. 请求处理模型

未来目标请求流：

```text
websocket frame
-> transport connection metadata
-> protobuf envelope
-> preliminary route request
-> session validation
-> request identity context
-> application dispatch
-> domain command or query handler
```

当前 runtime 状态：

- WebSocket transport 和 Protobuf request dispatch 已经存在。
- Inventory handlers 可以接收 `app.Session` 和 `app.Target` metadata。
- 已经存在 application-owned session validation hook，并提供 metadata-only default path。
- Inventory permission policies 已经能接收 application-owned request identity handoff context，但当前 bootstrap policy 仍然是显式 static。
- 还没有真实 authentication、token parsing、credential lookup、session persistence 或 player account lookup。
- Metadata-only request identity 不是 authenticated proof，不得满足 identity-aware privileged permission policy。

## 7. 权限边界

Permissions 是 module-owned business checks。

Authentication 和 session validation 回答：

```text
Who is the actor?
```

Module permission policies 回答：

```text
May this actor perform this operation on this module target?
```

Inventory 必须在未来 work item 中摆脱 bootstrap static allow behavior，但这个未来步骤仍然不得在未经单独确认时选择 authentication scheme。

对于 inventory：

- `GrantItem` 需要 `inventory_grant_item`。
- `GetInventory` 需要 `inventory_read`。
- Inventory permission context 携带 requested actor text、目标 `player_id` 和 `RequestIdentity`。
- Static bootstrap permission policy 可以显式允许当前 local proof-slice behavior，但它不是 production authorization。
- Identity-aware permission policies 必须把 `metadata_only` identity 当作 unauthenticated。
- 未来 player-bound read policy 可以允许玩家只读取自己的 inventory。
- Service 或 admin actors 需要明确 actor/permission modeling 后，privileged grants 才能成为生产行为。

## 8. 延后决策

本标准有意不决定：

- Token format。
- JWT、OAuth、OIDC、password、device、guest、social 或 external-provider authentication。
- Credential storage。
- Session persistence store。
- Session expiration 和 refresh model。
- Reconnect replay behavior。
- Player account database schema。
- Player account migration files。
- Protobuf envelope changes。
- WebSocket handshake authentication contract。
- Presence、parties、rooms、matches 或 broadcast group behavior。

每一项都需要单独 change spec；涉及架构时还需要 ADR。

## 9. Application Handoff

当前 Go runtime 在以下位置定义 application-owned `RequestIdentity` handoff：

```text
runtime/internal/app
```

规则：

- `RouteRequest.Identity` 是 command 和 query handlers 面向 application 的 identity context。
- `ApplicationResult.Identity` 会保留 request identity context，供下游 protocol mapping 和未来 auditing 使用。
- `MetadataOnlyIdentityFromSession` 可以把现有 session metadata normalize 到 request identity context。
- `MetadataOnlySessionValidator` 会保留当前 metadata-only behavior，但不会认证 clients。
- `SessionValidatingDispatcher` 是 application-owned hook，它在 protocol decoding 之后、module handlers 收到 request 之前运行。
- Metadata-only identity 不是 authenticated identity。
- 在未来 session validator 真正完成验证前，`PlayerIDValidated` 和 `SessionValidated` 保持 false。
- `SessionValidationResult` 是 handoff vocabulary 和 hook output。它不是 authentication implementation、token contract 或 session store。
- 未来 real session validation 应在这个 hook 后面实现，并在验证成功时用 validated identity 替换 metadata-only identity。
- Inventory 只通过自己的 permission handoff context 使用 `RequestIdentity`。Inventory 不得验证 token、查询 player account、拥有 session，或 import player module。

## 10. 已完成的边界序列

M-003 boundary sequence 已完成：

1. Player identity 和 session responsibilities 已在本标准和 `ADR-0021` 中分离。
2. `modules/player/module.yaml` 和 module agent guides 声明 player identity 与 account lifecycle ownership，但不包含 runtime code、migrations、credential storage 或 authentication providers。
3. Application-owned request identity 和 session validation handoff types 已存在于 `runtime/internal/app`。
4. `MetadataOnlySessionValidator` 和 `SessionValidatingDispatcher` 提供未来 validation hook，同时保持 metadata-only behavior。
5. Inventory permission policies 通过 `PermissionContext` 接收 `RequestIdentity`，并且 `MetadataOnlyDenyPermissionPolicy` 防止 metadata-only identity 满足 privileged grants。
6. `runtime.identity_boundary` repository checks 保护该边界，避免常见 transport、domain、generated Protobuf、auth dependency、player runtime、未经 ratify 的 player Protobuf 和 player persistence regressions。

下一步 implementation 必须等待 maintainer 确认 major direction。

可能的下一阶段 milestone 包括：

- 确认 player account 和 session public contracts。
- 确认 authentication、token 和 session validation design。
- 在 production authentication 前继续扩展 game-domain breadth，例如 item catalog、currency、rewards、quests 或 match sessions。
- 在扩展 runtime features 前改进 generators 和 contract tooling。

Agents 必须在任何选择具体 authentication mechanism、token format、credential store、player account schema、session persistence model、Protobuf envelope change 或 WebSocket handshake contract 的步骤前停止。

## 11. 验证

当前 repository checks：

```bash
node tools/vibit check architecture
node tools/vibit check protocol
node tools/vibit check runtime
node tools/vibit check work
node tools/vibit check all
```

`node tools/vibit check runtime` 包含 `runtime.identity_boundary` repository check。

当前该 check 会验证：

- WebSocket transport 不 import runtime domain modules、player runtime packages、inventory runtime packages、generated Protobuf packages 或 Protobuf runtime dependencies。
- Domain modules 不 import WebSocket transport、generated Protobuf packages、Protobuf runtime dependencies，或已知 authentication、token、OAuth、OIDC、credential、password-hashing dependencies。
- 在 public player contracts 被确认前，runtime player module code 仍然不存在。
- Player Protobuf source roots 只允许用于已 ratify 的 player wire contracts，并且不得暗示 runtime handlers、authentication 或 persistence 已存在。
- 在 schema ratification 前，PostgreSQL migrations 不引入 player identity 或 player account persistence。
- `modules/player/module.yaml` 在 semantic 和 wire contracts 被 ratify 后仍保持 identity/session boundary markers。

有意延后的 checks：

- 证明 semantic authentication correctness。当前 repository 还没有 authentication implementation 可验证。
- 证明所有可能的 third-party authentication libraries。当前 check 只阻断已知 dependency families；当具体 authentication decision 被确认后必须扩展。
- 检查第一版已 ratify player account Protobuf package 之外的 generated player/session output。Runtime session generated output 目前还不存在。
- 证明 production authorization 的 permission policy completeness。当前 policies 只定义 handoff 和 metadata-only guard boundary。

## 12. Agent Rules

Agents 必须：

- 在添加 player、account、authentication、session、permission 或 request identity 行为前阅读本标准。
- 保持 player identity、account lifecycle、authentication、runtime session、transport connection 和 request identity context 分离。
- 在实现 public behavior 前新增或更新 manifests 和 contracts。
- 在相关 change spec 或 ADR 中记录 adopted、adapted、deferred 或 rejected 的 Nakama/Pitaya pattern。
- 在跨越 deferred decision boundary 前询问。

Agents 不得：

- 在 session validation 存在前，把 client envelope 提供的 `player_id` 当作 authenticated proof。
- 在 boundary-only work 中添加 player account migrations。
- 未经确认添加 JWT、OAuth、OIDC、password hashing、password storage 或 external identity provider dependencies。
- 把 inventory ownership 移入 player module。
- 把 authentication 或 permission checks 隐藏在 WebSocket transport 或 Protobuf bridge code 中。
