# Plan

1. Read the protocol route gate, generated-output standard, runtime protocol adapter boundary, existing storage route implementation, and friends application service.
2. Add failing tests for friends Protobuf bridge mapping and bootstrap route handlers.
3. Add `proto/vibit/friends/v1/friends.proto` and generate Go Protobuf output with `buf generate`.
4. Add friends route keys under `runtime/internal/app/friends/routes.go`.
5. Implement `runtime/internal/platform/protocol/protobuf/friends_bridge.go` and register it from the existing payload bridge path.
6. Implement `runtime/internal/app/bootstrap/friends.go` with route registration, validated identity handoff, malformed payload rejection, and redacted public error mapping.
7. Register the friends service and handlers in PostgreSQL startup composition and add command bypass routes for service-owned unit-of-work behavior.
8. Add ADR-0149, conversation memory, and W-0241 change metadata.
9. Update architecture manifests, module guides/manifests, README/alpha docs, rules, and tooling checks so W-0241 is recorded as completed and W-0242 is next-ready.
10. Run targeted Go tests, Buf checks, repository checks, full Go tests, and whitespace verification.
