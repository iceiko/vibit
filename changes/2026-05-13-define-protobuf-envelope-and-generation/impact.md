# Impact Analysis

## Affected Modules

- `inventory`

## Module Ownership Impact

No module ownership changes.

The inventory module remains the semantic owner of inventory commands, queries, events, errors, permissions, and invariants. The new `.proto` file defines inventory wire message shape only. The protocol envelope remains owned by platform protocol code, not by the inventory domain module.

## Public Contract Impact

No semantic command, query, event, error, or permission contract is added or removed.

This change adds Protobuf wire schemas aligned with existing registered contracts:

- `GrantItem`
- `GetInventory`
- `ItemGranted`
- `inventory_errors`
- `inventory_permissions`

The wire schema is compatibility-sensitive once public clients exist.

## Data And Migration Impact

No data ownership, persistence, PostgreSQL schema, or migration impact.

## Test Impact

No Go runtime tests are added because this change does not implement runtime behavior.

Tooling verification is updated so `node tools/vibit check protocol` validates:

- Buf config.
- Envelope source shape.
- Protocol enum values.
- Go package options.
- Inventory message and field alignment.
- Source trace markers.

## Documentation Impact

English and Simplified Chinese public docs are updated where the protocol source status changed from planned to created:

- Repository README.
- Agent operating guide.
- Architecture README.
- Game protocol standard.
- Protobuf README.

An ADR and conversation log are added.

## Compatibility Risks

No existing public client compatibility risk because no public clients exist yet.

The main future risk is Protobuf compatibility. Field numbers, package names, and `go_package` options should now be treated as intentional and changed only through a change spec and ADR.
