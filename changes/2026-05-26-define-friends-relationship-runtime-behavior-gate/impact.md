# Impact

## Added

- `docs/friends-relationship-runtime-behavior-gate.md`
- `docs/friends-relationship-runtime-behavior-gate.zh-CN.md`
- `ADR-0146`
- `runtime.friends_relationship_runtime_behavior_gate`
- `M-167/W-0239 Implement friends relationship runtime behavior` as the next-ready follow-up

## Behavioral Impact

No runtime behavior is added.

The new standard defines future application service ownership, validated actor identity derivation, actor-relative public status behavior, permission and route-policy posture, unit-of-work handoff, conflict mapping, redaction, and test expectations only. It does not register handlers, expose routes, execute application behavior, parse transport credentials, validate sessions, or change startup composition.

## Data Impact

No migration is added or changed.

The existing `friend_relationships` table and implemented PostgreSQL adapter remain the persistence source for future runtime behavior, but this slice does not change repository interfaces, SQL mapping, or adapter behavior.

## Protocol Impact

No protocol route, Protobuf source, generated output, bridge, route key, or client SDK is added.

## Security And Privacy Impact

The gate refuses metadata-only `player_id` and `session_id` as proof, requires actor identity to come from validated request identity, and treats private social graph data as not log-safe by default. Future runtime errors must collapse hidden graph details, storage errors, actor ids, target ids, relationship ids, token material, verifier digests, transport metadata, and SQL details into redacted application classes.

## Reference Alignment

Nakama guides the product capability need for durable friends relationship behavior. vibit adapts that need into an application-owned runtime behavior gate before implementation. Pitaya remains deferred as a future distributed architecture reference.

