# Agent-Native Feature Request Scaffolding Gate

状态：Accepted v0.1
最后更新：2026-05-24
范围：为未来 feature request scaffolding 定义 gate-only boundary，使用户后端需求能在实现前转成 bounded change artifacts
依赖：`docs/agent-native-feature-request-test-workflow.md`、`docs/change-spec.md`、`docs/workflow.md`、`docs/nakama-pitaya-product-parity-roadmap.md`、`decisions/ADR-0135-agent-native-feature-request-scaffolding-selection.md`
权威决策：`ADR-0136`

说明：本文件是 `docs/agent-native-feature-request-scaffolding-gate.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文定义 agent-native feature request scaffolding gate。它是 gate artifact。它不实现 scaffolding，不添加 templates，不添加 `tools/vibit` scaffold commands，不改变 runtime behavior，不添加 protocol routes，不添加 Protobuf source，不改变 generated output，不添加 migrations、dependencies、persistence，不改变 startup wiring，不改变 authentication/session behavior，不添加 delivery guarantees、stream subscriptions、chat rooms、groups、broadcast fanout、matchmaking、match runtime、operations/admin behavior，不发布 SDK，不生成 client libraries，不添加 hosted deployments、release artifacts、public announcements、paid promotion、Pitaya-style distributed architecture 或 direct Nakama/Pitaya API compatibility。

## 1. 核心规则

Agent-native feature request scaffolding gate 记录为：

```yaml
agent_native_feature_request_scaffolding_gate: defined
completed_work_item: W-0228
decision: ADR-0136
check_rule: runtime.agent_native_feature_request_scaffolding_gate
source_selection_decision: ADR-0135
source_workflow_decision: ADR-0128
standard: docs/agent-native-feature-request-scaffolding-gate.md
translation: docs/agent-native-feature-request-scaffolding-gate.zh-CN.md
selected_nakama_capability_family: agent_native_requirement_test_implementation_workflow
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
scaffolding_posture: source_first_change_artifact_scaffolding
future_scaffold_command_candidate: tools/vibit scaffold feature
future_template_directory_candidate: changes/_template/feature-request/
future_implementation_work_item: W-0229
future_implementation_direction: implement_agent_native_feature_request_scaffolding
required_scaffold_artifacts:
  - request.md
  - spec.yaml
  - impact.md
  - plan.md
  - checklist.md
  - verification.md
implementation_added_by_this_gate: false
scaffolding_implementation_added: false
template_files_added: false
scaffold_command_added: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
persistence_added: false
startup_wiring_added: false
authentication_session_behavior_changed: false
delivery_guarantees_added: false
stream_subscription_added: false
chat_added: false
group_messaging_added: false
broadcast_fanout_added: false
matchmaking_added: false
match_runtime_added: false
operations_admin_behavior_added: false
sdk_added: false
generated_client_library_added: false
hosted_deployment_added: false
release_artifact_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. 目的

`ADR-0135` 在 local alpha example path 之后选择了 `agent_native_requirement_test_implementation_workflow` 作为下一项 Nakama-first prototype-ready capability family。现有 workflow standard 已定义一个良好的 feature change 应包含什么，但未来用户可见的产品闭环还需要更可执行的 intake path。

未来 scaffolding 应让 agent 从用户需求出发，在实现开始前创建正确的 source-first change artifact 形状。scaffold 必须强化产品目的：

```text
user requirement
-> bounded request record
-> spec with Nakama capability mapping
-> acceptance criteria
-> test plan
-> implementation boundaries
-> verification plan
-> durable memory expectations
```

本 gate 故意保持窄范围。它只定义 posture 和 stop conditions，不实现 scaffold command、template directory 或 check。

## 3. 选定的 Scaffolding 形状

第一版接受的形状是 source-first repository-local change artifact scaffolding。

未来 implementation candidates：

```text
tools/vibit scaffold feature
changes/_template/feature-request/
docs/agent-native-feature-request-scaffolding-gate.md
```

未来 implementation 可以添加：

- 用于 feature-request artifacts 的 `tools/vibit` scaffold command；
- `changes/_template/` 下的 feature-request template directory；
- 验证 template 存在并包含 required phase markers 的 repository checks；
- 说明 agent 如何在 coding 前使用 scaffold 的文档。

未来 implementation 不得生成 runtime code、protocol source、generated output、migrations、dependencies、hosted surfaces、SDKs、client libraries 或 direct compatibility shims。

## 4. 必需 Scaffold Artifacts

未来 scaffold 必须创建或引导创建这些文件：

```text
request.md
spec.yaml
impact.md
plan.md
checklist.md
verification.md
```

scaffolded `spec.yaml` 必须包含 placeholders 或 fields：

