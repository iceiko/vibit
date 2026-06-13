# currency Module Agent Guide

Status: Draft v0.1

## When To Use This Module

Use this module for currency wallet repository vocabulary, storage-neutral value types, application-owned wallet behavior, and the currently ratified protected currency wallet protocol route family.

The current implemented slice is intentionally narrow:

- `runtime/internal/modules/currency.Repository`
- `CurrencyWallet`, `CurrencyWalletOwner`, `CurrencyWalletBalance`, `CurrencyWalletTransaction`, actor, lifecycle, idempotency, amount, and version value types
- create, wallet lookup, owner lookup, balance list, grant record, spend record, and transaction list input/result types
- conflict classes and redacted repository errors
- normalization helpers and focused Go tests

`M-224 Currency Wallet Repository Interface Implementation` is completed by `W-0296`. The check rule is `runtime.currency_wallet_repository_interface_implementation`.

`W-0297 Define currency wallet PostgreSQL adapter gate` is completed. It accepted `ADR-0205` and registered `runtime.currency_wallet_postgresql_adapter_gate`.

`W-0298 Implement currency wallet PostgreSQL adapter` is completed. It accepted `ADR-0206`, registered `runtime.currency_wallet_postgresql_adapter_implementation`, added `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`, focused fake-executor tests, and `UnitOfWork.NewCurrencyWalletRepository`.

`W-0299 Define currency wallet runtime behavior gate` is completed. It accepted `ADR-0207`, registered `runtime.currency_wallet_runtime_behavior_gate`, and defined `docs/currency-wallet-runtime-behavior-gate.md`.

`W-0300 Implement currency wallet runtime behavior` is completed. It accepted `ADR-0208`, registered `runtime.currency_wallet_runtime_behavior_implementation`, and added `runtime/internal/app/currency/service.go` plus focused service tests.

`W-0301 Define currency wallet protocol route gate` is completed. It accepted `ADR-0209`, registered `runtime.currency_wallet_protocol_route_gate`, and defined `docs/currency-wallet-protocol-route-gate.md`.

`W-0302 Implement currency wallet protocol route` is completed. It accepted `ADR-0210`, registered `runtime.currency_wallet_protocol_route_implementation`, added `proto/vibit/currency/v1/currency.proto`, generated Go output, route keys, bootstrap handlers, protocol bridge mapping, payload registry dispatch, server-side grant policy enforcement, focused tests, and the authenticated local route proof.

No repository next work item is currently ready. Do not invent `W-0303`; keep reward integration, inventory integration, purchase behavior, catalog/event tables, payment behavior, migrations, authentication/session behavior changes, SDK publication, hosted deployments, distributed runtime, and direct Nakama/Pitaya API compatibility behind later bounded work items.

## When Not To Use This Module

Do not use this module for:

- WebSocket or HTTP transport behavior.
- Protobuf or generated wire behavior outside the ratified `vibit.currency.v1` route artifacts.
- PostgreSQL adapter implementation or SQL execution under this module.
- Runtime wallet behavior outside `runtime/internal/app/currency`.
- Reward, inventory, purchase, or payment behavior.
- Player account lifecycle.
- Authentication, token formats, credential storage, or session validation.
- Currency catalog management, event/audit tables, reservations, settlement, refunds, transfers, paid currency, matchmaking, or match runtime.
- Direct Nakama or Pitaya public API compatibility.

If a requirement needs one of those concepts, create or update the owning boundary instead of adding hidden ownership here.

## Extension Points

- Repository interface: `runtime/internal/modules/currency.Repository`
- Repository value types: `CurrencyWallet`, `CurrencyWalletBalance`, `CurrencyWalletTransaction`, `CurrencyWalletOwner`, `CurrencyWalletActor`, `CurrencyWalletVersion`, `CurrencyBalanceVersion`
- Lifecycle vocabulary: `active`, `suspended`, `closed`
- Transaction vocabulary: `grant`, `spend`
- Actor vocabulary: `system`, `player`, `operation`
- Normalizers: wallet records, balance records, transaction records, list results, owner identity, actor identity, idempotency fields, metadata JSON, and repository inputs
- Tests: `runtime/internal/modules/currency/repository_test.go`
- PostgreSQL adapter owner candidate: `runtime/internal/platform/persistence/postgres`
- PostgreSQL adapter implementation: `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`
- Runtime behavior gate: `runtime.currency_wallet_runtime_behavior_gate`
- Runtime behavior implementation: `runtime.currency_wallet_runtime_behavior_implementation`
- Runtime behavior implementation source: `runtime/internal/app/currency/service.go`
- Protocol route gate: `runtime.currency_wallet_protocol_route_gate`
- Protocol route implementation: `runtime.currency_wallet_protocol_route_implementation`
- Protocol route source: `proto/vibit/currency/v1/currency.proto`
- Protocol route bridge: `runtime/internal/platform/protocol/protobuf/currency_bridge.go`
- Protocol route handler: `runtime/internal/app/bootstrap/currency.go`

Future runtime behavior must derive owner and actor identity from validated request identity or service-authoritative context before calling this repository interface; client-supplied player ids are not authentication proof.

## Forbidden Shortcuts

- Do not bypass boundaries declared in `module.yaml`.
- Do not add unregistered public commands, queries, events, errors, or permissions.
- Do not add PostgreSQL adapter code under this module.
- Do not import `pgx`, `database/sql`, WebSocket packages, generated Protobuf packages, SDK packages, or distributed runtime packages into this module.
- Do not execute SQL or mention concrete SQL statements in currency module source.
- Do not change migrations from this module.
- Do not wire transport behavior from this module.
- Do not add new protocol routes, Protobuf sources, or generated output beyond the ratified `W-0302` currency wallet route family.
- Do not store raw credentials, raw tokens, verifier material, lookup digests, verifier digests, cookies, headers, transport subprotocols, connection metadata, payment provider payloads, Nakama API paths, Pitaya server ids, or direct compatibility fields in currency value types.
- Do not treat wallet ids, owner ids, actor ids, `player_id`, `session_id`, idempotency keys, or transport metadata as authenticated proof.

## Required Tests

See `tests.required` in `module.yaml`.

For the current repository interface slice, tests must cover:

- Repository interface storage neutrality.
- Closed owner kind, lifecycle state, transaction kind, actor kind, and conflict vocabularies.
- Wallet record normalization.
- Balance record normalization.
- Transaction record normalization, idempotency, actor, metadata, and delta validation.
- Create/get/owner/list input normalization.
- Grant/spend input normalization and expected-version pointer copying.
- List result copying.
- Redacted conflict and repository errors.
- Absence of secret, transport, protocol, PostgreSQL, session, distributed, payment, and direct compatibility fields.

Run `node tools/vibit check runtime` after changing currency runtime source. When Go is available, also run `cd runtime && go test ./...`.
