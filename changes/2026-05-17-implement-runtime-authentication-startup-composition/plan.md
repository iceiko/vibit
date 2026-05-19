# Plan

1. Update `runtime/cmd/vibit-server` to accept route-protector injection into the Protobuf frame handler.
2. Add startup authentication configuration parsing for verifier keys, access-token TTL, and token audience.
3. Compose `authentication.Service` for the PostgreSQL runtime store with PostgreSQL unit-of-work, `crypto/rand.Reader`, system clock, token record id generator, token lifetime, and token audience.
4. Build `authentication.NewRouteAccessTokenValidator` and `app.NewRouteProtector`.
5. Preserve memory startup without authentication repository capability.
6. Add focused tests for default/configured auth startup config, missing verifier config, invalid TTL, route protector injection, token record id generation, and metadata-only memory behavior.
7. Run Go tests and repository verification.
