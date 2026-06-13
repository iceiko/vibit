# Impact: Implement Currency Wallet Runtime Behavior

## Affected Areas

- `runtime/internal/app/currency/service.go`
- `runtime/internal/app/currency/service_test.go`
- `runtime/internal/modules/currency/repository.go`
- `runtime/internal/platform/persistence/postgres/currency_wallet_repository.go`
- `.arch/`
- `modules/currency/`
- `tools/vibit`
- `rules/check-rules.json`
- repository guides and continuation docs

## Behavior Added

- Application-owned currency wallet service.
- Validated player identity owner derivation.
- Metadata-only identity refusal before unit-of-work access.
- Player wallet ensure/get operations.
- Own wallet balance and transaction listing.
- Server-authoritative grants using a system actor.
- Player-authorized spends using the validated player actor.
- Redacted repository conflict to public error mapping.
- Unit-of-work repository handoff through `NewCurrencyWalletRepository`.

## Not Added

- Protocol routes.
- Protobuf source.
- Generated output.
- Startup wiring.
- Repository interface changes.
- PostgreSQL adapter changes.
- Migrations.
- Dependencies.
- Authentication/session behavior changes.
- Reward, inventory, purchase, catalog, event/audit, payment, reservation, settlement, refund, or transfer behavior.
- Hosted surfaces, SDKs, release artifacts, distributed runtime, or direct Nakama/Pitaya API compatibility.
