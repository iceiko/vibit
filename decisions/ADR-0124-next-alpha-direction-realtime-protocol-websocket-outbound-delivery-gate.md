# ADR-0124: Next Alpha Direction Realtime Protocol WebSocket Outbound Delivery Gate

Status: Accepted
Date: 2026-05-23
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-23-confirm-next-alpha-direction-after-realtime-runtime-slice/`

Related conversations:

- `conversations/2026-05-23-next-alpha-direction-realtime-protocol-websocket-outbound-delivery-gate.md`

Related artifacts:

- `runtime/internal/app/realtime/service.go`
- `runtime/internal/app/realtime/service_test.go`
- `docs/first-server-push-realtime-messaging-gate.md`
- `decisions/ADR-0123-first-server-push-realtime-messaging-runtime-slice.md`
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

`ADR-0123` completed the first server push and realtime messaging runtime slice under `runtime/internal/app/realtime`. The slice is deliberately application-owned: it validates server-authored outbound intents, resolves active recipients through the connection registry, and returns redacted delivery intents. It does not write socket frames, add Protobuf realtime messages, add a protocol bridge, wire startup, persist messages, add delivery guarantees, or expose public client publish routes.

The next prototype-ready gap is therefore the boundary between that application-owned delivery intent and future client-visible outbound delivery. If the project jumps straight to socket writes, later agents could mix protocol payload decisions, delivery authorization, WebSocket connection mechanics, and domain policy into one place.

Nakama provides the product pressure: useful shared online services eventually need notifications, streams, chat, and presence-adjacent outbound delivery. Pitaya provides the architecture pressure: acceptors, sessions, handlers, protocol serialization, backend service intent, push, groups, broadcast, and later cluster/RPC concerns must stay separated.

This work item is a direction-selection step only. It does not implement realtime protocol or WebSocket outbound delivery.

## Decision

Select:

```text
define_realtime_protocol_websocket_outbound_delivery_gate
```

as the next prototype-ready alpha direction after the first server push and realtime messaging runtime slice.

Create `M-145/W-0217` as the next ready bounded gate. The next work item should define, but not yet implement, the realtime protocol and WebSocket outbound delivery gate.

The future gate should decide the first narrow bridge from application-owned delivery intents to protocol/transport planning while preserving:

- WebSocket transport as byte movement and connection lifecycle, not message policy;
- protocol adapters as serializers and payload mappers, not recipient authorization owners;
- application-owned realtime service behavior as the authority for outbound intent validation and recipient resolution;
- server-observed connection id and epoch authority;
- no Protobuf source or generated output until a later explicitly authorized implementation slice;
- no direct Nakama/Pitaya public API compatibility;
- no premature stream subscription persistence, chat, groups, broadcast fanout, delivery guarantees, cluster/RPC, service discovery, or match runtime.

This direction selection completes `M-144/W-0216`.

## Alternatives Considered

- Strengthen concurrency and failure-path verification around the authenticated gameplay loop.
- Define a minimal operations inspection surface for active connections and delivery intents.
- Add a clearer example client or example app path.
- Define stream subscription ownership first.
- Jump directly to concrete WebSocket outbound socket writes.
- Add Protobuf realtime messages and generated output immediately.
- Jump directly to chat, groups, broadcast fanout, matchmaking, match runtime, or distributed runtime.
- Add direct Nakama/Pitaya API compatibility.

## Rationale

The previous runtime slice gave vibit a stable application-layer shape for server-authored outbound intents. The smallest useful continuation is to define the protocol/transport gate before adding any wire shape or socket write behavior.

This selection fits Nakama because it moves the project toward client-visible outbound delivery needed by notifications, streams, chat, and presence-adjacent features. It fits Pitaya because it preserves the separation between acceptor/transport mechanics, session and connection state, protocol serialization, backend service intent, push/group/broadcast vocabulary, and later cluster concerns.

The selection intentionally avoids implementing delivery guarantees, stream subscriptions, chat semantics, groups, broadcast fanout, or distributed runtime. Those are valid Nakama/Pitaya-class targets, but they should sit on an explicit protocol and transport delivery boundary.

## Agent Reasoning Summary

The maintainer asked to continue, keep Nakama/Pitaya alignment explicit, commit, and push. After the realtime runtime slice, the next bounded step is a gate that defines how future protocol and WebSocket delivery should relate to the application-owned realtime service. Implementing protocol or socket writes in this direction-selection slice would skip the project's contract-first and gate-density discipline.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  protocol_transport_boundary_clarity: high
  implementation_scope_control: high
  generated_output_risk: none_in_this_step
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `M-144/W-0216` records the selected next alpha direction.
- `M-145/W-0217` becomes the next ready work item.
- The next work item should define the realtime protocol and WebSocket outbound delivery gate before any implementation.
- WebSocket outbound delivery, concrete socket writes, Protobuf source, generated output, protocol bridge, application bootstrap handler, startup wiring, persistence, migrations, dependencies, authentication/session behavior changes, route-protection changes, stream subscriptions, chat, groups, broadcast fanout, offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, backpressure, distributed fanout, matchmaking, match runtime, hosted deployments, release artifacts, public announcements, paid promotion, blob/S3 storage, and direct Nakama/Pitaya API compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- external alpha feedback shows operations visibility, example clients, or failure verification is a stronger blocker than outbound delivery planning;
- the application-owned realtime service handoff shape changes before protocol delivery planning begins;
- the maintainer explicitly selects stream subscription ownership, chat, groups, operations inspection, SDK/example work, matchmaking, match runtime, or direct compatibility instead;
- a later ADR changes the prototype-ready execution plan.

## Follow-Up

- Complete `W-0217`: define the realtime protocol and WebSocket outbound delivery gate.
- Keep Protobuf source, generated output, protocol bridges, WebSocket outbound writers, startup wiring, stream subscriptions, chat semantics, offline inboxes, delivery guarantees, distributed fanout, matchmaking, match runtime, and direct compatibility behind later bounded work items.
