# ADR-0122: First Server Push And Realtime Messaging Gate

Status: Accepted
Date: 2026-05-23
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-23-define-first-server-push-realtime-messaging-gate/`

Related conversations:

- `conversations/2026-05-23-first-server-push-realtime-messaging-gate.md`

Related artifacts:

- `docs/first-server-push-realtime-messaging-gate.md`
- `docs/first-server-push-realtime-messaging-gate.zh-CN.md`
- `decisions/ADR-0121-next-alpha-direction-first-server-push-realtime-messaging-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0121` selected `define_first_server_push_realtime_messaging_gate` as the next prototype-ready alpha direction after the storage objects local proof.

The project now needs a first explicit outbound realtime vocabulary before implementing server push, broadcast, streams, chat, notifications, or presence-adjacent server facts. Without this boundary, future work could put recipient policy in WebSocket transport, let protocol adapters own business authorization, or copy Nakama/Pitaya public APIs instead of adapting their capability and architecture pressure.

## Decision

Accept `docs/first-server-push-realtime-messaging-gate.md` as the gate for the first server push and realtime messaging vocabulary.

The gate defines:

- future application ownership under `runtime/internal/app/realtime`;
- future protocol and generated-output candidates under `proto/vibit/realtime/v1/` and `runtime/internal/generated/proto/vibit/realtime/v1/`;
- future protocol bridge, application handler, and transport delivery candidates;
- first message intent vocabulary: `server_notice`, `domain_event_push`, `stream_message`, and `presence_signal`;
- first target vocabulary: `connection_id_and_epoch`, `player_current_connections`, and `stream_subscribers`;
- conservative sender authority and identity rules;
- Nakama/Pitaya reference mapping without direct compatibility;
- stop conditions that keep implementation, protocol, generated output, persistence, delivery guarantees, distributed fanout, matchmaking, match runtime, and direct compatibility deferred.

Create `M-143/W-0215 Implement first server push and realtime messaging runtime slice` as the next ready work item.

This decision completes `M-142/W-0214`.

## Alternatives Considered

- Implement server push immediately inside the WebSocket transport.
- Add Protobuf realtime messages and generated output before the application boundary is defined.
- Copy Nakama notification/channel/stream APIs or Pitaya push/group/broadcast APIs directly.
- Defer outbound realtime messaging entirely and move next to operations, SDK/example work, matchmaking, or match runtime.

## Rationale

Nakama shows that a useful game backend eventually needs outbound realtime surfaces such as notifications, streams, chat, and presence-adjacent messages.

Pitaya shows that Go game server frameworks need clear separation among acceptors, sessions, handlers, push, groups, broadcast, backend services, and cluster/RPC concerns.

vibit should adopt the capability and layering lessons without copying public APIs. A gate first is the smallest safe continuation because it records ownership, vocabulary, allowed future files, verification expectations, and stop conditions before any runtime or protocol behavior exists.

## Agent Reasoning Summary

The maintainer asked to continue while staying aligned with Nakama and Pitaya. After storage objects were proven through the local alpha route flow, the next useful Nakama-class family is outbound realtime capability: notifications, streams, chat, and presence-adjacent server facts. Pitaya adds the matching architecture pressure: push, groups, broadcast, sessions, handlers, protocol adaptation, backend services, and future cluster/RPC concerns must stay separated.

The smallest correct continuation is therefore a gate, not implementation. This records vibit's own vocabulary and future W-0215 slice while preserving deferrals for protocol expansion, generated output, transport delivery, persistence, distributed fanout, broad social modules, matchmaking, match runtime, and direct Nakama/Pitaya compatibility.

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

- `M-142/W-0214` is completed.
- `M-143/W-0215` becomes next-ready.
- `runtime.first_server_push_realtime_messaging_gate` becomes the repository check rule for this boundary.
- Runtime behavior, transport delivery, Protobuf source, generated output, protocol routes, startup wiring, persistence, migrations, dependencies, authentication/session changes, hosted deployment, release artifacts, public announcements, paid promotion, matchmaking, match runtime, distributed runtime, broad social modules, blob/S3 storage, and direct Nakama/Pitaya API compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- prototype feedback shows operations, SDK/client ergonomics, or failure verification is a stronger blocker than outbound realtime behavior;
- existing connection lifecycle behavior cannot support a safe single-process outbound runtime slice;
- a later protocol ADR changes the envelope model before realtime payloads are implemented;
- the maintainer explicitly chooses matchmaking, match runtime, operations inspection, SDK work, or direct compatibility instead.

## Follow-Up

- Complete `W-0215`: implement the smallest first server push and realtime messaging runtime slice authorized by this gate.
- Keep broad chat/social features, stream subscriptions, offline inboxes, delivery guarantees, distributed fanout, protocol expansion, generated output, matchmaking, match runtime, and direct compatibility behind later bounded work items unless explicitly authorized.
