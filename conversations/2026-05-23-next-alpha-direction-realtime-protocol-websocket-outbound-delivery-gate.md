# Conversation: Next Alpha Direction Realtime Protocol WebSocket Outbound Delivery Gate

Date: 2026-05-23
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-23-confirm-next-alpha-direction-after-realtime-runtime-slice/`

Related artifacts:

- `decisions/ADR-0124-next-alpha-direction-realtime-protocol-websocket-outbound-delivery-gate.md`
- `runtime/internal/app/realtime/service.go`
- `runtime/internal/app/realtime/service_test.go`
- `docs/first-server-push-realtime-messaging-gate.md`
- `docs/product-maturity-milestones.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/prototype-ready-foundation-execution-plan.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-143/W-0215` completed the first server push and realtime messaging runtime slice. The next-ready item was `W-0216 Confirm next alpha direction after realtime runtime slice`.

The requested continuation needed to keep Nakama/Pitaya alignment explicit, then commit and push the result.

## Maintainer Narrative

The maintainer asked:

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。提交所用的Key在Git忽略的文件里有
```

English summary: continue the work, keep Nakama/Pitaya alignment explicit, commit and push the result, and use the Git credential material from an ignored local file.

## Agent Response Summary

The agent advanced one bounded direction-selection work item. The selected next direction is:

```text
define_realtime_protocol_websocket_outbound_delivery_gate
```

The change completes `M-144/W-0216`, accepts `ADR-0124`, adds `runtime.next_alpha_direction_after_realtime_runtime_slice`, and opens `M-145/W-0217 Define realtime protocol and WebSocket outbound delivery gate` as the next-ready item.

## Decisions

- Complete `M-144/W-0216`.
- Accept `ADR-0124`.
- Select the realtime protocol and WebSocket outbound delivery gate as the next prototype-ready alpha direction.
- Open `M-145/W-0217`.
- Keep the next step as a gate before implementation.
- Keep WebSocket outbound delivery, socket writes, Protobuf source, generated output, protocol bridge, application bootstrap handler, startup wiring, persistence, migrations, dependencies, authentication/session changes, route-protection changes, stream subscriptions, chat, groups, broadcast fanout, delivery guarantees, distributed runtime, matchmaking, match runtime, hosted deployments, release artifacts, public announcements, paid promotion, blob/S3 storage, and direct compatibility deferred.

## Nakama And Pitaya Reference Basis

Nakama guided the capability choice: after server-authored delivery intents exist in the application layer, the next useful prototype-ready step is to define bounded protocol and outbound transport delivery planning for future notifications, streams, chat, and presence-adjacent outbound behavior.

Pitaya guided the layering choice: acceptor, session, handler, protocol serialization, backend service intent, push, group, broadcast, and later cluster/RPC concerns must stay separated before any concrete socket writes, group fanout, service discovery, or remote call behavior.

Direct compatibility remains deferred.

## Artifacts

- `decisions/ADR-0124-next-alpha-direction-realtime-protocol-websocket-outbound-delivery-gate.md`
- `changes/2026-05-23-confirm-next-alpha-direction-after-realtime-runtime-slice/`
- `conversations/2026-05-23-next-alpha-direction-realtime-protocol-websocket-outbound-delivery-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- The future gate must decide the first outbound protocol payload posture, if any, without adding generated output in this slice.
- The future gate must decide the planned transport delivery adapter boundary before socket writes exist.
- Persistence, offline inboxes, acknowledgements, ordering, retries, backpressure, distributed fanout, group/room scopes, chat semantics, and client SDK behavior remain deferred.

## Follow-Up

- Advance `W-0217 Define realtime protocol and WebSocket outbound delivery gate`.
- Preserve Nakama/Pitaya alignment as capability and layering guidance, not direct API compatibility.

## Redaction Notes

No secrets, GitHub tokens, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, DSNs with credentials, or raw storage object values from a real user are recorded in this conversation log.
