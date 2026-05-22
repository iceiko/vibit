# Nakama 与 Pitaya 产品同级路线图

状态：Draft v0.1
最后更新：2026-05-18
范围：vibit 面向 Nakama/Pitaya 同级游戏后端目标的产品路线图标准

英文文件 `docs/nakama-pitaya-product-parity-roadmap.md` 是权威版本。本文是简体中文译本。

## 1. 目的

本标准把“参考 Nakama 和 Pitaya”升级为明确的产品目标：

```text
nakama_pitaya_product_parity_roadmap: ratified
completed_work_item: W-0168
decision: ADR-0078
check_rule: runtime.reference_product_parity_roadmap
parity_goal: nakama_pitaya_same_class_common_capability_coverage
api_compatibility_goal: false
direct_nakama_pitaya_api_compatibility_added: false
implementation_authorized_by_this_standard: roadmap_only
```

vibit 应成为 Nakama/Pitaya 同级别的开源游戏后端框架，也就是 Nakama/Pitaya-class game backend product。这意味着 vibit 必须覆盖游戏团队通常期待这些系统提供的常用能力族，同时保留 vibit 的核心差异化：通过 explicit contracts、manifests、generation、tests、ADRs 和 repository checks 实现 agent-native maintainability。

产品同级意味着能力覆盖和运维实用性接近。它不意味着直接 API 兼容，不意味着复制公开 routes、data models、clustering internals，也不意味着照搬任一项目的实现细节。

## 2. 产品目标

Nakama 仍是 broad game backend product coverage 的主要参考：

- Accounts、authentication、users、sessions。
- Storage objects 和 durable game state。
- Friends、groups、parties、chat、status、presence、notifications。
- Leaderboards、tournaments、rewards、currencies 以及其他 metagame systems。
- Match listing、matchmaking、realtime multiplayer、authoritative match runtime。
- Server runtime customization、hooks、RPC-like extension points、streams、console、metrics、operations。
- Client SDK 和 sample application ergonomics。

Pitaya 仍是 Go game server architecture vocabulary 的主要参考：

- Acceptors 和 connection lifecycle。
- Sessions、binding、kick/disconnect、session data。
- Handler routing、pipelines、serializers、message forwarding。
- Groups、broadcast、multicast、push。
- Frontend/backend server roles。
- Server-to-server RPC、service discovery、cluster mode、monitoring、tracing。

vibit 必须把这些能力适配进自己的模型：

- Contract-first public behavior。
- Module-owned invariants。
- Generated repeatable structure。
- Application-owned lifecycle policy。
- Transport/protocol/domain separation。
- Repository and persistence boundaries。
- Agent-readable guides 和 checkable architecture。

## 3. 同级能力族

以下能力族属于一等路线图范围：

```text
parity_capability_families:
  - identity_authentication_sessions
  - connection_lifecycle_reconnect_logout
  - storage_objects_and_durable_game_state
  - presence_status_and_notifications
  - chat_streams_and_realtime_messaging
  - friends_groups_and_parties
  - leaderboards_tournaments_and_competitive_systems
  - economy_inventory_rewards_currencies_and_progression
  - matchmaking_match_listing_and_room_lifecycle
  - realtime_multiplayer_and_authoritative_match_runtime
  - server_runtime_hooks_rpc_and_custom_logic
  - admin_console_metrics_observability_and_operations
  - client_sdks_examples_and_developer_experience
  - distributed_runtime_frontend_backend_rpc_and_service_discovery
```

每个能力族最终都应具备：

- Module 或 runtime subsystem owner。
- 公开行为存在时要有 semantic contract surface。
- 存在 client/server messages 时要有 protocol surface。
- 存在 durable state 时要有 storage boundary。
- 针对 invariants 和 error behavior 的 tests。
- 针对关键边界的 repository checks 或 architecture checks。
- 面向公开读者时同时维护英文与简体中文文档。

## 4. 当前状态

已完成的基础：

- Agent-native 项目治理、change specs、ADRs、conversation memory 和 checks。
- Go runtime layout。
- WebSocket transport 和 Protobuf envelope。
- PostgreSQL persistence 和 migration tooling。
- Inventory proof slice。
- Player account repository 和 PostgreSQL adapter。
- Device credential login、opaque access-token validation、logout、runtime session persistence、session validation、route protection。
- First-message connection binding。
- Single-process active connection registry。
- Single-process WebSocket close policy，可以 invalidates registry records，但不执行 concrete socket close handoff。

