# Conversation: Currency Wallet PostgreSQL Adapter Implementation

Date: 2026-06-07
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-currency-wallet-postgresql-adapter/`

Related artifacts:

- `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`
- `runtime/internal/platform/persistence/postgres/currency_wallet_repository_test.go`
- `runtime/internal/platform/persistence/postgres/runner.go`
- `decisions/ADR-0206-currency-wallet-postgresql-adapter-implementation.md`
- `changes/2026-05-26-implement-currency-wallet-postgresql-adapter/`
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

`M-225/W-0297` completed the currency wallet PostgreSQL adapter gate. It authorized a later bounded implementation under `runtime/internal/platform/persistence/postgres` and preserved runtime wallet behavior, protocol routes, Protobuf sources, generated output, dependencies, migrations, authentication/session changes, reward integration, inventory integration, purchase behavior, catalog/event tables, hosted surfaces, SDKs, release artifacts, distributed runtime, and direct Nakama/Pitaya API compatibility as deferred concerns.

The next-ready work item was `W-0298 Implement currency wallet PostgreSQL adapter`.

## Maintainer Narrative

The maintainer asked:

```text
继续20步
```

The standing product direction was Nakama-class backend capability coverage adapted to vibit's agent-native maintainability model, with Pitaya as a distributed-runtime reference and direct API compatibility deferred.

## Agent Response Summary

The agent advanced one bounded work item and implemented the currency wallet PostgreSQL adapter under `runtime/internal/platform/persistence/postgres` using the existing `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions` migration source.

The work added:

- `CurrencyWalletRepository`;
- `NewCurrencyWalletRepositoryForUnitOfWork`;
- `UnitOfWork.NewCurrencyWalletRepository`;
- SQL mapping for wallet create, wallet lookup, owner lookup, balance listing, grant recording, spend recording, and transaction listing;
- idempotency, duplicate owner, foreign-key, version, insufficient-balance, not-found, invalid-row, and storage-unavailable error mapping into redacted currency repository errors;
- row scanning through currency module normalizers;
- caller-owned transaction control with no `BEGIN`, `COMMIT`, or `ROLLBACK` SQL in the adapter;
- focused fake-executor adapter tests;
- ADR, change spec, manifest, check-rule, and continuation updates.

## TDD Notes

A failing test was written before implementation. The first focused run failed because `CurrencyWalletRepository`, `NewCurrencyWalletRepositoryForUnitOfWork`, and `UnitOfWork.NewCurrencyWalletRepository` did not exist. After the initial implementation, a scanning panic exposed a row-helper type mismatch, so the adapter was changed to scan database string columns into local `string` variables before converting to currency value types. A later failure showed SQL arguments carrying named string types into fake-executor assertions; the adapter now converts domain string aliases to plain `string` at the SQL boundary.

## Decisions

- Complete `M-226/W-0298`.
- Accept `ADR-0206`.
- Add `runtime.currency_wallet_postgresql_adapter_implementation`.
- Keep currency wallet runtime behavior out of this slice.
- Select `M-227/W-0299 Define currency wallet runtime behavior gate` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: durable wallet and balance state needs a concrete persistence path before useful grant, spend, balance, transaction, reward, inventory, purchase, or economy-facing runtime behavior can be defined.

Pitaya guided the layering pressure: persistence concerns should remain below handlers, routes, RPC, distributed lifecycle, and cluster behavior.

vibit adapted those lessons into its own model: a PostgreSQL adapter implementing a storage-neutral repository interface, with no direct public API compatibility and no runtime/protocol behavior in this slice.

## Artifacts

- `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`
- `runtime/internal/platform/persistence/postgres/currency_wallet_repository_test.go`
- `runtime/internal/platform/persistence/postgres/runner.go`
- `decisions/ADR-0206-currency-wallet-postgresql-adapter-implementation.md`
- `changes/2026-05-26-implement-currency-wallet-postgresql-adapter/`
- `conversations/2026-05-26-currency-wallet-postgresql-adapter-implementation.md`
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

- Currency wallet runtime behavior remains deferred to `W-0299` and later implementation work.
- Protocol routes and Protobuf messages remain deferred.
- Reward, inventory, purchase, catalog, event/audit, payment, reservation, settlement, refund, and transfer behavior remain deferred.
- Authentication/session behavior and route protection changes remain deferred.
- Hosted surfaces, SDK publication, release artifacts, distributed runtime, and direct compatibility remain deferred.
- Live PostgreSQL adapter verification remains optional and unavailable in this default fake-executor slice.

## Follow-Up

- Define the currency wallet runtime behavior gate.
- Only after that gate, implement currency wallet runtime behavior in a separate bounded slice.
- Only after runtime behavior is ratified, define protocol routes and generated output.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, raw wallet material, or real private wallet state are recorded in this conversation log.
