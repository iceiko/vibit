# ADR-0208: Currency Wallet Runtime Behavior Implementation

Status: Accepted
Date: 2026-06-07
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-currency-wallet-runtime-behavior/`

Related conversations:

- `conversations/2026-05-26-currency-wallet-runtime-behavior-implementation.md`

Related artifacts:

- `docs/currency-wallet-runtime-behavior-gate.md`
- `runtime/internal/app/currency/service.go`
- `runtime/internal/app/currency/service_test.go`
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

`M-227/W-0299` defined the currency wallet runtime behavior gate after the storage-neutral repository, migration source, and PostgreSQL adapter were in place. The gate authorized a later bounded application-owned runtime behavior implementation under `runtime/internal/app`, with `runtime/internal/app/currency` as the package candidate.

The currency module owns storage-neutral wallet value types, normalizers, lifecycle vocabulary, transaction vocabulary, actor vocabulary, idempotency vocabulary, and repository error vocabulary. The PostgreSQL adapter owns SQL mapping and unit-of-work repository construction. Runtime behavior must compose these pieces through application dependencies and the unit-of-work boundary without importing protocol, transport, generated Protobuf, migration, reward, inventory, purchase, payment, catalog, event/audit, hosted, SDK, or distributed runtime packages.

Nakama motivates durable wallets, balances, grants, spends, and transaction history as a core economy primitive. Pitaya reinforces separating route/session context, handlers, persistence, and distributed runtime behavior. vibit adapts those references through validated request identity and application-owned behavior, not direct public API compatibility.

## Decision

Implement currency wallet runtime behavior under:

```text
runtime/internal/app/currency
```

The implementation adds:

- `Service` and `NewService`;
- `EnsurePlayerWallet`;
- `GetOwnWallet`;
- `ListOwnWalletBalances`;
- `GrantCurrency`;
- `SpendCurrency`;
- `ListOwnWalletTransactions`;
- request and result vocabulary for application callers;
- stable redacted public error codes for later handlers to map;
- validated player owner derivation from `app.RequestIdentity`;
- metadata-only `player_id` and `session_id` refusal before unit-of-work access or repository calls;
- server-owned wallet id generation for wallet creation;
- server-owned transaction id generation for grants and spends;
- system actor posture for server-authoritative grants;
- player actor posture for player-authorized spends;
- unit-of-work currency repository handoff through `NewCurrencyWalletRepository`;
- redacted conflict mapping for invalid request, unauthenticated, not found, already exists, inactive wallet, insufficient balance, duplicate transaction, version mismatch, and unavailable cases;
- focused fake-repository tests.

This ADR does not add runtime protocol handlers, route registration, startup wiring, Protobuf sources, generated output, repository interface changes, PostgreSQL adapter changes, dependencies, migration changes, authentication/session behavior changes, reward integration, inventory integration, purchase behavior, catalog tables, event/audit tables, payment behavior, reservation behavior, settlement behavior, refund behavior, transfer behavior, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

Open the next bounded work item:

```text
M-229/W-0301 Define currency wallet protocol route gate
```

## Alternatives Considered

- Add WebSocket/Protobuf routes together with the service.
- Put the service under `runtime/internal/modules/currency`.
- Let client payloads provide wallet owner ids, wallet ids, or actor ids.
- Treat metadata-only envelope/session fields as sufficient owner proof.
- Expose repository or SQL details directly to future callers.
- Implement reward, inventory, purchase, payment, catalog, event/audit, reservation, settlement, refund, or transfer behavior together with runtime behavior.
- Copy Nakama economy API semantics directly.

## Rationale

The useful next step after the runtime behavior gate is an application service that future protocol slices can call. This keeps identity enforcement, input validation, unit-of-work orchestration, public failure classes, idempotency handoff, and redaction in the application layer while preserving protocol and generated-output decisions for later gates.

Rejecting metadata-only identity before repository access prevents a client from turning unauthenticated envelope or session metadata into wallet authority. Server-derived owner identity prevents cross-owner wallet access in the first posture. Keeping grants and spends inside a narrow service makes economy behavior testable before route and protocol exposure.

## Agent Reasoning Summary

The safest continuation from `W-0299` was an application service only. It makes the currency wallet repository usable inside the runtime while keeping protocol exposure, generated files, migrations, dependencies, reward/inventory/purchase coupling, catalog/event tables, payment behavior, and broader economy features behind later work items.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  identity_safety: high
  economy_integrity: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: medium
  protocol_risk: deferred
  dependency_expansion: low
confidence: high
```

## Consequences

- `runtime/internal/app/currency/service.go` exists.
- `runtime/internal/app/currency/service_test.go` exists.
- `runtime.currency_wallet_runtime_behavior_implementation` becomes the repository check rule for this slice.
- `M-228/W-0300` is completed.
- `M-229/W-0301 Define currency wallet protocol route gate` becomes the next-ready work item.
- Protocol routes, Protobuf source/generated output, startup wiring, migrations, dependencies, authentication/session semantics, reward integration, inventory integration, purchase behavior, catalog/event tables, payment behavior, SDKs, hosted deployment, release artifacts, distributed runtime, and direct compatibility remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- request identity validation semantics change materially;
- the currency wallet repository interface changes materially;
- route policy selects a different first protected-route posture for currency wallets;
- wallet ownership expands beyond player-owned wallets before protocol exposure;
- reward, inventory, purchase, catalog, event/audit, or payment behavior becomes mandatory before route exposure;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Define currency wallet protocol routes and generated output in a later gate before exposing this service over WebSocket/Protobuf.
- Preserve metadata-only identity refusal in future handlers.
- Keep reward, inventory, purchase, catalog, event/audit, payment, reservation, settlement, refund, transfer, SDKs, hosted surfaces, distributed runtime, and direct compatibility behind later gates.
