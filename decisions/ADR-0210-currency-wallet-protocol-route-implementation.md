# ADR-0210: Currency Wallet Protocol Route Implementation

Status: Accepted
Date: 2026-06-08
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-currency-wallet-protocol-route/`

Related conversations:

- `conversations/2026-05-26-currency-wallet-protocol-route-implementation.md`

Related artifacts:

- `proto/vibit/currency/v1/currency.proto`
- `runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go`
- `runtime/internal/app/currency/routes.go`
- `runtime/internal/app/currency/routes_test.go`
- `runtime/internal/app/bootstrap/currency.go`
- `runtime/internal/app/bootstrap/currency_test.go`
- `runtime/internal/platform/protocol/protobuf/currency_bridge.go`
- `runtime/internal/platform/protocol/protobuf/currency_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`
- `runtime/internal/platform/protocol/protobuf/payload_registry.go`
- `runtime/internal/platform/protocol/protobuf/frame_handler.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `modules/currency/module.yaml`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-229/W-0301` accepted `ADR-0209` and defined the currency wallet protocol route gate. The gate authorized a later bounded implementation of protected currency wallet command/query routes over the existing WebSocket/Protobuf request flow.

The application-owned currency wallet service already exists under `runtime/internal/app/currency`. It derives wallet ownership from validated request identity, refuses metadata-only identity before unit-of-work access, records server-authoritative grants, records player-authorized spends, lists own balances and transactions, uses `runtime/internal/modules/currency.Repository` through unit-of-work handoff, and maps repository failures to redacted public `CURRENCY_WALLET_*` errors.

Nakama provides capability pressure: durable wallets, balances, grants, spends, and transaction history are common game backend economy surfaces. Pitaya provides layering pressure: acceptors, sessions, protocol serializers, route handlers, backend service behavior, and persistence remain separate. vibit adapts both references without direct public API compatibility.

## Decision

Implement the protected currency wallet protocol route family authorized by `ADR-0209`.

The implementation adds:

- Protobuf source `proto/vibit/currency/v1/currency.proto`;
- generated Go Protobuf output `runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go` through Buf;
- route keys under `runtime/internal/app/currency/routes.go`;
- routes `currency.EnsurePlayerWallet`, `currency.GetOwnWallet`, `currency.ListOwnWalletBalances`, `currency.GrantCurrency`, `currency.SpendCurrency`, and `currency.ListOwnWalletTransactions`;
- protocol bridge mapping under `runtime/internal/platform/protocol/protobuf/currency_bridge.go`;
- payload registry integration for currency request and response messages;
- bootstrap handlers under `runtime/internal/app/bootstrap/currency.go`;
- a server-side `CurrencyGrantPolicy` requirement before `GrantCurrency` route execution;
- focused route, bridge, and bootstrap tests;
- an authenticated local gameplay proof in `TestCurrencyWalletProtocolRouteLocalAlphaFlow`;
- repository check coverage through `runtime.currency_wallet_protocol_route_implementation`.

The Protobuf payloads intentionally omit client-supplied owner ids, player ids, session ids, wallet lookup ids, actor ids, raw access tokens, credential material, lookup digests, verifier digests, payment provider payloads, catalog ids, reward ids, and inventory item ids. Owner identity continues to come from authenticated request identity, and grant actor identity comes from server-side route policy.

The authenticated frame handler now accepts empty inner payload bytes only when the registered Protobuf descriptor for the inner payload type has zero fields. This preserves malformed-payload rejection for non-empty messages while allowing proto3 zero-field requests such as `EnsurePlayerWalletRequest` and `GetOwnWalletRequest`.

This ADR does not add reward integration, inventory integration, purchase behavior, catalog tables, event/audit tables, payment behavior, reservation behavior, settlement behavior, refund behavior, transfer behavior, operations/admin behavior, migrations, dependency additions, startup wiring, authentication/session semantic changes, hosted deployment, SDK publication, generated client libraries, release artifacts, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Keep currency wallet behavior application-only until a broader economy proof exists.
- Add route keys and handlers but defer Protobuf source and generated output.
- Let client payloads include owner ids, wallet ids, or actor ids.
- Expose `GrantCurrency` without a server-side grant policy.
- Put route behavior in WebSocket transport, generated Protobuf output, or PostgreSQL adapter packages.
- Combine currency wallet routes with reward, inventory, purchase, catalog, payment, event/audit, reservation, transfer, or admin behavior.
- Treat empty authenticated inner payload bytes as malformed for every message, including zero-field proto3 requests.

## Rationale

The route gate already selected the safe first surface: validated-player wallet operations and a server-authoritative grant posture. Implementing that route family makes the already-ratified currency service visible through the same protected WebSocket/Protobuf architecture used by inventory, storage, and friends routes without broadening economy scope.

The most important safety properties are identity and minting control. Wallet owner proof stays in validated authenticated request identity. The grant route cannot mint currency unless an application-owned route policy explicitly allows the operation and supplies a system actor id.

Allowing zero-field proto3 request payloads is necessary because valid proto3 messages with no fields encode to zero bytes. The implementation checks the inner payload type descriptor before allowing an empty payload, so the authenticated wrapper still rejects empty payloads for messages that have fields.

## Agent Reasoning Summary

The smallest product-useful continuation after the route gate is to wire the existing currency wallet service into the established route/protocol/bootstrap layers and prove the authenticated flow locally. This advances Nakama-class wallet capability while preserving Pitaya-style layering and vibit's generated-output, identity-safety, and boundary-check rules.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  identity_safety: high
  economy_integrity: high
  protocol_compatibility_control: high
  boundary_clarity: high
  implementation_risk: medium
  generated_output_risk: medium
confidence: high
```

## Consequences

- `proto/vibit/currency/v1/currency.proto` exists.
- `runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go` exists and is generated through Buf.
- `runtime/internal/app/currency/routes.go` exposes the currency wallet route keys.
- `runtime/internal/app/bootstrap/currency.go` registers application handlers and requires grant policy for `GrantCurrency`.
- `runtime/internal/platform/protocol/protobuf/currency_bridge.go` maps currency Protobuf payloads.
- `TestCurrencyWalletProtocolRouteLocalAlphaFlow` proves the protected route family in a local authenticated gameplay path.
- `runtime.currency_wallet_protocol_route_implementation` becomes the repository check rule for this slice.
- `M-230/W-0302` is completed.
- No repository next work item is currently ready.
- Broader economy features, startup wiring, hosted deployment, release artifact expansion, SDK/client libraries, distributed runtime, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- route policy for protected gameplay requests changes materially;
- request identity validation semantics change materially;
- currency wallet behavior changes its owner derivation model;
- generated-output standards change the Protobuf source or generated Go output path;
- direct Nakama or Pitaya public API compatibility becomes an explicit future goal through a separate ADR;
- a later local or startup proof shows that the route surface cannot support the intended alpha request flow without changing protocol shape;
- the project authorizes reward, inventory, purchase, catalog, event/audit, payment, reservation, settlement, refund, transfer, or admin economy behavior.

## Follow-Up

- Do not invent `W-0303` without a bounded change spec and maintainer direction.
- Preserve the current currency route Protobuf shape unless a later proof finds a concrete compatibility issue.
- Keep reward integration, inventory integration, purchase behavior, catalog/event/payment behavior, migrations, dependencies, startup wiring, hosted deployment, SDK/client libraries, distributed runtime, and direct compatibility behind later bounded work items.
