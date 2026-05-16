# Request

## Original Request

The maintainer said "继续" under the established work continuation standard.

## Clarified Requirement

Advance one `next_ready` work item. The current work item is `W-0106 Define authentication service behavior implementation gate`.

This change must define the implementation gate for future authentication service behavior without adding service code.

## User-Visible Outcome

The repository gains a durable standard, ADR, manifest markers, agent guidance, and a check rule for future authentication service behavior implementation.

## Non-Goals

- Do not implement authentication service behavior.
- Do not execute device credential login.
- Do not validate access tokens.
- Do not add logout, refresh, or cleanup behavior.
- Do not add Protobuf authentication messages.
- Do not add WebSocket proof carriers.
- Do not change public contracts.
- Do not change repository interfaces or PostgreSQL adapters.
- Do not change migrations.
- Do not wire startup.
- Do not add authentication dependencies.
- Do not add production authentication behavior.

## Unknowns

- None for this gate-only change.

## Acceptance Criteria

- The gate names future service behavior ownership and files.
- The gate defines repository handoff through the application unit-of-work boundary.
- The gate maps existing helper slices into future login and access-token validation composition flows.
- The gate defines public error collapse and redaction posture.
- The gate defines focused tests for future service behavior.
- The gate preserves all login, token validation, logout, refresh, cleanup, protocol, repository, migration, startup, dependency, and production behavior deferrals.
- No service code is added.
