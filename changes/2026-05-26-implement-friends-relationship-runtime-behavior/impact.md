# Impact

## Affected Modules

- `runtime`: gains an application-owned friends service package.
- `friends`: its storage-neutral repository vocabulary is consumed by application behavior.

## Module Ownership Impact

The friends domain module continues to own storage-neutral relationship vocabulary and repository interfaces. The application layer owns request identity enforcement, unit-of-work orchestration, public error mapping, and actor-relative result shaping.

## Public Contract Impact

No public command, query, event, permission, Protobuf, generated contract, or route contract is added.

## Data And Migration Impact

No migration or data shape changes are added. The service uses the existing repository interface and existing PostgreSQL unit-of-work handoff.

## Runtime Impact

The runtime gains an internal application service only. No startup wiring, protocol route registration, WebSocket behavior, generated output, or transport behavior changes are added.

## Test Impact

Focused package tests use fake unit-of-work and fake repository dependencies. They do not require live PostgreSQL, protocol adapters, WebSocket transport, generated clients, or startup composition.

## Compatibility Risks

The implementation is internal application behavior and does not change wire compatibility. The main risks are identity leakage and social graph leakage; the service maps failures through redacted application errors and rejects metadata-only identity before repository access.
