# ADR-0057: First Message Connection Binding Gate

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-session-handshake-ratification/`
- `changes/2026-05-17-define-first-message-connection-binding-gate/`

Related conversations:

- `conversations/2026-05-17-session-handshake-next-direction-and-first-message-binding-gate.md`

Related artifacts:

- `docs/first-message-connection-binding-gate.md`
- `docs/first-message-connection-binding-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0056` ratified request-level access-token validation as the current path, kept WebSocket transport credential-neutral, kept the existing Protobuf envelope unchanged, and named first-message protocol/application binding as the preferred future connection-level gate.

The work queue then reached `M-051/W-0123`, asking for the next major direction. The maintainer asked the agent to recommend the next ten steps and continue according to that recommendation, while continuing to treat Nakama and Pitaya as key references.

## Decision

Select `define_first_message_connection_binding_gate` and define the gate-only standard for future first-message connection binding.

The selected future route is:

```text
runtime.authentication.BindConnection
```

It is a `system` message, not a domain command. The future payload candidates are:

```text
vibit.authentication.v1.BindConnectionRequest
vibit.authentication.v1.BindConnectionResponse
```

The future bind message carries opaque access-token proof in a Protobuf payload inside the existing WebSocket frame loop. The WebSocket transport continues to ignore credential carriers in HTTP headers, bearer values, cookies, query strings, and subprotocols. The existing Protobuf envelope remains unchanged.

This ADR does not implement connection binding, add `.proto` messages, generate Protobuf output, create a connection-bound identity registry, change route policy, add session persistence, add migrations, change repositories, add dependencies, implement logout/revocation, implement reconnect/epoch behavior, or adopt direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Parse access tokens in WebSocket handshake headers, cookies, query strings, or subprotocols.
- Add authentication proof or binding fields to the existing Protobuf envelope.
- Treat connection binding as a gameplay domain command.
- Add session persistence before defining connection binding behavior.
- Require binding as the first frame before allowing public login.
- Implement single-socket replacement, kick, reconnect, or epoch behavior in the first gate.
- Adopt Nakama or Pitaya public API compatibility directly.

## Rationale

Nakama guides the capability sequence: authenticate first, then use realtime socket features with authenticated state. Nakama also shows that session/token lifecycle and active socket lifecycle are related but not identical.

Pitaya guides the architecture vocabulary: connection acceptors, session binding, handlers, groups, and cluster behaviors are separate concerns. Pitaya exposes session binding and session lookup concepts, but vibit should adapt them into explicit protocol/application boundaries.

The first-message approach uses the already ratified WebSocket Protobuf loop without binding authentication to browser-specific handshake carrier behavior. It also leaves space for clients to call public login on an anonymous connection before binding the connection with the newly issued access token.

## Agent Reasoning Summary

After public login, request-level route protection, and session/handshake ratification, the next high-value step is to specify how a connection becomes associated with validated identity. This must be done before presence, rooms, parties, matches, groups, reconnect, or active-connection invalidation can be designed coherently.

## Decision Weights

```yaml
decision_weights:
  transport_credential_neutrality: high
  compatibility_with_existing_request_level_validation: high
  future_presence_room_match_foundation: high
  nakama_pitaya_alignment: high
  implementation_risk_reduction: high
  direct_feature_delivery: medium
  reversibility: medium
confidence: high
```

## Consequences

- `docs/first-message-connection-binding-gate.md` becomes the standard for future first-message connection binding.
- `runtime.first_message_connection_binding_gate` becomes the repository check rule.
- `M-051/W-0123` is closed with the selected direction.
- `M-052/W-0124` is completed as a gate-only milestone.
- The work queue blocks again at `M-053/W-0125` before implementation or another major direction.
- Existing runtime behavior is not changed by this ADR.

## Reversal Conditions

Revisit this decision if:

- A future client compatibility requirement requires handshake-level authentication before protocol binding.
- The existing envelope must version before system binding messages can be carried.
- A future durable session schema makes connection binding depend on session records first.
- A later ADR adopts direct Nakama or Pitaya public API compatibility.
- Operational load from anonymous WebSocket connections becomes a primary near-term constraint.

## Follow-Up

- Define a first-message connection binding implementation gate before adding Protobuf messages, generated output, connection registries, route policy changes, or startup composition.
- Define PostgreSQL session persistence before durable session ids or session validation become production behavior.
- Define logout/revocation active-connection behavior before logout can invalidate bound connections.
- Define reconnect and connection epoch behavior before duplicate connection replacement or resume behavior.
