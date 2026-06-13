# ADR-0203: Currency Wallet Repository Boundary

Status: Accepted
Date: 2026-06-07
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-currency-wallet-repository-boundary/`

Related conversations:

- `conversations/2026-05-26-currency-wallet-repository-boundary.md`

Related artifacts:

- `docs/currency-wallet-repository-boundary.md`
- `docs/currency-wallet-repository-boundary.zh-CN.md`
- `runtime/migrations/postgres/000008_create_currency_wallets.sql`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0202` added the PostgreSQL migration source for `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions`. The project is sequencing this economy path as lifecycle gate, persistence schema gate, migration source, repository boundary, repository interface, adapter gate, adapter implementation, runtime behavior, protocol route, and local proof.

The work queue reached `M-223/W-0295 Define currency wallet repository boundary`. This slice is a boundary-only step in the Nakama/Hiro-aligned `economy_inventory_rewards_currencies_and_progression` capability family. Pitaya remains deferred as a future distributed architecture reference.

## Decision

Accept `docs/currency-wallet-repository-boundary.md` as the storage-neutral repository boundary for future currency wallet behavior.

The boundary records:

- future repository owner candidate `runtime/internal/modules/currency`;
- future interface candidate `runtime/internal/modules/currency.Repository`;
- future PostgreSQL adapter owner `runtime/internal/platform/persistence/postgres`;
- source migration `runtime/migrations/postgres/000008_create_currency_wallets.sql`;
- logical tables `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions`;
- candidate value types for wallets, owners, balances, currency codes, transactions, actors, idempotency, inputs, conflicts, and repository errors;
- candidate capabilities for wallet creation, wallet lookup, owner lookup, balance listing, grant recording, spend recording, and transaction listing;
- request identity and actor handoff from already validated application identity or server operation identity;
- transaction and idempotency posture;
- optimistic version and typed conflict posture;
- private economy redaction posture;
- PostgreSQL adapter expectations tied to the indexes created by `ADR-0202`;
- future implementation queue boundaries.

This ADR does not add repository interface implementation, PostgreSQL adapter behavior, runtime wallet behavior, grant/spend execution, protocol routes, Protobuf source, generated output, dependencies, migrations, startup wiring, automatic migration application, authentication/session changes, reward integration, inventory integration, purchase behavior, paid-currency behavior, catalog tables, event/audit tables, reservations, settlement, refunds, transfers, SDK publication, hosted deployments, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

Open the next bounded work item:

```text
M-224/W-0296 Implement storage-neutral currency wallet repository interface
```

## Alternatives Considered

- Add the Go repository interface in the same slice.
- Add the PostgreSQL adapter in the same slice.
- Add runtime grant, spend, balance read, wallet read, or transaction read behavior in the same slice.
- Add protocol routes or Protobuf source with the repository boundary.
- Add currency catalog or event/audit tables before repository vocabulary.
- Treat actor ids supplied by clients as proof.
- Copy external Nakama or Pitaya API compatibility.
- Introduce distributed economy routing before single-process repository semantics are proven.

## Rationale

Nakama and Hiro show that wallet, currency, and economy state are core game/backend capabilities. vibit adapts that product need into an agent-native sequence where each layer becomes explicit and checkable before the next layer depends on it.

The repository boundary is needed before interface implementation because future agents need a concise vocabulary for ownership, value types, capabilities, transaction/idempotency behavior, conflicts, redaction, and adapter expectations. Defining that vocabulary now reduces the risk that W-0296 invents ad hoc SQL-shaped types or mixes identity proof with economy state.

Keeping the slice boundary-only reduces risk. It lets static checks prove that no runtime behavior, adapter behavior, protocol shape, generated output, or direct external compatibility was added.

## Agent Reasoning Summary

The agent continued from `W-0294`, kept Nakama and Hiro as economy capability references, and converted the migration source into a storage-neutral future repository plan. The change preserves the user's AI-native product purpose: a future user requirement should become a bounded spec, acceptance criteria, test plan, tests, implementation, verification, and durable memory.

## Decision Weights

```yaml
decision_weights:
  nakama_hiro_product_alignment: high
  ai_native_requirement_test_workflow_alignment: high
  repository_testability: high
  agent_readability: high
  transaction_idempotency_correctness: high
  privacy_and_redaction: high
  runtime_behavior_risk: low
  dependency_expansion: low
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `docs/currency-wallet-repository-boundary.md` and its Simplified Chinese translation exist.
- `runtime.currency_wallet_repository_boundary` becomes the repository check rule for this slice.
- `M-223/W-0295` is completed as a repository-boundary-only milestone.
- The work queue advances to `M-224/W-0296 Implement storage-neutral currency wallet repository interface`.
- Existing runtime behavior is not changed by this ADR.
- Currency wallets are not yet exposed through repository interfaces, PostgreSQL adapters, runtime services, protocol routes, generated output, SDKs, hosted surfaces, or direct external compatibility.

## Reversal Conditions

Revisit this decision if:

- future repository interface implementation proves the candidate vocabulary cannot express wallet, balance, transaction, or idempotency semantics safely;
- privacy or compliance requirements demand a different redaction posture before interface implementation;
- a future runtime behavior gate selects a materially different conflict or leakage model;
- a future ADR explicitly authorizes a different economy module storage model or external compatibility target.

## Follow-Up

- Implement the storage-neutral currency wallet repository interface as `W-0296`.
- Keep PostgreSQL adapter behavior, runtime wallet behavior, protocol routes, Protobuf source, generated output, startup wiring, dependencies, catalog tables, event/audit tables, rewards, purchases, inventory integration, SDKs, hosted surfaces, distributed runtime, and direct compatibility behind later bounded work items.
