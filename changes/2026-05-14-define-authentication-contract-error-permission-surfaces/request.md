# Request

Define the authentication contract, error, permission, and audit surfaces required by M-013 before implementing the selected first login and token posture.

The work must preserve current runtime code, Protobuf behavior, WebSocket behavior, generated output, schemas, migrations, and implementation deferral.

## Maintainer Context

The maintainer asked the agent to continue advancing the queue according to professional judgment unless a true decision boundary required confirmation.

## Scope

- Define semantic contract sources for selected authentication and token surfaces.
- Define error and permission catalogs.
- Define audit-oriented event surfaces.
- Update machine-readable manifests and inspection tooling.
- Preserve implementation deferral.

## Out Of Scope

- Runtime authentication code.
- Login handlers.
- Token generation, parsing, validation, refresh, revocation, rotation, replay handling, or storage.
- Credential, token, external identity, or session tables.
- Migrations.
- Protobuf messages or generated Protobuf output.
- Generated Go contract shape files.
- WebSocket handshake authentication.
- Runtime player handlers.
- WebSocket routes.
- Major dependencies.
