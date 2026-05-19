# Impact

## Affected Modules

- `runtime`: records the next runtime direction after login route exposure.
- `authentication`: records that session and handshake work remains outside the storage-neutral authentication module until later bounded gates.

## Module Ownership Impact

No ownership moves in this confirmation step.

Future session and handshake work must keep:

- WebSocket transport as a credential-neutral acceptor unless a later transport boundary grants a narrow extraction role.
- Authentication validation under application-owned or authentication-owned validation contracts.
- Domain modules dependent only on normalized `RequestIdentity`.

## Public Contract Impact

No commands, queries, events, errors, permissions, or Protobuf messages are added by this confirmation step.

## Data And Migration Impact

No tables, migrations, repositories, or PostgreSQL adapters are added.

## Reference Impact

Nakama is used for the login-session-socket capability sequence. Pitaya is used for session binding, handler routing, and connection/application separation vocabulary.

Neither reference becomes a public API compatibility target.

## Compatibility Risks

Low. This change records direction only.

The next gate must continue to avoid accidental token carriers in WebSocket transport or the existing envelope.

## Test Impact

No Go tests are required for the confirmation step. Repository checks verify the work queue, manifests, and change spec.
