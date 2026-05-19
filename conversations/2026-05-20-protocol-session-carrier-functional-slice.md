# Conversation: Protocol Session Carrier Functional Slice

Date: 2026-05-20
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0084-protocol-session-carrier-functional-slice.md`
Related changes:

- `changes/2026-05-20-define-protocol-session-carrier-functional-slice/`

Related artifacts:

- `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`
- `runtime/internal/platform/protocol/protobuf/authentication_bridge.go`
- `runtime/internal/platform/protocol/protobuf/inventory_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_request_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Context

The repository had completed `M-103/W-0175`, the reconnect connection epoch functional slice. The next ready work item was `M-104/W-0176`, a Tier 2 protocol session carrier functional slice under `ADR-0082`.

Completed prerequisites included login-created runtime sessions, runtime session validation, route policy vocabulary, first-message connection binding, active connection registry, close policy, concrete transport close handoff, and server-observed connection epoch progression.

The remaining lifecycle gap was that successful login created a runtime session id but protocol clients could not yet receive it in a bounded carrier.

## Maintainer Narrative

The maintainer asked in Chinese to continue advancing and emphasized that vibit needs to reach Nakama/Pitaya-class server capability faster:

```text
继续推进，我们要能够保证我们的服务器能较快地达到nakama pitaya的水平，否则我们永远做不完，没法上线。
```

The maintainer also asked to reduce confirmation gates. This slice follows that direction by embedding the boundary in the change spec and implementing the bounded behavior in the same work item.

## Agent Response Summary

The agent advanced `W-0176` as a Tier 2 functional slice instead of creating another pure confirmation milestone.

The implementation reused the existing `Envelope.Session` metadata for response session carriers. Successful login responses now carry the server-created runtime session id and authenticated player id. The Protobuf response builder can also derive response session metadata from already validated application identity, while metadata-only identity remains metadata-only.

## Decisions

- Complete `M-104/W-0176`.
- Accept `ADR-0084`.
- Add `runtime.protocol_session_carrier_functional_slice` as the repository check rule for the implementation state.
- Reuse the existing envelope session metadata rather than adding Protobuf fields or generated output.
- Keep reconnect tokens, resume behavior, handshake authentication, session id as proof, logout-triggered close, runtime session revocation, presence lifecycle, operations/admin disconnect, dependencies, and direct Nakama/Pitaya API compatibility deferred.
- Move the work queue to `M-105/W-0177`, a Tier 2 presence lifecycle functional slice.

## Artifacts

- `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`
- `runtime/internal/platform/protocol/protobuf/authentication_bridge.go`
- `runtime/internal/platform/protocol/protobuf/inventory_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_request_test.go`
- `changes/2026-05-20-define-protocol-session-carrier-functional-slice/`
- `decisions/ADR-0084-protocol-session-carrier-functional-slice.md`
- `conversations/2026-05-20-protocol-session-carrier-functional-slice.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- Presence lifecycle remains next.
- Session id remains metadata, not proof.
- Reconnect tokens and resume tokens remain deferred.
- WebSocket handshake authentication remains deferred.
- Logout-triggered socket close remains deferred.
- Runtime session revocation remains deferred.
- Direct Nakama/Pitaya API compatibility remains deferred.

## Follow-Up

Continue with `M-105/W-0177 define_presence_lifecycle_functional_slice` as the next Tier 2 lifecycle/product-parity step, unless the maintainer redirects the queue.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, database credentials, close reason text, remote addresses, headers, or GitHub tokens are recorded in this conversation log.
