# ADR-0029: Credential Token Session Schema Gates

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-credential-token-session-schema-gates/`

Related conversations:

- `conversations/2026-05-14-credential-token-session-schema-gates.md`

Related artifacts:

- `docs/credential-token-session-schema-gates.md`
- `docs/credential-token-session-schema-gates.zh-CN.md`
- `docs/token-lifecycle-storage-implications.md`
- `docs/authentication-contract-error-permission-surfaces.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`

## Context

M-013 has ratified the first login method, opaque access-token posture, token lifecycle, storage implications, and semantic authentication contract surfaces. The next risk is that a future agent may treat those decisions as permission to add credential tables, token tables, session tables, player account columns, repositories, migrations, or runtime authentication behavior in one broad implementation step.

W-0071 defines schema gates so future work can proceed in bounded, verifiable stages.

## Decision

Define credential, token verifier, external identity, runtime session, and audit persistence schema gates.

The first posture requires future schema ratification for:

- Credential records for `device_credential_login`.
- Token verifier records for opaque access-token validation, presented-token logout, revocation, expiration, and rotation support.

The first posture does not require:

- External identity link records.
- Runtime session records.
- Refresh-token storage.
- Session-token vocabulary.
- WebSocket connection binding.

PostgreSQL remains the default durable target for the first required credential and token verifier schema gates. Redis-like storage remains deferred.

Player account lifecycle tables remain credential-free, token-free, external-identity-free, session-free, WebSocket-state-free, and request-validation-free.

No migration, repository interface, PostgreSQL adapter, runtime lookup behavior, Protobuf field, WebSocket behavior, handler, route, or authentication implementation is authorized by this decision.

## Alternatives Considered

- Adding credential and token tables immediately.
- Adding only prose notes without manifest markers.
- Treating token verifier storage as a runtime memory concern.
- Combining credential, token verifier, external identity, and session persistence into one future table family.
- Requiring runtime session persistence for the first access-token posture.
- Allowing player account lifecycle tables to carry credential or token state.
- Selecting a Redis-like store before distributed runtime requirements exist.

## Rationale

The selected first posture needs credential verifier storage and token verifier storage before production authentication can exist. Those are security-sensitive schemas and should be designed explicitly before migrations or runtime behavior appear.

External identity linking and runtime session persistence are important game-server capabilities, as seen in Nakama and Pitaya vocabulary, but they are not required for the first `device_credential_login` plus opaque access-token posture. Deferring them keeps the first implementation slice smaller and prevents session vocabulary from being conflated with access-token validation.

PostgreSQL is already the ratified authoritative durable store. It is the correct default target for the first required schema gates until a later distributed runtime or performance decision justifies another store.

## Agent Reasoning Summary

This decision keeps the project self-bootstrapping and controllable. Future agents get enough schema direction to avoid shortcuts, but they still cannot create tables, repositories, migrations, or runtime authentication behavior without a later bounded work item.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  security_boundary_clarity: high
  schema_reversibility: high
  implementation_deferral: high
  production_safety: high
  game_backend_capability_alignment: high
  dependency_load: low
  protocol_stability: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- Future credential schema work must be explicit before `device_credential_login` implementation.
- Future token verifier schema work must be explicit before opaque access-token validation.
- Refresh-token storage remains forbidden for the first posture.
- Runtime session persistence remains deferred.
- External identity linking remains deferred.
- Player account lifecycle storage remains isolated.
- W-0072 can add narrow repository checks for selected login/token boundaries.
- M-013 can close after checks confirm that implementation remains deferred.

## Reversal Conditions

Revisit this decision if:

- The first implementation milestone selects a different login method.
- A security review requires a different token storage model.
- A distributed runtime decision requires Redis-like token/session storage.
- A later compatibility goal with Nakama, Pitaya, or another framework changes the token/session/session-store model.
- A future protocol decision moves proof carrier behavior into the Protobuf envelope, first system message, or WebSocket handshake.

## Follow-Up

- Add repository checks for selected login/token boundaries in W-0072.
- Close M-013 after checks pass.
- Start runtime authentication implementation only after a future implementation milestone explicitly authorizes schema ratification, migrations, repositories, adapters, tests, and runtime wiring.
