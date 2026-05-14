# Impact

## Affected Modules

- `runtime`
- `player`
- `inventory`

The change affects repository checks and standards. It does not add runtime authentication behavior.

## Module Ownership Impact

WebSocket transport remains the owner of connection mechanics only. It must not import domain modules, player/inventory runtime packages, generated Protobuf packages, or Protobuf runtime dependencies.

Domain modules remain owners of module business behavior only. They must not import WebSocket transport, generated Protobuf, Protobuf runtime, authentication providers, token libraries, credential stores, or password-hashing dependencies.

The player module remains a boundary-only manifest owner until public player contracts, account schemas, persistence, and authentication decisions are ratified.

## Public Contract Impact

No public commands, queries, events, errors, permissions, Protobuf messages, or WebSocket routes are added or changed.

## Runtime Impact

`node tools/vibit check runtime` now includes a `runtime.identity_boundary` rule. The rule inspects existing Go source files, dependency declarations, player manifest markers, player Protobuf source absence, and player migration absence.

## Protocol Impact

No Protobuf envelope fields change. No WebSocket handshake behavior changes.

## Data And Migration Impact

No database schema or migration is added. The new check fails if player/account persistence appears under PostgreSQL migrations before ratification.

## Test Impact

No Go runtime behavior changes, but `check runtime` still runs `go test ./...`. The new static checks are covered through repository verification.

## Documentation Impact

The player identity/session boundary standard and runtime architecture manifest record the new repository check and the intentionally deferred checks.

## Compatibility Risks

The check is intentionally conservative. It may need refinement when player public contracts are ratified, because today it treats player runtime implementation, player Protobuf sources, and player/account migrations as premature.

To avoid false authority, the check does not try to prove semantic authentication correctness. It blocks the common structural mistakes that would let metadata-only identity become de facto authentication.
