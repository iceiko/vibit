# Impact

## Affected Modules

The `inventory` module remains the first proof slice.

No Inventory behavior changes.

## Module Ownership Impact

No module ownership moves.

The check reinforces existing ownership:

- Semantic contracts remain under `contracts/`.
- Protobuf source files belong under `proto/vibit/<module>/v1/`.
- Generated Go Protobuf output belongs under `runtime/internal/generated/proto/`.

## Public Contract Impact

No public contract shape changes.

The check derives expected Protobuf message names from existing command, query, and event contract IDs.

## Event Impact

No event contract changes.

The `ItemGranted` event receives an expected future Protobuf message mapping.

## Permission Impact

No permission contract changes.

## Data And Migration Impact

No migrations are added.

## Test Impact

No Go tests are added.

The existing Node CLI checks exercise the new protocol check.

## Tooling Impact

Adds a dedicated `check protocol` command.

The current state passes when expected `.proto` files are absent and the contract registry still declares Protobuf alignment as planned. Once `.proto` files exist, the check validates file path, package, source trace, expected message declarations, and expected field names.

## Documentation Impact

Repository guides, README files, protocol README files, and architecture manifests are updated in English and Simplified Chinese where public-facing text changes.

## Compatibility Risks

The main risk is making the first alignment rule too strict before the first `.proto` file exists. This is reduced by treating missing `.proto` files as acceptable only while the registry status is explicitly planned, then enforcing concrete checks once files appear.
