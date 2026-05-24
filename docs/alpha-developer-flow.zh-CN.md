# Alpha Developer Flow 中文版

状态：Draft v0.1
最后更新：2026-05-21
范围：vibit v0.1 alpha path 的 packaged local developer journey
说明：本文件是 `docs/alpha-developer-flow.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文档把现有 local alpha entry points 整理成一条 developer journey。它不是 release declaration，也不授权 release publishing、release packaging、hosted deployment、runtime behavior changes、protocol changes、generated output changes、migrations、dependencies、broad operations/admin behavior、product module expansion 或 direct Nakama/Pitaya API compatibility。

## 1. 目的

Local alpha path 现在已经有一个具备技术能力的 contributor 检查 vibit 所需的组件：

- `README.md` 中的 project positioning；
- `docs/runtime-runbook.md` 中的 runtime startup 和 verification notes；
- `examples/local-alpha-request-loop.sh` 中的 redacted request-loop script；
- `/healthz`、`/readyz`、`/version` 和 `/configz` local status endpoints；
- `docs/alpha-acceptance-checklist.md` 中的 acceptance criteria；
- `.arch/work-items.yaml` 中的 continuation state。

本文档把这些组件连接起来，让 contributor 不需要在 project memory 中来回查找也能沿着同一条 sequence 操作。

## 2. 当前 Package State

```text
local_alpha_developer_flow_packaged: true
release_declared: false
release_publishing_authorized_by_this_document: false
release_packaging_authorized_by_this_document: false
release_publishing_decision_gate_defined: true
release_execution_preparation_gate_defined: true
release_execution_authorization_gate_defined: true
release_execution_maintainer_decision_recorded: true
release_identifier_artifact_plan_defined: true
release_execution_final_authorization_recorded: true
proposed_release_identifier: v0.1.0-alpha.1
authorized_release_identifier: v0.1.0-alpha.1
prototype_ready_local_development_path_package_implemented: true
prototype_ready_local_development_path_package: docs/prototype-ready-local-development-path-package.md
storage_objects_protocol_route_gate_defined: true
storage_objects_protocol_route_gate: docs/storage-objects-protocol-route-gate.md
storage_objects_protocol_route_gate_decision: ADR-0118
storage_objects_protocol_route_implementation_completed: true
storage_objects_protocol_route_implementation_decision: ADR-0119
storage_objects_protocol_route_local_proof_completed: true
storage_objects_protocol_route_local_proof_decision: ADR-0120
first_server_push_realtime_messaging_gate_defined: true
first_server_push_realtime_messaging_gate_decision: ADR-0122
first_server_push_realtime_messaging_runtime_slice_completed: true
first_server_push_realtime_messaging_runtime_slice_decision: ADR-0123
next_alpha_direction_after_realtime_runtime_slice_selected: true
next_alpha_direction_after_realtime_runtime_slice_decision: ADR-0124
realtime_protocol_websocket_outbound_delivery_gate_defined: true
realtime_protocol_websocket_outbound_delivery_gate_decision: ADR-0125
realtime_protocol_websocket_outbound_delivery_implementation_completed: true
realtime_protocol_websocket_outbound_delivery_implementation_decision: ADR-0126
next_alpha_direction_after_realtime_outbound_delivery_slice_decision: ADR-0127
agent_native_feature_request_test_workflow_defined: true
agent_native_feature_request_test_workflow_decision: ADR-0128
agent_native_feature_request_test_workflow_standard: docs/agent-native-feature-request-test-workflow.md
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
next_direction: define_friends_relationship_repository_boundary
next_work_status: next_ready
```

Repository 已获得 `v0.1.0-alpha.1` 的 source-first alpha authorization。Packaged flow 已准备好进行 local review，授权 release surface 仍只包括 Git tag、GitHub release record、release notes 和 GitHub source archive。Binaries、packages、containers、checksums、signing/provenance artifacts、hosted deployment、install scripts、registry publication，以及 GitHub release record 之外的 public announcements 继续 deferred。

Prototype-ready local development path package 已记录在 `docs/prototype-ready-local-development-path-package.md`。它的 example 和 local placeholder configuration 入口是 `examples/README.md` 和 `examples/local.prototype.env.example`。

当前 packaged flow 的追溯引用包括 `docs/first-alpha-user-discovery-loop.md`、`docs/first-alpha-feedback-intake-surfaces.md`、`docs/product-maturity-milestones.md`、`docs/prototype-ready-foundation-execution-plan.md`、`docs/prototype-ready-local-development-path-gate.md`、`docs/storage-objects-behavior-gate.md` 和 `docs/storage-objects-persistence-schema-gate.md`。

## 3. 推荐 Journey

1. 阅读 `README.md`，理解 vibit、pre-alpha 状态，以及 agent-native server framework 目标。
2. 阅读 `docs/v0.1-alpha-goal.md`，了解短期 `v0.1 alpha` 目标和长期 Nakama/Pitaya-class 方向。
3. 阅读 `docs/alpha-acceptance-checklist.md`，确认哪些 alpha items 是 ready、manual、deferred 或 blocked。
4. 安装 local prerequisites：Go、Node.js；测试 persistent runtime path 时需要 PostgreSQL；只有重新生成 Protobuf output 时才需要 Buf/Protobuf tooling。
5. 运行 static repository checks：

```bash
node tools/vibit inspect next
node tools/vibit check work --json
node tools/vibit check all --json
```

6. 运行 Go tests：

```bash
cd runtime
go test ./...
```

7. 从 repository root 运行 redacted local alpha request loop：

```bash
examples/local-alpha-request-loop.sh
```

8. 手动启动 server 或评估 PostgreSQL runtime path 时，使用 `docs/runtime-runbook.md`。
9. 评估可重复 source-first local path 时，阅读 `docs/prototype-ready-local-development-path-package.md`。
10. 使用 `.arch/work-items.yaml` 和 `node tools/vibit inspect next` 找到下一步 bounded contribution。

## 4. Runtime Entry Points

当前 process 暴露：

```text
/v1/ws
/healthz
/readyz
/version
/configz
```

`/v1/ws` 是 gameplay WebSocket endpoint。它期望 binary `vibit.protocol.v1.Envelope` Protobuf bytes。该 endpoint 不接受 JSON。

`/healthz`、`/readyz`、`/version` 和 `/configz` 是 local alpha troubleshooting endpoints。`/configz` 报告 redacted runtime posture，并包含 `secrets_redacted: true`。这些 endpoints 不是 production operations API、admin console、metrics backend、gameplay protocol route、release artifact 或 hosted deployment surface。

## 5. Local Proof Flow

Packaged local proof flow 是：

```text
local onboarding
-> device credential login
-> connection binding
-> protected inventory grant/read
-> protected presence query
-> protected own-player storage object put/get/list/delete
-> logout
-> post-logout protected request rejection
```

可执行入口是：

```bash
examples/local-alpha-request-loop.sh
```

该 script 包装 focused Go E2E proof：

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run 'TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout|TestStorageObjectsProtocolRouteLocalAlphaFlow' -v
```

