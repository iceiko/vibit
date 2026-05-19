# Impact

Runtime impact:

- Adds the first public authentication command route for `AuthenticateWithDeviceCredential`.
- Adds Protobuf request and response payloads for the public login route.
- Adds protocol bridge behavior between `vibit.authentication.v1` messages and `runtime/internal/app/authentication` request/result types.
- Adds application bootstrap route handlers that call the existing authentication service.
- Registers the public login route in the PostgreSQL runtime startup path.
- Adds a narrow transaction-wrapper bypass for the public login route.

Architecture impact:

- Authentication service behavior remains application-owned.
- The protocol adapter only bridges payloads; it does not validate credentials or call repositories.
- WebSocket transport remains credential-neutral.
- The existing Protobuf envelope remains unchanged.
- Memory startup remains a bootstrap path without durable authentication repository capability.
- Nakama informs the authenticate-before-gameplay sequence; Pitaya informs the transport/session/route/handler layering.

Data impact:

- No migrations are added or changed.
- No repository interfaces are changed.
- No PostgreSQL adapters are changed.
- No raw credential or token material is stored outside the already-ratified service behavior.

Security impact:

- `credential_proof` and `access_token` are treated as secret values.
- Public error mapping avoids leaking credential lookup state, verifier details, raw proof material, or raw token material.
- Protected routes still require the existing request-level access-token wrapper.

Compatibility impact:

- Existing inventory routes and protected route behavior remain unchanged.
- Existing WebSocket envelope fields are unchanged.
- Clients must use the new Protobuf command payload for login; no HTTP or handshake credential carrier is added.
