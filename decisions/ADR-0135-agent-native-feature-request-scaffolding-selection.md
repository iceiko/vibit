# ADR-0135: Agent-Native Feature Request Scaffolding Selection

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-select-next-nakama-prototype-ready-capability-after-local-alpha-example-client-path/`

Related conversations:

- `conversations/2026-05-24-agent-native-feature-request-scaffolding-selection.md`

Related artifacts:

- `docs/agent-native-feature-request-test-workflow.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/product-maturity-milestones.md`
- `docs/alpha-developer-flow.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0134` implemented the source-first local alpha example client path. The repository now has a readable way to exercise the local alpha loop for onboarding, device credential login, connection binding, protected inventory, protected presence, own-player storage objects, logout, rejected post-logout requests, and failure-path redaction.

The maintainer clarified that Nakama should remain the primary reference and Pitaya should be deferred. The maintainer also clarified the product purpose: development and testing should be AI-native. A user should state a backend requirement, and AI should carry that requirement through specification, acceptance criteria, tests, implementation, verification, and durable project memory.

The next prototype-ready decision should choose the highest-leverage Nakama-first gap without prematurely adding broad chat, groups, matchmaking, match runtime, operations/admin behavior, SDK publication, hosted releases, or distributed runtime.

Nakama capability family:

```text
agent_native_requirement_test_implementation_workflow
```

Pitaya remains deferred as a future distributed architecture reference.

## Decision

Select:

```text
define_agent_native_feature_request_scaffolding_gate
```

as the next direction.

Open:

```text
M-156/W-0228 Define agent-native feature request scaffolding gate
```

as the next-ready work item.

The follow-up gate must define the acceptable scaffolding posture for future feature requests: required artifacts, template or tool ownership, redaction rules, check coverage, test expectations, stop conditions, and implementation deferrals. This decision does not define the gate itself, implement scaffolding, add runtime behavior, add protocol routes, change generated output, add migrations, add dependencies, publish SDKs, add hosted surfaces, or add direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Define the minimum operations inspection surface.
- Define notifications or realtime delivery flow proof.
- Define chat or group messaging.
- Define friends, groups, parties, leaderboards, matchmaking, or match runtime gates.
- Define a public client package boundary.
- Reactivate Pitaya distributed architecture review.

## Rationale

The current local alpha can now be read and exercised. The next gap is making vibit's core product promise executable for users: a new requirement should reliably become a bounded spec, acceptance criteria, tests, implementation boundary, verification record, and durable memory before coding expands the server.

Broad Nakama-style services remain important, but adding them before the feature-request scaffolding path would increase product breadth while leaving the AI-native development loop too dependent on agent discipline and chat context. The scaffolding gate should make that loop explicit enough for future product slices.

Operations inspection is important before production-candidate maturity. Public SDKs and client packages are useful later. Chat, groups, matchmaking, and match runtime are product-class capabilities. Pitaya distributed architecture is still future work. None of those should outrank the product's requirement-to-test workflow at this stage.

## Agent Reasoning Summary

The maintainer asked to keep advancing toward Nakama and emphasized that AI-native development and testing are the product design purpose. After the example path, the next highest-leverage step is to make future feature intake executable and checkable so later Nakama-style capability work starts from a requirement spec and tests rather than ad hoc implementation.

## Decision Weights

```yaml
decision_weights:
  ai_native_product_purpose: high
  nakama_product_alignment: high
  future_feature_velocity: high
  test_workflow_clarity: high
  implementation_boundedness: high
  avoids_premature_product_breadth: high
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `M-155/W-0227` is completed.
- `M-156/W-0228` is next-ready.
- The selected capability family is `agent_native_requirement_test_implementation_workflow`.
- The next work should define a gate for feature request scaffolding before implementing templates or tooling.
- Runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, SDK publication, hosted deployments, release artifacts, chat, groups, matchmaking, match runtime, operations/admin behavior, distributed runtime, and direct compatibility remain deferred.
- `runtime.agent_native_feature_request_scaffolding_selection` checks the selection records and next-ready state.

## Reversal Conditions

Revisit this decision if:

- external alpha users report operations/admin inspection as a stronger blocker than feature intake scaffolding;
- an explicit maintainer decision authorizes a broad product module before scaffolding work;
- the future gate finds that useful scaffolding requires a larger tooling architecture decision first.

## Follow-Up

- Complete `W-0228`: define the agent-native feature request scaffolding gate.
- Keep scaffolding implementation, runtime behavior, protocol changes, generated output, migrations, dependencies, SDK publication, hosted demos, distributed runtime, and direct compatibility behind later bounded work items.
