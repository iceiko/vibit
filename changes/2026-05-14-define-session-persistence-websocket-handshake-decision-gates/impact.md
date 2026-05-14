# Impact

## Affected Modules

`runtime` is affected because the standard defines future gates for application validation, session persistence, reconnect behavior, WebSocket handshake authentication, and Protobuf envelope interaction.

`player` is affected only as the owner of stable player identity. No player account lifecycle behavior, schema, handler, or route changes are introduced.

## Module Ownership Impact

No ownership moves.

Application dispatch remains the owner of normalized request identity handoff. WebSocket transport remains credential-neutral. Protobuf protocol adapters continue to treat session fields as metadata carriers until validation exists.

## Public Contract Impact

No command, query, event, error, or permission contract is added or changed.

The change defines future artifact gates that must exist before session persistence contracts, handshake contracts, transport/auth handoff contracts, protocol changes, or route authentication behavior are implemented.

## Data And Migration Impact

No database schema or migration is added.

No session tables, token tables, credential tables, external identity tables, or player account lifecycle changes are introduced.

## Protocol Impact

No Protobuf source, generated Protobuf output, envelope field, message kind, handshake message, system message, or WebSocket route behavior is changed.

## Reference Alignment Impact

Nakama session token, refresh token, session expiration, revocation, and realtime socket binding concepts are mapped into deferred or adapted vibit decision gates.

Pitaya session binding and handler session-context vocabulary are retained as architecture vocabulary without adopting its API surface.

## Compatibility Risks

This is a decision-gate standard and manifest update. It should not break runtime behavior.

The main risk is implying that one validation model is selected. The standard explicitly labels all validation models as future options.

## Test Impact

No Go runtime code changes are required.

Repository verification should cover runtime boundary checks, architecture checks, work queue consistency, memory/change-spec checks, and the full repository check.
