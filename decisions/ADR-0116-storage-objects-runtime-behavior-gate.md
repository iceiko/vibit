# ADR-0116: Storage Objects Runtime Behavior Gate

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-define-storage-objects-runtime-behavior-gate/`

Related conversations:

- `conversations/2026-05-22-storage-objects-runtime-behavior-gate.md`

Related artifacts:

- `docs/storage-objects-runtime-behavior-gate.md`
- `docs/storage-objects-runtime-behavior-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/storage/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-135/W-0207` implemented the storage objects PostgreSQL adapter under `runtime/internal/platform/persistence/postgres`. The repository now has a storage-neutral repository interface, a PostgreSQL table migration, and a platform adapter with unit-of-work repository handoff.

The next bounded step is to define how runtime behavior may later use those persistence pieces. That behavior must derive ownership from validated request identity, preserve metadata-only identity protections, keep permission and route policy application-owned, and avoid protocol or generated-output changes until a later protocol slice.

Nakama keeps durable storage-object-like game state as a core backend capability. Pitaya reinforces separating route/session context, handlers, and persistence. vibit adapts those lessons by defining an application-owned runtime behavior gate, not by copying public APIs.

## Decision

Accept `docs/storage-objects-runtime-behavior-gate.md` as the gate for future storage objects runtime behavior.

The gate records:

- future application owner `runtime/internal/app`;
- future package candidate `runtime/internal/app/storage`;
- future service source and test candidates;
- validated request identity as the only first-posture owner source;
- explicit refusal to treat metadata-only `player_id` or `session_id` as proof;
- first owner kind `player`;
- route-policy expectation `request_token_required`;
- permission, validation, conflict, redaction, and unit-of-work handoff expectations;
- future runtime behavior implementation direction `storage_objects_runtime_behavior_implementation`;
- stop conditions before runtime implementation, handlers, startup wiring, protocol routes, Protobuf sources, generated output, repository/adapter changes, dependencies, migrations, authentication/session changes, hosted deployment, release artifacts, public announcements, paid promotion, object/blob storage, S3-compatible storage, or direct compatibility.

This ADR does not add runtime behavior implementation, runtime handlers, startup wiring, protocol routes, Protobuf sources, generated output, repository interface changes, PostgreSQL adapter changes, dependencies, migration changes, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, broad product module expansion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement runtime storage object behavior immediately after the PostgreSQL adapter.
- Add protocol routes and Protobuf messages together with behavior.
- Let client payloads supply owner ids.
- Treat envelope metadata `player_id` as sufficient owner proof.
- Put runtime behavior inside the storage domain module or PostgreSQL adapter package.
- Add direct Nakama-compatible storage object routes.

## Rationale

Storage object runtime behavior has identity and permission risk even though it is a prototype-friendly feature. The first posture must not let clients select arbitrary owners or turn metadata-only player ids into proof. A gate-only ADR keeps the next implementation slice focused on application behavior and makes route/protocol expansion a separate explicit decision.

Separating this gate from implementation also lets checks reject accidental runtime handlers, Protobuf shapes, or generated output before those surfaces are ratified.

## Agent Reasoning Summary

The safest continuation from the PostgreSQL adapter is an application behavior gate. It defines how future services should derive owner identity, require protected routes, validate inputs, map conflicts, and use unit-of-work repository handoff while preserving protocol, authentication/session, and broad product deferrals.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  identity_safety: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: low
  protocol_risk: deferred
  dependency_expansion: low
confidence: high
```

## Consequences

- `docs/storage-objects-runtime-behavior-gate.md` and `.zh-CN.md` exist.
- `runtime.storage_objects_runtime_behavior_gate` becomes the repository check rule for this slice.
- `M-136/W-0208` is completed.
- `M-137/W-0209 Implement storage objects runtime behavior` becomes the next-ready work item.
- Runtime behavior implementation, protocol behavior, generated output, migrations, dependencies, and authentication/session behavior remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- request identity validation semantics change materially;
- route policy selects a different first protected-route posture for storage objects;
- storage object ownership expands beyond player-owned objects before implementation;
- the repository interface changes before runtime behavior implementation;
- storage objects become large blob or S3-compatible object storage instead of small JSON game state;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Implement storage objects runtime behavior only after this gate is accepted.
- Keep protocol routes, Protobuf messages, generated output, startup wiring, public ACLs, admin search, group/guild scopes, large blob storage, S3-compatible storage, and direct compatibility behind later gates.
- Preserve metadata-only identity refusal in the implementation tests.
