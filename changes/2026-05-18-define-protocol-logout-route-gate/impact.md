# Impact Analysis

## Affected Modules

- `runtime`: defines the future route and protocol adapter boundary for logout.
- `authentication`: participates through the existing `LogoutAccessToken` service behavior but does not own protocol route registration, Protobuf generation, socket close, reconnect, or session carrier behavior.

## Module Ownership Impact

No runtime owner changes are implemented. The gate records future ownership:

- Authentication service behavior stays under `runtime/internal/app/authentication`.
- Route handler registration belongs under `runtime/internal/app/bootstrap`.
- Protocol mapping belongs under `runtime/internal/platform/protocol/protobuf`.
- Startup composition belongs under `runtime/cmd/vibit-server`.
- WebSocket transport remains credential-neutral.

## Public Contract Impact

No command, query, event, error, permission, Protobuf, or generated contract shape changes are made in this gate.

## Data And Migration Impact

No data model, repository, migration, durable session, token schema, or adapter behavior changes are made.

## Test Impact

No Go tests are required for this gate-only standard. The future implementation slice must add tests for protocol messages, generated output, route registration, transaction bypass, bridge mapping, redaction, and socket/session deferrals.

## Documentation Impact

Adds an English standard, Simplified Chinese translation, ADR, conversation log, and change spec for the protocol logout route gate.

## Compatibility Risks

No runtime compatibility risk because this change does not alter running behavior. The gate explicitly avoids direct Nakama/Pitaya API compatibility and keeps the existing Protobuf envelope unchanged.
