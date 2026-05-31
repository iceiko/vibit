# Conversation: Friends Relationship Runtime Behavior Gate

Date: 2026-05-26
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-friends-relationship-runtime-behavior-gate/`

Related artifacts:

- `docs/friends-relationship-runtime-behavior-gate.md`
- `docs/friends-relationship-runtime-behavior-gate.zh-CN.md`
- `decisions/ADR-0146-friends-relationship-runtime-behavior-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Context

`M-165/W-0237` implemented the friends relationship PostgreSQL adapter under `runtime/internal/platform/persistence/postgres`. The repository now had a storage-neutral friends repository interface, a PostgreSQL `friend_relationships` migration source, and a platform adapter with unit-of-work repository handoff.

The next-ready work item was `W-0238 Define friends relationship runtime behavior gate`.

## Maintainer Narrative

The user asked to continue the repository's next-ready work item.

```text
继续
```

The standing product direction was Nakama-first, with Pitaya deferred as a future distributed architecture reference, and with AI-native requirement-to-spec-to-tests-to-implementation-to-verification as the product design goal.

## Agent Response Summary

The agent advanced one bounded work item and defined the friends relationship runtime behavior gate without adding runtime behavior implementation, handlers, startup wiring, protocol routes, Protobuf source, generated output, repository or adapter changes, migrations, dependencies, authentication/session behavior changes, event/audit tables, broader social features, hosted or SDK scope, distributed runtime, or direct Nakama/Pitaya API compatibility.

The work added the gate standard, Simplified Chinese translation, ADR, change record, rule registration, manifest updates, guide updates, and next-ready queue entry for `M-167/W-0239 Implement friends relationship runtime behavior`.

The gate records:

- future owner `runtime/internal/app`;
- future package candidate `runtime/internal/app/friends`;
- future service source and test candidates;
- actor identity derived from validated request identity;
- metadata-only `player_id` and `session_id` refusal;
- actor-relative public status behavior;
- conservative `request_token_required` route-policy posture;
- unit-of-work handoff to `friends.Repository`;
- conflict and redaction expectations;
- future implementation test expectations;
- stop conditions that preserve runtime implementation, protocol, generated output, event/audit tables, broad social feature, distributed runtime, hosted, SDK, and direct compatibility deferrals.

## Decisions

- Complete `M-166/W-0238`.
- Accept `ADR-0146`.
- Add `runtime.friends_relationship_runtime_behavior_gate`.
- Keep future friends runtime behavior application-owned under `runtime/internal/app`.
- Keep metadata-only `player_id` and `session_id` unusable as actor proof.
- Select `M-167/W-0239 Implement friends relationship runtime behavior` as the next-ready work item.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: friend request, accept, reject, remove, block, unblock, list, and status behavior are core social graph behavior and should be prepared after repository and adapter persistence exist.

Pitaya guided the layering pressure: route/session context, handlers, and persistence should remain separate, with runtime service behavior above persistence and below protocol route expansion.

vibit adapted those references into its own model: an application-owned runtime behavior gate with no direct public API compatibility and no protocol or generated-output behavior in this slice.

## Artifacts

- `docs/friends-relationship-runtime-behavior-gate.md`
- `docs/friends-relationship-runtime-behavior-gate.zh-CN.md`
- `decisions/ADR-0146-friends-relationship-runtime-behavior-gate.md`
- `changes/2026-05-26-define-friends-relationship-runtime-behavior-gate/`
- `conversations/2026-05-26-friends-relationship-runtime-behavior-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/friends/module.yaml`
- `modules/friends/AGENTS.md`
- `modules/friends/AGENTS.zh-CN.md`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- Friends relationship runtime behavior implementation remains deferred to `W-0239`.
- Protocol routes and Protobuf messages remain deferred.
- Permission model details remain deferred to the application behavior implementation.
- Event/audit tables, chat, groups, parties, matchmaking, match runtime, SDKs, hosted surfaces, distributed runtime, and direct compatibility remain deferred.

## Follow-Up

The next-ready work item is:

```text
M-167/W-0239 Implement friends relationship runtime behavior
```

That next slice may implement application-owned friends relationship runtime behavior only within the gate boundaries accepted by `ADR-0146`.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, or raw private social graph data from a real user are recorded in this conversation log.