该 proof 使用现有 runtime protocol handlers。它不需要 live PostgreSQL、committed verifier keys、raw credentials、raw access tokens、DSNs、digests 或手写 WebSocket client。

## 6. PostgreSQL Path

PostgreSQL runtime path 是当前 alpha runtime composition，但它有 manual setup requirements：

- 准备 local PostgreSQL database；
- 显式 apply 或 verify SQL migrations；
- 设置 `VIBIT_RUNTIME_STORE=postgres`；
- 设置 `VIBIT_POSTGRES_DSN`；
- 设置全部 authentication verifier key environment variables；
- 不提交 local verifier keys 或 DSNs。

普通 server startup 不会自动 apply migrations。Optional live PostgreSQL verification 仍通过 `VIBIT_POSTGRES_TEST_DSN` 和 disposable database opt-in。

## 7. Redaction Contract

不要记录或提交：

- raw device credential text 或 bytes；
- raw access tokens；
- credential 或 token lookup digests；
- credential 或 token verifier digests；
- HMAC input 或 output bytes；
- verifier key values；
- concrete verifier key set ids；
- 带 credentials 的 PostgreSQL DSNs；
- 可能携带 secrets 的 headers、cookies、query strings、WebSocket subprotocol values 或 remote addresses。

Request-loop script 和 `/configz` surface 都属于该 redaction posture。

## 8. Contribution Entry Point

下一步 contribution path 始终是 machine-readable：

```bash
node tools/vibit inspect next
```

Release execution、first alpha discovery/feedback、prototype-ready execution plan、storage objects path、realtime outbound delivery path、Nakama-first direction、agent-native feature request/test workflow、Nakama-aligned presence/status workflow pilot、presence/status local proof hardening、authenticated gameplay failure-path verification、next Nakama prototype-ready capability selection、local alpha example client path gate、local alpha example client path implementation、feature request scaffolding selection、feature request scaffolding gate、feature request scaffolding implementation、scaffolded Nakama feature request intake pilot、friends relationship lifecycle gate、friends relationship persistence schema gate 和 friends relationship migration source 继续记录在既有 artifacts 中。Agent-native feature request and test workflow 已记录在 `docs/agent-native-feature-request-test-workflow.md` 和 `ADR-0128`，scaffolded intake pilot decision 已记录在 `ADR-0138`，friends lifecycle gate decision 已记录在 `ADR-0139`，friends persistence schema gate decision 已记录在 `ADR-0140`，friends migration source decision 已记录在 `ADR-0141`。下一步 work 是 `W-0234 Define friends relationship repository boundary`；它只应定义 friends relationship repository boundary，同时继续延后 Pitaya、SDK publication、hosted surfaces、broad product/runtime scope、runtime behavior、protocol routes、Protobuf source、generated output、persistence execution、repository implementation、PostgreSQL adapter implementation、event/audit tables 和 direct compatibility，直到后续明确 work item 授权。

## 9. Deferred Work

以下内容继续 deferred，直到后续明确 work item：

- 发布 `v0.1 alpha`；
- 选择 release identifiers for execution；
- 创建 release tags、binaries、archives、containers、packages、checksums、provenance files 或 hosted deployments；
- 添加 public local onboarding protocol route；
- 添加 production signup、external identity providers、password login、account recovery、account merge 或 multi-device linking；
- 添加 broad operations/admin behavior、metrics backend integration 或 production observability；
- 添加 chat、friends、groups、parties、matchmaking、match runtime、SDKs、distributed runtime 或 direct Nakama/Pitaya API compatibility。

## 10. Verification

检查 packaged flow 时使用这组命令：

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./cmd/vibit-server
cd runtime && go test ./...
examples/local-alpha-request-loop.sh
git diff --check
```
