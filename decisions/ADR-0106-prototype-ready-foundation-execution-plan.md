# ADR-0106: Prototype-Ready Foundation Execution Plan

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-define-prototype-ready-foundation-execution-plan/`

Related conversations:

- `conversations/2026-05-22-prototype-ready-foundation-execution-plan.md`

Related artifacts:

- `docs/prototype-ready-foundation-execution-plan.md`
- `docs/prototype-ready-foundation-execution-plan.zh-CN.md`
- `docs/product-maturity-milestones.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`v0.1.0-alpha.1` is now a source-first alpha with a real authenticated gameplay loop and a repository-owned feedback intake surface. `ADR-0105` records the product maturity stages: source-first alpha, prototype-ready foundation, single-node production-candidate foundation, and Nakama/Pitaya-class product.

The maintainer asked to keep pushing toward a real product and productivity stage. The next necessary step is to turn the Stage 2 ambition into an ordered execution plan without silently changing runtime scope or overstating production readiness.

## Decision

Define `docs/prototype-ready-foundation-execution-plan.md` as the first execution plan for moving from source-first alpha to prototype-ready foundation.

Select the next first execution slice as:

```text
W-0199 Define prototype-ready local development path gate
```

The first slice is a gate because vibit should make setup, migration, configuration, and example flow repeatable before adding broader shared online-service behavior. Storage objects, realtime messaging/server push, failure/concurrency verification, and minimal operations inspection remain planned later families.

The repository check rule for this decision is:

```text
runtime.prototype_ready_foundation_execution_plan
```

Do not implement runtime behavior, add protocol routes, add Protobuf source or generated output, add migrations, add dependencies, broaden operations/admin behavior, add hosted deployments, create release artifacts, run public announcements, run paid promotion, change authentication/session behavior, add broad product modules, or add direct Nakama/Pitaya API compatibility from this planning slice.

## Decision Record

```yaml
prototype_ready_foundation_execution_plan: defined
completed_work_item: W-0198
decision: ADR-0106
check_rule: runtime.prototype_ready_foundation_execution_plan
source_stage: source_first_alpha
source_release_identifier: v0.1.0-alpha.1
target_stage: prototype_ready_foundation
execution_plan_standard: docs/prototype-ready-foundation-execution-plan.md
recommended_sequence_recorded: true
candidate_work_items_recorded: true
maturity_stage_mapping_recorded: true
nakama_pitaya_capability_mapping_recorded: true
success_criteria_recorded: true
stop_conditions_recorded: true
selected_first_execution_slice: prototype_ready_local_development_path_gate
next_work_item: W-0199
next_direction: prototype_ready_local_development_path_gate
planning_only: true
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
direct_nakama_pitaya_api_compatibility_added: false
```

## Alternatives Considered

- Jump directly to storage objects.
- Jump directly to chat, streams, or push.
- Run broader public outreach before reducing local setup friction.
- Start packaging binaries or containers before the source-first local path is stronger.
- Treat the current alpha as production-adjacent.

## Rationale

The most product-useful next step is not a broad feature leap. It is making vibit easier to run and evaluate as a prototype base. A reproducible local development path and a richer example flow will make later storage, push, messaging, operations, and verification work easier to test and easier for external developers to understand.

Nakama and Pitaya remain capability baselines, but the first prototype-ready step should strengthen vibit's own source-first developer experience and agent-native workflow before broadening product modules.

## Agent Reasoning Summary

The agent treated `W-0198` as a planning and sequencing slice. The plan keeps the current alpha honest, records the Stage 2 path, chooses a low-risk first gate around local development and examples, and preserves all runtime, protocol, migration, dependency, release, hosted deployment, outreach, authentication/session, broad product, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  developer_activation: high
  prototype_usefulness: high
  sequencing_clarity: high
  runtime_scope_change: none
  artifact_scope_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `W-0198` completes as the durable prototype-ready execution plan.
- `W-0199` becomes the next ready work item.
- The next work should define a local development path gate before behavior changes.
- Storage objects, realtime messaging/server push, failure/concurrency verification, and operations inspection remain planned but not yet authorized implementation.
- No runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, hosted deployment, operations/admin expansion, authentication/session behavior change, broad product module, extra release artifact, public announcement, paid promotion, or direct Nakama/Pitaya API compatibility is added by this decision.

## Reversal Conditions

Revisit this decision if early feedback shows that a different first prototype-ready blocker is more urgent than local development and example ergonomics, or if the maintainer explicitly authorizes a different first Stage 2 execution slice.

## Follow-Up

- Define the prototype-ready local development path gate.
- Keep later implementation work bounded and explicit.
- Convert real alpha feedback into work items when available.
- Ask for explicit maintainer authorization before runtime, protocol, migration, dependency, release artifact, hosted deployment, public outreach, broad product, authentication/session, or direct compatibility expansion.

