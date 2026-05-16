# Impact

This change implements the first real login behavior inside the application authentication service.

## Affected Areas

- `runtime/internal/app/authentication/service.go` now executes `AuthenticateWithDeviceCredential`.
- `runtime/internal/app/authentication/service_test.go` now covers the service behavior slice.
- Architecture manifests record `W-0109` and `M-037` as completed and open `W-0110`.
- Repository checks gain `runtime.device_credential_login_service_behavior_implementation`.
- Agent guides record the implemented login slice and remaining deferrals.

## Unaffected Areas

- No public command, query, event, error, or permission contract source changes.
- No Protobuf or WebSocket proof carrier changes.
- No access-token validation, logout, refresh, cleanup, session persistence, or route protection.
- No repository interface changes.
- No PostgreSQL adapter changes.
- No migrations.
- No generated files.
- No external dependencies.
- No startup wiring.

## Compatibility

The Go service constructor now requires the dependencies needed by login behavior. There are no external runtime callers yet, and existing non-login service methods remain fail-closed.
