# Conversation: Currency Wallet Protocol Route Gate

Date: 2026-06-08
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-currency-wallet-protocol-route-gate/`

Related artifacts:

- `docs/currency-wallet-protocol-route-gate.md`
- `docs/currency-wallet-protocol-route-gate.zh-CN.md`
- `decisions/ADR-0209-currency-wallet-protocol-route-gate.md`
- `changes/2026-05-26-define-currency-wallet-protocol-route-gate/`
- `conversations/2026-05-26-currency-wallet-protocol-route-gate.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/currency/module.yaml`

## Context

`M-228/W-0300` implemented application-owned currency wallet runtime behavior under `runtime/internal/app/currency`. The service covered own-player wallet ensure/get, balance listing, server-authoritative grants, player-authorized spends, and transaction listing with validated request identity handoff, metadata-only identity refusal, unit-of-work repository access, and redacted public errors.

The next-ready work item was `W-0301 Define currency wallet protocol route gate`.

## Maintainer Narrative

The maintainer asked to continue up to twenty next-ready steps.

```text
继续20步
```

The standing product direction remained Nakama-class game/backend capability coverage adapted to vibit's agent-native maintainability model, with Pitaya retained as a distributed-runtime reference and direct public API compatibility deferred.

## Agent Response Summary

The agent advanced one bounded work item and defined a protocol route gate only. The slice recorded the future protected currency wallet route family, candidate Protobuf package and generated-output posture, request identity handoff, public error mapping, service handoff, local proof expectations, Nakama/Pitaya reference mapping, verification expectations, and stop conditions.

The future route family is:

- `currency.EnsurePlayerWallet`
- `currency.GetOwnWallet`
- `currency.ListOwnWalletBalances`
- `currency.GrantCurrency`
- `currency.SpendCurrency`
- `currency.ListOwnWalletTransactions`

The gate keeps wallet owner identity derived only from validated authenticated request identity. Metadata-only `player_id` and `session_id` remain insufficient proof, and client payloads must not provide owner ids, wallet ids for lookup, or actor ids as proof.

## Decisions

- Complete `M-229/W-0301`.
- Accept `ADR-0209`.
- Add `runtime.currency_wallet_protocol_route_gate`.
- Record `docs/currency-wallet-protocol-route-gate.md` and `docs/currency-wallet-protocol-route-gate.zh-CN.md`.
- Preserve route implementation, Protobuf source, generated output, startup wiring, repository interface changes, adapter changes, dependencies, migrations, authentication/session behavior changes, reward/inventory/purchase behavior, catalog/event/payment behavior, hosted surfaces, SDKs, distributed runtime, and direct Nakama/Pitaya API compatibility as deferred concerns.
- Select `M-230/W-0302 Implement currency wallet protocol route` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the route family pressure: durable wallets, balance reads, grants, spends, and transaction history are core game economy access patterns.

Pitaya guided the layering pressure: route handlers, request/session identity, protocol adapter mapping, application behavior, and persistence must remain separated.

vibit adapted those references into protected request-token route planning with validated authenticated request identity, no metadata-only proof, no client-supplied owner/wallet/actor proof, and no direct public API compatibility.

## Artifacts

- `docs/currency-wallet-protocol-route-gate.md`
- `docs/currency-wallet-protocol-route-gate.zh-CN.md`
- `decisions/ADR-0209-currency-wallet-protocol-route-gate.md`
- `changes/2026-05-26-define-currency-wallet-protocol-route-gate/`
- `conversations/2026-05-26-currency-wallet-protocol-route-gate.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/currency/module.yaml`

## Open Questions

- Protobuf request and response field details remain deferred to `W-0302`.
- Protocol bridge implementation and bootstrap route registration remain deferred.
- Route-specific generated output remains deferred.
- Grant route exposure must prove it cannot become client-authoritative minting.
- Reward, inventory, purchase, catalog, event/audit, payment, reservation, settlement, refund, transfer, SDK publication, hosted surfaces, distributed runtime, and direct compatibility remain deferred.

## Follow-Up

The next-ready work item is:

```text
M-230/W-0302 Implement currency wallet protocol route
```

That next slice may implement only the protected currency wallet route family authorized by `ADR-0209`.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, raw wallet material, payment provider payloads, or real private wallet state are recorded in this conversation log.
