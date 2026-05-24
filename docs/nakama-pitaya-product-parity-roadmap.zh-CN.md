# Nakama 优先产品能力路线图

状态：Draft v0.2
最后更新：2026-05-24
范围：vibit 面向 Nakama 优先游戏后端目标的产品路线图标准

英文文件 `docs/nakama-pitaya-product-parity-roadmap.md` 是权威版本。本文是简体中文译本。

## 1. 目的

本标准最初把“参考 Nakama 和 Pitaya”升级为明确产品目标。`ADR-0127` 进一步细化该姿态：Nakama 现在是主要产品能力参考，Pitaya 暂缓为未来 distributed Go game server 架构参考。

```text
nakama_pitaya_product_parity_roadmap: ratified
completed_work_item: W-0168
decision: ADR-0078
check_rule: runtime.reference_product_parity_roadmap
reference_posture_update: ADR-0127
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
parity_goal: nakama_first_same_class_common_capability_coverage
api_compatibility_goal: false
direct_nakama_pitaya_api_compatibility_added: false
implementation_authorized_by_this_standard: roadmap_only
```

历史 `ADR-0078` 检查 marker 会保留用于追踪 lineage。它们作为当前 planning guidance 已被 `ADR-0127` superseded，但旧 repository checks 仍用这些 marker 确认原始 roadmap 来源：

```text
parity_goal: nakama_pitaya_same_class_common_capability_coverage
recommended_next_direction: define_protocol_logout_route_gate
first_module_expansion_after_lifecycle: define_presence_lifecycle_gate
Nakama/Pitaya-class game backend product
```

保留历史 `W-0220` workflow markers 供 repository checks 追踪。它们记录的是当时打开、现在已经完成的 workflow pilot 方向，不是当前 next-ready item：

```text
recommended_next_direction: pilot_nakama_aligned_feature_request_workflow
```

vibit 应成为 Nakama-class 的开源游戏后端框架，并把 AI-native development 与 AI-native testing 作为产品目的。这意味着 vibit 必须覆盖游戏团队通常期待 Nakama-style 系统提供的常用能力族，同时保留 vibit 的核心差异化：用户用普通产品语言描述 backend requirement，AI agents 把需求转成 bounded specs、acceptance criteria、tests、implementation、verification records、ADRs、manifests 和 repository checks。

产品同级意味着能力覆盖和运维实用性接近。它不意味着直接 API 兼容，不意味着复制公开 routes、data models、clustering internals，也不意味着照搬 Nakama 的实现细节。Pitaya-style frontend/backend、RPC、service discovery、cluster groups 和 distributed topology 仍是未来事项，直到后续 ADR 重新把 Pitaya 激活为 active architecture reference。

## 2. 产品目标

Nakama 是 broad game backend product coverage 的主要参考：

- Accounts、authentication、users、sessions。
- Storage objects 和 durable game state。
- Friends、groups、parties、chat、status、presence、notifications。
- Leaderboards、tournaments、rewards、currencies 以及其他 metagame systems。
- Match listing、matchmaking、realtime multiplayer、authoritative match runtime。
- Server runtime customization、hooks、RPC-like extension points、streams、console、metrics、operations。
- Client SDK 和 sample application ergonomics。

Pitaya 不再是当前 product planning driver。它暂缓为未来 distributed Go game server topology 的架构参考：

- Acceptors 和 connection lifecycle。
- Sessions、binding、kick/disconnect、session data。
- Handler routing、pipelines、serializers、message forwarding。
- Groups、broadcast、multicast、push。
- Frontend/backend server roles。
- Server-to-server RPC、service discovery、cluster mode、monitoring、tracing。

Agents 不应使用 Pitaya 把 cluster/RPC/frontend-backend 工作提前拉入当前 prototype-ready foundation。Pitaya vocabulary 在解释 transport、session、protocol、application、backend service concerns 为何分离时仍可作为辅助，但当前 roadmap 的产品能力优先级由 Nakama 决定。

vibit 必须把这些能力适配进自己的模型：

- Contract-first public behavior。
- Module-owned invariants。
- Generated repeatable structure。
- Application-owned lifecycle policy。
- Transport/protocol/domain separation。
- Repository and persistence boundaries。
- Agent-readable guides 和 checkable architecture。
- AI-native requirement intake、test planning、implementation 和 verification。

## 2.1 AI-Native 产品目的

vibit 的产品目的不只是运行一个 game backend。它的产品目的还包括让后端研发本身变成 AI-native。

