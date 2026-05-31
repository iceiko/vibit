# Request

## Original Request

Continue the next ready repository work item.

## Clarified Requirement

Advance `M-169/W-0241 Implement friends relationship protocol route`.

Implement only the protected friends relationship protocol route family authorized by `ADR-0148`, using the existing application-owned friends service from `ADR-0147`.

## User-Visible Outcome

The runtime can expose friends relationship operations through the existing WebSocket/Protobuf protocol path:

- `friends.SendFriendRequest`
- `friends.AcceptFriendRequest`
- `friends.RejectFriendRequest`
- `friends.RemoveFriend`
- `friends.BlockPlayer`
- `friends.UnblockPlayer`
- `friends.ListFriendRelationships`
- `friends.GetFriendRelationshipStatus`

## Non-Goals

- No repository interface changes.
- No PostgreSQL adapter changes.
- No migration changes.
- No dependency additions.
- No authentication or session semantic changes.
- No route-protection semantic changes.
- No delivery guarantees, stream subscriptions, chat rooms, groups, parties, broadcast fanout, matchmaking, or match runtime.
- No operations/admin behavior.
- No SDK publication or generated client libraries.
- No hosted deployments, release artifacts, public announcements, or paid promotion.
- No event/audit tables.
- No Pitaya-style distributed architecture.
- No direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- Add `proto/vibit/friends/v1/friends.proto` and generated Go Protobuf output under `runtime/internal/generated/proto/vibit/friends/v1/friends.pb.go`.
- Add route key helpers under `runtime/internal/app/friends/routes.go`.
- Add protocol bridge mapping under `runtime/internal/platform/protocol/protobuf/friends_bridge.go`.
- Add bootstrap route handlers under `runtime/internal/app/bootstrap/friends.go`.
- Register friends route handlers in PostgreSQL startup composition.
- Preserve request-token protected route posture through the existing authenticated wrapper flow.
- Derive actor identity only from validated `app.RequestIdentity`; do not add client-supplied actor id fields.
- Preserve metadata-only identity refusal in the application service.
- Map public service errors to redacted `FRIENDSHIP_*` application errors.
- Add focused bridge, handler, and startup tests.
