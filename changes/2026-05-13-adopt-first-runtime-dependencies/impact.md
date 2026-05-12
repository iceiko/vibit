# Impact

## Affected Modules

No domain module ownership changes.

The first inventory proof slice remains the preferred first runtime capability. This change affects platform and tooling boundaries that the inventory runtime will use later.

## Module Ownership Impact

Domain modules remain forbidden from directly importing transport, Protobuf, PostgreSQL, migration, S3, MinIO, observability, or framework dependencies.

Accepted dependencies are owned by platform adapters or generation tooling:

- WebSocket library: platform transport adapter.
- Protobuf runtime and generator: generation tooling and generated protocol packages.
- PostgreSQL driver: platform persistence adapter.
- Migration tooling: platform migration tooling.

## Public Contract Impact

No command, query, event, error, or permission contract changes.

The semantic contract source remains vibit manifest YAML. Protobuf remains the wire schema format and must align with semantic contracts.

## Event Impact

No event contract changes.

## Permission Impact

No permission changes.

## Data And Migration Impact

No database migrations are added.

The change accepts PostgreSQL driver and migration tooling choices for future implementation, but no runtime database schema is introduced yet.

## Test Impact

No runtime tests are added because server runtime code has not started.

The future runtime test path should use Go standard-library `testing` first. External test framework adoption remains deferred until a concrete need appears.

## Documentation Impact

This change updates architecture manifests, AGENTS guides, conversation memory, and ADRs.

## Compatibility Risks

The main risk is premature dependency lock-in. The risk is bounded by:

- Keeping imports behind platform adapters and generation tooling.
- Recording replacement paths.
- Leaving unneeded categories deferred.
- Avoiding Go implementation code in this change.
