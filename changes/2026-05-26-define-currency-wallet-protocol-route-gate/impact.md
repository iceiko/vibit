# Impact Analysis

## Affected Modules

- `currency`
- `runtime`
- `workflow`
- `reference`

## Module Ownership Impact

The currency module remains the semantic owner of storage-neutral wallet vocabulary and repository types. Application currency behavior remains under `runtime/internal/app/currency`. Future protocol bridge ownership is planned under `runtime/internal/platform/protocol/protobuf`, future route handler ownership under `runtime/internal/app/bootstrap`, and future Protobuf source ownership under `proto/vibit/currency/v1`.

No runtime ownership changes are implemented by this gate.

## Public Contract Impact

No public command, query, event, error, permission, or Protobuf contract is added by this gate.

The gate records candidate future protocol routes:

- `currency.EnsurePlayerWallet`
- `currency.GetOwnWallet`
- `currency.ListOwnWalletBalances`
- `currency.GrantCurrency`
- `currency.SpendCurrency`
- `currency.ListOwnWalletTransactions`

## Data And Migration Impact

No migrations or data model changes are added. Existing currency wallet persistence remains the PostgreSQL `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions` tables introduced by `W-0294`.

## Test Impact

No Go tests are required because this is a gate-only documentation and manifest change. Future implementation tests are recorded in the gate standard.

## Documentation Impact

Adds:

- `docs/currency-wallet-protocol-route-gate.md`
- `docs/currency-wallet-protocol-route-gate.zh-CN.md`
- `ADR-0209`
- conversation log
- change spec artifacts

Updates architecture manifests, currency module manifest/guides, repository docs, and continuation pointers.

## Compatibility Risks

This gate reduces compatibility risk by preventing ad hoc currency wallet routes, client-authoritative minting, or copied Nakama/Pitaya API shapes before protocol contracts are ratified.

No wire compatibility changes occur because no `.proto` or generated output is added.
