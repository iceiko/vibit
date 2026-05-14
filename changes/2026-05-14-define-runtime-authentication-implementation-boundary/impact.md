# Impact

## Architecture

This change adds a boundary standard for future runtime authentication implementation.

It defines that future authentication orchestration is application-owned under `runtime/internal/app`, while the authentication module remains the storage-neutral repository owner and the PostgreSQL adapter remains persistence-only.

## Runtime

No runtime behavior is added.

No login, token generation, token validation, verifier comparison, logout execution, cleanup job, Protobuf message, WebSocket carrier, generated authentication shape, or authentication dependency is added.

## Contracts

No contract sources are changed.

The existing `contracts/runtime/authentication/` family remains the semantic source for future runtime behavior.

## Data

No migrations are added or changed.

Credential and token verifier migration sources remain as previously ratified.

## Compatibility

No public API, event, wire schema, database schema, or runtime behavior compatibility is changed.

## Agent Workflow

Future agents get an explicit implementation boundary and queue before touching runtime authentication code.
