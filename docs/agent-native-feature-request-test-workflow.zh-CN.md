# Agent-Native Feature Request And Test Workflow

状态：Accepted v0.1
最后更新：2026-05-24
范围：把用户后端需求转成 bounded specification、tests、implementation、verification 和 durable project memory 的默认流程
依赖：`CONSTITUTION.md`、`docs/change-spec.md`、`docs/workflow.md`、`docs/nakama-pitaya-product-parity-roadmap.md`、`docs/reference-game-server-alignment.md`、`decisions/ADR-0127-nakama-first-ai-native-requirement-test-workflow-direction.md`
权威决策：`ADR-0128`

说明：本文件是 `docs/agent-native-feature-request-test-workflow.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本文定义 agent-native feature request and test workflow。它不添加 runtime behavior、protocol routes、Protobuf source、generated output、migrations、dependencies、startup wiring、broad product modules、Pitaya-style distributed architecture、hosted deployments、release artifacts、public announcements、paid promotion 或 direct Nakama/Pitaya API compatibility。

## 1. 核心规则

Agent-native feature request and test workflow 记录为：

```yaml
agent_native_feature_request_test_workflow: defined
completed_work_item: W-0220
decision: ADR-0128
check_rule: runtime.agent_native_feature_request_test_workflow
source_direction_decision: ADR-0127
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
implementation_authorized_by_this_standard: false
future_pilot_work_item: W-0221
future_pilot_direction: pilot_nakama_aligned_feature_request_workflow
workflow_phases:
  - user_requirement
  - requirement_spec
  - nakama_capability_mapping
  - acceptance_criteria
  - test_plan
  - tests
  - implementation_boundaries
  - verification
  - durable_memory
required_artifacts:
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
distributed_runtime_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. 目的

维护者方向已经明确：vibit 应以 Nakama-class product surface 为目标，但差异化能力是 AI-native development 和 AI-native testing。用户应该可以说出一个后端需求，AI 应继续完成 specification、acceptance criteria、test planning、tests、implementation、verification 和 durable repository memory。

在扩展 broad new product modules 之前，需要先定义 workflow standard。否则未来 agent 可能写出局部看似合理的代码，却跳过让后续 agent 可维护该代码的 artifacts。

这个流程把产品命题落到执行层：

```text
user requirement
-> AI-written bounded requirement spec
-> Nakama capability mapping
-> AI-written acceptance criteria
-> AI-written test plan
-> AI-written or updated tests
-> implementation inside declared boundaries
-> AI-run verification
-> AI-updated docs, manifests, ADRs, change records, and conversation memory
```

## 3. 适用范围

以下非平凡 user-facing backend feature work 必须使用本流程：

- 新 capability family；
- 新 command、query、event、permission、error、route、protocol payload、migration、repository behavior、runtime service 或 operational surface；
- 用户、client、operator 或 contributor 可观察到的行为变化；
- Nakama-aligned product feature slice；
- acceptance tests、permission behavior、persistence behavior、protocol behavior 或 failure behavior 重要的 feature。

小 typo、formatting-only edit 和纯机械文档修正可以用更轻量的 change record。agent 仍必须记录改了什么以及如何验证。

## 4. 流程阶段

### 4.1 User Requirement

在 `request.md` 中捕获用户原始请求。保留能说明意图的用户语言，但要把实现假设转成明确 unknowns 或 non-goals。

最小输出：

```yaml
user_requirement: <plain requirement>
user_visible_outcome: <what changes for a developer, player, operator, or agent>
non_goals:
  - <explicitly out of scope>
unknowns:
  - <question or assumption>
```

### 4.2 Requirement Spec

在 `spec.yaml` 中记录 bounded feature spec。spec 必须足够窄，使另一个 agent 能实现或 review，而不需要重新打开整个产品方向。

必需期待：

```yaml
user_requirement: <plain requirement>
nakama_capability_family: <roadmap family or no_mapping_applies>
acceptance_criteria:
  - <observable condition>
test_plan:
  - <test class or command>
implementation_boundaries:
  allowed:
    - <allowed file, package, module, or artifact family>
  forbidden:
    - <forbidden scope>
verification:
  required:
    - <command>
memory_updates:
  - <docs, manifests, ADRs, conversations, or module guides to update>
```

### 4.3 Nakama Capability Mapping

每个 major feature 必须映射到一个 Nakama-style capability family，或明确说明 no mapping applies。

接受的 mapping families 是 `docs/nakama-pitaya-product-parity-roadmap.md` 中的 roadmap families，包括 identity/auth/session、storage、presence/status/notifications、chat/realtime messaging、friends/groups/parties、leaderboards/tournaments、economy/progression、matchmaking、match runtime、operations、SDK/developer experience 和 agent-native requirement/test workflow。

规则：

- Nakama 是 primary product capability reference。
- mapping 解释 product intent，不表示 API compatibility。
- 未经未来明确 ADR，不要复制 Nakama public routes、payloads、storage models、runtime API names 或 compatibility shims。
- Pitaya 暂缓为 future distributed architecture reference。

### 4.4 Acceptance Criteria

Acceptance criteria 必须是用户可判断、可测试的。它们应描述 observable outcomes，而不是内部代码偏好。

对 feature behavior，要包含 positive 和 negative criteria。对 standards 或 planning changes，要包含 artifact 和 checkability criteria。

### 4.5 Test Plan

非平凡行为需要在实现之前或实现同时写测试。test plan 应在代码修改前记录；如果实现过程中发现更好的测试边界，再更新它。

使用相关类别：