- user requirement；
- user-visible outcome；
- Nakama capability family；
- Pitaya status；
- acceptance criteria；
- test plan；
- implementation boundaries；
- generated output posture；
- migration posture；
- dependency posture；
- redaction posture；
- verification commands；
- durable memory updates；
- direct Nakama/Pitaya compatibility status。

## 5. 所有权

未来 implementation ownership：

```yaml
scaffold_command_owner: tools/vibit
template_owner: changes/_template/feature-request/
workflow_standard_owner: docs/agent-native-feature-request-test-workflow.md
gate_standard_owner: docs/agent-native-feature-request-scaffolding-gate.md
check_rule_owner: tools/vibit
change_artifact_owner: changes/<date>-<change-id>/
runtime_owner: not_applicable
protocol_owner: not_applicable
```

规则：

- `tools/vibit` 只有在未来 implementation work item 授权后，才能生成 text artifacts。
- Templates 必须保持 source-first 和 repository-local。
- Scaffolding 必须足够 deterministic，便于 agent review diffs。
- Scaffolding 不得绕过 change-spec workflow。
- 除非后续 work item 明确授权，scaffolding 不得创建 production secrets、`.env` 文件、generated runtime outputs、migrations 或 code stubs。

## 6. Redaction

未来 scaffold 不得记录或要求 agent 记录：

- raw device credential text or bytes；
- raw access tokens；
- verifier keys；
- credential 或 token lookup digests；
- credential 或 token verifier digests；
- HMAC inputs 或 outputs；
- 带凭据的 PostgreSQL DSNs；
- GitHub tokens 或 local secret values；
- headers、cookies、query strings、WebSocket subprotocol values、remote addresses 或 concrete transport metadata；
- 超出用户明确要求记录范围的 private user data。

允许的 scaffold 内容包括 placeholders、redaction reminders、route names、capability family names、module names、high-level status classes 和 explicit not-applicable rationale。

## 7. Nakama 映射

Nakama reference mapping：

- 本 gate 覆盖 `agent_native_requirement_test_implementation_workflow` capability family。
- 它通过要求每个未来 major feature 从 bounded request、acceptance criteria 和 test plan 开始，支撑 Nakama-first capability growth。
- 它不复制 Nakama public REST paths、client package names、runtime helper names、wire payloads、storage APIs、session token shapes 或 compatibility promises。

Pitaya reference mapping：

- Pitaya 保持为 future distributed architecture reference。
- 本 gate 不引入 frontend/backend server roles、RPC、service discovery、groups、cluster routing 或 distributed session behavior。

## 8. 未来 Implementation Work

打开：

```text
M-157/W-0229 Implement agent-native feature request scaffolding
```

未来 work item 可以：

- 添加 `changes/_template/feature-request/`；
- 添加窄口径 `tools/vibit scaffold feature` command；
- 为 scaffold output shape 添加 focused tests 或 repository checks；
- 如果 command usage 需要被引用，更新 `docs/agent-native-feature-request-test-workflow.md`；
- 更新 repository checks 和 durable memory。

未来 work item 不得：

- 添加 runtime behavior；
- 添加 protocol routes；
- 添加 Protobuf source 或 generated output；
- 添加 migrations、persistence、repository interfaces、adapters、dependencies、startup wiring、SDK publication、generated client libraries、hosted demos、release artifacts、direct compatibility 或 Pitaya-style distributed architecture。

## 9. Verification Expectations

本 gate 应通过以下命令验证：

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.agent_native_feature_request_scaffolding_gate
node tools/vibit check change define-agent-native-feature-request-scaffolding-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

本 gate 不改变 Go source 或 runtime behavior，因此不要求 Go tests。

## 10. Stop Conditions

如果 scaffolding 需要以下内容，应停止并创建单独 gate：

- runtime code generation；
- protocol source generation；
- generated Go output；
- migrations；
- dependency adoption；
- startup wiring；
- public SDK 或 client package behavior；
- hosted demo behavior；
- release artifact behavior；
- 超出 explicit change text 的 user private data ingestion；
- direct Nakama/Pitaya API compatibility；
- Pitaya-style distributed runtime、frontend/backend split、RPC、service discovery、groups 或 cluster routing。

## 11. Future Implementation Acceptance Criteria

只有满足以下条件，未来 implementation 才应被接受：

- scaffold 创建或说明六个 required change artifacts；
- scaffold 包含 Nakama capability mapping 和 Pitaya deferral fields；
- scaffold 要求 implementation 前先写 acceptance criteria 和 test plan；
- scaffold 记录 implementation boundaries 和 forbidden scope；
- scaffold 记录 verification commands 和 not-applicable rationale；
- scaffold 包含 redaction reminders；
- repository checks 验证 scaffolded shape 或 template shape；
- 没有添加 runtime、protocol、generated output、migration、dependency、hosted、SDK、distributed runtime 或 direct compatibility scope。
