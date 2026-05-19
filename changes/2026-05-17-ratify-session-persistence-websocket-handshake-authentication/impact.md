# Impact

## Affected Modules

- `runtime`: owns the current request-level validation path, future session binding posture, and future session store gates.
- `authentication`: remains storage-neutral and does not own WebSocket transport, Protobuf framing, runtime session persistence, or direct public API compatibility.

## Module Ownership Impact

The gate clarifies that future session behavior is application/runtime-owned. The authentication service can validate proof and produce request identity, but it does not own WebSocket connections, first-message protocol negotiation, or long-lived connection membership.

## Public Contract Impact

No commands, queries, events, errors, permissions, or wire messages are added.

Future first-message binding will require a separate protocol/system-message contract before implementation.

## Data And Migration Impact

No data schema is added.

Future session persistence may prefer PostgreSQL as the first durable target because PostgreSQL is already the accepted durable store for module state. That preference does not add a table, migration, repository, or adapter in this gate.

## Runtime Impact

Current runtime behavior remains:

- Login uses the public `runtime.authentication.AuthenticateWithDeviceCredential` command route.
- Protected gameplay routes use request-level validation through `vibit.authentication.v1.AuthenticatedRequest`.
- `RequestIdentity.SessionValidated` remains false.
- WebSocket transport ignores credential carriers.

## Reference Impact

Nakama is used for session/socket lifecycle coverage: authenticate first, then associate a realtime connection with validated identity/session state.

Pitaya is used for architecture vocabulary: connection acceptors, session binding, route handlers, groups, and later frontend/backend separation must stay layered.

## Compatibility Risks

Low. This is a standards-only ratification.

The primary risk is future accidental implementation of token carriers in transport or envelope metadata; the new check rule guards against that.
