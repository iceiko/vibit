# Impact Analysis

## Affected Modules

- `inventory`

Inventory is affected only through tests that prove public inventory application errors can be encoded into protocol error envelopes.

## Module Ownership Impact

No module ownership changes are expected.

The protocol adapter owns the error envelope mapping. Application dispatch and domain modules continue to expose vibit-owned error values without importing generated Protobuf packages.

## Public Contract Impact

No public command, query, event, error, permission, or schema contracts change.

This change uses existing `protocolv1.Error` envelope fields and existing application/domain error codes.

## Data And Migration Impact

No database schema, migration, persistence ownership, or durable data behavior changes.

## Test Impact

Add tests for:

- Building an error envelope from an application error result.
- Preserving request, route, target, and session metadata.
- Encoding domain-owned inventory errors such as `INVENTORY_PERMISSION_DENIED`.
- Keeping successful inventory bridge behavior unchanged.

## Documentation Impact

Update runtime protocol adapter docs and runtime state docs to record the first application error envelope mapping.

Update the Simplified Chinese translations in the same change.

## Compatibility Risks

Compatibility risk is low because the wire schema already contains `Error`.

The main risk is leaking internal implementation details. This change maps only `ApplicationError.Message`, which is currently the public application/domain error message. Internal non-application errors still remain outside this mapping until a separate protocol error handling decision exists.