- positive behavior；
- negative behavior；
- permission 和 authentication behavior；
- persistence 和 transaction behavior；
- protocol encoding、decoding 和 route mapping；
- failure paths 和 redaction；
- 相关时的 concurrency、idempotency 或 ordering behavior；
- feature 跨边界时的 integration 或 end-to-end proof；
- architecture、generated output、manifests 和 docs 的 repository checks。

如果 tests 不适用，`verification.md` 必须记录明确原因。只写 "No tests" 不够。

### 4.6 Tests

测试应放在能证明行为的最低边界：

- business invariants 用 domain/application service tests；
- persistence 和 SQL behavior 用 repository adapter tests；
- payload mapping 用 protocol bridge tests；
- cross-boundary request flow 用 handler 或 E2E tests；
- architecture rules 和 artifact presence 用 repository checks。

默认测试不应依赖 live services，除非仓库已把它们标记为 opt-in 并记录 prerequisites。

### 4.7 Implementation Boundaries

实现前要记录：

- owning module 或 runtime subsystem；
- allowed file/package areas；
- forbidden file/package areas；
- generated outputs 及其 source of truth；
- migrations 和 persistence ownership；
- dependency adoption status；
- public contract 或 protocol gate status；
- redaction 和 log-safety requirements；
- direct Nakama/Pitaya compatibility status。

实现必须留在声明的边界内。如果 feature 需要更宽边界，停止并创建或更新 ADR，而不是直接穿透边界写代码。

### 4.8 Verification

`verification.md` 必须列出实际运行的命令及结果。

规则：

- 不要声称运行了未运行的验证。
- unavailable verification 要记录原因和 follow-up path。
- architecture-affecting changes 应包含 `node tools/vibit check all --json` 或更窄但有理由的命令集合。
- runtime behavior changes 应包含 Go tests。
- 涉及 Protobuf 或 generated sources 时包含 generated-output checks。
- 涉及 SQL sources 时包含 migration checks。

### 4.9 Durable Memory

当变更影响 intent、architecture、product direction 或 continuation state 时，要更新 durable project memory。

Memory 可以包括：

- change spec files；
- ADRs；
- conversation logs；
- `.arch/` manifests；
- module manifests；
- module 或 repository `AGENTS.md` guides；
- README 和 public docs；
- rule catalog entries；
- `tools/vibit` check logic。

不要把重要产品方向只留在 chat history 里。

## 5. 必需 Artifacts

对非平凡 user-facing feature work，要创建或更新：

- `request.md`：original request、clarified requirement、non-goals、unknowns、acceptance criteria。
- `spec.yaml`：machine-readable scope、Nakama mapping、tests、boundaries、verification、memory updates。
- `impact.md`：module ownership、public behavior、contracts、data、protocol、tests、docs、compatibility。
- `plan.md`：files to create/edit、generated outputs、handwritten logic、tests、verification commands、rollback notes。
- `checklist.md`：tracked completion tasks。
- `verification.md`：实际运行的命令、skipped checks、not-applicable rationale。

同时在需要时创建或更新：

- 当 direction、architecture、standards、public contracts、generated conventions、dependencies 或 major boundaries 变化时创建 ADR；
- 当 maintainer intent 或 product direction 重要时创建 conversation log；
- 当 continuation、ownership 或 agent behavior 变化时更新相关 manifests 和 `AGENTS.md`；
- tests 或明确的 not-applicable rationale。

## 6. 测试政策

默认是 test-first 或 test-with-implementation。

规则：

- 非平凡 behavior change 必须在同一 change 中拥有 tests，最好先写 tests。
- public protocol 或 route change 必须有 protocol mapping tests，并且除非 ADR 明确延期，应至少有一个 handler 或 flow proof。
- persistence change 必须有 migration checks 和 repository/adapter tests。
- permission/authentication change 应在适用时测试 allowed、denied、missing proof、invalid proof 和 redacted failure behavior。
- concurrency-sensitive change 必须在可行的最低边界测试 stale state、duplicate action、ordering 或 race-sensitive behavior。
- docs-only standard change 可以用 repository checks 和明确的 not-applicable test rationale。

## 7. Nakama-First 姿态

Nakama 是 active primary product reference。它指导哪些 backend capability families 重要，以及用户如何识别产品有用性。

Pitaya 已暂缓。它只能在未来 ADR 重新激活后，用于 distributed architecture vocabulary。

当前姿态是：

```yaml
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
api_compatibility_goal: false
direct_nakama_pitaya_api_compatibility_added: false
```

Feature work 应说 "Nakama-style capability family"，不要说 "Nakama-compatible API"，除非未来 compatibility ADR 改变目标。

## 8. 停止条件

遇到以下情况要停止并请求维护者方向：

- direct Nakama 或 Pitaya API compatibility；
- Pitaya-style cluster/RPC/frontend-backend/service-discovery work；
- 当前 work item 之外的新 runtime behavior；
- 没有 protocol work item 的新 protocol route 或 Protobuf source；
- 没有 source schema 和 generation checks 的 generated output；
- 没有 persistence/schema work item 的 migrations；
- 没有 dependency adoption 的新 dependencies；
- 没有 bounded work item 的 broad product modules，例如 chat、groups、matchmaking、match runtime、leaderboards、economy 或 operations；
- hosted deployments、release artifacts、public announcements 或 paid promotion。

## 9. 下一项 Pilot

下一项 work item 是：

```text
M-149/W-0221 Pilot Nakama-aligned feature request workflow
```

该 pilot 应使用本流程选择并塑造下一个 bounded Nakama-aligned prototype-ready feature slice。pilot 应证明这个 workflow 能把 product requirement 转成 spec、acceptance criteria、test plan、implementation boundary、verification plan 和 durable memory，而不是直接跳到 broad runtime scope。

