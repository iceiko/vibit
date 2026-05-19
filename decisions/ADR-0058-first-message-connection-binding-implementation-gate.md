# ADR-0058: First Message Connection Binding Implementation Gate

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-first-message-connection-binding-gate/`
- `changes/2026-05-17-define-first-message-connection-binding-implementation-gate/`

Related conversations:

- `conversations/2026-05-17-first-message-binding-next-direction-and-implementation-gate.md`

Related artifacts:

- `docs/first-message-connection-binding-implementation-gate.md`
- `docs/first-message-connection-binding-implementation-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0057` selected first-message connection binding through a future `runtime.authentication.BindConnection` system route. It intentionally did not authorize Protobuf source changes, generated output, connection binding registries, route-policy use of bound identity, startup composition, or session persistence.

The work queue then reached `M-053/W-0125`. The maintainer again asked the agent to recommend the next ten steps and proceed according to the recommendation, while continuing to treat Nakama and Pitaya as key references.

## Decision

Select:

```text
define_first_message_connection_binding_implementation_gate
```

Define a gate-only standard for a later bounded implementation slice.

The future implementation slice may add:

- `BindConnectionRequest`, `BindConnectionResponse`, and `ConnectionBindingStatus` to `proto/vibit/authentication/v1/authentication.proto`.
- Regenerated Go Protobuf output through Buf.
- A Protobuf protocol adapter bridge for `runtime.authentication.BindConnection`.
- Application-owned connection binding types and optional process-local registry.
- PostgreSQL startup composition for the binder when the existing authentication service is composed.
- Focused tests for success, failure collapse, redaction, transport neutrality, envelope stability, and route-policy compatibility.

This ADR does not implement `BindConnection`, add `.proto` messages, run `buf generate`, create a registry, change route policy, add session persistence, add migrations, change repositories, add dependencies, implement logout/revocation, implement reconnect/epoch behavior, or adopt direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Implement `BindConnection` immediately without a separate implementation gate.
- Add access-token proof to the WebSocket handshake.
- Add binding or proof fields to the existing Protobuf envelope.
- Make bound identity immediately satisfy all protected routes without an explicit route-policy update.
- Add durable session persistence before process-local connection binding.
- Adopt Nakama session/socket API compatibility or Pitaya session APIs directly.

## Rationale

Nakama shows that authenticated server features and realtime socket usage need a coherent session/socket lifecycle. Its model is useful as a capability reference, but vibit should not adopt Nakama's public API shape or JWT/session semantics by implication.

Pitaya shows a useful separation between acceptors, sessions, route handlers, groups, and cluster concerns. vibit should preserve that layering by keeping WebSocket transport credential-neutral and making binding an explicit protocol/application concern.

A separate implementation gate reduces ambiguity before crossing generated output, protocol adapter, app registry, startup composition, and route-policy boundaries.

## Agent Reasoning Summary

The highest-value next step after selecting the future `BindConnection` route is to define the exact implementation boundary. Without that boundary, adding protocol messages or route behavior would risk mixing transport authentication, session persistence, and route-policy changes in one uncontrolled slice.

## Decision Weights

```yaml
decision_weights:
  transport_credential_neutrality: high
  generated_output_control: high
  route_policy_safety: high
  nakama_pitaya_alignment: high
  implementation_risk_reduction: high
  direct_feature_delivery: medium
  reversibility: medium
confidence: high
```

## Consequences

- `docs/first-message-connection-binding-implementation-gate.md` becomes the standard for the future implementation slice.
- `runtime.first_message_connection_binding_implementation_gate` becomes the repository check rule.
- `M-053/W-0125` is closed with the selected direction.
- `M-054/W-0126` is completed as a gate-only milestone.
- The work queue blocks again at `M-055/W-0127` before actual implementation or another major direction.
- Existing runtime behavior is not changed by this ADR.

## Reversal Conditions

Revisit this decision if:

- A future client compatibility requirement demands handshake-level authentication.
- The existing envelope must be versioned before system binding messages can be represented safely.
- Durable session persistence becomes a prerequisite for any connection binding.
- A future ADR adopts direct Nakama or Pitaya public API compatibility.
- Operational load from anonymous WebSocket connections becomes a primary near-term constraint.

## Follow-Up

- Implement first-message connection binding in a bounded slice if the maintainer selects `implement_first_message_connection_binding`.
- Define PostgreSQL session persistence before durable session ids or session validation become production behavior.
- Define logout/revocation active-connection behavior before logout can invalidate bound connections.
- Define reconnect and connection epoch behavior before duplicate connection replacement or resume behavior.
