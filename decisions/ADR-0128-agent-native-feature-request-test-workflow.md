# ADR-0128: Agent-Native Feature Request And Test Workflow

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-define-agent-native-feature-request-test-workflow/`

Related conversations:

- `conversations/2026-05-24-agent-native-feature-request-test-workflow.md`

Related artifacts:

- `docs/agent-native-feature-request-test-workflow.md`
- `docs/agent-native-feature-request-test-workflow.zh-CN.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/v0.1-alpha-goal.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0127` selected a Nakama-first direction and made AI-native development and AI-native testing the product purpose. The maintainer clarified that users should be able to state a backend requirement and have AI handle the requirement spec, tests, implementation, verification, and project memory.

The repository already has change specs, ADRs, manifests, repository checks, runtime tests, protocol boundaries, generated-output rules, and a Nakama-first roadmap. What was missing was a single default workflow that tells future agents how to convert a user-facing backend request into bounded artifacts before implementation.

## Decision

Accept `docs/agent-native-feature-request-test-workflow.md` as the default workflow standard for non-trivial user-facing backend feature work.

The workflow phases are:

```text
user_requirement
requirement_spec
nakama_capability_mapping
acceptance_criteria
test_plan
tests
implementation_boundaries
verification
durable_memory
```

Required change artifacts are:

```text
request.md
spec.yaml
impact.md
plan.md
checklist.md
verification.md
```

The required `spec.yaml` expectations include:

```text
user_requirement
nakama_capability_family
acceptance_criteria
test_plan
implementation_boundaries
verification
memory_updates
```

Nakama is the primary product capability reference. Pitaya remains deferred as a future distributed architecture reference. Direct Nakama/Pitaya API compatibility remains out of scope.

This ADR opens:

```text
M-149/W-0221 Pilot Nakama-aligned feature request workflow
```

as the next-ready work item.

This decision does not add runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, startup wiring, chat, groups, matchmaking, match runtime, distributed runtime, hosted deployments, release artifacts, public announcements, paid promotion, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Continue implementing Nakama-aligned modules directly.
- Treat the existing change-spec format as sufficient without a more specific feature request/test workflow.
- Keep Nakama and Pitaya as co-equal near-term references.
- Start with a concrete operations or concurrency verification slice before defining the workflow.
- Define the workflow only in AGENTS instructions instead of an ADR-backed standard.

## Rationale

The product thesis is that vibit should be a server framework built for AI agents to maintain and extend. That thesis is not satisfied by runtime capabilities alone. Future feature work needs a repeatable path from user language to testable, bounded repository changes.

Nakama is the right near-term capability reference because it represents the game/backend product surface users recognize. Pitaya remains valuable for future distributed topology, but keeping it active now would pull scope toward cluster/RPC/frontend-backend concerns before the AI-native product loop is explicit.

Recording the workflow as a standard and check rule gives future agents a concrete intake path and lets repository checks detect drift.

## Agent Reasoning Summary

The maintainer asked to continue toward the Nakama target and previously narrowed the reference strategy to Nakama-first. The bounded continuation was to complete W-0220 by defining the workflow standard, update manifests and durable docs, add a repository check, and open a pilot work item rather than adding runtime behavior.

## Decision Weights

```yaml
decision_weights:
  ai_native_product_purpose: high
  nakama_product_alignment: high
  future_agent_reliability: high
  test_workflow_clarity: high
  runtime_scope_control: high
  implementation_cost: medium
confidence: high
```

## Consequences

- `M-148/W-0220` is completed.
- The repository has an explicit standard for user requirement -> spec -> tests -> implementation -> verification -> memory.
- Future non-trivial feature work must map to a Nakama-style capability family or explicitly say no mapping applies.
- Tests are required before or with non-trivial behavior implementation, or an explicit not-applicable rationale must be recorded.
- `runtime.agent_native_feature_request_test_workflow` checks the standard and next-ready state.
- `M-149/W-0221` is next-ready as the first pilot of the workflow.

## Reversal Conditions

Revisit this decision if:

- the workflow proves too heavy for common feature work and needs a tiered lightweight version;
- Nakama stops being the best primary product capability reference;
- Pitaya is reactivated by a later ADR for distributed architecture planning;
- repository checks cannot enforce the workflow without unacceptable maintenance overhead.

## Follow-Up

- Complete `W-0221`: apply the workflow to pilot a Nakama-aligned feature request.
- Keep broad product module expansion, runtime behavior, protocol changes, generated output, migrations, dependencies, distributed runtime, and direct compatibility behind later bounded work items.

