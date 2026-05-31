# Impact

## Architecture

- Adds the friends relationship PostgreSQL adapter implementation under `runtime/internal/platform/persistence/postgres`.
- Keeps SQL mapping in the platform persistence package.
- Preserves `runtime/internal/modules/friends` as the owner of storage-neutral repository vocabulary and error types.
- Adds unit-of-work repository handoff through the existing PostgreSQL unit-of-work shape.
- Opens the next bounded work item for a runtime behavior gate.

## Runtime

The adapter implementation adds no friends relationship runtime handlers, services, startup wiring, route policy, WebSocket behavior, or protocol dispatch.

The only runtime Go behavior added is platform persistence adapter code and its unit-of-work repository factory helper.

## Protocol And Generated Output

No protocol route, Protobuf source, generated output, command, query, event, permission, or public error contract is added.

## Data

No migration is added or changed. The adapter targets the existing `friend_relationships` table.

## Operations

No hosted deployment, release artifact, install script, registry publication, public announcement, paid promotion, or operations/admin behavior is added.

## Risks

- Default verification uses fake-executor adapter tests, not live PostgreSQL.
- Runtime friends relationship behavior remains deferred, so no public route can use the adapter yet.
- Future runtime behavior must still define actor identity derivation, actor-relative public status, permission, conflict, and route-policy rules before use.
