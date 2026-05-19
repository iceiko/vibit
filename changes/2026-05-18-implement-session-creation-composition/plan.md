# Plan

1. Extend the authentication service dependency set with a server-owned session id generator.
2. Extend the login unit-of-work capability with `NewSessionRepository`.
3. After token storage succeeds, generate a runtime session id and call `CreateRuntimeSession`.
4. Keep the runtime session active, player-owned, linked to `access_token_record_id`, and aligned to the access-token lifetime.
5. Return client-visible token material and application session metadata only after unit-of-work success.
6. Add focused authentication service tests for success and failure paths.
7. Add startup session id generator shape coverage.
8. Update ADR, conversation memory, work queue, manifests, guides, rule catalog, and check implementation.
9. Run Go tests and repository checks.
10. Fix any failures without expanding scope into route policy, protocol, transport, logout, or reconnect behavior.
