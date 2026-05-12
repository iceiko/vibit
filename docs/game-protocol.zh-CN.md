# Game Protocol Framework Standard 中文版

状态：Draft v0.1  
最后更新：2026-05-13  
范围：vibit 第一版 gameplay/client protocol framework  
权威决策：`ADR-0015`  
说明：本文件是 `docs/game-protocol.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本标准定义 vibit 第一版面向游戏客户端的协议，在创建 `.proto` files、runtime transport handlers 或 generated clients 之前应具备什么形态。

目标不是克隆某个现有游戏后端。目标是把成熟 game-server patterns 吸收到一个 agent-native 的 protocol surface 中，让 agents 能检查、扩展和验证。

## Problem

游戏后端协议不只是 RPC layer。

真实游戏需要 long-lived connections、authenticated sessions、player-scoped requests、server push、error semantics、reconnect behavior、room 或 match targeting、authoritative input handling，以及未来的 allocation 或 matchmaking flows。

如果 vibit 从一个通用 request/response protocol 开始，后续 agents 可能会把游戏行为藏在 route strings、transport handlers 或 module internals 里。这会削弱项目核心承诺：agents 应该从显式 contracts 和 manifests 理解并修改服务器。

## External References

成熟系统提供了有价值的模式：

- Nakama 区分 ordinary APIs 和 realtime socket 使用，并建模 presence、streams、parties、matchmaker flows 和 server-authoritative multiplayer。
- Colyseus 以 rooms、server-owned state、同步 state changes，以及 joined room 内的 client messages 为 multiplayer behavior 中心。
- Pitaya 和 Pomelo 体现了 connector 或 frontend network ownership、route-based dispatch、connection sessions、server push、handlers 和 backend remotes 的价值。
- Skynet 的 service/message-dispatch 风格对 server architecture 和 concurrency boundaries 有参考价值，但不应定义 vibit 第一版 client wire protocol。
- Agones 和 Open Match 表明 game server allocation 与 match construction 是独立的 lifecycle concerns。它们之后很重要，但不应被塞进第一版 in-game WebSocket envelope。

vibit 应借鉴词汇和经过验证的边界，而不是整体复制某个框架。

参考阅读：

- Nakama Multiplayer Engine：`https://heroiclabs.com/docs/nakama/concepts/multiplayer/`
- Colyseus State Synchronization：`https://docs.colyseus.io/state`
- Pitaya Documentation：`https://pitaya.readthedocs.io/`
- Pomelo repository：`https://github.com/NetEase/pomelo`
- Skynet repository：`https://github.com/cloudwu/skynet`
- Agones GameServerAllocation：`https://agones.dev/site/docs/reference/gameserverallocation/`
- Open Match Match Function：`https://openmatch.dev/site/docs/guides/matchmaker/matchfunction/`

## Protocol Model

vibit 第一版游戏协议是 WebSocket-framed Protobuf envelope。

Envelope 由 platform protocol code 拥有，不由 domain modules 拥有。Domain modules 通过 vibit manifests 和 contract sources 拥有 semantic commands、queries、events、invariants、permissions 和 errors。

第一版 protocol surface 应使用一个 gameplay WebSocket endpoint，初始 endpoint path 计划为：

```text
/v1/ws
```

Endpoint path 可在第一版 transport implementation 开始时最终确定。客户端存在后再修改它，属于 compatibility-sensitive protocol change。

## Envelope Responsibilities

Envelope 必须显式表达这些概念：

- Protocol version。
- Message kind。
- Request correlation。
- Module identity。
- Operation name。
- Game target scope。
- Session metadata。
- Payload encoding identity。
- Error mapping。

Envelope 不得包含 domain business logic。

初始 message kinds 是：

```text
command
query
event
error
ack
heartbeat
system
input
state
```

第一版 inventory slice 只需要使用 `command`、`query`、`event`、`error` 和 `system`。`ack`、`heartbeat`、`input` 和 `state` 先保留，因为成熟游戏协议需要它们，但应等待具体需求再实现。

## Routing

Routing 必须 agent-readable。

第一版 route identity 使用结构化字段：

```text
kind
module
name
```

Semantic route key 可以渲染为：

```text
<module>.<name>
```

示例：

```text
inventory.GrantItem
inventory.GetInventory
inventory.ItemGranted
```

不要把业务含义埋进 opaque path string。Transport handlers 解码 envelope，并把 route 交给 application dispatch。Domain modules 不解析 WebSocket frames 或 route strings。

## Session Model

Protocol model 区分 transport connection、logical session 和 player identity：

