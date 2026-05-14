# ADR-0021: Player Identity And Session Boundary

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-player-identity-and-session-boundary/`

Related conversations:

- Maintainer continuation after `M-002 Durable Inventory Runtime` completion.

Related artifacts:

- `docs/player-identity-session-boundary.md`
- `docs/player-identity-session-boundary.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/protocol.yaml`
- `.arch/modules.yaml`
- `.arch/work-items.yaml`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`
- `docs/reference-game-server-alignment.md`
- `decisions/ADR-0015-game-protocol-framework.md`
- `decisions/ADR-0019-nakama-and-pitaya-reference-baseline.md`

## Context

The durable inventory runtime now has WebSocket, Protobuf, application dispatch, inventory handlers, PostgreSQL persistence, migration tooling, and live PostgreSQL request-loop verification.

The next architectural risk is player identity drift. Inventory already accepts `player_id`, and the Protobuf envelope already has session metadata fields, but vibit has not yet defined which layer owns player accounts, authentication, runtime sessions, WebSocket connection identity, or application request identity context.

Nakama treats accounts, authentication, users, and sessions as first-class game backend concerns. Pitaya treats acceptors, sessions, routes, and handlers as separate server framework concerns. vibit should align with those separations while preserving its own agent-native boundaries.

## Decision

Define player identity, player accounts, authentication, runtime sessions, transport connections, and request identity context as separate concerns.

The future `player` module owns player identity and player account lifecycle.

The WebSocket transport owns only transport-local connection mechanics and metadata. It must not authenticate players, parse credentials, own player accounts, or enforce domain permissions.

The Protobuf protocol adapter decodes and preserves existing envelope session metadata, but it does not own long-lived sessions or choose authentication schemes.

Application dispatch owns the future session-validation handoff and the application-facing request identity context passed to domain handlers.

Inventory may reference `player_id` for player-owned inventory aggregates, but it does not own player accounts, authentication, session validation, or token formats.

This decision does not choose a concrete authentication scheme, token format, credential store, player account database schema, WebSocket handshake contract, or Protobuf envelope change.

## Alternatives Considered

- Let inventory continue treating `player_id` as trusted request data.
- Put authentication and session binding in the WebSocket transport adapter.
- Put token parsing and player account lookup in the Protobuf bridge.
- Make the player module own inventory state.
- Immediately copy Nakama account/session APIs.
- Immediately copy Pitaya session API shape.
- Add JWT, OAuth/OIDC, guest login, or password-based accounts now.

## Rationale

Trusting client-supplied `player_id` would let the first runtime keep moving, but it would create a weak foundation for permissions, account lifecycle, player-bound reads, service actors, and future multiplayer behavior.

Putting identity behavior in transport or protocol code would make the implementation locally convenient but would weaken vibit's core architecture: transport should handle connections, protocol should decode envelopes, application should orchestrate request context, and domain modules should enforce business rules.

Copying Nakama or Pitaya APIs now would prematurely lock vibit into an external surface. Their concepts are useful references, but vibit's differentiator is explicit, generated, machine-checkable agent-native architecture.

Deferring authentication scheme and credential decisions keeps the project from accidentally accepting long-term security, dependency, and compatibility commitments during a boundary-definition step.

## Agent Reasoning Summary

The next safe step after durable inventory is not to implement login. It is to make the identity/session boundary explicit enough that future login, permissions, and player-account work can be added without rewriting transport, protocol, application dispatch, or inventory ownership.

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

- Future player/session work must read `docs/player-identity-session-boundary.md`.
- Inventory remains the owner of inventory state and must not take over player account behavior.
- The player module is the planned owner of player identity and account lifecycle.
- WebSocket transport and Protobuf adapters remain identity-neutral until explicit validation handoff contracts exist.
- The existing envelope session fields remain unchanged.
- Adding authentication providers, token formats, credential storage, player account migrations, or WebSocket handshake authentication requires separate ratification.

## Reversal Conditions

Revisit this decision if:

- The first real authentication implementation proves the player/session split is too heavy.
- A future compatibility goal requires a Nakama-like or Pitaya-like public API surface.
- Distributed runtime routing requires session ownership to move to a different platform boundary.
- The current Protobuf envelope cannot support session validation without a compatibility-sensitive change.

## Follow-Up

- Add a player module manifest and module agent guides without account migrations.
- Define application-owned request identity/session handoff types.
- Define session validation result vocabulary before implementing authentication.
- Replace inventory bootstrap permission policy with session-context-aware permission behavior.
- Add player account persistence only after account schema and authentication model are ratified.
