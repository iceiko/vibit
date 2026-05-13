# Request

## Original Request

```text
继续
```

## Clarified Requirement

Continue the first Go runtime path by adding a narrow application dispatch skeleton after the protocol handoff slice.

This slice should:

- Add pure application-layer command/query dispatch types under `runtime/internal/app/`.
- Keep route registration explicit and test-covered.
- Preserve request correlation metadata when a route is dispatched.
- Return machine-readable application errors for invalid routes, duplicate routes, unknown routes, and nil handlers.
- Strengthen runtime boundary checks so application and domain packages do not depend backward on platform adapters or generated Protobuf packages.

## User-Visible Outcome

Future agents should have a stable, small, platform-free dispatch entry point between protocol adaptation and domain module handlers.

The repository should make it clear that WebSocket transport, Protobuf wire conversion, application dispatch, domain logic, and generated output are separate layers.

## Non-Goals

- Do not implement WebSocket transport.
- Do not implement PostgreSQL repositories, migrations, or transaction wiring.
- Do not implement inventory business behavior.
- Do not introduce authentication or session validation.
- Do not change public command, query, event, error, or permission contracts.
- Do not generate route registration yet.

## Unknowns

- Generated route registration remains deferred until the generator can use contract and module manifests as sources.
- The first inventory runtime handler shape remains deferred until repository and policy interfaces are declared.
- Transaction boundaries remain deferred until persistence and unit-of-work interfaces are introduced.

## Acceptance Criteria

- [ ] `runtime/internal/app/` exposes a dispatcher with explicit route registration.
- [ ] Dispatch supports command and query route requests and rejects unsupported message kinds.
- [ ] Duplicate route registration, invalid routes, unknown routes, and nil handlers produce stable application errors.
- [ ] Dispatch preserves request, route, target, and session metadata in application results.
- [ ] App-layer code does not import WebSocket, Protobuf, platform adapters, or generated Protobuf packages.
- [ ] Runtime checks include a layer-boundary check for app/domain package imports.
- [ ] Go runtime tests and `node tools/vibit check all --json` pass.
