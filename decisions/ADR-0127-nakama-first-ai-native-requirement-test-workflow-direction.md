# ADR-0127: Nakama-First AI-Native Requirement Test Workflow Direction

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-23-confirm-next-alpha-direction-after-realtime-outbound-delivery-slice/`

Related conversations:

- `conversations/2026-05-24-nakama-first-ai-native-feature-workflow-direction.md`

Related artifacts:

- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`
- `docs/reference-game-server-alignment.md`
- `docs/reference-game-server-alignment.zh-CN.md`
- `docs/v0.1-alpha-goal.md`
- `docs/v0.1-alpha-goal.zh-CN.md`
- `docs/product-maturity-milestones.md`
- `docs/product-maturity-milestones.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

After `ADR-0126`, vibit has a narrow realtime Protobuf and WebSocket outbound delivery slice. The next-ready item was `W-0219 Confirm next alpha direction after realtime outbound delivery slice`.

The maintainer clarified the product direction:

- reference Nakama now;
- leave Pitaya for later;
- make AI-native development and AI-native testing the product design purpose;
- make all architecture and implementation serve the workflow where a user states a requirement and AI handles the rest.

Earlier roadmap language treated Nakama and Pitaya as active paired references for product planning. That pairing helped early layering decisions, but it now risks splitting near-term scope between product capability breadth and distributed Go game server topology. The project needs a clearer current driver before adding more modules.

## Decision

Select:

```text
define_agent_native_feature_request_test_workflow
```

as the next prototype-ready alpha direction.

This opens:

```text
M-148/W-0220 Define agent-native feature request and test workflow
```

as the next-ready work item.

The reference posture is refined:

- Nakama is the primary product capability reference.
- Pitaya is deferred as a future architecture reference for distributed runtime concerns.
- vibit remains vibit-native and does not pursue direct Nakama or Pitaya API compatibility.

The product purpose is refined:

```text
user requirement -> AI-written bounded requirement spec -> AI-written acceptance criteria -> AI-written test plan -> AI-written or updated tests -> AI implementation -> AI-run verification -> AI-updated docs/manifests/ADRs/change records
```

AI-native development and AI-native testing are not side features. They are the reason the architecture exists.

This ADR does not add runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, startup wiring, chat, groups, matchmaking, match runtime, distributed runtime, hosted deployments, release artifacts, public announcements, paid promotion, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Strengthen authenticated local loop concurrency and failure verification next.
- Add a clearer example client or example app path next.
- Define a minimal operations inspection surface next.
- Define notifications or realtime delivery flow proof next.
- Define chat, groups, or matchmaking next.
- Start Pitaya-style distributed architecture, RPC, frontend/backend roles, groups, broadcast, or service discovery next.
- Keep Nakama and Pitaya as co-equal near-term references.

## Rationale

Nakama is the better near-term reference because it represents the broad product surface users expect from an open-source game/backend server: users, sessions, storage, realtime messaging, social features, matchmaking, authoritative matches, leaderboards, tournaments, operations, and developer experience.

Pitaya is still useful, but mostly for later distributed architecture vocabulary: acceptors, sessions, routes, handlers, remotes/RPC, frontend/backend roles, groups, broadcast, serializers, and service discovery. Pulling that in now would compete with the immediate product promise.

The maintainer's current product thesis is not merely "build another game server." It is "build a server framework where AI can safely take a user requirement through specification, tests, implementation, verification, and durable memory." A Nakama-first roadmap without an AI-native requirement/test workflow would still leave the core differentiation underspecified.

The selected direction gives future feature work a standard intake path before broad modules expand.

## Agent Reasoning Summary

The maintainer explicitly corrected the reference strategy and product purpose. The safe continuation was to complete the current direction-selection item, avoid runtime scope, update the durable roadmap and manifests, and open a workflow-standard work item. That keeps the next step focused on how AI should turn user requests into tested, verified changes before implementing more product features.

## Decision Weights

```yaml
decision_weights:
  product_focus: high
  ai_native_differentiation: high
  test_workflow_importance: high
  nakama_product_reference_value: high
  pitaya_near_term_scope_risk: high
  runtime_scope_control: high
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `M-147/W-0219` is completed.
- `M-148/W-0220` is next-ready.
- Roadmap documents now describe Nakama as the primary product capability reference.
- Pitaya is deferred until a later ADR reactivates it for distributed architecture planning.
- Future non-trivial user-facing feature work should be shaped around requirement spec, acceptance criteria, test plan, tests, implementation, verification, and memory updates.
- Repository checks track this direction-selection slice.

## Reversal Conditions

Revisit this decision if:

- near-term contributors need distributed runtime topology before product feature workflow is usable;
- Nakama stops being the best broad product capability reference;
- a later architecture decision explicitly reactivates Pitaya for cluster/RPC/frontend-backend planning;
- the AI-native workflow proves too heavy and needs a thinner tiered version for small changes.

## Follow-Up

- Complete `W-0220`: define the agent-native feature request and test workflow.
- Keep runtime behavior, protocol changes, generated output, migrations, dependencies, broad modules, Pitaya-style distributed architecture, and direct compatibility behind later bounded work items.
