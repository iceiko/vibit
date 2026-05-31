# Conversation: Friends Relationship Protocol Route Local Proof

Date: 2026-05-31
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-prove-friends-relationship-protocol-route-local-flow/`

Related artifacts:

- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `examples/local-alpha-example-client.sh`
- `examples/local-alpha-request-loop.sh`
- `examples/README.md`
- `examples/README.zh-CN.md`
- `examples/local-alpha-client/README.md`
- `examples/local-alpha-client/README.zh-CN.md`
- `decisions/ADR-0150-friends-relationship-protocol-route-local-proof.md`
- `changes/2026-05-26-prove-friends-relationship-protocol-route-local-flow/`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`

## Context

`M-169/W-0241` completed the friends relationship protocol route implementation. The next-ready work item was `W-0242 Prove friends relationship protocol route in local alpha request flow`.

The requested continuation needed to commit and push the result. The maintainer also noted that the push key is available in a Git-ignored local file.

## Maintainer Narrative

The maintainer asked:

```text
继续推进，注意提交和推送。提交和推送的key在Git忽略的文件里有，你找一下。
```

English summary: continue the work, then commit and push. The commit and push key is in a Git-ignored file; find it.

## Agent Response Summary

The agent advanced one bounded work item and proved the completed friends relationship protocol route family in the local alpha request flow.

The work added:

- friends route registration to the authenticated gameplay E2E fixture;
- deterministic multi-player local alpha fixture support;
- a test-only in-memory friends relationship repository implementing the existing friends repository interface;
- `TestFriendsRelationshipProtocolRouteLocalAlphaFlow`;
- local proof coverage for onboarding, login, first-message connection binding, authenticated friends send/status/accept/list/remove/block/unblock/resend/reject behavior;
- redaction assertions for the friends proof path;
- local alpha example script updates;
- English and Chinese examples README updates;
- ADR, change spec, manifest, and check-rule updates.

## Decisions

- Complete `M-170/W-0242`.
- Accept `ADR-0150`.
- Add `runtime.friends_relationship_protocol_route_local_proof`.
- Keep the proof local and test-only.
- Keep actor identity out of client payloads and derive it from validated request identity.
- Keep protocol shape, generated output, service behavior, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior, event/audit tables, broad social features, hosted surfaces, SDK publication, and direct compatibility deferred.
- Add `W-0243 Select next Nakama prototype-ready capability after friends relationship route proof` as the next bounded continuation item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability being proven: friends relationship request, accept/reject, list/status, remove, block, and unblock operations are common social graph primitives.

Pitaya guided the layering being proven: transport, session metadata, protocol serialization, route protection, application handlers, backend service behavior, and repository handoff remain separated.

vibit adapts those lessons into its own WebSocket/Protobuf route model and application-owned service boundary. This slice does not add direct Nakama or Pitaya public API compatibility.

## Artifacts

- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `examples/local-alpha-example-client.sh`
- `examples/local-alpha-request-loop.sh`
- `examples/README.md`
- `examples/README.zh-CN.md`
- `examples/local-alpha-client/README.md`
- `examples/local-alpha-client/README.zh-CN.md`
- `decisions/ADR-0150-friends-relationship-protocol-route-local-proof.md`
- `changes/2026-05-26-prove-friends-relationship-protocol-route-local-flow/`
- `conversations/2026-05-26-friends-relationship-protocol-route-local-proof.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`

## Open Questions

- The next product direction after friends relationship route local proof still needs selection.
- Event/audit tables, groups, parties, chat, matchmaking, match runtime, SDK publication, hosted deployment, and direct compatibility remain deferred.
- Production memory friends behavior remains deferred.

## Follow-Up

- Advance `W-0243 Select next Nakama prototype-ready capability after friends relationship route proof`.
- Preserve the current friends route Protobuf shape unless a later bounded work item records a concrete reason to change it.
- Keep Nakama/Pitaya alignment explicit as capability and layering guidance, not direct API compatibility.
- Use the local Git-ignored credential only for push and do not print, commit, or record its secret value.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, DSNs with credentials, or raw private social graph data are recorded in this conversation log.
