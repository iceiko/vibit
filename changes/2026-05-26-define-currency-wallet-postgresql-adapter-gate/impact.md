# Impact

## Added

- `docs/currency-wallet-postgresql-adapter-gate.md`
- `docs/currency-wallet-postgresql-adapter-gate.zh-CN.md`
- `ADR-0205`
- `runtime.currency_wallet_postgresql_adapter_gate`
- `M-226/W-0298 Implement currency wallet PostgreSQL adapter` as the next-ready follow-up

## Behavioral Impact

No runtime behavior is added.

The new standard defines future adapter ownership, SQL mapping posture, unit-of-work handoff, idempotency mapping, version mapping, insufficient-balance conflict posture, redaction, and test expectations only. It does not register handlers, expose routes, execute SQL, parse transport credentials, validate sessions, mutate rewards or inventory, process purchases, or change startup composition.

## Data Impact

No migration is added or changed.

The existing `runtime/migrations/postgres/000008_create_currency_wallets.sql` remains the source migration for future adapter work, but this slice does not execute it or map SQL rows.

## Protocol Impact

No protocol route, Protobuf source, generated output, bridge, route key, or client SDK is added.

## Security And Privacy Impact

The gate treats wallet ids, owner ids, transaction ids, actor ids, idempotency material, balances, metadata, SQL, DSNs, credentials, token material, and driver details as not log-safe by default. Future adapter errors must collapse driver, constraint, row, and private wallet details into redacted currency module errors.

## Reference Alignment

Nakama and Hiro guide the product capability need for durable currency wallet state. vibit adapts that need into a PostgreSQL adapter gate before implementation. Pitaya remains deferred as a future distributed architecture reference.
