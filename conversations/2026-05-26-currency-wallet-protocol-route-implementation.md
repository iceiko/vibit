# Conversation: Currency Wallet Protocol Route Implementation

Date: 2026-05-26
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-currency-wallet-protocol-route/`

Related artifacts:

- `proto/vibit/currency/v1/currency.proto`
- `runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go`
- `runtime/internal/app/currency/routes.go`
- `runtime/internal/app/currency/routes_test.go`
- `runtime/internal/app/bootstrap/currency.go`
- `runtime/internal/app/bootstrap/currency_test.go`
- `runtime/internal/platform/protocol/protobuf/currency_bridge.go`
- `runtime/internal/platform/protocol/protobuf/currency_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`
- `runtime/internal/platform/protocol/protobuf/payload_registry.go`
- `runtime/internal/platform/protocol/protobuf/frame_handler.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `decisions/ADR-0210-currency-wallet-protocol-route-implementation.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/currency/module.yaml`
- `modules/currency/AGENTS.md`
- `modules/currency/AGENTS.zh-CN.md`

## Context

`M-229/W-0301` completed the currency wallet protocol route gate. It authorized a bounded route implementation for the currency wallet route family while preserving reward integration, inventory integration, purchase behavior, catalog/event/payment behavior, migrations, dependencies, startup wiring, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, and direct Nakama/Pitaya API compatibility as deferred concerns.

The next-ready work item was `W-0302 Implement currency wallet protocol route`.

## Maintainer Narrative

The maintainer asked:

```text
继续20步
```

English summary: continue up to 20 next-ready steps, stopping early at blockers, ask-first boundaries, verification failure, or no next-ready work item.

## Agent Response Summary

The agent advanced the next bounded work item and implemented the protected currency wallet protocol route family.

The work added:

- `vibit.currency.v1` Protobuf request and response messages;
- generated Go Protobuf output through Buf;
- route keys for the six currency wallet routes;
- protocol bridge mapping between Protobuf payloads and application currency payloads;
- bootstrap route handlers that inject validated request identity and map redacted currency errors;
- a server-side grant policy requirement before `GrantCurrency` route execution;
- focused route, bridge, bootstrap, and authenticated local gameplay tests;
- a narrow authenticated wrapper fix for zero-field proto3 request messages;
- ADR, change spec, manifest, check-rule, and continuation updates.

## Decisions

- Complete `M-230/W-0302`.
- Accept `ADR-0210`.
- Add `runtime.currency_wallet_protocol_route_implementation`.
- Keep owner identity out of Protobuf payloads and derive it from validated request identity.
- Keep metadata-only `player_id` and `session_id` as non-proof.
- Require server-side grant policy before grant route execution.
- Allow zero-byte authenticated inner payload only when the inner Protobuf descriptor has zero fields.
- Keep reward, inventory, purchase, catalog, event/audit, payment, reservation, settlement, refund, transfer, operations/admin, SDK, hosted deployment, release artifact, distributed runtime, and direct compatibility scope deferred.
- Record that no repository next work item is currently ready after W-0302.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: durable wallets, balances, grants, spends, and transaction history are common game backend economy primitives.

Pitaya guided the layering pressure: route registration, authenticated request/session context, protocol serialization, application handlers, backend service behavior, and persistence remain separated.

vibit adapts those lessons into its own WebSocket/Protobuf route model and application-owned service boundary. This slice does not add direct Nakama or Pitaya public API compatibility.

## Artifacts

- `proto/vibit/currency/v1/currency.proto`
- `runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go`
- `runtime/internal/app/currency/routes.go`
- `runtime/internal/app/currency/routes_test.go`
- `runtime/internal/app/bootstrap/currency.go`
- `runtime/internal/app/bootstrap/currency_test.go`
- `runtime/internal/platform/protocol/protobuf/currency_bridge.go`
- `runtime/internal/platform/protocol/protobuf/currency_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`
- `runtime/internal/platform/protocol/protobuf/payload_registry.go`
- `runtime/internal/platform/protocol/protobuf/frame_handler.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `decisions/ADR-0210-currency-wallet-protocol-route-implementation.md`
- `changes/2026-05-26-implement-currency-wallet-protocol-route/`
- `conversations/2026-05-26-currency-wallet-protocol-route-implementation.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/currency/module.yaml`
- `modules/currency/AGENTS.md`
- `modules/currency/AGENTS.zh-CN.md`

## Open Questions

- No `W-0303` is currently defined. A later continuation needs a bounded maintainer direction or newly declared next-ready work item.
- Startup wiring for currency routes remains deferred.
- Reward integration, inventory integration, purchase behavior, catalog/event/payment behavior, transfers, reservations, settlement, refunds, operations/admin tools, SDKs, hosted deployment, and direct compatibility remain deferred.

## Follow-Up

- Do not invent a new next-ready work item without a bounded change spec and maintainer direction.
- Preserve the current currency route Protobuf shape unless a later proof reveals a concrete issue.
- Keep Nakama/Pitaya alignment explicit as capability and layering guidance, not direct API compatibility.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, DSNs with credentials, or raw private wallet data from a real user are recorded in this conversation log.
