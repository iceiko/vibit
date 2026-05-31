# ADR-0155: Pitaya-Aligned Distributed Runtime Vocabulary Source-First Map

Status: Accepted
Date: 2026-05-31
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-implement-pitaya-aligned-distributed-runtime-vocabulary-source-first-map/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-distributed-runtime-vocabulary-source-first-map.md`

Related artifacts:

- `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`
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

`ADR-0154` reactivated a narrow Pitaya-aligned distributed runtime vocabulary for future architecture planning. The maintainer asked to keep moving toward Pitaya, but the repository still needs implementation boundaries to stay source-first and single-process until later gates authorize behavior.

The safe next step is to make the vocabulary inspectable through the repository tooling, not to implement frontend/backend server roles, RPC/remotes, service discovery, distributed groups, broadcast fanout, or cluster-safe session routing.

## Decision

Implement `node tools/vibit inspect pitaya-vocabulary --json` as the source-first repository map for Pitaya-aligned distributed runtime vocabulary.

Register `runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map` as the repository check rule.

The inspection output records:

- the ADR-0154 allowed vocabulary: `acceptor`, `frontend_server`, `backend_server`, `route_handler`, `session_binding`, `server_to_server_rpc`, `remote_call`, `service_discovery`, `distributed_group`, `room_broadcast`, and `cluster_safe_session_routing`;
- the current single-process mapping for WebSocket acceptors, first-message connection binding, application dispatch/protocol bridge, deferred frontend/backend roles, deferred RPC/remotes, deferred groups/broadcast, and deferred service discovery;
- source surfaces under `.arch/`, docs, `tools/vibit`, and `rules/check-rules.json`;
- explicit false flags for distributed runtime implementation, frontend/backend role implementation, RPC/remotes, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, and direct compatibility;
- redaction flags proving the map does not expose raw credentials, raw tokens, digests, verifier keys, DSNs, database payloads, or local secret file contents.

Complete `M-175/W-0247` and open `M-176/W-0248 Define Pitaya-aligned frontend/backend role boundary gate` as next-ready.

This decision does not add distributed runtime behavior, frontend/backend server role implementation, server-to-server RPC, remote call behavior, service discovery, distributed groups, room broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement a Pitaya-style frontend/backend role split immediately.
- Add server-to-server RPC or service discovery now.
- Keep the vocabulary only in ADR-0154 without a tool inspection surface.
- Extend `inspect operations` with more Pitaya fields instead of a dedicated vocabulary command.

## Rationale

A dedicated source-first inspection command gives agents a stable place to inspect the Pitaya vocabulary without confusing it with runtime operations behavior. It also keeps future distributed architecture planning explicit, checkable, and redacted.

This preserves vibit's agent-native model: vocabulary, mapping, deferrals, source surfaces, checks, and memory before implementation.

## Agent Reasoning Summary

The active work item is an implementation of an inspection map, not a distributed runtime slice. The correct continuation is to add the `tools/vibit` command, repository check coverage, change artifacts, ADR, and continuation pointers. The follow-up should define frontend/backend role vocabulary boundaries before any role implementation.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  architecture_clarity: high
  implementation_boundedness: high
  inspection_surface_value: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `node tools/vibit inspect pitaya-vocabulary --json` becomes the source-first Pitaya vocabulary map.
- `runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map` becomes the check rule for W-0247.
- `M-175/W-0247` is completed.
- `M-176/W-0248 Define Pitaya-aligned frontend/backend role boundary gate` becomes next-ready.
- Distributed runtime implementation, frontend/backend role implementation, RPC/remotes, service discovery, distributed groups, broadcast fanout, cluster-safe session routing, protocol changes, generated output, persistence, dependencies, hosted deployment, SDKs, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete distributed runtime implementation model;
- the command output creates confusion with public API compatibility;
- single-process session, routing, or group semantics change enough to require remapping;
- future Pitaya-aligned planning needs separate inspect commands for topology, routing, groups, or service discovery.

## Follow-Up

- Complete `W-0248`: define a Pitaya-aligned frontend/backend role boundary gate.
- Keep frontend/backend role implementation, cluster/RPC/service-discovery implementation, distributed groups, cluster-safe session routing, runtime endpoint expansion, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