当用户描述一个新的 backend requirement，预期产品工作流是：

```text
user requirement
-> AI-written bounded requirement spec
-> AI-written acceptance criteria
-> AI-written test plan
-> AI-written or updated tests
-> AI implementation inside declared boundaries
-> AI-run verification
-> AI-updated docs, manifests, ADRs, and change records
```

架构存在的原因就是让这个工作流可靠。Contracts、module ownership、generated files、repository checks、redaction rules 和 verification commands 不是内部官僚流程；它们是让 AI agent 安全地把用户需求转成已测试 backend behavior 的机制。

未来每个有意义的 feature 都应说明：

- 它满足的用户需求；
- 它映射到的 Nakama-style product capability family；
- 用户可判断的 acceptance criteria；
- 证明行为的正向和负向 tests；
- implementation boundary；
- 已运行的 verification commands；
- 已更新哪些 artifacts，方便未来 AI agents 安全继续。

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
  - agent_native_requirement_test_implementation_workflow
```

每个能力族最终都应具备：

- Module 或 runtime subsystem owner。
- 公开行为存在时要有 semantic contract surface。
- 存在 client/server messages 时要有 protocol surface。
- 存在 durable state 时要有 storage boundary。
- 针对 invariants 和 error behavior 的 tests。
- 针对关键边界的 repository checks 或 architecture checks。
- 面向公开读者时同时维护英文与简体中文文档。
- 非平凡用户可见变更应有 requirement、acceptance 和 test-plan 记录。

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

- Source-first alpha 已经可见，prototype-ready foundation execution plan 也已经推进到 storage objects 和第一版 realtime outbound delivery foundation。
- 下一项缺口是在扩大 product modules 前，明确 AI-native requirement-to-test-to-implementation loop。否则 agents 可以继续增加 slice，但用户还不能获得“说出需求，AI 负责 spec、tests、implementation、verification 和 records”的产品承诺。

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

这一阶段吸收 Nakama 对 account/session/socket lifecycle 的显式语义压力，然后再让更高层的 social 或 multiplayer modules 依赖这些生命周期行为。Transport、protocol、application 和 backend service separation 仍是 vibit 自身架构规则；它不再需要 Pitaya 成为 active near-term planning driver。

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
storage_objects_runtime_behavior_implementation: runtime/internal/app/storage/service.go
storage_objects_protocol_route_gate: docs/storage-objects-protocol-route-gate.md
storage_objects_protocol_route_gate_decision: ADR-0118
storage_objects_protocol_route_implementation: proto/vibit/storage/v1/storage.proto
storage_objects_protocol_route_implementation_decision: ADR-0119
storage_objects_protocol_route_local_proof: runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go
storage_objects_protocol_route_local_proof_decision: ADR-0120
first_server_push_realtime_messaging_gate: docs/first-server-push-realtime-messaging-gate.md
first_server_push_realtime_messaging_gate_decision: ADR-0122
first_server_push_realtime_messaging_runtime_slice: runtime/internal/app/realtime/service.go
first_server_push_realtime_messaging_runtime_slice_decision: ADR-0123
next_alpha_direction_after_realtime_runtime_slice_decision: ADR-0124
realtime_protocol_websocket_outbound_delivery_implementation: proto/vibit/realtime/v1/realtime.proto
realtime_protocol_websocket_outbound_delivery_implementation_decision: ADR-0126
agent_native_feature_request_test_workflow: docs/agent-native-feature-request-test-workflow.md
agent_native_feature_request_test_workflow_decision: ADR-0128
next_nakama_prototype_ready_capability_selection_decision: ADR-0132
next_work_item: W-0228 Define agent-native feature request scaffolding gate
```

目标：从开发者可以检查的 source-first alpha，推进到可以用于严肃小型 prototype 的 foundation。

候选工作：

