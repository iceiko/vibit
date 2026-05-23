# Conversation: Storage Objects Protocol Route Local Proof

Date: 2026-05-23
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-22-prove-storage-objects-protocol-route-local-flow/`

Related artifacts:

- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `examples/local-alpha-request-loop.sh`
- `examples/README.md`
- `examples/README.zh-CN.md`
- `decisions/ADR-0120-storage-objects-protocol-route-local-proof.md`
- `changes/2026-05-22-prove-storage-objects-protocol-route-local-flow/`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`

## Context

`M-139/W-0211` completed the own-player storage object protocol route implementation. The next-ready work item was `W-0212 Prove storage objects protocol route in local alpha request flow`.

The requested continuation needed to keep Nakama/Pitaya alignment explicit and then commit and push the result.

## Maintainer Narrative

The maintainer asked:

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。
```

English summary: continue the work, keep Nakama/Pitaya alignment explicit, then commit and push the result.

## Agent Response Summary

The agent advanced one bounded work item and proved the completed storage object protocol route family in the local alpha request flow.

The work added:

- storage route registration to the authenticated gameplay E2E fixture;
- a test-only in-memory storage repository implementing the existing storage repository interface;
- `TestStorageObjectsProtocolRouteLocalAlphaFlow`;
- local proof coverage for onboarding, login, first-message connection binding, authenticated storage object put/get/list/delete, and post-delete not-found behavior;
- redaction assertions for the storage error path;
- local alpha request-loop script updates;
- English and Chinese examples README updates;
- ADR, change spec, manifest, and check-rule updates.

## Decisions

- Complete `M-140/W-0212`.
- Accept `ADR-0120`.
- Add `runtime.storage_objects_protocol_route_local_proof`.
- Keep the proof local and test-only.
- Keep storage object routes own-player only.
- Keep owner identity out of client payloads and derive it from validated request identity.
- Keep public ACLs, cross-owner access, batch writes, JSON patch/merge, TTL, script hooks, blob/S3 storage, repository/adapter/migration/dependency changes, authentication/session changes, production memory storage behavior, and direct compatibility deferred.
- Add `W-0213 Confirm next alpha direction after storage objects local proof` as the next bounded continuation item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability being proven: player-owned storage objects are a common durable game-state primitive with collection/key/value/version and put/get/list/delete semantics.

Pitaya guided the layering being proven: transport, session metadata, protocol serialization, route protection, application handlers, backend service behavior, and repository handoff remain separated.

vibit adapts those lessons into its own WebSocket/Protobuf route model and application-owned service boundary. This slice does not add direct Nakama or Pitaya public API compatibility.

## Artifacts

- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `examples/local-alpha-request-loop.sh`
- `examples/README.md`
- `examples/README.zh-CN.md`
- `decisions/ADR-0120-storage-objects-protocol-route-local-proof.md`
- `changes/2026-05-22-prove-storage-objects-protocol-route-local-flow/`
- `conversations/2026-05-23-storage-objects-protocol-route-local-proof.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`

## Open Questions

- The next product direction after storage object local proof still needs confirmation.
- Public permissions, ACLs, admin search, group/guild/party/room/match scopes, batch writes, JSON patch, merge semantics, TTL, script hooks, blob/S3 storage, and direct compatibility remain deferred.
- Production memory storage behavior remains deferred.

## Follow-Up

- Advance `W-0213 Confirm next alpha direction after storage objects local proof`.
- Preserve the current storage route Protobuf shape unless a later bounded work item records a concrete reason to change it.
- Keep Nakama/Pitaya alignment explicit as capability and layering guidance, not direct API compatibility.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, DSNs with credentials, or raw storage object values from a real user are recorded in this conversation log.
