# ADR-0024: Login Method Token Format Ratification Boundary

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-confirm-next-direction-after-authentication-design/`
- `changes/2026-05-14-define-login-method-token-format-ratification-standard/`

Related conversations:

- `conversations/2026-05-14-login-method-token-format-direction.md`
- `conversations/2026-05-14-login-method-token-format-ratification-standard.md`

Related artifacts:

- `docs/login-method-token-format-ratification.md`
- `docs/login-method-token-format-ratification.zh-CN.md`
- `docs/authentication-token-session-validation.md`
- `docs/credential-storage-external-identity-linking-boundaries.md`
- `docs/session-persistence-websocket-handshake-decision-gates.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `AGENTS.md`
- `runtime/AGENTS.md`

## Context

The repository has completed the authentication and token/session validation design milestone, but it has intentionally not selected a concrete login method, token format, proof carrier, credential schema, session store, Protobuf envelope change, WebSocket handshake behavior, runtime player handler, or WebSocket route.

The maintainer authorized the agent's recommendation to make login-method and token-format ratification the next major direction. This is the correct moment to convert the broad design boundary into a narrower decision process before implementation agents start writing login code.

Without this boundary, a future agent could choose a convenient login method, place tokens in existing metadata fields, add JWT or password hashing dependencies prematurely, create credential tables inside player account lifecycle storage, or copy Nakama/Pitaya API shape without an explicit compatibility decision.

## Decision

Create a dedicated login-method and token-format ratification standard before implementation.

The standard requires future work to compare:

- First login-method candidate families.
- Token kinds and token formats.
- Proof carrier options.
- Token lifecycle semantics.
- Storage, migration, dependency, error, permission, test, and repository-check implications.
- Nakama capability coverage and Pitaya session vocabulary in vibit-native terms.

The standard also establishes the M-013 work queue:

1. Define the ratification standard.
2. Compare first login-method candidates.
3. Ratify the first login-method set.
4. Compare token format and token carrier options.
5. Ratify the first token format and proof carrier posture.
6. Define token lifecycle and storage implications.
7. Define authentication contract, error, and permission surfaces.
8. Define credential, token, and session schema gates.
9. Add repository checks for selected login/token boundaries.
10. Close the milestone and create the next implementation gate.

This decision does not implement authentication, select a final token parser, add credential storage, add external identity linking, add session persistence, change the Protobuf envelope, change WebSocket handshake authentication, add runtime player handlers, or add WebSocket routes.

## Alternatives Considered

- Implement device or guest login immediately.
- Choose JWT immediately because it is common.
- Choose opaque tokens immediately because they are operationally simple.
- Add email/password first because it is familiar.
- Add WebSocket handshake authentication first to reject unauthenticated connections early.
- Reuse existing Protobuf `Session.session_id` as proof.
- Let each future implementation work item decide its own authentication shape.
- Copy Nakama's account/auth/session public API directly.

## Rationale

Login methods and token formats are high-impact choices. They affect public contracts, persistence, secret handling, dependencies, replay controls, operations, client ergonomics, and future compatibility.

The first implementation should be small, but the decision process must still be production-minded. A small first slice is useful only if agents can see why the selected method is safe enough, what it intentionally does not solve, and what artifacts must exist before code.

Nakama is useful for capability coverage because it shows common game backend authentication and session behavior. Pitaya is useful for session binding and handler/session vocabulary. vibit should learn from both without importing their API shape or cluster assumptions.

## Agent Reasoning Summary

The next safe step is to define how login and token choices will be made, then advance through bounded comparison and ratification work items. This keeps the project moving while preserving agent-native maintainability and preventing shortcut security implementations.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: low
  reversibility: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- Future login/token work must read `docs/login-method-token-format-ratification.md`.
- M-013 can advance without asking the maintainer for every technical comparison, but implementation boundaries remain protected.
- Metadata-only `player_id`, `session_id`, and `connection_id` remain unauthenticated.
- Token carrier behavior remains separate from current Protobuf metadata fields and WebSocket handshake behavior.
- Credential storage, external identity linking, session persistence, and major authentication dependencies remain deferred until separately ratified.
- Future agents must leave ratification packets when selecting login methods, token models, carriers, or lifecycle rules.

## Reversal Conditions

Revisit this decision if:

- The first authentication implementation proves the ratification packet is too heavy for a small production slice.
- A compatibility goal with Nakama, Pitaya, or another framework is explicitly ratified.
- A future security review requires selecting a token format before login-method selection.
- A future protocol decision changes the envelope or handshake model enough that token carrier choices must be reordered.

## Follow-Up

- Compare first login-method candidates.
- Ratify the first login-method set.
- Compare token format and token carrier options.
- Ratify the first token format and carrier posture.
- Define lifecycle, storage, contract, schema, check, and implementation gates before code.
