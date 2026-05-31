# Request

## Original Request

Continue the next ready repository work item.

## Clarified Requirement

Advance `M-167/W-0239 Implement friends relationship runtime behavior`.

Implement only the application-owned friends relationship runtime behavior service selected by `ADR-0146`, using validated request identity and the storage-neutral friends repository through the unit-of-work boundary.

## User-Visible Outcome

Future protocol or handler slices can call a focused application service for friend request, accept, reject, remove, block, unblock, list, and status behavior without owning persistence, route policy, or transport behavior.

## Non-Goals

- No protocol routes.
- No Protobuf source.
- No generated output.
- No dependency additions.
- No migration changes.
- No startup wiring beyond package-local service tests.
- No event/audit table.
- No chat, groups, parties, broadcast fanout, matchmaking, match runtime, SDK, hosted deployment, release artifact, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- The service lives under `runtime/internal/app/friends`.
- The service rejects missing, metadata-only, unvalidated, non-player, or mismatched request identity before repository mutation.
- The service obtains `runtime/internal/modules/friends.Repository` through unit-of-work capabilities.
- Send, accept, reject, remove, block, unblock, list, and status operations map to friends repository operations.
- Results expose actor-relative public relationship status.
- Repository conflicts map to stable redacted application errors.
- Focused tests cover identity refusal, self-targeting, actor-relative statuses, conflict mapping, unit-of-work use, redaction, and dependency boundaries.
