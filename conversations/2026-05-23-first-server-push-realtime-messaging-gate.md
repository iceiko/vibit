# Conversation: First Server Push And Realtime Messaging Gate

Date: 2026-05-23
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-23-define-first-server-push-realtime-messaging-gate/`

Related artifacts:

- `docs/first-server-push-realtime-messaging-gate.md`
- `docs/first-server-push-realtime-messaging-gate.zh-CN.md`
- `decisions/ADR-0122-first-server-push-realtime-messaging-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-141/W-0213` selected the first server push and realtime messaging gate as the next prototype-ready direction after the storage objects local proof. The current next-ready item was `W-0214 Define first server push and realtime messaging gate`.

The requested continuation again required keeping Nakama/Pitaya alignment explicit and then committing and pushing the result.

## Maintainer Narrative

The maintainer asked:

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。
```

English summary: continue the work, keep Nakama/Pitaya alignment explicit, then commit and push the result.

## Agent Response Summary

The agent advanced one bounded gate work item. The completed gate is:

```text
W-0214 Define first server push and realtime messaging gate
```

The change accepts `ADR-0122`, adds `docs/first-server-push-realtime-messaging-gate.md`, registers `runtime.first_server_push_realtime_messaging_gate`, and opens `M-143/W-0215 Implement first server push and realtime messaging runtime slice` as next-ready.

## Decisions

- Complete `M-142/W-0214`.
- Accept `ADR-0122`.
- Define the first outbound realtime vocabulary and ownership boundary.
- Keep server push and realtime implementation behind `W-0215`.
- Keep protocol routes, Protobuf source, generated output, persistence, delivery guarantees, distributed fanout, matchmaking, match runtime, broad social modules, and direct compatibility deferred.

## Nakama And Pitaya Reference Basis

Nakama guided the capability vocabulary: notifications, streams, chat, and presence-adjacent outbound messages are common game-backend surfaces after durable storage exists.

Pitaya guided the layering vocabulary: acceptors, sessions, handlers, push, groups, broadcast, backend services, and later cluster/RPC concerns must remain separated.

This gate adapts those references into vibit-owned terms: `server_notice`, `domain_event_push`, `stream_message`, `presence_signal`, `connection_id_and_epoch`, `player_current_connections`, and `stream_subscribers`.

Direct compatibility remains deferred.

## Artifacts

- `docs/first-server-push-realtime-messaging-gate.md`
- `docs/first-server-push-realtime-messaging-gate.zh-CN.md`
- `decisions/ADR-0122-first-server-push-realtime-messaging-gate.md`
- `changes/2026-05-23-define-first-server-push-realtime-messaging-gate/`
- `conversations/2026-05-23-first-server-push-realtime-messaging-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- `W-0215` must choose the smallest runtime behavior shape and tests for first outbound realtime delivery.
- Protocol payloads, generated output, transport delivery, persistence, acknowledgements, ordering, retries, backpressure, stream subscriptions, chat semantics, and distributed fanout remain behind later explicit authorization unless `W-0215` narrows them further.

## Follow-Up

- Complete `W-0215`: implement the smallest first server push and realtime messaging runtime slice authorized by `ADR-0122`.
- Keep direct Nakama/Pitaya API compatibility, broad chat/social modules, matchmaking, match runtime, protocol expansion, generated output, persistence, and distributed fanout behind later bounded work items unless explicitly authorized.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, DSNs with credentials, raw storage object values, or concrete transport metadata from a real user are recorded in this conversation log.
