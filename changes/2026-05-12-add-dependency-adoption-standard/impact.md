# Impact

## Architecture Impact

This change adds a dependency adoption standard and `.arch/dependencies.yaml` registry.

It prepares upcoming dependency choices without choosing any dependency. The registry records slots for WebSocket, Protobuf tooling, PostgreSQL driver, migration tooling, S3 client, MinIO server, test framework, and observability.

## Module Impact

No module ownership changes.

Domain modules remain forbidden from importing foundational transport, protocol, persistence, object-storage, or framework dependencies directly.

## Public Contract Impact

No command, query, event, error, or permission contracts change.

## Data And Migration Impact

No data model or migration is introduced.

This change prepares the dependency adoption path for future PostgreSQL driver and migration tooling decisions.

## Dependency Impact

No dependency is added.

No dependency slot is accepted. All slots remain `proposed` or `deferred`.

## Tooling Impact

`tools/vibit check architecture` now checks that the dependency adoption registry, standard, and template exist and are referenced.

## Documentation Impact

New public English documentation has a paired Simplified Chinese translation.

README, AGENTS, and architecture README files are updated to point future agents to `.arch/dependencies.yaml`.

## Compatibility

No public runtime API exists yet, so there is no compatibility break.

This change is intentionally conservative: it creates the evaluation path for future branch points without taking a branch.
