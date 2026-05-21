# Impact Analysis

## Affected Modules

- `runtime`
- `workflow`
- `reference`

## Module Ownership Impact

No runtime ownership changes are introduced.

The runbook documents current ownership:

- WebSocket transport remains credential-neutral.
- Protobuf adaptation remains under `runtime/internal/platform/protocol/protobuf`.
- Application authentication and onboarding behavior remains under `runtime/internal/app/authentication`.
- Local onboarding remains application-service-only, not protocol-owned.

## Public Contract Impact

No commands, queries, events, errors, permissions, Protobuf sources, generated output, or wire envelope fields change.

The runbook now explicitly names existing protocol routes and the existing `AuthenticatedRequest` carrier for protected routes.

## Data And Migration Impact

No migrations or repository interfaces change.

The runbook reiterates that normal server startup does not apply migrations automatically and that PostgreSQL runtime operation requires explicit migration preparation.

## Test Impact

No tests are added or changed.

The runbook points developers to the existing focused E2E proof:

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
```

## Documentation Impact

Updated:

- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`

Added:

- `decisions/ADR-0093-runtime-runbook-alpha-path-refresh.md`
- `conversations/2026-05-21-runtime-runbook-alpha-path-refresh.md`
- `changes/2026-05-21-refresh-runtime-runbook-for-alpha/`

## Compatibility Risks

There is no API, data, wire, migration, dependency, or release compatibility risk because this is a documentation and check-rule slice.

The main risk is overclaiming public onboarding availability. The runbook explicitly states that local onboarding exists as an application service and is not yet a public protocol, HTTP, CLI, or startup auto-creation surface.
