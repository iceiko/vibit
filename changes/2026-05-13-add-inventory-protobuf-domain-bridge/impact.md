# Impact Analysis

## Affected Modules

- `inventory`

## Module Ownership Impact

No module ownership changes are expected.

The inventory domain module continues to own inventory behavior, repository interfaces, policy interfaces, and invariants. The Protobuf protocol adapter owns generated Protobuf payload conversion and remains the only handwritten runtime layer allowed to import generated Protobuf packages for this bridge.

## Public Contract Impact

No public command, query, event, error, permission, or schema contracts change.

The change adapts existing generated Protobuf message shapes to existing handwritten inventory runtime structs.

## Data And Migration Impact

No database schema, migration, persistence ownership, or durable data behavior changes.

## Test Impact

Add protocol adapter tests for:

- `GrantItemRequest` Protobuf payload to inventory runtime request mapping.
- `GetInventoryRequest` Protobuf payload to inventory runtime request mapping.
- Inventory runtime response to Protobuf response mapping.
- Inventory runtime event to Protobuf event mapping.
- End-to-end envelope, bridge, dispatcher, and response envelope behavior.

## Documentation Impact

Update runtime and module guides to record that the first inventory Protobuf/domain bridge now exists.

Update the runtime protocol adapter standard to document the narrow bridge rule.

## Compatibility Risks

Compatibility risk is low because the change does not alter public contracts or generated wire schemas.

The main architecture risk is generated Protobuf types leaking into domain or application packages. Runtime import-boundary checks and tests should guard against this.
