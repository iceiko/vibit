# Impact

## Runtime

Defines a future application-owned authentication service interface boundary under `runtime/internal/app`.

No Go runtime code is added or changed.

## Authentication

Maps existing generated authentication contract shapes and semantic authentication contracts to future service-level request/result vocabulary, redaction expectations, error mapping, permission mapping, audit handoff, and request identity handoff.

The authentication repository remains storage-neutral. The PostgreSQL adapter remains persistence-only.

## Protocol And Transport

No Protobuf authentication messages are added. No WebSocket proof carriers are added. The WebSocket handshake and Protobuf envelope remain unchanged.

## Persistence

No migration source is added or changed. No repository interface is changed.

## Tooling

Adds a focused runtime check rule:

```text
runtime.application_authentication_service_interface_boundary
```

The rule verifies that the new boundary standard, ADR, manifests, guides, and status markers exist and that service code and runtime authentication behavior remain deferred.
