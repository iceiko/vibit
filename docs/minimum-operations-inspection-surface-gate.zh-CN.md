# Minimum Operations Inspection Surface Gate 中文版

状态：Accepted v0.1
最后更新：2026-05-31
范围：在 friends relationship route proof 之后，为第一组最小 source-first operations inspection surface 定义 gate-only boundary
依赖：`decisions/ADR-0151-select-next-nakama-prototype-ready-capability-after-friends-route-proof.md`、`docs/runtime-runbook.md`、`docs/alpha-developer-flow.md`、`docs/alpha-acceptance-checklist.md`、`docs/nakama-pitaya-product-parity-roadmap.md`、`docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0152`

英文版 `docs/minimum-operations-inspection-surface-gate.md` 是权威版本。本文是配套简体中文翻译。

本文定义 minimum operations inspection surface gate。它是 gate artifact。本 slice 不实现 operations/admin endpoints、metrics endpoints、observability pipelines、dashboards、runtime behavior、protocol messages 或 routes、Protobuf source、generated output、repository interfaces、PostgreSQL adapters、migrations、dependencies、authentication/session behavior changes、event/audit tables、groups、parties、chat、stream subscriptions、matchmaking、match runtime、SDK publication、hosted deployments、release artifacts、public announcements、paid promotion、distributed runtime 或 direct Nakama/Pitaya API compatibility。

## 1. Core Rule

Minimum operations inspection surface gate 记录是：

```yaml
minimum_operations_inspection_surface_gate: defined
completed_work_item: W-0244
decision: ADR-0152
check_rule: runtime.minimum_operations_inspection_surface_gate
source_selection_decision: ADR-0151
source_friends_route_proof_decision: ADR-0150
source_workflow_decision: ADR-0128
standard: docs/minimum-operations-inspection-surface-gate.md
translation: docs/minimum-operations-inspection-surface-gate.zh-CN.md
selected_nakama_capability_family: admin_console_metrics_observability_and_operations
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
surface_posture: source_first_local_operations_inspection
admin_console_posture: deferred
metrics_endpoint_posture: deferred
observability_pipeline_posture: deferred
future_implementation_work_item: W-0245
future_implementation_direction: implement_minimum_operations_inspection_source_first_surface
future_cli_inspection_candidate: tools/vibit inspect operations
future_docs_candidate: docs/runtime-runbook.md
future_acceptance_checklist_candidate: docs/alpha-acceptance-checklist.md
accepted_existing_runtime_surfaces:
  - /healthz
  - /readyz
  - /version
  - /configz
accepted_existing_source_surfaces:
  - .arch/work-items.yaml
  - .arch/runtime.yaml
  - docs/runtime-runbook.md
  - docs/alpha-developer-flow.md
  - docs/alpha-acceptance-checklist.md
  - examples/local-alpha-example-client.sh
  - examples/local-alpha-request-loop.sh
minimum_inspectable_state_categories:
  - process_liveness_and_readiness
  - runtime_version_and_release_posture
  - runtime_store_and_configuration_posture
  - protocol_route_family_inventory
  - local_alpha_flow_steps
  - persistence_and_migration_posture
  - authentication_session_connection_posture
  - generated_output_and_proto_posture
  - repository_check_and_verification_posture
  - deferred_operations_surfaces
redaction_required: true
runtime_behavior_added_by_this_gate: false
operations_admin_endpoint_added: false
admin_console_added: false
metrics_endpoint_added: false
observability_pipeline_added: false
dashboard_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
migration_added: false
dependency_added: false
authentication_session_behavior_changed: false
event_audit_table_added: false
hosted_deployment_added: false
sdk_added: false
distributed_runtime_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. 目的

`ADR-0151` 在 protected friends relationship route family 已通过本地 proof 之后，选择 `admin_console_metrics_observability_and_operations` 作为下一项 Nakama-first prototype-ready capability family。

当前 local alpha 已暴露小型 HTTP troubleshooting endpoints：

```text
/healthz
/readyz
/version
/configz
```

它也可以通过 architecture manifests、runbooks、example scripts、repository checks 和 focused Go tests 进行 source-first inspection。缺口是：还没有一个清晰的 minimum operations inspection posture，告诉 prototype authors 哪些状态可检查、哪些必须 redacted、哪些 operations surface 仍是未来工作。

第一步 operations inspection 因此应保持 source-first 和 local。它应帮助 developer 或 agent 理解当前 server state 和 verification posture，但不创建 admin console、metrics backend、telemetry pipeline、sensitive state dump 或 compatibility promise。

## 3. 最小可检查类别

未来 implementation 应用 source-first 方式暴露或记录以下类别：

```yaml
minimum_inspectable_state_categories:
  process_liveness_and_readiness:
    existing_runtime_surface:
      - /healthz
      - /readyz
  runtime_version_and_release_posture:
    existing_runtime_surface:
      - /version
  runtime_store_and_configuration_posture:
    existing_runtime_surface:
      - /configz
  protocol_route_family_inventory:
    existing_source_surface:
      - runtime route constants
      - Protobuf bridge registrations
      - docs/runtime-runbook.md
  local_alpha_flow_steps:
    existing_source_surface:
      - examples/local-alpha-example-client.sh
      - examples/local-alpha-request-loop.sh
      - runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go
  persistence_and_migration_posture:
    existing_source_surface:
      - runtime/migrations/postgres/
      - docs/runtime-runbook.md
  authentication_session_connection_posture:
    existing_source_surface:
      - docs/runtime-runbook.md
      - .arch/runtime.yaml
  generated_output_and_proto_posture:
    existing_source_surface:
      - proto/
      - runtime/internal/generated/proto/
      - docs/generated-output.md
  repository_check_and_verification_posture:
    existing_source_surface:
      - tools/vibit
      - rules/check-rules.json
  deferred_operations_surfaces:
    allowed_future_source_inspection:
      - admin console、metrics endpoint、observability pipeline、live dashboards、player/session/token inspectors、hosted operations、SDKs 和 distributed runtime 继续 deferred
