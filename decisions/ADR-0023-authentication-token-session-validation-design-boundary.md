# ADR-0023: Authentication Token Session Validation Design Boundary

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-confirm-next-direction-after-player-account-postgresql-persistence/`
- `changes/2026-05-14-define-authentication-token-session-validation-design-standard/`
- `changes/2026-05-14-ratify-authentication-proof-token-session-contract-dimensions/`

Related conversations:

- `conversations/2026-05-14-authentication-token-session-validation-direction.md`
- `conversations/2026-05-14-authentication-token-session-validation-design-standard.md`

Related artifacts:

- `docs/authentication-token-session-validation.md`
- `docs/authentication-token-session-validation.zh-CN.md`
- `docs/authentication-proof-token-session-contract-dimensions.md`
- `docs/authentication-proof-token-session-contract-dimensions.zh-CN.md`
- `docs/player-identity-session-boundary.md`
- `docs/player-account-session-contracts.md`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `AGENTS.md`
- `runtime/AGENTS.md`

## Context

The repository now has player account lifecycle contracts, Protobuf wire messages, PostgreSQL lifecycle schema, migration source, repository interface, and PostgreSQL adapter implementation. It also has an application-owned session validation hook, but the current validator is metadata-only and does not authenticate clients.

The maintainer selected the next direction as authentication and token/session validation design before runtime player handlers or WebSocket routes.

Without a design boundary, a future agent could implement login by placing token parsing in WebSocket transport, storing credentials in player account tables, treating envelope metadata as proof, or mixing Nakama/Pitaya concepts directly into vibit API shape.

## Decision

Ratify an authentication, token, and session validation design standard before implementation.

The design separates:

- Authentication proof.
- Login methods.
- Tokens.
- Credentials.
- External identity links.
- Runtime sessions.
- Request identity.
- WebSocket handshake authentication.
- Player account lifecycle.
- Transport connection metadata.
- Protobuf envelope metadata.

The current metadata-only validator remains a non-authenticated bootstrap path. Metadata-only `player_id`, `session_id`, and `connection_id` must not satisfy production permissions.

Nakama remains the reference for account/auth/session token/refresh/realtime socket capability coverage. Pitaya remains the reference for session binding, handler context, frontend/backend, and realtime session vocabulary. vibit adapts those concepts into its own agent-native ownership, contract, and verification model.

This decision does not choose or implement a login method, token format, refresh behavior, signing, expiration, revocation, credential storage, external identity linking, session persistence, Protobuf envelope change, WebSocket handshake authentication behavior, runtime player handler, or WebSocket route.

## Alternatives Considered

- Implement runtime player account handlers before authentication design.
- Choose JWT or opaque tokens immediately.
- Add guest/device login as a convenient first authentication path.
- Store credentials, token state, or external identity subjects in `player_accounts`.
- Put token parsing in the WebSocket transport or Protobuf adapter.
- Copy Nakama or Pitaya public API shape.
- Treat envelope `player_id` and `session_id` as sufficient for production permissions.

## Rationale

Authentication decisions are security-sensitive and architecture-sensitive. Login methods affect public contracts, persistence, dependencies, migration shape, error behavior, replay handling, and operations.

Separating design dimensions first lets agents continue making progress without accidentally committing vibit to a token format, credential table, provider dependency, or transport handshake model.

The separation also keeps vibit's main differentiator intact: future authentication work should be explicit, contract-first, manifest-visible, generated where appropriate, and checkable by narrow repository tooling.

## Agent Reasoning Summary

The safe next step after player account persistence is not to implement login. It is to define which concepts must stay separate, map mature game server references into vibit-native terms, and create a bounded queue for future checks and contracts.

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

- Future authentication work must read `docs/authentication-token-session-validation.md`.
- Metadata-only identity remains unauthenticated until a real validator is separately ratified and implemented.
- Player account lifecycle persistence remains free of credentials, tokens, external identity links, session rows, WebSocket state, and request validation results.
- WebSocket transport and Protobuf adapters remain credential-neutral until a future protocol or handshake decision changes their role.
- M-011 proceeds through bounded design, check, and contract steps before implementation.
- Any concrete token, login, credential, external identity, session persistence, Protobuf envelope, or WebSocket handshake choice remains an ask-first boundary.

## Reversal Conditions

Revisit this decision if:

- A future first authentication implementation proves that the boundary prevents a simple, verifiable production path.
- A compatibility goal with Nakama, Pitaya, or another framework is explicitly ratified.
- The current Protobuf envelope cannot support any reasonable validated session flow without a compatibility-sensitive change.
- Distributed runtime routing requires session ownership to move to a different platform boundary.

## Follow-Up

- Add architecture checks enforcing the authentication/token/session design boundary.
- Ratify semantic contract dimensions for authentication proof and token/session validation without choosing a concrete token format. Completed by `changes/2026-05-14-ratify-authentication-proof-token-session-contract-dimensions/`.
- Define credential storage and external identity linking boundaries.
- Define session persistence and WebSocket handshake decision gates.
- Create the next implementation confirmation gate after M-011 is complete.
