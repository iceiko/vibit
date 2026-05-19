# Request

## Original Request

The maintainer clarified that vibit should be replanned as a Nakama/Pitaya-class game backend product that covers their common capability families, then asked to continue development.

## Clarified Requirement

Implement the bounded `W-0170` protocol logout route slice authorized by `docs/protocol-logout-route-gate.md` and `ADR-0079`.

The implementation must expose the existing `authentication.Service.LogoutAccessToken` behavior through the WebSocket Protobuf route surface as:

```text
runtime.authentication.LogoutAccessToken
```

This work continues the replanned lifecycle-first route: close the authentication/session/connection lifecycle loop before broad presence, chat, social, matchmaking, match runtime, operations, SDK, or distributed runtime expansion.

## User-Visible Outcome

Clients using the PostgreSQL runtime composition can send a Protobuf logout command carrying the presented opaque access token and receive a logout response or redacted public authentication error.

Agents and maintainers can verify that this route is protocol-visible without mixing token revocation with socket close, runtime session revocation, reconnect behavior, protocol session carriers, or direct Nakama/Pitaya API compatibility.

## Non-Goals

- Do not close WebSocket sockets from the logout route.
- Do not revoke runtime sessions from the logout route.
- Do not invalidate active connection registry records from the logout route.
- Do not add reconnect, resume, duplicate replacement, or connection epoch behavior.
- Do not add protocol session carriers or change the existing Protobuf envelope.
- Do not add presence, chat, friends, groups, parties, leaderboards, tournaments, matchmaking, match runtime, admin disconnect, SDK, cluster, or distributed runtime behavior.
- Do not add new dependencies.
- Do not provide direct Nakama/Pitaya API compatibility.
- Do not make the memory runtime provide durable logout behavior.

## Unknowns

- Whether a later transport close handoff should close the current socket after logout or only after explicit disconnect/kick policy.
- Whether runtime session revocation should be a separate client-facing route, an admin operation, or a policy side effect after reconnect semantics are ratified.
- Whether `token_record_id` remains present in all protocol logout responses or is later restricted to audit/admin surfaces.

## Acceptance Criteria

- [ ] Logout request and response messages are added to `proto/vibit/authentication/v1/authentication.proto`.
- [ ] Go Protobuf output is regenerated through Buf rather than hand-edited.
- [ ] The protocol bridge maps logout request and response payloads.
- [ ] Application bootstrap registers `runtime.authentication.LogoutAccessToken` only when an authentication service is composed.
- [ ] PostgreSQL startup composition registers the route and bypasses the outer transactional dispatcher for logout.
- [ ] Focused tests cover route registration, bridge mapping, route policy posture, transaction bypass, redacted errors, unchanged envelope, and no socket/session/reconnect side effects.
- [ ] Repository manifests, work queue, AGENTS guides, and checks record `W-0170` completion and the next lifecycle direction.
