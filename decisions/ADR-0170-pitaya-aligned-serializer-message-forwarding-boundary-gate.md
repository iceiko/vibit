# ADR-0170: Pitaya-Aligned Serializer And Message Forwarding Boundary Gate

Status: Accepted
Date: 2026-06-01
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-26-define-pitaya-aligned-serializer-message-forwarding-boundary-gate/`

Related conversations:

- `conversations/2026-05-26-pitaya-aligned-serializer-message-forwarding-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md`
- `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.zh-CN.md`
- `decisions/ADR-0169-select-next-pitaya-aligned-direction-after-route-handler-pipeline-map.md`
- `decisions/ADR-0168-pitaya-aligned-route-handler-pipeline-source-first-map.md`
- `docs/pitaya-aligned-route-handler-pipeline-boundary-gate.md`
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

`ADR-0169` selected `define_pitaya_aligned_serializer_message_forwarding_boundary_gate` as the next bounded Pitaya-aligned direction after the route handler pipeline source-first map.

Pitaya-style serializer and message forwarding vocabulary is useful architecture pressure for future distributed runtime planning. vibit's current runtime, however, has explicit Protobuf bridge functions, application-owned outbound message intent, metadata-only target scopes, and single-process WebSocket delivery.

## Decision

Accept `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md` and its Simplified Chinese translation as the gate for Pitaya-aligned serializer and message forwarding vocabulary.

Register `runtime.pitaya_aligned_serializer_message_forwarding_boundary_gate` as the repository check rule.

Complete `M-190/W-0262` and open `M-191/W-0263 Implement Pitaya-aligned serializer and message forwarding source-first map` as next-ready.

This decision does not add serializer behavior, message forwarding behavior, route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, backend route targeting, cluster-safe session routing behavior, distributed session routing, service discovery implementation, RPC, remote calls, frontend/backend role behavior, distributed runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement serializer behavior or message forwarding immediately.
- Add a source-first serializer and forwarding map without first defining a gate.
- Fold serializer and forwarding vocabulary back into the route handler pipeline map.
- Return directly to backend route targeting, service discovery, RPC, or distributed runtime implementation.

## Rationale

Serializer and message forwarding vocabulary is high-risk because it can be mistaken for permission to add codecs, serializer registries, forwarding workers, backend route targeting, RPC, remote delivery, protocol changes, dependencies, or distributed runtime behavior. Defining a gate first lets vibit preserve Pitaya-aligned planning vocabulary while keeping the concrete runtime single-process and source-first.

## Agent Reasoning Summary

The active work item is a gate. The correct continuation is to write the standard, ADR, change artifacts, repository checks, and manifest updates. The follow-up should implement a source-first serializer and message forwarding map, not serializer behavior, message forwarding behavior, route handlers, backend targeting, RPC, remote calls, protocol changes, generated output, or distributed runtime behavior.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  serializer_forwarding_clarity: high
  implementation_boundedness: high
  distributed_runtime_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `serializer_boundary`, `serializer_format`, `encode_boundary`, `decode_boundary`, `message_forwarding`, `forwarding_target`, `forwarding_envelope`, and `delivery_handoff` are allowed as future architecture vocabulary only.
- `runtime.pitaya_aligned_serializer_message_forwarding_boundary_gate` becomes the check rule for W-0262.
- `M-190/W-0262` is completed.
- `M-191/W-0263 Implement Pitaya-aligned serializer and message forwarding source-first map` becomes next-ready.

## Reversal Conditions

Revisit this decision if a later architecture ADR selects a concrete serializer or forwarding model, if the vocabulary creates confusion with public API compatibility, or if protocol bridge/outbound delivery ownership changes enough to require remapping.

## Follow-Up

- Complete `W-0263`: implement a source-first Pitaya-aligned serializer and message forwarding map.
- Keep serializer behavior, message forwarding behavior, backend route targeting, service discovery, RPC, remote calls, frontend/backend role behavior, distributed runtime behavior, protocol changes, generated output, persistence, dependencies, hosted deployment, SDKs, and direct compatibility behind later bounded work items.
