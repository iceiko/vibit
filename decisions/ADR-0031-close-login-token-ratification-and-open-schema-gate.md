# ADR-0031: Close Login Token Ratification And Open Schema Gate

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-close-login-method-token-format-ratification-milestone/`

Related conversations:

- `conversations/2026-05-14-login-method-token-format-ratification-closeout.md`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/reference.yaml`
- `docs/login-method-token-format-ratification.md`
- `docs/login-method-token-format-ratification.zh-CN.md`
- `docs/credential-token-session-schema-gates.md`
- `docs/selected-login-token-boundary-checks.md`
- `ADR-0024`
- `ADR-0025`
- `ADR-0026`
- `ADR-0027`
- `ADR-0028`
- `ADR-0029`
- `ADR-0030`

## Context

M-013 exists to ratify the first login method and token posture before implementation. Its work queue is complete:

- First login-method candidates were compared.
- `device_credential_login` was ratified as the first login-method set.
- Token format and carrier options were compared.
- Opaque high-entropy access tokens, login-command token issuance, and explicit request proof payloads were ratified as the first posture.
- Token lifecycle and storage implications were defined.
- Authentication contract, error, permission, and audit surfaces were defined.
- Credential, token verifier, external identity, runtime session, and audit schema gates were defined.
- Selected login/token repository checks were added.

The remaining risk is that closing the milestone could accidentally become permission to implement authentication. It must instead open the next bounded gate.

## Decision

Close `M-013 Login Method And Token Format Ratification`.

Open `M-014 Credential And Token Verifier Schema Ratification`.

The first next work item is:

```text
W-0074 Define credential record schema boundary
```

M-014 is a schema-ratification gate. It may define credential and token verifier record semantics, persistence ownership, migration prerequisites, repository boundaries, adapter prerequisites, redaction rules, and verification expectations. It must not add migrations, repositories, PostgreSQL adapters, runtime lookup behavior, login handlers, token validators, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket routes, WebSocket handshake auth, new auth dependencies, or production authentication behavior unless a later work item explicitly authorizes those steps.

## Alternatives Considered

- Start runtime authentication immediately after M-013.
- Create a maintainer confirmation gate before any schema work.
- Ratify credential and token verifier schemas in one large work item.
- Add credential and token migrations directly.
- Start with WebSocket proof carrier or Protobuf authentication message work.
- Copy Nakama's authentication/session API shape.
- Adopt a Pitaya-style session binding implementation before schema ratification.

## Rationale

The selected posture requires persistent credential verifier records and token verifier records before production authentication can be implemented safely. Schema is the next correct layer because vibit's constitution requires public behavior and data contracts to precede code, and because the user explicitly prefers preparation over rushing into the first minimal example.

Credential records and token verifier records are related but separable. Starting with credential schema boundary keeps the next step narrow and lets later work define token verifier schema with the credential relationship visible.

Nakama and Pitaya remain useful references: Nakama for login/session capability coverage, Pitaya for Go server session vocabulary. Neither should determine vibit's database shape or public API.

## Agent Reasoning Summary

Closing M-013 should move the project from selected posture into schema ratification, not runtime implementation. This keeps the project self-bootstrapping, verifiable, and agent-safe while advancing toward production authentication in the correct order.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  schema_first_discipline: high
  security_boundary_clarity: high
  implementation_deferral: high
  reversibility: high
  game_backend_reference_alignment: medium
  default_verification_cost: low
  long_term_maintainability: high
confidence: high
```

## Consequences

- M-013 is completed.
- M-014 becomes active.
- Future continuation begins with credential record schema boundary work.
- Production authentication remains unimplemented.
- Credential and token verifier migrations remain forbidden until schema ratification creates a later migration work item.
- Runtime handlers, WebSocket routes, Protobuf changes, WebSocket handshake auth, generated authentication output, and auth dependencies remain deferred.

## Reversal Conditions

Revisit this decision if:

- A security review requires token verifier schema to be ratified before credential schema.
- The maintainer decides to prioritize a non-authentication game-server capability before credential/token schema work.
- Direct compatibility with Nakama, Pitaya, or another framework is explicitly ratified.
- A later implementation plan proves that credential and token verifier schema must be ratified together to preserve atomicity.

## Follow-Up

- Define credential record schema boundary in W-0074.
- Define token verifier record schema boundary after credential schema boundary.
- Continue to keep migrations, repositories, adapters, runtime lookup, handlers, routes, generated authentication output, Protobuf changes, WebSocket changes, and authentication implementation behind separate gates.
