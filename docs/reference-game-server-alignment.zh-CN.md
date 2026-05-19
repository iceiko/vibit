# Reference Game Server Alignment 中文版

状态：Draft v0.1
最后更新：2026-05-13
范围：vibit game server capability planning 的主动参考基线
说明：本文件是 `docs/reference-game-server-alignment.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文档记录 vibit 应如何使用成熟 game server projects 作为参考。

`ADR-0078` 和 `docs/nakama-pitaya-product-parity-roadmap.md` 已把这个 reference baseline 细化为明确产品目标：vibit 应通过覆盖 common capability families，成为 Nakama/Pitaya-class game backend framework，同时保留 vibit 的 agent-native constraints。本文档仍定义 reference roles；product parity roadmap 定义分阶段执行。

## 1. 目的

vibit 不只是一个 inventory proof slice。

vibit 最终应覆盖与 Nakama、Pitaya 等成熟 game backend 和 game server frameworks 相同的大问题域：

- Player identity 和 sessions。
- Social 和 realtime communication。
- Storage 和 durable game state。
- Matchmaking 和 match/session lifecycle。
- Realtime multiplayer 和 authoritative server behavior。
- Leaderboards、rewards、currencies 等 metagame systems。
- Operational visibility 和 production maintainability。
- Scalable routing，以及后续 distributed server topology。

区别不是这些 game server capabilities 是否存在。区别是 vibit 必须把这些能力做成 Agent-Native：

- 显式 module ownership。
- Contract-first public behavior。
- Generated repeatable structure。
- Machine-readable manifests。
- 窄 runtime boundaries。
- Durable decision records。
- Agents 在变更前后都能运行的 verification commands。

Nakama 记录了一套 broad game server product surface，包括 user accounts、authentication、storage、friends、groups、chat、leaderboards、tournaments、matchmaking、realtime multiplayer 和 authoritative match runtime concepts。Pitaya 记录了一套围绕 client acceptors、sessions、route handlers、remote calls、groups、serializers 和 cluster vocabulary 的 Go game server framework shape。vibit 应同时向这两类 surface 学习，但每个被采纳的 capability 都必须显式、必要时可生成，并且机器可检查。

当前产品目标标记：

```text
nakama_pitaya_product_parity_roadmap: ratified
decision: ADR-0078
check_rule: runtime.reference_product_parity_roadmap
parity_goal: nakama_pitaya_same_class_common_capability_coverage
api_compatibility_goal: false
recommended_next_direction: define_protocol_logout_route_gate
```

## 2. Reference Roles

### Nakama

把 Nakama 作为 broad game backend product capability surface 的主要参考。

参考领域：

- Authentication、users、accounts 和 sessions。
- Friends、groups、chat、parties、leaderboards 和 tournaments 等 social systems。
- Storage objects 和 server-side runtime customization。
- Match listing、matchmaker、realtime multiplayer 和 authoritative match logic。
- Dashboard、metrics 和 operational visibility。
- 在观察更广义 Heroic Labs 生态时，参考 economy、rewards、currencies 和 LiveOps-style capability families。

Nakama 应指导一个有用的 general game backend 最终应该支持什么。

除非未来 ADR 明确采纳 compatibility surface，否则 Nakama 不应成为 vibit 的治理性 API shape。

### Pitaya

把 Pitaya 作为 Go game server framework architecture vocabulary 的主要参考。

参考领域：

- WebSocket 和 TCP 等 client connection acceptors。
- User sessions 和 session binding。
- 面向 client messages 的 handler routing。
- 面向 server-to-server communication 的 remote calls。
- Cluster mode 下的 frontend 和 backend server roles。
- 用于 rooms 等 broadcast/multicast use cases 的 groups。
- Message forwarding、serializers、RPC services 和 cluster service discovery。

Pitaya 应指导 Go game server frameworks 如何分离 transport、session、route、server role、RPC 和 group concerns。

Pitaya 不应迫使 vibit 在 modular monolith proof slice 健康之前进入 distributed clustering。

## 3. Capability Matrix

以下矩阵是 planning tool，不承诺所有能力必须立即实现。

| Capability | Reference | Vibit Direction |
| --- | --- | --- |
| Accounts and authentication | Nakama | `player` 或 `identity` module，带显式 auth/session contracts。 |
| Sessions and connection identity | Nakama, Pitaya | Platform session adapter 加 app-owned session context；不得使用 transport shortcut identity。 |
| Storage objects | Nakama | PostgreSQL-backed module state first；object storage 只用于 large artifacts。 |
| Inventory | Common game backend need | 第一 proof module；必须证明 contract -> generated shape -> handler -> tests。 |
| Currency and wallets | Nakama/Hiro capability family | 未来 `currency` module，带 transactional invariants。 |
| Rewards and claims | Nakama/Hiro capability family | 未来 `reward` module，带 eligibility、idempotency 和 event tests。 |
| Friends, groups, parties | Nakama | 未来 social modules，显式 ownership membership 和 events。 |
| Chat and realtime messaging | Nakama | 未来 realtime module；不得隐藏在 WebSocket transport 中。 |
| Presence and status | Nakama, Pitaya | 未来 platform/application capability，带显式 lifecycle semantics。 |
| Matchmaking | Nakama | 未来 `matchmaking` module；实现前先有 query 和 criteria contracts。 |
| Match/session lifecycle | Nakama, Pitaya | 未来 `match` module；authoritative match behavior 与 transport 分离。 |
| Rooms and broadcast groups | Pitaya | 未来 group/room abstraction；target scopes 已保留 `room` 和 `match`。 |
| Leaderboards and tournaments | Nakama | 未来 competitive modules；ranking 和 reset rules 必须 contract-first。 |
| Authoritative realtime simulation | Nakama | 未来 match runtime；server remains authoritative。 |
| Cluster frontend/backend split | Pitaya | 延后，直到 single-process boundaries 被证明。 |
| Server-to-server RPC | Pitaya | 延后；引入时不得绕过 module contracts。 |
| Dashboard and operations | Nakama | 未来 admin/inspection surface；不得混入 gameplay protocol。 |
| Metrics and observability | Nakama, Pitaya | 延后的 dependency decision；必须由 platform 拥有。 |

## 4. Phased Roadmap

### Phase 0: Agent-Native Foundation

当前阶段。

目标：

- Constitution、AGENTS guides、change specs、conversation logs、ADRs。
- Architecture manifests 和 checks。
- Go/WebSocket/Protobuf/PostgreSQL direction。
- Runtime handoff、protocol adapter 和 application dispatch skeleton。
- 第一版 inventory semantic 和 wire contracts。

退出标准：

- 一个小 backend change 能走完 spec、contract、generated shape、handwritten logic、tests、verification 和 docs。

### Phase 1: First Vertical Game Backend Slice

目标：

- Inventory command/query/event handler boundary。
- Repository 和 policy interfaces。
- 第一版 durable module state 的 PostgreSQL persistence。
- 针对 invariants、errors、events 和 repository behavior 的 focused tests。
- 通过现有 dispatcher 映射 protocol response。

参考：

- Nakama storage 和 custom server logic concepts。
- Pitaya handler routing separation。

### Phase 2: Player, Session, And Transport

目标：

- Player/account/session module boundaries。
- WebSocket transport adapter。
- Session validation 和 connection lifecycle。
- Protocol error response mapping。
- Minimal playable client/server request flow。

参考：

- Nakama authentication 和 session model。
- Pitaya session binding 和 WebSocket acceptor separation。

### Phase 3: Core Game Backend Modules

目标：

- Currency/wallet。
- Rewards/claims。
- Friends/groups/party。
- Presence/status。
- Chat 或 realtime messaging。
- Leaderboards。

参考：

- Nakama social、competitive 和 metagame capabilities。

### Phase 4: Multiplayer And Match Runtime

目标：

- Matchmaking criteria 和 ticket lifecycle。
- Match/session lifecycle。
- Room/group broadcast semantics。
- Authoritative match loop contracts。
- Reconnect/replay decisions。

参考：

- Nakama matchmaker 和 authoritative/relayed multiplayer distinction。
- Pitaya groups 和 route/handler model。

### Phase 5: Distributed Runtime

目标：

- Frontend/backend server role split。
- Server-to-server RPC。
- Service discovery。
- Distributed groups/rooms。
- Cluster-safe session 和 routing semantics。

参考：

- Pitaya cluster architecture。

该阶段不得在 single-process module、transaction、protocol 和 verification boundaries 稳定前开始。

## 5. Agent Rules

Agents 必须：

- 在提出新的 game server modules 或 runtime subsystems 前阅读本文档。
- 检查拟议 capability 是否映射到已知 Nakama/Pitaya capability family。
- 即使匹配参考功能，也保留 vibit 的 Agent-Native constraints。
- 记录某个 reference pattern 是被采纳、改造、推迟还是拒绝，以及原因。
- 优先增加小而可执行的 manifest/check，而不是宽泛愿景文字。

Agents 不得：

- 在没有明确 compatibility ADR 的情况下复制外部 APIs。
- 把 Nakama-like 或 Pitaya-like feature 直接加入 transport handlers。
- 在 modular monolith proof slice 稳定前加入 cluster/RPC/service-discovery work。
- 让 feature parity 优先于 explicit contracts、module ownership 或 verification。
- 用 references 作为绕过 vibit generated-file 和 boundary rules 的理由。

## 6. Next Planning Implications

如果 handler、repository 和 policy boundaries 还不清楚，下一步 runtime work 仍不应急着写 WebSocket server。

推荐近期顺序：

1. 定义 inventory runtime repository 和 policy interfaces。
2. 通过 application dispatcher 实现第一版 inventory command/query handler。
3. 在 interfaces 稳定后添加 PostgreSQL migration 和 repository behavior。
4. 在 request -> dispatch -> handler -> result path 已被测试后添加 WebSocket transport。
5. 在把 WebSocket connection identity 当作 durable player identity 之前，先添加 player/session/auth boundaries。

这个顺序让项目与 Nakama/Pitaya 功能面保持一致，同时保留 vibit 的主要差异。

## 7. References

- Nakama documentation home: `https://heroiclabs.com/docs/`
- Nakama getting started: `https://heroiclabs.com/docs/nakama/getting-started/`
- Nakama multiplayer concepts: `https://heroiclabs.com/docs/nakama/concepts/multiplayer/`
- Nakama GitHub repository: `https://github.com/heroiclabs/nakama`
- Pitaya overview: `https://pitaya.readthedocs.io/en/latest/overview.html`
- Pitaya features: `https://pitaya.readthedocs.io/en/stable/features.html`
- Pitaya GitHub repository: `https://github.com/topfreegames/pitaya`
