# Impact

## Module Impact

`inventory` remains the preferred first runtime proof slice, but its runtime implementation status returns to `not_started`.

The existing semantic contracts remain in place:

- `GrantItem`
- `GetInventory`
- `ItemGranted`
- `inventory_errors`
- `inventory_permissions`

Their implementation metadata now points toward planned Go runtime extension points and Protobuf wire schema outputs rather than removed TypeScript files.

## Runtime Impact

Go is now the ratified first server runtime implementation language.

The previous TypeScript runtime slice and npm package baseline are removed from the mainline so future agents do not mistake them for approved server architecture.

`tools/vibit` remains a Node.js standard-library CLI for architecture checks and inspection. It is explicitly not the server runtime.

## Protocol Impact

WebSocket is the first gameplay/client protocol.

Protobuf is the first client/server wire message format.

vibit manifests and contract files remain the business semantics source of truth. Protobuf files will own wire message shape and compatibility after a formal proto layout and tooling decision.

## Dependency Impact

No new external dependency is added in this change.

Future WebSocket libraries, Protobuf tooling, and foundational runtime dependencies require an explicit adoption record. Popularity is useful evidence but not enough by itself.

## Compatibility

No public runtime API exists yet, so there is no external API, event, or data compatibility break.

This is a repository-direction correction before Go runtime implementation starts.

## Risk

The main risk is historical confusion because previous change specs still record the removed TypeScript work. This is intentional: historical records are preserved, while current ADRs and manifests supersede the earlier direction.

Future agents must follow `ADR-0008`, `ADR-0009`, `ADR-0010`, and `.arch/runtime.yaml` rather than treating earlier TypeScript runtime changes as current direction.
