# ADR-0037: Close Runtime Auth Boundary And Open Generated Shape Gate

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-close-runtime-authentication-implementation-boundary-planning-milestone/`

Related conversations:

- `conversations/2026-05-14-runtime-authentication-boundary-closeout.md`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/runtime-authentication-implementation-boundary.md`
- `decisions/ADR-0036-runtime-authentication-implementation-boundary.md`

## Context

`M-016 Runtime Authentication Implementation Boundary Planning` defined the first runtime authentication implementation boundary after the authentication PostgreSQL adapter.

The milestone clarified application ownership, repository usage, persistence-only PostgreSQL adapter rules, request identity handoff, error and permission mapping, Nakama/Pitaya alignment, and separately gated future implementation steps. It did not add runtime authentication behavior.

The next useful step should continue preparation without collapsing into login or token validation implementation.

## Decision

Close `M-016 Runtime Authentication Implementation Boundary Planning`.

Open `M-017 Authentication Generated Contract Shape Gate`.

The first next work item is:

```text
W-0088 Decide authentication generated contract shape timing
```

`M-017` may decide whether generated Go authentication contract shapes are required before application service interfaces and runtime behavior. It may plan generator/check updates for authentication contract shapes. It must not implement login, token generation, verifier comparison, token validation, logout execution, cleanup jobs, Protobuf authentication messages, WebSocket proof carriers, generated authentication shapes, authentication dependencies, repository interface changes, or migration schema changes unless a later bounded work item explicitly authorizes that exact slice.

## Alternatives Considered

- Start `AuthenticateWithDeviceCredential` runtime behavior immediately.
- Start token verifier algorithm selection immediately.
- Add Protobuf authentication messages before deciding generated contract shape timing.
- Add generated authentication contract shapes in the closeout change.
- Ask for maintainer confirmation before a non-behavioral generated-shape timing decision.

## Rationale

Generated contract shapes are the next non-behavioral question in the W-0086 queue. Deciding their timing first keeps future runtime authentication code tied to semantic contracts and reduces the chance that agents invent ad hoc application payloads.

This also preserves vibit's "generated shape, handwritten logic" principle without changing generated file conventions or adding generated files implicitly.

## Agent Reasoning Summary

The project should continue moving, but the next step should still be preparation. Generated authentication contract shape timing is the narrowest next gate because it affects future runtime code structure without adding authentication behavior.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  schema_before_code: high
  generated_shape_discipline: high
  security_boundary_clarity: high
  implementation_cost: low
  reversibility: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- `M-016` is completed.
- `M-017` becomes active.
- `W-0088` becomes the next ready work item.
- Runtime authentication behavior remains deferred.
- Generated authentication shapes remain deferred until a later work item explicitly authorizes them.
- Protobuf messages, WebSocket proof carriers, and authentication dependencies remain deferred.

## Reversal Conditions

Revisit this decision if:

- A future implementation plan proves that verifier algorithm selection must precede generated shape timing.
- Generated authentication contract shapes are found to conflict with the existing generated-output standard.
- The maintainer explicitly chooses Protobuf authentication messages, WebSocket carriers, or runtime login behavior as the next priority.

## Follow-Up

- Advance W-0088.
- Keep generated file changes, runtime code, protocol code, transport code, and dependency adoption behind separate gates.
