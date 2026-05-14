# Request

## Original Request

The maintainer asked to continue multiple work items unless a real decision point or blocker is encountered.

## Clarified Requirement

Advance `W-0061`: define decision gates for session persistence and WebSocket handshake authentication without implementing either path and without choosing a production validation model, session store, Protobuf envelope change, handshake/system message, or route-level authentication behavior.

## User-Visible Outcome

Agents gain a durable standard that separates request-level validation, first-message validation, handshake-level validation, every-request validation, and hybrid validation as future choices. Future session persistence, reconnect, connection epoch, token/session carrier, Protobuf envelope, and WebSocket handshake work now has explicit artifact gates before implementation.

## Non-Goals

- Do not choose request-level, first-message, handshake-level, every-request, or hybrid validation as the production model.
- Do not choose a session store.
- Do not add session tables or migrations.
- Do not choose token/session carrier behavior.
- Do not change Protobuf envelope fields.
- Do not add handshake or system messages.
- Do not add WebSocket handshake authentication behavior.
- Do not add route-level authentication implementation.
- Do not treat metadata-only identity as production proof.

## Acceptance Criteria

- Add an English standard and Simplified Chinese translation for session persistence and WebSocket handshake decision gates.
- Define the future validation-model options without selecting one.
- Define future artifact gates for session store, expiration, refresh, revocation, reconnect, connection epoch, Protobuf envelope interaction, and WebSocket handshake behavior.
- Update manifests and guides so future agents can discover the decision gates.
- Mark `W-0061` completed and `W-0062` next-ready.
- Run repository verification and record results.