```

第一项 implementation 应优先选择 `tools/vibit inspect operations` 或等价的 source-first repository inspection，而不是新增 runtime endpoints。

## 4. Redaction

任何未来 inspection surface 都不得 print、persist、log、serialize 或 record：

- raw device credential material；
- raw access tokens；
- credential 或 token lookup digests；
- credential 或 token verifier digests；
- verifier key values；
- concrete verifier key set ids；
- HMAC inputs 或 outputs；
- 带 credentials 的 PostgreSQL DSNs；
- passwords、connection strings 或 local secret file contents；
- HTTP headers、cookies、query strings、WebSocket subprotocol values、remote addresses 或 concrete transport metadata；
- 完整 player、session、token、connection、relationship、storage 或 presence identifiers，除非后续 redaction policy 明确认定某个 surface 中安全；
- database row payloads 或 arbitrary JSON object values。

允许的 public output 包括 route names、endpoint names、file paths、high-level capability status、redacted configuration posture、verification command names 和 broad state categories。

## 5. Ownership

未来 implementation ownership：

```yaml
source_first_inspection_owner: tools/vibit
runtime_status_endpoint_owner: runtime/cmd/vibit-server
runbook_owner: docs/runtime-runbook.md
alpha_flow_owner: docs/alpha-developer-flow.md
acceptance_checklist_owner: docs/alpha-acceptance-checklist.md
architecture_memory_owner:
  - .arch/work-items.yaml
  - .arch/runtime.yaml
  - .arch/reference.yaml
runtime_behavior_owner: unchanged
protocol_owner: unchanged
persistence_owner: unchanged
```

规则：

- `tools/vibit` 可以总结已提交 source、manifests、docs 和 check status。
- Runtime endpoints 在后续 bounded implementation 明确改变前，仍限定为现有 local-alpha troubleshooting surface。
- Runtime behavior、protocol payloads、persistence、migrations、repository interfaces、adapters、authentication/session semantics 和 generated output 继续由既有 boundaries 管理。
- 任何 domain module 都不应成为 broad operations inspection owner。

## 6. Nakama 和 Pitaya 映射

Nakama reference mapping：

- 本 gate 覆盖 `admin_console_metrics_observability_and_operations` capability family。
- 它采纳 backend developer 在评估 framework 时需要检查 health、version、configuration posture、feature availability 和 operational state 的产品压力。
- 它不复制 Nakama console routes、REST paths、metrics names、dashboard behavior、SDK shapes、account/session/player inspectors 或 compatibility promises。

Pitaya reference mapping：

- Pitaya 继续作为 future distributed architecture reference deferred。
- 本 gate 不引入 frontend/backend server roles、RPC、service discovery、cluster metrics、distributed session inspection 或 distributed operations behavior。

## 7. 未来 Implementation Work

打开：

```text
M-173/W-0245 Implement minimum operations inspection source-first surface
```

未来 work item 可以：

- 添加 source-first `tools/vibit inspect operations` command 或等价 inspection subcommand；
- 更新 `docs/runtime-runbook.md` 及其翻译，描述 minimum inspection workflow；
- 更新 `docs/alpha-developer-flow.md` 和 checklist references；
- 添加 repository checks，验证 redaction 和 accepted categories；
- 总结现有 local alpha routes、endpoints、manifests 和 verification posture。

未来 work item 不得：

- 添加 admin console、metrics endpoint、observability pipeline、dashboard、hosted operations surface、player/session/token inspector endpoint 或 database-state dump；
- 添加 runtime behavior、新 HTTP endpoint behavior、protocol messages 或 routes、Protobuf source、generated output、migrations、dependencies、repository interfaces、PostgreSQL adapters、authentication/session behavior changes、event/audit tables、SDK publication、hosted deployments、release artifacts、distributed runtime 或 direct Nakama/Pitaya API compatibility。

## 8. Verification Expectations

未来 implementation 应验证：

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.minimum_operations_inspection_surface_gate
node tools/vibit check change define-minimum-operations-inspection-surface-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

如果未来 implementation 修改 Go runtime behavior、protocol code 或 tests，还必须运行 focused Go tests 和 `cd runtime && go test ./...`。本 gate 不添加 Go behavior，因此不要求 Go tests。

## 9. Stop Conditions

如果 implementation 需要以下内容，请停止并创建独立 gate：

- 新 HTTP operations/admin endpoint；
- metrics、tracing、logging 或 observability backend dependencies；
- live player、session、token、connection、relationship 或 storage inspectors；
- database row inspection 或 arbitrary state dumps；
- public 或 hosted admin UI behavior；
- WebSocket 或 Protobuf operations routes；
- authentication/session/route protection behavior changes；
- startup wiring 或 process lifecycle changes；
- event/audit tables；
- generated client libraries、SDK publication、hosted deployment、release artifacts、public announcements、paid promotion、distributed runtime 或 direct Nakama/Pitaya API compatibility。
