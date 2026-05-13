# ADR-0019: Nakama And Pitaya Reference Baseline

Status: Accepted
Date: 2026-05-13
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-13-align-with-nakama-and-pitaya/`

Related conversations:

- `conversations/2026-05-13-nakama-and-pitaya-reference-alignment.md`

Related artifacts:

- `docs/reference-game-server-alignment.md`
- `docs/reference-game-server-alignment.zh-CN.md`
- `.arch/reference.yaml`
- `README.md`
- `AGENTS.md`

## Context

The project has started its Go runtime foundation with generated Protobuf output, protocol handoff types, and an application dispatch skeleton.

The maintainer clarified that vibit should actively and repeatedly reference Nakama and Pitaya. The intended game server functionality is broadly the same problem class; vibit's main difference is that it is Agent-Native.

This needs a durable decision so future agents do not treat the first inventory slice as the product boundary, and also do not blindly copy external APIs or rush into distributed runtime work.

## Decision

Use Nakama and Pitaya as active reference baselines for game server capability planning.

Nakama is the primary reference for broad game backend product capability surface:

- Accounts, authentication, users, sessions.
- Storage, social systems, chat, groups, parties.
- Matchmaker, realtime multiplayer, authoritative matches.
- Leaderboards, tournaments, economy-style and LiveOps capability families.
- Dashboard, metrics, and operations.

Pitaya is the primary reference for Go game server framework architecture vocabulary:

- Client acceptors such as WebSocket and TCP.
- Sessions and route handlers.
- Server-to-server RPC.
- Frontend/backend server roles.
- Groups, rooms, broadcast, serializers, and cluster service discovery vocabulary.

These projects are references, not governing standards. vibit's constitution, ADRs, manifests, generated boundaries, and verification commands remain authoritative.

## Alternatives Considered

- Continue only with the inventory proof slice and defer product capability alignment.
- Pick Nakama as the sole reference.
- Pick Pitaya as the sole reference.
- Attempt API compatibility with one of the references now.
- Start with distributed Pitaya-style runtime before proving single-process module boundaries.

## Rationale

Nakama and Pitaya cover different useful parts of the same game server problem space.

Nakama helps prevent vibit from becoming too narrow by reminding agents of the complete backend capability surface that real games need. Pitaya helps prevent vibit from inventing Go server framework vocabulary without learning from a mature route/session/RPC/group architecture.

Using both as references gives vibit a stronger product compass while preserving its real differentiator: making the server architecture unusually legible, bounded, generated, and verifiable for coding agents.

## Agent Reasoning Summary

The correct response to the maintainer's clarification is to add a planning and architecture standard, not to immediately copy APIs or add runtime dependencies. The standard should shape future module choices and implementation order while keeping the current proof slice focused and controllable.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: low
  reversibility: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- Future agents should consult `docs/reference-game-server-alignment.md` and `.arch/reference.yaml` before adding new game server modules or runtime subsystems.
- The roadmap should cover Nakama-like broad backend capabilities and Pitaya-like server framework architecture concerns over time.
- Feature parity must not override vibit's contract-first, manifest-first, generated, and checked architecture.
- Pitaya-style distributed runtime remains deferred until the modular monolith proof slice is stable.
- Any future external API compatibility decision requires a separate ADR.

## Reversal Conditions

Revisit this decision if vibit intentionally pivots away from game servers toward generic backend-only scope, if a different reference framework becomes clearly more relevant, or if compatibility with a specific external API becomes a ratified product goal.

## Follow-Up

- Use the reference capability matrix before proposing new modules.
- Keep the near-term runtime sequence focused on inventory handler, repository, policy, and persistence boundaries.
- Revisit player/session/auth as the next major capability family after the first vertical inventory slice is stable.
- Defer cluster, RPC, groups, and distributed routing until single-process contracts and verification are proven.
