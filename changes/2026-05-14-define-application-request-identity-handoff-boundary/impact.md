# Impact

## Affected Modules

- `runtime/internal/app`
- `runtime/internal/platform/protocol/protobuf`
- `runtime/internal/app/bootstrap`
- `runtime/internal/modules/inventory`
- `.arch/runtime.yaml`
- `.arch/protocol.yaml`
- `.arch/conventions.yaml`
- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `docs/player-identity-session-boundary.md`
- `docs/player-identity-session-boundary.zh-CN.md`

The `player` module manifest is not changed by this step.

## Module Ownership Impact

Application dispatch now owns the request identity handoff shape. This follows `docs/player-identity-session-boundary.md`: transport connection metadata, Protobuf session metadata, application request identity, player identity, and future authentication remain separate concerns.

Inventory remains the owner of inventory state only. It does not gain player account, authentication, session validation, or player module dependencies.

## Public Contract Impact

No public commands, queries, events, errors, permissions, Protobuf messages, or WebSocket routes are added or changed.

## Runtime Impact

`app.RouteRequest` now carries `RequestIdentity`. Current session fields from the envelope are converted to `metadata_only` identity. Existing `player_id` and `session_id` values are normalized but not validated.

`app.ApplicationResult` now carries `RequestIdentity` so downstream protocol and future permission work can preserve application identity context.

The runtime, protocol, and conventions manifests now record that application request identity handoff exists but real authentication and real session validation do not.

## Protocol Impact

The existing Protobuf envelope is unchanged. The WebSocket handshake is unchanged. The Protobuf adapter only maps existing envelope session metadata into the new application-owned handoff type.

## Data And Migration Impact

No database schema, migration, credential store, token store, or session store is added.

## Compatibility Risk

This is an internal Go runtime type change. The client wire format is unchanged. Existing Go code that constructs `RouteRequest` can omit identity and rely on application dispatch to populate metadata-only identity.

## Test Impact

Focused Go tests cover:

- Metadata-only identity normalization.
- Empty-player metadata behavior.
- Validated identity helper behavior.
- Dispatch-time identity defaulting and result propagation.
- Envelope-to-route identity conversion.
- Frame connection metadata refresh for metadata-only identity.
