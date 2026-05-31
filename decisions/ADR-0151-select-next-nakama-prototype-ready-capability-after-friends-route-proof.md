# ADR-0151: Select Next Nakama Prototype-Ready Capability After Friends Route Proof

Status: Accepted
Date: 2026-05-31
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-select-next-nakama-prototype-ready-capability-after-friends-route-proof/`

Related conversations:

- `conversations/2026-05-26-select-next-nakama-prototype-ready-capability-after-friends-route-proof.md`

Related artifacts:

- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/product-maturity-milestones.md`
- `docs/alpha-developer-flow.md`
- `docs/v0.1-alpha-goal.md`
- `examples/local-alpha-example-client.sh`
- `examples/local-alpha-request-loop.sh`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0150` completed the friends relationship protocol route local proof and opened `M-171/W-0243` as a selection slice. The local alpha now proves login, connection binding, protected inventory, presence, own-player storage objects, logout and revoked-token rejection, realtime outbound delivery, source-first example-client flow, feature request scaffolding, and protected friends relationship routes.

The next prototype-ready decision should choose the most useful Nakama-first gap without prematurely adding groups, parties, chat, leaderboards, matchmaking, match runtime, operations implementation, SDK publication, hosted releases, or distributed runtime.

The strongest current gap is minimum operations inspection. The alpha has `/healthz`, `/readyz`, `/version`, and `/configz`, but it does not yet define what a prototype author can inspect when troubleshooting players, sessions, tokens, connections, routes, storage objects, friends state, realtime state, or configuration posture.

Nakama capability family:

```text
admin_console_metrics_observability_and_operations
```

Pitaya remains deferred as a future distributed architecture reference.

## Decision

Select:

```text
define_minimum_operations_inspection_surface_gate
```

as the next direction.

Open:

```text
M-172/W-0244 Define minimum operations inspection surface gate
```

as the next-ready work item.

The follow-up gate must define the acceptable source-first operations inspection posture, inspectable state categories, redaction rules, test expectations, ownership boundaries, and stop conditions before implementation. This decision does not implement an operations surface, admin console, metrics endpoint, observability pipeline, protocol route, runtime behavior, generated output, dependency, hosted deployment, SDK, release artifact, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Define groups relationship lifecycle.
- Define parties lifecycle.
- Define chat streams.
- Define leaderboards or tournaments.
- Define matchmaking.
- Define authoritative match runtime.
- Publish an SDK or generated client library.
- Reactivate Pitaya distributed architecture review.

## Rationale

The current source-first alpha is now broad enough that a developer can exercise meaningful local gameplay behavior. That breadth also creates a practical debugging problem: before adding more product modules, vibit should define how a local prototype author and future agent inspect operational state without leaking sensitive material or inventing ad hoc admin surfaces.

Groups, parties, chat, leaderboards, matchmaking, and match runtime are product-useful, but they add more state and failure modes. A minimum operations inspection gate creates the vocabulary for safe introspection first. SDK or generated client publication is still premature while the local alpha surface is changing quickly. Pitaya-style distributed architecture remains a later architecture reference.

## Agent Reasoning Summary

The maintainer asked to continue and to commit and push. After the friends route proof, the next highest-leverage Nakama-aligned step is not another domain feature; it is a bounded operations inspection gate. This keeps the alpha easier to understand while preserving the repository's agent-native requirement/spec/test/verification flow.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  prototype_ready_value: high
  operations_debuggability_value: high
  implementation_boundedness: high
  redaction_importance: high
  avoids_premature_product_breadth: high
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `M-171/W-0243` is completed.
- `M-172/W-0244` is next-ready.
- The selected capability family is `admin_console_metrics_observability_and_operations`.
- The next work should define a gate for a source-first minimum operations inspection surface.
- Runtime behavior, operations/admin endpoints, metrics endpoints, observability pipelines, protocol routes, Protobuf source, generated output, migrations, dependencies, persistence, startup wiring, repository interface changes, PostgreSQL adapter changes, SDK publication, hosted deployments, release artifacts, chat, groups, parties, matchmaking, match runtime, distributed runtime, and direct compatibility remain deferred.
- `runtime.next_nakama_prototype_ready_capability_after_friends_route_proof` checks the selection records and next-ready state.

## Reversal Conditions

Revisit this decision if:

- external alpha users report a specific product capability, such as groups or matchmaking, as a stronger blocker than operations inspection;
- an explicit maintainer decision authorizes a broader product module before operations inspection work;
- the future gate finds that useful inspection requires a prerequisite runtime state ownership boundary that should be split first.

## Follow-Up

- Complete `W-0244`: define the minimum operations inspection surface gate.
- Keep operations implementation, admin consoles, metrics endpoints, observability pipelines, hosted operations surfaces, protocol changes, runtime changes, dependencies, distributed runtime, SDK publication, and direct compatibility behind later bounded work items.
