# Request

## Original Request

The maintainer selected:

```text
wire_runtime_authentication_startup_composition
```

and asked to continue while strongly referencing Nakama and Pitaya.

## Clarified Requirement

Implement the bounded startup composition slice authorized by `ADR-0054`:

- Wire existing authentication service validation into the PostgreSQL runtime startup path.
- Load verifier keys from explicit environment lookup.
- Compose a route access-token validator and route protector.
- Inject route protection into the Protobuf frame handler.
- Keep the memory path as a metadata-only bootstrap path.

## User-Visible Outcome

When `VIBIT_RUNTIME_STORE=postgres`, process startup now requires authentication verifier key configuration and composes request-level route protection into the WebSocket Protobuf request loop.

## Non-Goals

- Do not add WebSocket handshake authentication.
- Do not add session persistence.
- Do not add authentication command routes or new Protobuf messages.
- Do not change repository interfaces, PostgreSQL adapters, migrations, generated files, or dependencies.
- Do not add logout, refresh, cleanup, token rotation, token validation audit mutation, or broader production authentication behavior.

## Acceptance Criteria

- [x] PostgreSQL startup composes the existing authentication service.
- [x] Verifier keys are loaded through the existing environment loader.
- [x] Missing or invalid verifier key configuration fails startup closed.
- [x] Token lifetime and audience use defaults with optional environment overrides.
- [x] Route protector is injected into the Protobuf frame handler.
- [x] Memory startup remains available without route protector injection.
- [x] WebSocket transport remains credential-neutral.
- [x] Existing Protobuf envelope remains unchanged.
- [x] Focused tests pass.
- [x] Verification is recorded.
