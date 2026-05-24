# ADR-0129: Nakama-Aligned Presence Status Workflow Pilot

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-pilot-nakama-aligned-feature-request-workflow/`

Related conversations:

- `conversations/2026-05-24-nakama-presence-status-workflow-pilot.md`

Related artifacts:

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

`ADR-0128` defined the agent-native feature request and test workflow. The next-ready work item, `M-149/W-0221`, required applying that workflow to one concrete Nakama-aligned capability request and selecting a bounded follow-up slice.

The repository already has a protected self-presence query, active connection registry, first-message connection binding, authentication route protection, and a local authenticated gameplay E2E path that observes online presence. Presence/status is a recognizable Nakama-style online-service capability and is already close enough to the alpha flow that it can be hardened without opening broad product scope.

## Decision

Pilot the agent-native workflow on this selected user requirement:

```text
As a developer using the local alpha request flow, I want the server to prove a player's self-presence/status through authenticated requests, including online after binding and offline after close or invalidation, so I can trust the foundation before broader Nakama-style realtime, social, and multiplayer features build on it.
```

The selected Nakama capability family is:

```text
presence_status_and_notifications
```

The selected follow-up direction is:

```text
presence_status_local_proof_hardening
```

Open:

```text
M-150/W-0222 Harden presence status local proof through close and offline cases
```

as the next-ready work item.

This pilot records the request, Nakama mapping, acceptance criteria, test plan, implementation boundaries, verification expectations, and durable memory. It does not add runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, startup wiring, persistence, presence subscriptions, status broadcast fanout, chat, groups, matchmaking, match runtime, distributed runtime, hosted deployments, release artifacts, public announcements, paid promotion, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Pilot the workflow on chat or stream subscriptions.
- Pilot the workflow on matchmaking or match runtime.
- Pilot the workflow on operations/admin inspection.
- Pilot the workflow on a clearer example client or SDK path.
- Reactivate Pitaya for distributed architecture planning.
- Treat the existing presence proof as enough and move directly to a broad social module.

## Rationale

Presence/status is a good first pilot because it is product-recognizable in Nakama terms and already has local alpha foundations in vibit. It lets the workflow produce a useful, bounded next slice without forcing a new protocol surface or broad module.

Chat, groups, matchmaking, match runtime, SDKs, and operations are valid future Nakama-style families, but they are broader and should not be the first test of the workflow. Pitaya remains deferred because this decision is about product capability proof, not distributed topology.

The selected follow-up asks for proof hardening, not feature expansion. That keeps the next step test-heavy and compatible with the maintainer's requirement that AI-native development and testing be central to the product.

## Agent Reasoning Summary

The maintainer asked for sustained progress toward Nakama. The active work item required applying the newly defined AI-native workflow to a concrete capability. The existing codebase showed that presence/status already exists in a minimal form, so the best next move is to harden and prove that capability rather than invent a new broad module.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  alpha_usefulness: high
  testability: high
  bounded_scope: high
  reuse_existing_foundation: high
  runtime_scope_control_for_pilot: high
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `M-149/W-0221` is completed.
- `M-150/W-0222` is next-ready.
- The first workflow pilot maps to `presence_status_and_notifications`.
- Future presence/status hardening should focus on online, close, invalidation, offline, redaction, and local alpha proof behavior.
- Broad subscriptions, broadcasts, chat, groups, matchmaking, match runtime, persistence, distributed runtime, SDKs, operations, and direct compatibility remain deferred.
- `runtime.nakama_aligned_feature_request_workflow_pilot` checks the pilot records and next-ready state.

## Reversal Conditions

Revisit this decision if:

- presence/status proof hardening proves already sufficient and another capability is a better first implementation pilot;
- alpha users report a more urgent Nakama-style blocker;
- the workflow is too heavy and needs a lighter tier before runtime work;
- a later ADR explicitly reactivates Pitaya for distributed architecture planning.

## Follow-Up

- Complete `W-0222`: harden presence status local proof through close and offline cases.
- Keep new protocol routes, Protobuf, generated output, migrations, dependencies, subscriptions, broadcasts, chat, groups, matchmaking, match runtime, persistence, distributed runtime, and direct compatibility behind later bounded work items.
