# Impact Analysis

## Affected Modules

- Runtime application layer: adds `ConnectionBinder`.
- Runtime protocol adapter: adds `BindConnection` system-route handling.
- Runtime WebSocket transport: carries server-observed `connection_epoch` metadata without parsing credentials.
- Runtime startup: injects binder in the PostgreSQL authentication composition.
- Authentication module manifest and guide: records that storage-neutral authentication does not own connection binding.

## Module Ownership Impact

Connection binding is application-owned under `runtime/internal/app`. The authentication module remains storage-neutral and still owns credential/token record vocabulary only. The Protobuf adapter owns wire conversion, not validation logic. The WebSocket transport still owns only connection/frame metadata and opaque frame handoff.

## Public Contract Impact

Adds authentication Protobuf wire messages:

- `vibit.authentication.v1.BindConnectionRequest`
- `vibit.authentication.v1.BindConnectionResponse`
- `vibit.authentication.v1.ConnectionBindingStatus`

Adds binding-specific public error codes in application vocabulary:

- `CONNECTION_BINDING_TOKEN_MISSING`
- `CONNECTION_BINDING_TOKEN_MALFORMED`
- `CONNECTION_BINDING_TOKEN_INVALID`
- `CONNECTION_BINDING_UNAVAILABLE`
- `CONNECTION_BINDING_REQUIRED`

No semantic command/query/event/permission contracts are changed in this slice.

## Data And Migration Impact

No durable session or connection-binding state is added. No migration, repository interface, or PostgreSQL adapter changes are included.

## Test Impact

Adds focused tests for:

- Application binder success, failure collapse, metadata-only rejection, server-observed connection metadata, and redaction.
- Protobuf adapter success and error envelope mapping.
- Envelope stability.
- Startup binder injection.
- WebSocket connection epoch metadata handoff and credential-neutral handshake behavior.

## Documentation Impact

Updates change specs, conversation log, work queue, manifests, AGENTS guides, rules, and repository checks.

## Compatibility Risks

Existing command/query request-level route protection remains unchanged. `BindConnection` does not make bound identity satisfy normal protected routes yet. Existing clients that do not use the new system route continue to use the request-level `AuthenticatedRequest` wrapper for protected routes.
