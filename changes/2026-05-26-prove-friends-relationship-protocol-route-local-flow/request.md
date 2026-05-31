# Request

## Original Request

```text
继续推进，注意提交和推送。提交和推送的key在Git忽略的文件里有，你找一下。
```

English summary: continue the next bounded work item, then commit and push. The push key is in a Git-ignored local file; locate it without exposing the secret.

## Clarified Requirement

Advance `W-0242 Prove friends relationship protocol route in local alpha request flow` by proving the completed protected friends relationship route family through the existing local WebSocket/Protobuf request path.

The slice must stay proof-only:

- use the existing `vibit.friends.v1` Protobuf payloads;
- use the existing `friends.SendFriendRequest`, `friends.AcceptFriendRequest`, `friends.RejectFriendRequest`, `friends.RemoveFriend`, `friends.BlockPlayer`, `friends.UnblockPlayer`, `friends.ListFriendRelationships`, and `friends.GetFriendRelationshipStatus` routes;
- use the existing authenticated request wrapper and request-token protected route policy;
- prove a two-player local alpha request flow;
- record prerequisites, commands, request/response shape, redaction expectations, and Nakama/Pitaya reference alignment;
- avoid direct Nakama/Pitaya public API compatibility.

## User-Visible Outcome

`examples/local-alpha-example-client.sh` and `examples/local-alpha-request-loop.sh` now include the friends relationship route proof alongside the existing authenticated gameplay, storage object, presence/status, and failure-path local alpha proofs.

The examples README files document the friends proof path and its redaction expectations.

## Non-Goals

- Adding new protocol messages or routes.
- Changing friends relationship service behavior.
- Changing the friends relationship repository interface.
- Changing the PostgreSQL adapter.
- Changing migrations.
- Adding dependencies.
- Changing authentication/session behavior or route-protection semantics.
- Adding WebSocket handshake authentication.
- Adding event/audit tables.
- Adding groups, parties, chat, matchmaking, or match runtime.
- Adding stream subscriptions, broadcast fanout, delivery guarantees, or group messaging.
- Adding hosted deployments.
- Creating release binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, generated client libraries, or SDK packages.
- Executing public announcements beyond the GitHub release record.
- Running paid promotion.
- Adding direct Nakama/Pitaya API compatibility.

## Unknowns

- The next prototype-ready capability after proving the friends relationship route remains a separate selection step.
- Live PostgreSQL proof remains opt-in and is not required for this local alpha test-only proof.
- Git push credential location is local and Git-ignored; the secret must not be printed or committed.

## Acceptance Criteria

- [x] A focused E2E test proves the protected friends relationship route family through `FrameHandler`.
- [x] The proof uses the existing WebSocket/Protobuf envelope path and authenticated request wrapper.
- [x] The proof registers existing friends route handlers against the local alpha dispatcher fixture.
- [x] The proof uses two authenticated and connection-bound local players.
- [x] The proof covers send, status, accept, list, remove, block, unblock, resend, and reject behavior.
- [x] The proof checks redaction expectations.
- [x] `examples/local-alpha-example-client.sh` and `examples/local-alpha-request-loop.sh` include the friends proof.
- [x] `examples/README.md`, `examples/README.zh-CN.md`, `examples/local-alpha-client/README.md`, and `examples/local-alpha-client/README.zh-CN.md` document the proof path.
- [x] `ADR-0150` records the proof decision.
- [x] `runtime.friends_relationship_protocol_route_local_proof` check coverage exists.
- [x] `W-0242` is completed and `W-0243` is next-ready.
