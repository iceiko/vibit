# Impact

## Affected Modules

- `player`: receives its first Protobuf wire schema and generated Go Protobuf output for existing semantic account contracts.
- `runtime`: generated Protobuf output grows, but no handwritten runtime handler is added.

## Module Ownership Impact

No ownership moves.

The player module continues to own player identity and account lifecycle semantics. Protobuf source defines wire message shape only. Runtime protocol adapter packages may import generated Protobuf packages later, but domain modules must not import generated Protobuf directly.

## Public Contract Impact

No new semantic command, query, event, error, or permission is added.

This change adds wire messages aligned with existing player account semantic contracts:

- `CreatePlayerAccountRequest`
- `CreatePlayerAccountResponse`
- `GetPlayerAccountRequest`
- `GetPlayerAccountResponse`
- `PlayerAccountCreated`

## Data And Migration Impact

No database schema, migration, repository, credential store, token store, or session store is added.

## Protocol Impact

Adds a player module Protobuf source package:

```text
vibit.player.v1
```

The envelope is unchanged. The WebSocket endpoint and handshake are unchanged. No route is mounted.

## Generated Output Impact

Adds generated Go Protobuf output under:

```text
runtime/internal/generated/proto/vibit/player/v1/
```

Generated output must be produced by `buf generate` and must not be hand-edited.

## Test Impact

No runtime behavior tests are added because no handwritten runtime behavior changes.

Repository checks must validate protocol alignment, generated-output traceability, runtime boundaries, and the work queue.

## Compatibility Risks

Low to moderate. This creates an initial wire schema for player account contracts, so later changes to these messages become compatibility-sensitive.

The main architecture risks are premature runtime handler implementation, authentication decisions, or persistence schema decisions. This change explicitly keeps those deferred.
