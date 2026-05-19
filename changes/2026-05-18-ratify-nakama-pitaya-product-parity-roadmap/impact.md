# Impact Analysis

## Affected Modules

- `runtime`: the near-term execution order changes to lifecycle closure before feature module expansion.
- `reference`: Nakama/Pitaya alignment becomes a product parity roadmap rather than background reference only.

## Module Ownership Impact

No module ownership changes are implemented. Future features still need their own module or runtime owner decisions.

## Public Contract Impact

No command, query, event, error, permission, Protobuf, or generated contract shape changes are made in this slice.

## Data And Migration Impact

No data model, repository, migration, storage adapter, or durable session behavior changes are made.

## Test Impact

No Go tests are required because this is a roadmap and check-rule change. Repository checks are updated to verify the new roadmap markers.

## Documentation Impact

Adds a new English roadmap standard and Simplified Chinese translation, plus ADR and conversation memory.

## Compatibility Risks

No runtime compatibility risk. The roadmap explicitly rejects accidental direct Nakama/Pitaya API compatibility.
