# Feature Request

## Original Request

Define a Nakama-first friends relationship lifecycle gate for vibit so a future player social graph can support sending friend requests, accepting, rejecting, removing, blocking, unblocking, listing relationships, and reading relationship status through server-authoritative rules before runtime behavior, protocol routes, persistence, generated output, or broader social features are implemented.

## Clarified Requirement

Define a gate-only friends relationship lifecycle standard for vibit. The gate must translate the Nakama-style friends/social graph capability into vibit-owned, contract-first semantics for future request, accept, reject, remove, block, unblock, list, and status-read behavior. It must define future command/query/event/error/permission vocabulary, actor-relative relationship states, identity requirements, lifecycle invariants, redaction, test expectations, non-goals, verification, and the next bounded persistence schema gate before any runtime, protocol, generated, migration, repository, adapter, dependency, or broader social feature work.

## User-Visible Outcome

A developer or future agent can inspect the friends relationship lifecycle gate and know exactly what social graph behavior must be specified and tested before implementation. The next continuation step becomes `W-0232 Define friends relationship persistence schema gate`.

## Nakama Capability Mapping

- Capability family: `friends_groups_and_parties`
- Product intent: Friendship lifecycle is a core social graph primitive in Nakama-class backends. It should precede richer social surfaces such as groups, parties, chat targeting, invites, matchmaking filters, and match social context.
- Adopted: friend request, accept, reject, remove, block, unblock, list, and status-read capability semantics.
- Adapted: vibit defines its own contract-first, actor-relative, validated-identity lifecycle instead of external API compatibility.
- API compatibility: This gate does not authorize direct Nakama API compatibility.

## Pitaya Status

Pitaya remains a deferred future distributed architecture reference. Do not use this request to introduce frontend/backend roles, RPC, service discovery, groups, cluster routing, or distributed session behavior.

## Non-Goals

- Runtime friendship behavior.
- Protocol routes.
- Protobuf source.
- Generated output.
- Migrations.
- Repository interfaces.
- PostgreSQL adapters.
- Dependencies.
- Startup wiring.
- Authentication/session behavior changes.
- Delivery guarantees.
- Stream subscriptions.
- Chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, or operations/admin behavior.
- SDK publication, generated client libraries, hosted deployments, release artifacts, distributed runtime, or direct Nakama/Pitaya compatibility.

## Unknowns

- [x] Exact persistence table names are deferred to `W-0232`.
- [x] Canonical pair identity encoding is deferred to `W-0232`.
- [x] Duplicate request idempotency remains a future semantic choice before runtime implementation.
- [x] Concurrency conflict behavior remains deferred until schema/repository/runtime gates.
- [x] Protocol route and payload shape remain deferred.

## Acceptance Criteria

- [x] Future request, accept, reject, remove, block, unblock, list, and status-read semantics are defined.
- [x] Future command, query, event, error, permission, actor-relative state, invariant, redaction, and test vocabulary is recorded.
- [x] Validated player identity is required and metadata-only `player_id` or `session_id` is rejected as proof.
- [x] Nakama `friends_groups_and_parties` capability mapping is recorded while Pitaya remains deferred.
- [x] `W-0232 Define friends relationship persistence schema gate` is opened.
- [x] Runtime behavior, protocol routes, Protobuf source, generated output, migrations, repository interfaces, PostgreSQL adapters, dependencies, startup wiring, broader social features, hosted surfaces, distributed runtime, and direct compatibility are not added.

## Test Expectations

- [x] Positive future tests identified for request, accept, reject, remove, block, unblock, list, and status-read behavior.
- [x] Negative future tests identified for self-targeting, duplicate request behavior, invalid transition, blocked relationship interaction, missing target, missing relationship, and metadata-only identity.
- [x] Permission/authentication tests identified for validated player identity and client-supplied actor id handling.
- [x] Persistence tests are deferred to the next schema gate.
- [x] Protocol and E2E tests are deferred until route/runtime slices are authorized.

## Redaction Notes

Do not record raw device credentials, raw access tokens, verifier keys, credential or token digests, HMAC inputs or outputs, PostgreSQL DSNs with credentials, GitHub tokens, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or private user data beyond explicit request text.

Relationship graph details, target ids, actor ids, private statuses, and hidden conflict details are not log-safe by default.
