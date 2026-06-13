# ADR-0206: Currency Wallet PostgreSQL Adapter Implementation

Status: Accepted
Date: 2026-06-07
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-currency-wallet-postgresql-adapter/`

Related conversations:

- `conversations/2026-05-26-currency-wallet-postgresql-adapter-implementation.md`

Related artifacts:

- `docs/currency-wallet-postgresql-adapter-gate.md`
- `decisions/ADR-0205-currency-wallet-postgresql-adapter-gate.md`
- `runtime/internal/modules/currency/repository.go`
- `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`
- `runtime/internal/platform/persistence/postgres/currency_wallet_repository_test.go`
- `runtime/internal/platform/persistence/postgres/runner.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/currency/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-225/W-0297` accepted the currency wallet PostgreSQL adapter gate in `ADR-0205`. The gate authorized the next bounded implementation slice to map the storage-neutral `runtime/internal/modules/currency.Repository` interface to the existing PostgreSQL currency wallet tables while preserving runtime, protocol, generated output, dependency, migration, authentication/session, reward, inventory, purchase, catalog/event, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

The migration source already defines `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions`. The currency module already owns wallet, owner, balance, transaction, actor, lifecycle, idempotency, amount, version, conflict, normalizer, and redacted repository-error vocabulary. The PostgreSQL platform package already owns `pgx`, executor handoff, fake-executor tests, and `UnitOfWork` repository factory helpers for other modules.

Nakama and Hiro remain product references for durable economy and currency wallet capability. Pitaya remains a deferred future distributed architecture reference. This implementation adapts the capability through vibit's bounded persistence adapter model, not direct public API compatibility.

## Decision

Implement the currency wallet PostgreSQL adapter under `runtime/internal/platform/persistence/postgres`.

The implementation adds:

- `CurrencyWalletRepository`;
- `NewCurrencyWalletRepositoryForUnitOfWork`;
- compile-time conformance to `currency.Repository`;
- wallet creation against `currency_wallets`;
- wallet lookup by wallet id;
- wallet lookup by owner identity;
- wallet-scoped balance listing with deterministic ordering and bounded pagination;
- grant recording that upserts a balance and inserts one transaction fact through the caller-supplied executor;
- spend recording that checks active wallet and sufficient balance before inserting one transaction fact;
- wallet-scoped transaction listing with optional currency filter and `(created_at, transaction_id)` cursor posture;
- row scanning through currency module normalizers;
- redacted PostgreSQL error mapping for missing rows, uniqueness, foreign-key/check/not-null failures, invalid rows, insufficient balance, and duplicate idempotency keys;
- `UnitOfWork.NewCurrencyWalletRepository`;
- focused fake-executor tests that do not require live PostgreSQL by default.

This ADR also registers `runtime.currency_wallet_postgresql_adapter_implementation` as the check rule for the slice and opens:

```text
M-227/W-0299 Define currency wallet runtime behavior gate
```

## Non-Goals

This implementation does not add:

- runtime currency wallet services, handlers, or route dispatch;
- grant/spend execution above the adapter;
- protocol routes;
- Protobuf sources;
- generated output;
- dependencies;
- migrations or schema changes;
- authentication/session behavior changes;
- reward, inventory, purchase, payment, catalog, event/audit, reservation, settlement, refund, or transfer behavior;
- SDK publication;
- hosted deployments;
- release artifacts;
- public announcements or paid promotion;
- Pitaya-style distributed architecture;
- direct Nakama, Hiro, or Pitaya API compatibility.

## Alternatives Considered

- Defer the adapter and define another gate first.
- Implement runtime wallet behavior together with the adapter.
- Implement protocol messages or routes together with the adapter.
- Put SQL under `runtime/internal/modules/currency`.
- Make the adapter own its own transaction or connection pool.
- Require live PostgreSQL verification for default tests.
- Add catalog, reward, inventory, purchase, or event/audit behavior in the same slice.
- Copy external Nakama, Hiro, or Pitaya public APIs.

## Rationale

The repository interface and migration source are already in place, and the accepted gate narrowed the implementation surface to the PostgreSQL adapter. Keeping the adapter in the platform persistence package preserves domain storage neutrality and keeps `pgx` ownership out of module code.

Fake-executor and query-capture tests are sufficient for the default verification posture because this slice verifies SQL mapping, executor handoff, result scanning, conflict mapping, redaction, and absence of transaction-control SQL without requiring a live database. Live PostgreSQL verification remains opt-in and can be added later if the repository ratifies a disposable database check path for this capability.

Opening a runtime behavior gate next preserves the existing progression used by storage objects and friends relationships: schema and repository first, PostgreSQL adapter second, runtime behavior only after an explicit gate.

## Agent Reasoning Summary

The bounded implementation that best follows `ADR-0205` is a platform-only adapter with focused tests and no behavior above the repository interface. The next useful work item is not a route or generated protocol change; it is a runtime behavior gate that can decide how validated identity, service-authoritative actor context, application unit-of-work composition, and public route behavior will call the adapter later.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  hiro_economy_reference_alignment: medium
  ai_native_requirement_test_workflow_alignment: high
  boundary_clarity: high
  implementation_risk: medium
  runtime_behavior_risk: deferred
  protocol_risk: deferred
  integration_risk: deferred
  dependency_expansion: none
confidence: high
```

## Consequences

- `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go` exists and implements `currency.Repository`.
- `runtime/internal/platform/persistence/postgres/currency_wallet_repository_test.go` exists and covers fake-executor SQL mapping, scanning, error mapping, redaction, and unit-of-work handoff.
- `runtime/internal/platform/persistence/postgres/runner.go` exposes `UnitOfWork.NewCurrencyWalletRepository`.
- `runtime.currency_wallet_postgresql_adapter_implementation` becomes the repository check rule for this slice.
- `M-226/W-0298` is completed after verification.
- `M-227/W-0299 Define currency wallet runtime behavior gate` becomes the next-ready work item after verification.
- Runtime behavior, protocol behavior, generated output, migrations, dependencies, integrations, hosted surfaces, release artifacts, distributed runtime, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- the storage-neutral `currency.Repository` interface changes materially;
- the `currency_wallets`, `currency_wallet_balances`, or `currency_wallet_transactions` schema changes materially;
- a future live PostgreSQL verification reveals SQL semantics that fake-executor tests could not catch;
- the project chooses a non-PostgreSQL first durable wallet store;
- runtime wallet behavior requires a different repository transaction handoff;
- direct Nakama, Hiro, or Pitaya compatibility becomes an explicit project goal.

## Follow-Up

- Define the currency wallet runtime behavior gate in `W-0299`.
- Keep protocol routes, generated output, reward/inventory/purchase integration, catalog/event/audit behavior, paid currency, settlement, refund, transfer, SDKs, hosted surfaces, distributed runtime, and direct compatibility behind future work items.
