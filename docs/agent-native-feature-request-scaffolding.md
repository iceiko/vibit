# Agent-Native Feature Request Scaffolding

Status: Accepted v0.1
Last updated: 2026-05-24
Scope: Source-first scaffold for turning user backend requirements into bounded change artifacts
Depends on: `docs/agent-native-feature-request-scaffolding-gate.md`, `docs/agent-native-feature-request-test-workflow.md`, `docs/change-spec.md`, `docs/workflow.md`
Canonical decision: `ADR-0137`

The paired Simplified Chinese translation is `docs/agent-native-feature-request-scaffolding.zh-CN.md`. The English file is authoritative.

This document records the implemented feature request scaffold. It adds templates and a `tools/vibit scaffold feature` command. Pitaya remains deferred as a future architecture reference while Nakama remains the primary product capability reference. It does not add runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, authentication/session behavior changes, delivery guarantees, stream subscriptions, chat rooms, groups, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The scaffold implementation record is:

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

## 2. Command

Use:

```bash
node tools/vibit scaffold feature <change-id> --request "<original request>" --summary "<one-line summary>"
```

Optional:

```bash
node tools/vibit scaffold feature <change-id> --date YYYY-MM-DD
node tools/vibit scaffold feature <change-id> --dry-run
```

The command creates:

```text
changes/YYYY-MM-DD-<change-id>/
  request.md
  spec.yaml
  impact.md
  plan.md
  checklist.md
  verification.md
```

It refuses to overwrite an existing change directory. `--dry-run` validates the target path and template replacement without writing files.

## 3. Template Requirements

The feature-request template must include:

- original request;
- clarified requirement;
- user-visible outcome;
- Nakama capability mapping;
- Pitaya deferred status;
- non-goals and unknowns;
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

## 4. Agent Rules

Agents should use the scaffold before implementing non-trivial user-facing backend feature work.

Rules:

- Fill the scaffold before coding.
- Keep Nakama as the primary product capability reference.
- Keep Pitaya deferred unless a later ADR explicitly reactivates it.
- Record tests or a concrete not-applicable rationale before implementation.
- Do not treat scaffold creation as permission to add runtime, protocol, generated output, migration, dependency, SDK, hosted, distributed runtime, or direct compatibility scope.
- Do not record raw credentials, tokens, verifier keys, digests, DSNs with credentials, GitHub tokens, transport metadata, headers, cookies, query strings, WebSocket subprotocols, remote addresses, or private user data beyond explicit request text.

## 5. Verification

The implementation is verified by:

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
