# Request

## Original Request

The maintainer asked to continue advancing the project toward Nakama/Pitaya-class server capability and specifically called out that excessive confirmation gates are hurting development velocity.

## Clarified Requirement

Advance `M-103/W-0175` as the first prospective Tier 2 functional slice under `ADR-0082`, embedding the reconnect/connection epoch boundary in this change spec instead of creating another pure confirmation gate.

The implementation must stay bounded to server-observed connection epoch semantics.

## User-Visible Outcome

The application-owned active connection registry now has a concrete first epoch progression rule:

- A newer server-observed epoch for the same connection id supersedes earlier active records.
- A stale or repeated epoch after a newer epoch has been observed is rejected with the redacted `connection_epoch_stale` error.
- Superseded records remain inspectable by connection id and epoch but are excluded from active targeting lists.

This gives future reconnect, presence, and match-runtime work a tested lifecycle primitive without adding protocol resume behavior yet.

## Non-Goals

- No reconnect token.
- No resume token.
- No durable or distributed reconnect behavior.
- No duplicate replacement socket close.
- No WebSocket close code mapping.
- No player-visible close reason text.
- No logout-triggered socket close.
- No runtime session revocation.
- No protocol session carrier.
- No Protobuf source, generated output, or existing envelope change.
- No presence lifecycle, operations/admin disconnect, or broader product module behavior.
- No new dependencies.
- No direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- [x] The change spec embeds the reconnect/connection epoch boundary and records Tier 2 gate density.
- [x] The application registry adds a `superseded` terminal state for older active epochs.
- [x] Newer server-observed epochs supersede older active records for the same connection id.
- [x] Stale or repeated epochs after a newer observed epoch fail closed with `connection_epoch_stale`.
- [x] Superseded records are visible through lifecycle inspection and excluded from active target lists.
- [x] Focused tests cover epoch supersession, stale rejection, lifecycle inspection, and active-list exclusion.
- [x] Repository checks verify no protocol, transport credential, logout, session revocation, dependency, or direct compatibility side effects.
