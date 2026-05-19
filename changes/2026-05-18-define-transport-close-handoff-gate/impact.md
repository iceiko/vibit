# Impact

## Affected Modules

- `runtime`
- `authentication`
- `reference`

## Runtime Impact

This is a gate-only standard. It adds no Go runtime behavior.

No WebSocket sockets are closed by this change. No route behavior, protocol behavior, authentication service behavior, session behavior, registry behavior, or generated output changes are made.

## Ownership Impact

The gate records:

- Application close policy remains under `runtime/internal/app/connection`.
- Active connection registry remains under `runtime/internal/app/connection`.
- Future concrete close mechanics may live under `runtime/internal/platform/transport/ws`.
- Protocol adapters must not close sockets directly.
- Authentication service must not call transport.
- Domain modules must not import WebSocket transport.

## Product Roadmap Impact

This keeps the Nakama/Pitaya-class roadmap on lifecycle closure:

```text
protocol logout route -> transport close handoff -> reconnect/epoch -> protocol session carrier -> presence
```

## Compatibility Impact

No API, Protobuf, data, or runtime behavior compatibility change is made.

Direct Nakama/Pitaya API compatibility remains out of scope.
