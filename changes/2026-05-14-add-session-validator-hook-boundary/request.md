# Request

## Original Request

Continue advancing the project work queue. The active next-ready work item is `W-0026 Add session validator hook boundary`.

## Clarified Requirement

Add an application-owned session validation hook that can normalize request identity before route handlers run, while preserving the current metadata-only behavior and not implementing authentication.

The hook must be explicit enough for future authenticated session validation to plug in later, but this change must not choose or imply any authentication, token, credential, session store, player account lookup, Protobuf envelope, or WebSocket handshake design.

## User-Visible Outcome

The request path now has a durable application boundary where future session validation belongs. Agents can add future validation behavior through an injected `SessionValidator` instead of placing it in WebSocket transport, Protobuf decoding, inventory handlers, or player module internals.

## Non-Goals

- No real authentication.
- No token parsing.
- No credential lookup.
- No session persistence.
- No player account lookup.
- No Protobuf envelope changes.
- No WebSocket handshake changes.
- No inventory dependency on the player module.

## Unknowns

- The future production authentication scheme remains undecided.
- The future token format remains undecided.
- The future session store, expiration, refresh, and reconnect behavior remain undecided.
- The future player account schema remains undecided.

## Acceptance Criteria

- `runtime/internal/app` defines an application-owned `SessionValidator` boundary.
- A metadata-only default path preserves current behavior without authenticating clients.
- A dispatcher wrapper can run validation before route handlers receive requests.
- Injected validators can replace metadata-only identity with a validated identity.
- Invalid validation results can stop dispatch without reaching handlers.
- Focused Go tests cover metadata-only pass-through, injected replacement, invalid stop behavior, and missing dependency behavior.
