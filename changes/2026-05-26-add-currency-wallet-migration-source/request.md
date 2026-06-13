# Request

Advance `W-0294 Add currency wallet migration source`.

## User Requirement

Add the first currency wallet PostgreSQL migration source for vibit after the currency wallet persistence schema gate, without adding runtime behavior, protocol routes, generated output, repository interfaces, PostgreSQL adapters, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct compatibility.

## Expected Outcome

The repository records:

- `runtime/migrations/postgres/000008_create_currency_wallets.sql`;
- `ADR-0202`;
- `runtime.currency_wallet_migration_source`;
- change and conversation records for W-0294;
- `W-0295 Define currency wallet repository boundary` as the next-ready work item.

## Non-Goals

- Implement currency wallet behavior.
- Implement wallet transaction behavior.
- Add reward, inventory, purchase, grant, or spend behavior.
- Add protocol routes, Protobuf sources, or generated output.
- Add repository interfaces or PostgreSQL adapters.
- Add dependencies or startup wiring.
- Change authentication/session behavior.
- Add hosted surfaces, SDKs, release artifacts, distributed runtime, or direct Nakama/Pitaya API compatibility.
