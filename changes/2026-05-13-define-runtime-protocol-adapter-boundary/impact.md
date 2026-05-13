# Impact

## Architecture Impact

This change defines the runtime handoff between WebSocket transport, Protobuf protocol adaptation, application dispatch, generated code, and domain modules before Go runtime implementation starts.

It strengthens the existing ADR-0014 and ADR-0015 package boundaries without changing the accepted runtime language, protocol, wire format, or persistence direction.

## Contract Impact

No command, query, event, error, or permission contracts change.

## Runtime Impact

No Go source files are added.

Future runtime implementation must keep:

- WebSocket frame handling in `runtime/internal/platform/transport/ws/`.
- Protobuf envelope conversion in `runtime/internal/platform/protocol/protobuf/`.
- Command/query dispatch in `runtime/internal/app/`.
- Domain behavior in `runtime/internal/modules/<module>/`.
- Generated output under `runtime/internal/generated/`.

## Protocol Impact

No `.proto` files change.

The change clarifies how decoded envelope metadata becomes an application route request and where protocol errors are mapped.

## Documentation Impact

Adds the runtime protocol adapter standard and Chinese translation. Updates architecture manifests, runtime guides, protocol docs, README, rule catalog, and CLI checks.

## Compatibility Impact

No public API, event, or data compatibility impact.

Future runtime code may fail architecture checks if it places behavior in the wrong layer.
