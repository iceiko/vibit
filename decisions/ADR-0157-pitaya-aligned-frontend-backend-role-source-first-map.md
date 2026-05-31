# ADR-0157: Pitaya-Aligned Frontend Backend Role Source-First Map

Status: Accepted
Date: 2026-05-31
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-frontend-backend-role-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-frontend-backend-role-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-frontend-backend-role-boundary-gate.md`
- `docs/pitaya-aligned-frontend-backend-role-boundary-gate.zh-CN.md`
- `decisions/ADR-0156-pitaya-aligned-frontend-backend-role-boundary-gate.md`
- `decisions/ADR-0155-pitaya-aligned-distributed-runtime-vocabulary-source-first-map.md`
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

ADR-0156 defined a gate-only Pitaya-aligned frontend/backend role boundary. It allowed `frontend_server` and `backend_server` as future architecture vocabulary and mapped current vibit single-process acceptor, session binding, protocol bridge, application dispatch, and module handler responsibilities to deferred role concepts.

The safe next step is source-first inspection. Agents need a concrete repository command that can summarize the role vocabulary and current single-process mapping before any runtime topology, RPC, service discovery, or process-role implementation is authorized.

## Decision

Implement `node tools/vibit inspect pitaya-roles --json` as the source-first Pitaya-aligned frontend/backend role map for `M-177/W-0249`.

The command reports:

- ADR-0156 as the source gate and ADR-0157 as the implementation decision.
- `runtime.pitaya_aligned_frontend_backend_role_source_first_map` as the check rule.
- Allowed role vocabulary: `frontend_server` and `backend_server`.
- Related vocabulary: `acceptor`, `session_binding`, and `route_handler`.
- Current single-process mappings for frontend and backend role planning.
- Explicit false deferrals for runtime behavior, role implementation, distributed runtime implementation, server-to-server RPC, service discovery, protocol, generated output, persistence, dependency, hosted, SDK, and direct compatibility surfaces.
- Redaction flags for credential, token, digest, key, DSN, database payload, and local secret file contents.
- `M-178/W-0250 Define Pitaya-aligned server-to-server RPC boundary gate` as the next-ready follow-up.

## Alternatives Considered

- Implement frontend/backend server roles immediately.
- Add server-to-server RPC or service discovery while adding the role map.
- Keep the role mapping only in ADR-0156 without a tool inspection surface.
- Fold role vocabulary into `node tools/vibit inspect pitaya-vocabulary --json` instead of adding a focused command.

## Rationale

Frontend/backend role vocabulary is useful for Pitaya-aligned planning, but it is also easy to confuse with permission to split processes or add cluster behavior. A dedicated source-first inspection command gives agents a precise place to inspect role vocabulary, current single-process mappings, deferrals, and redaction posture before any implementation gate.

This preserves vibit's agent-native model: inspectable source surfaces, explicit check rules, bounded vocabulary, and continuation memory before runtime behavior.

## Agent Reasoning Summary

The active work item is an inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all runtime topology, RPC, service discovery, protocol, persistence, dependency, hosted, SDK, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  role_mapping_clarity: high
  implementation_boundedness: high
  inspection_surface_value: high
  topology_change_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned frontend/backend role vocabulary without reading every architecture document.

This decision does not add distributed runtime behavior, frontend/backend server role implementation, server-to-server RPC, remote call behavior, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

Future server-to-server RPC planning must start with W-0250 as a boundary gate rather than implementation.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete frontend/backend role topology;
- the role inspection output creates confusion with public API compatibility;
- single-process acceptor, session binding, dispatch, protocol bridge, or handler ownership changes enough to require remapping;
- future Pitaya-aligned planning needs separate inspection surfaces for RPC, service discovery, groups, room broadcast, or cluster-safe routing.

## Follow-Up

- Complete `W-0250`: define a Pitaya-aligned server-to-server RPC boundary gate.
- Keep frontend/backend role implementation, server-to-server RPC implementation, remote calls, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint expansion, protocol shape, generated output, persistence, dependencies, hosted deployment, SDKs, and direct compatibility behind later bounded work items.
