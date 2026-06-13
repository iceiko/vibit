# Request

Advance `W-0293 Define currency wallet persistence schema gate`.

## User Requirement

Define a currency wallet persistence schema gate for vibit so future currency wallet, balance, transaction, and idempotency persistence can be specified before SQL migration source, repositories, adapters, runtime behavior, protocol routes, generated output, or broader economy features are implemented.

## Expected Outcome

The repository records:

- the currency wallet persistence schema gate standard;
- the paired Simplified Chinese translation;
- ADR-0201;
- future PostgreSQL table candidates for wallets, balances, and transactions;
- future migration source candidate `runtime/migrations/postgres/000008_create_currency_wallets.sql`;
- repository check coverage;
- `W-0294 Add currency wallet migration source` as the next-ready work item.

## Non-Goals

- Add SQL migration source.
- Create currency wallet, balance, transaction, catalog, idempotency, reward, purchase, inventory, or audit tables.
- Implement currency wallet runtime behavior.
- Add grant/spend execution.
- Add protocol routes, Protobuf sources, or generated output.
- Add repository interfaces or PostgreSQL adapters.
- Add dependencies or startup wiring.
- Change authentication/session behavior.
- Add hosted surfaces, SDKs, release artifacts, distributed runtime, or direct Nakama/Pitaya API compatibility.
