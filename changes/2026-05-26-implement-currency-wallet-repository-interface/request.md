# Request

Implement the storage-neutral currency wallet repository interface for the economy/inventory/rewards/currencies/progression capability family.

The implementation must stay inside the repository-interface slice:

- add module-owned Go repository vocabulary under `runtime/internal/modules/currency`;
- preserve wallet, balance, transaction, idempotency, validated identity handoff posture, lifecycle state vocabulary, version/conflict semantics, and redacted errors from `ADR-0203`;
- add focused tests for normalization, idempotency input handling, copying, redaction, and forbidden dependency posture;
- add module manifest and module AGENTS guidance;
- keep PostgreSQL adapter behavior, SQL execution, runtime wallet behavior, protocol routes, Protobuf source, generated output, dependencies, migration changes, reward integration, inventory integration, purchase behavior, SDKs, hosted deployment, distributed runtime, and direct Nakama/Pitaya API compatibility deferred.

## User-Facing Need

The project is aiming at a Nakama-class game/backend server framework shaped for AI-native development and testing. Currency wallets are a core economy primitive, but future runtime behavior should not depend on ad hoc SQL-shaped or transport-shaped vocabulary.

This slice gives future agents a typed, storage-neutral repository contract before adapter and runtime behavior are authorized.

## Stop Conditions

- PostgreSQL adapter implementation or SQL execution is required.
- Runtime wallet creation, grant, spend, balance read, or transaction read behavior is required.
- Protocol routes, Protobuf source, or generated output are required.
- Authentication/session behavior, request identity validation, or startup wiring changes are required.
- Reward integration, inventory integration, purchase behavior, payment behavior, currency catalog tables, event/audit tables, SDKs, hosted deployment, distributed runtime, or direct compatibility are required.
- Verification fails.
