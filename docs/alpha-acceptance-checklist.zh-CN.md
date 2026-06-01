# Alpha Acceptance Checklist 中文版

状态：Draft v0.1
最后更新：2026-05-21
范围：vibit 的本地 v0.1 alpha acceptance criteria
说明：本文件是 `docs/alpha-acceptance-checklist.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本 checklist 不是 release declaration。它是给 maintainer 和 contributor 使用的验收清单，用于判断本地 `v0.1 alpha` developer flow 是否已经准备好进行 local inspection，并供后续考虑 publishing。

## 1. 目的

第一版 alpha 在发布前应先做到 developer-usable。本 checklist 定义本地 flow 要被接受为 alpha-ready 需要满足的条件：

- contributor 能理解 vibit 是什么；
- 能准备 local prerequisites 和 configuration；
- 能运行或验证当前 runtime；
- 能理解 authenticated gameplay flow；
- 能运行 repository checks；
- 能确认 secrets 仍被 redacted；
- 能找到下一步 contribution entry point。

本 checklist 不授权 release publishing、release packaging、runtime behavior changes、protocol changes、generated output changes、migrations、dependencies、broad operations/admin behavior、product module expansion 或 direct Nakama/Pitaya API compatibility。

## 2. Checklist 状态

Review alpha flow 时使用这些状态：

- `Ready`：repository 已包含所需 artifact 或 behavior，并且由 tests、checks 或明确 documentation 验证。
- `Manual`：该步骤需要 developer 执行 local setup，例如安装工具或 apply migrations。
- `Deferred`：该项有意不属于 first alpha，或等待后续 work item。
- `Blocked`：该项解决前不能接受 alpha。

## 3. Repository Intake

- [x] `README.md` 和 `README.zh-CN.md` 说明 vibit 是 agent-native server framework，并说明当前状态是 pre-alpha。
- [x] `docs/v0.1-alpha-goal.md` 和 `docs/v0.1-alpha-goal.zh-CN.md` 定义短期 `v0.1 alpha` 目标。
- [x] `AGENTS.md`、`AGENTS.zh-CN.md`、`runtime/AGENTS.md` 和 `runtime/AGENTS.zh-CN.md` 指向当前 continuation queue。
- [x] `.arch/work-items.yaml` 记录当前 continuation state。
- [x] `node tools/vibit inspect next` 能识别当前 `W-0258 Select next Pitaya-aligned direction after cluster-safe session routing map` continuation step。
- [x] `node tools/vibit inspect operations --json` 会报告 source-first local operations posture、route families、redaction flags 和 Pitaya deferred architecture mapping。
- [x] `node tools/vibit inspect pitaya-vocabulary --json` 会报告 source-first Pitaya vocabulary map 和 deferrals。
- [x] `node tools/vibit inspect pitaya-roles --json` 会报告 source-first Pitaya frontend/backend role map 和 deferrals。
- [x] `node tools/vibit inspect pitaya-rpc --json` 会报告 source-first Pitaya server-to-server RPC map 和 deferrals。
- [x] `node tools/vibit inspect pitaya-discovery --json` 会报告 source-first Pitaya service discovery map 和 deferrals。
- [x] `node tools/vibit inspect pitaya-groups --json` 会报告 source-first Pitaya distributed group and broadcast map 和 deferrals。
- [x] `node tools/vibit inspect pitaya-sessions --json` 会报告 source-first Pitaya cluster-safe session routing map 和 deferrals。
- [x] 当前 `next_direction: select_next_pitaya_aligned_direction_after_cluster_safe_session_routing_map` 接续已完成的 friends relationship protocol route local proof、W-0243 post-proof selection、W-0244 minimum operations inspection surface gate、W-0245 source-first operations inspection implementation、W-0246 Pitaya vocabulary gate、W-0247 source-first vocabulary map、W-0248 frontend/backend role boundary gate、W-0249 source-first role map、W-0250 server-to-server RPC boundary gate、W-0251 server-to-server RPC source-first map、W-0252 service discovery boundary gate、W-0253 service discovery source-first map、W-0254 distributed group/broadcast boundary gate、W-0255 distributed group/broadcast source-first map、W-0256 cluster-safe session routing boundary gate 和 W-0257 cluster-safe session routing source-first map。
- [x] `docs/prototype-ready-local-development-path-package.md` 已打包可重复 source-first local path。

## 4. Local Prerequisites

- [ ] Manual：安装 Go，用于 runtime tests 和 local server execution。
- [ ] Manual：安装 Node.js，用于 `tools/vibit` checks。
- [ ] Manual：测试 persistent runtime path 时可用 PostgreSQL。
- [ ] Manual：重新生成 Protobuf output 时可用 Buf 和 Protobuf tooling。
- [x] `docs/runtime-runbook.md` 已记录 memory 和 PostgreSQL startup paths。
- [x] 默认 repository checks 不要求 live PostgreSQL。

## 5. Configuration And Secret Handling

- [x] `VIBIT_RUNTIME_STORE=memory` 仍是默认 bootstrap path。
- [x] `VIBIT_RUNTIME_STORE=postgres` 已记录为当前 alpha runtime composition。
- [x] PostgreSQL runtime startup 要求 `VIBIT_POSTGRES_DSN` 和 verifier key environment variables。
- [x] `examples/local.prototype.env.example` 提供 placeholder-only local configuration checklist。
- [x] `.gitignore` 会忽略 `.vibit.local.env`、`.env.local` 和 `.env.*.local`。
- [x] Verifier key material、raw device credentials、raw access tokens、DSNs、digests、headers、cookies、query strings、subprotocol values、remote addresses 和 concrete transport metadata 都被记录为不是 log-safe。
- [x] Local alpha request-loop script 不打印 raw credentials、raw tokens、verifier keys、DSNs、digests 或 transport metadata。
- [x] `/configz` 只报告 redacted posture，并包含 `secrets_redacted: true`。

## 6. Database And Migration Posture

- [x] PostgreSQL migration sources 位于 `runtime/migrations/postgres`。
- [x] 普通 runtime startup 不会自动 apply migrations。
- [x] Repository runtime tooling 已提供 migration apply/status tooling。
- [ ] Manual：使用 `VIBIT_RUNTIME_STORE=postgres` 前准备 local PostgreSQL database。
- [ ] Manual：在 fresh PostgreSQL database 上使用前，apply 或 verify required SQL migrations。
- [ ] Manual：optional live PostgreSQL verification 使用 `VIBIT_POSTGRES_TEST_DSN` 和 disposable database。

## 7. Runtime Surface

- [x] Gameplay WebSocket endpoint 是 `/v1/ws`。
- [x] `/v1/ws` 期望 binary `vibit.protocol.v1.Envelope` Protobuf bytes，不接受 JSON。
- [x] `/healthz` 报告 process health。
- [x] `/readyz` 报告 readiness posture、runtime store 和 WebSocket path。
- [x] `/version` 报告 pre-alpha runtime version。
- [x] `/configz` 只报告 redacted runtime posture。
- [x] 这些 HTTP status endpoints 是 local troubleshooting surfaces，不是 production operations API、admin console、metrics backend 或 gameplay protocol route。

## 8. Authenticated Gameplay Flow

- [x] Local onboarding 作为 application service behavior 存在：`OnboardLocalPlayerWithDeviceCredential`。
- [x] Local onboarding 没有暴露为 public WebSocket、Protobuf、HTTP、CLI 或 startup auto-creation surface。
- [x] Device credential login 通过 `runtime.authentication.AuthenticateWithDeviceCredential` protocol route 暴露。
- [x] Login 返回 opaque access token 和 runtime session metadata。
- [x] First-message connection binding 通过 `runtime.authentication.BindConnection` protocol route 暴露。
- [x] Protected inventory grant/read requests 使用 `AuthenticatedRequest`。
- [x] Protected presence query 可通过 `runtime.presence.GetPlayerPresence` 使用。
- [x] Protected own-player storage object put/get/list/delete 已通过 `storage.GetOwnStorageObject`、`storage.ListOwnStorageObjects`、`storage.PutOwnStorageObject` 和 `storage.DeleteOwnStorageObject` 证明。
- [x] Logout 通过 `runtime.authentication.LogoutAccessToken` protocol route 暴露。
- [x] 使用同一 token 的 logout 后 protected request 会被拒绝。
- [x] Focused authenticated gameplay E2E test 证明 onboarding -> login -> connection binding -> protected inventory -> presence query -> logout -> post-logout rejection。

## 9. Verification Commands

除特别说明外，从 repository root 运行：

```bash
node tools/vibit inspect next
node tools/vibit inspect operations --json
node tools/vibit check work --json
node tools/vibit check all --json
git diff --check
cd runtime && go test ./...
examples/local-alpha-request-loop.sh
```

Optional focused checks：

```bash
cd runtime && go test ./cmd/vibit-server
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
cd runtime && go test ./internal/platform/protocol/protobuf -run TestStorageObjectsProtocolRouteLocalAlphaFlow -v
```

Optional live PostgreSQL verification 仍是 opt-in，并要求 disposable database：

```bash
cd runtime && VIBIT_POSTGRES_TEST_DSN='postgres://user:pass@127.0.0.1:5432/vibit_test?sslmode=disable' VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1 go test ./internal/platform/protocol/protobuf -run TestPostgresPersistentInventoryRequestLoop -v
```

## 10. Contribution Entry Point

Alpha developer flow 现在已经 packaged 在：

```text
docs/alpha-developer-flow.md
```

Release、discovery、prototype-ready、storage objects、realtime outbound delivery、Nakama-first direction、agent-native feature request/test workflow、Nakama-aligned presence/status pilot、presence/status local proof hardening、authenticated gameplay failure-path verification、next Nakama prototype-ready capability selection、local alpha example client path gate、local alpha example client path implementation、feature request scaffolding、scaffolded Nakama feature request intake pilot、friends relationship lifecycle gate、friends relationship persistence schema gate、friends relationship migration source、friends relationship repository boundary、friends relationship repository interface implementation、friends relationship PostgreSQL adapter gate、friends relationship PostgreSQL adapter implementation、friends relationship runtime behavior gate、friends relationship runtime behavior implementation、friends relationship protocol route gate、friends relationship protocol route implementation、friends relationship protocol route local proof、friends route proof 后的 next Nakama prototype-ready capability selection、minimum operations inspection surface gate、minimum operations inspection source-first surface implementation、Pitaya-aligned distributed runtime vocabulary reactivation gate、Pitaya-aligned distributed runtime vocabulary source-first map、Pitaya-aligned frontend/backend role boundary gate、Pitaya-aligned frontend/backend role source-first map、Pitaya-aligned server-to-server RPC boundary gate、Pitaya-aligned server-to-server RPC source-first map、Pitaya-aligned service discovery boundary gate、Pitaya-aligned service discovery source-first map、Pitaya-aligned distributed group and broadcast boundary gate、Pitaya-aligned distributed group and broadcast source-first map、Pitaya-aligned cluster-safe session routing boundary gate 和 Pitaya-aligned cluster-safe session routing source-first map artifacts 继续记录在既有文档中。Friends route local proof 已记录在 `ADR-0150`，检查规则是 `runtime.friends_relationship_protocol_route_local_proof`；post-proof selection 已记录在 `ADR-0151`，检查规则是 `runtime.next_nakama_prototype_ready_capability_after_friends_route_proof`；operations gate 已记录在 `ADR-0152`，检查规则是 `runtime.minimum_operations_inspection_surface_gate`；operations inspection implementation 已记录在 `ADR-0153`，检查规则是 `runtime.minimum_operations_inspection_source_first_surface_implementation`；Pitaya vocabulary gate 已记录在 `ADR-0154`，检查规则是 `runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate`；source-first Pitaya vocabulary map 已记录在 `ADR-0155`，检查规则是 `runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map`；frontend/backend role boundary gate 已记录在 `ADR-0156`，检查规则是 `runtime.pitaya_aligned_frontend_backend_role_boundary_gate`；frontend/backend role source-first map 已记录在 `ADR-0157`，检查规则是 `runtime.pitaya_aligned_frontend_backend_role_source_first_map`，命令是 `node tools/vibit inspect pitaya-roles --json`；server-to-server RPC boundary gate 已记录在 `ADR-0158`，检查规则是 `runtime.pitaya_aligned_server_to_server_rpc_boundary_gate`；server-to-server RPC source-first map 已记录在 `ADR-0159`，检查规则是 `runtime.pitaya_aligned_server_to_server_rpc_source_first_map`，命令是 `node tools/vibit inspect pitaya-rpc --json`；service discovery boundary gate 已记录在 `ADR-0160`，检查规则是 `runtime.pitaya_aligned_service_discovery_boundary_gate`；service discovery source-first map 已记录在 `ADR-0161`，检查规则是 `runtime.pitaya_aligned_service_discovery_source_first_map`，命令是 `node tools/vibit inspect pitaya-discovery --json`；distributed group and broadcast boundary gate 已记录在 `ADR-0162`，检查规则是 `runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate`；distributed group and broadcast source-first map 已记录在 `ADR-0163`，检查规则是 `runtime.pitaya_aligned_distributed_group_broadcast_source_first_map`，命令是 `node tools/vibit inspect pitaya-groups --json`；cluster-safe session routing boundary gate 已记录在 `ADR-0164`，检查规则是 `runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate`；cluster-safe session routing source-first map 已记录在 `ADR-0165`，检查规则是 `runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map`，命令是 `node tools/vibit inspect pitaya-sessions --json`。下一步 work 是 `W-0258 Select next Pitaya-aligned direction after cluster-safe session routing map`；next_direction: select_next_pitaya_aligned_direction_after_cluster_safe_session_routing_map。

Alpha acceptance 的追溯引用包括 `docs/first-alpha-user-discovery-loop.md`、`docs/first-alpha-feedback-intake-surfaces.md`、`docs/product-maturity-milestones.md`、`docs/prototype-ready-foundation-execution-plan.md`、`docs/prototype-ready-local-development-path-gate.md`、`docs/storage-objects-behavior-gate.md` 和 `docs/storage-objects-persistence-schema-gate.md`。

## 11. Release Deferrals

以下内容继续 deferred，直到后续明确 work item 授权：

- 创建 GitHub source archive 之外的 additional release artifacts。
- 创建 release binaries、packages、containers、checksums、provenance files、signing artifacts、install scripts、registry publications 或 hosted deployments。
- 添加 public local onboarding protocol route。
- 添加 production signup、external identity providers、password login、account recovery、account merge 或 multi-device linking。
- 添加 broad operations/admin behavior、metrics backend integration 或 production observability。
- 添加 chat、friends、groups、parties、matchmaking、match runtime、SDKs、distributed runtime 或 direct Nakama/Pitaya API compatibility。

## 12. Current Acceptance Result

Local alpha flow 已经可以检查，并已 packaged 供 local developer review；但 repository 仍是 pre-alpha，直到后续 release-publishing work item 明确声明。

当前结果：

```text
local_alpha_flow_checkable: true
local_alpha_developer_flow_packaged: true
release_publishing_decision_gate_defined: true
release_execution_preparation_gate_defined: true
release_execution_authorization_gate_defined: true
release_execution_maintainer_decision_recorded: true
release_identifier_artifact_plan_defined: true
release_execution_final_authorization_recorded: true
proposed_release_identifier: v0.1.0-alpha.1
authorized_release_identifier: v0.1.0-alpha.1
release_declared: true
release_publishing_authorized_by_this_checklist: true
prototype_ready_local_development_path_package_implemented: true
storage_objects_protocol_route_gate_defined: true
storage_objects_protocol_route_implementation_completed: true
storage_objects_protocol_route_local_proof_completed: true
first_server_push_realtime_messaging_gate_defined: true
first_server_push_realtime_messaging_runtime_slice_completed: true
realtime_protocol_websocket_outbound_delivery_gate_defined: true
realtime_protocol_websocket_outbound_delivery_implementation_completed: true
agent_native_feature_request_test_workflow_defined: true
operations_inspection_source_first_surface_implemented: true
pitaya_deferred_architecture_map_recorded: true
next_direction: select_next_pitaya_aligned_direction_after_cluster_safe_session_routing_map
next_work_status: next_ready
```
