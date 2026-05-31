# Conversation: Friends Relationship Runtime Behavior Implementation

Date: 2026-05-26
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-friends-relationship-runtime-behavior/`

Related artifacts:

- `runtime/internal/app/friends/service.go`
- `runtime/internal/app/friends/service_test.go`
- `decisions/ADR-0147-friends-relationship-runtime-behavior-implementation.md`
- `changes/2026-05-26-implement-friends-relationship-runtime-behavior/`
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

`M-166/W-0238` completed the friends relationship runtime behavior gate. It authorized a bounded implementation under `runtime/internal/app/friends` while preserving protocol routes, Protobuf source/generated output, startup wiring, repository interface changes, PostgreSQL adapter changes, migrations, dependencies, authentication/session behavior changes, event/audit tables, broader social features, hosted surfaces, SDKs, distributed runtime, and direct Nakama/Pitaya API compatibility as deferred concerns.

The next-ready work item was `W-0239 Implement friends relationship runtime behavior`.

## Maintainer Narrative

The maintainer asked:

```text
继续
```

The agent continued using the repository workflow, with test-first implementation already established in this work item and verification-before-completion for the final status.

## Agent Response Summary

The agent advanced one bounded work item and implemented application-owned friends relationship runtime behavior under `runtime/internal/app/friends`.

The work added:

- `Service` and dependency validation;
- friend request send, accept, reject, remove, block, unblock, list, and status methods;
- validated request identity actor derivation;
- metadata-only identity rejection before id generation, unit-of-work access, or repository mutation;
- server-owned relationship id generation;
- incoming pending enforcement for accept/reject behavior;
- actor-relative public status mapping;
- unit-of-work friends repository handoff;
- redacted public error mapping;
- focused fake-repository tests;
- ADR, change spec, manifest, check-rule, and continuation updates.

## Decisions

- Complete `M-167/W-0239`.
- Accept `ADR-0147`.
- Add `runtime.friends_relationship_runtime_behavior_implementation`.
- Keep friends protocol routes and generated output out of this slice.
- Select a later friends relationship protocol route gate as the next bounded direction.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: friends relationships are a common game backend social graph primitive.

Pitaya guided the layering pressure: handler/session context and persistence must remain separated.

vibit adapted those lessons into a validated-request-identity application service, with no direct public API compatibility and no protocol route exposure in this slice.

## Artifacts

- `runtime/internal/app/friends/service.go`
- `runtime/internal/app/friends/service_test.go`
- `decisions/ADR-0147-friends-relationship-runtime-behavior-implementation.md`
- `changes/2026-05-26-implement-friends-relationship-runtime-behavior/`
- `conversations/2026-05-26-friends-relationship-runtime-behavior-implementation.md`
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

- Protocol route names and Protobuf request/response shapes remain deferred.
- Startup wiring and route registration remain deferred.
- Event/audit tables remain deferred.
- Chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime social context, SDK publication, hosted surfaces, distributed runtime, and direct compatibility remain deferred.
- Live PostgreSQL verification remains unnecessary for this application-service slice because persistence mapping is covered by the adapter tests.

## Follow-Up

- Define the friends relationship protocol route gate before exposing the service through WebSocket/Protobuf.
- Preserve metadata-only identity refusal and actor-relative public status in any future handler.
- Keep route policy on protected request-token behavior unless a later route-policy ADR changes that posture.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, or private real social graph data are recorded in this conversation log.
