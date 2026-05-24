# ADR-0136: Agent-Native Feature Request Scaffolding Gate

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-define-agent-native-feature-request-scaffolding-gate/`

Related conversations:

- `conversations/2026-05-24-agent-native-feature-request-scaffolding-gate.md`

Related artifacts:

- `docs/agent-native-feature-request-scaffolding-gate.md`
- `docs/agent-native-feature-request-scaffolding-gate.zh-CN.md`
- `docs/agent-native-feature-request-test-workflow.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0135` selected `agent_native_requirement_test_implementation_workflow` as the next Nakama-first prototype-ready capability family after the source-first local alpha example client path.

The maintainer clarified that Nakama should be the primary reference and Pitaya should be deferred. The maintainer also clarified the product purpose: development and testing should be AI-native. A user should be able to state a backend requirement, and AI should help carry it through specification, acceptance criteria, tests, implementation, verification, and durable memory.

The repository already has a change-spec workflow and an agent-native feature request/test workflow standard. The next missing piece is an executable intake shape. Before implementing templates or commands, the repository needs a gate that defines what future scaffolding may create and what it must not create.

## Decision

Accept the agent-native feature request scaffolding gate.

The selected posture is:

```text
source_first_change_artifact_scaffolding
```

Future implementation candidates:

```text
tools/vibit scaffold feature
changes/_template/feature-request/
```

Future scaffold output must create or guide these artifacts:

```text
request.md
spec.yaml
impact.md
plan.md
checklist.md
verification.md
```

Open:

```text
M-157/W-0229 Implement agent-native feature request scaffolding
```

as the next-ready work item.

This decision does not implement scaffolding, add templates, add a scaffold command, add runtime behavior, add protocol routes, add Protobuf source, change generated output, add migrations, add dependencies, publish SDKs, add hosted surfaces, add Pitaya-style distributed architecture, or add direct Nakama/Pitaya API compatibility.

## Rationale

The product promise is not only a game server capability surface. It is a maintainable, agent-native workflow where user requirements become specs, tests, implementation, verification, and durable memory without relying on invisible chat context.

Nakama remains the primary product capability reference. Future Nakama-style features such as chat, groups, matchmaking, match runtime, and operations/admin behavior will be safer to add once the feature intake path requires explicit acceptance criteria and tests before implementation.

Pitaya remains valuable later for distributed runtime vocabulary, but it should not shape this near-term scaffold. The scaffold is about source-first repository workflow, not frontend/backend server roles, RPC, service discovery, groups, or cluster routing.

## Alternatives Considered

- Implement scaffolding immediately without a gate.
- Continue directly to chat, groups, matchmaking, match runtime, or operations/admin gates.
- Publish a public SDK or generated client package before feature intake scaffolding.
- Reactivate Pitaya distributed architecture review.

## Consequences

- `M-156/W-0228` is completed.
- `M-157/W-0229` is next-ready.
- `runtime.agent_native_feature_request_scaffolding_gate` becomes the repository check rule for this gate.
- The future implementation is bounded to docs, templates, tooling, and repository checks.
- Runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, SDK publication, generated client libraries, hosted deployments, release artifacts, chat, groups, matchmaking, match runtime, operations/admin behavior, distributed runtime, and direct compatibility remain deferred.

## Agent Reasoning Summary

The maintainer asked to continue toward Nakama and emphasized that AI-native development and testing are the product design purpose. The gate turns that product purpose into a concrete future implementation boundary: the next slice can build a source-first feature request scaffold without accidentally broadening into runtime or protocol behavior.

## Decision Weights

```yaml
decision_weights:
  ai_native_product_purpose: high
  nakama_product_alignment: high
  future_feature_velocity: high
  test_workflow_clarity: high
  implementation_boundedness: high
  direct_api_compatibility: low
confidence: high
```

## Reversal Conditions

Revisit this decision if:

- future scaffolding requires a larger command architecture than `tools/vibit scaffold feature`;
- external users need a different artifact shape than the six required change files;
- the repository adopts a broader product management workflow that supersedes source-first change artifacts.

## Follow-Up

- Complete `W-0229`: implement the agent-native feature request scaffolding inside docs, templates, tooling, and checks only.
- Keep runtime behavior, protocol changes, generated output, migrations, dependencies, SDK publication, hosted demos, distributed runtime, and direct compatibility behind later bounded work items.

