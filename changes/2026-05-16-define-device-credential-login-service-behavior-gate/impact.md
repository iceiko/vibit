# Impact

This change adds a gate-only standard for future device credential login behavior.

## Affected Areas

- Architecture manifests now record `ADR-0051` and the completed `W-0108` gate.
- Agent guides now point future agents to the device credential login gate before implementing login.
- The authentication module manifest records that future login may use existing repository methods only through application unit-of-work capabilities.
- Repository checks gain a dedicated rule for the gate.
- The next work item becomes `W-0109`, a bounded service-only implementation slice.

## Unaffected Areas

- No Go authentication service behavior is added.
- No service method signatures are changed.
- No repository interfaces, PostgreSQL adapters, migrations, Protobuf sources, generated files, WebSocket carriers, startup wiring, or external dependencies are changed.
- Access-token validation, logout, refresh, cleanup, protocol carriers, and broader production authentication behavior remain deferred.
