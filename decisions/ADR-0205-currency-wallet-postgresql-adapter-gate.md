# ADR-0205: Currency Wallet PostgreSQL Adapter Gate

Status: Accepted
Date: 2026-06-07
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-currency-wallet-postgresql-adapter-gate/`

Related conversations:

- `conversations/2026-05-26-currency-wallet-postgresql-adapter-gate.md`

Related artifacts:

- `docs/currency-wallet-postgresql-adapter-gate.md`
- `docs/currency-wallet-postgresql-adapter-gate.zh-CN.md`
- `runtime/internal/modules/currency/repository.go`
- `runtime/migrations/postgres/000008_create_currency_wallets.sql`
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

`M-224/W-0296` implemented the storage-neutral currency wallet repository interface under `runtime/internal/modules/currency`. The next bounded step is to define the PostgreSQL adapter gate before any adapter code, SQL execution, runtime behavior, protocol surface, or integration behavior is added.

The existing `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions` migration source already defines the first durable table set, and the repository interface already defines wallet creation, wallet lookup, owner lookup, balance listing, grant recording, spend recording, transaction listing, idempotency, version, and conflict vocabulary. A separate adapter gate keeps SQL mapping, transaction handoff, affected-row handling, conflict mapping, redaction, and tests explicit before implementation.

Nakama and Hiro remain the product capability references for durable economy and currency wallet state. Pitaya remains deferred as a future distributed architecture reference. vibit adapts the capability through explicit adapter boundaries and checkable manifests, not direct public API compatibility.

## Decision

Accept `docs/currency-wallet-postgresql-adapter-gate.md` as the gate for the future currency wallet PostgreSQL adapter.

The gate records:

- future adapter owner `runtime/internal/platform/persistence/postgres`;
- future source candidate `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`;
- future test candidate `runtime/internal/platform/persistence/postgres/currency_wallet_repository_test.go`;
- constructor candidate `NewCurrencyWalletRepositoryForUnitOfWork`;
- repository interface source `runtime/internal/modules/currency.Repository`;
- SQL mapping posture for `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions`;
- constructor and caller-supplied executor expectations;
- unit-of-work and transaction handoff expectations;
- redacted driver-error, affected-row, idempotency, version, and insufficient-balance mapping;
- focused adapter implementation test expectations;
- stop conditions before implementation, runtime behavior, grant/spend execution above the adapter, protocol routes, generated output, dependencies, migration changes, authentication/session behavior, reward integration, inventory integration, purchase behavior, catalog/event tables, hosted deployment, release artifacts, public announcements, paid promotion, distributed runtime, or direct compatibility.

This ADR does not add PostgreSQL currency wallet adapters, SQL execution behavior, unit-of-work factory wiring, runtime wallet handlers, grant/spend execution, protocol routes, Protobuf sources, generated output, dependencies, migration changes, authentication/session behavior changes, reward integration, inventory integration, purchase behavior, catalog tables, event/audit tables, SDK publication, hosted deployments, release artifacts, public announcements, paid promotion, distributed runtime, Pitaya-style architecture, or direct Nakama/Pitaya API compatibility.

Open the next bounded work item:

```text
M-226/W-0298 Implement currency wallet PostgreSQL adapter
```

## Alternatives Considered

- Implement the PostgreSQL adapter immediately after the repository interface.
- Reuse the storage objects or friends adapter implementation shape without a currency-specific gate.
- Put SQL execution under `runtime/internal/modules/currency`.
- Add runtime grant/spend behavior together with adapter implementation.
- Add Protobuf messages or public routes immediately.
- Add catalog, event/audit, purchase, reward, inventory, settlement, refund, or transfer behavior in the same slice.
- Copy external Nakama, Hiro, or Pitaya public API compatibility.

## Rationale

The currency wallet repository has conflict, idempotency, and redaction pressure that should be explicit before implementation: owner uniqueness, wallet/balance version checks, grant/spend direction, idempotency uniqueness, affected-row interpretation, insufficient-balance handling, foreign-key outcomes, transaction fact insertion, balance update ordering, metadata privacy, and driver-detail redaction.

A gate-only ADR keeps the next implementation slice bounded and makes repository checks able to reject accidental SQL, protocol, runtime, generated output, integration, or hosted feature behavior before that work item is authorized.

## Agent Reasoning Summary

The safest continuation from `W-0296` is a platform adapter gate. It gives future implementation a precise owner, constructor posture, SQL mapping checklist, conflict mapping posture, and test list while preserving the separation between currency module vocabulary, PostgreSQL adapter behavior, runtime routing, protocol shape, broader economy integrations, and product-scope expansion.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  hiro_economy_reference_alignment: medium
  ai_native_requirement_test_workflow_alignment: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: low
  adapter_risk: contained_by_next_work_item
  protocol_risk: deferred
  integration_risk: deferred
  dependency_expansion: low
confidence: high
```

## Consequences

- `docs/currency-wallet-postgresql-adapter-gate.md` and `.zh-CN.md` exist.
- `runtime.currency_wallet_postgresql_adapter_gate` becomes the repository check rule for this slice.
- `M-225/W-0297` is completed.
- `M-226/W-0298 Implement currency wallet PostgreSQL adapter` becomes the next-ready work item.
- Existing runtime behavior, protocol behavior, migrations, dependencies, generated output, integrations, and SQL execution behavior remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- `runtime/internal/modules/currency.Repository` changes before adapter implementation;
- the `currency_wallets`, `currency_wallet_balances`, or `currency_wallet_transactions` migration source changes materially;
- the project selects a different first currency wallet persistence engine;
- transaction ownership moves away from caller-supplied unit-of-work boundaries;
- catalog, event/audit, purchase, reward, inventory, settlement, refund, or transfer behavior becomes mandatory before adapter work;
- direct Nakama, Hiro, or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Implement the currency wallet PostgreSQL adapter only after this gate is accepted.
- Keep runtime behavior, permissions, protocol routes, generated output, local proof, and economy integrations behind later gates.
- Keep catalog, event/audit tables, purchase, reward, inventory, payment, settlement, refunds, transfers, SDKs, hosted surfaces, distributed runtime, and direct compatibility behind future gates.
