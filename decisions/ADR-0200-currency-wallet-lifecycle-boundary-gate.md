# ADR-0200: Currency Wallet Lifecycle Boundary Gate

Status: Accepted
Date: 2026-06-07
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-currency-wallet-lifecycle-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-currency-wallet-lifecycle-boundary-gate.md`

Related artifacts:

- `docs/currency-wallet-lifecycle-boundary-gate.md`
- `docs/currency-wallet-lifecycle-boundary-gate.zh-CN.md`
- `decisions/ADR-0199-select-next-pitaya-aligned-direction-after-startup-shutdown-hook-map.md`
- `docs/reference-game-server-alignment.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0199` completed the next-direction selection after the Pitaya-aligned startup and shutdown hook source-first map. It selected `define_currency_wallet_lifecycle_boundary_gate` and opened `M-220/W-0292 Define currency wallet lifecycle boundary gate`.

The product roadmap is now back in Phase 3 core game backend modules. Currency and wallets are part of the economy, inventory, rewards, currencies, and progression capability family. They carry balance correctness, idempotency, privacy, permission, transaction, and future integration risks that should be made explicit before any table, repository, protocol, or runtime behavior exists.

## Decision

Define the currency wallet lifecycle boundary gate in:

```text
docs/currency-wallet-lifecycle-boundary-gate.md
docs/currency-wallet-lifecycle-boundary-gate.zh-CN.md
```

The gate records the future semantic scope:

```text
currency catalog read, wallet read, balance read, grant currency, spend currency, balance change recording, transaction read
```

It records future command, query, event, error, permission, lifecycle state, invariant, idempotency, redaction, and test-plan vocabulary without creating contract source files, runtime behavior, protocol routes, generated output, migrations, repository interfaces, PostgreSQL adapters, dependencies, startup wiring, SDKs, hosted surfaces, distributed runtime, or direct Nakama/Pitaya API compatibility.

Every future player-owned wallet read and player-initiated spend must require validated player identity. Metadata-only `player_id` or `session_id` is not proof. Future grants must be service-authoritative and must not be client-authoritative.

The repository check rule is:

```text
runtime.currency_wallet_lifecycle_boundary_gate
```

Open the next bounded work item:

```text
M-221/W-0293 Define currency wallet persistence schema gate
```

## Alternatives Considered

- Implement currency wallet runtime behavior immediately.
- Add balance tables or wallet transaction tables immediately.
- Add protocol routes and Protobuf messages immediately.
- Start with rewards, purchases, inventory pricing, transfers, reservations, or live operations instead.
- Copy an external economy API shape for compatibility.
- Reopen Pitaya distributed runtime implementation after the startup/shutdown vocabulary pass.

## Rationale

Currency wallets are high-value product capability, but they are unsafe to implement without a stable vocabulary for identity, balance mutation authority, idempotency, transaction recording, insufficient balance behavior, wallet lifecycle states, and redaction.

The selected step keeps the repository aligned with Nakama/Hiro-class economy coverage while preserving vibit's contract-first and agent-native maintainability model. It also keeps currency ownership separate from transport, protocol adapters, authentication, storage objects, inventory, rewards, purchases, operations dashboards, and distributed runtime.

The next step should be a persistence schema gate, not migration source or runtime implementation, because wallet correctness depends on durable balance, transaction, and duplicate-prevention posture.

## Agent Reasoning Summary

The agent continued from `W-0292` and treated it as gate-only. It defined the future lifecycle vocabulary, identity posture, invariants, redaction rules, and test expectations. It opened a schema gate as the next bounded continuation and did not add currency wallet behavior, balance tables, wallet transaction behavior, reward integration, protocol routes, generated output, persistence code, repository interfaces, adapters, dependencies, hosted surfaces, SDKs, distributed runtime, or direct compatibility.

## Decision Weights

```yaml
decision_weights:
  nakama_hiro_product_alignment: high
  economy_invariant_clarity: high
  contract_first_safety: high
  identity_and_privacy_risk_control: high
  idempotency_and_transaction_risk_control: high
  future_testability: high
  persistence_readiness: medium
  runtime_scope_change: none
  protocol_scope_change: none
  migration_scope_change: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- `M-220/W-0292` is completed.
- `runtime.currency_wallet_lifecycle_boundary_gate` is registered.
- The currency wallet lifecycle semantic standard and Simplified Chinese translation exist.
- Future currency commands, queries, events, errors, permissions, lifecycle states, invariants, redaction posture, and tests are defined as planning vocabulary.
- `M-221/W-0293 Define currency wallet persistence schema gate` becomes next-ready.
- Runtime behavior, balance tables, wallet transaction behavior, reward integration, inventory integration, purchase behavior, grant/spend execution, audit/event tables, protocol routes, Protobuf source, generated output, migrations, repository interfaces, PostgreSQL adapters, dependencies, startup wiring, authentication/session changes, SDK publication, hosted surfaces, distributed runtime, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- the next persistence schema gate shows that wallet lifecycle and transaction vocabulary need separate gates;
- the roadmap needs rewards, purchases, or inventory pricing before base wallet semantics;
- a later ADR chooses a different economy module ownership model;
- a later ADR explicitly authorizes an external compatibility surface;
- a later ADR authorizes overdraft, transfer, reservation, settlement, paid-currency, or live-operations behavior earlier than planned.

## Follow-Up

- Complete `W-0293`: define the currency wallet persistence schema gate.
- Keep runtime behavior, balance tables, wallet transaction behavior, reward integration, inventory integration, purchase behavior, grant/spend execution, audit/event tables, protocol routes, Protobuf source, generated output, migrations, repository interfaces, PostgreSQL adapters, dependencies, startup wiring, SDKs, hosted surfaces, distributed runtime, and direct compatibility behind later bounded work items.
