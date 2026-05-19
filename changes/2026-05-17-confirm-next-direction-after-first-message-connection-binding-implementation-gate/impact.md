# Impact Analysis

## Affected Modules

- Runtime workflow and manifests.
- Authentication module guidance, because the selected route belongs to `runtime.authentication` while implementation ownership remains outside the storage-neutral authentication module.

## Module Ownership Impact

This direction confirmation does not change runtime ownership by itself. It authorizes the next bounded implementation milestone to work in the application, protocol adapter, generated Protobuf, startup composition, and WebSocket metadata handoff boundaries defined by `ADR-0058`.

## Public Contract Impact

No public contract or wire schema is changed by this confirmation step. The next implementation slice may add `BindConnectionRequest`, `BindConnectionResponse`, and `ConnectionBindingStatus`.

## Data And Migration Impact

No migrations, repositories, session tables, or persistent connection binding state are selected.

## Test Impact

The implementation slice must add focused app, protocol adapter, startup, generated output, and WebSocket transport neutrality tests.

## Documentation Impact

The work queue, manifests, AGENTS guides, rules, tools, change spec, and conversation log must record the selected direction and the later implementation result.

## Compatibility Risks

The direction is intentionally narrow. It does not change WebSocket handshake authentication, the existing Protobuf envelope, request-level route protection, session persistence, logout/revocation, reconnect/epoch behavior, or direct Nakama/Pitaya compatibility.