当前缺口：

- Source-first alpha 已经可见，但下一产品阶段仍需要明确的 prototype-ready foundation execution plan。
- Prototype-ready foundation execution plan 已记录，下一项 work 必须先定义 local development path gate，再进入扩大 shared online services 的 implementation slices。

## 5. 阶段计划

### Phase 2R：Runtime Lifecycle Closure

状态：

```text
phase_2r_runtime_lifecycle_closure: active
current_near_term_priority: protocol_logout_and_connection_lifecycle
```

目标：在扩展产品模块前，完成 login、route protection、runtime session、connection binding、logout、close intent、concrete close、reconnect、session-carrier 的闭环。

近期必须推进的 gates：

1. `define_protocol_logout_route_gate`
2. `define_transport_close_handoff_gate`
3. `define_reconnect_connection_epoch_gate`
4. `define_protocol_session_carrier_gate`
5. `define_presence_lifecycle_gate`
6. `strengthen_operations_observability_and_admin_tooling`

这一阶段吸收 Nakama 对 account/session/socket lifecycle 的显式语义压力，也吸收 Pitaya 对 acceptor/session/handler 分层的经验，然后再让更高层的 social 或 multiplayer modules 依赖这些生命周期行为。

### Phase 3：Shared Online Services

目标：补齐许多上层功能依赖的常驻在线服务。

候选工作：

- 超出 module-local inventory state 的 storage object contracts 和 permissions。
- Presence 和 status lifecycle。
- Notifications。
- Chat、streams、realtime messaging。
- Server push 和 broadcast vocabulary。
- 面向 players、sessions、tokens、active connections 的 admin inspection。

### Phase 2P：Prototype-Ready Foundation

状态：

```text
phase_2p_prototype_ready_foundation: next_product_stage
standard: docs/product-maturity-milestones.md
execution_plan: docs/prototype-ready-foundation-execution-plan.md
local_development_path_gate: docs/prototype-ready-local-development-path-gate.md
local_development_path_package: docs/prototype-ready-local-development-path-package.md
storage_objects_behavior_gate: docs/storage-objects-behavior-gate.md
storage_objects_persistence_schema_gate: docs/storage-objects-persistence-schema-gate.md
storage_objects_repository_boundary: docs/storage-objects-repository-boundary.md
next_work_item: W-0205 Implement storage-neutral storage objects repository interface
```

目标：从开发者可以检查的 source-first alpha，推进到可以用于严肃小型 prototype 的 foundation。

候选工作：

- 降低 local setup、migration 和 configuration friction。已由 `W-0200` 完成。
- 添加更清晰的 example client 或 example app path。
- 定义超出 inventory proof slice 的 first general storage-object behavior。已由 `W-0201` 完成。
- 定义第一版 storage objects persistence schema posture。已由 `W-0202` 完成。
- 添加第一版 storage objects migration source。已由 `W-0203` 完成。
- 定义 storage objects repository boundary。已由 `W-0204` 完成。
- 实现 storage-neutral storage objects repository interface。Next。
- 定义第一版 realtime messaging、stream、broadcast 或 server-push behavior。
- 加强 authenticated gameplay loop 周围的 concurrency 和 failure-path verification。
- 定义 serious prototype 使用前需要的最小 operations inspection surface。

这一阶段不声明 production readiness。它选择下一批最小 product-useful slices，同时保留 source-alpha honesty 和现有 ask-first boundaries。

### Phase 4：Social And Competitive Modules

目标：覆盖 Nakama 风格的 metagame 和 social surfaces。

候选模块：

- Friends。
- Groups。
- Parties。
- Leaderboards。
- Tournaments。
- Wallet/currency。
- Rewards/claims。
- Quests/progression。

### Phase 5：Matchmaking And Match Runtime

目标：覆盖 Nakama 期望的 multiplayer 能力族，并吸收 Pitaya 的 routing/group semantics。

候选工作：

- Match listing。
- Matchmaking tickets 和 criteria。
- Room lifecycle。
- Broadcast groups 和 target scopes。
- Authoritative match runtime contracts。
- Relayed realtime multiplayer contracts。
- Reconnect/replay decisions。