- 降低 local setup、migration 和 configuration friction。已由 `W-0200` 完成。
- 定义并添加更清晰的 example client 或 example app path。已由 `W-0225` 和 `W-0226` 完成。
- 定义超出 inventory proof slice 的 first general storage-object behavior。已由 `W-0201` 完成。
- 定义第一版 storage objects persistence schema posture。已由 `W-0202` 完成。
- 添加第一版 storage objects migration source。已由 `W-0203` 完成。
- 定义 storage objects repository boundary。已由 `W-0204` 完成。
- 实现 storage-neutral storage objects repository interface。已由 `W-0205` 完成。
- 定义并实现 storage objects PostgreSQL adapter。已由 `W-0206` 和 `W-0207` 完成。
- 定义并实现 storage objects runtime behavior。已由 `W-0208` 和 `W-0209` 完成。
- 定义并实现 storage objects protocol route family。已由 `W-0210` 和 `W-0211` 完成。
- 通过 local alpha request flow 证明 storage object routes。已由 `W-0212` 完成。
- 在 storage object local proof 后确认下一项 alpha direction。已由 `W-0213` 完成。
- 定义第一版 server push and realtime messaging gate。已由 `W-0214` 完成。
- 实现第一版 server push and realtime messaging runtime slice。已由 `W-0215` 完成。
- 确认 realtime runtime slice 后的下一项 alpha direction。已由 `W-0216` 完成。
- 实现 realtime protocol and WebSocket outbound delivery slice。已由 `W-0218` 完成。
- 确认 realtime outbound delivery 后的下一项 alpha direction。已由 `W-0219` 完成。
- 定义 agent-native feature request and test workflow。已由 `W-0220` 完成。
- 试点 Nakama-aligned feature request workflow。已由 `W-0221` 完成。
- Harden presence/status local proof through close and offline cases。Next。
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

目标：覆盖 Nakama 期望的 multiplayer 能力族。Routing、room 和 broadcast implementation details 继续由 vibit 自己拥有，直到未来 architecture ADR 重新引入 Pitaya-style distributed vocabulary。

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

目标：只在 single-process semantics 稳定后，并且后续 ADR 决定 Pitaya 是否重新成为 active architecture reference 后，再引入 distributed topology。

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

1. Requirement intake：复述用户需求、用户可见行为、non-goals 和 risks。
2. Nakama reference review：记录正在覆盖的 Nakama product capability；如果没有 Nakama 映射，也要显式说明。
3. Acceptance criteria：实现前定义用户可见的成功和失败行为。
4. Test plan：定义该 slice 需要的 positive、negative、permission、failure-path、persistence、protocol、integration tests。
5. Vibit ownership：决定 module、runtime、platform、generated、operations owners。
6. Semantic contract：定义 commands、queries、events、errors、permissions、invariants。
7. Protocol contract：semantic behavior 稳定后再定义 wire messages。
8. Persistence boundary：在 adapters 前定义 tables、repositories、indexes、redaction。
9. Application behavior：实现最小 vertical slice。
10. Operations surface：当行为需要 operator 时，补 inspection、metrics、admin actions 或 runbook。
11. Verification：补 focused Go tests 和 repository checks，然后运行。
12. Memory：更新 change specs、ADRs、manifests、conversation logs。

首选实现形态仍然是：

```text
user requirement -> spec -> acceptance criteria -> test plan -> tests -> contract -> generated shape -> logic -> checks -> docs
```

## 7. 近期建议

Authenticated failure-path proof、next capability selection、example client path gate、example client path implementation 和 follow-up scaffolding selection 之后，下一个具体工作不应直接跳到 chat、groups、matchmaking、match runtime、SDK publication、hosted demos 或 distributed runtime。下一个具体工作应定义 feature request scaffolding gate：

```text
recommended_next_direction: define_agent_native_feature_request_scaffolding_gate
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
```

理由：

- 项目已经有明确的用户可见 AI-native development workflow，并已通过 presence/status 和 authenticated failure-path proof 试点。
- 现有 source-first alpha capabilities 已有可读的 local example path；下一项 gap 是让新的 user requirements 能 scaffold 成 specs、acceptance criteria、tests、implementation boundaries、verification 和 durable memory，再扩展 broad modules。
- 如果没有 requirement 和 test workflow，后续 feature work 可以技术上正确，但无法兑现“AI 帮用户完成 specification、tests、implementation、verification”的产品承诺。
- Nakama-first product planning 可以避免近期 scope 被产品广度和 Pitaya-style distributed architecture 同时拉扯。
- Pitaya-style cluster/RPC/frontend-backend concerns 应继续延后，直到 single-process behavior 和 AI-native feature workflow 被证明。

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
- 在后续 ADR 重新激活前，把 Pitaya 当成当前 product driver。
- 让 AI 在没有 spec、acceptance criteria、test plan、tests 或明确 verification rationale 的情况下实现非平凡用户需求。

## 9. Agent 规则

Agents 必须：

- 把 Nakama-class 常用能力覆盖当成产品要求，而不是背景阅读。
- 把 AI-native requirement intake、testing、implementation 和 verification 当成产品目的，而不是内部 workflow。
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
- 对非平凡用户可见需求跳过 acceptance criteria 或 test planning。
