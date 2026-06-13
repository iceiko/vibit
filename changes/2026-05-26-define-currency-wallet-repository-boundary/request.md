# Request

Advance `W-0295 Define currency wallet repository boundary`.

## User Requirement

Define the currency wallet storage-neutral repository boundary after the migration source, without adding repository interface implementation, PostgreSQL adapters, runtime behavior, protocol routes, generated output, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct compatibility.

## Expected Outcome

The repository records:

- the currency wallet repository boundary standard;
- the paired Simplified Chinese translation;
- ADR-0203;
- future repository owner, interface candidate, value types, capabilities, conflicts, transaction/idempotency handoff, PostgreSQL adapter expectations, and tests;
- repository check coverage;
- `W-0296 Implement storage-neutral currency wallet repository interface` as the next-ready work item.

## Non-Goals

- Add Go repository interface implementation.
- Add PostgreSQL adapter behavior.
- Execute SQL from Go.
- Implement currency wallet runtime behavior.
- Add grant/spend execution.
- Add reward, inventory, purchase, catalog, event/audit, reservation, settlement, refund, or transfer behavior.
- Add protocol routes, Protobuf sources, or generated output.
- Add dependencies or startup wiring.
- Change authentication/session behavior.
- Add hosted surfaces, SDKs, release artifacts, distributed runtime, or direct Nakama/Pitaya API compatibility.
