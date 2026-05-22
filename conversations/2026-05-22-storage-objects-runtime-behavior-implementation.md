# Conversation: Storage Objects Runtime Behavior Implementation

Date: 2026-05-22
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-22-implement-storage-objects-runtime-behavior/`

Related artifacts:

- `runtime/internal/app/storage/service.go`
- `runtime/internal/app/storage/service_test.go`
- `decisions/ADR-0117-storage-objects-runtime-behavior-implementation.md`
- `changes/2026-05-22-implement-storage-objects-runtime-behavior/`
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

`M-136/W-0208` completed the storage objects runtime behavior gate. It authorized a bounded implementation under `runtime/internal/app/storage` while preserving protocol routes, Protobuf sources/generated output, startup wiring, repository interface changes, PostgreSQL adapter changes, migrations, dependencies, authentication/session behavior changes, blob/S3 storage, and direct Nakama/Pitaya API compatibility as deferred concerns.

The next-ready work item was `W-0209 Implement storage objects runtime behavior`.

## Maintainer Narrative

The maintainer asked:

```text
继续推进。
```

The agent continued using the repository workflow and Superpowers-style discipline manually, including test-first implementation and verification-before-completion.

## Agent Response Summary

The agent advanced one bounded work item and implemented application-owned storage object runtime behavior under `runtime/internal/app/storage`.

The work added:

- `Service` and dependency validation;
- own-object get, list, put, and delete methods;
- validated request identity owner derivation;
- metadata-only identity rejection before unit-of-work or repository access;
- input validation for collection, key, JSON object values, value size, list pagination, and expected versions;
- unit-of-work storage repository handoff;
- create-or-replace `put` behavior when no expected version is supplied;
- expected-version update and delete behavior;
- redacted public error mapping;
- focused fake-repository tests;
- ADR, change spec, manifest, check-rule, and continuation updates.

## Decisions

- Complete `M-137/W-0209`.
- Accept `ADR-0117`.
- Add `runtime.storage_objects_runtime_behavior_implementation`.
- Keep storage object protocol routes and generated output out of this slice.
- Select a later storage object protocol route gate as the next bounded direction.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: durable player-owned storage objects are a common game backend primitive.

Pitaya guided the layering pressure: handler/session context and persistence must remain separated.

vibit adapted those lessons into a validated-request-identity application service, with no direct public API compatibility and no protocol route exposure in this slice.

## Artifacts

- `runtime/internal/app/storage/service.go`
- `runtime/internal/app/storage/service_test.go`
- `decisions/ADR-0117-storage-objects-runtime-behavior-implementation.md`
- `changes/2026-05-22-implement-storage-objects-runtime-behavior/`
- `conversations/2026-05-22-storage-objects-runtime-behavior-implementation.md`
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

- Protocol route names and Protobuf request/response shapes remain deferred.
- Startup wiring and route registration remain deferred.
- Public permissions, ACLs, admin search, group/guild scopes, batch writes, JSON patch, merge semantics, TTL, script hooks, blob/S3 storage, and direct compatibility remain deferred.
- Live PostgreSQL verification remains unnecessary for this application-service slice because persistence mapping is covered by the adapter tests.

## Follow-Up

- Define the storage object protocol route gate before exposing the service through WebSocket/Protobuf.
- Preserve metadata-only identity refusal in any future handler.
- Keep route policy on protected request-token behavior unless a later route-policy ADR changes that posture.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, or raw storage object values from a real user are recorded in this conversation log.
