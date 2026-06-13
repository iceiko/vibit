# ADR-0215: Pitaya-Aligned Frontend Message Forwarding Boundary Gate

Status: Accepted
Date: 2026-06-13
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-06-13-define-pitaya-aligned-frontend-message-forwarding-boundary-gate/`

Related conversations:

- `conversations/2026-06-13-define-pitaya-aligned-frontend-message-forwarding-boundary-gate.md`

Related artifacts:

- `docs/pitaya-aligned-frontend-message-forwarding-boundary-gate.md`
- `docs/pitaya-aligned-frontend-message-forwarding-boundary-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

The maintainer confirmed a ten-step continuation to move vibit closer to Pitaya-class architecture planning after `ADR-0210` completed the currency wallet protocol route implementation. The sequence is source-first and must not implement distributed runtime behavior, RPC behavior, frontend/backend behavior, service discovery behavior, protocol changes, migrations, dependencies, or direct Nakama/Pitaya API compatibility.

## Decision

Define `docs/pitaya-aligned-frontend-message-forwarding-boundary-gate.md` and register `runtime.pitaya_aligned_frontend_message_forwarding_boundary_gate` as the gate-only repository check rule for `M-235/W-0307`.

This decision records `W-0307` as completed, with `ADR-0215` as the decision and `runtime.pitaya_aligned_frontend_message_forwarding_boundary_gate` as the check rule.

This decision does not add frontend message forwarding behavior, backend targeting behavior, frontend/backend role runtime behavior, runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Implement frontend message forwarding behavior immediately.
- Add concrete distributed runtime APIs or direct Pitaya-compatible names.
- Return to broad product module expansion without source-first Pitaya planning.
- Defer the confirmed ten-step sequence until after another product slice.

## Rationale

Frontend Message Forwarding is useful Pitaya-class vocabulary, but it can easily imply distributed runtime or RPC behavior. A narrow gate lets agents reason about future architecture without changing the running server.

## Agent Reasoning Summary

The correct continuation is to add bounded docs, ADRs, change artifacts, manifests, check rules, and inspection output while preserving all runtime, protocol, generated, persistence, dependency, hosted, SDK, release, distributed runtime, and direct compatibility deferrals.

## Decision Weights

```yaml
decision_weights:
  pitaya_alignment_value: high
  source_first_clarity: high
  implementation_boundedness: high
  runtime_behavior_risk: none_in_this_step
  dependency_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

Agents can now inspect or reason about Pitaya-aligned frontend message forwarding vocabulary without treating that vocabulary as live runtime behavior.

The repository remains single-process and source-first for this concern.

## Reversal Conditions

Revisit this decision if a later ADR selects a concrete frontend message forwarding, distributed runtime, service discovery, RPC, frontend/backend, event bus, protocol, persistence, hosted, or direct compatibility implementation model.

## Follow-Up

- Continue through the confirmed ten-step Pitaya sequence until `W-0312` is complete.
- Keep runtime behavior, protocol changes, generated output, persistence changes, dependencies, hosted deployment, SDKs, release artifacts, distributed runtime behavior, and direct compatibility behind later bounded work items.

## Verification Markers

```yaml
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
dependency_added: false
distributed_runtime_implementation_added: false
direct_nakama_pitaya_api_compatibility_added: false
```
