# Impact

## Runtime

Adds `OnboardLocalPlayerWithDeviceCredential` to the application authentication service under `runtime/internal/app/authentication`.

The service is local-only application behavior. It is not registered as a public WebSocket, HTTP, CLI, or Protobuf route.

Startup authentication composition now provides the service with a device credential random source and local onboarding id generators. This makes the service complete when authentication is composed, but it does not create accounts automatically and does not expose a new surface.

## Data

No migrations were added. Existing repository interfaces are reused:

- `player.Repository.CreatePlayerAccount`
- `authentication.Repository.StoreCredential`

The service stores only credential lookup and verifier digests plus metadata. Raw device credential text and bytes are not stored.

## Protocol

No Protobuf source, generated output, envelope field, route registration, WebSocket behavior, or authenticated request wrapper behavior changed.

## Authentication

The existing `AuthenticateWithDeviceCredential` login route remains proof-based and does not create accounts from `AccountCreationIntent`.

Onboarding does not return an access token and does not create a runtime session. Login remains the step that authenticates with the issued credential.

## Reference Alignment

Nakama alignment: the runtime now has a first local path to obtain player credential material before gameplay.

Pitaya alignment: onboarding remains application service behavior and does not move into transport or protocol handlers.
