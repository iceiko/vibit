# Conversation: Currency Wallet PostgreSQL Adapter Gate

Date: 2026-06-07
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-currency-wallet-postgresql-adapter-gate/`

Related artifacts:

- `docs/currency-wallet-postgresql-adapter-gate.md`
- `docs/currency-wallet-postgresql-adapter-gate.zh-CN.md`
- `decisions/ADR-0205-currency-wallet-postgresql-adapter-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/currency/module.yaml`
- `modules/currency/AGENTS.md`
- `modules/currency/AGENTS.zh-CN.md`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Context

`M-224/W-0296` completed the storage-neutral `runtime/internal/modules/currency.Repository` interface. That slice defined wallet, balance, transaction, actor, idempotency, version, and error vocabulary without adding PostgreSQL SQL execution, adapter implementation, runtime wallet behavior, protocol routes, generated output, dependencies, reward integration, inventory integration, purchase behavior, hosted surfaces, SDKs, distributed runtime, or direct Nakama/Pitaya API compatibility.

The next-ready work item was `M-225/W-0297 Define currency wallet PostgreSQL adapter gate`.

## Maintainer Narrative

The maintainer asked to continue multiple bounded steps:

```text
继续20步
```

The standing product direction remained Nakama/Hiro-aligned economy capability coverage, adapted to vibit's agent-native maintainability model, with Pitaya retained as a distributed-runtime reference but not as a direct compatibility target.

## Agent Response Summary

The agent defined the currency wallet PostgreSQL adapter gate after the storage-neutral repository interface existed.

The gate records future ownership under `runtime/internal/platform/persistence/postgres`, source and test candidates, the planned `NewCurrencyWalletRepositoryForUnitOfWork` constructor, unit-of-work transaction handoff, SQL mapping posture for `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions`, idempotency and version-conflict behavior, affected-row handling, insufficient-balance mapping, redaction requirements, and focused implementation test expectations.

The work also added the check rule `runtime.currency_wallet_postgresql_adapter_gate` and advanced repository pointers to `M-226/W-0298 Implement currency wallet PostgreSQL adapter`.

## Decisions

- Complete `M-225/W-0297`.
- Accept `ADR-0205`.
- Add `runtime.currency_wallet_postgresql_adapter_gate`.
- Keep the current slice gate-only and documentation/check-rule focused.
- Select `M-226/W-0298 Implement currency wallet PostgreSQL adapter` as the next-ready work item.

## Artifacts

- `docs/currency-wallet-postgresql-adapter-gate.md`
- `docs/currency-wallet-postgresql-adapter-gate.zh-CN.md`
- `decisions/ADR-0205-currency-wallet-postgresql-adapter-gate.md`
- `changes/2026-05-26-define-currency-wallet-postgresql-adapter-gate/`
- `conversations/2026-05-26-currency-wallet-postgresql-adapter-gate.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/currency/module.yaml`
- `README.md`
- `README.zh-CN.md`

## Open Questions

- The PostgreSQL adapter implementation remains deferred to `W-0298`.
- Unit-of-work factory wiring remains deferred.
- Runtime wallet creation, balance reads, grant execution, spend execution, and transaction reads remain deferred.
- Protocol routes and Protobuf messages remain deferred.
- Reward, inventory, purchase, payment, catalog, event/audit, reservation, settlement, refund, transfer, SDK, hosted surface, distributed runtime, and direct compatibility scope remain deferred.

## Follow-Up

- Implement the currency wallet PostgreSQL adapter in the next bounded slice.
- Keep adapter behavior constrained to the gate until a later runtime behavior gate authorizes application service execution.
- Preserve protocol, generated output, hosted, SDK, reward, inventory, purchase, and distributed-runtime work behind later explicit gates.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, payment provider payloads, raw wallet records from a real user, or raw private economy data are recorded in this conversation log.
