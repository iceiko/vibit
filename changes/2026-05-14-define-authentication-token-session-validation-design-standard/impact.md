# Impact

## Affected Modules

Affected module boundaries:

- `runtime`: owns application session validation handoff and future validation orchestration.
- `player`: owns stable player identity and player account lifecycle, but not credentials, tokens, external identity links, runtime sessions, WebSocket connections, or request validation results.

No runtime Go behavior is changed by this work item.

## Module Ownership Impact

This change clarifies ownership only.

It preserves:

- WebSocket transport owns connections and frame IO only.
- Protobuf protocol adapter owns envelope conversion only.
- Application dispatch owns request identity and session validation handoff.
- Player module owns account lifecycle only.
- Domain modules consume `RequestIdentity` for permissions but do not validate tokens or credentials.

## Public Contract Impact

No public command, query, event, error, or permission contract is added or changed.

Future authentication and token/session work must add contracts before implementation.

## Event Impact

No event contract is added or changed.

## Permission Impact

No permission contract is added or changed.

The change reinforces that metadata-only identity does not grant production permissions.

## Data And Migration Impact

No data migration is added.

The standard explicitly keeps credentials, tokens, external identity links, runtime sessions, WebSocket connection state, and request validation results out of `player_accounts` and `player_account_events`.

## Protocol Impact

No Protobuf source, generated output, envelope shape, or WebSocket handshake behavior is changed.

The standard defines future ask-first boundaries for any protocol or handshake changes.

## Reference Impact

The change maps Nakama account/auth/session token/refresh/realtime socket concepts and Pitaya session binding/handler/frontend/backend/realtime session vocabulary into vibit-native adopted, adapted, deferred, or rejected terms.

## Test Impact

No Go runtime tests are required because no Go runtime code changes.

Repository checks must verify architecture, runtime boundary, memory, work queue, and change-spec consistency.

## Compatibility Risks

No runtime behavior changes, so no client compatibility impact is expected.

The main risk is future drift if checks are not added. The next M-011 work item should add architecture checks for this boundary.
