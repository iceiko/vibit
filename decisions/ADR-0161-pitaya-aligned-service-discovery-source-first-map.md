# ADR-0161: Pitaya-Aligned Service Discovery Source-First Map

Status: Accepted
Date: 2026-05-31
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-service-discovery-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-service-discovery-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-service-discovery-boundary-gate.md`
- `docs/pitaya-aligned-service-discovery-boundary-gate.zh-CN.md`
- `decisions/ADR-0160-pitaya-aligned-service-discovery-boundary-gate.md`
- `decisions/ADR-0159-pitaya-aligned-server-to-server-rpc-source-first-map.md`
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

`ADR-0160` defined a gate-only Pitaya-aligned service discovery boundary. It allowed `service_discovery`, `service_registry`, `service_instance`, and `service_selector` as future architecture vocabulary and mapped current static single-process composition, startup wiring, direct dispatch, route handlers, and module handlers to deferred discovery concepts.

The safe next step is source-first inspection. Agents need a concrete repository command that summarizes service discovery vocabulary and current single-process mapping before any registry storage, selectors, node identity, membership, topology, RPC, remote calls, protocol carriers, persistence, or dependencies are authorized.

## Decision

Implement `node tools/vibit inspect pitaya-discovery --json` as the source-first Pitaya-aligned service discovery map for `M-181/W-0253`.

The command reports:

- ADR-0160 as the source gate and ADR-0161 as the implementation decision.
- `runtime.pitaya_aligned_service_discovery_source_first_map` as the check rule.
- Allowed service discovery vocabulary: `service_discovery`, `service_registry`, `service_instance`, and `service_selector`.
- Related vocabulary: `frontend_server`, `backend_server`, `server_to_server_rpc`, `remote_call`, `route_handler`, `module_handler`, `static_process_composition`, `distributed_group`, and `room_broadcast`.
- Current single-process mappings for discovery, registry, instance, selector, route handler, module handler, distributed group, and room broadcast planning.
- Explicit false deferrals for service discovery implementation, registries, selectors, node identity, RPC implementation, remote calls, frontend/backend role implementation, distributed runtime implementation, distributed groups, room broadcast fanout, cluster-safe routing, protocol, generated output, persistence, dependency, hosted, SDK, and direct compatibility surfaces.
- Redaction flags for credential, token, digest, key, DSN, database payload, local secret file, node credential, and transport metadata contents.
- `M-182/W-0254 Define Pitaya-aligned distributed group and broadcast boundary gate` as the next-ready follow-up.

## Alternatives Considered

- Implement service discovery immediately.
- Add service registries, selectors, or node identity while adding the map.
- Keep the service discovery mapping only in ADR-0160 without a tool inspection surface.
- Fold service discovery vocabulary into `node tools/vibit inspect pitaya-rpc --json` instead of adding a focused command.

## Rationale

Service discovery vocabulary is useful for Pitaya-aligned planning, but it is easy to confuse with permission to add registries, membership, selection algorithms, topology, dependencies, or distributed runtime behavior. A dedicated source-first inspection command gives agents a precise place to inspect vocabulary, current single-process mappings, deferrals, and redaction posture before any implementation gate.

This preserves vibit's agent-native model: inspectable source surfaces, explicit check rules, bounded vocabulary, and continuation memory before runtime behavior.

## Agent Reasoning Summary

The active work item is an inspection-map implementation. The correct continuation is to add the `tools/vibit` command, repository check rule, change artifacts, ADR, and memory updates while preserving all service discovery implementation, registry, selector, node identity, RPC, remote-call, topology, protocol, persistence, dependency, hosted, SDK, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  service_discovery_mapping_clarity: high
  implementation_boundedness: high
  inspection_surface_value: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect Pitaya-aligned service discovery vocabulary without reading every architecture document.

This decision does not add service discovery implementation, service registries, service selectors, node identity, server-to-server RPC implementation, remote call behavior, frontend/backend server role implementation, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

Future distributed group and room broadcast planning must start with W-0254 as a boundary gate rather than implementation.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete discovery or distributed runtime model;
- the discovery inspection output creates confusion with public API compatibility;
- static composition, direct dispatch, route handler, module handler, or service discovery ownership changes enough to require remapping;
- future Pitaya-aligned planning needs separate inspection surfaces for distributed groups, room broadcast, or cluster-safe routing.

## Follow-Up

- Complete `W-0254`: define a Pitaya-aligned distributed group and broadcast boundary gate.
- Keep service discovery implementation, registries, selectors, node identity, server-to-server RPC implementation, remote calls, frontend/backend role implementation, distributed runtime behavior, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint expansion, protocol shape, generated output, persistence, dependencies, hosted deployment, SDKs, and direct compatibility behind later bounded work items.
