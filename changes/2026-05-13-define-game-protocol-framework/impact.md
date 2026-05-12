# Impact

## Affected Modules

The `inventory` module remains the first proof slice.

This change defines the game protocol framework that future inventory WebSocket/Protobuf work must use. It does not change inventory's semantic command, query, event, error, permission, or module ownership contracts.

## Module Ownership Impact

No data or module ownership moves.

The change strengthens runtime boundaries:

- WebSocket transport belongs to `runtime/internal/platform/transport/ws/`.
- Protobuf envelope handling belongs to `runtime/internal/platform/protocol/protobuf/`.
- Application dispatch belongs to `runtime/internal/app/`.
- Inventory domain behavior belongs to `runtime/internal/modules/inventory/`.
- Protobuf sources belong to `proto/vibit/inventory/v1/`.

## Public Contract Impact

No public command, query, event, error, or permission contracts change.

The new standard affects future wire protocol contracts by defining envelope, routing, session, target, authority, and error expectations before `.proto` files are created.

## Event Impact

No event contracts change.

The protocol standard clarifies that domain events pushed to clients are server facts and clients must not publish domain facts directly.

## Permission Impact

No permission contracts change.

Future protocol work must still map permission failures into registered error catalogs.

## Data And Migration Impact

No data migrations are added.

No persistence model changes are introduced.

## Test Impact

No runtime tests are added because no Go runtime protocol code is added.

Verification is limited to documentation, manifests, CLI syntax, protocol manifest checks, and existing repository checks.

## Documentation Impact

This change adds a protocol standard, its Simplified Chinese translation, a machine-readable protocol manifest, an ADR, and conversation memory. It also updates repository guides so future agents can discover the protocol standard.

## Compatibility Risks

The main risk is over-shaping the protocol before the first runtime handler exists.

The risk is bounded by:

- Keeping `/v1/ws` planned until transport implementation begins.
- Reserving room, match, stream, input, state, heartbeat, and reconnect behavior instead of implementing it now.
- Avoiding `.proto` files and generated output in this change.
- Recording reversal conditions in ADR-0015.
