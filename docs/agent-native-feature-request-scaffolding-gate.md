# Agent-Native Feature Request Scaffolding Gate

Status: Accepted v0.1
Last updated: 2026-05-24
Scope: Gate-only boundary for future feature request scaffolding that turns user backend requirements into bounded change artifacts before implementation
Depends on: `docs/agent-native-feature-request-test-workflow.md`, `docs/change-spec.md`, `docs/workflow.md`, `docs/nakama-pitaya-product-parity-roadmap.md`, `decisions/ADR-0135-agent-native-feature-request-scaffolding-selection.md`
Canonical decision: `ADR-0136`

The paired Simplified Chinese translation is `docs/agent-native-feature-request-scaffolding-gate.zh-CN.md`. The English file is authoritative.

This document defines the agent-native feature request scaffolding gate. It is a gate artifact. It does not implement scaffolding, add templates, add `tools/vibit` scaffold commands, change runtime behavior, add protocol routes, add Protobuf source, change generated output, add migrations, add dependencies, add persistence, change startup wiring, change authentication/session behavior, add delivery guarantees, add stream subscriptions, add chat rooms, add groups, add broadcast fanout, add matchmaking, add match runtime, add operations/admin behavior, publish SDKs, generate client libraries, add hosted deployments, create release artifacts, run public announcements, run paid promotion, add Pitaya-style distributed architecture, or add direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The agent-native feature request scaffolding gate record is:

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

## 2. Purpose

`ADR-0135` selected `agent_native_requirement_test_implementation_workflow` as the next Nakama-first prototype-ready capability family after the local alpha example path. The existing workflow standard defines what a good feature change must contain, but a future user-facing product loop needs a more executable intake path.

The future scaffolding should let an agent start from a user requirement and create the correct source-first change artifact shape before implementation begins. The scaffold must reinforce the product purpose:

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

The gate keeps this work deliberately narrow. It defines the posture and stop conditions before any scaffold command, template directory, or check implementation exists.

## 3. Selected Scaffolding Shape

The first accepted shape is source-first repository-local change artifact scaffolding.

Future implementation candidates:

```text
tools/vibit scaffold feature
changes/_template/feature-request/
docs/agent-native-feature-request-scaffolding-gate.md
```

The future implementation may add:

- a `tools/vibit` scaffold command for feature-request artifacts;
- a feature-request template directory under `changes/_template/`;
- repository checks that verify the template exists and contains the required phase markers;
- docs that explain how agents should use the scaffold before coding.

The future implementation must not generate runtime code, protocol source, generated output, migrations, dependencies, hosted surfaces, SDKs, client libraries, or direct compatibility shims.

## 4. Required Scaffold Artifacts

The future scaffold must create or guide creation of these files:

```text
request.md
spec.yaml
impact.md
plan.md
checklist.md
verification.md
```

The scaffolded `spec.yaml` must include placeholders or fields for:

- user requirement;
- user-visible outcome;
- Nakama capability family;
- Pitaya status;
- acceptance criteria;
- test plan;
- implementation boundaries;
- generated output posture;
- migration posture;
- dependency posture;
- redaction posture;
- verification commands;
- durable memory updates;
- direct Nakama/Pitaya compatibility status.

## 5. Ownership

Future implementation ownership:

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

Rules:

- `tools/vibit` may generate text artifacts only after a future implementation work item authorizes it.
- Templates must remain source-first and repository-local.
- Scaffolding must be deterministic enough for agents to review diffs.
- Scaffolding must not bypass the change-spec workflow.
- Scaffolding must not create production secrets, `.env` files, generated runtime outputs, migrations, or code stubs unless a later work item explicitly authorizes that scope.

## 6. Redaction

The future scaffold must not record or ask agents to record:

- raw device credential text or bytes;
- raw access tokens;
- verifier keys;
- credential or token lookup digests;
- credential or token verifier digests;
- HMAC inputs or outputs;
- PostgreSQL DSNs with credentials;
- GitHub tokens or local secret values;
- headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or concrete transport metadata;
- private user data beyond what the user explicitly asks to record in the change request.

Allowed scaffold content includes placeholders, redaction reminders, route names, capability family names, module names, high-level status classes, and explicit not-applicable rationale.

## 7. Nakama Mapping

Nakama reference mapping:

- The gate covers the `agent_native_requirement_test_implementation_workflow` capability family.
- It supports Nakama-first capability growth by requiring every future major feature to start from a bounded request, acceptance criteria, and test plan.
- It does not copy Nakama public REST paths, client package names, runtime helper names, wire payloads, storage APIs, session token shapes, or compatibility promises.

Pitaya reference mapping:

- Pitaya remains deferred as a future distributed architecture reference.
- This gate does not introduce frontend/backend server roles, RPC, service discovery, groups, cluster routing, or distributed session behavior.

## 8. Future Implementation Work

Open:

```text
M-157/W-0229 Implement agent-native feature request scaffolding
```

The future work item may:

- add `changes/_template/feature-request/`;
- add a narrow `tools/vibit scaffold feature` command;
- add focused tests or repository checks for scaffold output shape;
- update `docs/agent-native-feature-request-test-workflow.md` if command usage needs to be referenced;
- update repository checks and durable memory.

The future work item must not:

- add runtime behavior;
- add protocol routes;
- add Protobuf source or generated output;
- add migrations, persistence, repository interfaces, adapters, dependencies, startup wiring, SDK publication, generated client libraries, hosted demos, release artifacts, direct compatibility, or Pitaya-style distributed architecture.

## 9. Verification Expectations

This gate should be verified with:

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

Go tests are not required for this gate because it does not change Go source or runtime behavior.

## 10. Stop Conditions

Stop and create a separate gate if scaffolding requires:

- runtime code generation;
- protocol source generation;
- generated Go output;
- migrations;
- dependency adoption;
- startup wiring;
- public SDK or client package behavior;
- hosted demo behavior;
- release artifact behavior;
- user private data ingestion beyond explicit change text;
- direct Nakama/Pitaya API compatibility;
- Pitaya-style distributed runtime, frontend/backend split, RPC, service discovery, groups, or cluster routing.

## 11. Future Implementation Acceptance Criteria

The future implementation should be accepted only if:

- the scaffold creates or documents the six required change artifacts;
- the scaffold includes Nakama capability mapping and Pitaya deferral fields;
- the scaffold requires acceptance criteria and test planning before implementation;
- the scaffold records implementation boundaries and forbidden scope;
- the scaffold records verification commands and not-applicable rationale;
- the scaffold includes redaction reminders;
- repository checks verify the scaffolded shape or template shape;
- no runtime, protocol, generated output, migration, dependency, hosted, SDK, distributed runtime, or direct compatibility scope is added.
