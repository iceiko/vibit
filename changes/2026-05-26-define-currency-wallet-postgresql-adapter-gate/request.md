# Request

Define the currency wallet PostgreSQL adapter gate for the Nakama/Hiro-aligned economy, currency, and wallet capability family.

The gate must stay inside the adapter-boundary slice:

- define the future PostgreSQL adapter owner, source candidate, test candidate, and constructor posture;
- preserve the storage-neutral `runtime/internal/modules/currency.Repository` interface as the adapter contract;
- map the future SQL posture to the existing `currency_wallets`, `currency_wallet_balances`, and `currency_wallet_transactions` migration source without adding SQL execution behavior;
- define transaction/unit-of-work handoff, idempotency mapping, version mapping, insufficient-balance conflict mapping, redaction, and test expectations;
- open a bounded follow-up for implementation;
- keep actual PostgreSQL adapter code, runtime wallet behavior, grant/spend execution, protocol routes, Protobuf source, generated output, dependencies, migration changes, authentication/session behavior, reward integration, inventory integration, purchase behavior, catalog/event tables, SDKs, hosted deployment, distributed runtime, and direct Nakama/Pitaya API compatibility deferred.

## User-Facing Need

The project is aiming at a Nakama-class game/backend server framework shaped for AI-native development and testing. Currency wallets are a core economy primitive, but adapter behavior should be explicit and testable before it is implemented.

This slice gives future agents a precise persistence-adapter contract so a later implementation can be written from a bounded spec and test plan.

## Stop Conditions

- PostgreSQL adapter implementation or SQL execution is required.
- Runtime wallet behavior, grant/spend execution above the adapter, or integration behavior is required.
- Protocol routes, Protobuf source, or generated output are required.
- Authentication/session behavior, request identity validation, or startup wiring changes are required.
- Reward integration, inventory integration, purchase behavior, catalog/event tables, SDKs, hosted deployment, distributed runtime, or direct compatibility are required.
- Verification fails.
