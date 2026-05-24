# Feature Request

## Original Request

Define a Nakama-first friends relationship persistence schema gate before SQL migration source, repositories, adapters, runtime behavior, protocol routes, generated output, or broader social features are implemented.

## Clarified Requirement

Define a gate-only persistence schema standard for the future friends relationship social graph state selected by `ADR-0139`. The gate must record future table candidates, canonical pair identity, lifecycle state representation, actor-specific block representation, indexes, uniqueness, timestamps, event/audit posture, redaction, migration-source candidate, repository owner, PostgreSQL adapter owner, stop conditions, and the next bounded migration-source work item.

## User-Visible Outcome

Future agents can continue from a concrete, Nakama-first persistence target for player friendship relationships without inventing ad hoc social graph tables or skipping tests and boundaries. The next step is explicitly limited to a migration-source-only slice.

## Nakama Capability Mapping

- Capability family: `friends_groups_and_parties`
- Product intent: Friends relationships are a core Nakama-class social graph primitive. A durable relationship schema is required before future request, accept, reject, remove, block, unblock, list, and status behavior can be implemented safely.
- API compatibility: This gate does not authorize direct Nakama API compatibility.

## Pitaya Status

Pitaya remains a deferred future distributed architecture reference. This gate does not introduce frontend/backend roles, RPC, service discovery, groups, cluster routing, distributed sessions, or distributed social graph routing.

## Non-Goals

- Add SQL migration source or create `friend_relationships`.
- Add repository interfaces, PostgreSQL adapters, runtime friendship behavior, protocol routes, Protobuf source, generated output, dependencies, or startup wiring.
- Add chat, groups, parties, matchmaking, match runtime, operations/admin behavior, SDK publication, hosted deployments, release artifacts, distributed runtime, event/audit tables, or direct Nakama/Pitaya compatibility.

## Unknowns

- Exact SQL constraint and index names are deferred to the migration-source slice.
- Repository interface shape is deferred.
- PostgreSQL adapter SQL and error mapping are deferred.
- Runtime behavior, idempotency, conflict handling, route mapping, and protocol payloads are deferred.
- Event/audit table and outbox posture are deferred.

## Acceptance Criteria

- [x] The gate defines future friends relationship table candidates, canonical pair identity, lifecycle state, block representation, indexes, uniqueness, timestamps, event/audit posture, and redaction expectations before SQL.
- [x] The gate records a future migration source candidate, repository owner, PostgreSQL adapter owner, transaction posture, test expectations, and stop conditions.
- [x] The gate keeps Nakama as the primary product capability reference and Pitaya deferred.
- [x] The gate opens `W-0233 Add friends relationship migration source`.
- [x] The gate adds no migration source, relationship table, repository interface, adapter implementation, runtime behavior, protocol route, Protobuf source, generated output, dependency, hosted surface, SDK, distributed runtime, or direct compatibility scope.

## Test Expectations

- Positive behavior tests are deferred because this gate adds no runtime behavior.
- Negative behavior tests are deferred because this gate adds no runtime behavior.
- Permission/authentication tests remain required for future behavior but are not changed by this schema gate.
- Persistence tests are planned for the future migration-source slice and must cover table, columns, checks, indexes, forbidden columns, and static boundary checks.
- Repository checks must cover this gate through `runtime.friends_relationship_persistence_schema_gate`.

## Redaction Notes

Do not record raw device credentials, raw access tokens, verifier keys, credential or token digests, HMAC inputs or outputs, PostgreSQL DSNs with credentials, GitHub tokens, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or private social graph records beyond explicit planning vocabulary.
