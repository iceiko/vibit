# ADR-0209: Currency Wallet Protocol Route Gate

Status: Accepted
Date: 2026-06-08
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-currency-wallet-protocol-route-gate/`

Related conversations:

- `conversations/2026-05-26-currency-wallet-protocol-route-gate.md`

Related artifacts:

- `docs/currency-wallet-protocol-route-gate.md`
- `docs/currency-wallet-protocol-route-gate.zh-CN.md`
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

`M-228/W-0300` implemented application-owned currency wallet runtime behavior under `runtime/internal/app/currency`. The service ensures and reads the validated player's wallet, lists own balances and transactions, records server-authoritative grants, records player-authorized spends, derives owner identity from validated `app.RequestIdentity`, rejects metadata-only identity before unit-of-work access, uses `runtime/internal/modules/currency.Repository` through the application unit-of-work handoff, and maps repository conflicts to redacted public `CURRENCY_WALLET_*` errors.

The next bounded step is to define how this service may later become client-facing over the existing WebSocket/Protobuf protocol. That requires route names, message-shape posture, protected-route policy, protocol adapter ownership, generated-output expectations, public error mapping, local proof expectations, and stop conditions before `.proto` or route code exists.

Nakama is the reference for the product capability: durable wallets, balances, grants, spends, and transaction history are common game backend economy surfaces. Pitaya is the reference for layering: acceptors, sessions, route handlers, serializers, and backend services remain separate. vibit adapts both references without copying either public API.

## Decision

Accept `docs/currency-wallet-protocol-route-gate.md` as the gate for future currency wallet protocol routes.

The gate records:

- future route family `currency.EnsurePlayerWallet`, `currency.GetOwnWallet`, `currency.ListOwnWalletBalances`, `currency.GrantCurrency`, `currency.SpendCurrency`, and `currency.ListOwnWalletTransactions`;
- command/query split for mutation and read routes;
- future Protobuf source candidate `proto/vibit/currency/v1/currency.proto`;
- future generated Go output candidate `runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go`;
- candidate request/response messages and fields;
- unchanged envelope posture;
- `request_token_required` protected-route posture with authenticated wrapper requirement;
- validated authenticated request identity as the only owner source;
- refusal to treat metadata-only `player_id` or `session_id` as proof;
- no client-supplied owner id, wallet id lookup proof, or actor id proof in the first posture;
- server-authoritative grant policy requirement before any future grant route can be exposed;
- future protocol bridge and application handler ownership;
- generated-output and no-hand-edit expectations;
- public service error mapping for the `CURRENCY_WALLET_*` public error family;
- Nakama/Pitaya reference mapping;
- required future tests and local proof expectations;
- stop conditions before route implementation, Protobuf source, generated output, startup wiring, repository/adapter/migration/dependency changes, authentication/session changes, reward/inventory/purchase behavior, catalog/event/payment behavior, hosted deployment, release artifacts, public announcements, paid promotion, SDK work, distributed runtime, or direct compatibility.

This ADR does not add protocol route implementation, Protobuf source, generated output, startup wiring, runtime handlers, repository interface changes, PostgreSQL adapter changes, dependencies, migration changes, authentication/session behavior changes, reward integration, inventory integration, purchase behavior, catalog tables, event/audit tables, payment behavior, reservation behavior, settlement behavior, refund behavior, transfer behavior, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add currency `.proto` messages and route handlers in the same slice.
- Copy Nakama economy route names and response models directly.
- Use a Pitaya-style opaque route string instead of vibit's `kind/module/name` envelope route identity.
- Allow client payloads to include owner ids, wallet ids, or actor ids.
- Treat envelope `player_id` or `session_id` metadata as sufficient wallet owner proof.
- Expose `GrantCurrency` as an unauthenticated or client-authoritative mint route.
- Put currency wallet route behavior in WebSocket transport, generated bridge, or PostgreSQL adapter packages.
- Add reward, inventory, purchase, catalog, event/audit, payment, transfer, or settlement behavior before route boundaries are ratified.

## Rationale

The service behavior exists, but exposing it over the client protocol creates compatibility, identity, economy-integrity, and generated-output risk. A gate-only slice lets the project preserve useful Nakama-style economy capability while keeping Pitaya-style route/session/handler separation and vibit's own contract-first protocol discipline.

The future route family should be own-player scoped because the service is already built around validated player identity. Grants require extra care: the current service uses a system actor for grant transactions, so future protocol exposure must prove that a route cannot become client-authoritative minting. Broader economy features such as rewards, inventory grants, purchase catalogs, payment providers, event/audit history, transfers, reservations, settlement, refunds, operations/admin tools, and SDK/client libraries are product-useful, but they need separate contracts and checks.

## Agent Reasoning Summary

The safest continuation from `W-0300` is a protocol route gate. It records the future route and message shape, connects the work to Nakama/Pitaya capability mapping, and prevents agents from jumping directly into `.proto`, generated code, startup wiring, reward/purchase/payment behavior, event/audit tables, or direct compatibility.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  identity_safety: high
  economy_integrity: high
  protocol_compatibility_control: high
  boundary_clarity: high
  implementation_risk: low
  generated_output_risk: deferred
confidence: high
```

## Consequences

- `docs/currency-wallet-protocol-route-gate.md` and `.zh-CN.md` exist.
- `runtime.currency_wallet_protocol_route_gate` becomes the repository check rule for this slice.
- `M-229/W-0301` is completed.
- The next bounded direction is currency wallet protocol route implementation.
- Protobuf source/generated output, route implementation, startup wiring, repository/adapter/migration/dependency changes, authentication/session semantics, reward/inventory/purchase behavior, catalog/event/payment behavior, hosted deployment, release artifacts, SDK/client libraries, distributed runtime, and direct compatibility remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- route policy changes the first protected-route posture for currency wallets;
- request identity validation semantics change materially;
- currency wallet runtime behavior changes its owner derivation model;
- grant behavior becomes explicitly admin-only, local-proof-only, or excluded from public routes before protocol implementation;
- the project chooses direct Nakama or Pitaya API compatibility through a future ADR;
- currency wallet capability expands to rewards, inventory integration, purchases, catalogs, event/audit history, payment providers, transfers, reservations, settlement, refunds, operations/admin surfaces, or SDK/client libraries before protocol exposure;
- generated-output standards change before the future implementation slice.

## Follow-Up

- Implement currency wallet protocol routes only after this gate is accepted.
- Keep Protobuf source/generated output, startup wiring, reward integration, inventory integration, purchase behavior, catalog/event tables, payment behavior, transfer behavior, SDK publication, generated client libraries, hosted deployment, and direct compatibility behind later bounded work items.
- Preserve metadata-only identity refusal in future route tests.
- Prove grant route policy before exposing a grant route to ordinary client flows.
