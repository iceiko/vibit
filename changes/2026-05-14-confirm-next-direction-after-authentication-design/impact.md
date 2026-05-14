# Impact Analysis

## Affected Modules

- `runtime`
- `player`

## Module Ownership Impact

No ownership changes are introduced.

Runtime remains the owner of application session-validation handoff, request identity, dispatch, protocol adaptation, transport, and platform concerns.

The player module remains the owner of stable player identity and account lifecycle, but this direction does not add login handlers or account routes.

## Public Contract Impact

No public commands, queries, events, errors, or permissions are added, changed, or removed.

The selected next milestone will ratify login and token behavior before contracts are implemented.

## Data And Migration Impact

No migration source is added or changed.

Credential, external identity, token, and session persistence schemas remain deferred until future work explicitly ratifies them.

## Runtime Impact

No Go runtime behavior is added.

The current metadata-only session validator remains non-authenticated and must not be treated as production proof.

## Protocol Impact

No Protobuf source, generated output, envelope behavior, token carrier, system message, WebSocket subprotocol, header, cookie, or route behavior is changed.

## Reference Alignment Impact

Nakama remains the capability reference for authentication methods, session tokens, refresh, expiration, revocation, and realtime socket/session relationships.

Pitaya remains the vocabulary reference for session binding and handler session context.

Neither reference governs vibit's API shape.

## Compatibility Risks

No API, event, data, protocol envelope, or WebSocket handshake compatibility risk is introduced because no public runtime behavior changes.
