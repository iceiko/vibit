# Impact

## Affected Modules

- `runtime`: gains friends Protobuf bridge mapping, bootstrap handlers, route keys, generated output registration, and PostgreSQL startup registration.
- `friends`: its application service becomes reachable through protocol routes.

## Module Ownership Impact

The friends domain module continues to own storage-neutral relationship vocabulary and repository interfaces. The application service continues to own identity validation, unit-of-work orchestration, actor-relative public status, and public error mapping. The protocol adapter owns generated Protobuf payload conversion only.

## Public Contract Impact

This change adds first `vibit.friends.v1` Protobuf request/response messages and route keys for eight protected friends relationship routes. It does not add semantic contract source files, SDKs, generated client libraries, or direct external API compatibility.

## Data And Migration Impact

No migration or data shape changes are added. The route implementation uses existing application service behavior and existing repository/adapter/migration boundaries.

## Runtime Impact

The PostgreSQL-backed runtime registers friends route handlers and constructs the friends service with a random relationship id generator. Friends command routes are bypassed by the outer transactional dispatcher because the application service owns its own unit of work.

## Test Impact

Focused tests cover Protobuf request mapping, optional expected-version mapping, response mapping, route registration, validated identity handoff, redacted public error mapping, malformed payload rejection, and startup id generation.

## Compatibility Risks

The protocol envelope is unchanged. The new wire package is additive under `vibit.friends.v1`. The main risks are actor proof confusion and social graph leakage; this slice keeps actor proof out of payloads, relies on the existing protected-route wrapper, and maps failures to redacted public error messages.
