# Request

## Original Request

Continue advancing the project work queue. The active next-ready work item is `W-0028 Add metadata-only identity repository checks`.

## Clarified Requirement

Add repository checks that protect the current metadata-only player/session boundary from common regression paths.

The checks must make it hard for a future agent to accidentally turn client-supplied metadata into authentication by placing player, inventory, authentication, token, credential, generated Protobuf, or transport behavior in the wrong runtime layer.

## User-Visible Outcome

`node tools/vibit check runtime` now verifies that the metadata-only identity boundary remains explicit. The project can continue adding player/session work without relying on maintainer memory to catch transport/domain/authentication shortcuts.

## Non-Goals

- No real authentication implementation.
- No token format selection.
- No credential, password, OAuth, OIDC, JWT, or session-store dependency.
- No player account database schema or migrations.
- No Protobuf envelope change.
- No WebSocket handshake change.
- No broad static analysis dependency.

## Unknowns

- The future production authentication model remains undecided.
- The future session persistence model remains undecided.
- The future player public contracts remain unratified.
- Some semantic checks remain intentionally deferred until contract schemas and generators become richer.

## Acceptance Criteria

- Runtime checks fail if WebSocket transport imports domain modules, player/inventory runtime packages, generated Protobuf packages, or Protobuf runtime dependencies.
- Runtime checks fail if domain modules import WebSocket, generated Protobuf, Protobuf runtime, known authentication, token, credential, or password-hashing dependencies.
- Runtime checks fail if player runtime code, player Protobuf sources, or player/account PostgreSQL migrations appear before public contracts and schemas are ratified.
- Runtime checks verify that `modules/player/module.yaml` still declares boundary-only public API markers.
- New rule metadata is registered in `rules/check-rules.json`.
- Deferred checks that would currently be too brittle are recorded.
