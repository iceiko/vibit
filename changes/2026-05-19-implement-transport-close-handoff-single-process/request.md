# Request

## Original Request

The maintainer asked to inspect the project directory, continue advancing the project, commit promptly, and push.

## Clarified Requirement

Advance `M-102/W-0174` by implementing the bounded `implement_transport_close_handoff_single_process` slice selected after the transport close handoff gate.

The implementation must add only a WebSocket transport-owned, single-process concrete close handoff by server-observed `connection_id` and `connection_epoch`.

## User-Visible Outcome

Maintainers and agents can call a narrow `runtime/internal/platform/transport/ws` handoff to request concrete close for an already accepted socket by server-owned connection id and epoch.

The result reports redacted, policy-neutral outcomes:

- `close_requested`
- `socket_not_found`
- `epoch_mismatch`
- `already_closed`
- `close_failed`

## Non-Goals

- No WebSocket close code mapping.
- No player-visible close reason text.
- No logout-triggered socket close.
- No runtime session revocation.
- No reconnect, resume, or duplicate replacement behavior.
- No protocol session carriers.
- No Protobuf source, generated output, or existing envelope changes.
- No authentication service call to WebSocket transport.
- No operations/admin disconnect surface.
- No new dependencies.
- No durable or distributed close handoff.
- No direct Nakama/Pitaya API compatibility.

## Unknowns

- Future close code and close reason policy remains unrated.
- Future reconnect/epoch behavior remains unrated.
- Future protocol session carrier behavior remains unrated.
- Future logout-triggered close coupling remains unrated.

## Acceptance Criteria

- [x] WebSocket transport exposes a narrow in-process handoff keyed by server-observed connection id and epoch.
- [x] The handoff can request concrete socket close for an already accepted socket.
- [x] Stale epoch, missing socket, already closed socket, and close failure outcomes are redacted and policy-neutral.
- [x] Application close policy and active connection registry remain owners of close decisions and registry lifecycle markers.
- [x] Tests cover target matching, stale epoch protection, redaction, transport credential neutrality, and unchanged protocol/logout/session behavior.
