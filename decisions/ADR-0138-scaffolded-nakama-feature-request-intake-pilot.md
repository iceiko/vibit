# ADR-0138: Scaffolded Nakama Feature Request Intake Pilot

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-pilot-scaffolded-nakama-feature-request-intake/`

Related conversations:

- `conversations/2026-05-24-scaffolded-nakama-feature-request-intake-pilot.md`

Related artifacts:

- `docs/agent-native-feature-request-scaffolding.md`
- `docs/agent-native-feature-request-test-workflow.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0137` implemented the source-first feature request scaffold. The next-ready work item, `M-158/W-0230`, required using that scaffold on one bounded Nakama-style request intake before broad product module expansion.

The maintainer has clarified two durable directions:

- Nakama is the primary product capability reference.
- AI-native development and AI-native testing are the product purpose: user requirement to spec, acceptance criteria, tests, implementation, verification, and durable memory.

Pitaya remains deferred as a future distributed architecture reference. This pilot must not add runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, SDKs, hosted surfaces, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Decision

Pilot the scaffold on this selected user requirement:

```text
As a game developer building on vibit, I want a future player friendship relationship lifecycle so players can request, accept, reject, remove, and block social relationships through server-authoritative rules before broader groups, parties, chat, matchmaking, or match runtime features depend on social graph state.
```

The selected Nakama capability family is:

```text
friends_groups_and_parties
```

The selected follow-up direction is:

```text
define_friends_relationship_lifecycle_gate
```

Open:

```text
M-159/W-0231 Define friends relationship lifecycle gate
```

as the next-ready work item.

This pilot records the request, Nakama mapping, acceptance criteria, test plan, implementation boundaries, verification expectations, and durable memory. It does not define the gate itself, implement friendship behavior, add protocol routes, add Protobuf source, change generated output, add migrations, add dependencies, change startup wiring, add persistence, publish SDKs, add hosted surfaces, add distributed runtime, or add direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Pilot the scaffold on chat or group messaging.
- Pilot the scaffold on leaderboards or tournaments.
- Pilot the scaffold on matchmaking.
- Pilot the scaffold on authoritative match runtime.
- Pilot the scaffold on minimal operations inspection.
- Reactivate Pitaya distributed architecture review.

## Rationale

Friend relationship lifecycle is a useful next Nakama-first intake because it is a core social graph primitive and can be gated before implementation. It naturally precedes groups, parties, chat targeting, invites, matchmaking filters, and match social context, but it can still remain small enough for a contract-first gate.

Chat, matchmaking, and match runtime are product-class capabilities, but they are broader and depend on stronger social graph and realtime semantics. Operations inspection remains important for later production-candidate maturity. Pitaya remains deferred because this decision is about single-process product capability planning, not distributed topology.

The selected follow-up is a gate, not implementation. That keeps the next step aligned with the agent-native workflow: define requirements, acceptance criteria, contracts, tests, data posture, protocol posture, and stop conditions before adding code.

## Agent Reasoning Summary

The maintainer asked to continue toward Nakama and to make AI-native development/testing the organizing product principle. After implementing the scaffold, the next best proof is to use it on a concrete product capability and open a bounded gate. Friendship lifecycle is central enough to Nakama-class game backends and narrow enough to gate safely.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  ai_native_workflow_proof: high
  prototype_ready_value: high
  implementation_boundedness: high
  future_testability: high
  avoids_premature_product_breadth: high
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `M-158/W-0230` is completed.
- `M-159/W-0231` is next-ready.
- The selected capability family is `friends_groups_and_parties`.
- The next work should define a friendship lifecycle gate before implementation.
- Runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, groups, parties, chat, matchmaking, match runtime, operations/admin behavior, SDK publication, hosted deployments, release artifacts, distributed runtime, and direct compatibility remain deferred.
- `runtime.scaffolded_nakama_feature_request_intake_pilot` checks the pilot records and next-ready state.

## Reversal Conditions

Revisit this decision if:

- alpha users report a more urgent Nakama-style blocker than friendship lifecycle;
- the future gate finds that friendship lifecycle depends on an unresolved player/profile/account contract gap;
- an explicit maintainer decision authorizes a different social, realtime, or operations capability before friendship work;
- a later ADR reactivates Pitaya for distributed architecture planning.

## Follow-Up

- Complete `W-0231`: define the friends relationship lifecycle gate.
- Keep friendship implementation, protocol changes, generated output, migrations, dependencies, SDK publication, hosted demos, distributed runtime, and direct compatibility behind later bounded work items.
