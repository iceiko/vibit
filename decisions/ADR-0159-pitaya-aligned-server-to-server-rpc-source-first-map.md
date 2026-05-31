# ADR-0159: Pitaya-Aligned Server To Server RPC Source-First Map

Status: Accepted
Date: 2026-05-31
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-server-to-server-rpc-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-server-to-server-rpc-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.md`
- `docs/pitaya-aligned-server-to-server-rpc-boundary-gate.zh-CN.md`
- `decisions/ADR-0158-pitaya-aligned-server-to-server-rpc-boundary-gate.md`
- `decisions/ADR-0157-pitaya-aligned-frontend-backend-role-source-first-map.md`
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

ADR-0158 defined a gate-only Pitaya-aligned server-to-server RPC boundary. It allowed `server_to_server_rpc` and `remote_call` as future architecture vocabulary and mapped current vibit single-process application dispatch, protocol bridge, route handler, module handler, and absent service discovery responsibilities to deferred RPC concepts.

The safe next step is source-first inspection. Agents need a concrete repository command that can summarize RPC vocabulary and current single-process mapping before any RPC transport, remote call behavior, service discovery, topology split, protocol carrier, persistence, or dependency is authorized.

## Decision

Implement `node tools/vibit inspect pitaya-rpc --json` as the source-first Pitaya-aligned server-to-server RPC map for `M-179/W-0251`.

The command reports:

- ADR-0158 as the source gate and ADR-0159 as the implementation decision.
- `runtime.pitaya_aligned_server_to_server_rpc_source_first_map` as the check rule.
- Allowed RPC vocabulary: `server_to_server_rpc` and `remote_call`.
- Related vocabulary: `route_handler`, `module_handler`, `application_dispatch`, and `service_discovery`.
- Current single-process mappings for RPC and remote-call planning.
- Explicit false deferrals for runtime behavior, RPC implementation, remote calls, service discovery, frontend/backend role implementation, distributed runtime implementation, groups, broadcast fanout, cluster-safe routing, protocol, generated output, persistence, dependency, hosted, SDK, and direct compatibility surfaces.
- Redaction flags for credential, token, digest, key, DSN, database payload, and local secret file contents.
- `M-180/W-0252 Define Pitaya-aligned service discovery boundary gate` as the next-ready follow-up.

## Alternatives Considered

- Implement server-to-server RPC immediately.
- Add remote-call behavior or service discovery while adding the RPC map.
- Keep the RPC mapping only in ADR-0158 without a tool inspection surface.
- Fold RPC vocabulary into `node tools/vibit inspect pitaya-vocabulary --json` instead of adding a focused command.

## Rationale

Server-to-server RPC vocabulary is useful for Pitaya-aligned planning, but it is easy to confuse with permission to add remoting, discovery, transport behavior, or a distributed runtime. A dedicated source-first inspection command gives agents a precise place to inspect RPC vocabulary, current single-process mappings, deferrals, and redaction posture before any implementation gate.

This preserves vibit's agent-native model: inspectable source surfaces, explicit check rules, bounded vocabulary, and continuation memory before runtime behavior.

## Agent Reasoning Summary

The active work item is an inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all RPC implementation, remote-call behavior, service discovery, topology, protocol, persistence, dependency, hosted, SDK, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  rpc_mapping_clarity: high
  implementation_boundedness: high
  inspection_surface_value: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned server-to-server RPC and remote-call vocabulary without reading every architecture document.

This decision does not add server-to-server RPC implementation, remote call behavior, service discovery, frontend/backend server role implementation, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

Future service discovery planning must start with W-0252 as a boundary gate rather than implementation.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete RPC or distributed runtime model;
- the RPC inspection output creates confusion with public API compatibility;
- single-process dispatch, protocol bridge, route handler, module handler, or service discovery ownership changes enough to require remapping;
- future Pitaya-aligned planning needs separate inspection surfaces for service discovery, groups, room broadcast, or cluster-safe routing.

## Follow-Up

- Complete `W-0252`: define a Pitaya-aligned service discovery boundary gate.
- Keep server-to-server RPC implementation, remote calls, service discovery, frontend/backend role implementation, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint expansion, protocol shape, generated output, persistence, dependencies, hosted deployment, SDKs, and direct compatibility behind later bounded work items.
