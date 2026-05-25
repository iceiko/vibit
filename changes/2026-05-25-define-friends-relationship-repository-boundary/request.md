# Feature Request

## Original Request

Define the storage-neutral friends relationship repository boundary after the friend_relationships migration source, preserving Nakama-first planning and deferring implementation.

## Clarified Requirement

Define the friends relationship repository boundary as a documentation, ADR, manifest, and static-check slice only. The boundary must record the future repository owner, interface candidate, value types, lifecycle method vocabulary, conflict classes, transaction handoff, redaction, PostgreSQL adapter expectations, tests, and stop conditions before any repository interface implementation.

## User-Visible Outcome

A future agent can inspect the repository and understand exactly how to implement the first friends relationship repository interface without inventing SQL-shaped types, treating metadata as identity proof, or jumping to runtime/protocol behavior. `node tools/vibit inspect next --json` should advance to `W-0235 Implement storage-neutral friends relationship repository interface`.

## Nakama Capability Mapping

- Capability family: `friends_groups_and_parties`
- Product intent: Friends relationship state is a core Nakama-class social graph primitive. This request prepares the storage-neutral repository seam for future request, accept, reject, remove, block, unblock, list, and status behavior.
- API compatibility: This change does not authorize direct Nakama API compatibility.

## Pitaya Status

Pitaya remains a deferred future distributed architecture reference. This change does not introduce frontend/backend roles, RPC, service discovery, groups, cluster routing, distributed session behavior, or distributed social graph routing.

## Non-Goals

- Repository interface implementation.
- PostgreSQL adapter behavior or SQL execution beyond the existing W-0233 migration source.
- Runtime friendship behavior.
- Protocol routes, Protobuf source, or generated output.
- New migrations or event/audit tables.
- Dependencies, startup wiring, authentication/session behavior changes, delivery guarantees, stream subscriptions, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, SDK publication, generated client libraries, hosted deployments, release artifacts, public announcements, paid promotion, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Unknowns

- Exact Go interface method names remain a W-0235 implementation choice, bounded by the candidate vocabulary in the boundary.
- Public protocol error mapping remains deferred to a future protocol gate.
- Event/audit storage remains deferred to a later bounded work item.

## Acceptance Criteria

- `docs/friends-relationship-repository-boundary.md` and `docs/friends-relationship-repository-boundary.zh-CN.md` exist and define a boundary-only repository plan.
- `decisions/ADR-0142-friends-relationship-repository-boundary.md` records the accepted decision.
- `runtime.friends_relationship_repository_boundary` is registered and checks the boundary artifacts.
- `.arch/work-items.yaml` marks `W-0234` completed and opens `M-163/W-0235` as next-ready.
- No friends repository interface implementation, PostgreSQL adapter behavior, runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, event/audit table, distributed runtime, or direct compatibility is added.

## Test Expectations

- Positive behavior tests are not applicable because no runtime behavior is added.
- Negative behavior is covered by static checks that reject forbidden scope.
- Permission/authentication tests are not applicable because identity proof behavior remains deferred.
- Persistence/protocol/integration tests are not applicable because no adapter, protocol route, migration, or runtime behavior is added.
- Repository checks must cover W-0234 artifacts and absence of premature friends implementation files.

## Redaction Notes

Do not record raw device credentials, raw access tokens, verifier keys, credential or token digests, HMAC inputs or outputs, PostgreSQL DSNs with credentials, GitHub tokens, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or private social graph data beyond explicit request text and redacted placeholders.