### Phase 6：Runtime Extensibility And Developer Experience

目标：让 vibit 成为真正可用的 framework，而不只是一个 server codebase。

候选工作：

- Server runtime hooks。
- Server-side custom logic/RPC surface。
- Module scaffolding 和 generator hardening。
- Client SDKs 和 examples。
- Local development runbook。
- Admin console 和 CLI workflows。

### Phase 7：Distributed Runtime

目标：只在 single-process semantics 稳定后，再引入 Pitaya 风格的 distributed topology。

候选工作：

- Frontend/backend server role split。
- Server-to-server RPC。
- Service discovery。
- Distributed groups 和 push。
- Cluster-safe session routing。
- Multi-node presence 和 matchmaking behavior。

这一阶段不能通过削弱 single-process contracts 来开始。它应该把已经验证的单进程语义提升到 distributed adapters。

## 6. 开发方式

每个能力族应按以下顺序推进：

1. Reference review：记录正在覆盖的 Nakama/Pitaya 能力。
2. Vibit ownership：决定 module、runtime、platform、generated、operations owners。
3. Semantic contract：定义 commands、queries、events、errors、permissions、invariants。
4. Protocol contract：semantic behavior 稳定后再定义 wire messages。
5. Persistence boundary：在 adapters 前定义 tables、repositories、indexes、redaction。
6. Application behavior：实现最小 vertical slice。
7. Operations surface：当行为需要 operator 时，补 inspection、metrics、admin actions 或 runbook。
8. Verification：补 focused Go tests 和 repository checks。
9. Memory：更新 change specs、ADRs、manifests、conversation logs。

首选实现形态仍然是：

```text
requirement -> reference mapping -> spec -> contract -> generated shape -> logic -> tests -> checks -> docs
```

## 7. 近期建议

本路线图 ratify 后，下一个具体工作不应直接跳到 chat、groups、matchmaking 或 match runtime。下一个具体工作应先完成 runtime lifecycle foundation：

```text
recommended_next_direction: define_protocol_logout_route_gate
second_direction: define_transport_close_handoff_gate
third_direction: define_reconnect_connection_epoch_gate
fourth_direction: define_protocol_session_carrier_gate
first_module_expansion_after_lifecycle: define_presence_lifecycle_gate
```

理由：

- Logout 已存在 service behavior，但尚未暴露为 client protocol route。
- Close policy 可以 invalidates registry records，但还不能关闭 concrete WebSocket sockets。
- Reconnect 和 duplicate connection behavior 不应等 presence 或 match runtime 依赖它之后才设计。
- Protocol session carriers 必须先明确，客户端才能安全理解 runtime sessions。
- Presence 和 chat 依赖 connection/session lifecycle semantics。

## 8. 非目标

本路线图不做以下事情：

- 自身实现任何新的 runtime behavior。
- 添加 protocol logout routes。
- 添加 concrete WebSocket close handoff。
- 添加 reconnect、resume、duplicate replacement 或 session carrier behavior。
- 添加 presence、chat、friends、groups、parties、leaderboards、tournaments、matchmaking 或 match runtime 代码。
- 添加 admin console、SDKs、cluster、RPC、service discovery 或 distributed runtime behavior。
- 添加 dependencies。
- 承诺 Nakama 或 Pitaya API compatibility。

## 9. Agent 规则

Agents 必须：

- 把 Nakama/Pitaya 同级常用能力覆盖当成产品要求，而不是背景阅读。
- 把新的 major work 映射到一个 roadmap family。
- 在 high-level social 或 multiplayer modules 之前，优先完成 lifecycle 和 shared services。
- 即使匹配 reference capabilities，也要保留 vibit 的 agent-native constraints。
- 在 ADRs 或 change specs 中记录 adopted、adapted、deferred、rejected reference patterns。
- 把 direct API compatibility 放在未来显式 compatibility ADR 后面。

Agents 禁止：

- 把 product features 直接加进 WebSocket transport handlers。
- 意外添加 direct external API compatibility。
- 因为 reference system 已有类似功能就跳过 semantic contracts。
- 在 single-process semantics 测试稳定前启动 distributed runtime work。
- 把 product parity 当成削弱 generated-file、redaction、permission、session 或 module-boundary rules 的许可。
