# Impact

## Added

- `runtime/internal/modules/friends/repository.go`
- `runtime/internal/modules/friends/repository_test.go`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`
- Standard module subdirectories under `modules/friends/`
- `ADR-0143`
- `runtime.friends_relationship_repository_interface_implementation`

## Behavioral Impact

No runtime behavior is added.

The new Go package defines storage-neutral repository vocabulary and validation helpers only. It does not register handlers, expose routes, execute SQL, parse transport credentials, validate sessions, or change startup composition.

## Data Impact

No migration is added or changed.

The existing `runtime/migrations/postgres/000007_create_friend_relationships.sql` remains the source migration for future adapter work, but this slice does not execute it or map SQL rows.

## Protocol Impact

No protocol route, Protobuf source, generated output, bridge, route key, or client SDK is added.

## Security And Privacy Impact

The repository vocabulary treats social graph data as not log-safe by default. Errors are typed and redacted. Actor ids and pair member ids remain data, not authentication proof.

## Reference Alignment

Nakama guides the product capability need for friends relationship social graph state. vibit adapts that need into module-owned, storage-neutral repository vocabulary before adapter and runtime behavior. Pitaya remains deferred as a future distributed architecture reference.
