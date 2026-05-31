# ADR-0154: Pitaya-Aligned Distributed Runtime Vocabulary Reactivation Gate

Status: Accepted
Date: 2026-05-31
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`

Related artifacts:

- `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md`
- `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.zh-CN.md`
- `decisions/ADR-0153-minimum-operations-inspection-source-first-surface-implementation.md`
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

`ADR-0153` implemented the source-first operations inspection command and recorded a Pitaya deferred architecture map. The maintainer asked to continue multiple steps and move closer to Pitaya.

The safe next step is not distributed runtime implementation. The current repository posture still keeps Pitaya as a deferred future architecture reference, while Nakama remains the primary product capability reference. However, the project now benefits from a stricter vocabulary gate so future agents can talk about acceptors, frontend/backend roles, RPC/remotes, service discovery, distributed groups, broadcast, and cluster-safe session routing without treating those words as implementation permission.

## Decision

Accept `docs/pitaya-aligned-distributed-runtime-vocabulary-reactivation-gate.md` and its Simplified Chinese translation as the gate for reactivating Pitaya-aligned distributed runtime vocabulary.

Register `runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate` as the repository check rule.

The gate defines:

- allowed vocabulary: `acceptor`, `frontend_server`, `backend_server`, `route_handler`, `session_binding`, `server_to_server_rpc`, `remote_call`, `service_discovery`, `distributed_group`, `room_broadcast`, and `cluster_safe_session_routing`;
- current single-process mapping for WebSocket acceptors, first-message connection binding, application dispatch/protocol bridge, deferred frontend/backend roles, deferred RPC/remotes, deferred groups/broadcast, and deferred service discovery;
- ownership for architecture vocabulary and a future source-first vocabulary map;
- stop conditions for any implementation behavior.

Complete `M-174/W-0246` and open `M-175/W-0247 Implement Pitaya-aligned distributed runtime vocabulary source-first map` as next-ready.

This decision does not add distributed runtime behavior, frontend/backend server roles, server-to-server RPC, remote calls, service discovery, distributed groups, broadcast fanout, cluster-safe session routing, runtime endpoint behavior, metrics endpoints, observability pipelines, dashboards, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, SDK publication, hosted deployments, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Jump directly to a Pitaya-style frontend/backend role implementation.
- Add server-to-server RPC or service discovery now.
- Add distributed groups or broadcast fanout now.
- Keep Pitaya vocabulary hidden until a much later distributed runtime milestone.
- Fold the vocabulary into the operations inspection command only.

## Rationale

The maintainer's Pitaya direction is valid architecture pressure, but the repository's current maturity still requires single-process semantics to stay stable. A vocabulary gate lets agents make future distributed architecture discussions more precise while preventing accidental behavior expansion.

This preserves vibit's agent-native model: define vocabulary, ownership, deferrals, stop conditions, checks, and memory before implementation.

## Agent Reasoning Summary

The active work item is a gate. The correct continuation is to write the standard, ADR, change artifacts, repository checks, and manifest updates. The next implementation slice should be a source-first vocabulary map, not cluster/RPC/service-discovery behavior.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  architecture_clarity: high
  implementation_boundedness: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

- Pitaya vocabulary is reactivated for future architecture planning only.
- `runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate` becomes the check rule for W-0246.
- `M-174/W-0246` is completed.
- `M-175/W-0247 Implement Pitaya-aligned distributed runtime vocabulary source-first map` becomes next-ready.
- Distributed runtime implementation, frontend/backend server roles, RPC/remotes, service discovery, distributed groups, broadcast fanout, cluster-safe session routing, protocol changes, generated output, persistence, dependencies, hosted deployment, SDKs, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- a later architecture ADR selects a concrete distributed runtime implementation model;
- the vocabulary creates confusion with public API compatibility;
- single-process session, routing, or group semantics change enough to require remapping;
- prototype feedback shows a stronger near-term need for a different product capability family.

## Follow-Up

- Complete `W-0247`: implement a source-first Pitaya-aligned distributed runtime vocabulary map.
- Keep cluster/RPC/service-discovery implementation, distributed groups, cluster-safe session routing, runtime endpoint expansion, metrics, dashboards, SDKs, hosted deployment, and direct compatibility behind later bounded work items.
