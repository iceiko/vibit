# ADR-0202: Currency Wallet Migration Source

Status: Accepted
Date: 2026-06-07
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-add-currency-wallet-migration-source/`

Related conversations:

- `conversations/2026-05-26-currency-wallet-migration-source.md`

Related artifacts:

- `runtime/migrations/postgres/000008_create_currency_wallets.sql`
- `docs/currency-wallet-persistence-schema-gate.md`
- `docs/currency-wallet-persistence-schema-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0201` completed the currency wallet persistence schema gate. It selected PostgreSQL as the first durable currency wallet store and named `runtime/migrations/postgres/000008_create_currency_wallets.sql` as the future migration source candidate.

The next bounded step is migration-source-only: add the SQL source for wallets, balances, and transactions while preserving repository, adapter, runtime, protocol, generated output, dependency, hosted, SDK, distributed runtime, and direct compatibility deferrals.

## Decision

Add the PostgreSQL migration source:

```text
runtime/migrations/postgres/000008_create_currency_wallets.sql
```

The migration defines three logical tables:

```text
currency_wallets
currency_wallet_balances
currency_wallet_transactions
```

The migration records the first player-owned wallet identity, balance row identity, transaction fact, idempotency uniqueness, timestamp, and redaction posture from `ADR-0201`.

The repository check rule is:

```text
runtime.currency_wallet_migration_source
```

This ADR does not add currency wallet repository interfaces, PostgreSQL adapters, SQL execution code, runtime wallet behavior, grant/spend execution, reward integration, inventory integration, purchase behavior, protocol routes, Protobuf source, generated output, dependencies, startup wiring, authentication/session behavior changes, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

Open the next bounded work item:

```text
M-223/W-0295 Define currency wallet repository boundary
```

## Alternatives Considered

- Add repository interfaces in the same slice.
- Add PostgreSQL adapter code in the same slice.
- Add runtime grant/spend behavior with the migration.
- Add a currency catalog table immediately.
- Add event/audit tables immediately.
- Copy a direct Nakama or Pitaya economy API or schema compatibility surface.

## Rationale

Migration source belongs before repository and adapter code because it fixes durable table names, key shapes, constraints, and index posture. Keeping this slice migration-source-only makes later repository boundary and adapter work easier to verify.

The first migration includes a current-state wallet table, a current-state balance table, and a transaction fact table. The transaction table owns idempotency uniqueness for future grant/spend behavior without adding that behavior now.

The currency catalog remains deferred because normalization, display metadata, precision, mutability, and paid-currency rules need their own bounded decision before catalog tables exist.

## Agent Reasoning Summary

The agent continued from `W-0294`, followed `ADR-0201`, added only the PostgreSQL migration source, registered the repository check, recorded W-0294 completion, and opened `W-0295` as the repository boundary follow-up.

## Decision Weights

```yaml
decision_weights:
  migration_source_first_workflow_alignment: high
  ai_native_verifiability: high
  durable_balance_correctness: high
  idempotency_and_supportability: high
  runtime_scope_control: high
  protocol_scope_control: high
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `M-222/W-0294` is completed.
- `runtime.currency_wallet_migration_source` is registered.
- `runtime/migrations/postgres/000008_create_currency_wallets.sql` exists.
- `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions` become planned durable SQL source tables.
- `M-223/W-0295 Define currency wallet repository boundary` becomes next-ready.
- Repository interfaces, PostgreSQL adapters, runtime behavior, protocol routes, Protobuf source, generated output, startup wiring, dependencies, rewards, purchases, inventory integration, event/audit tables, hosted surfaces, SDKs, distributed runtime, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- the repository boundary proves the table identity cannot express future wallet, balance, transaction, or idempotency semantics;
- catalog ownership must be selected before balance rows can safely exist;
- future runtime behavior requires reservation, settlement, refund, transfer, or event/audit tables before repository interfaces;
- a later ADR authorizes a different economy module ownership model or external compatibility surface.

## Follow-Up

- Complete `W-0295`: define the currency wallet repository boundary.
- Keep repository implementation, PostgreSQL adapters, runtime behavior, protocol routes, Protobuf source, generated output, startup wiring, dependencies, rewards, purchases, inventory integration, event/audit tables, hosted surfaces, SDKs, distributed runtime, and direct compatibility behind later bounded work items.
