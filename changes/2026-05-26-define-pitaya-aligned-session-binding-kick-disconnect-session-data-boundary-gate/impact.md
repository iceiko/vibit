# Impact

## Runtime

The runtime impact is vocabulary-only. Existing authentication, request identity validation, first-message binding, active connection registry, logout behavior, WebSocket close handling, and presence lifecycle behavior remain unchanged.

No runtime behavior is added.

## Protocol

No protocol shape changes are made. No protocol routes, Protobuf sources, generated output, wire messages, or proof carriers are added.

## Persistence

No session data persistence, repository interface, PostgreSQL adapter, migration, session table, connection table, or distributed owner registry is added.

## Operations

No metrics endpoint, tracing pipeline, dashboard, hosted surface, SDK, release artifact, or operational behavior is added.

## Reference Alignment

This slice moves vibit closer to Pitaya architecture vocabulary by naming session binding, server-initiated kick/disconnect, session data, unbind, close reason, and presence handoff concepts while preserving vibit-owned boundaries and avoiding direct Nakama/Pitaya API compatibility.

The follow-up is `W-0269 Implement Pitaya-aligned session binding, kick/disconnect, and session data source-first map`.
