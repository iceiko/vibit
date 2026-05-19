# ADR-0084: Protocol Session Carrier Functional Slice

Status: Accepted
Date: 2026-05-20
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-20-define-protocol-session-carrier-functional-slice/`

Related conversations:

- `conversations/2026-05-20-protocol-session-carrier-functional-slice.md`

Related artifacts:

- `runtime/internal/platform/protocol/protobuf/inventory_bridge.go`
- `runtime/internal/platform/protocol/protobuf/authentication_bridge.go`
- `runtime/internal/platform/protocol/protobuf/inventory_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_request_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `rules/check-rules.json`

## Context

`ADR-0068` made successful device-credential login create a durable runtime session row and return the created session id in the application-owned `AuthenticationResult`. That session id intentionally was not protocol-visible at the time because the protocol session carrier boundary had not been selected.

`ADR-0082` later reduced confirmation-gate density for non-security functional slices. `ADR-0083` then completed the server-observed connection epoch primitive. The next lifecycle closure gap was protocol visibility: clients could receive an access token but did not have a bounded way to receive the runtime session id created by login.

## Decision

Select:

```text
define_protocol_session_carrier_functional_slice
```

as a Tier 2 functional slice and implement the smallest checkable behavior directly.

The Protobuf adapter now reuses existing `vibit.protocol.v1.Envelope.Session` metadata for response session carriers:

- successful `runtime.authentication.AuthenticateWithDeviceCredential` responses carry the created runtime `session_id`,
- those responses carry the authenticated `player_id`,
- existing server-observed `connection_id` and `connection_epoch` remain available when already present in request/frame metadata,
- response session metadata may also derive from already validated application identity,
- metadata-only identity remains metadata-only and is not upgraded.

The implementation is in:

```text
runtime/internal/platform/protocol/protobuf/inventory_bridge.go
runtime/internal/platform/protocol/protobuf/authentication_bridge.go
```

## Boundaries

This ADR keeps these boundaries:

- Protobuf adapter owns only wire metadata mapping.
- Authentication service owns login execution and runtime session creation composition.
- Runtime `session_id` is client-visible metadata after successful login; `session_id is not proof` by itself.
- WebSocket transport remains credential-neutral and does not authenticate handshakes.
- Application route policy remains the owner of whether a route requires request-token proof, bound connection identity, session-validated identity, or bound-session identity.

This decision does not add new Protobuf fields, generated output, reconnect tokens, resume tokens, WebSocket handshake authentication, durable/distributed session routing, close code mapping, close reason text, logout-triggered socket close, runtime session revocation, presence lifecycle, operations/admin disconnect, dependencies, broad product modules, or direct Nakama/Pitaya API compatibility.

## Nakama And Pitaya Mapping

Nakama informs the product pressure: authentication should yield session material that clients can carry into later realtime lifecycle behavior.

Pitaya informs the architecture pressure: session context should be available to routing/session layers without moving authentication into the frontend acceptor.

vibit adapts those lessons by exposing only the already-created runtime session id through existing envelope metadata. It does not copy either project's public API.

## Alternatives Considered

- Add new fields to `AuthenticateWithDeviceCredentialResponse`.
- Add a new Protobuf session message.
- Change the envelope schema.
- Add a first-message session binding route.
- Add WebSocket handshake authentication.
- Treat the session id as proof for route authorization.
- Jump directly to presence lifecycle without making runtime sessions visible.

## Rationale

The smallest useful protocol carrier is to reuse the envelope session metadata that already exists. It avoids generated-output churn, keeps the public wire shape stable, and gives clients the runtime session id produced by the server after successful login.

This is enough to close the immediate lifecycle gap before presence work: clients can now observe the server-created session id, while future route policy and session validation still decide when that metadata becomes useful for authorization.

## Agent Reasoning Summary

The maintainer asked for faster progress toward Nakama/Pitaya-class capability. The useful compromise is not another gate and not a broad session feature. The bounded step is to expose already-existing runtime session metadata in the existing envelope response path, then move on to presence lifecycle work.

## Decision Weights

```yaml
decision_weights:
  development_velocity: high
  lifecycle_closure: high
  protocol_shape_stability: high
  transport_auth_separation: high
  route_policy_safety: high
  direct_api_compatibility: low
  reconnect_resume_now: low
confidence: high
```

## Consequences

- `runtime.protocol_session_carrier_functional_slice` becomes the repository check rule for this slice.
- Successful login responses can carry runtime session id through `Envelope.Session`.
- Response envelopes can derive session metadata from validated identity where available.
- Metadata-only identity still does not become proof.
- Future presence lifecycle work can rely on clients having an explicit session id carrier without adding handshake auth or reconnect/resume behavior.

## Reversal Conditions

Revisit this decision if a future ADR replaces the envelope session model, if direct Nakama/Pitaya API compatibility is explicitly selected and requires different fields, or if security review decides that client-visible runtime session ids require a different carrier posture.

## Follow-Up

- Define the first presence lifecycle functional slice.
- Keep reconnect/resume tokens, handshake authentication, route-level session proof changes, logout-triggered close, and runtime session revocation behind explicit future work.
