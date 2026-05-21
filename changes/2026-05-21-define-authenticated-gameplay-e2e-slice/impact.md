# Impact

## Affected Areas

- `runtime/internal/platform/protocol/protobuf/`: adds a focused Go E2E proof test that drives the existing Protobuf frame handler.
- `runtime/internal/app/authentication`: used by the E2E test through the existing local onboarding, login, validation, and logout service methods.
- `runtime/internal/modules/inventory`: used by the E2E test through the existing protected inventory routes.
- `runtime/internal/app/presence` and `runtime/internal/app/connection`: used by the E2E test through the existing protected presence query and in-memory active connection registry.
- `.arch/`, `AGENTS*`, `README*`, `docs/v0.1-alpha-goal*`, `tools/vibit`, and `rules/check-rules.json`: updated so future agents can find the completed E2E proof and next alpha work item.

## Behavior

This slice adds no production runtime behavior. It proves that existing behavior composes into one local authenticated gameplay path.

The proof:

- calls local onboarding directly through the application service;
- sends a device credential login frame through the existing protocol route;
- sends a `BindConnection` system frame through the existing binding route;
- sends protected inventory grant/read frames through the existing authenticated request wrapper;
- sends a protected self-presence query through the existing authenticated request wrapper;
- sends logout through the existing service-validated logout route;
- confirms the revoked access token no longer satisfies a protected inventory request.

## Contracts

No public command, query, event, error, permission, Protobuf, generated output, or database contract changes are introduced.

## Compatibility

- No breaking API change.
- No event compatibility change.
- No data compatibility change.
- No Protobuf source or generated output change.
- No migration or dependency change.
- No direct Nakama/Pitaya API compatibility.

## Known Deferrals

- Runtime runbook refresh remains a follow-up work item.
- Minimal example client or request-loop script remains a follow-up work item.
- Health/readiness/version/config surfaces remain deferred.
- Alpha acceptance checklist remains deferred.
- Public signup, production identity providers, password login, account recovery, account merge, multi-device linking, chat/social/matchmaking/match runtime modules, SDKs, and direct compatibility remain out of scope.
- Presence currently proves online bound connection state. Runtime session id linkage inside the presence snapshot remains a separate future decision because ordinary request-token validation and binding do not yet use persisted session identity for route policy.
