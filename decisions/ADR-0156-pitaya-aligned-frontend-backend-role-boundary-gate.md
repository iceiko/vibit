# ADR-0156: Pitaya-Aligned Frontend Backend Role Boundary Gate

Status: Accepted
Date: 2026-05-31
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-frontend-backend-role-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-frontend-backend-role-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-frontend-backend-role-boundary-gate.md`
- `docs/pitaya-aligned-frontend-backend-role-boundary-gate.zh-CN.md`
- `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`
- `decisions/ADR-0155-pitaya-aligned-distributed-runtime-vocabulary-source-first-map.md`
- `decisions/ADR-0154-pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `modules/friends/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0155` implemented `node tools/vibit inspect pitaya-vocabulary --json` as the source-first map for the Pitaya-aligned distributed runtime vocabulary. The maintainer asked to continue moving toward Pitaya, but the repository still must preserve source-first boundaries and avoid runtime topology changes until separately authorized.

The most implementation-sensitive vocabulary in the current map is `frontend_server` and `backend_server`. Those terms are useful for future architecture planning, but they can also be mistaken for permission to split processes, add listeners, add RPC/remoting, or bypass vibit module contracts.

## Decision

Accept `docs/pitaya-aligned-frontend-backend-role-boundary-gate.md` and its Simplified Chinese translation as the gate for Pitaya-aligned frontend/backend role vocabulary.

Register `runtime.pitaya_aligned_frontend_backend_role_boundary_gate` as the repository check rule.

The gate defines:

- allowed role vocabulary: `frontend_server` and `backend_server`;
- related vocabulary: `acceptor`, `session_binding`, and `route_handler`;
- current single-process mapping for acceptor/session ingress, application dispatch, protocol bridge, and module handlers;
- ownership for role vocabulary and a future source-first role map;
- stop conditions for any implementation behavior.

Complete `M-176/W-0248` and open `M-177/W-0249 Implement Pitaya-aligned frontend/backend role source-first map` as next-ready.

This decision does not add distributed runtime behavior, frontend/backend server role implementation, server-to-server RPC, remote call behavior, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement a frontend/backend process split immediately.
- Add a source-first role inspection command without first defining a role-specific gate.
- Keep frontend/backend vocabulary only in the broader Pitaya vocabulary map.
- Switch near-term focus away from Pitaya vocabulary and back to another product capability family.

## Rationale

The broader Pitaya vocabulary map is useful, but role vocabulary deserves its own gate because it implies topology and ownership. Defining the gate first makes future planning precise while keeping vibit's current single-process runtime implementation stable.

This preserves vibit's agent-native model: vocabulary, ownership, deferrals, stop conditions, checks, and memory before implementation.

## Agent Reasoning Summary

The active work item is a gate. The correct continuation is to write the standard, ADR, change artifacts, repository checks, and manifest updates. The follow-up should implement a source-first role map, not frontend/backend server roles or distributed runtime behavior.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  role_boundary_clarity: high
  implementation_boundedness: high
  topology_change_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `frontend_server` and `backend_server` are allowed as future architecture vocabulary only.
- `runtime.pitaya_aligned_frontend_backend_role_boundary_gate` becomes the check rule for W-0248.
- `M-176/W-0248` is completed.
- `M-177/W-0249 Implement Pitaya-aligned frontend/backend role source-first map` becomes next-ready.
- Frontend/backend server role implementation, distributed runtime behavior, RPC/remotes, service discovery, distributed groups, broadcast fanout, cluster-safe session routing, protocol changes, generated output, persistence, dependencies, hosted deployment, SDKs, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete distributed runtime or role topology model;
- the role vocabulary creates confusion with public API compatibility;
- single-process acceptor, dispatch, or handler ownership changes enough to require remapping;
- prototype feedback shows a stronger near-term need for a different product capability family.

## Follow-Up

- Complete `W-0249`: implement a source-first Pitaya-aligned frontend/backend role map.
- Keep frontend/backend role implementation, cluster/RPC/service-discovery implementation, distributed groups, cluster-safe session routing, runtime endpoint expansion, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
