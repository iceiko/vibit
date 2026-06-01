# friends 模块 Agent 指南

状态：Draft v0.1

## 什么时候使用本模块

本模块用于 friends relationship social graph 的 repository vocabulary 和 storage-neutral value types。

当前已实现的 slice 刻意保持狭窄：

- `runtime/internal/modules/friends.Repository`
- `FriendRelationship`、canonical unordered pair、actor、block state、lifecycle state、public status 和 version value types
- request、accept、reject、remove、block、unblock、pair lookup、player-scoped list input/result types
- optimistic conflict classes 和 redacted repository errors
- normalization helpers 和 focused Go tests

`M-163 Friends Relationship Repository Interface Implementation` 已由 `W-0235` 完成。检查规则是 `runtime.friends_relationship_repository_interface_implementation`。

`M-164 Friends Relationship PostgreSQL Adapter Gate` 已由 `W-0236` 完成。检查规则是 `runtime.friends_relationship_postgresql_adapter_gate`。

`M-165 Friends Relationship PostgreSQL Adapter Implementation` 已由 `W-0237` 完成。检查规则是 `runtime.friends_relationship_postgresql_adapter_implementation`。

`M-166 Friends Relationship Runtime Behavior Gate` 已由 `W-0238` 完成。检查规则是 `runtime.friends_relationship_runtime_behavior_gate`。

`M-167 Friends Relationship Runtime Behavior Implementation` 已由 `W-0239` 完成。检查规则是 `runtime.friends_relationship_runtime_behavior_implementation`。

`M-168 Friends Relationship Protocol Route Gate` 已由 `W-0240` 完成。检查规则是 `runtime.friends_relationship_protocol_route_gate`。

`M-169 Friends Relationship Protocol Route Implementation` 已由 `W-0241` 完成。检查规则是 `runtime.friends_relationship_protocol_route_implementation`。

`M-170 Friends Relationship Protocol Route Local Proof` 已由 `W-0242` 完成。检查规则是 `runtime.friends_relationship_protocol_route_local_proof`。

repository 下一项 work item 是 `W-0268 Define Pitaya-aligned session binding, kick/disconnect, and session data boundary gate`。该 work 不属于本模块且必须保持 gate-only；在后续 bounded work item 明确授权前，不要在本模块添加 protocol shape changes、repository interface changes、PostgreSQL adapter changes、migrations、dependencies、event/audit tables、groups、parties、chat、matchmaking、match runtime、SDK publication、hosted deployments、session binding behavior、kick/disconnect behavior、session data behavior 或 persistence、acceptor behavior、TCP acceptors、WebSocket behavior changes、connection lifecycle behavior changes、route handler implementation、handler routing behavior、handler pipeline behavior、serializer behavior、message forwarding behavior、backend route targeting、distributed runtime behavior、frontend/backend server role implementation、server-to-server RPC behavior、remote calls、service discovery implementation、service registries、service selectors、distributed groups、group membership registries、stream subscriptions、room broadcast fanout、broadcast delivery guarantees、cluster-safe session routing behavior、session location registries、connection owner node registries、routing epoch behavior、session route targets、remote connection handoff、distributed session routing、operations/admin implementation 或 direct Nakama/Pitaya API compatibility。

## 什么时候不要使用本模块

不要用本模块处理：

- WebSocket、HTTP、Protobuf 或 generated wire behavior。
- 本模块下的 PostgreSQL adapter implementation 或 SQL execution。
- runtime friend request、accept、reject、remove、block、unblock、list 或 status behavior。
- player account lifecycle。
- authentication、token formats、credential storage 或 session validation。
- storage objects、inventory、presence、realtime messaging、chat、groups、parties、matchmaking 或 match runtime。
- friendship history 的 event/audit tables。
- direct Nakama 或 Pitaya public API compatibility。

如果需求涉及这些概念，应创建或更新对应 owner boundary，而不是在本模块隐藏 ownership。

## 扩展点

