# Conversation: First Server Push And Realtime Messaging Runtime Slice

Date: 2026-05-23
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-23-implement-first-server-push-realtime-messaging-runtime-slice/`

Related artifacts:

- `runtime/internal/app/realtime/service.go`
- `runtime/internal/app/realtime/service_test.go`
- `decisions/ADR-0123-first-server-push-realtime-messaging-runtime-slice.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-142/W-0214` defined the first server push and realtime messaging gate. The current next-ready item was `W-0215 Implement first server push and realtime messaging runtime slice`.

The requested continuation required keeping Nakama/Pitaya alignment explicit and then committing and pushing the result.

## Maintainer Narrative

The maintainer asked:

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。
```

English summary: continue the work, keep Nakama/Pitaya alignment explicit, then commit and push the result.

## Agent Response Summary

The agent advanced one bounded runtime work item:

```text
W-0215 Implement first server push and realtime messaging runtime slice
```

The change adds `runtime/internal/app/realtime/service.go` and focused tests. The service accepts server-authored realtime message intents, validates allowed intent and target vocabulary, requires validated service/admin sender authority, rejects metadata-only and player-authored publish attempts before registry resolution, resolves active bound recipients from the single-process connection registry, and returns redacted delivery intents without writing sockets.

## Decisions

- Complete `M-143/W-0215`.
- Accept `ADR-0123`.
- Add application-owned realtime service behavior under `runtime/internal/app/realtime`.
- Register `runtime.first_server_push_realtime_messaging_runtime_slice`.
- Open `M-144/W-0216 Confirm next alpha direction after realtime runtime slice` as next-ready.
- Keep WebSocket outbound delivery, Protobuf source, generated output, protocol routes, startup wiring, persistence, delivery guarantees, distributed fanout, matchmaking, match runtime, broad social modules, and direct compatibility deferred.

## Nakama And Pitaya Reference Basis

Nakama guided the capability family: notifications, streams, chat, and presence-adjacent outbound messages are common game-backend needs after durable storage and connection identity exist.

Pitaya guided the layering: acceptors, sessions, handlers, backend services, push, groups, broadcast, and cluster/RPC concerns must stay separated.

This slice adapts those references by placing recipient policy and delivery-intent construction in application code, using server-observed connection registry state, and leaving protocol and WebSocket delivery adapters for later bounded work.

Direct Nakama/Pitaya API compatibility remains deferred.

## Artifacts

- `runtime/internal/app/realtime/service.go`
- `runtime/internal/app/realtime/service_test.go`
- `decisions/ADR-0123-first-server-push-realtime-messaging-runtime-slice.md`
- `changes/2026-05-23-implement-first-server-push-realtime-messaging-runtime-slice/`
- `conversations/2026-05-23-first-server-push-realtime-messaging-runtime-slice.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- The next direction after W-0215 must be confirmed before adding protocol payloads, generated output, WebSocket outbound delivery, operations inspection, failure/concurrency verification, or another prototype-ready family.
- Stream subscription ownership remains undefined; `stream_subscribers` is still future vocabulary only.
- Concrete delivery guarantees, offline inboxes, ordering, retries, backpressure, distributed fanout, and chat semantics remain separate future decisions.

## Follow-Up

- Complete `W-0216`: confirm the next alpha direction after the realtime runtime slice.
- Keep direct Nakama/Pitaya API compatibility, broad chat/social modules, matchmaking, match runtime, protocol expansion, generated output, persistence, and distributed fanout behind later bounded work items unless explicitly authorized.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, DSNs with credentials, raw storage object values, raw realtime payload from a real user, or concrete transport metadata from a real user are recorded in this conversation log.
