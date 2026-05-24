# Agent-Native Feature Request Scaffolding

状态：Accepted v0.1
最后更新：2026-05-24
范围：把用户 backend requirement 转成 bounded change artifacts 的 source-first scaffold
依赖：`docs/agent-native-feature-request-scaffolding-gate.md`、`docs/agent-native-feature-request-test-workflow.md`、`docs/change-spec.md`、`docs/workflow.md`
Canonical decision：`ADR-0137`

英文文件 `docs/agent-native-feature-request-scaffolding.md` 是权威版本。本文是简体中文译本。

本文记录已实现的 feature request scaffold。它添加 templates 和 `tools/vibit scaffold feature` 命令。Pitaya remains deferred，作为未来 architecture reference；Nakama 保持主要 product capability reference。它不添加 runtime behavior、protocol routes、Protobuf source、generated output、migrations、dependencies、persistence、startup wiring、authentication/session behavior changes、delivery guarantees、stream subscriptions、chat rooms、groups、broadcast fanout、matchmaking、match runtime、operations/admin behavior、SDK publication、generated client libraries、hosted deployments、release artifacts、Pitaya-style distributed architecture 或 direct Nakama/Pitaya API compatibility。

## 1. 核心规则

Scaffold implementation 记录为：

```yaml
agent_native_feature_request_scaffolding: implemented
completed_work_item: W-0229
decision: ADR-0137
check_rule: runtime.agent_native_feature_request_scaffolding_implementation
source_gate_decision: ADR-0136
source_selection_decision: ADR-0135
source_workflow_decision: ADR-0128
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
scaffold_command: tools/vibit scaffold feature
template_directory: changes/_template/feature-request/
required_scaffold_artifacts:
  - request.md
  - spec.yaml
  - impact.md
  - plan.md
  - checklist.md
  - verification.md
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
sdk_added: false
hosted_deployment_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. 命令

使用方式：

```bash
node tools/vibit scaffold feature <change-id> --request "<original request>" --summary "<one-line summary>"
```

可选参数：

```bash
node tools/vibit scaffold feature <change-id> --date YYYY-MM-DD
node tools/vibit scaffold feature <change-id> --dry-run
```

该命令创建：

```text
changes/YYYY-MM-DD-<change-id>/
  request.md
  spec.yaml
  impact.md
  plan.md
  checklist.md
  verification.md
```

如果目标 change directory 已存在，命令会拒绝覆盖。`--dry-run` 会校验目标路径和 template replacement，但不写入文件。

## 3. Template 要求

Feature-request template 必须包含：

- original request；
- clarified requirement；
- user-visible outcome；
- Nakama capability mapping；
- Pitaya deferred status；
- non-goals 和 unknowns；
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

## 4. Agent 规则

Agents 在实现非平凡、用户可见的 backend feature work 前，应先使用 scaffold。

规则：

- 先填写 scaffold，再写代码。
- Nakama 保持主要 product capability reference。
- Pitaya 继续 deferred，除非后续 ADR 明确重新激活。
- 在 implementation 前记录 tests 或具体 not-applicable rationale。
- 不要把 scaffold creation 当作添加 runtime、protocol、generated output、migration、dependency、SDK、hosted、distributed runtime 或 direct compatibility scope 的授权。
- 不要记录 raw credentials、tokens、verifier keys、digests、带凭证的 DSNs、GitHub tokens、transport metadata、headers、cookies、query strings、WebSocket subprotocols、remote addresses，或超出明确 request text 的 private user data。

## 5. 验证

Implementation 使用以下命令验证：

```bash
node -c tools/vibit
node tools/vibit scaffold feature scaffold-smoke --date 2026-05-24 --request "Smoke test feature request scaffold." --summary "Smoke test feature request scaffold." --dry-run
node tools/vibit inspect rule runtime.agent_native_feature_request_scaffolding_implementation
node tools/vibit check change implement-agent-native-feature-request-scaffolding --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```
