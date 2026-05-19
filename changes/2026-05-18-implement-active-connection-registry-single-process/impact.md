# Impact

## Affected Modules

- `runtime`: Adds application-owned active connection registry primitives.
- `authentication`: Remains token lifecycle owner but does not own active connection registries or socket policy.

## Module Ownership Impact

The registry belongs to `runtime/internal/app/connection`. It is not part of WebSocket transport, Protobuf protocol adapters, authentication repositories, session repositories, PostgreSQL adapters, or domain modules.

WebSocket transport remains credential-neutral. Future transport code may report server-observed open/close lifecycle facts through a later narrow handoff, but this change does not wire that handoff.

## Public Contract Impact

No public commands, queries, events, permissions, Protobuf sources, generated output, WebSocket handshake behavior, close codes, close reasons, or database schemas are added or changed.

## Runtime Behavior Impact

The new registry is a local application primitive. It does not yet participate in startup, Protobuf frame handling, logout, route protection, or WebSocket socket closure.

The registry stores server-observed connection id and epoch, active/terminal state, validated player linkage, optional runtime session id, optional access-token record id, and lifecycle timestamps. It does not store raw proof material or transport credential material.

## Reference Alignment

Nakama alignment: the implementation acknowledges that authenticated session material and realtime sockets need coordinated lifecycle state, but keeps logout, refresh, session management, and realtime socket invalidation as separate future surfaces.

Pitaya alignment: the implementation keeps acceptor/transport concerns, session context, route handlers, and connection management separated. Registry state is application-owned and policy-neutral.

## Compatibility Risks

The main risk is accidentally treating a registry record as authentication proof or adding hidden socket side effects. The implementation requires binding from validated player identity and keeps invalidation distinct from concrete socket close behavior.
