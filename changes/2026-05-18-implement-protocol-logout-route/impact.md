# Impact Analysis

## Affected Modules

- `runtime`
- `authentication`

## Module Ownership Impact

Authentication service behavior remains owned by `runtime/internal/app/authentication`.

Protocol mapping is owned by `runtime/internal/platform/protocol/protobuf`.

Application route registration is owned by `runtime/internal/app/bootstrap`.

Route policy remains owned by `runtime/internal/app`.

Process startup composition remains owned by `runtime/cmd/vibit-server`.

WebSocket transport remains credential-neutral and does not own logout proof validation, token revocation, close policy, reconnect, or session behavior.

## Public Contract Impact

This change adds a protocol-visible command route:

```text
runtime.authentication.LogoutAccessToken
```

The route maps to existing semantic behavior from `contracts/runtime/authentication/commands/LogoutAccessToken.yaml`.

New wire messages:

- `vibit.authentication.v1.LogoutAccessTokenRequest`
- `vibit.authentication.v1.LogoutAccessTokenResponse`

No new semantic command, event, permission, or error family is introduced. Existing public authentication token error codes are reused.

## Data And Migration Impact

No migrations are added or changed.

Token revocation persistence already exists in the authentication service behavior and PostgreSQL adapter path. This change only exposes that behavior through a protocol route in the PostgreSQL runtime composition.

## Test Impact

Focused tests are required for:

- Protobuf bridge request and response mapping.
- Application bootstrap route registration and redacted error mapping.
- Route policy classification for service-validated logout.
- Transaction bypass for authentication lifecycle routes.
- Frame handling for public/service-validated logout route without `AuthenticatedRequest`.
- Startup composition exposing the route as service-validated/public.
- Existing envelope shape remaining unchanged.

## Documentation Impact

The existing gate standard remains the canonical route standard.

The change spec records this implementation slice. AGENTS guides and architecture manifests must be updated so future agents see that protocol logout route exposure is complete and that the next lifecycle work remains transport close handoff, reconnect/epoch, protocol session carrier, presence lifecycle, or operations posture.

## Compatibility Risks

The Protobuf authentication package gains new messages. This is additive to the current wire surface.

The existing Protobuf envelope is unchanged.

The route is explicit and service-validated. Normal protected gameplay routes still require the existing request-level authenticated wrapper unless separately classified.
