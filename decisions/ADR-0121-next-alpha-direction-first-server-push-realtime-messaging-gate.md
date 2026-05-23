# ADR-0121: Next Alpha Direction First Server Push Realtime Messaging Gate

Status: Accepted
Date: 2026-05-23
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-23-confirm-next-alpha-direction-after-storage-objects-local-proof/`

Related conversations:

- `conversations/2026-05-23-next-alpha-direction-first-server-push-realtime-messaging-gate.md`

Related artifacts:

- `docs/product-maturity-milestones.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/prototype-ready-foundation-execution-plan.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `AGENTS.md`
- `runtime/AGENTS.md`
- `README.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-140/W-0212` proved the completed own-player storage object route family in the local alpha WebSocket/Protobuf request flow. That closed the first general durable game-state capability path beyond inventory.

The prototype-ready execution plan now points at the next shared online-service gap. A serious small prototype needs more than durable state and request/response commands. It needs a first explicit vocabulary for outbound realtime behavior so later work can reason about notifications, streams, chat, presence-adjacent signals, and broadcast without hiding those concerns inside the WebSocket transport.

Nakama provides the product capability pressure: storage, presence, status, notifications, chat, streams, and realtime messaging are common game-backend services. Pitaya provides the architecture pressure: acceptors, sessions, handlers, push, groups, broadcast, serializers, backend services, and later cluster/RPC topology must remain distinct.

This work item is a direction-selection step only. It does not implement server push or realtime messaging.

## Decision

Select:

```text
define_first_server_push_realtime_messaging_gate
```

as the next prototype-ready alpha direction after storage object local proof.

Create `M-142/W-0214` as the next ready bounded gate. The next work item should define, but not yet implement, the first server push / realtime messaging gate.

The future gate should decide the first narrow vocabulary for outbound realtime behavior while preserving:

- WebSocket transport as credential-neutral byte movement and connection lifecycle;
- protocol adapters as serializers and route/payload mappers;
- application-owned policy for who can receive a message;
- backend service ownership for message intent and invariants;
- no direct Nakama/Pitaya public API compatibility;
- no premature cluster/RPC, distributed groups, or broad chat/social/matchmaking implementation.

This direction selection completes `M-141/W-0213`.

## Alternatives Considered

- Strengthen concurrency and failure-path verification around the existing authenticated loop.
- Define the minimum operations inspection surface needed before serious prototype use.
- Add a clearer example client or example app path.
- Expand presence/status semantics first.
- Jump directly to chat.
- Jump directly to matchmaking or match runtime.
- Start Pitaya-style distributed frontend/backend, RPC, groups, or cluster service discovery.
- Add direct Nakama/Pitaya API compatibility.

## Rationale

Storage objects made vibit more useful for durable game state. The next prototype-ready gap is outbound realtime usefulness. Choosing a gate first keeps the work small and contract-first: the project can decide vocabulary, ownership, non-goals, and verification expectations before adding any protocol shape or runtime behavior.

This selection fits Nakama because it moves toward shared online-service coverage such as notifications, streams, chat, and presence-adjacent realtime behavior. It fits Pitaya because it forces push/broadcast/group vocabulary to be defined across transport, session, handler, protocol, and backend layers before implementation.

The selection intentionally avoids jumping to matchmaking, match runtime, broad chat, or distributed runtime. Those are valid Nakama/Pitaya-class targets, but they should not sit on an undefined outbound messaging foundation.

## Agent Reasoning Summary

The maintainer asked to continue and keep Nakama/Pitaya alignment explicit. After the storage object local proof, the prototype-ready roadmap lists the first server push, stream, broadcast, or realtime messaging vocabulary as the next product-useful family. The smallest correct continuation is to open a gate for that family, not to implement it directly.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  boundary_clarity: high
  implementation_scope_control: high
  protocol_change_risk: none_in_this_step
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `M-141/W-0213` records the selected next alpha direction.
- `M-142/W-0214` becomes the next ready work item.
- The next work item should define the first server push / realtime messaging gate before any implementation.
- Runtime behavior, protocol routes, Protobuf messages, generated output, storage behavior, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session changes, hosted deployments, release artifacts, public announcements, paid promotion, matchmaking, match runtime, distributed runtime, broad social/competitive modules, blob/S3 storage, and direct Nakama/Pitaya API compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- external alpha feedback shows setup, operations, or failure verification is a stronger blocker than outbound realtime vocabulary;
- the existing authenticated loop cannot safely support any outbound realtime gate without additional lifecycle hardening;
- the maintainer explicitly selects matchmaking, match runtime, operations inspection, SDK/example work, or direct compatibility instead;
- a later ADR changes the prototype-ready execution plan.

## Follow-Up

- Complete `W-0214`: define the first server push and realtime messaging gate.
- Keep implementation, protocol messages/routes, generated output, dependencies, persistence, delivery guarantees, offline inboxes, acknowledgements, ordering, retries, backpressure, distributed fanout, chat, matchmaking, match runtime, and direct compatibility behind later bounded work items.
