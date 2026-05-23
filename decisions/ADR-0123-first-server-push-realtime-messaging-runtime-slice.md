# ADR-0123: First Server Push And Realtime Messaging Runtime Slice

Status: Accepted
Date: 2026-05-23
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-23-implement-first-server-push-realtime-messaging-runtime-slice/`

Related conversations:

- `conversations/2026-05-23-first-server-push-realtime-messaging-runtime-slice.md`

Related artifacts:

- `runtime/internal/app/realtime/service.go`
- `runtime/internal/app/realtime/service_test.go`
- `docs/first-server-push-realtime-messaging-gate.md`
- `decisions/ADR-0122-first-server-push-realtime-messaging-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0122` defined the first server push and realtime messaging gate. It authorized a later bounded runtime slice under `runtime/internal/app/realtime` while keeping protocol routes, Protobuf source, generated output, WebSocket outbound delivery, startup wiring, persistence, migrations, dependencies, delivery guarantees, distributed fanout, broad social modules, matchmaking, match runtime, and direct Nakama/Pitaya API compatibility deferred.

The runtime already has a single-process active connection registry under `runtime/internal/app/connection`. That registry records server-observed connection ids and epochs, validated player linkage, runtime session ids, and access-token record ids without making transport own identity or policy.

The next useful step is an application-owned realtime service that validates server-authored outbound message intents and resolves allowed active recipients into redacted delivery intents. It should give future protocol and transport adapter work a stable handoff shape without sending socket frames yet.

Nakama provides the product pressure: useful game backends need notifications, streams, chat, and presence-adjacent outbound messages. Pitaya provides the layering pressure: acceptors, sessions, handlers, backend services, push, groups, broadcast, and later cluster/RPC concerns must remain separate. vibit adapts both references without direct public API compatibility.

## Decision

Implement the first server push and realtime messaging runtime slice under:

```text
runtime/internal/app/realtime
```

The implementation adds:

- `Service` and `NewService`;
- `AcceptServerMessage`;
- vibit-native intent kinds `server_notice`, `domain_event_push`, `stream_message`, and `presence_signal`;
- target kinds `connection_id_and_epoch`, `player_current_connections`, and future-only `stream_subscribers`;
- delivery outcomes `accepted`, `no_active_recipient`, `recipient_not_authorized`, `payload_invalid`, and `delivery_unavailable`;
- application-owned recipient resolution from `runtime/internal/app/connection` registry state;
- server-authorized sender validation using validated service or admin identity;
- metadata-only and validated-player sender refusal before registry resolution;
- single-process bound-connection and player-current-connection recipient resolution;
- future-only stream target rejection until subscription ownership is separately defined;
- redacted delivery results and service errors;
- copied payload bytes in returned delivery intents;
- focused tests for recipient validation, metadata-only identity refusal, no active recipient, redaction, and result copying.

This ADR does not add WebSocket outbound delivery code, concrete socket writes, protocol routes, Protobuf source, generated output, startup wiring, persistence, migrations, dependencies, authentication/session behavior changes, route-protection changes, public client publish routes, stream subscription persistence, offline inboxes, acknowledgements, retries, ordering guarantees, durable offsets, backpressure, distributed fanout, frontend/backend split, RPC, service discovery, cluster groups, broad chat/social behavior, matchmaking, match runtime, hosted deployments, release artifacts, public announcements, paid promotion, blob/S3 storage, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add WebSocket outbound delivery in the same slice.
- Add Protobuf realtime messages and generated output together with the service.
- Expose a public client publish or subscribe route immediately.
- Treat player-authenticated requests as sufficient authority to publish realtime facts.
- Use stream subscribers as an implemented target without defining subscription ownership.
- Copy Nakama notification/channel/stream or Pitaya push/group/broadcast APIs directly.

## Rationale

The smallest runtime behavior authorized by `ADR-0122` is an application service that accepts server-authored intents and resolves server-observed recipients. That gives future protocol and transport slices a concrete handoff shape while preserving the transport and protocol deferrals.

Requiring a validated service or admin identity keeps this first posture server-authored. A validated player identity is not enough to publish server facts, and metadata-only `player_id`, `session_id`, or connection metadata is never proof.

Using the existing connection registry preserves Pitaya-style separation: acceptors and WebSocket transport own connection mechanics, the registry owns server-observed connection state, the realtime service owns policy and recipient resolution, and later protocol/transport adapters can map and write frames without deciding authorization.

The slice moves vibit toward Nakama-class outbound realtime usefulness without broad chat, stream subscriptions, offline inboxes, delivery guarantees, or compatibility promises.

## Agent Reasoning Summary

The maintainer asked to continue and keep Nakama/Pitaya alignment explicit. After the gate, implementing a small application-owned runtime service is the next bounded step. It is narrow enough to verify recipient targeting and redaction while avoiding the ask-first boundaries around WebSocket writes, Protobuf shape, generated output, startup wiring, persistence, and broader social/realtime features.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  transport_protocol_app_separation: high
  identity_safety: high
  implementation_scope_control: high
  protocol_change_risk: none_in_this_step
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `runtime/internal/app/realtime/service.go` exists.
- `runtime/internal/app/realtime/service_test.go` exists.
- `runtime.first_server_push_realtime_messaging_runtime_slice` becomes the repository check rule for this slice.
- `M-143/W-0215` is completed.
- The next bounded direction is `W-0216 Confirm next alpha direction after realtime runtime slice`.
- WebSocket outbound delivery, Protobuf source, generated output, protocol routes, startup wiring, persistence, delivery guarantees, distributed fanout, broad chat/social modules, matchmaking, match runtime, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- a later protocol ADR changes outbound envelope semantics before transport delivery is implemented;
- connection registry semantics change away from server-observed connection id and epoch;
- the first concrete transport delivery slice needs a different handoff shape;
- subscription ownership is ratified and changes the stream target model;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal through a later ADR.

## Follow-Up

- Confirm the next alpha direction after the realtime runtime slice.
- Keep protocol payloads, generated output, WebSocket outbound delivery, stream subscriptions, chat semantics, offline inboxes, delivery guarantees, distributed fanout, matchmaking, match runtime, and direct compatibility behind later bounded work items.
