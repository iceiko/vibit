# ADR-0207: Currency Wallet Runtime Behavior Gate

Status: Accepted
Date: 2026-06-07
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-currency-wallet-runtime-behavior-gate/`

Related conversations:

- `conversations/2026-05-26-currency-wallet-runtime-behavior-gate.md`

Related artifacts:

- `docs/currency-wallet-runtime-behavior-gate.md`
- `docs/currency-wallet-runtime-behavior-gate.zh-CN.md`
- `runtime/internal/modules/currency/repository.go`
- `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`
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

`M-226/W-0298` implemented the currency wallet PostgreSQL adapter under `runtime/internal/platform/persistence/postgres`. The repository now has a storage-neutral currency wallet repository interface, a PostgreSQL migration source for `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions`, and a platform adapter with unit-of-work repository handoff.

The next bounded step is to define how runtime behavior may later use those persistence pieces. That behavior must derive player wallet ownership from validated request identity, preserve metadata-only identity protections, keep permission and route policy application-owned, map wallet conflicts into stable public error classes, and avoid protocol or generated-output changes until a later protocol slice.

Nakama keeps wallet, balance, grant, spend, and transaction history capability as a core economy primitive. Pitaya reinforces separating route/session context, handlers, and persistence. vibit adapts those lessons by defining an application-owned runtime behavior gate, not by copying public APIs.

## Decision

Accept `docs/currency-wallet-runtime-behavior-gate.md` as the gate for future currency wallet runtime behavior.

The gate records:

- future application owner `runtime/internal/app`;
- future package candidate `runtime/internal/app/currency`;
- future service source and test candidates;
- validated request identity as the only first-posture player wallet owner source;
- explicit refusal to treat metadata-only `player_id` or `session_id` as proof;
- first owner kind `player`;
- actor derivation from validated request identity or service-authoritative server operation;
- route-policy expectation `request_token_required`;
- candidate operations for ensuring a player wallet, reading the player's wallet, listing balances, granting currency, spending currency, and listing transactions;
- validation, permission, conflict, redaction, idempotency, and unit-of-work handoff expectations;
- future runtime behavior implementation direction `implement_currency_wallet_runtime_behavior`;
- stop conditions before runtime implementation, handlers, startup wiring, protocol routes, Protobuf sources, generated output, repository/adapter changes, dependencies, migrations, authentication/session changes, reward/inventory/purchase integration, catalog/event tables, hosted deployment, release artifacts, public announcements, paid promotion, payment/reservation/settlement/refund/transfer behavior, or direct compatibility.

This ADR does not add runtime behavior implementation, runtime handlers, startup wiring, protocol routes, Protobuf sources, generated output, repository interface changes, PostgreSQL adapter changes, dependencies, migration changes, authentication/session behavior changes, reward integration, inventory integration, purchase behavior, catalog tables, event/audit tables, payment behavior, reservation behavior, settlement behavior, refund behavior, transfer behavior, hosted deployments, release artifacts, public announcements, paid promotion, SDKs, distributed runtime, or direct Nakama/Pitaya API compatibility.

Open the next bounded work item:

```text
M-228/W-0300 Implement currency wallet runtime behavior
```

## Alternatives Considered

- Implement runtime currency wallet behavior immediately after the PostgreSQL adapter.
- Add protocol routes and Protobuf messages together with behavior.
- Let client payloads supply wallet owner ids or actor ids.
- Treat envelope metadata `player_id` as sufficient owner proof.
- Put runtime behavior inside the currency domain module or PostgreSQL adapter package.
- Add reward, inventory, purchase, payment, catalog, reservation, settlement, refund, or transfer behavior in the first runtime slice.
- Add direct Nakama-compatible wallet or economy routes.

## Rationale

Currency wallet runtime behavior has identity, economy, and conflict-mapping risk. The first posture must not let clients select arbitrary owners, turn metadata-only player ids into proof, mutate another player's wallet, leak private wallet state, or hide economy behavior inside transport, protocol, or persistence layers. A gate-only ADR keeps the next implementation slice focused on application behavior and makes route/protocol/economy integration a separate explicit decision.

Separating this gate from implementation also lets checks reject accidental runtime handlers, Protobuf shapes, generated output, catalog/event tables, reward/inventory/purchase coupling, payment behavior, or direct compatibility before those surfaces are ratified.

## Agent Reasoning Summary

The safest continuation from the PostgreSQL adapter is an application behavior gate. It defines how future services should derive wallet owner identity, derive actors, require protected routes, validate inputs, map conflicts, and use unit-of-work repository handoff while preserving protocol, authentication/session, reward, inventory, purchase, catalog, event/audit, payment, and broad product deferrals.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  identity_safety: high
  economy_integrity: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: low
  protocol_risk: deferred
  dependency_expansion: low
confidence: high
```

## Consequences

- `docs/currency-wallet-runtime-behavior-gate.md` and `.zh-CN.md` exist.
- `runtime.currency_wallet_runtime_behavior_gate` becomes the repository check rule for this slice.
- `M-227/W-0299` is completed.
- `M-228/W-0300 Implement currency wallet runtime behavior` becomes the next-ready work item.
- Runtime behavior implementation, protocol behavior, generated output, migrations, dependencies, reward/inventory/purchase integration, catalog/event tables, payment behavior, and authentication/session behavior remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- request identity validation semantics change materially;
- route policy selects a different first protected-route posture for currency wallets;
- wallet ownership expands beyond player-owned wallets before implementation;
- grant/spend behavior must be coupled to reward, inventory, purchase, catalog, or payment before runtime behavior implementation;
- the repository interface or PostgreSQL adapter changes before runtime behavior implementation;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Implement currency wallet runtime behavior only after this gate is accepted.
- Keep protocol routes, Protobuf messages, generated output, startup wiring, reward integration, inventory integration, purchase behavior, catalog tables, event/audit tables, payment provider behavior, reservation, settlement, refund, transfer, SDKs, hosted surfaces, distributed runtime, and direct compatibility behind later gates.
- Preserve metadata-only identity refusal and server-derived wallet owner behavior in implementation tests.
