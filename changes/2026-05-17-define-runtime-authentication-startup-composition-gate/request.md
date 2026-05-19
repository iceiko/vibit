# Request

## Original Request

The maintainer selected:

```text
wire_runtime_authentication_startup_composition
```

and asked to continue while strongly referencing Nakama and Pitaya.

## Clarified Requirement

Define a gate for runtime authentication startup composition before implementation. The gate must specify how the existing authentication service, environment verifier key loader, PostgreSQL unit-of-work runner, route access-token validator, route protector, and Protobuf frame handler may be composed at process startup.

## User-Visible Outcome

The repository has a written standard and ADR that authorize a bounded implementation slice in `runtime/cmd/vibit-server`.

## Non-Goals

- Do not implement startup composition in the gate-only step.
- Do not add WebSocket handshake authentication.
- Do not add session persistence.
- Do not add authentication command Protobuf messages or login route registration.
- Do not change repository interfaces, PostgreSQL adapters, migrations, generated files, or dependencies.
- Do not add logout, refresh, cleanup, token rotation, token validation audit mutation, or broader production authentication behavior.

## Acceptance Criteria

- [x] English and Simplified Chinese startup composition gate docs exist.
- [x] `ADR-0054` records the startup composition decision.
- [x] Manifests record the gate, selected PostgreSQL-only first path, default lifetime, default audience, and deferrals.
- [x] Repository check rule `runtime.authentication_startup_composition_gate` exists.
- [x] `W-0116` is completed and `W-0117` is available for implementation.
- [x] Verification is recorded.
