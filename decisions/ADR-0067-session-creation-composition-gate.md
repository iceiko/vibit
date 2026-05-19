# ADR-0067: Session Creation Composition Gate

Status: Accepted
Date: 2026-05-17
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-runtime-session-validation-implementation/`
- `changes/2026-05-17-define-session-creation-composition-gate/`

Related conversations:

- `conversations/2026-05-17-session-creation-composition-gate.md`

Related artifacts:

- `docs/session-creation-composition-gate.md`
- `docs/session-creation-composition-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`W-0144` implemented an application-owned persistent runtime session validator after durable `runtime_sessions`, the storage-neutral `session.Repository`, and the PostgreSQL adapter were already in place.

The work queue reached `M-073/W-0145`, a confirmation gate. The maintainer asked the agent to recommend the next ten steps and continue, with Nakama and Pitaya as key reference baselines.

The implementation gap is now explicit: login can issue an access token, and persisted session validation exists, but no production path creates durable runtime session rows.

## Decision

Select:

```text
define_session_creation_composition_gate
```

Define a gate-only session creation composition standard in:

```text
docs/session-creation-composition-gate.md
docs/session-creation-composition-gate.zh-CN.md
```

The future composition owner is:

```text
runtime/internal/app
```

The first future composition candidate is:

```text
runtime/internal/app/authentication.AuthenticateWithDeviceCredential
```

Future session creation may use:

```text
runtime/internal/app/session.Repository.CreateRuntimeSession
```

through an application unit-of-work capability after a later implementation gate authorizes code.

This ADR does not add session creation code, modify the authentication service, generate session ids, create sessions at login or `BindConnection`, change runtime session validation, change route policy, change WebSocket handshake authentication, add transport credential carriers, add Protobuf session messages, change the existing Protobuf envelope, add generated output, add logout/revocation active-connection invalidation, add reconnect/epoch behavior, add cleanup jobs, add dependencies, add memory durable session behavior, or adopt direct Nakama/Pitaya public API compatibility.

## Alternatives Considered

- Implement login-created runtime sessions immediately.
- Create runtime sessions during `BindConnection`.
- Treat first-message bound identity as the session creation point.
- Make a persisted session row satisfy route policy as soon as it is created.
- Put session creation behavior inside the PostgreSQL adapter.
- Put session creation behavior inside WebSocket transport or Protobuf protocol adapters.
- Define logout/revocation or reconnect behavior before creation composition.
- Copy Nakama or Pitaya session APIs directly.

## Rationale

Nakama shows that authentication and session creation sit on a lifecycle path with expiration, refresh, logout, and session management implications. That pressure should shape vibit's session creation order, token linkage, and commit semantics.

Pitaya shows that acceptors, session context, and handlers should remain separate. The useful vibit adaptation is to keep durable session creation in application composition, not in transport or protocol code.

Defining the gate before implementation prevents later agents from accidentally treating session creation as route authorization, transport authentication, logout disconnect, or reconnect policy.

## Agent Reasoning Summary

After runtime session validation, the highest-leverage next step is not route policy. Route policy cannot reasonably depend on session validation until there is a defined way to create durable sessions. Session creation also crosses token issuance, unit-of-work composition, session id generation, redaction, and protocol-carrier decisions.

The gate gives the project a precise next implementation target while preserving vibit's separation between authentication proof validation, durable session lifecycle storage, request identity construction, transport connection plumbing, and protocol routing.

## Decision Weights

```yaml
decision_weights:
  lifecycle_correctness: high
  nakama_pitaya_alignment: high
  unit_of_work_atomicity: high
  transport_protocol_app_separation: high
  future_route_policy_clarity: high
  immediate_feature_delivery: medium
  direct_nakama_pitaya_api_compatibility: low
confidence: high
```

## Consequences

- `docs/session-creation-composition-gate.md` and `docs/session-creation-composition-gate.zh-CN.md` exist.
- `runtime.session_creation_composition_gate` becomes the repository check rule for this slice.
- `M-073/W-0145` is closed with `define_session_creation_composition_gate`.
- `M-074/W-0146` is completed as a gate-only milestone.
- The work queue blocks again at the next confirmation gate before session creation implementation, route-policy use of session or bound identity, logout/revocation active-connection behavior, reconnect/epoch behavior, operations, memory durable session behavior, direct Nakama/Pitaya compatibility, or broader game backend expansion.
- Existing runtime behavior is not changed by this ADR.

## Reversal Conditions

Revisit this decision if:

- A future ADR chooses handshake-level authentication as the primary creation point for runtime sessions.
- A future ADR chooses `BindConnection` rather than login as the first durable session creation trigger.
- A future token/session lifecycle ADR requires sessions to outlive or rotate independently from access tokens.
- Route-policy requirements force a different session creation and validation composition model.
- A later ADR adopts direct Nakama or Pitaya public API compatibility.

## Follow-Up

- Implement session creation composition in a bounded slice before any production path expects durable sessions to exist after login.
- Define bound/session identity route policy before protected routes can rely on session-validated identity.
- Define logout/revocation active-connection behavior before revocation closes or invalidates WebSocket connections.
- Define reconnect and connection epoch behavior before duplicate replacement or resume behavior.
