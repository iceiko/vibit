# Conversation: Currency Wallet Runtime Behavior Gate

Date: 2026-06-07
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-currency-wallet-runtime-behavior-gate/`

Related artifacts:

- `docs/currency-wallet-runtime-behavior-gate.md`
- `docs/currency-wallet-runtime-behavior-gate.zh-CN.md`
- `decisions/ADR-0207-currency-wallet-runtime-behavior-gate.md`
- `runtime/internal/modules/currency/repository.go`
- `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`
- `changes/2026-05-26-define-currency-wallet-runtime-behavior-gate/`
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

`M-226/W-0298` completed the currency wallet PostgreSQL adapter implementation. It added platform persistence code for `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions`, plus fake-executor tests and `UnitOfWork.NewCurrencyWalletRepository`.

The next-ready work item was `W-0299 Define currency wallet runtime behavior gate`.

## Maintainer Narrative

The maintainer asked:

```text
继续20步
```

The standing product direction was Nakama-class game/backend capability coverage adapted to vibit's agent-native maintainability model, with Pitaya retained as a distributed-runtime reference and direct public API compatibility deferred.

## Agent Response Summary

The agent advanced one bounded work item and defined the currency wallet runtime behavior gate.

The work added:

- `docs/currency-wallet-runtime-behavior-gate.md`;
- `docs/currency-wallet-runtime-behavior-gate.zh-CN.md`;
- `ADR-0207`;
- `runtime.currency_wallet_runtime_behavior_gate`;
- change spec and conversation artifacts;
- manifest and guide updates that move continuation to `W-0300 Implement currency wallet runtime behavior`;
- repository checks that keep W-0299 gate-only.

## Decisions

- Complete `M-227/W-0299`.
- Accept `ADR-0207`.
- Add `runtime.currency_wallet_runtime_behavior_gate`.
- Keep currency wallet runtime behavior implementation out of this slice.
- Select `M-228/W-0300 Implement currency wallet runtime behavior` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: durable wallets, balances, grants, spends, and transaction history need application-owned behavior before protocol routes or economy integrations become useful.

Pitaya guided the layering pressure: persistence, route/session context, handlers, and distributed runtime behavior must stay separated.

vibit adapted those lessons into its own model: a gate-only application behavior boundary with explicit identity, route-policy, idempotency, redaction, and unit-of-work handoff expectations.

## Artifacts

- `docs/currency-wallet-runtime-behavior-gate.md`
- `docs/currency-wallet-runtime-behavior-gate.zh-CN.md`
- `decisions/ADR-0207-currency-wallet-runtime-behavior-gate.md`
- `changes/2026-05-26-define-currency-wallet-runtime-behavior-gate/`
- `conversations/2026-05-26-currency-wallet-runtime-behavior-gate.md`
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

- Currency wallet runtime behavior implementation remains deferred to `W-0300`.
- Protocol routes and Protobuf messages remain deferred.
- Reward, inventory, purchase, catalog, event/audit, payment, reservation, settlement, refund, and transfer behavior remain deferred.
- Authentication/session behavior and route protection changes remain deferred.
- Hosted surfaces, SDK publication, release artifacts, distributed runtime, and direct compatibility remain deferred.

## Follow-Up

- Implement currency wallet runtime behavior in a separate bounded slice.
- Only after runtime behavior is ratified, define protocol routes and generated output.
- Keep reward, inventory, purchase, catalog, event/audit, payment, reservation, settlement, refund, and transfer behavior behind later gates.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, raw wallet material, payment provider payloads, or real private wallet state are recorded in this conversation log.
