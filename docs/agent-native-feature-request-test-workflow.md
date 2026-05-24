# Agent-Native Feature Request And Test Workflow

Status: Accepted v0.1
Last updated: 2026-05-24
Scope: Default workflow for turning user backend requirements into bounded specifications, tests, implementation, verification, and durable project memory
Depends on: `CONSTITUTION.md`, `docs/change-spec.md`, `docs/workflow.md`, `docs/nakama-pitaya-product-parity-roadmap.md`, `docs/reference-game-server-alignment.md`, `decisions/ADR-0127-nakama-first-ai-native-requirement-test-workflow-direction.md`
Canonical decision: `ADR-0128`

The paired Simplified Chinese translation is `docs/agent-native-feature-request-test-workflow.zh-CN.md`. The English file is authoritative.

This document defines the agent-native feature request and test workflow. It does not add runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, startup wiring, broad product modules, Pitaya-style distributed architecture, hosted deployments, release artifacts, public announcements, paid promotion, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The agent-native feature request and test workflow record is:

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

## 2. Purpose

The maintainer direction is now explicit: vibit should target a Nakama-class product surface, but its differentiator is AI-native development and AI-native testing. A user should be able to state a backend requirement, and AI should carry the work through specification, acceptance criteria, test planning, tests, implementation, verification, and durable repository memory.

That outcome requires a workflow standard before broad new product modules expand. Without a standard, future agents may add locally plausible code while skipping the artifacts that make the code maintainable by later agents.

This workflow makes the product thesis operational:

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

## 3. Applicability

Use this workflow for non-trivial user-facing backend feature work, including:

- a new capability family;
- a new command, query, event, permission, error, route, protocol payload, migration, repository behavior, runtime service, or operational surface;
- a behavioral change that users, clients, operators, or contributors can observe;
- a Nakama-aligned product feature slice;
- any feature where acceptance tests, permission behavior, persistence behavior, protocol behavior, or failure behavior matters.

Small typo fixes, formatting-only edits, and purely mechanical documentation corrections may use a lighter change record. The agent must still record what changed and how it was verified.

## 4. Workflow Phases

### 4.1 User Requirement

Capture the user's original request in `request.md`. Preserve the user's language where it clarifies intent, but convert implementation assumptions into explicit unknowns or non-goals.

Minimum output:

```yaml
user_requirement: <plain requirement>
user_visible_outcome: <what changes for a developer, player, operator, or agent>
non_goals:
  - <explicitly out of scope>
unknowns:
  - <question or assumption>
```

### 4.2 Requirement Spec

Record the bounded feature spec in `spec.yaml`. The spec must be narrow enough that another agent can implement or review it without reopening the entire product direction.

Required expectations:

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

Every major feature must map to a Nakama-style capability family or explicitly say that no mapping applies.

Accepted mapping families are the roadmap families in `docs/nakama-pitaya-product-parity-roadmap.md`, including identity/auth/session, storage, presence/status/notifications, chat/realtime messaging, friends/groups/parties, leaderboards/tournaments, economy/progression, matchmaking, match runtime, operations, SDK/developer experience, and agent-native requirement/test workflow.

Rules:

- Nakama is the primary product capability reference.
- The mapping explains product intent, not API compatibility.
- Do not copy Nakama public routes, payloads, storage models, runtime API names, or compatibility shims without a future explicit ADR.
- Pitaya remains deferred as a future distributed architecture reference.

### 4.4 Acceptance Criteria

Acceptance criteria must be user-judgable and testable. They should describe observable outcomes, not internal code preferences.

For feature behavior, include positive and negative criteria. For standards or planning changes, include artifact and checkability criteria.

### 4.5 Test Plan

Non-trivial behavior needs tests before or with implementation. The test plan should be recorded before code changes begin, then updated if implementation reveals a better test boundary.

Use the relevant classes:

- positive behavior;
- negative behavior;
- permission and authentication behavior;
- persistence and transaction behavior;
- protocol encoding, decoding, and route mapping;
- failure paths and redaction;
- concurrency, idempotency, or ordering behavior where relevant;
- integration or end-to-end proof where the feature crosses boundaries;
- repository checks for architecture, generated output, manifests, and docs.

If tests are not applicable, `verification.md` must record the explicit rationale. "No tests" is not enough.

### 4.6 Tests

Tests should live at the lowest boundary that proves the behavior:

- domain/application service tests for business invariants;
- repository adapter tests for persistence and SQL behavior;
- protocol bridge tests for payload mapping;
- handler or E2E tests for cross-boundary request flow;
- repository checks for architecture rules and artifact presence.