- Repository interface：`runtime/internal/modules/friends.Repository`
- Repository value types：`FriendRelationship`、`FriendRelationshipPair`、`FriendRelationshipActor`、`FriendRelationshipBlockState`、`FriendRelationshipVersion`
- Lifecycle vocabulary：`pending`、`friends`、`rejected`、`removed`
- Public status vocabulary：`pending`、`friends`、`blocked`、`ended`
- Normalizers：relationship records、list results、pair identity、actors、block state 和 repository inputs
- Tests：`runtime/internal/modules/friends/repository_test.go`
- PostgreSQL adapter owner candidate：`runtime/internal/platform/persistence/postgres`
- PostgreSQL adapter gate：`docs/friends-relationship-postgresql-adapter-gate.md`
- PostgreSQL adapter implementation：`runtime/internal/platform/persistence/postgres/friend_relationship_repository.go`
- PostgreSQL adapter tests：`runtime/internal/platform/persistence/postgres/friend_relationship_repository_test.go`
- Runtime behavior gate：`docs/friends-relationship-runtime-behavior-gate.md`
- Runtime behavior implementation：`runtime/internal/app/friends/service.go`
- Runtime behavior tests：`runtime/internal/app/friends/service_test.go`
- Protocol route gate：`docs/friends-relationship-protocol-route-gate.md`
- Protocol route implementation：`runtime.friends_relationship_protocol_route_implementation`
- Protocol route local proof：`W-0242 Prove friends relationship protocol route in local alpha request flow`
- Pitaya-aligned acceptor and connection lifecycle source-first map：`W-0266 Implement Pitaya-aligned acceptor and connection lifecycle source-first map`
- Pitaya-aligned session binding、kick/disconnect 和 session data boundary gate：`W-0268 Define Pitaya-aligned session binding, kick/disconnect, and session data boundary gate`

第一批 public runtime commands 和 queries 仍然 deferred。未来 runtime behavior 必须先从 validated request identity 派生 actor identity，再调用 repository interface；client-supplied player ids 不是 authentication proof。

## 禁止的捷径

- 不要绕过 `module.yaml` 声明的边界。
- 不要添加未注册的 public commands、queries、events、errors 或 permissions。
- 不要在本模块下添加 PostgreSQL adapter code。
- 不要在本模块导入 `pgx`、`database/sql`、WebSocket packages、generated Protobuf packages、SDK packages 或 distributed runtime packages。
- 不要在 friends module source 中执行 SQL 或写入具体 SQL statements。
- 不要从本模块修改 migrations。
- 不要从本模块接线 runtime handlers、startup composition、route policy、protocol adapters 或 transport behavior。
- 没有后续 bounded work item 时，不要添加 protocol routes、Protobuf sources 或 generated output。
- 不要在 friends value types 中存放 raw credentials、raw tokens、verifier material、lookup digests、verifier digests、cookies、transport subprotocols、connection metadata、chat、group、party、matchmaking、match runtime、Pitaya server routing 或 direct Nakama/Pitaya compatibility fields。
- 不要把 pair member ids、actor ids、`player_id`、`session_id` 或 transport metadata 当成 authenticated proof。

## 必需测试

见 `module.yaml` 中的 `tests.required`。

当前 repository interface slice 的测试必须覆盖：

- Repository interface storage neutrality。
- lifecycle state、public status 和 conflict vocabulary 的 closed set。
- canonical unordered pair normalization 和 self-relationship rejection。
- returned record normalization。
- returned list result normalization 和 slice copying。
- send/block/unblock self-target rejection。
- accept/reject/remove actor-in-pair requirements。
- expected version validation 和 pointer copying。
- list pagination bounds。
- redacted conflict and repository errors。
- 不存在 secret、transport、protocol、distributed、chat、group、party、match 和 direct compatibility fields。

修改 friends runtime source 后运行 `node tools/vibit check runtime`。Go 可用时也运行 `cd runtime && go test ./...`。

## Repository Continuation

`ADR-0181` 此前完成 `W-0273 Select next Pitaya-aligned direction after runtime observability map`，并选择 `define_pitaya_aligned_metrics_tracing_boundary_gate` 作为 follow-up。这仍是 repository historical context，不会新增 friends module scope。

`ADR-0183` 已注册 `runtime.pitaya_aligned_metrics_tracing_source_first_map`，完成 `W-0275 Implement Pitaya-aligned metrics and tracing source-first map`，实现 `node tools/vibit inspect pitaya-metrics-tracing --json`，并打开 `W-0276 Select next Pitaya-aligned direction after metrics and tracing map` 作为 repository next work item。`node tools/vibit inspect pitaya-observability --json` inspection 仍可用。该 work 仍不属于 friends module。

`ADR-0184` 已注册 `runtime.next_pitaya_aligned_direction_after_metrics_tracing_map`，完成 `W-0276 Select next Pitaya-aligned direction after metrics and tracing map`，选择 `define_pitaya_aligned_dashboard_admin_operations_boundary_gate`，并打开 `W-0277 Define Pitaya-aligned dashboard and admin operations boundary gate` 作为 repository next work item。该 work 仍不属于 friends module，也不添加 dashboard、admin、runtime endpoint、protocol、generated output、persistence、dependency、distributed runtime 或 direct compatibility scope。
