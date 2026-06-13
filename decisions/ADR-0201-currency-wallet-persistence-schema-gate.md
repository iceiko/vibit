# ADR-0201: Currency Wallet Persistence Schema Gate

Status: Accepted
Date: 2026-06-07
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-currency-wallet-persistence-schema-gate/`

Related conversations:

- `conversations/2026-05-26-currency-wallet-persistence-schema-gate.md`

Related artifacts:

- `docs/currency-wallet-persistence-schema-gate.md`
- `docs/currency-wallet-persistence-schema-gate.zh-CN.md`
- `docs/currency-wallet-lifecycle-boundary-gate.md`
- `docs/postgresql-persistence-boundary.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0200` completed the currency wallet lifecycle boundary gate. It defined future currency catalog read, wallet read, balance read, grant, spend, balance change recording, and transaction read semantics without adding runtime behavior, balance tables, wallet transaction behavior, protocol routes, generated output, migrations, repositories, adapters, dependencies, or direct compatibility scope.

The next step is data-first: future balance mutation correctness requires stable durable wallet, balance, transaction, and idempotency posture before SQL, repositories, adapters, protocol, or runtime handlers exist.

## Decision

Define the currency wallet persistence schema gate in:

```text
docs/currency-wallet-persistence-schema-gate.md
docs/currency-wallet-persistence-schema-gate.zh-CN.md
```

The first persistence target is PostgreSQL. The future migration source candidate is:

```text
runtime/migrations/postgres/000008_create_currency_wallets.sql
```

The future logical table candidates are:

```text
currency_wallets
currency_wallet_balances
currency_wallet_transactions
```

The gate records a player-owned wallet identity posture, balance row identity, transaction or ledger fact posture, idempotency uniqueness posture, catalog deferral, timestamps, redaction, and repository/adapter ownership candidates.

The repository check rule is:

```text
runtime.currency_wallet_persistence_schema_gate
```

This ADR is a schema gate only. It does not add the SQL migration source, create tables, add repositories or adapters, implement runtime behavior, expose protocol routes, add Protobuf source or generated output, add dependencies, wire startup, add reward, purchase, inventory, paid-currency, event/audit, reservation, settlement, refund, transfer, or direct compatibility behavior.

Open the next bounded work item:

```text
M-222/W-0294 Add currency wallet migration source
```

## Alternatives Considered

- Add the SQL migration source in the same change.
- Add repository interfaces before a schema gate.
- Add runtime grant/spend behavior before persistence is explicit.
- Store balances only and omit transaction records.
- Store a full currency catalog in the first migration.
- Copy an external Nakama or Pitaya API or schema compatibility surface.

## Rationale

Wallet and balance state need a transaction record to make future grant/spend behavior inspectable, idempotent, and supportable. A separate balance table keeps reads efficient, while the transaction table records durable mutation facts and duplicate-prevention keys.

The first catalog posture remains deferred because code normalization, precision, display metadata, and mutability rules need their own bounded decision before catalog tables exist. The future migration can still enforce non-blank bounded currency codes and later reference a catalog after a catalog gate authorizes it.

Using PostgreSQL follows the existing durable runtime posture and avoids introducing a document database, graph database, cache dependency, payment dependency, or external economy service before local schema correctness is proven.

## Agent Reasoning Summary

The agent continued from `W-0293`, treated it as schema-gate-only, selected PostgreSQL as the first persistence target, recorded future `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions` table candidates, anchored idempotency in transaction records, deferred catalog/event/reward/purchase/inventory integrations, and preserved all implementation deferrals.

## Decision Weights

```yaml
decision_weights:
  nakama_hiro_economy_alignment: high
  ai_native_requirement_test_workflow_alignment: high
  contract_first_safety: high
  durable_balance_correctness: high
  idempotency_and_auditability: high
  privacy_and_redaction_control: high
  future_testability: high
  migration_scope_change: none
  runtime_scope_change: none
  protocol_scope_change: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `M-221/W-0293` is completed.
- `runtime.currency_wallet_persistence_schema_gate` is registered.
- The currency wallet persistence schema standard and Simplified Chinese translation exist.
- The future migration source candidate is `runtime/migrations/postgres/000008_create_currency_wallets.sql`.
- The future table candidates are `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions`.
- `M-222/W-0294 Add currency wallet migration source` becomes next-ready.
- Migration source creation, tables, repository interfaces, PostgreSQL adapters, runtime behavior, protocol routes, Protobuf source, generated output, startup wiring, dependencies, reward integration, purchase behavior, inventory integration, event/audit tables, hosted surfaces, SDKs, distributed runtime, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- the migration-source slice shows the table candidates cannot preserve wallet, balance, and transaction consistency;
- idempotency needs a separate table before the first migration can be safe;
- catalog ownership must be selected before balance rows can be created;
- future runtime behavior needs reservation, settlement, refund, or transfer tables earlier than expected;
- a later ADR authorizes a different economy module ownership model or external compatibility surface.

## Follow-Up

- Complete `W-0294`: add the currency wallet migration source.
- Keep repository interfaces, PostgreSQL adapters, runtime behavior, protocol routes, Protobuf source, generated output, startup wiring, dependencies, rewards, purchases, inventory integration, event/audit tables, hosted surfaces, SDKs, distributed runtime, and direct compatibility behind later bounded work items.