- `connection_id` 是 transport-local，reconnect 后可能变化。
- `session_id` 是存在 authentication 后的 authenticated logical session。
- `player_id` 是可用时的玩家 domain identity。
- `connection_epoch` 或等价 reconnect version 为未来 reconnect rules 保留。

在 player/auth module 存在前，inventory protocol work 可以把 authenticated identity 当作计划中的 context，而不要发明 authentication shortcut。

## Target Model

游戏消息的目标经常不只是单个 RPC handler。

Envelope 应保留 target metadata，并包含这些 scope values：

```text
global
player
party
room
match
stream
system
```

第一版 inventory slice 应对 player-owned inventory operations 使用 `player` scope。

Room、match、party 和 stream behavior 先保留。在相关 module、contracts 和 runtime lifecycle rules 存在前，不得用隐藏的 ad hoc fields 实现它们。

## Client-To-Server Messages

Client-to-server messages 可以表示：

- Commands：client 改变 server state 的意图。
- Queries：client 读取 server state 的请求。
- Inputs：realtime authoritative gameplay input，保留到 match 或 room module 存在后。
- System messages：protocol-level negotiation 或 lifecycle messages。

Client messages 是 requests，不是 facts。服务器决定 command、query 或 input 是否有效。

## Server-To-Client Messages

Server-to-client messages 可以表示：

- 与 `request_id` 关联的 command 或 query responses。
- Domain event pushes。
- Error envelopes。
- System/session messages。
- State snapshots 或 patches，保留到 room 或 match state-sync standard 存在后。
- Acknowledgements 和 heartbeats，保留到 transport lifecycle 需要 WebSocket ping/pong 之上的协议层能力时。

Domain events 是 server facts。Client code 不允许直接向 domain model 发布 facts。

## Error Model

当 protocol errors 指向 public module behavior 时，必须映射到已登记的 vibit error catalogs。

第一版 wire error shape 应包含：

- Stable error code。
- Human-readable message。
- 可用时的 related `request_id`。
- Retryability signal。
- Optional structured details。

Transport errors、malformed envelopes、permission errors、invariant failures 和 unknown routes 必须可区分。

## Compatibility Rules

第一版 Protobuf package version 仍为：

```text
vibit.<module>.v1
```

Envelope versioning 与 module message versioning 相关但分离：

- Envelope version 管 framing、routing、target、session 和 system semantics。
- Module package version 管 module command、query、event 和 payload wire shapes。

在创建第一批 `.proto` files 之前，protocol standard 应定义 reserved fields 和 extension discipline。一旦 `.proto` files 存在：

- 不要复用已删除的 Protobuf field numbers。
- 不要随意重命名 public field semantics。
- 优先 additive evolution。
- Breaking changes 需要 change spec 和 ADR。

## Agent Rules

在新增或修改 `.proto` files、runtime protocol handlers 或 generated protocol code 前，agents 必须阅读：

- `docs/game-protocol.md`
- `.arch/protocol.yaml`
- `ADR-0009`
- `ADR-0014`
- `ADR-0015`
- `.arch/contracts.yaml`
- 受影响的 module manifest 和 contract files

Agents 必须保持这些边界：

- WebSocket transport 属于 `runtime/internal/platform/transport/ws/`。
- Protobuf envelope encode/decode 属于 `runtime/internal/platform/protocol/protobuf/`。
- Application dispatch 属于 `runtime/internal/app/`。
- Domain behavior 属于 `runtime/internal/modules/<module>/`。
- Generated Protobuf output 属于 `runtime/internal/generated/proto/`。
- Protobuf source files 属于 `proto/vibit/<module>/v1/`。

## First Implementation Guidance

第一版 inventory runtime slice 应证明 protocol framework，而不是假装实现所有 multiplayer concerns。

推荐第一条 flow：

```text
WebSocket frame
-> Protobuf envelope
-> kind/module/name dispatch
-> generated command or query payload
-> application dispatcher
-> inventory handler
-> response envelope
```

除非新的 change spec 和 ADR 把 room state sync、matchmaking、allocation、reconnect replay、presence 或 streams 纳入目标，否则不要在 inventory slice 中实现它们。

## Verification Direction

当前 verification：

```bash
node tools/vibit check protocol
```

Protocol check 应验证：

- `.arch/protocol.yaml` 存在。
- Manifest 引用 `ADR-0015`。
- Manifest 记录 WebSocket、Protobuf、envelope routing、session identity、target scopes、message kinds 和 server authority。
- 一旦 Protobuf source files 存在，它们与已登记 contracts 对齐。

未来 checks 应验证 `.proto` envelope shape、generated output traceability、reserved field policy 和 runtime dispatch ownership。
