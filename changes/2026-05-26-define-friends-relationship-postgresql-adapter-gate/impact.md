# Impact

## Added

- `docs/friends-relationship-postgresql-adapter-gate.md`
- `docs/friends-relationship-postgresql-adapter-gate.zh-CN.md`
- `ADR-0144`
- `runtime.friends_relationship_postgresql_adapter_gate`
- `M-165/W-0237 Implement friends relationship PostgreSQL adapter` as the next-ready follow-up

## Behavioral Impact

No runtime behavior is added.

The new standard defines future adapter ownership, SQL mapping posture, unit-of-work handoff, conflict mapping, redaction, and test expectations only. It does not register handlers, expose routes, execute SQL, parse transport credentials, validate sessions, or change startup composition.

## Data Impact

No migration is added or changed.

The existing `runtime/migrations/postgres/000007_create_friend_relationships.sql` remains the source migration for future adapter work, but this slice does not execute it or map SQL rows.

## Protocol Impact

No protocol route, Protobuf source, generated output, bridge, route key, or client SDK is added.

## Security And Privacy Impact

The gate treats social graph data as not log-safe by default. Future adapter errors must collapse driver, constraint, player id, private relationship state, SQL, DSN, credential, token, and digest details into redacted friends module errors.

## Reference Alignment

Nakama guides the product capability need for durable friends relationship social graph state. vibit adapts that need into a PostgreSQL adapter gate before implementation. Pitaya remains deferred as a future distributed architecture reference.

