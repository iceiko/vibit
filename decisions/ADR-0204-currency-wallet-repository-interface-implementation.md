# ADR-0204: Currency Wallet Repository Interface Implementation

Status: Accepted
Date: 2026-06-07
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-currency-wallet-repository-interface/`

Related conversations:

- `conversations/2026-05-26-currency-wallet-repository-interface-implementation.md`

Related artifacts:

- `runtime/internal/modules/currency/repository.go`
- `runtime/internal/modules/currency/repository_test.go`
- `modules/currency/module.yaml`
- `modules/currency/AGENTS.md`
- `modules/currency/AGENTS.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-223/W-0295` defined the currency wallet repository boundary after the PostgreSQL `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions` migration source. The next bounded step was to turn that boundary into storage-neutral Go repository vocabulary without adding PostgreSQL adapters, SQL execution, runtime wallet behavior, protocol routes, Protobuf source, generated output, dependencies, authentication/session behavior changes, reward integration, inventory integration, purchase behavior, distributed runtime, or direct external compatibility.

Nakama and Hiro guide the product capability pressure: currency wallets are a common economy/progression primitive. Pitaya remains deferred as a future distributed architecture reference. vibit adapts the capability through module-owned repository interfaces, explicit tests, and repository checks instead of copying public APIs.

## Decision

Implement the storage-neutral repository interface under:

```text
runtime/internal/modules/currency
```

The package defines:

- `runtime/internal/modules/currency.Repository`;
- wallet, owner, balance, transaction, actor, currency code, idempotency, amount, and version value types;
- closed first-posture vocabularies for owner kind `player`, lifecycle states `active`, `suspended`, and `closed`, transaction kinds `grant` and `spend`, and actor kinds `system`, `player`, and `operation`;
- create wallet, get wallet, get wallet for owner, list balances, record grant, record spend, and list transactions input/result types;
- conflict classes including insufficient balance, duplicate transaction, stale versions, invalid idempotency, and storage unavailable;
- redacted repository error types;
- normalization helpers for records, list results, owner identity, actor identity, idempotency fields, metadata JSON, and repository inputs;
- focused tests for storage neutrality, closed vocabularies, normalization, idempotency input handling, expected-version copying, result copying, redaction, and forbidden coupling.

Add the first currency module manifest and module AGENTS guides so future agents can discover ownership before adding adapters or runtime behavior.

This ADR does not add PostgreSQL currency adapters, SQL execution behavior, unit-of-work factory wiring, runtime wallet handlers, grant/spend execution, reward integration, inventory integration, purchase behavior, payment behavior, protocol routes, Protobuf sources, generated output, dependencies, migration changes, authentication/session behavior changes, currency catalog tables, event/audit tables, reservations, settlement, refunds, transfers, SDK publication, hosted deployments, release artifacts, distributed runtime, or direct Nakama/Pitaya API compatibility.

Open the next bounded work item:

```text
M-225/W-0297 Define currency wallet PostgreSQL adapter gate
```

## Alternatives Considered

- Implement the PostgreSQL adapter in the same slice.
- Add runtime wallet creation, grant, spend, balance read, or transaction read behavior with the repository interface.
- Add Protobuf messages or public routes immediately.
- Place the interface under `runtime/internal/app`.
- Treat owner ids, actor ids, wallet ids, or metadata as authentication proof.
- Add currency catalog, event/audit, reward, inventory, purchase, or payment behavior before adapter gates.
- Copy external Nakama or Pitaya public API compatibility.

## Rationale

The repository boundary already selected the owner candidate and capability vocabulary. Implementing only the storage-neutral interface now reduces future adapter ambiguity while keeping SQL, behavior, protocol, generated output, and broader economy concerns behind later gates.

Putting the interface in `runtime/internal/modules/currency` makes currency wallets a first-class domain module without making it own player accounts, authentication, sessions, transport behavior, inventory, rewards, purchases, payments, matchmaking, match runtime, or distributed topology.

## Agent Reasoning Summary

The safest continuation from `W-0295` was an interface-only code slice. It gives later PostgreSQL adapter work a stable typed contract, adds tests for redaction and storage neutrality, and preserves stop conditions that keep protocol/runtime behavior from leaking into this package.

## Decision Weights

```yaml
decision_weights:
  nakama_hiro_product_alignment: high
  ai_native_requirement_test_workflow_alignment: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: low
  adapter_risk: deferred
  protocol_risk: low
  dependency_expansion: low
confidence: high
```

## Consequences

- `runtime/internal/modules/currency/repository.go` exists.
- `runtime/internal/modules/currency/repository_test.go` exists.
- `modules/currency/module.yaml` and paired module guides exist.
- `runtime.currency_wallet_repository_interface_implementation` becomes the repository check rule for this slice.
- `M-224/W-0296` is completed.
- `M-225/W-0297 Define currency wallet PostgreSQL adapter gate` becomes the next-ready work item.
- Existing runtime behavior, protocol behavior, migrations, generated output, and dependencies are unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- currency wallets stop being a module-owned economy capability;
- the first adapter needs a different repository owner;
- the `currency_wallets`, `currency_wallet_balances`, or `currency_wallet_transactions` table shape changes in a way that invalidates the value vocabulary;
- a later ADR selects a materially different conflict, idempotency, or leakage model;
- catalog, payment, purchase, reward, or inventory requirements become mandatory before adapter work;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Define the currency wallet PostgreSQL adapter gate.
- Implement the adapter only after the gate is accepted.
- Define runtime behavior, permissions, protocol routes, generated output, and local proof only after repository and adapter boundaries are accepted.
- Keep currency catalogs, event/audit tables, rewards, purchases, payments, inventory integration, SDKs, hosted surfaces, distributed runtime, and direct compatibility behind future gates.