Tests must not rely on live services by default unless the repository already marks them opt-in and documents prerequisites.

### 4.7 Implementation Boundaries

Before implementation, record:

- owning module or runtime subsystem;
- allowed file/package areas;
- forbidden file/package areas;
- generated outputs and their source of truth;
- migrations and persistence ownership;
- dependency adoption status;
- public contract or protocol gate status;
- redaction and log-safety requirements;
- direct Nakama/Pitaya compatibility status.

Implementation must stay inside the declared boundaries. If the feature needs a broader boundary, stop and create or update an ADR before coding through it.

### 4.8 Verification

`verification.md` must list commands actually run and their result.

Rules:

- Do not claim verification that did not run.
- Record unavailable verification with a reason and a follow-up path.
- Include `node tools/vibit check all --json` or a narrower justified set for architecture-affecting changes.
- Include Go tests for runtime behavior changes.
- Include generated-output checks when Protobuf or generated sources are involved.
- Include migration checks when SQL sources are involved.

### 4.9 Durable Memory

Update durable project memory when the change alters intent, architecture, product direction, or continuation state.

Memory may include:

- change spec files;
- ADRs;
- conversation logs;
- `.arch/` manifests;
- module manifests;
- module or repository `AGENTS.md` guides;
- README and public docs;
- rule catalog entries;
- check logic in `tools/vibit`.

Do not leave important product direction only in chat history.

## 5. Required Artifacts

For non-trivial user-facing feature work, create or update:

- `request.md`: original request, clarified requirement, non-goals, unknowns, acceptance criteria.
- `spec.yaml`: machine-readable scope, Nakama mapping, tests, boundaries, verification, memory updates.
- `impact.md`: module ownership, public behavior, contracts, data, protocol, tests, docs, compatibility.
- `plan.md`: files to create/edit, generated outputs, handwritten logic, tests, verification commands, rollback notes.
- `checklist.md`: tracked completion tasks.
- `verification.md`: commands actually run, skipped checks, not-applicable rationale.

Also create or update:

- an ADR when direction, architecture, standards, public contracts, generated conventions, dependencies, or major boundaries change;
- a conversation log when maintainer intent or product direction matters;
- relevant manifests and `AGENTS.md` files when continuation, ownership, or agent behavior changes;
- tests or explicit not-applicable rationale.

## 6. Test Policy

The default is test-first or test-with-implementation.

Rules:

- A non-trivial behavior change must have tests before or in the same change.
- A public protocol or route change must have protocol mapping tests and at least one handler or flow proof unless an ADR defers it.
- A persistence change must have migration checks and repository/adapter tests.
- A permission/authentication change must test allowed, denied, missing proof, invalid proof, and redacted failure behavior where applicable.
- A concurrency-sensitive change must test stale state, duplicate action, ordering, or race-sensitive behavior at the narrowest feasible boundary.
- A docs-only standard change may use repository checks and explicit not-applicable test rationale.

## 7. Nakama-First Posture

Nakama is the active primary product reference. It guides which backend capability families matter and how users will recognize product usefulness.

Pitaya is deferred. It may inform future distributed architecture vocabulary only after a later ADR reactivates it.

The current posture is:

```yaml
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
api_compatibility_goal: false
direct_nakama_pitaya_api_compatibility_added: false
```

Feature work should say "Nakama-style capability family" rather than "Nakama-compatible API" unless a future compatibility ADR changes the goal.

## 8. Stop Conditions

Stop and ask for maintainer direction before:

- direct Nakama or Pitaya API compatibility;
- Pitaya-style cluster/RPC/frontend-backend/service-discovery work;
- new runtime behavior outside the current work item;
- new protocol route or Protobuf source without a protocol work item;
- generated output without source schema and generation checks;
- migrations without a persistence/schema work item;
- new dependencies without dependency adoption;
- broad product modules such as chat, groups, matchmaking, match runtime, leaderboards, economy, or operations without a bounded work item;
- hosted deployments, release artifacts, public announcements, or paid promotion.

## 9. Next Pilot

The next work item is:

```text
M-149/W-0221 Pilot Nakama-aligned feature request workflow
```

The pilot should apply this workflow to choose and shape the next bounded Nakama-aligned prototype-ready feature slice. The pilot should prove that the workflow can turn a product requirement into spec, acceptance criteria, test plan, implementation boundary, verification plan, and durable memory without jumping directly to broad runtime scope.

