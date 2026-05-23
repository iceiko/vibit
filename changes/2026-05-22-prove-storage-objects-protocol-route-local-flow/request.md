# Request

## Original Request

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。
```

English summary: continue the next bounded work item, keep Nakama/Pitaya alignment explicit, then commit and push the result.

## Clarified Requirement

Advance `W-0212 Prove storage objects protocol route in local alpha request flow` by proving the completed own-player storage object route family through the existing local WebSocket/Protobuf request path.

The slice must stay proof-only:

- use the existing `vibit.storage.v1` Protobuf payloads;
- use the existing `storage.GetOwnStorageObject`, `storage.ListOwnStorageObjects`, `storage.PutOwnStorageObject`, and `storage.DeleteOwnStorageObject` routes;
- use the existing authenticated request wrapper and request-token protected route policy;
- record prerequisites, commands, request/response shape, redaction expectations, and Nakama/Pitaya reference alignment;
- avoid direct Nakama/Pitaya public API compatibility.

## User-Visible Outcome

`examples/local-alpha-request-loop.sh` now exercises both:

- the existing authenticated gameplay proof: onboarding -> login -> bind connection -> inventory -> presence -> logout -> revoked-token rejection;
- the storage object proof: onboarding -> login -> bind connection -> own-player storage object put/get/list/delete over the same WebSocket/Protobuf route flow.

The examples README documents the storage proof path and its redaction expectations.

## Non-Goals

- Adding new protocol messages or routes.
- Changing storage object service behavior.
- Changing the storage object repository interface.
- Changing the PostgreSQL adapter.
- Changing migrations.
- Adding dependencies.
- Changing authentication/session behavior or route-protection semantics.
- Adding WebSocket handshake authentication.
- Adding public ACLs or cross-owner access.
- Adding admin search.
- Adding group/guild/party/room/match storage scopes.
- Adding batch writes.
- Adding JSON patch or merge semantics.
- Adding TTL or script hooks.
- Adding broad operations/admin behavior.
- Adding hosted deployments.
- Creating release binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, or SDK packages.
- Executing public announcements beyond the GitHub release record.
- Running paid promotion.
- Adding large object/blob storage or S3-compatible object storage.
- Adding direct Nakama/Pitaya API compatibility.

## Unknowns

- The next product direction after proving storage objects remains a separate confirmation step.
- Live PostgreSQL proof remains opt-in and is not required for this local alpha test-only proof.

## Acceptance Criteria

- [x] A focused E2E test proves authenticated own-player storage object put/get/list/delete through `FrameHandler`.
- [x] The proof uses the existing WebSocket/Protobuf envelope path and authenticated request wrapper.
- [x] The proof registers existing storage route handlers against the local alpha dispatcher fixture.
- [x] The proof checks a post-delete not-found error and redaction expectations.
- [x] `examples/local-alpha-request-loop.sh` runs both the existing gameplay proof and the new storage proof.
- [x] `examples/README.md` and `examples/README.zh-CN.md` document the proof path.
- [x] `ADR-0120` records the proof decision.
- [x] `runtime.storage_objects_protocol_route_local_proof` check coverage exists.
- [x] `W-0212` is completed and `W-0213` is next-ready.
