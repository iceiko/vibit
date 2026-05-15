# Request

## Original Request

The maintainer said "继续" under the established work continuation standard.

## Clarified Requirement

Advance one `next_ready` work item. The current work item is `W-0104 Define verifier digest comparison helper gate`.

This change must define the implementation gate for future constant-time verifier digest comparison helpers without adding comparison code.

## User-Visible Outcome

The repository gains a durable standard, ADR, manifest markers, agent guidance, and a check rule for future verifier digest comparison helper implementation.

## Non-Goals

- Do not implement verifier comparison code.
- Do not implement authentication service behavior.
- Do not execute login.
- Do not validate access tokens.
- Do not add logout, refresh, or cleanup behavior.
- Do not add Protobuf authentication messages.
- Do not add WebSocket proof carriers.
- Do not change repositories.
- Do not change migrations.
- Do not wire startup.
- Do not add authentication dependencies.
- Do not add production authentication behavior.

## Unknowns

- None for this gate-only change.

## Acceptance Criteria

- The gate names future helper ownership and files.
- The gate requires constant-time verifier digest comparison.
- The gate rejects raw material, lookup digest, key id, account id, and protocol metadata comparison.
- The gate defines redacted failure posture.
- The gate defines focused tests for the future helper.
- The gate preserves all service, protocol, repository, migration, startup, dependency, and production behavior deferrals.
- No verifier comparison code is added.
