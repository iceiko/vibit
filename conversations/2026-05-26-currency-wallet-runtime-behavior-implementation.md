# Conversation: Currency Wallet Runtime Behavior Implementation

Date: 2026-06-07
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-currency-wallet-runtime-behavior/`

Related artifacts:

- `runtime/internal/app/currency/service.go`
- `runtime/internal/app/currency/service_test.go`
- `decisions/ADR-0208-currency-wallet-runtime-behavior-implementation.md`
- `docs/currency-wallet-runtime-behavior-gate.md`
- `runtime/internal/modules/currency/repository.go`
- `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`
- `changes/2026-05-26-implement-currency-wallet-runtime-behavior/`
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

`M-227/W-0299` completed the currency wallet runtime behavior gate. The next-ready work item was `W-0300 Implement currency wallet runtime behavior`.

## Maintainer Narrative

The maintainer asked:

```text
继续20步
```

The standing product direction was Nakama-class game/backend capability coverage adapted to vibit's agent-native maintainability model, with Pitaya retained as a distributed-runtime reference and direct public API compatibility deferred.

## Agent Response Summary

The agent advanced one bounded implementation work item and implemented the application-owned currency wallet runtime behavior service.

The work added:

- `runtime/internal/app/currency/service.go`;
- `runtime/internal/app/currency/service_test.go`;
- `ADR-0208`;
- `runtime.currency_wallet_runtime_behavior_implementation`;
- change spec and conversation artifacts;
- manifest and guide updates that move continuation to `W-0301 Define currency wallet protocol route gate`;
- repository checks that keep W-0300 application-service-only.

## Decisions

- Complete `M-228/W-0300`.
- Accept `ADR-0208`.
- Add `runtime.currency_wallet_runtime_behavior_implementation`.
- Keep protocol routes, Protobuf sources, generated output, startup wiring, dependency, migration, reward, inventory, purchase, catalog/event, payment, hosted, SDK, distributed runtime, and direct compatibility scope out of this slice.
- Select `M-229/W-0301 Define currency wallet protocol route gate` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: durable wallets, balances, grants, spends, and transaction history need application-owned behavior before protocol routes or economy integrations become useful.

Pitaya guided the layering pressure: route/session context, handlers, persistence, and distributed runtime behavior must remain separate.

vibit adapted those lessons into its own model: a tested application service with explicit identity, route-policy, idempotency, redaction, and unit-of-work handoff expectations.

## Artifacts

- `runtime/internal/app/currency/service.go`
- `runtime/internal/app/currency/service_test.go`
- `decisions/ADR-0208-currency-wallet-runtime-behavior-implementation.md`
- `changes/2026-05-26-implement-currency-wallet-runtime-behavior/`
- `conversations/2026-05-26-currency-wallet-runtime-behavior-implementation.md`
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

- Currency wallet protocol routes and Protobuf messages remain deferred to `W-0301` and later work.
- Reward, inventory, purchase, catalog, event/audit, payment, reservation, settlement, refund, and transfer behavior remain deferred.
- Authentication/session behavior and route protection changes remain deferred.
- Hosted surfaces, SDK publication, release artifacts, distributed runtime, and direct compatibility remain deferred.

## Follow-Up

- Define the currency wallet protocol route gate in a separate bounded slice.
- Only after the protocol gate is ratified, add Protobuf routes and generated output.
- Keep reward, inventory, purchase, catalog, event/audit, payment, reservation, settlement, refund, and transfer behavior behind later gates.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, raw wallet material, payment provider payloads, or real private wallet state are recorded in this conversation log.
