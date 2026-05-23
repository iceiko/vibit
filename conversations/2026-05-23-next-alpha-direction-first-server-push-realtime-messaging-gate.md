# Conversation: Next Alpha Direction First Server Push Realtime Messaging Gate

Date: 2026-05-23
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-23-confirm-next-alpha-direction-after-storage-objects-local-proof/`

Related artifacts:

- `decisions/ADR-0121-next-alpha-direction-first-server-push-realtime-messaging-gate.md`
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

`M-140/W-0212` completed the storage objects protocol route local proof. The next-ready item was `W-0213 Confirm next alpha direction after storage objects local proof`.

The requested continuation needed to keep Nakama/Pitaya alignment explicit and then commit and push the result.

## Maintainer Narrative

The maintainer asked:

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。
```

English summary: continue the work, keep Nakama/Pitaya alignment explicit, then commit and push the result.

## Agent Response Summary

The agent advanced one bounded direction-selection work item. The selected next direction is:

```text
define_first_server_push_realtime_messaging_gate
```

The change completes `M-141/W-0213`, accepts `ADR-0121`, adds `runtime.next_alpha_direction_after_storage_objects_local_proof`, and opens `M-142/W-0214 Define first server push and realtime messaging gate` as the next-ready item.

## Decisions

- Complete `M-141/W-0213`.
- Accept `ADR-0121`.
- Select the first server push / realtime messaging gate as the next prototype-ready alpha direction.
- Open `M-142/W-0214`.
- Keep the next step as a gate before implementation.
- Keep runtime behavior, protocol routes, Protobuf messages, generated output, storage behavior, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session changes, hosted deployments, release artifacts, public announcements, paid promotion, matchmaking, match runtime, distributed runtime, blob/S3 storage, and direct compatibility deferred.

## Nakama And Pitaya Reference Basis

Nakama guided the capability choice: after durable storage objects, a prototype-ready backend needs an outbound realtime vocabulary that can later support notifications, streams, chat, and presence-adjacent shared online services.

Pitaya guided the layering choice: push, broadcast, group, stream, session, handler, protocol, backend service, and later cluster/RPC concerns must remain separated and explicitly gated.

This direction does not add direct Nakama or Pitaya public API compatibility.

## Artifacts

- `decisions/ADR-0121-next-alpha-direction-first-server-push-realtime-messaging-gate.md`
- `changes/2026-05-23-confirm-next-alpha-direction-after-storage-objects-local-proof/`
- `conversations/2026-05-23-next-alpha-direction-first-server-push-realtime-messaging-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- The future gate must choose the first vocabulary shape for server push or realtime messaging.
- Persistence, delivery guarantees, offline inboxes, acknowledgements, ordering, retries, backpressure, distributed fanout, group/room scopes, chat semantics, and client SDK behavior remain deferred.

## Follow-Up

- Advance `W-0214 Define first server push and realtime messaging gate`.
- Preserve Nakama/Pitaya alignment as capability and layering guidance, not direct API compatibility.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, DSNs with credentials, or raw storage object values from a real user are recorded in this conversation log.
