# Request: Implement Currency Wallet Runtime Behavior

Implement the bounded application-owned currency wallet runtime behavior slice after the runtime behavior gate.

The implementation must stay under `runtime/internal/app/currency`, use the storage-neutral currency wallet repository only through unit-of-work handoff, derive player wallet ownership from validated request identity, reject metadata-only identity before repository access, and keep protocol routes, Protobuf source, generated output, dependencies, migrations, startup wiring, reward/inventory/purchase integration, catalog/event tables, payment behavior, hosted surfaces, SDKs, distributed runtime, and direct Nakama/Pitaya API compatibility deferred.
